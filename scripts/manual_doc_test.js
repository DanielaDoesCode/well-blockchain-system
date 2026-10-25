import "dotenv/config";
import hre from "hardhat";
import { execFileSync } from "child_process";
import path from "path";
import { fileURLToPath } from "url";
import { readFileSync, writeFileSync, mkdirSync } from "fs";
import { createHash } from "crypto";
import { precliSetup, precliSelf, precliVerify } from "./lib/precli.js";

// manual_doc_test.js
//
// End-to-end manual test of the full ZK-policy flow against a REAL document.
// Run it via `npx hardhat run` (hardhat does not forward CLI args, so the
// values are read from environment variables):
//
//   MANUAL_DOC_FILE=./my-report.pdf \
//     MANUAL_DOC_DEPARTMENT=cardiology \
//     MANUAL_DOC_INSTITUTION=general \
//     npx hardhat run scripts/manual_doc_test.js --network localhost
//
// Optional env vars:
//   MANUAL_DOC_FILE         path to the real document to upload (required)
//   MANUAL_DOC_DEPARTMENT   attribute name to mint (default: cardiology)
//   MANUAL_DOC_INSTITUTION  attribute name to mint (default: general)
//
// On success the fetched document, a copy of the original, and a SHA256SUMS
// manifest are written to manual-test-output/ so the round-trip can be
// verified with:  cd manual-test-output && sha256sum -c SHA256SUMS
//
// The stack (hardhat node + postgres + ipfs + vault + backend) must already be
// running. Use scripts/run_manual_doc_test.sh to boot everything fresh.
//
// Flow: register users -> ZK sessions -> upload real file -> create record on
// chain -> store C_patient -> admin issues/assigns attributes (Vault) ->
// real Groth16 policy proof (attributes fetched from Vault via API) ->
// setPolicy -> grantAccessWithProof (self-grant) -> PRE -> download the file
// back and assert it matches the original byte-for-byte (sha256).

const { ethers } = hre;

const API_BASE = "http://localhost:8081";
const MAX_UPLOAD_BYTES = 100 << 20; // backend ParseMultipartForm limit (100 MB)
const __dirname = path.dirname(fileURLToPath(import.meta.url));
const GEN_PROOF = path.join(__dirname, "..", "zk-circuits", "gen_policy_proof.js");

// The fetched document, the original, and a SHA256SUMS manifest are written
// here (project root / manual-test-output).
const OUTPUT_DIR = path.join(__dirname, "..", "manual-test-output");

function parseArgs(argv) {
  const args = {};
  for (let i = 0; i < argv.length - 1; i += 1) {
    if (argv[i].startsWith("--")) {
      args[argv[i].slice(2)] = argv[i + 1];
    }
  }
  return args;
}

