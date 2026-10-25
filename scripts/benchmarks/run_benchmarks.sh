#!/usr/bin/env bash
set -euo pipefail

# Runs the full local benchmark suite against the geth (Clique PoA) stack.
# Pipeline steps are implemented in scripts/benchmarks/lib-run.sh so the
# latency-injected suite (run_wan_benchmarks.sh) reuses the same reset ->
# init-geth -> deploy -> backend lifecycle.
#
# Prerequisites: Docker with the docker compose plugin, Node deps installed
# (npm install), zk-circuits build artifacts present (make zkauth/key variants),
# and Rust circom available if artifacts must be rebuilt.

cd "$(dirname "$0")/../.."
source scripts/benchmarks/lib-run.sh

bench_kill_stale_backend
trap bench_stop_backend EXIT

echo "=== [1/12] ==="; bench_reset
echo "=== [2/12] ==="; bench_init_geth
echo "=== [3/12] ==="; bench_start_stack
echo "=== [4/12] ==="; bench_wait_geth
echo "=== [5/12] ==="; bench_apply_schema
echo "=== [6/12] ==="; bench_deploy
echo "=== [7/12] ==="
# geth exposes WebSocket on 8546 in this stack (deploy_mock's ETH_RPC points at
# the Hardhat-node default ws:8545), so override unconditionally.
set -a; . ./.env; set +a
export ETH_RPC="ws://localhost:8546"
bench_start_backend
bench_build_precli

echo "=== [8/12] PRE primitive benchmark (AFGH, in-process) ==="
node scripts/benchmarks/bench_pre.js

echo "=== [9/12] ZK proof generation benchmark ==="
npx hardhat run scripts/benchmarks/bench_zk.js --network localhost

echo "=== [10/12] Latency + gas benchmark ==="
npx hardhat run scripts/benchmarks/bench_latency.js --network localhost

echo "=== [11/12] Throughput benchmark ==="
npx hardhat run scripts/benchmarks/bench_throughput.js --network localhost

echo "=== [12/12] Ledger integrity check ==="
bench_verify_ledger "$RESULTS_DIR"

bench_results_banner