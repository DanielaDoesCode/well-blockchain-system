import "@nomicfoundation/hardhat-ethers";
import "dotenv/config";

// Deploy/funder account override for the devnet L1/L2. ops-bedrock devnets
// vary on which key holds the gas; when required, export DEV_KEY=<privkey> so
// deploy_mock's signers[0] is that funded key. Bench signers stay the well-known
// mnemonic accounts (topped up via the bridge by l2/fund-accounts.js).
const NETWORK_ACCOUNTS = process.env.DEV_KEY ? [process.env.DEV_KEY] : undefined;

// L2 devnet sequencer endpoint (default: legacy ops-bedrock compose layout;
// Kurtosis devnets: export L2_RPC_HTTP).
const L2_RPC_HTTP = process.env.L2_RPC_HTTP || "http://127.0.0.1:9545";

export default {
  solidity: "0.8.20",
  paths: {
    sources: "./src",
    artifacts: "./artifacts",
    cache: "./cache",
  },
  networks: {
    localhost: {
      url: "http://127.0.0.1:8545",
      accounts: NETWORK_ACCOUNTS,
    },
    "localhost-wan": {
      // the latency-injected suite routes JSON-RPC through the WAN RPC relay
      // (18545 -> geth 8545)
      url: "http://127.0.0.1:18545",
      accounts: NETWORK_ACCOUNTS,
    },
    "localhost-l2": {
      // OP Stack devnet L2 sequencer (op-geth)
      url: L2_RPC_HTTP,
      accounts: NETWORK_ACCOUNTS,
    },
    "localhost-l2-wan": {
      // latency-injected suite against the L2: same relay port, upstream = L2 RPC
      url: "http://127.0.0.1:18545",
      accounts: NETWORK_ACCOUNTS,
    },
    hardhat: {
      allowUnlimitedContractSize: true,
      blockGasLimit: 1000000000,
    },
  },
};