async function apiPost(pathname, body, extraHeaders = {}) {
  const headers = { "Content-Type": "application/json", ...extraHeaders };
  const res = await fetch(`${API_BASE}${pathname}`, {
    method: "POST", headers, body: JSON.stringify(body),
  });
  const data = await res.json();
  if (!res.ok) throw new Error(`POST ${pathname}: ${res.status} ${JSON.stringify(data)}`);
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

async function apiGet(pathname, sessionToken, userAddress) {
  const res = await fetch(`${API_BASE}${pathname}`, {
    headers: { "X-User-Address": userAddress, "X-Session-Token": sessionToken },
  });
  return { status: res.status, data: await res.json() };
}

function sha256hex(buf) {
  return createHash("sha256").update(buf).digest("hex");
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

async function main() {
  const args = parseArgs(process.argv);

  const filePath = args.file ?? process.env.MANUAL_DOC_FILE;
  if (!filePath) {
    throw new Error("Missing document path. Set MANUAL_DOC_FILE=<path> (or pass --file when running under plain node).");
  }

  console.log("=== Loading real document ===");
  let fileData;
  try {
    fileData = readFileSync(filePath);
  } catch (e) {
    throw new Error(`cannot read file "${filePath}" (working directory: ${process.cwd()}): ${e.message}`);
  }
  if (fileData.length > MAX_UPLOAD_BYTES) {
    throw new Error(`file is ${fileData.length} bytes; backend upload limit is ${MAX_UPLOAD_BYTES} bytes (100 MB)`);
  }
  const fileHash = sha256hex(fileData);
  console.log(`  file:   ${filePath}`);
  console.log(`  size:   ${fileData.length} bytes`);
  console.log(`  sha256: 0x${fileHash}`);

  const deptName = args.department || process.env.MANUAL_DOC_DEPARTMENT || "cardiology";
  const instName = args.institution || process.env.MANUAL_DOC_INSTITUTION || "general";

  const [deployer, doctor, patient] = await ethers.getSigners();

  const addr = {
    registry:       process.env.REGISTRY_ADDRESS,
    userRegistry:   process.env.USER_REGISTRY_ADDRESS,
    historyFactory: process.env.HISTORY_FACTORY_ADDRESS,
    zkAuth:         process.env.ZKAUTH_VERIFIER_ADDRESS,
  };
  for (const [k, v] of Object.entries(addr)) {
    if (!v) throw new Error(`missing ${k} in .env (run scripts/deploy_mock.js against the running node first)`);
  }

  const registry       = await ethers.getContractAt("WellRegistry", addr.registry);
  const historyFactory = await ethers.getContractAt("HistoryFactory", addr.historyFactory);
  const zkAuth         = await ethers.getContractAt("ZKAuthVerifier", addr.zkAuth);

  console.log("\n=== Register users ===");

  const adminCommitment = ethers.hexlify(ethers.randomBytes(32));
  const adminInternalID = ethers.hexlify(ethers.randomBytes(32));
  let tx = await registry.connect(deployer).register(adminCommitment, adminInternalID, 3);
  await tx.wait();

  const doctorCommitment = ethers.hexlify(ethers.randomBytes(32));
  const doctorInternalID = ethers.hexlify(ethers.randomBytes(32));
  tx = await registry.connect(deployer).registerByAdmin(doctor.address, doctorCommitment, doctorInternalID, 1);
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

  console.log("\n=== Upload real document ===");
  const uploadResult = await apiUpload(fileData, doctorSession, doctor.address);
  console.log("  IPFS CID:", uploadResult.ipfs_cid);
  console.log("  Content hash:", uploadResult.content_hash);

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

  const onChainRef = await record.storageReference();
  const onChainHash = await record.contentHash();
  console.log("  On-chain storageRef matches:", onChainRef === uploadResult.ipfs_cid ? "YES" : "NO");
  console.log("  On-chain contentHash matches:", onChainHash === uploadResult.content_hash ? "YES" : "NO");
  if (onChainRef !== uploadResult.ipfs_cid || onChainHash !== uploadResult.content_hash) {
    throw new Error("on-chain record metadata does not match uploaded document");
  }

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

  console.log(`\n=== Issue & assign attributes (admin, Vault): ${deptName} / ${instName} ===`);
  const oneYearSec = BigInt(365 * 24 * 3600);
  const validUntil = String(BigInt(Math.floor(Date.now() / 1000)) + oneYearSec);
  const adminHeaders = { "X-User-Address": deployer.address, "X-Session-Token": adminSession };
  const dept = await apiPost("/api/v1/attributes/issue", { kind: "department", name: deptName }, adminHeaders);
  const inst = await apiPost("/api/v1/attributes/issue", { kind: "institution", name: instName }, adminHeaders);
  console.log("  department", deptName, dept.uuid);
  console.log("  institution", instName, inst.uuid);

  await apiPost("/api/v1/attributes/assign", {
    doctor_address: doctor.address.toLowerCase(),
    department: dept.uuid,
    institution: inst.uuid,
    valid_until: validUntil,
  }, adminHeaders);
  console.log("  assigned to", doctor.address);

  const attrsCheck = await apiGet(`/api/v1/attributes/${doctor.address.toLowerCase()}`, doctorSession, doctor.address);
  if (attrsCheck.status !== 200 || attrsCheck.data.department !== dept.uuid.toLowerCase()) {
    throw new Error("assigned attributes not readable via API");
  }
  console.log("  [PASS] attributes stored in Vault and served by API");

  console.log("\n=== Generate policy proof (attributes fetched from Vault via API) ===");
  const proof = await genProof({
    doctor: doctor.address, apiBase: API_BASE, sessionToken: doctorSession,
  });
  console.log("  deptCommitment    ", proof.deptCommitment);
  console.log("  instCommitment    ", proof.instCommitment);
  console.log("  doctorCommitment  ", proof.doctorCommitment);
  console.log("  evaluatedAt       ", parseInt(proof.evaluatedAt, 16), "(" + new Date(parseInt(proof.evaluatedAt, 16) * 1000).toISOString() + ")");

  console.log("\n=== Set policy & doctor self-grants via proof ===");
  tx = await record.connect(patient).setPolicy(proof.deptCommitment, proof.instCommitment);
  await tx.wait();

  tx = await record.connect(doctor).grantAccessWithProof(
    proof.proof.pA, proof.proof.pB, proof.proof.pC,
    proof.deptCommitment, proof.instCommitment, proof.doctorCommitment, proof.evaluatedAt
  );
  const grantReceipt = await tx.wait();
  const grantLog = grantReceipt.logs
    .map(l => { try { return record.interface.parseLog({ topics: [...l.topics], data: l.data }); } catch { return null; } })
    .find(e => e?.name === "AccessGranted");
  if (!grantLog) throw new Error("AccessGranted not emitted on self-grant");
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

  console.log("\n=== Download the document back ===");
  const dlRes = await fetch(
    `${API_BASE}/api/v1/records/${recordAddr.toLowerCase()}/download`,
    { headers: { "X-User-Address": doctor.address, "X-Session-Token": doctorSession } }
  );
  if (!dlRes.ok) throw new Error(`download: ${dlRes.status}`);
  const dlBuf = Buffer.from(await dlRes.arrayBuffer());
  console.log(`  Downloaded ${dlBuf.length} bytes`);

  const dlHash = sha256hex(dlBuf);

  console.log("\n=== Saving fetched document to output directory ===");
  mkdirSync(OUTPUT_DIR, { recursive: true });
  const baseName = path.basename(filePath);
  const originalOut = path.join(OUTPUT_DIR, `original-${baseName}`);
  const downloadedOut = path.join(OUTPUT_DIR, `downloaded-${baseName}`);
  writeFileSync(originalOut, fileData);
  writeFileSync(downloadedOut, dlBuf);
  const sumsPath = path.join(OUTPUT_DIR, "SHA256SUMS");
  writeFileSync(sumsPath, `${fileHash}  original-${baseName}\n${dlHash}  downloaded-${baseName}\n`);
  console.log("  Output dir:    " + OUTPUT_DIR);
  console.log("  original-...:  " + originalOut);
  console.log("  downloaded-...:" + downloadedOut);

  if (dlHash !== fileHash) {
    throw new Error(`downloaded document does NOT match the original (got 0x${dlHash}, want 0x${fileHash})`);
  }
  if (uploadResult.content_hash !== "0x" + fileHash) {
    throw new Error("upload content_hash does not match sha256 of the source file");
  }
  console.log("  [PASS] sha256(downloaded) == sha256(original) == " + uploadResult.content_hash);
  console.log(`  [PASS] hashes match. Verify anytime with:  cd ${OUTPUT_DIR} && sha256sum -c SHA256SUMS`);

  console.log("\n=== Manual document test passed! ===");
  console.log("\nSummary (reuse for manual curl):");
  console.log(`  RECORD=${recordAddr}`);
  console.log(`  DOCTOR=${doctor.address}`);
  console.log(`  SESSION_TOKEN=${doctorSession}`);
  console.log(`  IPFS_CID=${uploadResult.ipfs_cid}`);
  console.log(`  IPFS_GATEWAY_URL=http://localhost:8080/ipfs/${uploadResult.ipfs_cid}`);
  console.log(`  OUTPUT_DIR=${OUTPUT_DIR}  (cd there && sha256sum -c SHA256SUMS)`);
}

function genProof({ doctor, apiBase, sessionToken }) {
  const args = [
    "--api-base", apiBase,
    "--user-address", doctor,
    "--session-token", sessionToken,
  ];
  const out = execFileSync("node", [GEN_PROOF, ...args], {
    encoding: "utf8", maxBuffer: 1024 * 1024,
    stdio: ["ignore", "pipe", "pipe"],
  });
  return JSON.parse(out);
}

main().catch(err => {
  console.error("\nManual document test failed:", err.message);
  process.exit(1);
});
