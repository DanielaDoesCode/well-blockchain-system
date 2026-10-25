import "dotenv/config";
import hre from "hardhat";
import {
  readAddrs, buildPoseidon, zkauthInputs, proveFor, realZkAuthSession, apiPost, apiGet, apiUpload,
  provisionIdentity, stat, printTable, writeJSON, SNAPSHOT, waitForAccessLog, toHex,
} from "./lib.js";

const { ethers } = hre;
import crypto from "crypto";
const rndHex = (n) => "0x" + crypto.randomBytes(n).toString("hex");
const randBig = () => BigInt("0x" + crypto.randomBytes(31).toString("hex")) + 1n;
const nowSec = () => Math.floor(Date.now() / 1000);

const TX_COUNT = Number(process.env.BENCH_TX_COUNT || 20);
const CONC = Number(process.env.BENCH_CONC || 8);
const LOAD_MS = Number(process.env.BENCH_LOAD_MS || 5000);

const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

async function workerLoad(label, fn, { conc, durationMs }) {
  const lat = [];
  let ok = 0, err = 0;
  let stop = false;
  const worker = async () => {
    while (!stop) {
      const t0 = Date.now();
      try {
        const res = await fn();
        lat.push(Date.now() - t0);
        if (res && res.status >= 200 && res.status < 300) ok++; else err++;
      } catch (e) {
        err++;
      }
    }
  };
  const workers = Array.from({ length: conc }, () => worker());
  await sleep(durationMs);
  stop = true;
  await Promise.all(workers);
  const s = stat(lat);
  const rps = (ok + err) / (durationMs / 1000);
  console.log(`${label}  ` +
    `${rps.toFixed(1)} req/s (ok=${ok} err=${err})  ` +
    `latency mean=${s.mean.toFixed(1)}ms p50=${s.p50.toFixed(1)} p95=${s.p95.toFixed(1)} p99=${s.p99.toFixed(1)}`);
  return { label, rps, ok, err, durationMs, latency: { ...s } };
}

