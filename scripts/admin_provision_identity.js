import "dotenv/config";
import hre from "hardhat";

const { ethers } = hre;

const API_BASE = `http://localhost:${process.env.AGGREGATOR_PORT || 8081}`;

// Mints a ZK-auth session for the admin signer (dummy proof: dev harness uses
// the MockVerifier). In production the proof would cover secret/commitment.
async function adminSession(zkAuth, signer) {
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
  if (!authLog) throw new Error("Authenticated event not found for admin " + signer.address);
  return authLog.args.sessionToken;
}

// Usage:
//   AGGREGATOR_PORT=8081 node scripts/admin_provision_identity.js \
//     --national-id PAT-2024-0001 --name "Alice Patient" --role patient \
//     --address 0x... [--summary 0x...] [--commitment 0x...] [--internal-id 0x...]
async function main() {
  const args = {};
  for (let i = 2; i < process.argv.length; i++) {
    const k = process.argv[i].replace(/^--/, "");
    args[k] = process.argv[++i];
  }

  const nationalID = args["national-id"];
  const name = args.name;
  const role = args.role;
  const address = args.address;
  if (!nationalID || !name || !role || !address) {
    throw new Error("--national-id, --name, --role and --address are required");
  }

  const zkAuthAddress = process.env.ZKAUTH_VERIFIER_ADDRESS;
  if (!zkAuthAddress) throw new Error("missing ZKAUTH_VERIFIER_ADDRESS in .env");
  const zkAuth = await ethers.getContractAt("ZKAuthVerifier", zkAuthAddress);

  const [deployer] = await ethers.getSigners();
  const token = await adminSession(zkAuth, deployer);

  const body = {
    national_id: nationalID,
    display_name: name,
    role,
    address,
    summary_address: args.summary || "",
    commitment: args.commitment || "",
    internal_id: args["internal-id"] || "",
  };

  const res = await fetch(`${API_BASE}/api/v1/identities`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-User-Address": deployer.address,
      "X-Session-Token": token,
    },
    body: JSON.stringify(body),
  });
  const data = await res.json();
  if (res.status !== 201) {
    throw new Error(`provision failed (${res.status}): ${JSON.stringify(data)}`);
  }
  console.log("[PROVISIONED]", JSON.stringify(data.identity, null, 2));
  console.log(`\nNow a doctor/responder can resolve it with:\n` +
    `  GET ${API_BASE}/api/v1/identities?national_id=${nationalID}`);
}

main().catch(err => {
  console.error("Failed:", err.message);
  process.exit(1);
});