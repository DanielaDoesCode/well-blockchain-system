#!/usr/bin/env bash
set -euo pipefail

# Deploy the Well contracts to BOTH layers of the local devnet so the bridge
# demo and the L1-vs-L2 comparison have live contracts on each:
#   L1 devnet -> .env-l1   (DEPLOY_ETH_RPC=ws://127.0.0.1:8546)
#   L2 devnet -> .env-l2   (DEPLOY_ETH_RPC=ws://127.0.0.1:9546)
# .env-l1/.env-l2 are gitignored and feed well-backend + the bench scripts.

cd "$(dirname "$0")/.."
source scripts/benchmarks/lib-run.sh
source l2/env.sh

l2_check

npx hardhat compile --quiet

echo
echo "=== Deploying to devnet L1 (chainid $L1_CHAIN_ID) ==="
timeout 480 env DEPLOY_ENV_FILE=".env-l1" DEPLOY_ETH_RPC="$L1_RPC_WS" \
  npx hardhat run scripts/deploy_mock.js --network localhost || {
    echo "L1 deploy failed/timed out"; exit 1; }

echo
echo "=== Deploying to devnet L2 (chainid $L2_CHAIN_ID) ==="
timeout 480 env DEPLOY_ENV_FILE=".env-l2" DEPLOY_ETH_RPC="$L2_RPC_WS" \
  npx hardhat run scripts/deploy_mock.js --network localhost-l2 || {
    echo "L2 deploy failed/timed out"; exit 1; }

echo
echo "=== Done. .env-l1 and .env-l2 written ==="
echo "Next: node l2/fund-accounts.js  (top up L2 signers via the bridge)"
echo "      bash l2/run-eval.sh      (loopback eval on the L2)"