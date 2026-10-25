import "dotenv/config";
import hre from "hardhat";
import { precliSetup, precliSelf, precliVerify } from "./lib/precli.js";
const { ethers } = hre;

const API_BASE = "http://localhost:8081";

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

async function main() {
  const [deployer, doctor, patient] = await ethers.getSigners();

  const addr = {
    registry:            process.env.REGISTRY_ADDRESS,
    userRegistry:        process.env.USER_REGISTRY_ADDRESS,
    historyFactory:      process.env.HISTORY_FACTORY_ADDRESS,
    zkAuth:              process.env.ZKAUTH_VERIFIER_ADDRESS,
    emergencySummary:    process.env.EMERGENCY_SUMMARY_ADDRESS,
    incidentManager:     process.env.INCIDENT_MANAGER_ADDRESS,
    riskAssessment:      process.env.RISK_ASSESSMENT_ADDRESS,
    researchContributer: process.env.RESEARCH_CONTRIBUTER_ADDRESS,
  };
  for (const [k, v] of Object.entries(addr)) {
    if (!v) throw new Error(`missing ${k} in .env`);
  }

  const registry       = await ethers.getContractAt("WellRegistry",addr.registry);
  const userReg        = await ethers.getContractAt("UserRegistry",addr.userRegistry);
  const historyFactory = await ethers.getContractAt("HistoryFactory",addr.historyFactory);
  const zkAuth         = await ethers.getContractAt("ZKAuthVerifier",addr.zkAuth);

  console.log("Registry:",addr.registry);
  console.log("ZKAuthVerifier:",addr.zkAuth);
  console.log("HistoryFactory:",addr.historyFactory);

  console.log("\n=== Register users ===");

  const adminCommitment = ethers.hexlify(ethers.randomBytes(32));
  const adminInternalID = ethers.hexlify(ethers.randomBytes(32));
  let tx = await registry.connect(deployer).register(adminCommitment, adminInternalID, 3);
  await tx.wait();
  console.log("  Admin registered");

  const doctorCommitment = ethers.hexlify(ethers.randomBytes(32));
  const doctorInternalID = ethers.hexlify(ethers.randomBytes(32));
  tx = await registry.connect(deployer).registerByAdmin(doctor.address, doctorCommitment, doctorInternalID, 1);
  await tx.wait();
  console.log("Doctor registered:", doctor.address);

  const patientCommitment = ethers.hexlify(ethers.randomBytes(32));
  const patientInternalID = ethers.hexlify(ethers.randomBytes(32));
  tx = await registry.connect(patient).register(patientCommitment, patientInternalID, 0);
  const patientReceipt = await tx.wait();
  console.log("  Patient registered:", patient.address);

  const hfEvents = await historyFactory.queryFilter(
    historyFactory.filters.HistoryCreated(),
    patientReceipt.blockNumber, patientReceipt.blockNumber
  );
  const historyAddr = hfEvents.find(e => e.args.patient === patient.address)?.args.historyContract;
  if (!historyAddr) throw new Error("could not find HistoryCreated event");
  const poseidonLibs = { PoseidonT3: process.env.POSEIDON_T3_ADDRESS };
  const History = await ethers.getContractFactory("History", { libraries: poseidonLibs });
  const history = History.attach(historyAddr);
  console.log("  Patient's History contract:", historyAddr);

  console.log("\n=== Create ZK sessions ===");

  const dummyProof = [[0, 0], [[0, 0], [0, 0]], [0, 0]];
  async function createSession(signer) {
    const t = await zkAuth.connect(signer).authenticate(
      signer.address,
      ethers.hexlify(ethers.randomBytes(32)),
      dummyProof[0], dummyProof[1], dummyProof[2],
      ethers.hexlify(ethers.randomBytes(32)),
      ethers.hexlify(ethers.randomBytes(32))
    );
    const receipt = await t.wait();
    const log = receipt.logs
      .map(l => { try { return zkAuth.interface.parseLog({ topics: [...l.topics], data: l.data }); } catch { return null; } })
      .find(e => e?.name === "Authenticated");
    if (!log) throw new Error("could not find Authenticated event");
    return log.args.sessionToken;
  }

  const sessionToken = await createSession(doctor);
  const patientSessionToken = await createSession(patient);
  console.log("  Doctor session token on-chain:", sessionToken);
  console.log("  Patient session token on-chain:", patientSessionToken);

  //waiting to populate sessions
  console.log("  Waiting for listener to process sessions (3s)...");
  await new Promise(r => setTimeout(r, 3000));

  //upload file
  console.log("\n=== Upload encrypted EHR ===");

  const ehrData = Buffer.from("PATIENT-DIAGNOSIS-REPORT-" + Date.now());
  const uploadResult = await apiUpload(ehrData, sessionToken, doctor.address);
  console.log("  IPFS CID:", uploadResult.ipfs_cid);
  console.log("  Content hash:", uploadResult.content_hash);

  //create record on chain
  console.log("\n=== Create record on-chain ===");

  const recordType = { recordType: 0, createdAt: 0, updatedAt: 0 };
  const contentHashHex = uploadResult.content_hash;
  tx = await registry.connect(patient).createRecord(
    doctor.address, recordType, contentHashHex, uploadResult.ipfs_cid
  );
  const createReceipt = await tx.wait();

  const recordLog = createReceipt.logs
    .map(l => {
      try { return history.interface.parseLog({ topics: [...l.topics], data: l.data }); }
      catch { return null; }
    })
    .find(e => e?.name === "RecordCreated");
  if (!recordLog) throw new Error("could not find RecordCreated event");
  const recordAddr = recordLog.args.record;
  console.log("  Record contract:", recordAddr);

  const Record = await ethers.getContractFactory("Record", { libraries: poseidonLibs });
  const record = Record.attach(recordAddr);
  const onChainRef = await record.storageReference();
  const onChainHash = await record.contentHash();
  console.log("  On-chain storageRef matches:", onChainRef === uploadResult.ipfs_cid ? "YES" : "NO");
  console.log("  On-chain contentHash matches:", onChainHash === contentHashHex ? "YES" : "NO");

  //seal a PRE lockbox and register the patient-to-doctor re-encryption key
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
    "X-Session-Token": patientSessionToken,
  });
  console.log("  C_patient (lockbox) stored");

  await apiPost(`/api/v1/records/${recordAddr.toLowerCase()}/delegate`, {
    patient_address: patient.address.toLowerCase(),
    doctor_address: doctor.address.toLowerCase(),
    re_encryption_key: preState.re_encryption_key,
  }, {
    "X-User-Address": patient.address,
    "X-Session-Token": patientSessionToken,
  });
  console.log("  rk_{patient->doctor} stored");

  //trigger PRE
  console.log("\n=== Grant access / PRE ===");

  tx = await history.connect(patient).grantRecordAccess(recordAddr, doctor.address);
  await tx.wait();
  console.log("  Access granted, waiting for PRE (3s)...");
  await new Promise(r => setTimeout(r, 3000));

  //retrieve C_doctor and verify the doctor can open it
  console.log("\n=== Retrieve C_doctor ===");

  const getKeyResult = await apiGet(
    `/api/v1/records/${recordAddr.toLowerCase()}/key?doctor_address=${doctor.address}`,
    sessionToken, doctor.address
  );
  console.log("  Status:", getKeyResult.status);
  console.log("  C_doctor retrieved:", getKeyResult.data.reencrypted_key ? "YES" : "NO");
  if (getKeyResult.status !== 200) throw new Error("re-encrypted key not found");

  const patientView = precliSelf(preState.statePath, preState.lockbox);
  const doctorView = precliVerify(preState.statePath, getKeyResult.data.reencrypted_key);
  console.log("  Patient self-recovery matches doctor view:", patientView === doctorView ? "YES" : "NO");
  if (patientView !== doctorView) throw new Error("PRE round-trip mismatch");

  //download block
  console.log("\n=== Download blob ===");

  const dlRes = await fetch(
    `${API_BASE}/api/v1/records/${recordAddr.toLowerCase()}/download`,
    { headers: { "X-User-Address": doctor.address, "X-Session-Token": sessionToken } }
  );
  if (!dlRes.ok) throw new Error(`download: ${dlRes.status}`);
  console.log("  Download status:", dlRes.status);

  //revoke access, check if key gets deleted
  console.log("\n=== Revoke access / verify cleanup ===");

  tx = await history.connect(patient).revokeRecordAccess(recordAddr, doctor.address);
  await tx.wait();
  console.log("  Access revoked, waiting (2s)...");
  await new Promise(r => setTimeout(r, 2000));

  const revokedResult = await apiGet(
    `/api/v1/records/${recordAddr.toLowerCase()}/key?doctor_address=${doctor.address}`,
    sessionToken, doctor.address
  );
  console.log("  After revoke, key lookup:", revokedResult.status);
  if (revokedResult.status !== 404) {
    throw new Error("expected 404 after revoke, got " + revokedResult.status);
  }

  console.log("\n=== E2E test passed! ===");
  console.log("\nSummary:");
  console.log("Record contract:", recordAddr);
  console.log("IPFS CID:", uploadResult.ipfs_cid);
  console.log("Content hash:", uploadResult.content_hash);
  console.log("History contract:", historyAddr);
}

main().catch(err => {
  console.error("\n E2E test failed:", err.message);
  process.exit(1);
});
