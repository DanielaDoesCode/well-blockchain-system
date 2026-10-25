import hre from "hardhat";
import { writeFileSync } from "fs";

const { ethers } = hre;

async function main() {
  const [deployer] = await ethers.getSigners();
  console.log("Deploying with account:", deployer.address);

  // ── Auth layer ──────────────────────────────────
  const Groth16Verifier = await ethers.getContractFactory("Groth16Verifier");
  const groth16Verifier = await Groth16Verifier.deploy();
  await groth16Verifier.waitForDeployment();
  console.log("Groth16Verifier deployed at:", await groth16Verifier.getAddress());

  const ZKAuthVerifier = await ethers.getContractFactory("ZKAuthVerifier");
  const zkVerifier = await ZKAuthVerifier.deploy(
    await groth16Verifier.getAddress(),
    3600 // 1 hour session duration
  );
  await zkVerifier.waitForDeployment();
  console.log("ZKAuthVerifier deployed at:", await zkVerifier.getAddress());

  // ── Identity / registry layer ───────────────────
  const UserRegistry = await ethers.getContractFactory("UserRegistry");
  const userRegistry = await UserRegistry.deploy();
  await userRegistry.waitForDeployment();
  console.log("UserRegistry deployed at:", await userRegistry.getAddress());

  const PoseidonT3 = await ethers.getContractFactory("PoseidonT3");
  const poseidonT3 = await PoseidonT3.deploy();
  await poseidonT3.waitForDeployment();
  console.log("PoseidonT3 deployed at:", await poseidonT3.getAddress());

  const PolicyVerifier = await ethers.getContractFactory("PolicyVerifier");
  const policyVerifier = await PolicyVerifier.deploy();
  await policyVerifier.waitForDeployment();
  console.log("PolicyVerifier deployed at:", await policyVerifier.getAddress());

  const HistoryFactory = await ethers.getContractFactory("HistoryFactory", {
    libraries: { PoseidonT3: await poseidonT3.getAddress() },
  });
  const historyFactory = await HistoryFactory.deploy(
    await userRegistry.getAddress()
  );
  await historyFactory.waitForDeployment();
  console.log("HistoryFactory deployed at:", await historyFactory.getAddress());

  const WellRegistry = await ethers.getContractFactory("WellRegistry");
  const wellRegistry = await WellRegistry.deploy(
    await userRegistry.getAddress(),
    await historyFactory.getAddress()
  );
  await wellRegistry.waitForDeployment();
  console.log("WellRegistry deployed at:", await wellRegistry.getAddress());

  // Wire HistoryFactory -> WellRegistry + PolicyVerifier
  let tx = await historyFactory.setWellRegistry(await wellRegistry.getAddress());
  await tx.wait();
  tx = await historyFactory.setPolicyVerifier(await policyVerifier.getAddress());
  await tx.wait();
  console.log("HistoryFactory.wellRegistry + policyVerifier set");

  // ── Emergency infrastructure ────────────────────
  const IncidentManager = await ethers.getContractFactory("IncidentManager");
  const incidentManager = await IncidentManager.deploy(
    await userRegistry.getAddress()
  );
  await incidentManager.waitForDeployment();
  console.log("IncidentManager deployed at:", await incidentManager.getAddress());

  const EmergencySummaryFactory = await ethers.getContractFactory("EmergencySummaryFactory");
  const emergencySummaryFactory = await EmergencySummaryFactory.deploy(
    await userRegistry.getAddress()
  );
  await emergencySummaryFactory.waitForDeployment();
  console.log("EmergencySummaryFactory deployed at:", await emergencySummaryFactory.getAddress());

  const RiskAssessment = await ethers.getContractFactory("RiskAssessment");
  const riskAssessment = await RiskAssessment.deploy(
    await incidentManager.getAddress(),
    await userRegistry.getAddress(),
    await emergencySummaryFactory.getAddress()
  );
  await riskAssessment.waitForDeployment();
  console.log("RiskAssessment deployed at:", await riskAssessment.getAddress());

  // Wire emergency contracts into WellRegistry
  tx = await wellRegistry.setIncidentManager(await incidentManager.getAddress());
  await tx.wait();
  tx = await wellRegistry.setEmergencySummaryFactory(await emergencySummaryFactory.getAddress());
  await tx.wait();
  tx = await wellRegistry.setRiskAssessment(await riskAssessment.getAddress());
  await tx.wait();
  console.log("WellRegistry emergency setters wired");

  // Wire EmergencySummaryFactory
  tx = await emergencySummaryFactory.setWellRegistry(await wellRegistry.getAddress());
  await tx.wait();
  tx = await emergencySummaryFactory.setRiskAssessment(await riskAssessment.getAddress());
  await tx.wait();
  console.log("EmergencySummaryFactory wired");

  // ── Attestation layer ───────────────────────────
  const AttestationManager = await ethers.getContractFactory("AttestationManager", {
    libraries: { PoseidonT3: await poseidonT3.getAddress() },
  });
  const attestationManager = await AttestationManager.deploy(
    await userRegistry.getAddress()
  );
  await attestationManager.waitForDeployment();
  console.log("AttestationManager deployed at:", await attestationManager.getAddress());

  // ── Research layer ──────────────────────────────
  const ResearchVerifier = await ethers.getContractFactory("ResearchVerifier");
  const researchVerifier = await ResearchVerifier.deploy();
  await researchVerifier.waitForDeployment();
  console.log("ResearchVerifier deployed at:", await researchVerifier.getAddress());

  const ResearchContributer = await ethers.getContractFactory("ResearchContributer");
  const researchContributer = await ResearchContributer.deploy(
    await researchVerifier.getAddress(),
    await wellRegistry.getAddress(),
    await attestationManager.getAddress()
  );
  await researchContributer.waitForDeployment();
  console.log("ResearchContributer deployed at:", await researchContributer.getAddress());

  // ── Log integrity layer ──────────────────────────
  const LogAnchor = await ethers.getContractFactory("LogAnchor");
  const logAnchor = await LogAnchor.deploy();
  await logAnchor.waitForDeployment();
  console.log("LogAnchor deployed at:", await logAnchor.getAddress());

  // ── Summary ─────────────────────────────────────
  const addrs = {
    REGISTRY_ADDRESS: await wellRegistry.getAddress(),
    USER_REGISTRY_ADDRESS: await userRegistry.getAddress(),
    HISTORY_FACTORY_ADDRESS: await historyFactory.getAddress(),
    ZKAUTH_VERIFIER_ADDRESS: await zkVerifier.getAddress(),
    INCIDENT_MANAGER_ADDRESS: await incidentManager.getAddress(),
    RISK_ASSESSMENT_ADDRESS: await riskAssessment.getAddress(),
    EMERGENCY_SUMMARY_ADDRESS: await emergencySummaryFactory.getAddress(),
    EMERGENCY_SUMMARY_FACTORY_ADDRESS: await emergencySummaryFactory.getAddress(),
    RESEARCH_CONTRIBUTER_ADDRESS: await researchContributer.getAddress(),
    POSEIDON_T3_ADDRESS: await poseidonT3.getAddress(),
    LOG_ANCHOR_ADDRESS: await logAnchor.getAddress(),
    // Devnet operator keys (hardhat account #0 by default). The backend signs
    // aggregator / ledger-anchor transactions with these. Override via env.
    AGGREGATOR_PRIVATE_KEY: process.env.AGGREGATOR_PRIVATE_KEY || "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
    ANCHOR_PRIVATE_KEY: process.env.ANCHOR_PRIVATE_KEY || "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
    AGGREGATOR_PORT: "8081",
    ETH_RPC: "ws://localhost:8545",
    DB_URL: "postgres://postgres:postgres@localhost:5432/well_db?sslmode=disable",
  };

  console.log("\n=== Deployed addresses ===");
  for (const [key, val] of Object.entries(addrs)) {
    console.log(`${key}=${val}`);
  }

  const envLines = Object.entries(addrs)
    .map(([k, v]) => `${k}=${v}`)
    .join("\n");
  writeFileSync(".env", envLines + "\n");
  console.log("\n.env written to project root");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
