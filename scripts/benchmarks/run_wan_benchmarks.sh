#!/usr/bin/env bash
set -euo pipefail

# Latency-injected (wide-area) benchmark suite.
#
# Reuses the exact loopback pipeline (lib-run.sh) so every scenario starts from
# a fully reset chain+DB, then routes traffic through per-leg transparent TCP
# relays that inject a one-way propagation delay drawn from the situation's
# geographic band. The network component is the only injected delay — there is
# no separate human-input model. See
#     docs/benchmark-methodology-wan.md
# for the methodology and what each delta means.
#
# Scenario matrix (defaults):
#   situations (geographic RTT band, jittered per connection):
#     same-region        band  1- 10 ms   (RTT 1-10, one-way 0.5-5ms)
#     cross-region       band 10- 40 ms   (RTT 10-40, one-way 5-20ms)
#     intercontinental   band 100-250 ms  (RTT 100-250, one-way 50-125ms)
#   modes     P (client-remote: bench->geth RPC, bench->API)
#             E (geo-distributed: additionally backend->geth WS, ->Vault, ->IPFS)
#
# Every incoming TCP connection draws ONE RTT from its situation band via the
# seeded PRNG in wan/relay.cjs; both directions of that connection use RTT/2.
# Fresh connections draw fresh samples, so round trips vary across the band.
# The relays share SEED_BASE, so each leg reproduces the same draw sequence.
#
# Knobs:
#   WAN_SITUATIONS "same-region cross-region"  override the situation list
#   WAN_MODES       "P"                        override modes
#   SEED_BASE       20260910                   PRNG seed for the relay bands
#   BENCH_*                                     passed through to the bench scripts
#
# Output: scripts/benchmarks/results/wan/<situation>/<mode>/

cd "$(dirname "$0")/../.."
source scripts/benchmarks/lib-run.sh

# the seeder/latency benchmark call the reference PRE client once per record
bench_build_precli

WAN_SITUATIONS="${WAN_SITUATIONS:-same-region cross-region intercontinental}"
WAN_MODES="${WAN_MODES:-P E}"
SEED_BASE="${SEED_BASE:-20260910}"

# RTT class the scenario lives in: every relay leg samples from this band.
declare -A WAN_BANDS=(
  [same-region]="1-10"
  [cross-region]="10-40"
  [intercontinental]="100-250"
)

# Upstream service ports the relays forward to. Defaults reproduce the loopback
# geth stack; override (e.g. UPSTREAM_RPC_PORT=9545 UPSTREAM_WS_PORT=9546) to run
# the same latency-injected matrix against a different (e.g. L2) node.
UPSTREAM_RPC_PORT="${UPSTREAM_RPC_PORT:-8545}"
UPSTREAM_WS_PORT="${UPSTREAM_WS_PORT:-8546}"
UPSTREAM_API_PORT="${UPSTREAM_API_PORT:-8081}"
UPSTREAM_VAULT_PORT="${UPSTREAM_VAULT_PORT:-8200}"
UPSTREAM_IPFS_PORT="${UPSTREAM_IPFS_PORT:-5001}"

RELAY_DIR="$LOGDIR/relays"
mkdir -p "$RELAY_DIR"
RELAY_PIDS=""
PROBE_LOG="$LOGDIR/wan-probes.log"
: > "$PROBE_LOG"

wan_stop_relays() {
  [ -n "$RELAY_PIDS" ] && kill $RELAY_PIDS 2>/dev/null || true
  RELAY_PIDS=""
}

wan_cleanup() {
  bench_stop_backend
  wan_stop_relays
}
trap wan_cleanup EXIT

# ports: client legs + backend service legs (always started; unused ones idle)
RELAY_RPC=18545; RELAY_API=18443; RELAY_WS=18546; RELAY_VAULT=18200; RELAY_IPFS=15001

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

# HTTP round-trip probe through a relay: proves the injected RTT is live.
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

