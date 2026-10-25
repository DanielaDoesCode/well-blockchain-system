import "dotenv/config";
import hre from "hardhat";
import { execFileSync } from "child_process";
import path from "path";
import { fileURLToPath } from "url";
import { precliSetup, precliSelf, precliVerify } from "./lib/precli.js";

const { ethers } = hre;

const API_BASE = "http://localhost:8081";
const __dirname = path.dirname(fileURLToPath(import.meta.url));
const GEN_PROOF = path.join(__dirname, "..", "zk-circuits", "gen_policy_proof.js");

// The doctor's "Vault attributes" are issued at runtime: the admin mints
// department/institution UUIDs and assigns them to the doctor via the API;
// the doctor's ZK proof is generated from those attributes (Vault -> API ->
// proof generator). Only their commitments ever reach the chain.

async function apiPost(path, body, extraHeaders = {}) {
  const headers = { "Content-Type": "application/json", ...extraHeaders };
  const res = await fetch(`${API_BASE}${path}`, {
    method: "POST", headers, body: JSON.stringify(body),
  });
  const data = await res.json();
  if (!res.ok) throw new Error(`POST ${path}: ${res.status} ${JSON.stringify(data)}`);
  return data;
}

async function apiUpload(fileData, sessionToken, userAddress) {
  const form = new FormData();
  form.append("file", new Blob([fileData]));
  const res = await fetch(`${API_BASE}/api/v1/records/upload`, {
    method: "POST",
    headers: { "X-User-Address": userAddress, "X-Session-Token": sessionToken },
    body: form,
  });
  const data = await res.json();
  if (!res.ok) throw new Error(`upload: ${res.status} ${JSON.stringify(data)}`);
  return data;
}

async function apiGet(path, sessionToken, userAddress) {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { "X-User-Address": userAddress, "X-Session-Token": sessionToken },
  });
  return { status: res.status, data: await res.json() };
}

function genProof({
  doctor, department, institution,
  requiredDept = department, requiredInst = institution,
  validUntil, evaluatedAt, sessionToken, apiBase,
}) {
  const args = [];
  if (apiBase && sessionToken) {
    args.push("--api-base", apiBase, "--user-address", doctor, "--session-token", sessionToken);
  }
  if (department !== undefined) args.push("--department", department);
  if (institution !== undefined) args.push("--institution", institution);
  if (requiredDept !== undefined) args.push("--requiredDept", requiredDept);
  if (requiredInst !== undefined) args.push("--requiredInst", requiredInst);
  if (doctor !== undefined) args.push("--doctorAddr", doctor);
  if (validUntil !== undefined) args.push("--validUntil", String(validUntil));
  if (evaluatedAt !== undefined) args.push("--evaluatedAt", String(evaluatedAt));
  const out = execFileSync("node", [GEN_PROOF, ...args], {
    encoding: "utf8", maxBuffer: 1024 * 1024,
    stdio: ["ignore", "pipe", "pipe"],
  });
  return JSON.parse(out);
}

async function zkAuthSession(zkAuth, signer) {
  const dummyProof = [[0, 0], [[0, 0], [0, 0]], [0, 0]];
  const tx = await zkAuth.connect(signer).authenticate(
    signer.address,
    ethers.hexlify(ethers.randomBytes(32)),
    dummyProof[0], dummyProof[1], dummyProof[2],
    ethers.hexlify(ethers.randomBytes(32)),
    ethers.hexlify(ethers.randomBytes(32))
  );
  const receipt = await tx.wait();
  const authLog = receipt.logs
    .map(l => { try { return zkAuth.interface.parseLog({ topics: [...l.topics], data: l.data }); } catch { return null; } })
    .find(e => e?.name === "Authenticated");
  if (!authLog) throw new Error("could not find Authenticated event for " + signer.address);
  return authLog.args.sessionToken;
}

async function expectRevert(label, fn, selector) {
  try {
    await fn();
  } catch (e) {
    const reason = e.reason ?? e.message ?? "";
    if (selector && !reason.includes(selector)) {
      throw new Error(`${label}: reverted with unexpected reason "${reason}" (wanted ${selector})`);
    }
    console.log(`  [PASS] ${label}`);
    return;
  }
  throw new Error(`${label}: expected revert, but transaction succeeded`);
}

