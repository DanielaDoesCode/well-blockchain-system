// Top up the hardhat signer accounts on the L2 via the L1->L2 bridge.
//
// The 10 well-known "test ... junk" accounts are prefunded by the ops-bedrock
// devnet on L1; on L2 they are usually seeded too, but the exact seed set varies
// between devnet versions. This script guarantees everyone the bench suite uses
// (signers 0..7) holds >= MIN_BAL on the L2, funding shortfall through the real
// bridge so the *delivered* L2 currency is genuinely cross-chain.
//
// Usage:  node l2/fund-accounts.js
// Idempotent: skips accounts already at/above the threshold.

import { ethers } from "ethers";
import { providers, MNEMONIC, accountAt, funderL1, balances, depositToL2, waitL2Credit, fmtEth } from "./lib-bridge.js";

const INITIAL_BAL = ethers.parseEther("0.5");
const TOP_UP = (i) => ethers.parseEther(i < 7 ? "2" : "1"); // heavier for bench signers

const { l1, l2 } = providers();
const funder = funderL1(l1);

async function main() {
  const [chainL1, chainL2] = await Promise.all([
    (await l1.getNetwork()).chainId,
    (await l2.getNetwork()).chainId,
  ]);
  const l2Name = await l2.send("web3_clientVersion").catch(() => "n/a");
  console.log(`L1 devnet chainid=${chainL1} | L2 chainid=${chainL2} (${l2Name})`);

  // Targets: the optional DEV_KEY deployer first, then the mnemonic accounts.
  const targets = [];
  if (process.env.DEV_KEY) {
    const k = new ethers.Wallet(process.env.DEV_KEY);
    targets.push({ label: "DEV_KEY(deployer)", address: k.address, i: 0 });
  }
  for (let i = 0; i < 10; i++) targets.push({ label: `mnemonic[${i}]`, address: accountAt(MNEMONIC, i).address, i });

  let funded = 0;
  for (const t of targets) {
    const b = await balances({ l1, l2 }, t.address);
    console.log(`\n[${t.label}] ${t.address}  L1=${fmtEth(b.l1)}  L2=${fmtEth(b.l2)}`);
    if (b.l2 >= INITIAL_BAL) {
      console.log("  already funded on L2, skipping");
      continue;
    }
    if (b.l1.isZero()) {
      console.log("  WARNING: no L1 balance to bridge from; devnet L1 seed differs — fund this one manually");
      continue;
    }
    const value = TOP_UP(t.i);
    console.log(`  bridging ${fmtEth(value)} L1 -> L2 ...`);
    const { via } = await depositToL2(l1, funder, t.address, value);
    const bal = await waitL2Credit(l2, t.address, b.l2, value);
    funded++;
    console.log(`  done via ${via}; L2 now ${fmtEth(bal)}`);
  }

  console.log(`\nbridge funding complete: ${funded} account(s) topped up`);
  if (funded === 0) console.log("all signer accounts already funded on the L2");
}

main().catch((e) => { console.error(e); process.exit(1); });