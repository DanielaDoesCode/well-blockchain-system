import "dotenv/config";
import hre from "hardhat";

const { ethers } = hre;

const API_BASE = `http://localhost:${process.env.AGGREGATOR_PORT || 8081}`;

async function apiPost(path, body, sessionToken, userAddress) {
  const res = await fetch(`${API_BASE}${path}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-User-Address": userAddress,
      "X-Session-Token": sessionToken,
    },
    body: JSON.stringify(body),
  });
  const data = await res.json();
  return { status: res.status, data };
}

async function apiGet(path, sessionToken, userAddress) {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { "X-User-Address": userAddress, "X-Session-Token": sessionToken },
  });
  const data = await res.json();
  return { status: res.status, data };
}

async function apiPatch(path, body, sessionToken, userAddress) {
  const res = await fetch(`${API_BASE}${path}`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
      "X-User-Address": userAddress,
      "X-Session-Token": sessionToken,
    },
    body: JSON.stringify(body),
  });
  const data = await res.json();
  return { status: res.status, data };
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
  const authLog = await parseEvent(receipt, zkAuth.interface, "Authenticated");
  if (!authLog) throw new Error("could not find Authenticated event for " + signer.address);
  return authLog.args.sessionToken;
}

async function expectRevert(label, fn, selector, contract) {
  try {
    await fn();
  } catch (e) {
    let reason = e.reason ?? "";
    // ethers v6 decodes custom errors into e.revert when it has the ABI.
    if (!reason && e.revert?.name) reason = e.revert.name;
    if (!reason) {
      const data = e.data ?? e.error?.data ?? e.info?.error?.data;
      if (data && contract?.interface) {
        try { reason = contract.interface.parseError(data)?.name ?? ""; } catch { /* keep reason empty */ }
      }
    }
    if (!reason && e.message) reason = e.message;
    if (selector && !reason.includes(selector)) {
      throw new Error(`${label}: reverted with unexpected reason "${reason}" (wanted ${selector})`);
    }
    console.log(`  [PASS] ${label}`);
    return;
  }
  throw new Error(`${label}: expected revert, but transaction succeeded`);
}

async function parseEvent(receipt, iface, name) {
  const log = receipt.logs
    .map(l => { try { return iface.parseLog({ topics: [...l.topics], data: l.data }); } catch { return null; } })
    .find(e => e?.name === name);
  return log;
}

async function main() {
  const [deployer, hospital, responder, otherResponder, doctor, patient, patient2] = await ethers.getSigners();

  const addr = {
    registry:         process.env.REGISTRY_ADDRESS,
    incidentManager:  process.env.INCIDENT_MANAGER_ADDRESS,
    riskAssessment:   process.env.RISK_ASSESSMENT_ADDRESS,
    emergencyFactory: process.env.EMERGENCY_SUMMARY_FACTORY_ADDRESS,
    zkAuth:           process.env.ZKAUTH_VERIFIER_ADDRESS,
  };
  for (const [k, v] of Object.entries(addr)) {
    if (!v) throw new Error(`missing ${k} in .env`);
  }

  const registry   = await ethers.getContractAt("WellRegistry", addr.registry);
  const incident   = await ethers.getContractAt("IncidentManager", addr.incidentManager);
  const riskAssess = await ethers.getContractAt("RiskAssessment", addr.riskAssessment);
  const esFactory  = await ethers.getContractAt("EmergencySummaryFactory", addr.emergencyFactory);
  const zkAuth     = await ethers.getContractAt("ZKAuthVerifier", addr.zkAuth);

  console.log("=== Register users ===");

  let tx = await registry.connect(deployer).register(
    ethers.hexlify(ethers.randomBytes(32)), ethers.hexlify(ethers.randomBytes(32)), 3
  );
  await tx.wait();
  console.log("  Admin registered");

  tx = await registry.connect(deployer).registerByAdmin(hospital.address,
    ethers.hexlify(ethers.randomBytes(32)), ethers.hexlify(ethers.randomBytes(32)), 4);
  await tx.wait();
  console.log("  Hospital registered");

  tx = await registry.connect(deployer).registerByAdmin(responder.address,
    ethers.hexlify(ethers.randomBytes(32)), ethers.hexlify(ethers.randomBytes(32)), 2);
  await tx.wait();
  tx = await registry.connect(deployer).registerByAdmin(otherResponder.address,
    ethers.hexlify(ethers.randomBytes(32)), ethers.hexlify(ethers.randomBytes(32)), 2);
  await tx.wait();
  console.log("  2 FirstResponders registered");

  tx = await registry.connect(deployer).registerByAdmin(doctor.address,
    ethers.hexlify(ethers.randomBytes(32)), ethers.hexlify(ethers.randomBytes(32)), 1);
  await tx.wait();
  console.log("  Doctor registered");

  const patientCommitment = ethers.hexlify(ethers.randomBytes(32));
  tx = await registry.connect(patient).register(
    patientCommitment, ethers.hexlify(ethers.randomBytes(32)), 0
  );
  await tx.wait();
  console.log("  Patient registered (EmergencySummary auto-created)");

  const summaryAddr = await registry.getEmergencySummaryByAddress(patient.address);
  if (summaryAddr === ethers.ZeroAddress) throw new Error("no EmergencySummary created for patient");
  const emergency = await ethers.getContractAt("EmergencySummary", summaryAddr);
  console.log("  EmergencySummary:", summaryAddr);

  const byAddr = await esFactory.getEmergencySummaryByAddress(patient.address);
  const byCommitment = await esFactory.getEmergencySummaryByCommitment(patientCommitment);
  if (byAddr !== summaryAddr || byCommitment !== summaryAddr) {
    throw new Error("EmergencySummaryFactory mapping mismatch");
  }
  console.log("  [PASS] factory mapping resolves patient -> summary");

  const patient2Commitment = ethers.hexlify(ethers.randomBytes(32));
  tx = await registry.connect(patient2).register(
    patient2Commitment, ethers.hexlify(ethers.randomBytes(32)), 0
  );
  await tx.wait();
  const patient2Summary = await registry.getEmergencySummaryByAddress(patient2.address);
  if (patient2Summary === ethers.ZeroAddress) throw new Error("no EmergencySummary created for patient2");
  console.log("  Patient2 registered (for cross-patient negative case):", patient2Summary);

  console.log("\n=== ZK auth sessions (for recovery-key flow) ===");
  const patientSession = await zkAuthSession(zkAuth, patient);
  const responderSession = await zkAuthSession(zkAuth, responder);
  const otherSession = await zkAuthSession(zkAuth, otherResponder);
  console.log("  sessions minted; waiting for listener to persist them (3s)...");
  await new Promise(r => setTimeout(r, 3000));

  console.log("\n=== Off-chain identity index: provision + resolve ===");
  const adminSession = await zkAuthSession(zkAuth, deployer);
  const PATIENT_NID = "PAT-2024-0001";
  const PATIENT2_NID = "PAT-2024-0002";

  let prov = await apiPost("/api/v1/identities",
    { national_id: PATIENT_NID, display_name: "Alice Patient", role: "patient",
      address: patient.address, summary_address: summaryAddr,
      commitment: patientCommitment,
      internal_id: ethers.hexlify(ethers.randomBytes(32)) },
    adminSession, deployer.address);
  if (prov.status !== 201) throw new Error("provision patient failed: " + JSON.stringify(prov));

  prov = await apiPost("/api/v1/identities",
    { national_id: PATIENT2_NID, display_name: "Bob Patient", role: "patient",
      address: patient2.address, summary_address: patient2Summary,
      commitment: patient2Commitment,
      internal_id: ethers.hexlify(ethers.randomBytes(32)) },
    adminSession, deployer.address);
  if (prov.status !== 201) throw new Error("provision patient2 failed: " + JSON.stringify(prov));
  console.log("  [PASS] admin provisioned 2 patient identities in the off-chain index");

  let lookup = await apiGet(`/api/v1/identities?national_id=${PATIENT_NID}`, responderSession, responder.address);
  if (lookup.status !== 200) throw new Error("lookup by national id failed: " + JSON.stringify(lookup));
  const resolvedPatient = lookup.data.identity;
  if (resolvedPatient.address.toLowerCase() !== patient.address.toLowerCase() ||
      resolvedPatient.summary_address.toLowerCase() !== summaryAddr.toLowerCase() ||
      resolvedPatient.commitment !== patientCommitment) {
    throw new Error("resolved identity mismatch: " + JSON.stringify(resolvedPatient));
  }
  console.log("  [PASS] responder resolved patient by national ID -> address/summary/commitment");
  const resolvedSummary = resolvedPatient.summary_address;
  const resolvedCommitment = resolvedPatient.commitment;

  lookup = await apiGet(`/api/v1/identities?name=${encodeURIComponent("Bob Patient")}`, responderSession, responder.address);
  if (lookup.status !== 200) throw new Error("lookup by name failed: " + JSON.stringify(lookup));
  const resolvedPatient2 = lookup.data.identity;
  if (resolvedPatient2.address.toLowerCase() !== patient2.address.toLowerCase()) {
    throw new Error("resolved patient2 mismatch: " + JSON.stringify(resolvedPatient2));
  }
  console.log("  [PASS] responder resolved patient2 by display name");
  const resolvedSummary2 = resolvedPatient2.summary_address;
  const resolvedCommitment2 = resolvedPatient2.commitment;

  lookup = await apiGet("/api/v1/identities?national_id=UNKNOWN-NID", responderSession, responder.address);
  if (lookup.status !== 404) throw new Error("expected 404 for unknown national id, got " + lookup.status);
  console.log("  [PASS] unknown national id -> 404");

  lookup = await apiGet(`/api/v1/identities?national_id=${PATIENT_NID}`, patientSession, patient.address);
  if (lookup.status !== 403) throw new Error("expected 403 for patient lookup, got " + lookup.status);
  console.log("  [PASS] patient role blocked from index (403)");

  console.log("\n=== Hospital creates incident ===");
  const now = Math.floor(Date.now() / 1000);
  const expiresAt = now + 3600;
  const locationHash = ethers.hexlify(ethers.randomBytes(32));
  tx = await incident.connect(hospital).createIncident(locationHash, expiresAt);
  const incRcpt = await tx.wait();
  const incLog = await parseEvent(incRcpt, incident.interface, "IncidentCreated");
  if (!incLog) throw new Error("IncidentCreated not emitted");
  const incidentId = incLog.args.id;
  if (!(await incident.isIncidentActive(incidentId))) throw new Error("incident not active");
  console.log("  [PASS] IncidentCreated (expires " + incLog.args.expiresAt.toString() + ")");

  console.log("\n=== Hospital adds responder ===");
  tx = await incident.connect(hospital).addResponders(incidentId, [responder.address]);
  const addRcpt = await tx.wait();
  const addLog = await parseEvent(addRcpt, incident.interface, "RespondersAdded");
  if (!addLog || !addLog.args.responders.includes(responder.address)) {
    throw new Error("RespondersAdded mismatch");
  }
  if (!(await incident.isResponderInIncident(incidentId, responder.address))) {
    throw new Error("responder not assigned to incident");
  }
  console.log("  [PASS] responder assigned to incident");

  console.log("\n=== Patient uploads encrypted emergency summary ===");
  const encryptedData = ethers.hexlify(ethers.randomBytes(64));
  const dataHash = ethers.keccak256(ethers.randomBytes(64));
  tx = await emergency.connect(patient).updateData(encryptedData, dataHash);
  const upRcpt = await tx.wait();
  const upLog = await parseEvent(upRcpt, emergency.interface, "DataUpdated");
  if (!upLog || upLog.args.dataHash !== dataHash) throw new Error("DataUpdated mismatch");
  console.log("  [PASS] DataUpdated (version=" + upLog.args.version.toString() + ")");

  await expectRevert("non-owner updateData rejected",
    () => emergency.connect(responder).updateData(encryptedData, dataHash),
    "Unauthorized", emergency);

  console.log("\n=== Responder registers public key + requests access ===");
  const pubKey = ethers.hexlify(ethers.randomBytes(32));
  tx = await emergency.connect(responder).registerPublicKey(pubKey);
  const pkrRcpt = await tx.wait();
  const pkrLog = await parseEvent(pkrRcpt, emergency.interface, "PublicKeyRegistered");
  if (!pkrLog || pkrLog.args.user !== responder.address) throw new Error("PublicKeyRegistered mismatch");
  console.log("  [PASS] PublicKeyRegistered");

  tx = await emergency.connect(responder).requestAccess(
    pubKey, "Street 1", "Cardiac arrest", "patient-id-001"
  );
  const reqRcpt = await tx.wait();
  const reqLog = await parseEvent(reqRcpt, emergency.interface, "AccessRequested");
  if (!reqLog || reqLog.args.responder !== responder.address) throw new Error("AccessRequested mismatch");
  console.log("  [PASS] AccessRequested (requestId=" + reqLog.args.requestId.toString() + ")");

  console.log("\n=== Patient approves access request ===");
  const encryptedKey = ethers.hexlify(ethers.randomBytes(32));
  await expectRevert("non-owner approveAccessRequest rejected",
    () => emergency.connect(responder).approveAccessRequest(reqLog.args.requestId, encryptedKey),
    "Unauthorized", emergency);

  tx = await emergency.connect(patient).approveAccessRequest(reqLog.args.requestId, encryptedKey);
  const apprRcpt = await tx.wait();
  const apprLog = await parseEvent(apprRcpt, emergency.interface, "AccessGranted");
  if (!apprLog || apprLog.args.user !== responder.address) throw new Error("AccessGranted mismatch");
  if (!(await emergency.hasAccess(responder.address))) {
    throw new Error("responder has no access after approval");
  }
  console.log("  [PASS] AccessGranted to responder (expires " + apprLog.args.expiresAt.toString() + ")");

  console.log("\n=== Responder reads encrypted summary ===");
  const [data, key, version, hash] = await emergency.connect(responder).getData();
  if (ethers.hexlify(data) !== encryptedData || hash !== dataHash) throw new Error("getData mismatch");
  if (key.length === 0) throw new Error("responder has no decryption key");
  console.log("  [PASS] responder retrieved encrypted summary (version=" + version.toString() + ")");

  await expectRevert("unauthorized getData rejected",
    () => emergency.connect(doctor).getData(),
    "Unauthorized", emergency);

  console.log("\n=== Owner direct grant + modifier enforcement ===");
  tx = await emergency.connect(patient).grantAccess(
    doctor.address, ethers.hexlify(ethers.randomBytes(32)), 3600
  );
  await tx.wait();
  if (!(await emergency.hasAccess(doctor.address))) throw new Error("direct grant failed");
  console.log("  [PASS] owner grantAccess works");

  await expectRevert("arbitrary grantAccess rejected (modifier fix)",
    () => emergency.connect(doctor).grantAccess(
      responder.address, ethers.hexlify(ethers.randomBytes(32)), 3600
    ),
    "Unauthorized", emergency);

  console.log("\n=== RiskAssessment emergencyGrant (assigned first responder) ===");
  tx = await riskAssess.connect(responder).emergencyGrant(
    resolvedSummary, incidentId, locationHash, resolvedCommitment,
    ethers.hexlify(ethers.randomBytes(32))
  );
  const egRcpt = await tx.wait();
  const egLog = await parseEvent(egRcpt, riskAssess.interface, "EmergencyGranted");
  if (!egLog || egLog.args.responder !== responder.address || egLog.args.incidentId !== incidentId) {
    throw new Error("EmergencyGranted mismatch");
  }
  console.log("  [PASS] RiskAssessment.EmergencyGranted");

  console.log("\n=== Negative: wrong locationHash emergencyGrant rejected ===");
  await expectRevert("wrong locationHash rejected",
    () => riskAssess.connect(responder).emergencyGrant(
      resolvedSummary, incidentId, ethers.hexlify(ethers.randomBytes(32)), resolvedCommitment,
      ethers.hexlify(ethers.randomBytes(32))
    ),
    "LocationMismatch", riskAssess);

  console.log("\n=== Negative: summary not matching commitment rejected ===");
  await expectRevert("summary not matching commitment rejected",
    () => riskAssess.connect(responder).emergencyGrant(
      resolvedSummary, incidentId, locationHash, ethers.hexlify(ethers.randomBytes(32)),
      ethers.hexlify(ethers.randomBytes(32))
    ),
    "InvalidPatient", riskAssess);

  console.log("\n=== Negative: cross-patient grant on same incident rejected ===");
  await expectRevert("cross-patient grant rejected",
    () => riskAssess.connect(responder).emergencyGrant(
      resolvedSummary2, incidentId, locationHash, resolvedCommitment2,
      ethers.hexlify(ethers.randomBytes(32))
    ),
    "PatientMismatch", riskAssess);

  console.log("\n=== Negative: unassigned responder emergencyGrant rejected ===");
  await expectRevert("unassigned responder rejected",
    () => riskAssess.connect(otherResponder).emergencyGrant(
      resolvedSummary, incidentId, locationHash, resolvedCommitment,
      ethers.hexlify(ethers.randomBytes(32))
    ),
    "NotAssignedResponder", riskAssess);

  console.log("\n=== Negative: non-first-responder emergencyGrant rejected ===");
  await expectRevert("doctor emergencyGrant rejected",
    () => riskAssess.connect(doctor).emergencyGrant(
      resolvedSummary, incidentId, locationHash, resolvedCommitment,
      ethers.hexlify(ethers.randomBytes(32))
    ),
    "NotFirstResponder", riskAssess);

  console.log("\n=== Patient denies a second access request ===");
  tx = await emergency.connect(responder).requestAccess(
    pubKey, "Street 2", "Trauma", "patient-id-001"
  );
  const req2Rcpt = await tx.wait();
  const req2Log = await parseEvent(req2Rcpt, emergency.interface, "AccessRequested");
  if (!req2Log) throw new Error("second AccessRequested not emitted");

  tx = await emergency.connect(patient).denyAccessRequest(req2Log.args.requestId);
  await tx.wait();
  const req2Info = await emergency.connect(responder).getAccessRequest(req2Log.args.requestId);
  if (Number(req2Info.status) !== 2) throw new Error("expected Denied status, got " + req2Info.status);
  console.log("  [PASS] request denied");

  console.log("\n=== Negative: emergencyGrant on deactivated incident rejected ===");
  tx = await incident.connect(hospital).deactivateIncident(incidentId);
  await tx.wait();
  await expectRevert("inactive incident rejected",
    () => riskAssess.connect(responder).emergencyGrant(
      resolvedSummary, incidentId, locationHash, resolvedCommitment,
      ethers.hexlify(ethers.randomBytes(32))
    ),
    "IncidentNotActive", riskAssess);

  console.log("\n=== Recovery key escrow (hybrid: Vault transit wrap + on-chain attestation) ===");
  const recoveryKey = ethers.hexlify(ethers.randomBytes(32));

  let store = await apiPost("/api/v1/emergency/recovery-key",
    { emergency_summary: summaryAddr, content_key: recoveryKey },
    patientSession, patient.address);
  if (store.status !== 200 || !store.data.recovery_key) {
    throw new Error("owner store recovery key failed: " + JSON.stringify(store));
  }
  const wrappedHex = store.data.recovery_key;
  console.log("  [PASS] owner escrowed recovery key via backend (wrapped len=" + wrappedHex.length + ")");

  let forbidden = await apiPost("/api/v1/emergency/recovery-key",
    { emergency_summary: summaryAddr, content_key: recoveryKey },
    responderSession, responder.address);
  if (forbidden.status !== 403) throw new Error("expected 403 for non-owner store, got " + forbidden.status);
  console.log("  [PASS] non-owner store rejected (403)");

  let got = await apiGet(`/api/v1/emergency/recovery-key?emergency_summary=${summaryAddr.toLowerCase()}`,
    patientSession, patient.address);
  if (got.status !== 200 || got.data.recovery_key !== wrappedHex) {
    throw new Error("owner read mismatch: " + JSON.stringify(got));
  }
  console.log("  [PASS] owner read back identical wrapped key");

  forbidden = await apiGet(`/api/v1/emergency/recovery-key?emergency_summary=${summaryAddr.toLowerCase()}`,
    responderSession, responder.address);
  if (forbidden.status !== 403) throw new Error("expected 403 for non-owner read, got " + forbidden.status);
  console.log("  [PASS] non-owner read rejected (403)");

  console.log("\n=== On-chain recovery key attestation ===");
  tx = await emergency.connect(patient).storeRecoveryKey("0x" + wrappedHex);
  const rkRcpt = await tx.wait();
  const rkLog = await parseEvent(rkRcpt, emergency.interface, "RecoveryKeyStored");
  if (!rkLog) throw new Error("RecoveryKeyStored not emitted");
  if ((await emergency.recoveryVersion()).toString() !== "1") throw new Error("recoveryVersion != 1");
  const chainKey = await emergency.connect(patient).getRecoveryKey();
  if (chainKey !== "0x" + wrappedHex) throw new Error("on-chain recovery key mismatch");
  console.log("  [PASS] RecoveryKeyStored (version=1), on-chain read matches backend blob");

  await expectRevert("non-owner on-chain store rejected",
    () => emergency.connect(responder).storeRecoveryKey("0x" + wrappedHex),
    "Unauthorized", emergency);
  await expectRevert("non-owner on-chain read rejected",
    () => emergency.connect(responder).getRecoveryKey(),
    "Unauthorized", emergency);

  console.log("\n=== Re-encrypt recovery key for granted responder ===");
  let re = await apiPost("/api/v1/emergency/recovery-key/re-encrypt",
    { emergency_summary: summaryAddr, recipient: responder.address },
    responderSession, responder.address);
  if (re.status !== 200 || !re.data.reencrypted_key) {
    throw new Error("granted re-encrypt failed: " + JSON.stringify(re));
  }
  console.log("  [PASS] granted responder re-encrypted recovery key for recipient");

  re = await apiPost("/api/v1/emergency/recovery-key/re-encrypt",
    { emergency_summary: summaryAddr, recipient: responder.address },
    otherSession, otherResponder.address);
  if (re.status !== 403) throw new Error("expected 403 for unprivileged re-encrypt, got " + re.status);
  console.log("  [PASS] unprivileged responder re-encrypt rejected (403)");

  console.log("\n=== E2E emergency test passed! ===");
}

main().catch(err => {
  console.error("\nE2E emergency test failed:", err.message);
  process.exit(1);
});
