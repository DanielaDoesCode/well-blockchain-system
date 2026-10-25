import "dotenv/config";
import hre from "hardhat";
import {
  API_BASE, readAddrs, buildPoseidon, zkauthInputs, policyInputs, researchInputs,
  proveFor, realZkAuthSession, measureTx, measureApi, apiPost, apiGet, apiUpload,
  provisionIdentity, resolveIdentity, waitForEvent, waitForAccessLog,
  metrics, stat, printTable, writeJSON, SNAPSHOT, toHex, addMetric, zkAuthVerifierDesc,
} from "./lib.js";
import { precliSetup } from "../lib/precli.js";

const { ethers } = hre;

import crypto from "crypto";
import os from "os";
const rndHex = (n) => "0x" + crypto.randomBytes(n).toString("hex");
const randBig = () => BigInt("0x" + crypto.randomBytes(31).toString("hex")) + 1n;
const nowSec = () => Math.floor(Date.now() / 1000);

const AUTH_REPS = Number(process.env.BENCH_AUTH_REPS || 3);
const RECORD_REPS = Number(process.env.BENCH_RECORD_REPS || 3);
// One client state file for the whole run: it only carries the secret keys
// the `self` / `verify` commands read, and this benchmark never calls them.
const PRE_STATE = `${os.tmpdir()}/well-pre-bench-${process.pid}.json`;
let nidSeq = 0;
const NID = (p) => `${p}-BENCH-${++nidSeq}`;

async function parseEvent(receipt, iface, name) {
  const log = receipt.logs
    .map((l) => { try { return iface.parseLog({ topics: [...l.topics], data: l.data }); } catch { return null; } })
    .find((e) => e?.name === name);
  return log;
}

async function expectOk(label, res) {
  if (res.status < 200 || res.status >= 300) {
    throw new Error(`${label} failed (${res.status}): ${JSON.stringify(res.data)}`);
  }
}

