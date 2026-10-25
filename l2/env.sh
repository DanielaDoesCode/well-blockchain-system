#!/usr/bin/env bash
# Shared configuration + helpers for the L2 (OP Stack devnet) evaluation suite.
#
# Sourced by the l2/ runners. Endpoints default to ops-bedrock `make devnet`
# and can be overridden with L1_RPC_HTTP / L1_RPC_WS / L2_RPC_HTTP / L2_RPC_WS.
# l2_check() verifies the layers are actually the ops-bedrock devnet (chain id
# 900 = L1 devnet, 901 = L2) so a stale/wrong node is never benchmarked.

set -euo pipefail

# --- devnet endpoints (ops-bedrock `make devnet` defaults) -------------------
L1_RPC_HTTP="${L1_RPC_HTTP:-http://127.0.0.1:8545}"
L1_RPC_WS="${L1_RPC_WS:-ws://127.0.0.1:8546}"
L2_RPC_HTTP="${L2_RPC_HTTP:-http://127.0.0.1:9545}"
L2_RPC_WS="${L2_RPC_WS:-ws://127.0.0.1:9546}"
L1_CHAIN_ID="${L1_CHAIN_ID:-900}"   # ops-bedrock devnet L1
L2_CHAIN_ID="${L2_CHAIN_ID:-901}"   # ops-bedrock devnet L2 (op-geth)

L2_RESULTS_DIR="${L2_RESULTS_DIR:-scripts/benchmarks/results/l2}"
OPTIMISM_DIR="${OPTIMISM_DIR:-$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)/../optimism}"

# Metadata stamped into the benchmark JSON so the L2 evidence is self-describing.
export L2_NODE_DESC="${BENCH_NODE_DESC:-Optimism devnet L2 (op-geth, ops-bedrock, chainid $L2_CHAIN_ID)}"

# read eth_chainId (hex JSON-RPC result) of $1 -> decimal ("" if unreachable)
l2_chainid() {
  curl -s -m 3 -H 'Content-Type: application/json' \
    --data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' "$1" 2>/dev/null \
  | node -e 'let s="";process.stdin.on("data",d=>s+=d);process.stdin.on("end",()=>{try{console.log(parseInt(JSON.parse(s).result,16))}catch(e){console.log("")}})'
}

l2_check() {
  local c1 c2
  c1=$(l2_chainid "$L1_RPC_HTTP")
  c2=$(l2_chainid "$L2_RPC_HTTP")
  echo "=== L2 devnet endpoints ==="
  [ -n "$c1" ] && echo "  L1 $L1_RPC_HTTP chainId=$c1" || echo "  L1 $L1_RPC_HTTP unreachable"
  [ -n "$c2" ] && echo "  L2 $L2_RPC_HTTP chainId=$c2" || echo "  L2 $L2_RPC_HTTP unreachable"
  if [ -z "$c1" ] || [ -z "$c2" ]; then
    echo "ERROR: both devnet layers must be reachable. Start them first:" >&2
    echo "  bash l2/setup-devnet.sh" >&2
    echo "Then check the endpoints: export L1_RPC_HTTP/L1_RPC_WS/L2_RPC_HTTP/L2_RPC_WS to match your devnet" >&2
    exit 1
  fi
  # The invariant that matters: the two endpoints are DISTINCT chains. The
  # ops-bedrock compose devnet uses 900/901; Kurtosis devnets use other ids.
  # Pin them via L1_CHAIN_ID/L2_CHAIN_ID if you want strict enforcement.
  if [ "$c1" = "$c2" ]; then
    echo "ERROR: L1 and L2 report the same chainId ($c1) — point L1_RPC_HTTP / L2_RPC_HTTP at the right layers" >&2
    exit 1
  fi
  echo "  devnet OK (L1 chainid=$c1, L2 chainid=$c2; expected defaults: L1=$L1_CHAIN_ID, L2=$L2_CHAIN_ID)"
  echo "  L1 WS: $L1_RPC_WS"
  echo "  L2 WS: $L2_RPC_WS"
}

# [postgres-only reset] Deterministic, loopback-parity DB reset while the chain
# state lives in the devnet (which we must NOT touch). Our own compose stack
# provides only postgres on L2 runs; wiping its volume gives a clean DB the way
# bench_reset does for the L1 pipeline.
l2_reset_db() {
  echo "=== Resetting Postgres (L2 chain state is left untouched) ==="
  $COMPOSE stop geth postgres >/dev/null 2>&1 || true
  docker volume rm -f well-blockchain_postgres_data >/dev/null 2>&1 || true
  for _ in $(seq 1 20); do
    if ! docker volume ls -q -f name=well-blockchain_postgres_data 2>/dev/null | grep -qx well-blockchain_postgres_data; then break; fi
    sleep 1
  done
  if docker volume ls -q -f name=well-blockchain_postgres_data 2>/dev/null | grep -qx well-blockchain_postgres_data; then
    echo "ERROR: postgres volume still present after reset"; exit 1
  fi
  $COMPOSE up -d postgres
  bench_apply_schema
}

# Deploy the Well contracts to the devnet L2 (addresses -> .env-l2).
# The backend + bench scripts read these from DEPLOY_ENV_FILE via deploy_mock.js.
l2_deploy_l2() {
  echo "=== Compiling + deploying contracts to the devnet L2 ==="
  npx hardhat compile --quiet
  export DEPLOY_AUTH_VERIFIER="${BENCH_AUTH_VERIFIER:-real}"
  timeout 480 env DEPLOY_ENV_FILE=".env-l2" DEPLOY_ETH_RPC="$L2_RPC_WS" \
    npx hardhat run scripts/deploy_mock.js --network localhost-l2 || {
      echo "L2 deploy failed/timed out"; exit 1; }
}

# Top up hardhat accounts on L2 via the L1->L2 bridge so bench signers are funded.
l2_ensure_funding() {
  echo "=== Ensuring hardhat signers are funded on the L2 (via L1->L2 bridge) ==="
  L1_RPC_HTTP="$L1_RPC_HTTP" L2_RPC_HTTP="$L2_RPC_HTTP" node l2/fund-accounts.js
}