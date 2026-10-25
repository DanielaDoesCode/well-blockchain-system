// Shared L1<->L2 bridge helpers for the OP Stack devnet evaluation.
//
// Endpoints and addresses are derived at runtime where possible so the helper
// is resilient to the devnet's exact contract versions:
//   - OptimismPortal (L1) canonical devnet address:
//       0x5FbDB2315678afecb367f032d93F642f64180aa3
//     (deterministically deployed by the OP Stack deployer)
//   - L1StandardBridge (L1) canonical devnet address:
//       0x99C9fc46f92E8a1c0deC1b1747d010903E884bE1
//   - L2ToL1MessagePasser (L2) predeploy:
//       0x4200000000000000000000000000000000000016
//
// "depositTransaction(address,uint256,uint64,bool,bytes)" is the classic, still
// present entry point on both the portal and the message passer.

import { ethers } from "ethers";

export const PORTAL_L1 = "0x5FbDB2315678afecb367f032d93F642f64180aa3";
export const BRIDGE_L1 = "0x99C9fc46f92E8a1c0deC1b1747d010903E884bE1";
export const MESSAGE_PASSER_L2 = "0x4200000000000000000000000000000000000016";

export const MNEMONIC = process.env.MNEMONIC || "test test test test test test test test test test test junk";
export const ACCOUNT_PATH = (i) => `m/44'/60'/0'/0/${i}`;
export const accountAt = (phrase, i) => ethers.HDNodeWallet.fromPhrase(phrase, undefined, ACCOUNT_PATH(i));

// The L1 funder: in ops-bedrock devnets the funded key varies; export
// DEV_KEY=<privkey> to override the default mnemonic account 0.
export const funderL1 = (l1Provider) =>
  process.env.DEV_KEY
    ? new ethers.Wallet(process.env.DEV_KEY, l1Provider)
    : accountAt(MNEMONIC, 0).connect(l1Provider);

export const providers = () => ({
  l1: new ethers.JsonRpcProvider(process.env.L1_RPC_HTTP || "http://127.0.0.1:8545"),
  l2: new ethers.JsonRpcProvider(process.env.L2_RPC_HTTP || "http://127.0.0.1:9545"),
});

export const fmtEth = (v) => `${ethers.formatEther(v)} ETH`;

export async function balances(providers_, addr) {
  return {
    l1: await providers_.l1.getBalance(addr),
    l2: await providers_.l2.getBalance(addr),
  };
}

// Detect which L1 bridge contract is present. Returns "portal" | "bridge" | null.
async function detectDepositTarget(l1) {
  if ((await l1.getCode(PORTAL_L1)).length > 2) return "portal";
  if ((await l1.getCode(BRIDGE_L1)).length > 2) return "bridge";
  return null;
}

const DEPOSIT_TX_IFACE = new ethers.Interface([
  "function depositTransaction(address _to, uint256 _value, uint64 _gasLimit, bool _isCreation, bytes _data) payable",
]);
const DEPOSIT_ETH_IFACE = new ethers.Interface([
  "function depositETH(uint32 _minGasLimit, bytes _extraData) payable",
]);

async function usableDepositTx(l1, portal, to, valueWei) {
  const data = DEPOSIT_TX_IFACE.encodeFunctionData("depositTransaction", [to, valueWei, 200000n, false, "0x"]);
  return l1.estimateGas({ to: portal, value: valueWei, data }).then(() => data).catch(() => null);
}

async function usableDepositEth(l1, bridge, to, valueWei) {
  const data = DEPOSIT_ETH_IFACE.encodeFunctionData("depositETH", [200000, "0x"]);
  return l1.estimateGas({ to: bridge, value: valueWei, data }).then(() => data).catch(() => null);
}

// Bridge `valueWei` ETH from `sender` (L1 signer) to `to` (L2 address).
// Returns the L1 transaction and which entry point was used.
export async function depositToL2(l1Provider, sender, to, valueWei) {
  const target = await detectDepositTarget(l1Provider);
  if (!target) {
    throw new Error(
      "no portal/bridge contract found on the L1 devnet — is the ops-bedrock devnet up? (bash l2/setup-devnet.sh)"
    );
  }

  let call = null;
  let via = null;
  if (target === "portal") {
    call = await usableDepositTx(l1Provider, PORTAL_L1, to, valueWei);
    via = "portal.depositTransaction";
    if (!call) {
      call = await usableDepositEth(l1Provider, PORTAL_L1, to, valueWei);
      via = "portal.depositETH";
    }
  }
  if (!call && target === "bridge") {
    call = await usableDepositEth(l1Provider, BRIDGE_L1, to, valueWei);
    via = "standardBridge.depositETH";
  }
  if (!call) {
    throw new Error(
      `deposit calls reverted on ${target === "portal" ? PORTAL_L1 : BRIDGE_L1} — ` +
        "this devnet likely uses an unexpected bridge ABI; check `make devnet` output"
    );
  }

  const tx = await sender.sendTransaction({ to: target === "portal" ? PORTAL_L1 : BRIDGE_L1, value: valueWei, data: call });
  const receipt = await tx.wait();
  return { via, receipt, tx };
}

// Wait until `addr`'s L2 balance reaches `from + minIncrease` (poll every 2 s).
export async function waitL2Credit(l2, addr, from, minIncrease, timeoutMs = 180000) {
  const deadline = Date.now() + timeoutMs;
  while (Date.now() < deadline) {
    const bal = await l2.getBalance(addr);
    if (bal >= from + minIncrease) return bal;
    await new Promise((r) => setTimeout(r, 2000));
  }
  throw new Error(`L1->L2 deposit not credited to ${addr} within ${timeoutMs}ms (raw balance stuck)`);
}

// Trigger an L2->L1 withdrawal of `valueWei` from `sender` (signer on L2) to
// `l1To`. Returns { nonce, withdrawalHash } parsed from the MessagePassed event.
export async function withdrawFromL2(l2Provider, sender, l1To, valueWei) {
  const iface = new ethers.Interface([
    "function depositTransaction(address _to, uint256 _value, uint64 _gasLimit, bool _isCreation, bytes _data) payable",
    "event MessagePassed(uint256 indexed nonce,address indexed sender,address indexed target,uint256 value,uint256 gasLimit,bytes data,bytes32 withdrawalHash)",
  ]);
  const data = iface.encodeFunctionData("depositTransaction", [l1To, valueWei, 500000n, false, "0x"]);
  const tx = await sender.sendTransaction({ to: MESSAGE_PASSER_L2, value: valueWei, data });
  const receipt = await tx.wait();
  const ev = receipt.logs
    .map((l) => { try { return iface.parseLog(l); } catch { return null; } })
    .find((e) => e?.name === "MessagePassed");
  return {
    receipt,
    nonce: ev?.args?.nonce?.toString() ?? "n/a",
    withdrawalHash: ev?.args?.withdrawalHash ?? "n/a",
  };
}