async function main() {
  const signers = await ethers.getSigners();
  // reuse the roles minted by bench_latency on the same chain:
  // signers[1]=hospital, signers[2]=responder, signers[4]=doctor, signers[5]=patient
  const [deployer, hospital, responder, , doctor, patient, patient2] = signers;

  const addr = readAddrs();
  const registry = await ethers.getContractAt("WellRegistry", addr.registry);
  const zkAuth = await ethers.getContractAt("ZKAuthVerifier", addr.zkAuth);
  const historyFactory = await ethers.getContractAt("HistoryFactory", addr.historyFactory);
  const poseidonLibs = { PoseidonT3: addr.poseidonT3 };
  const History = await ethers.getContractFactory("History", { libraries: poseidonLibs });
  const Record = await ethers.getContractFactory("Record", { libraries: poseidonLibs });

  const net = await ethers.provider.getNetwork();
  SNAPSHOT.provider = `${net.name}@${net.rpcUrl ?? ethers.provider.connection?.url ?? process.env.ETH_RPC ?? "n/a"}`;
  SNAPSHOT.client = await ethers.provider.send("web3_clientVersion").catch(() => "n/a");

  console.log("=== Setup (idempotent) ===");
  const hasUser = async (userAddr) => {
    try { await registry.getUserByAddress(userAddr); return true; } catch { return false; }
  };
  const safeReg = async (fn) => {
    try { const t = await fn(); await t.wait(); } catch { /* already registered (latency ran on this chain) */ }
  };
  if (!(await hasUser(deployer.address))) await safeReg(() => registry.connect(deployer).register(rndHex(32), rndHex(32), 3));
  for (const [who, role] of [[hospital, 4], [responder, 2], [doctor, 1]]) {
    if (!(await hasUser(who.address))) await safeReg(() => registry.connect(deployer).registerByAdmin(who.address, rndHex(32), rndHex(32), role));
  }
  if (!(await hasUser(patient.address))) await safeReg(() => registry.connect(patient).register(rndHex(32), rndHex(32), 0));
  if (!(await hasUser(patient2.address))) await safeReg(() => registry.connect(patient2).register(rndHex(32), rndHex(32), 0));

  const poseidon = await buildPoseidon();

  const adminTok = (await realZkAuthSession(zkAuth, deployer, poseidon, { secret: randBig(), nullifierSeed: randBig(), actionId: randBig() })).token;
  const responderTok = (await realZkAuthSession(zkAuth, responder, poseidon, { secret: randBig(), nullifierSeed: randBig(), actionId: randBig() })).token;
  const doctorTok = (await realZkAuthSession(zkAuth, doctor, poseidon, { secret: randBig(), nullifierSeed: randBig(), actionId: randBig() })).token;
  const patientTok = (await realZkAuthSession(zkAuth, patient, poseidon, { secret: randBig(), nullifierSeed: randBig(), actionId: randBig() })).token;
  await sleep(3000);

  const patNid = `LOAD-PAT-${crypto.randomBytes(4).toString("hex")}`;
  await provisionIdentity(adminTok, deployer.address, {
    national_id: patNid, display_name: "Load Patient", role: "patient",
    address: patient.address, summary_address: rndHex(20),
    commitment: rndHex(32), internal_id: rndHex(32),
  });

  // single record + grant so the authenticated key-read can be load-tested
  const up = await apiUpload(Buffer.from("LOAD-EHR-" + Date.now()), { token: doctorTok, user: doctor.address });
  if (up.status !== 200) throw new Error("upload: " + JSON.stringify(up.data));
  const { ipfs_cid, content_hash } = up.data;
  let tx;
  tx = await registry.connect(patient).createRecord(
    doctor.address, { recordType: 0, createdAt: 0, updatedAt: 0 }, content_hash, ipfs_cid
  );
  const rc = await tx.wait();
  const hf = await historyFactory.queryFilter(historyFactory.filters.HistoryCreated());
  const histAddr = hf.find((e) => e.args.patient.toLowerCase() === patient.address.toLowerCase()).args.historyContract;
  const history = History.attach(histAddr);
  const rcLogs = rc.logs.map((l) => { try { return history.interface.parseLog({ topics: [...l.topics], data: l.data }); } catch { return null; } });
  const recordAddr = rcLogs.find((e) => e?.name === "RecordCreated").args.record;
  const record = Record.attach(recordAddr);
  await apiPost(`/api/v1/records/${recordAddr.toLowerCase()}/key`, {
    record_address: recordAddr.toLowerCase(), ipfs_cid, content_hash,
    encrypted_key: rndHex(32), patient_address: patient.address.toLowerCase(),
  }, { token: doctorTok, user: doctor.address });
  // wait one block: the listener's per-record WatchAccessGranted subscription is
  // only created from the RecordCreated handler, so a same-block grant is lost.
  let head = rc.blockNumber;
  while (head <= rc.blockNumber) {
    await new Promise((r) => setTimeout(r, 200));
    head = await ethers.provider.getBlockNumber();
  }
  tx = await history.connect(patient).grantRecordAccess(recordAddr, doctor.address);
  await tx.wait();
  if (await waitForAccessLog(tx.hash, "grant", 20000) == null) {
    throw new Error("access_logs row missing for throughput grant");
  }
  await sleep(1000);
  console.log(`  record ${recordAddr} granted to doctor`);

  console.log(`\n=== Chain transaction throughput (geth, N=${TX_COUNT}, real Groth16 proofs) ===`);
  // Precompute one real zkauth proof per tx OFF the timing clock (proving is
  // pure local CPU and the metric is on-chain send+mine); each proof uses a
  // fresh secret/nullifier/actionId so the on-chain usedNullifiers mapping and
  // the Groth16Verifier (real, wired by the benchmark runner) both pass.
  const txPlans = [];
  for (let i = 0; i < TX_COUNT; i++) {
    txPlans.push(zkauthInputs(poseidon, {
      secret: randBig(), nullifierSeed: randBig(), actionId: randBig(),
    }));
  }
  const proofs = await Promise.all(txPlans.map((p) => proveFor("zkauth", p.inputs)));
  console.log(`  precomputed ${proofs.length} real zkauth proofs (off the timing clock)`);

  // Pin per-account nonces: under wide-area RTT the implicit per-send nonce
  // query can race (two same-account txs get the same nonce), and geth rejects
  // the second as "replacement transaction underpriced".
  const nonceMap = {};
  for (const s of signers) nonceMap[s.address.toLowerCase()] = await s.getNonce("pending");
  const txs = [];
  const tSend0 = Date.now();
  // send from the 10 funded geth accounts round-robin
  for (let i = 0; i < TX_COUNT; i++) {
    const s = signers[i % signers.length];
    const key = s.address.toLowerCase();
    const nonce = nonceMap[key];
    nonceMap[key] = nonce + 1;
    const p = proofs[i];
    txs.push(zkAuth.connect(s).authenticate(
      s.address,
      toHex(txPlans[i].commitment),
      p.formatted.pA, p.formatted.pB, p.formatted.pC,
      p.formatted.publicSignals[1], p.formatted.publicSignals[2],
      { nonce }
    ));
  }
  const sent = await Promise.all(txs);
  const sendMs = Date.now() - tSend0;
  const tMine0 = Date.now();
  const receipts = await Promise.all(sent.map((t) => t.wait()));
  const totalMs = Date.now() - tSend0;
  const mineMs = Date.now() - tMine0;
  const firstBlock = Math.min(...receipts.map((r) => r.blockNumber));
  const lastBlock = Math.max(...receipts.map((r) => r.blockNumber));
  const gasSum = receipts.reduce((a, r) => a + r.gasUsed, 0n);
  const chain = {
    txs: TX_COUNT,
    totalMs, sendMs, mineMs,
    tps: TX_COUNT / (totalMs / 1000),
    blocksSpan: lastBlock - firstBlock + 1,
    avgTxsPerBlock: TX_COUNT / (lastBlock - firstBlock + 1),
    avgGasPerTx: (gasSum / BigInt(TX_COUNT)).toString(),
    gasSum: gasSum.toString(),
  };
  console.log(`  ${chain.txs} txs in ${chain.totalMs.toFixed(0)}ms  ->  ${chain.tps.toFixed(2)} txs/s`);
  console.log(`  blocks spanned: ${chain.blocksSpan} (avg ${chain.avgTxsPerBlock} txs/min-block), avg gas/tx ${Number(chain.avgGasPerTx).toLocaleString("en-US")}`);

  console.log(`\n=== API read-path saturation (${CONC} concurrent, ${LOAD_MS}ms) ===`);
  const loadIdentity = await workerLoad("identity lookup (DB)",
    () => apiGet(`/api/v1/identities?national_id=${patNid}`, { token: responderTok, user: responder.address }),
    { conc: CONC, durationMs: LOAD_MS });
  const loadKey = await workerLoad("record key retrieve (DB+transit)",
    () => apiGet(`/api/v1/records/${recordAddr.toLowerCase()}/key?doctor_address=${doctor.address}`, { token: doctorTok, user: doctor.address }),
    { conc: CONC, durationMs: LOAD_MS });

  const json = {
    snapshot: SNAPSHOT,
    chain: chain,
    api: { identityLookup: loadIdentity, recordKey: loadKey },
  };
  const outPath = writeJSON("throughput.json", json);
  console.log(`\nWrote ${outPath}`);
  process.exit(0);
}

main().catch((err) => { console.error("bench_throughput failed:", err); process.exit(1); });