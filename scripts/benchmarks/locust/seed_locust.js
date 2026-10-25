// Seeds the on-chain + off-chain state needed by the Locust load test:
// registered users, authenticated sessions (genuine zkauth Groth16 proofs),
// the off-chain identity index, records (upload + key + ZK grant), and
// optional emergency recovery escrow — then dumps a CSV the locustfile eats.
//
// Assumes a fresh, reset stack (chain + postgres) with contracts deployed and
// the backend running — exactly the baseline the benchmark runner leaves
// behind. Run via `scripts/benchmarks/locust/run_locust_load.sh` or:
//
//   source ./.env && npx hardhat run scripts/benchmarks/locust/seed_locust.js --network localhost
//
// Knobs (env):
//   LOCUST_PATIENTS            default 25
//   LOCUST_DOCTORS             default 5
//   LOCUST_RECORDS_PER_PATIENT default 2
//   LOCUST_RECOVERY            default 1 (0 disables recovery-key escrow)
//   LOCUST_CSV                 default scripts/benchmarks/results/locust/seed.csv
//   API_BASE_URL               point at an API relay for latency-injected setups
//   SESSION_CONCURRENCY        default 1 (serial). >1 mints the zkauth sessions
//                               concurrently: proofs are precomputed off the
//                               timing clock, then the authenticate txs are
//                               submitted as a burst. Also drives the listener
//                               ingestion measurement.
//   INGEST_POLL_MS             default 200 (listener ingestion poll interval)
//   INGEST_TIMEOUT_MS          default 60000 (give up waiting for ingestion)
//   SESSION_INGEST_JSON        path for the burst/ingestion report
//   DELEGATION_INGEST_JSON     path for the per-record delegation lag report
//                               (grant tx confirmed -> reencrypted_keys row)
//
// Records are seeded on the proxy re-encryption path: the patient seals the
// content key into a PRE lockbox (C1), installs it, and delegates rk_{patient->
// doctor}. The listener's AccessGranted handler then re-encrypts it for the
// doctor. See docs/pre-design.md and docs/pre-evaluation.md.
import "dotenv/config";
import hre from "hardhat";
import crypto from "crypto";
import fs from "fs";
import os from "os";
import path from "path";
import {
  readAddrs, buildPoseidon, zkauthInputs, realZkAuthSession,
  apiGet, apiPost, apiUpload, provisionIdentity, toHex, pollDb, proveFor,
} from "./../lib.js";
import { precliSetup } from "../../lib/precli.js";

const { ethers } = hre;
const rndHex = (n) => "0x" + crypto.randomBytes(n).toString("hex");
const randBig = () => BigInt("0x" + crypto.randomBytes(31).toString("hex")) + 1n;
const nowSec = () => Math.floor(Date.now() / 1000);

const PATIENTS = Number(process.env.LOCUST_PATIENTS || 25);
const DOCTORS = Number(process.env.LOCUST_DOCTORS || 5);
const RECORDS_PER = Number(process.env.LOCUST_RECORDS_PER_PATIENT || 2);
const RECOVERY = (process.env.LOCUST_RECOVERY || "1") !== "0";
const CSV_PATH = process.env.LOCUST_CSV || path.join(__dirname, "..", "results", "locust", "seed.csv");

// Session-establishment burst knobs. Default 1 keeps the historical serial
// seeding, so the existing benchmark runs stay reproducible byte-for-byte.
const SESSION_CONCURRENCY = Math.max(1, Number(process.env.SESSION_CONCURRENCY || 1));
const BURST = SESSION_CONCURRENCY > 1;
const INGEST_POLL_MS = Number(process.env.INGEST_POLL_MS || 200);
const INGEST_TIMEOUT_MS = Number(process.env.INGEST_TIMEOUT_MS || 60000);
const INGEST_JSON = process.env.SESSION_INGEST_JSON || null;
const DELEGATION_INGEST_JSON = process.env.DELEGATION_INGEST_JSON || null;

const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

