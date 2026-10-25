#!/usr/bin/env bash
set -euo pipefail

# Latency-injected (wide-area) benchmark suite against the L2 (Optimism devnet).
#
# Same situation matrix and relay topology as scripts/benchmarks/run_wan_benchmarks.sh
# (modes P/E x situations {same-region, cross-region, intercontinental}, RTT
# jittered per connection from a seeded band), but the relays forward to the
# devnet L2 (RPC 9545, WS 9546) instead of the Clique geth, and the devnet chain
# state is never touched.
#
# Result layout (mirrors the L1 suite on purpose):
#   scripts/benchmarks/results/l2/wan/<situation>/<mode>/*.json
#
# Knobs: same as run_wan_benchmarks.sh (WAN_SITUATIONS, WAN_MODES, SEED_BASE).
# Usage:  bash l2/run-wan-eval.sh

cd "$(dirname "$0")/.."
source scripts/benchmarks/lib-run.sh
source l2/env.sh

WAN_SITUATIONS="${WAN_SITUATIONS:-same-region cross-region intercontinental}"
WAN_MODES="${WAN_MODES:-P E}"
declare -A WAN_BANDS=( [same-region]="1-10" [cross-region]="10-40" [intercontinental]="100-250" )
SEED_BASE="${SEED_BASE:-20260910}"

# same fixed client-leg relay ports as the L1 suite; upstreams = L2 services
RELAY_RPC=18545; RELAY_API=18443; RELAY_WS=18546; RELAY_VAULT=18200; RELAY_IPFS=15001

RELAY_DIR="$LOGDIR/relays-l2"
mkdir -p "$RELAY_DIR"
RELAY_PIDS=""
PROBE_LOG="$LOGDIR/wan-probes-l2.log"
: > "$PROBE_LOG"

wan_stop_relays() {
  [ -n "$RELAY_PIDS" ] && kill $RELAY_PIDS 2>/dev/null || true
  RELAY_PIDS=""
}

wan_start_relay() {
  local port="$1" target="$2" band="$3"
  node scripts/benchmarks/wan/relay.cjs "$port" 127.0.0.1 "$target" --band "$band" "$SEED_BASE" \
    > "$RELAY_DIR/relay-${port}.log" 2>&1 < /dev/null &
  RELAY_PIDS="$RELAY_PIDS $!"
}

wan_wait_port() {
  local port="$1" ok=0
  for _ in $(seq 1 50); do
    if (exec 3<>"/dev/tcp/127.0.0.1/$port") 2>/dev/null; then ok=1; break; fi
    sleep 0.2
  done
  if [ "$ok" != 1 ]; then
    echo "relay on :$port did not come up ($RELAY_DIR/relay-${port}.log)"; exit 1
  fi
}

wan_probe() {
  local label="$1" url="$2" data="$3"
  local t0 t1 ms
  t0=$(date +%s%N)
  if [ -n "$data" ]; then
    curl -s -o /dev/null -m 5 -H 'Content-Type: application/json' --data "$data" "$url" || true
  else
    curl -s -o /dev/null -m 5 "$url" || true
  fi
  t1=$(date +%s%N)
  ms=$(( (t1 - t0) / 1000000 ))
  echo "  $label round-trip = ${ms}ms"
  echo "$label $ms" >> "$PROBE_LOG"
}

wan_cleanup() { bench_stop_backend; wan_stop_relays; }
trap wan_cleanup EXIT

bench_kill_stale_backend
l2_check
l2_ensure_funding

# strip the L2 ports from the WS/HTTP URLs for the relays' upstream targets
L2_HTTP_PORT=${L2_RPC_HTTP##*:}
L2_WS_PORT=${L2_RPC_WS##*:}

bench_zk_done=""
for MODE in $WAN_MODES; do
  for SITUATION in $WAN_SITUATIONS; do
    BAND="${WAN_BANDS[$SITUATION]}"
    if [ -z "$BAND" ]; then
      echo "ERROR: unknown situation '$SITUATION' (no band defined)"; exit 1
    fi

    echo
    echo "############################################################"
    echo "### L2 WAN scenario: situation=$SITUATION mode=$MODE band=${BAND}ms (jittered per-connection)"
    echo "############################################################"

    echo "=== [reset] Postgres (devnet chain untouched) ==="
    l2_reset_db
    l2_deploy_l2

    echo "=== [relays] L2 upstreams, band ${BAND}ms seeded per-connection jitter ==="
    wan_start_relay "$RELAY_RPC" "$L2_HTTP_PORT" "$BAND"
    wan_start_relay "$RELAY_API" 8081 "$BAND"
    wan_start_relay "$RELAY_WS" "$L2_WS_PORT" "$BAND"
    wan_start_relay "$RELAY_VAULT" 8200 "$BAND"
    wan_start_relay "$RELAY_IPFS" 5001 "$BAND"
    for p in $RELAY_RPC $RELAY_API $RELAY_WS $RELAY_VAULT $RELAY_IPFS; do wan_wait_port "$p"; done

      echo "=== [backend] ==="
      set -a; . ./.env-l2; set +a
      if [ "$MODE" = "E" ]; then
        export ETH_RPC="ws://127.0.0.1:$RELAY_WS"
        export IPFS_API_URL="http://127.0.0.1:$RELAY_IPFS"
        export VAULT_ADDR="http://127.0.0.1:$RELAY_VAULT"
      else
        export ETH_RPC="ws://127.0.0.1:$L2_WS_PORT"
        export IPFS_API_URL="http://127.0.0.1:5001"
        export VAULT_ADDR="http://127.0.0.1:8200"
      fi
      export BENCH_NODE_DESC="$L2_NODE_DESC"
      bench_start_backend

      echo "=== [probe] injected RTT verification ==="
      wan_probe "rpc" "http://127.0.0.1:$RELAY_RPC" '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'
      wan_probe "api" "http://127.0.0.1:$RELAY_API/" ""

      RESULTS_BASE="scripts/benchmarks/results/l2/wan/${SITUATION}/${MODE}"
      mkdir -p "$RESULTS_BASE"
      export BENCH_RESULTS_DIR="$RESULTS_BASE"
      export BENCH_NET_LATENCY="${SITUATION} (band ${BAND}ms jittered per-connection, seed ${SEED_BASE})"
      export API_BASE_URL="http://127.0.0.1:$RELAY_API"

      echo "=== [ZK benchmark] (network-independent; ran once) ==="
      if [ -z "$bench_zk_done" ]; then
        timeout 600 npx hardhat run scripts/benchmarks/bench_zk.js --network localhost-l2-wan
        bench_zk_done=1
      fi

      echo "=== [latency + gas] ==="
      timeout 900 npx hardhat run scripts/benchmarks/bench_latency.js --network localhost-l2-wan

      echo "=== [throughput] ==="
      timeout 600 npx hardhat run scripts/benchmarks/bench_throughput.js --network localhost-l2-wan

      echo "=== [ledger integrity check] ==="
      export CHAIN_HTTP_RPC="$L2_RPC_HTTP"
      bench_verify_ledger "$RESULTS_BASE"

      unset BENCH_RESULTS_DIR BENCH_NET_LATENCY API_BASE_URL

      bench_stop_backend
      wan_stop_relays
    done
  done

echo
echo "=== L2 WAN matrix complete ==="
echo "Probes (RTT verification):"
cat "$PROBE_LOG"
echo
echo "Results:"
find scripts/benchmarks/results/l2/wan -name "*.json" 2>/dev/null | sort