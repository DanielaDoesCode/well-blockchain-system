// Full L1 <-> L2 bridge + app-layer comparison demo on the local devnet.
//
// Part 1: assets across the layers.
//   A. L1 -> L2  deposit  (OP Stack OptimismPortal.depositTransaction)
//   B. L2 -> L1  withdrawal initiation (L2ToL1MessagePasser.depositTransaction)
// Part 2: the same Well contract call (registry.register) on BOTH layers to show
//         gas / latency deltas between the devnet L1 and the L2.
//
// Requires: devnet up (bash l2/setup-devnet.sh) and both layers deployed
// (bash l2/deploy-layers.sh -> .env-l1 + .env-l2). Uses the first hardhat
// account as operator.
//
// Usage:  node l2/bridge-demo.js

import { readFileSync } from "fs";
import { ethers } from "ethers";
import { providers, MNEMONIC, accountAt, funderL1, balances, depositToL2, waitL2Credit, withdrawFromL2, fmtEth } from "./lib-bridge.js";

const { l1, l2 } = providers();
const operator = funderL1(l1);
const recipient = accountAt(MNEMONIC, 1);
const AMT_UP = ethers.parseEther("0.05");
const AMT_DOWN = ethers.parseEther("0.02");

const loadEnv = (file) => {
  const out = {};
  for (const line of readFileSync(file, "utf8").split("\n")) {
    const m = line.match(/^([A-Z_]+)=(.*)$/);
    if (m) out[m[1]] = m[2];
  }
  return out;
};

async function registerMetric(provider, wallet, registryAddr, label) {
  const registry = new ethers.Contract(
    registryAddr,
    ["function register(bytes32 _userId, bytes32 _secretHash, uint8 _role) external",
     "function getUserByAddress(address) external view"],
    wallet.connect(provider)
  );
  const pre = await registry.getUserByAddress(wallet.address).catch(() => null);
  if (pre) return { label, skipped: true, gasUsed: 0, blocks: 0, ms: 0 };
  const pending = await provider.getBlockNumber();
  const t0 = Date.now();
  const tx = await registry.register(ethers.id(`u-${label}-${Date.now()}`), ethers.id("s"), 3);
  const rc = await tx.wait();
  const ms = Date.now() - t0;
  return { label, skipped: false, gasUsed: rc.gasUsed.toString(), blocks: rc.blockNumber - pending, ms };
}

async function main() {
  const [chainL1, chainL2] = await Promise.all([(await l1.getNetwork()).chainId, (await l2.getNetwork()).chainId]);
  console.log(`L1 devnet chainid=${chainL1} | L2 chainid=${chainL2}\n`);

  const b0 = await balances({ l1, l2 }, operator.address);
  console.log("operator", operator.address);
  console.log(`  L1 ${fmtEth(b0.l1)}   L2 ${fmtEth(b0.l2)}`);

  // ---------- A: L1 -> L2 deposit ------------------------------------------
  console.log(`\n=== [A] L1 -> L2 deposit of ${fmtEth(AMT_UP)} to ${recipient.address} ===`);
  const r0 = await balances({ l1, l2 }, recipient.address);
  console.log(`recipient before:  L1=${fmtEth(r0.l1)}  L2=${fmtEth(r0.l2)}`);
  const { via, receipt } = await depositToL2(l1, operator.connect(l1), recipient.address, AMT_UP);
  const t0 = Date.now();
  const cred = await waitL2Credit(l2, recipient.address, r0.l2, AMT_UP);
  const secs = ((Date.now() - t0) / 1000).toFixed(1);
  console.log(`  portal tx ${receipt.hash} (L1 block ${receipt.blockNumber}, used ${receipt.gasUsed.toString()} gas)`);
  console.log(`  recipient after: L2=${fmtEth(cred)}  (credited after ~${secs}s)`);

  // ---------- B: L2 -> L1 withdrawal entry ------------------------------------
  console.log(`\n=== [B] L2 -> L1 withdrawal of ${fmtEth(AMT_DOWN)} (initiation) ===`);
  const { nonce, withdrawalHash } = await withdrawFromL2(l2, recipient.connect(l2), operator.address, AMT_DOWN);
  console.log("  message initiated on L2 (L2ToL1MessagePasser)");
  console.log(`    nonce = ${nonce}`);
  console.log(`    withdrawalHash = ${withdrawalHash}`);
  console.log("  The output root is proposed to the L1 OptimismPortal by the devnet");
  console.log("  proposer; after the (short) challenge window the withdrawal proof can");
  console.log("  be finalized on L1 (proveWithdrawalTransaction + finalizeWithdrawal).");

  // ---------- Part 2: app-layer L1-vs-L2 metrics ------------------------------
  console.log("\n=== [C] same Well call (registry.register) on both layers ===");
  const l1env = loadEnv("./.env-l1");
  const l2env = loadEnv("./.env-l2");
  if (!l1env.REGISTRY_ADDRESS || !l2env.REGISTRY_ADDRESS) {
    console.log("  .env-l1/.env-l2 missing — run bash l2/deploy-layers.sh first");
  } else {
    const rows = [
      await registerMetric(l1, operator, l1env.REGISTRY_ADDRESS, `L1 devnet (cid ${chainL1})`),
      await registerMetric(l2, operator, l2env.REGISTRY_ADDRESS, `L2 op (cid ${chainL2})`),
    ];
    for (const r of rows) {
      console.log(r.skipped
        ? `  ${r.label}: already registered (rerun on a fresh account or layer)`
        : `  ${r.label}: gas=${r.gasUsed}  confirm+wait=${r.ms}ms  blocks=${r.blocks}`);
    }
    console.log("\n  Note: op-geth gasUsed is pure EVM gas (code identical on both layers);");
    console.log("  the L2 total fee additionally includes the L1 data fee, not counted in gasUsed.");
  }

  console.log("\ndemo complete.");
}

main().catch((e) => { console.error(e); process.exit(1); });