function pct(values, p) {
  if (!values.length) return null;
  const s = [...values].sort((a, b) => a - b);
  const i = (s.length - 1) * (p / 100);
  const lo = Math.floor(i);
  const hi = Math.ceil(i);
  return lo === hi ? s[lo] : s[lo] + (s[hi] - s[lo]) * (i - lo);
}

function fmt(v) {
  return v == null ? "n/a" : v.toFixed(1);
}

async function parseEvent(receipt, iface, name) {
  const log = receipt.logs
    .map((l) => { try { return iface.parseLog({ topics: [...l.topics], data: l.data }); } catch { return null; } })
    .find((e) => e?.name === name);
  return log;
}

// Explicit per-account nonces: chain is fresh, account sends are strictly
// sequential, so pinning is deterministic and immune to slow nonce queries.
// Reciprocated against the node's pending count on every call (guards against
// txpool/pending-count lag under bursty seeding, e.g. the 50-patient soak).
function makeNoncer(provider, addr) {
  let next = null;
  return async () => {
    const pending = await provider.getTransactionCount(addr, "pending");
    next = next == null ? pending : Math.max(next, pending);
    const nonce = next;
    next += 1;
    return { nonce };
  };
}

// Retry well-known txpool races (nonce / replacement / "already known") with a
// freshly reconciled nonce each attempt; hardhat occasionally rejects fast
// bursts with "replacement transaction underpriced" even with correct nonces.
const POOL_ERR = /replacement transaction underpriced|replacement transaction fee too low|nonce too low|nonce too high|already known|already imported|transaction underpriced|replacement transaction cannot be set/i;
async function sendAndWait(makeTx, attempts = 6) {
  let lastErr;
  for (let a = 0; a < attempts; a++) {
    try {
      const tx = await makeTx();
      return await tx.wait();
    } catch (err) {
      lastErr = err;
      if (!POOL_ERR.test(String((err && err.message) || err))) throw err;
      await new Promise((r) => setTimeout(r, 300 * (a + 1)));
    }
  }
  throw lastErr;
}

// Is this session actually usable at the API yet?
//
// GET /api/v1/attributes/registry/{kind}/{name} runs requireAuth first and has
// no role gate and no side effects, so 401 is produced *only* by the session
// store missing the token. Any other status (404 here, because the random name
// does not resolve) therefore proves the listener has already persisted the
// session. That gives a black-box ingestion probe needing no backend change.
async function sessionUsable(user, token) {
  const res = await apiGet(`/api/v1/attributes/registry/nid/${rndHex(8)}`, { token, user });
  return res.status !== 401;
}