async function main() {
  const [deployer, hospital, responder, otherResponder, doctor, patient, patientR] = await ethers.getSigners();

  const addr = readAddrs();
  const registry = await ethers.getContractAt("WellRegistry", addr.registry);
  const historyFactory = await ethers.getContractAt("HistoryFactory", addr.historyFactory);
  const zkAuth = await ethers.getContractAt("ZKAuthVerifier", addr.zkAuth);
  const incident = await ethers.getContractAt("IncidentManager", addr.incident);
  const riskAssess = await ethers.getContractAt("RiskAssessment", addr.riskAssessment);
  const attestation = await ethers.getContractAt("AttestationManager", addr.attestation);
  const contributer = await ethers.getContractAt("ResearchContributer", addr.researchContributer);
  const poseidonLibs = { PoseidonT3: addr.poseidonT3 };
  const History = await ethers.getContractFactory("History", { libraries: poseidonLibs });
  const Record = await ethers.getContractFactory("Record", { libraries: poseidonLibs });

  const net = await ethers.provider.getNetwork();
  SNAPSHOT.provider = `${net.name}@${net.rpcUrl ?? ethers.provider.connection?.url ?? process.env.ETH_RPC ?? "n/a"}`;
  SNAPSHOT.client = await ethers.provider.send("web3_clientVersion").catch(() => "n/a");
  SNAPSHOT.environment = {
    networkLatency: process.env.BENCH_NET_LATENCY || "zero (loopback)",
    node: process.env.BENCH_NODE_DESC || "single local geth Clique PoA node (2s period), loopback only",
    zkAuthVerifier: zkAuthVerifierDesc(),
    vault: "dev-grade Vault (transit engine), used by the emergency-recovery escrow only",
    pre: "etclab/pre AFGH second attempt (BLS12-381), record delegation",
    ipfs: "single IPFS peer",
  };

  const poseidon = await buildPoseidon();

  console.log("=== Setup: register users (idempotent) ===");
  const hasUser = async (userAddr) => {
    try { await registry.getUserByAddress(userAddr); return true; } catch { return false; }
  };
  const safeReg = async (fn) => {
    try { const t = await fn(); await t.wait(); } catch { /* already registered on a reused chain */ }
  };
  if (!(await hasUser(deployer.address))) await safeReg(() => registry.connect(deployer).register(rndHex(32), rndHex(32), 3));
  for (const [who, role] of [[hospital, 4], [responder, 2], [otherResponder, 2], [doctor, 1]]) {
    if (!(await hasUser(who.address))) {
      await safeReg(() => registry.connect(deployer).registerByAdmin(who.address, rndHex(32), rndHex(32), role));
    }
  }

  const patientSecret = randBig();
  const { commitment: patientCommitment } = zkauthInputs(poseidon, { secret: patientSecret, nullifierSeed: 0n, actionId: 0n });

  console.log("\n=== [1] Registration flow (patient self-register + listener) ===");
  const hfExisting = (await historyFactory.queryFilter(historyFactory.filters.HistoryCreated()))
    .find((e) => e.args.patient.toLowerCase() === patient.address.toLowerCase());
  if (hfExisting) {
    console.log("  patient already registered on this chain; skipping measured metric (use a fresh chain for the number)");
  } else {
    const reg = await measureTx("registration.register", () =>
      registry.connect(patient).register(toHex(patientCommitment), rndHex(32), 0));
    const hfP = await waitForEvent(reg.receipt.hash, "HistoryFactory.HistoryCreated", 15000);
    if (hfP == null) throw new Error("HistoryCreated not indexed by listener");
    addMetric("registration.historyCreated_Listener", hfP);
  }

  const summaryAddr = await registry.getEmergencySummaryByAddress(patient.address);
  if (summaryAddr === ethers.ZeroAddress) throw new Error("no EmergencySummary for patient");
  const emergency = await ethers.getContractAt("EmergencySummary", summaryAddr);
  const hfEvents = hfExisting ? [hfExisting] : await historyFactory.queryFilter(historyFactory.filters.HistoryCreated());
  const hf = hfEvents.find((e) => e.args.patient.toLowerCase() === patient.address.toLowerCase());
  if (!hf) throw new Error("could not find HistoryCreated for patient");
  const history = History.attach(hf.args.historyContract);
  console.log(`  patient=${patient.address} history=${history.target} summary=${summaryAddr}`);

  console.log("\n=== Setup sessions via genuine zkauth proofs (admin/doctor/responder) ===");
  const adminTok = (await realZkAuthSession(zkAuth, deployer, poseidon, { secret: randBig(), nullifierSeed: randBig(), actionId: randBig() })).token;
  const doctorTok = (await realZkAuthSession(zkAuth, doctor, poseidon, { secret: randBig(), nullifierSeed: randBig(), actionId: randBig() })).token;
  const responderTok = (await realZkAuthSession(zkAuth, responder, poseidon, { secret: randBig(), nullifierSeed: randBig(), actionId: randBig() })).token;
  console.log("  sessions minted; waiting for listener persistence");
  await new Promise((r) => setTimeout(r, 3000));

  console.log("\n=== [2] Authentication flow (patient; witness+prove+tx+listener) ===");
  let patientTok;
  for (let i = 0; i < AUTH_REPS; i++) {
    const seed = randBig();
    const actionId = randBig();
    const proof = await proveFor("zkauth", zkauthInputs(poseidon, { secret: patientSecret, nullifierSeed: seed, actionId }).inputs);
    addMetric("zkAuth.prove.witness", proof.witnessMs);
    addMetric("zkAuth.prove.total", proof.totalMs);
    const t = await measureTx("zkAuth.authenticate", () => zkAuth.connect(patient).authenticate(
      patient.address,
      toHex(patientCommitment),
      proof.formatted.pA, proof.formatted.pB, proof.formatted.pC,
      proof.formatted.publicSignals[1], proof.formatted.publicSignals[2]
    ));
    const listener = await waitForEvent(t.receipt.hash, "ZKAuthVerifier.Authenticated", 15000);
    if (listener == null) throw new Error("Authenticated not indexed by listener");
    addMetric("zkAuth.sessionListener", listener);
    if (i === 0) {
      const log = await parseEvent(t.receipt, zkAuth.interface, "Authenticated");
      patientTok = log.args.sessionToken;
    }
    if (patientTok == null) throw new Error("no session token");
  }
  await new Promise((r) => setTimeout(r, 2000));

  console.log("\n=== [3] Off-chain identity index (enrollment + resolution) ===");
  const patNid = NID("PAT");
  let m = await measureApi("identity.provision", () => provisionIdentity(adminTok, deployer.address, {
    national_id: patNid, display_name: "Bench Patient", role: "patient",
    address: patient.address, summary_address: summaryAddr,
    commitment: toHex(patientCommitment), internal_id: rndHex(32),
  }));
  m = await measureApi("identity.resolve", () => resolveIdentity(responderTok, responder.address, { national_id: patNid }));
  if (m.value.status !== 200) throw new Error("identity resolve: " + JSON.stringify(m.data));

  console.log("\n=== [4] Record lifecycle: add record + ZK policy access + access logging ===");
  const policyProof = await proveFor("policy", policyInputs({
    department: 0x1111n, institution: 0x2222n,
    doctorAddr: BigInt(doctor.address),
    validUntil: BigInt(nowSec()) + 365n * 24n * 3600n,
    evaluatedAt: BigInt(nowSec()),
  }));
  const policy = {
    pA: policyProof.formatted.pA, pB: policyProof.formatted.pB, pC: policyProof.formatted.pC,
    deptC: policyProof.formatted.publicSignals[0],
    instC: policyProof.formatted.publicSignals[1],
    docC: policyProof.formatted.publicSignals[2],
    evaluatedAt: policyProof.formatted.publicSignals[3],
  };

  const keyRetrieve = async (recordAddr) => {
    const first = await measureApi("record.keyRetrieve", () => apiGet(
      `/api/v1/records/${recordAddr.toLowerCase()}/key?doctor_address=${doctor.address}`,
      { token: doctorTok, user: doctor.address }));
    if (first.value.status === 200) return first;
    for (let i = 0; i < 15; i++) {
      await new Promise((r) => setTimeout(r, 1000));
      const res = await apiGet(
        `/api/v1/records/${recordAddr.toLowerCase()}/key?doctor_address=${doctor.address}`,
        { token: doctorTok, user: doctor.address });
      if (res.status === 200) return first;
    }
    throw new Error(`record.keyRetrieve never reached 200 for ${recordAddr} (first ${first.value.status})`);
  };

  for (let r = 0; r < RECORD_REPS; r++) {
    console.log(`  record rep ${r + 1}/${RECORD_REPS}`);
    const up = await measureApi("record.upload", () => apiUpload(Buffer.from(`BENCH-EHR-${Date.now()}-${r}`), { token: doctorTok, user: doctor.address }));
    await expectOk("record.upload", up.value);
    const { ipfs_cid, content_hash } = up.value.data;

    const cr = await measureTx("record.create", () => registry.connect(patient).createRecord(
      doctor.address, { recordType: 0, createdAt: 0, updatedAt: 0 }, content_hash, ipfs_cid
    ));
    const crLog = await parseEvent(cr.receipt, history.interface, "RecordCreated");
    if (!crLog) throw new Error("RecordCreated not emitted");
    const record = Record.attach(crLog.args.record);
    const recordAddr = record.target;

    // Key storage, A/B on the same row: the legacy Vault-wrapped blob
    // (`record.storeKeyVault`, the pre-PRE shape) and the PRE lockbox
    // (`record.storeKeyLockbox`, written by the patient). Only the payload
    // differs, so any gap is the vault round-trip vs. a pure Postgres write.
    const skVault = await measureApi("record.storeKeyVault", () => apiPost(`/api/v1/records/${recordAddr.toLowerCase()}/key`, {
      record_address: recordAddr.toLowerCase(), ipfs_cid, content_hash,
      encrypted_key: rndHex(32), patient_address: patient.address.toLowerCase(),
    }, { token: doctorTok, user: doctor.address }));
    await expectOk("record.storeKeyVault", skVault.value);

    const pre = precliSetup(PRE_STATE);
    const skLock = await measureApi("record.storeKeyLockbox", () => apiPost(`/api/v1/records/${recordAddr.toLowerCase()}/key`, {
      record_address: recordAddr.toLowerCase(), ipfs_cid, content_hash,
      encrypted_key: pre.lockbox, lockbox_salt: pre.lockbox_salt,
      patient_public_key: pre.patient_public_key,
      patient_address: patient.address.toLowerCase(),
    }, { token: patientTok, user: patient.address }));
    await expectOk("record.storeKeyLockbox", skLock.value);

    const dg2 = await measureApi("record.delegate", () => apiPost(`/api/v1/records/${recordAddr.toLowerCase()}/delegate`, {
      patient_address: patient.address.toLowerCase(),
      doctor_address: doctor.address.toLowerCase(),
      re_encryption_key: pre.re_encryption_key,
    }, { token: patientTok, user: patient.address }));
    await expectOk("record.delegate", dg2.value);

    // Wait until the record's creation block is finalized past so the
    // listener's dynamic per-record WatchAccessGranted subscription (created
    // on RecordCreated) is active before we grant access.
    const recBlock = cr.receipt.blockNumber;
    let head = recBlock;
    while (head <= recBlock) {
      await new Promise((r) => setTimeout(r, 200));
      head = await ethers.provider.getBlockNumber();
    }

    await measureTx("record.setPolicy", () => record.connect(patient).setPolicy(policy.deptC, policy.instC));
    const zk = await measureTx("record.zkGrant", () => record.connect(doctor).grantAccessWithProof(
      policy.pA, policy.pB, policy.pC, policy.deptC, policy.instC, policy.docC, policy.evaluatedAt
    ));
    const zkEv = await waitForEvent(zk.receipt.hash, "Record.AccessGranted", 15000);
    if (zkEv == null) throw new Error("Record.AccessGranted not indexed by listener");
    addMetric("record.zkGrantListener", zkEv);

    // owner direct grant -> access_logs row ("grant") via History.AccessGranted
    const dg = await measureTx("record.directGrant", () => history.connect(patient).grantRecordAccess(recordAddr, otherResponder.address));
    const al = await waitForAccessLog(dg.receipt.hash, "grant", 15000);
    if (al == null) throw new Error("access_logs row missing for direct grant");
    addMetric("record.accessLogListener", al);

    const kk = await keyRetrieve(recordAddr);
    await expectOk("record.keyRetrieve", kk.value);

    // The doctor's own PRE request path: backend re-encrypts lockbox -> C_doctor.
    // (After `record.zkGrantListener` the listener has usually already done
    // this; the endpoint is idempotent, so the measured path is the full
    // parse + AFGH transform + upsert either way.)
    const re = await measureApi("record.reencrypt", () => apiPost(`/api/v1/records/${recordAddr.toLowerCase()}/reencrypt`, {
      doctor_address: doctor.address.toLowerCase(),
    }, { token: doctorTok, user: doctor.address }));
    await expectOk("record.reencrypt", re.value);

    await measureApi("record.download", async () => {
      const res = await fetch(`${API_BASE}/api/v1/records/${recordAddr.toLowerCase()}/download`, {
        headers: { "X-User-Address": doctor.address, "X-Session-Token": doctorTok },
      });
      return { status: res.status, data: res.ok ? { ok: true } : null };
    }).then(async (r) => { await expectOk("record.download", r.value); });
  }

  console.log("\n=== [5] Research flow (consent + attest + contribute) ===");
  const condGroup = 42n;
  const rProof = await proveFor("research", researchInputs(poseidon, {
    secret: randBig(), nullifierSeed: randBig(), conditionGroup: condGroup,
  }).inputs);
  const [rComm, rNull, rRoot] = rProof.formatted.publicSignals;
  const rCommHex = ethers.toBeHex(BigInt(rComm), 32);
  const rNullHex = ethers.toBeHex(BigInt(rNull), 32);
  if (!(await hasUser(patientR.address))) {
    await measureTx("research.register", () => registry.connect(patientR).register(rCommHex, rndHex(32), 0));
  } else {
    console.log("  patientR already registered (fresh-chain flow registers with the research commitment)");
  }
  await measureTx("research.optIn", () => attestation.connect(patientR).optIn(rCommHex));
  await measureTx("research.attest", () => attestation.connect(doctor).attest(rCommHex, ethers.toBeHex(condGroup, 32)));
  const rootOnChain = (await attestation.getRoot()).toString();
  if (rootOnChain !== BigInt(rRoot).toString()) throw new Error("tree root mismatch");
  await measureTx("research.contribute", () => contributer.connect(patientR).contribute(
    rProof.formatted.pA, rProof.formatted.pB, rProof.formatted.pC,
    rCommHex, rNullHex, BigInt(rRoot)
  ));
  if (!(await contributer.usedNullifiers(rNull))) throw new Error("nullifier not consumed");

  console.log("\n=== [6] Emergency flow ===");
  const locationHash = rndHex(32);
  const escrowKey = rndHex(32);

  const inc = await measureTx("emergency.createIncident", () => incident.connect(hospital).createIncident(locationHash, nowSec() + 3600));
  const incLog = await parseEvent(inc.receipt, incident.interface, "IncidentCreated");
  const incidentId = incLog.args.id;
  await measureTx("emergency.addResponders", () => incident.connect(hospital).addResponders(incidentId, [responder.address]));
  await measureTx("emergency.updateData", () => emergency.connect(patient).updateData(rndHex(64), rndHex(32)));
  await measureTx("emergency.registerKey", () => emergency.connect(responder).registerPublicKey(rndHex(32)));
  const req = await measureTx("emergency.requestAccess", () => emergency.connect(responder).requestAccess(rndHex(32), "St 1", "Cardiac arrest", "bench-id"));
  const reqLog = await parseEvent(req.receipt, emergency.interface, "AccessRequested");
  const ap = await measureTx("emergency.approveAccess", () => emergency.connect(patient).approveAccessRequest(reqLog.args.requestId, escrowKey));
  const apLog = await waitForAccessLog(ap.receipt.hash, "emergency_grant", 15000);
  if (apLog == null) throw new Error("emergency access_logs row missing");
  addMetric("emergency.accessLog_Listener", apLog);

  await measureTx("emergency.emergencyGrant", () => riskAssess.connect(responder).emergencyGrant(
    summaryAddr, incidentId, locationHash, toHex(patientCommitment), rndHex(32)
  ));

  await measureApi("emergency.getData", async () => {
    const [data] = await emergency.connect(responder).getData();
    return { status: data.length > 0 ? 200 : 404, data: { ok: true } };
  });

  const stb = await measureApi("emergency.storeRecoveryKeyBlob", () => apiPost("/api/v1/emergency/recovery-key",
    { emergency_summary: summaryAddr, content_key: escrowKey }, { token: patientTok, user: patient.address }));
  await expectOk("emergency.storeRecoveryKeyBlob", stb.value);
  const wrapped = stb.value.data.recovery_key;
  await measureTx("emergency.storeRecoveryKeyChain", () => emergency.connect(patient).storeRecoveryKey("0x" + wrapped));
  const rer = await measureApi("emergency.reencryptRecoveryKey", () => apiPost("/api/v1/emergency/recovery-key/re-encrypt",
    { emergency_summary: summaryAddr, recipient: responder.address }, { token: responderTok, user: responder.address }));
  await expectOk("emergency.reencryptRecoveryKey", rer.value);

  console.log("\n=== [7] Real on-chain Groth16 verification ===");
  const Groth16Verifier = await ethers.getContractFactory("Groth16Verifier");
  const vDeploy = await Groth16Verifier.deploy();
  await vDeploy.waitForDeployment();
  addMetric("groth16.deploy", 0, { gas: (await vDeploy.deploymentTransaction().wait()).gasUsed.toString() });
  const verifier = vDeploy;
  const vProof = await proveFor("zkauth", zkauthInputs(poseidon, { secret: randBig(), nullifierSeed: randBig(), actionId: randBig() }).inputs);
  const [c1, c2, c3] = vProof.formatted.publicSignals;
  const vArgs = [vProof.formatted.pA, vProof.formatted.pB, vProof.formatted.pC, [BigInt(c1), BigInt(c2), BigInt(c3)]];
  await measureApi("groth16.verifyProof", async () => {
    if (!(await verifier.verifyProof(...vArgs))) throw new Error("real zkauth proof rejected on-chain");
    return { status: 200, data: { ok: true } };
  });
  const estGas = await ethers.provider.estimateGas({
    to: verifier.target,
    data: verifier.interface.encodeFunctionData("verifyProof", vArgs),
  });
  addMetric("groth16.verifyProof", 0, { gas: estGas.toString() });

  // ---------------------------------------------------------------- report
  console.log("\n\n=== LATENCY (ms). For txs, total = RPC send + block inclusion ===");
  const keys = Object.keys(metrics);
  const rows = [["metric", "n", "mean", "p50", "p95", "p99", "min", "max"]];
  for (const k of keys) {
    const s = stat(metrics[k].map((e) => e.ms));
    if (!s) continue;
    const f = (x) => (x == null ? "-" : x.toFixed(1));
    rows.push([k, String(s.n), f(s.mean), f(s.p50), f(s.p95), f(s.p99), f(s.min), f(s.max)]);
  }
  printTable(rows);

  console.log("\n=== GAS (per operation) ===");
  const gasRows = [["metric", "gas used (mean)"]];
  for (const k of keys) {
    const gs = metrics[k].filter((e) => e.gas != null).map((e) => Number(e.gas));
    if (!gs.length) continue;
    const mean = Math.round(gs.reduce((a, b) => a + b, 0) / gs.length);
    gasRows.push([k, mean.toLocaleString("en-US")]);
  }
  printTable(gasRows);

  const latencyJson = {
    snapshot: SNAPSHOT,
    unit: "milliseconds",
    results: {},
  };
  for (const [k, entries] of Object.entries(metrics)) {
    const s = stat(entries.map((e) => e.ms));
    const gs = entries.filter((e) => e.gas != null).map((e) => e.gas);
    latencyJson.results[k] = { ...s, gasMean: gs.length ? (gs.reduce((a, b) => Number(a) + Number(b), 0) / gs.length).toFixed(0) : null, samples: entries };
  }
  const gasJson = { snapshot: SNAPSHOT, unit: "gas", results: {} };
  for (const [k, entries] of Object.entries(metrics)) {
    const gs = entries.filter((e) => e.gas != null).map((e) => e.gas);
    if (gs.length) gasJson.results[k] = { gasMean: (gs.reduce((a, b) => Number(a) + Number(b), 0) / gs.length).toFixed(0), samples: gs };
  }
  const lPath = writeJSON("latency.json", latencyJson);
  const gPath = writeJSON("gas.json", gasJson);
  console.log(`\nWrote ${lPath}\n      ${gPath}`);
}

main().then(() => process.exit(0)).catch((err) => { console.error("bench_latency failed:", err); process.exit(1); });