async function main() {
  const [deployer, doctor, otherDoctor, patient] = await ethers.getSigners();

  const addr = {
    registry:       process.env.REGISTRY_ADDRESS,
    userRegistry:   process.env.USER_REGISTRY_ADDRESS,
    historyFactory: process.env.HISTORY_FACTORY_ADDRESS,
    zkAuth:         process.env.ZKAUTH_VERIFIER_ADDRESS,
  };
  for (const [k, v] of Object.entries(addr)) {
    if (!v) throw new Error(`missing ${k} in .env`);
  }

  const registry       = await ethers.getContractAt("WellRegistry", addr.registry);
  const userReg        = await ethers.getContractAt("UserRegistry", addr.userRegistry);
  const historyFactory = await ethers.getContractAt("HistoryFactory", addr.historyFactory);
  const zkAuth         = await ethers.getContractAt("ZKAuthVerifier", addr.zkAuth);

  console.log("=== Register users ===");

  const adminCommitment = ethers.hexlify(ethers.randomBytes(32));
  const adminInternalID = ethers.hexlify(ethers.randomBytes(32));
  let tx = await registry.connect(deployer).register(adminCommitment, adminInternalID, 3);
  await tx.wait();

  const doctorCommitment = ethers.hexlify(ethers.randomBytes(32));
  const doctorInternalID = ethers.hexlify(ethers.randomBytes(32));
  tx = await registry.connect(deployer).registerByAdmin(doctor.address, doctorCommitment, doctorInternalID, 1);
  await tx.wait();

  tx = await registry.connect(deployer).registerByAdmin(otherDoctor.address,
    ethers.hexlify(ethers.randomBytes(32)), ethers.hexlify(ethers.randomBytes(32)), 1);
  await tx.wait();

  const patientCommitment = ethers.hexlify(ethers.randomBytes(32));
  const patientInternalID = ethers.hexlify(ethers.randomBytes(32));
  tx = await registry.connect(patient).register(patientCommitment, patientInternalID, 0);
  const patientReceipt = await tx.wait();

  const hfEvents = await historyFactory.queryFilter(
    historyFactory.filters.HistoryCreated(),
    patientReceipt.blockNumber, patientReceipt.blockNumber
  );
  const historyAddr = hfEvents.find(e => e.args.patient === patient.address)?.args.historyContract;
  if (!historyAddr) throw new Error("could not find HistoryCreated event");
  const poseidonLibs = { PoseidonT3: process.env.POSEIDON_T3_ADDRESS };
  const History = await ethers.getContractFactory("History", { libraries: poseidonLibs });
  const history = History.attach(historyAddr);
  console.log("  Patient's History:", historyAddr);

  console.log("\n=== ZK auth sessions (mock verifier) ===");
  const doctorSession = await zkAuthSession(zkAuth, doctor);
  const adminSession = await zkAuthSession(zkAuth, deployer);
  const patientSession = await zkAuthSession(zkAuth, patient);
  console.log("  doctor session:", doctorSession.slice(0, 10) + "...");
  console.log("  admin  session:", adminSession.slice(0, 10) + "...");
  console.log("  Waiting for listener to process sessions (3s)...");
  await new Promise(r => setTimeout(r, 3000));

  console.log("\n=== Upload encrypted EHR ===");
  const ehrData = Buffer.from("POLICY-PATIENT-REPORT-" + Date.now());
  const uploadResult = await apiUpload(ehrData, doctorSession, doctor.address);
  console.log("  IPFS CID:", uploadResult.ipfs_cid);

  console.log("\n=== Create record on-chain ===");
  const recordType = { recordType: 0, createdAt: 0, updatedAt: 0 };
  tx = await registry.connect(patient).createRecord(
    doctor.address, recordType, uploadResult.content_hash, uploadResult.ipfs_cid
  );
  const createReceipt = await tx.wait();
  const recordLog = createReceipt.logs
    .map(l => { try { return history.interface.parseLog({ topics: [...l.topics], data: l.data }); } catch { return null; } })
    .find(e => e?.name === "RecordCreated");
  if (!recordLog) throw new Error("could not find RecordCreated event");
  const recordAddr = recordLog.args.record;
  const Record = await ethers.getContractFactory("Record", { libraries: poseidonLibs });
  const record = Record.attach(recordAddr);
  const permissions = await ethers.getContractAt("RecordPermissions", await record.permissions());
  console.log("  Record:", recordAddr);

  console.log("\n=== Seal PRE lockbox ===");
  const preState = precliSetup();
  await apiPost(`/api/v1/records/${recordAddr.toLowerCase()}/key`, {
    record_address: recordAddr.toLowerCase(),
    ipfs_cid: uploadResult.ipfs_cid,
    content_hash: uploadResult.content_hash,
    encrypted_key: preState.lockbox,
    lockbox_salt: preState.lockbox_salt,
    patient_public_key: preState.patient_public_key,
    patient_address: patient.address.toLowerCase(),
  }, {
    "X-User-Address": patient.address,
    "X-Session-Token": patientSession,
  });
  await apiPost(`/api/v1/records/${recordAddr.toLowerCase()}/delegate`, {
    patient_address: patient.address.toLowerCase(),
    doctor_address: doctor.address.toLowerCase(),
    re_encryption_key: preState.re_encryption_key,
  }, {
    "X-User-Address": patient.address,
    "X-Session-Token": patientSession,
  });
  console.log("  C_patient and rk_{patient->doctor} stored");

  console.log("\n=== Issue & assign attributes (admin, on-chain role gate) ===");
  const oneYearSec = BigInt(365 * 24 * 3600);
  const validUntil = String(BigInt(Math.floor(Date.now() / 1000)) + oneYearSec);
  const adminHeaders = { "X-User-Address": deployer.address, "X-Session-Token": adminSession };
  const cardio = await apiPost("/api/v1/attributes/issue", { kind: "department", name: "cardiology" }, adminHeaders);
  const peds = await apiPost("/api/v1/attributes/issue", { kind: "department", name: "pediatrics" }, adminHeaders);
  const general = await apiPost("/api/v1/attributes/issue", { kind: "institution", name: "general" }, adminHeaders);
  console.log("  cardiology   ", cardio.uuid);
  console.log("  pediatrics   ", peds.uuid);
  console.log("  general      ", general.uuid);

  const regCheck = await apiGet("/api/v1/attributes/registry/department/cardiology", adminSession, deployer.address);
  if (regCheck.status !== 200 || regCheck.data.uuid !== cardio.uuid) {
    throw new Error("registry lookup mismatch");
  }
  console.log("  [PASS] registry lookup resolves issued UUID");

  await apiPost("/api/v1/attributes/assign", {
    doctor_address: doctor.address.toLowerCase(),
    department: cardio.uuid,
    institution: general.uuid,
    valid_until: validUntil,
  }, adminHeaders);
  console.log("  assigned cardiology+general to", doctor.address);

  const attrsCheck = await apiGet(`/api/v1/attributes/${doctor.address.toLowerCase()}`, adminSession, deployer.address);
  if (attrsCheck.status !== 200 || attrsCheck.data.department !== cardio.uuid.toLowerCase()) {
    throw new Error("assigned attributes not readable via API");
  }
  console.log("  [PASS] attributes stored in Vault and served by API");

  console.log("\n=== Generate policy proof (attributes fetched from Vault via API) ===");
  const proof = genProof({ doctor: doctor.address, apiBase: API_BASE, sessionToken: doctorSession });
  console.log("  deptCommitment    ", proof.deptCommitment);
  console.log("  instCommitment    ", proof.instCommitment);
  console.log("  doctorCommitment  ", proof.doctorCommitment);
  console.log("  evaluatedAt       ", parseInt(proof.evaluatedAt, 16), "(" + new Date(parseInt(proof.evaluatedAt, 16) * 1000).toISOString() + ")");

  console.log("\n=== Negative: wrong policy (peds dept) rejects proof ===");
  const pedsProof = genProof({
    doctor: doctor.address,
    department: peds.uuid, institution: general.uuid,
    requiredDept: peds.uuid, requiredInst: general.uuid,
  });
  tx = await record.connect(patient).setPolicy(pedsProof.deptCommitment, pedsProof.instCommitment);
  await tx.wait();
  await expectRevert("wrong-UUID proof rejected",
    () => record.connect(doctor).grantAccessWithProof(
      proof.proof.pA, proof.proof.pB, proof.proof.pC,
      proof.deptCommitment, proof.instCommitment, proof.doctorCommitment, proof.evaluatedAt
    ),
    "PolicyMismatch");
  const accessAfterWrong = await permissions.hasAccess(doctor.address);
  if (accessAfterWrong) throw new Error("doctor gained access despite wrong policy");

  console.log("\n=== Set correct policy ===");
  tx = await record.connect(patient).setPolicy(proof.deptCommitment, proof.instCommitment);
  await tx.wait();
  console.log("  Policy set to proof's commitments");

  console.log("\n=== Negative: replay by a different doctor rejected ===");
  await expectRevert("replay-by-other-doctor rejected",
    () => record.connect(otherDoctor).grantAccessWithProof(
      proof.proof.pA, proof.proof.pB, proof.proof.pC,
      proof.deptCommitment, proof.instCommitment, proof.doctorCommitment, proof.evaluatedAt
    ),
    "PolicyMismatch");

  console.log("\n=== Doctor self-grants via proof ===");
  tx = await record.connect(doctor).grantAccessWithProof(
    proof.proof.pA, proof.proof.pB, proof.proof.pC,
    proof.deptCommitment, proof.instCommitment, proof.doctorCommitment, proof.evaluatedAt
  );
  const grantReceipt = await tx.wait();
  const grantLog = grantReceipt.logs
    .map(l => { try { return record.interface.parseLog({ topics: [...l.topics], data: l.data }); } catch { return null; } })
    .find(e => e?.name === "AccessGranted");
  if (!grantLog) throw new Error("AccessGranted not emitted on self-grant");
  if (grantLog.args.doctor !== doctor.address) throw new Error("AccessGranted.doctor != caller");
  if (!await permissions.hasAccess(doctor.address)) {
    throw new Error("doctor not authorized after valid proof");
  }
  console.log("  [PASS] doctor authorized after self-grant (event: doctor=" + grantLog.args.doctor + ", grantedBy=" + grantLog.args.grantedBy + ")");

  console.log("\n=== PRE triggered by AccessGranted ===");
  console.log("  Waiting for listener to re-encrypt (3s)...");
  await new Promise(r => setTimeout(r, 3000));
  const keyResult = await apiGet(
    `/api/v1/records/${recordAddr.toLowerCase()}/key?doctor_address=${doctor.address}`,
    doctorSession, doctor.address
  );
  if (keyResult.status !== 200 || !keyResult.data.reencrypted_key) {
    throw new Error("C_doctor not available after self-grant (status " + keyResult.status + ")");
  }
  console.log("  [PASS] C_doctor retrieved via PRE");
  const patientView = precliSelf(preState.statePath, preState.lockbox);
  const doctorView = precliVerify(preState.statePath, keyResult.data.reencrypted_key);
  if (patientView !== doctorView) throw new Error("PRE round-trip mismatch");
  console.log("  [PASS] doctor recovered the sealed key through PRE");

  console.log("\n=== Negative: replayed (stale) proof after TTL rejected ===");
  await ethers.provider.send("evm_increaseTime", [Number(await record.POLICY_PROOF_TTL()) + 1]);
  await ethers.provider.send("evm_mine", []);
  await expectRevert("stale proof rejected",
    () => record.connect(doctor).grantAccessWithProof(
      proof.proof.pA, proof.proof.pB, proof.proof.pC,
      proof.deptCommitment, proof.instCommitment, proof.doctorCommitment, proof.evaluatedAt
    ),
    "StaleProof");

  console.log("\n=== Negative: expired credential cannot even produce a proof ===");
  const now = Math.floor(Date.now() / 1000);
  let expiredFailed = false;
  try {
    genProof({ doctor: doctor.address, validUntil: now - 3600, evaluatedAt: now });
  } catch (e) {
    expiredFailed = true;
    const lines = String(e.stderr).split("\n");
    const ourLine = lines.find(l => l.includes("gen_policy_proof error:"));
    const reason = (ourLine || lines[0]).replace(/^gen_policy_proof error: /, "").trim();
    console.log("  (reason: " + reason + ")");
  }
  if (!expiredFailed) throw new Error("expired credential produced a proof (circuit accepted validUntil < evaluatedAt)");
  console.log("  [PASS] expired credential rejected by circuit");

  console.log("\n=== Revoke access ===");
  tx = await history.connect(patient).revokeRecordAccess(recordAddr, doctor.address);
  await tx.wait();
  if (await permissions.hasAccess(doctor.address)) {
    throw new Error("doctor still authorized after revoke");
  }
  console.log("  Waiting for listener to delete C_doctor (2s)...");
  await new Promise(r => setTimeout(r, 2000));
  const revokedKey = await apiGet(
    `/api/v1/records/${recordAddr.toLowerCase()}/key?doctor_address=${doctor.address}`,
    doctorSession, doctor.address
  );
  if (revokedKey.status !== 404) {
    throw new Error("expected 404 after revoke, got " + revokedKey.status);
  }
  console.log("  [PASS] C_doctor deleted after revoke");

  console.log("\n=== E2E policy test passed! ===");
}

main().catch(err => {
  console.error("\nE2E policy test failed:", err.message);
  process.exit(1);
});