// Mints one real zkauth session per signer.
//
// Serial path (SESSION_CONCURRENCY=1) is the historical behaviour: prove, send,
// wait for the receipt, one user at a time.
//
// Burst path (>1) separates the two costs that dominate, because they have very
// different bottlenecks:
//   1. proof generation — pure client-side CPU, off the timing clock entirely;
//   2. transaction submission — bounded by the 2s Clique block floor and ~312k
//      gas of pairing verification per session.
// Sending them concurrently is what "session establishment under load" means,
// and it is the only way to see how many sessions the chain absorbs per block.
async function mintSessions(entries, zkAuth, poseidon, provider) {
  const t0 = Date.now();
  if (!BURST) {
    const out = [];
    for (const e of entries) {
      const s = await realZkAuthSession(zkAuth, e.signer, poseidon, {
        secret: randBig(), nullifierSeed: 0n, actionId: 0n,
      });
      out.push({ user: e.user, token: s.token, confirmAt: Date.now() });
    }
    return { sessions: out, proofMs: null, submitMs: Date.now() - t0 };
  }

  const prepared = [];
  for (const e of entries) {
    const secret = randBig();
    const { commitment, nullifier, inputs } = zkauthInputs(poseidon, { secret, nullifierSeed: 0n, actionId: 0n });
    const proof = await proveFor("zkauth", inputs);
    prepared.push({ ...e, commitment, nullifier, proof });
  }
  const proofMs = Date.now() - t0;
  console.log(`  precomputed ${prepared.length} real zkauth proofs in ${proofMs}ms (off the timing clock)`);

  // One funded wallet per signer, so there is no shared-nonce contention, but
  // pin each account's nonce anyway: geth's pending count can lag behind a wide
  // burst and hand two txs the same nonce.
  const noncers = new Map();
  for (const p of prepared) noncers.set(p.user, makeNoncer(provider, p.signer.address));

  const submitT0 = Date.now();
  const submitted = await Promise.all(prepared.map(async (p) => {
    const makeTx = () => zkAuth.connect(p.signer).authenticate(
      p.signer.address, toHex(p.commitment),
      p.proof.formatted.pA, p.proof.formatted.pB, p.proof.formatted.pC,
      toHex(p.nullifier), toHex(0n), noncers.get(p.user)()
    );
    const receipt = await sendAndWait(makeTx);
    const log = await parseEvent(receipt, zkAuth.interface, "Authenticated");
    if (!log) throw new Error(`genuine authenticate failed for ${p.user}`);
    return { user: p.user, token: log.args.sessionToken, confirmAt: Date.now(), block: receipt.blockNumber };
  }));
  const submitMs = Date.now() - submitT0;

  const blocks = submitted.map((s) => s.block);
  const span = blocks.length ? Math.max(...blocks) - Math.min(...blocks) + 1 : 0;
  console.log(`  burst of ${submitted.length} authenticate txs confirmed in ${submitMs}ms across ${span} block(s)`
    + ` (${(submitted.length / (submitMs / 1000)).toFixed(2)} txs/s, block floor ~2000ms)`);

  return { sessions: submitted, proofMs, submitMs, blocks: span };
}

// Polls every minted session until the listener has persisted it, recording
// per-session lag measured from that session's own block confirmation. This is
// the metric the old fixed `sleep(3000)` assumed but never measured.
async function measureIngestion(sessions) {
  const pending = new Set(sessions.map((s) => s.user));
  const lags = new Map();
  const start = Date.now();

  while (pending.size && Date.now() - start < INGEST_TIMEOUT_MS) {
    const batch = sessions.filter((s) => pending.has(s.user));
    const checks = await Promise.all(batch.map(async (s) => [s, await sessionUsable(s.user, s.token)]));
    for (const [s, ok] of checks) {
      if (!ok) continue;
      pending.delete(s.user);
      lags.set(s.user, Math.max(0, Date.now() - s.confirmAt));
    }
    if (pending.size) await sleep(INGEST_POLL_MS);
  }
  return { lags: [...lags.values()], missing: pending.size, wallMs: Date.now() - start };
}

function reportIngestion(stats, label) {
  const { lags, missing, wallMs } = stats;
  if (!lags.length) {
    console.log(`  listener ingestion: NONE of the sessions became usable within ${wallMs}ms (${missing} missing)`);
    return;
  }
  console.log(`  listener ingestion over ${lags.length} sessions (${missing} never appeared):`);
  console.log(`    min=${fmt(pct(lags, 0))}ms  p50=${fmt(pct(lags, 50))}ms  p95=${fmt(pct(lags, 95))}ms`
    + `  p99=${fmt(pct(lags, 99))}ms  max=${fmt(Math.max(...lags))}ms`);
  if (label) console.log(`  ${label}`);
}

