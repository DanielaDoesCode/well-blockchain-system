#!/usr/bin/env bash
set -euo pipefail

# Loopback benchmark suite run against the L2 (Optimism devnet).
#
# Mirrors scripts/benchmarks/run_benchmarks.sh (same pipeline, same bench
# scripts, same result schema) but targets the ops-bedrock devnet L2:
#   - never touches the devnet's chain state; only Postgres is reset
#   - signers are topped up through the real L1->L2 bridge
#   - contracts are deployed to the L2, addresses land in .env-l2
#   - results go to scripts/benchmarks/results/l2/*.json with metadata stamped
#     via BENCH_NODE_DESC / BENCH_NET_LATENCY
#
# Prereqs: devnet up (bash l2/setup-devnet.sh), Docker.
# Usage:   bash l2/run-eval.sh

cd "$(dirname "$0")/.."
source scripts/benchmarks/lib-run.sh
source l2/env.sh

bench_kill_stale_backend
trap bench_stop_backend EXIT

l2_check
l2_ensure_funding

echo "=== [1/6] Postgres reset (devnet chain untouched) ==="; l2_reset_db
echo "=== [2/6] Deploy to L2 ==="; l2_deploy_l2

echo "=== [3/6] Backend (ETH_RPC=$L2_RPC_WS) ==="
set -a; . ./.env-l2; set +a
export ETH_RPC="$L2_RPC_WS"
export BENCH_NODE_DESC="$L2_NODE_DESC"
export BENCH_NET_LATENCY="zero (loopback)"
export BENCH_RESULTS_DIR="$L2_RESULTS_DIR"
bench_start_backend

echo "=== [4/6] ZK proof generation benchmark ==="
npx hardhat run scripts/benchmarks/bench_zk.js --network localhost-l2

echo "=== [5/6] Latency + gas benchmark ==="
npx hardhat run scripts/benchmarks/bench_latency.js --network localhost-l2

echo "=== [6/7] Throughput benchmark ==="
npx hardhat run scripts/benchmarks/bench_throughput.js --network localhost-l2

echo "=== [7/7] Ledger integrity check ==="
export CHAIN_HTTP_RPC="$L2_RPC_HTTP"
bench_verify_ledger "$L2_RESULTS_DIR"

echo
echo "=== L2 (devnet) loopback results in $L2_RESULTS_DIR/ ==="
ls -1 "$L2_RESULTS_DIR/" 2>/dev/null || true
echo
echo "Compare with L1 baseline: scripts/benchmarks/results/"
echo "Backend log: $LOGDIR/backend.log"