bench_zk_done=""
for MODE in $WAN_MODES; do
  for SITUATION in $WAN_SITUATIONS; do
    BAND="${WAN_BANDS[$SITUATION]}"
    if [ -z "$BAND" ]; then
      echo "ERROR: unknown situation '$SITUATION' (no band defined)"; exit 1
    fi

    echo
    echo "############################################################"
    echo "### WAN scenario: situation=$SITUATION mode=$MODE band=${BAND}ms (jittered per-connection)"
    echo "############################################################"

    echo "=== [reset] deterministic baseline ==="
    bench_reset
    bench_init_geth
    bench_start_stack
    bench_wait_geth
    bench_apply_schema
    bench_deploy

    echo "=== [relays] band ${BAND}ms, seeded per-connection jitter ==="
    wan_start_relay "$RELAY_RPC" "$UPSTREAM_RPC_PORT" "$BAND"
    wan_start_relay "$RELAY_API" "$UPSTREAM_API_PORT" "$BAND"
    wan_start_relay "$RELAY_WS" "$UPSTREAM_WS_PORT" "$BAND"
    wan_start_relay "$RELAY_VAULT" "$UPSTREAM_VAULT_PORT" "$BAND"
    wan_start_relay "$RELAY_IPFS" "$UPSTREAM_IPFS_PORT" "$BAND"
    for p in $RELAY_RPC $RELAY_API $RELAY_WS $RELAY_VAULT $RELAY_IPFS; do wan_wait_port "$p"; done

    echo "=== [backend] ==="
    set -a; . ./.env; set +a
    if [ "$MODE" = "E" ]; then
      export ETH_RPC="ws://localhost:$RELAY_WS"
      export IPFS_API_URL="http://127.0.0.1:$RELAY_IPFS"
      export VAULT_ADDR="http://127.0.0.1:$RELAY_VAULT"
    else
      export ETH_RPC="ws://localhost:$UPSTREAM_WS_PORT"
      export IPFS_API_URL="http://127.0.0.1:$UPSTREAM_IPFS_PORT"
      export VAULT_ADDR="http://127.0.0.1:$UPSTREAM_VAULT_PORT"
    fi
    bench_start_backend

    echo "=== [probe] injected RTT verification ==="
    wan_probe "rpc" "http://127.0.0.1:$RELAY_RPC" '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'
    wan_probe "api" "http://127.0.0.1:$RELAY_API/" ""

    RESULTS_BASE="scripts/benchmarks/results/wan/${SITUATION}/${MODE}"
    mkdir -p "$RESULTS_BASE"
    export BENCH_RESULTS_DIR="$RESULTS_BASE"
    export BENCH_NET_LATENCY="${SITUATION} (band ${BAND}ms jittered per-connection, seed ${SEED_BASE})"
    export API_BASE_URL="http://127.0.0.1:$RELAY_API"

    echo "=== [ZK benchmark] (network-independent; ran once) ==="
    if [ -z "$bench_zk_done" ]; then
      timeout 600 npx hardhat run scripts/benchmarks/bench_zk.js --network localhost-wan
      bench_zk_done=1
    fi

    echo "=== [latency + gas] ==="
    timeout 900 npx hardhat run scripts/benchmarks/bench_latency.js --network localhost-wan

    echo "=== [throughput] ==="
    timeout 600 npx hardhat run scripts/benchmarks/bench_throughput.js --network localhost-wan

    echo "=== [ledger integrity check] ==="
    bench_verify_ledger "$RESULTS_BASE"

    # unset per-cycle exports so the next cycle re-derives them
    unset BENCH_RESULTS_DIR BENCH_NET_LATENCY API_BASE_URL

    bench_stop_backend
    wan_stop_relays
  done
done

echo
echo "=== WAN matrix complete ==="
echo "Probes (RTT verification):"
cat "$PROBE_LOG"
echo
echo "Results:"
find scripts/benchmarks/results/wan -name "*.json" | sort