async function main() {
  const [deployer] = await ethers.getSigners();
  const provider = ethers.provider;
  const dNonce = makeNoncer(provider, deployer.address);

  const addr = readAddrs();
  const registry = await ethers.getContractAt("WellRegistry", addr.registry);
  const historyFactory = await ethers.getContractAt("HistoryFactory", addr.historyFactory);
  const zkAuth = await ethers.getContractAt("ZKAuthVerifier", addr.zkAuth);
  const poseidonLibs = { PoseidonT3: addr.poseidonT3 };
  const History = await ethers.getContractFactory("History", { libraries: poseidonLibs });
  const Record = await ethers.getContractFactory("Record", { libraries: poseidonLibs });

  const poseidon = await buildPoseidon();

  console.log("=== seed: register dev admin (deployer) if needed ===");
  const hasUser = async (a) => { try { await registry.getUserByAddress(a); return true; } catch { return false; } };
  if (!(await hasUser(deployer.address))) {
    await sendAndWait(async () => registry.connect(deployer).register(rndHex(32), rndHex(32), 3, { nonce: (await dNonce(deployer)).nonce }));
  }

  console.log(`=== seed: create ${PATIENTS} patients + ${DOCTORS} doctors (fund + register) ===`);
  const wallets = async (n, isDoctor, registerFn) => {
    const out = [];
    for (let i = 0; i < n; i++) {
      const w = (await ethers.Wallet.createRandom(provider)).connect(provider);
      await sendAndWait(async () => deployer.sendTransaction({ to: w.address, value: ethers.parseEther("2"), nonce: (await dNonce(deployer)).nonce }));
      await registerFn(w);
      if (isDoctor) w.n = makeNoncer(provider, w.address);
      out.push(w);
    }
    return out;
  };
  const patients = await wallets(PATIENTS, false, async (w) => {
    const secret = randBig();
    const { commitment } = zkauthInputs(poseidon, { secret, nullifierSeed: 0n, actionId: 0n });
    const n = makeNoncer(provider, w.address);
    await sendAndWait(async () => registry.connect(w).register(toHex(commitment), rndHex(32), 0, await n(w)));
    w.commitment = toHex(commitment);
  });
  const doctors = await wallets(DOCTORS, true, async (w) => {
    await sendAndWait(async () => registry.connect(deployer).registerByAdmin(w.address, rndHex(32), rndHex(32), 1, await dNonce(deployer)));
  });

  console.log(`=== seed: mint real zkauth sessions (genuine Groth16 proofs) for ${PATIENTS + DOCTORS + 1} users ===`);
  if (BURST) {
    console.log(`  SESSION_CONCURRENCY=${SESSION_CONCURRENCY} (burst submission + measured listener ingestion)`);
  }
  // The deployer/admin session is minted serially and off the critical path:
  // it is needed for identity provisioning, and keeping it separate keeps the
  // burst measurement purely about the PATIENTS + DOCTORS population.
  const adminSession = (await realZkAuthSession(zkAuth, deployer, poseidon, {
    secret: randBig(), nullifierSeed: 0n, actionId: 0n,
  })).token;

  const burstEntries = [
    ...patients.map((p) => ({ user: p.address, signer: p })),
    ...doctors.map((d) => ({ user: d.address, signer: d })),
  ];
  const minted = await mintSessions(burstEntries, zkAuth, poseidon, provider);
  const byUser = new Map(minted.sessions.map((s) => [s.user, s.token]));
  const patientSessions = patients.map((p) => byUser.get(p.address));
  const doctorSessions = doctors.map((d) => byUser.get(d.address));

  // Replace the old fixed sleep with an actual measurement of how long the
  // listener takes to make each freshly-confirmed session usable.
  const ingest = await measureIngestion(minted.sessions);
  reportIngestion(ingest, BURST
    ? "  (lag is per-session, from its own block confirmation to first non-401 API response)"
    : null);
  if (ingest.missing > 0) {
    throw new Error(`${ingest.missing} session(s) never became usable at the API within ${INGEST_TIMEOUT_MS}ms`);
  }

  if (BURST) {
    const report = {
      concurrency: SESSION_CONCURRENCY,
      sessions: minted.sessions.length,
      proofMs: minted.proofMs,
      submitMs: minted.submitMs,
      blocksSpanned: minted.blocks ?? null,
      txsPerSec: minted.sessions.length / (minted.submitMs / 1000),
      ingestPollMs: INGEST_POLL_MS,
      ingestWallMs: ingest.wallMs,
      ingestMissing: ingest.missing,
      ingestLagMs: {
        min: pct(ingest.lags, 0), p50: pct(ingest.lags, 50), p95: pct(ingest.lags, 95),
        p99: pct(ingest.lags, 99), max: Math.max(...ingest.lags),
      },
      // The full sorted sample, not just percentiles: the distribution is
      // bimodal (a fast cluster plus a cluster near the block period), and
      // percentiles alone hide that.
      ingestLagSortedMs: [...ingest.lags].sort((a, b) => a - b),
    };
    const out = INGEST_JSON
      || path.join(path.dirname(CSV_PATH), "session_ingest.json");
    fs.mkdirSync(path.dirname(out), { recursive: true });
    fs.writeFileSync(out, JSON.stringify(report, null, 2) + "\n");
    console.log(`  wrote ${out}`);
    // Compare like with like: the burst is chain-side wall time, the lag is
    // per-session chain-to-API. ingestWallMs is only how long the poll loop
    // itself took, and is smaller than the max lag whenever early sessions
    // already satisfied their probe on the first pass.
    console.log(`  chain burst ${fmt(minted.submitMs)}ms vs listener ingestion`
      + ` p50=${fmt(pct(ingest.lags, 50))}ms p95=${fmt(pct(ingest.lags, 95))}ms`
      + ` — the ingestion tail, not its median, is what bounds how fast a`
      + ` freshly established session is usable`);
  }

  console.log("=== seed: off-chain identity index (admin provision) ===");
  const summaries = [];
  for (const [i, p] of patients.entries()) {
    const summary = await registry.getEmergencySummaryByAddress(p.address);
    if (summary === ethers.ZeroAddress) throw new Error(`no EmergencySummary for patient ${i}`);
    summaries.push(summary);
    await provisionIdentity(adminSession, deployer.address, {
      national_id: `LOCUST-P-${i}`, display_name: `Load Patient ${i}`, role: "patient",
      address: p.address, summary_address: summary, commitment: p.commitment, internal_id: rndHex(32),
    });
  }
  for (const [i, d] of doctors.entries()) {
    await provisionIdentity(adminSession, deployer.address, {
      national_id: `LOCUST-D-${i}`, display_name: `Load Doctor ${i}`, role: "doctor",
      address: d.address, summary_address: "", commitment: "", internal_id: rndHex(32),
    });
  }

  // per-patient History contract + per-record history/record instances
  const historyFor = async (patient) => {
    const evs = await historyFactory.queryFilter(historyFactory.filters.HistoryCreated());
    const ev = evs.find((e) => e.args.patient.toLowerCase() === patient.address.toLowerCase());
    if (!ev) throw new Error(`no HistoryCreated for ${patient.address}`);
    return History.attach(ev.args.historyContract);
  };
  const histories = [];
  for (const p of patients) histories.push(await historyFor(p));

  console.log(`=== seed: records (upload/lockbox/delegate/create/setPolicy/ZK grant) — ${PATIENTS * RECORDS_PER} records ===`);
  const patientRows = patients.map((p, i) => ({
    address: p.address, role: "patient", display_name: `Load Patient ${i}`,
    national_id: `LOCUST-P-${i}`, records: [], key_doctor: "", re_encryption_keys: [],
    summary: summaries[i], has_recovery: 0, token: patientSessions[i],
  }));
  const doctorRows = doctors.map((d, i) => ({
    address: d.address, role: "doctor", display_name: `Load Doctor ${i}`,
    national_id: `LOCUST-D-${i}`, records: [], key_doctor: d.address, re_encryption_keys: [],
    summary: "", has_recovery: 0, token: doctorSessions[i],
  }));

  // Delegation ingestion lag: grant tx confirmed -> reencrypted_keys row.
  // Mirrors the session-ingestion measurement above: the PRE transform itself
  // is single-digit milliseconds (see docs/pre-evaluation.md), so what bounds
  // a delegation's time-to-usability is the event-propagation path.
  const delegationLags = [];
  const delegationMissing = [];
  // One client state file for the whole seed: it only carries the secret keys
  // the `self` / `verify` commands read, and seeding never calls them.
  const PRE_STATE = path.join(os.tmpdir(), `well-pre-seed-${process.pid}.json`);

  for (let pi = 0; pi < patients.length; pi++) {
    const patient = patients[pi];
    const history = histories[pi];
    const doctor = doctors[pi % doctors.length];
    const di = doctors.indexOf(doctor);
    const pNonce = makeNoncer(provider, patient.address);
    const dNOnce = doctor.n;
    const doctorTok = doctorSessions[di];

    for (let r = 0; r < RECORDS_PER; r++) {
      const up = await apiUpload(Buffer.from(`LOCUST-EHR-${pi}-${r}-${Date.now()}`), { token: doctorTok, user: doctor.address });
      if (up.status !== 200) throw new Error(`upload failed: ${JSON.stringify(up)}`);
      const { ipfs_cid, content_hash } = up.data;

      const crTx = await sendAndWait(async () => registry.connect(patient).createRecord(
        doctor.address, { recordType: 0, createdAt: 0, updatedAt: 0 }, content_hash, ipfs_cid,
        await pNonce(patient)
      ));
      const crLog = await parseEvent(crTx, history.interface, "RecordCreated");
      if (!crLog) throw new Error(`RecordCreated not emitted for patient ${pi} record ${r}`);
      const recordAddr = crLog.args.record.toLowerCase();
      const record = Record.attach(recordAddr);

      // The patient seals the content key into PRE lockbox C1 and delegates
      // rk_{patient->doctor}. The backend stores both verbatim: it can
      // transform a lockbox but can never open one (docs/pre-design.md §7).
      const pre = precliSetup(PRE_STATE);
      const patientTok = patientSessions[pi];

      const sk = await apiPost(`/api/v1/records/${recordAddr}/key`, {
        record_address: recordAddr, ipfs_cid, content_hash,
        encrypted_key: pre.lockbox, lockbox_salt: pre.lockbox_salt,
        patient_public_key: pre.patient_public_key,
        patient_address: patient.address.toLowerCase(),
      }, { token: patientTok, user: patient.address });
      if (sk.status !== 200) throw new Error(`storeKey(lockbox) failed: ${JSON.stringify(sk)}`);

      const dg = await apiPost(`/api/v1/records/${recordAddr}/delegate`, {
        patient_address: patient.address.toLowerCase(),
        doctor_address: doctor.address.toLowerCase(),
        re_encryption_key: pre.re_encryption_key,
      }, { token: patientTok, user: patient.address });
      if (dg.status !== 200) throw new Error(`delegate failed: ${JSON.stringify(dg)}`);

      patientRows[pi].key_doctor = doctor.address;
      // one rk per record: precli derives a fresh pair per invocation, so the
      // rk is only valid for the lockbox it was derived alongside.
      patientRows[pi].re_encryption_keys.push(pre.re_encryption_key);

      // dynamic per-record WatchAccessGranted sub activates only after the
      // record's creation block; advance the head before granting.
      const recBlock = crTx.blockNumber;
      let head = recBlock;
      while (head <= recBlock) {
        await new Promise((r2) => setTimeout(r2, 200));
        head = await provider.getBlockNumber();
      }

      const policyInputs = {
        department: 0x1111n, institution: 0x2222n,
        requiredDept: 0n, requiredInst: 0n,
        doctorAddr: BigInt(doctor.address),
        validUntil: BigInt(nowSec()) + 365n * 24n * 3600n,
        evaluatedAt: BigInt(nowSec()),
      };
      const proof = await proveFor("policy", policyInputs);

      await sendAndWait(async () => record.connect(patient).setPolicy(
        proof.formatted.publicSignals[0], proof.formatted.publicSignals[1],
        await pNonce(patient)
      ));
      await sendAndWait(async () => record.connect(doctor).grantAccessWithProof(
        proof.formatted.pA, proof.formatted.pB, proof.formatted.pC,
        proof.formatted.publicSignals[0], proof.formatted.publicSignals[1], proof.formatted.publicSignals[2], proof.formatted.publicSignals[3],
        await dNOnce(doctor)
      ));
      const grantConfirmedAt = Date.now();

      const reKey = await pollDb(`SELECT reencrypted_key FROM reencrypted_keys WHERE record_address='${recordAddr}' AND doctor_address='${doctor.address.toLowerCase()}' LIMIT 1`);
      if (reKey == null) {
        delegationMissing.push(recordAddr);
        throw new Error(`reencrypted_keys row missing for ${recordAddr}->${doctor.address}`);
      }
      delegationLags.push(Date.now() - grantConfirmedAt);

      patientRows[pi].records.push(recordAddr);
      doctorRows[di].records.push(recordAddr);
      console.log(`  patient ${pi} record ${r}: ${recordAddr} granted to ${doctor.address.slice(0, 8)}...`);
    }
  }

  // Report the delegation lag distribution, in the same shape as the session
  // ingest report: percentiles plus the full sorted sample, because the
  // distribution inherits block structure and is not described by a mean.
  console.log("=== seed: delegation ingestion (grant confirmed -> C_doctor row) ===");
  if (delegationLags.length) {
    const sorted = [...delegationLags].sort((a, b) => a - b);
    console.log(`  over ${delegationLags.length} delegations (${delegationMissing.length} never appeared):`);
    console.log(`    min=${fmt(pct(sorted, 0))}ms  p50=${fmt(pct(sorted, 50))}ms  p95=${fmt(pct(sorted, 95))}ms`
      + `  p99=${fmt(pct(sorted, 99))}ms  max=${fmt(Math.max(...sorted))}ms`);
    const out = DELEGATION_INGEST_JSON
      || path.join(path.dirname(CSV_PATH), "delegation_ingest.json");
    fs.mkdirSync(path.dirname(out), { recursive: true });
    fs.writeFileSync(out, JSON.stringify({
      delegations: delegationLags.length,
      pollMs: 100,
      lagMs: {
        min: pct(sorted, 0), p50: pct(sorted, 50), p95: pct(sorted, 95),
        p99: pct(sorted, 99), max: Math.max(...sorted),
      },
      missing: delegationMissing,
      lagSortedMs: sorted,
    }, null, 2) + "\n");
    console.log(`  wrote ${out}`);
  }

  console.log("=== seed: emergency recovery escrow (optional) ===");
  if (RECOVERY) {
    for (let pi = 0; pi < patients.length; pi++) {
      const st = await apiPost("/api/v1/emergency/recovery-key",
        { emergency_summary: summaries[pi], content_key: rndHex(32) },
        { token: patientSessions[pi], user: patients[pi].address });
      if (st.status !== 200) throw new Error(`recovery escrow failed for patient ${pi}: ${JSON.stringify(st)}`);
      patientRows[pi].has_recovery = 1;
    }
  }

  console.log("=== seed: write CSV ===");
  fs.mkdirSync(path.dirname(CSV_PATH), { recursive: true });
  // re_encryption_keys is `|`-separated and parallel to `records`, so the
  // locustfile always pairs a record with the rk derived for *its* lockbox.
  const header = "address,role,display_name,national_id,records,key_doctor,summary,has_recovery,token,re_encryption_keys";
  const lines = [header];
  for (const row of [...patientRows, ...doctorRows]) {
    lines.push([row.address, row.role, row.display_name, row.national_id, row.records.join("|"), row.key_doctor, row.summary, row.has_recovery, row.token, row.re_encryption_keys.join("|")].join(","));
  }
  fs.writeFileSync(CSV_PATH, lines.join("\n") + "\n");
  console.log(`  wrote ${CSV_PATH} (${lines.length - 1} users)`);
  console.log("  spot-check: sessions active in backend (listener log shows added sessions)");
}

main().then(() => process.exit(0)).catch((err) => { console.error("seed_locust failed:", err); process.exit(1); });