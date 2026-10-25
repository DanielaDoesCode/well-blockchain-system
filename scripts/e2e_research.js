import "dotenv/config";
import hre from "hardhat";
import { execFileSync } from "child_process";
import path from "path";
import { fileURLToPath } from "url";

const { ethers } = hre;

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const GEN_PROOF = path.join(__dirname, "..", "zk-circuits", "gen_merkle_proof.js");

function genProof(extraArgs = []) {
  const out = execFileSync("node", [GEN_PROOF, ...extraArgs], {
    encoding: "utf8", maxBuffer: 1024 * 1024,
    stdio: ["ignore", "pipe", "pipe"],
  });
  return JSON.parse(out);
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
  const [deployer, doctor, patient, patientNoConsent] = await ethers.getSigners();

  const addr = {
    registry:          process.env.REGISTRY_ADDRESS,
    userRegistry:      process.env.USER_REGISTRY_ADDRESS,
    attestationManager: process.env.ATTENTATION_MANAGER_ADDRESS,
    researchVerifier:  process.env.RESEARCH_VERIFIER_ADDRESS,
    researchContributer: process.env.RESEARCH_CONTRIBUTER_ADDRESS,
  };
  for (const [k, v] of Object.entries(addr)) {
    if (!v) throw new Error(`missing ${k} in .env`);
  }

  const registry        = await ethers.getContractAt("WellRegistry", addr.registry);
  const userReg         = await ethers.getContractAt("UserRegistry", addr.userRegistry);
  const attestation     = await ethers.getContractAt("AttestationManager", addr.attestationManager);
  const researchVerifier = await ethers.getContractAt("ResearchVerifier", addr.researchVerifier);
  const contributer     = await ethers.getContractAt("ResearchContributer", addr.researchContributer);

  console.log("=== Generate reference proof (default inputs) ===");
  const proof = genProof();
  const commitment = proof.commitment;//Poseidon(secret, 1)
  const nullifier = proof.nullifier;//Poseidon(secret, ns, 42)
  const merkleRoot = proof.merkleRoot;//tree root for leaf Poseidon(commitment, 42)
  const cg = BigInt(proof.conditionGroup);
  console.log("  commitment :", commitment);
  console.log("  nullifier  :", nullifier);
  console.log("  merkleRoot :", merkleRoot);
  console.log("  condition  :", cg.toString());

  console.log("\n=== Register users ===");
  const adminCommitment = ethers.hexlify(ethers.randomBytes(32));
  const adminInternalID = ethers.hexlify(ethers.randomBytes(32));
  let tx = await registry.connect(deployer).register(adminCommitment, adminInternalID, 3);
  await tx.wait();

  const doctorCommitment = ethers.hexlify(ethers.randomBytes(32));
  const doctorInternalID = ethers.hexlify(ethers.randomBytes(32));
  tx = await registry.connect(deployer).registerByAdmin(doctor.address, doctorCommitment, doctorInternalID, 1);
  await tx.wait();
  console.log("  doctor registered via admin");

  tx = await registry.connect(patient).register(commitment, ethers.hexlify(ethers.randomBytes(32)), 0);
  await tx.wait();
  console.log("  patient self-registered with proof commitment");
  if (!(await userReg.getCommitmentByAddress(patient.address)) === commitment) {
    throw new Error("patient commitment mismatch after register");
  }

  const otherCommitment = ethers.hexlify(ethers.randomBytes(32));
  tx = await registry.connect(patientNoConsent).register(otherCommitment, ethers.hexlify(ethers.randomBytes(32)), 0);
  await tx.wait();
  console.log("  no-consent patient registered");

  console.log("\n=== Negative: attest without consent ===");
  await expectRevert("attest-without-consent rejected",
    () => attestation.connect(doctor).attest(otherCommitment, ethers.toBeHex(cg, 32)),
    "consent required");

  console.log("\n=== Patient opts in, doctor attests ===");
  tx = await attestation.connect(patient).optIn(commitment);
  await tx.wait();
  console.log("  patient consented");

  tx = await attestation.connect(doctor).attest(commitment, ethers.toBeHex(cg, 32));
  await tx.wait();

  const root = await attestation.getRoot();
  if (root !== BigInt(merkleRoot)) {
    throw new Error(`tree root mismatch: contract ${root} vs proof ${merkleRoot}`);
  }
  if (!(await attestation.isKnownRoot(merkleRoot))) {
    throw new Error("attestation root not known after insert");
  }
  console.log("  [PASS] on-chain tree root == proof merkleRoot", merkleRoot);

  console.log("\n=== On-chain proof verification (direct) ===");
  const pubSignals = [BigInt(commitment), BigInt(nullifier), BigInt(merkleRoot)];
  if (!(await researchVerifier.verifyProof(proof.proof.pA, proof.proof.pB, proof.proof.pC, pubSignals))) {
    throw new Error("ResearchVerifier rejected the reference proof");
  }
  console.log("  [PASS] ResearchVerifier.verifyProof == true");

  console.log("\n=== Contribute happy path ===");
  tx = await contributer.connect(patient).contribute(
    proof.proof.pA, proof.proof.pB, proof.proof.pC,
    commitment, nullifier, BigInt(merkleRoot)
  );
  await tx.wait();
  if (!(await contributer.usedNullifiers(nullifier))) {
    throw new Error("nullifier not marked used after contribute");
  }
  console.log("  [PASS] contribution accepted, nullifier consumed");

  console.log("\n=== Negative: nullifier reuse ===");
  await expectRevert("nullifier-reuse rejected",
    () => contributer.connect(patient).contribute(
      proof.proof.pA, proof.proof.pB, proof.proof.pC,
      commitment, nullifier, BigInt(merkleRoot)
    ),
    "NullifierAlreadyUsed");

  console.log("\n=== Negative: tampered proof ===");
  const proof2 = genProof(["--nullifierSeed", "111111111"]);
  if (proof2.commitment !== commitment || proof2.merkleRoot !== merkleRoot) {
    throw new Error("second proof should share commitment+root (same secret/conditionGroup)");
  }
  const tampered = JSON.parse(JSON.stringify(proof2.proof));
  tampered.pA[0] = (BigInt(tampered.pA[0]) + 1n).toString();
  await expectRevert("tampered-proof rejected",
    () => contributer.connect(patient).contribute(
      tampered.pA, tampered.pB, tampered.pC,
      commitment, proof2.nullifier, BigInt(merkleRoot)
    ),
    "InvalidProof");

  console.log("\n=== Negative: stale / invalid root ===");
  await expectRevert("invalid-root rejected",
    () => contributer.connect(patient).contribute(
      proof2.proof.pA, proof2.proof.pB, proof2.proof.pC,
      commitment, proof2.nullifier, 1n
    ),
    "StaleOrInvalidRoot");

  console.log("\n=== Negative: commitment not registered ===");
  const unknown = ethers.toBeHex(BigInt(commitment) ^ 1n, 32);
  await expectRevert("unknown-commitment rejected",
    () => contributer.connect(patient).contribute(
      proof2.proof.pA, proof2.proof.pB, proof2.proof.pC,
      unknown, proof2.nullifier, BigInt(merkleRoot)
    ),
    "CommitmentNotRegistered");

  console.log("\n=== E2E research test passed! ===");
}

main().catch(err => {
  console.error("\nE2E research test failed:", err.message);
  process.exit(1);
});
