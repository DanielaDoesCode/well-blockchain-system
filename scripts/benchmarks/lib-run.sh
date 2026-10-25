#!/usr/bin/env bash
# Shared functions for the benchmark runners (loopback + latency-injected).
# Sources must `cd` to the repo root then `source scripts/benchmarks/lib-run.sh`.

set -euo pipefail

COMPOSE="docker compose -f docker-compose.geth.yml"
export DB_POLL_CMD="${DB_POLL_CMD:-docker exec postgres psql -U postgres -d well_db -t -A}"
export BENCH_TMP_DIR="${BENCH_TMP_DIR:-/tmp/well-bench}"
LOGDIR="${BENCH_LOG_DIR:-/tmp/well-bench}"
RESULTS_DIR="${RESULTS_DIR:-scripts/benchmarks/results}"
mkdir -p "$LOGDIR"

BACKEND_PID=""

# kill any stale backend binary left over from a previous (interrupted) run
bench_kill_stale_backend() {
  pkill -f "^$LOGDIR/backend-bin" 2>/dev/null || true
}

# [1] Reset geth chain + Postgres state for a deterministic run.
bench_reset() {
  echo "=== Resetting geth chain + Postgres (deterministic run) ==="
  $COMPOSE stop geth postgres >/dev/null 2>&1 || true
  $COMPOSE rm -f geth postgres >/dev/null 2>&1 || true
  PROJECT_VOLS="$(docker volume ls -q -f label=com.docker.compose.project=well-blockchain 2>/dev/null || true)"
  VOL_TO_WIPE=""
  for VOL in $PROJECT_VOLS; do
    case "$VOL" in
      *geth_data|*postgres_data) VOL_TO_WIPE="$VOL_TO_WIPE $VOL" ;;
    esac
  done
  # also wipe by explicit name so a bare/detached volume can never leak state
  for NAME in well-blockchain_geth_data well-blockchain_postgres_data; do
    case " $VOL_TO_WIPE " in *" $NAME "*) ;; *) VOL_TO_WIPE="$VOL_TO_WIPE $NAME" ;; esac
  done
  if [ -n "$VOL_TO_WIPE" ]; then
    docker volume rm -f $VOL_TO_WIPE >/dev/null 2>&1 || true
  fi
  # verify: both volumes must be gone before a run starts, else the schema and
  # chain state are silently stale (a leaked volume caused a broken WAN cycle)
  for NAME in well-blockchain_geth_data well-blockchain_postgres_data; do
    for _ in $(seq 1 20); do
      if ! docker volume ls -q -f name="$NAME" 2>/dev/null | grep -qx "$NAME"; then break; fi
      sleep 1
    done
    if docker volume ls -q -f name="$NAME" 2>/dev/null | grep -qx "$NAME"; then
      echo "ERROR: volume $NAME still present after reset"; exit 1
    fi
  done
  echo "  wiped geth + postgres volumes (ipfs/vault data retained)"
}

# [2] geth datadir init + account key import (BEFORE up: if geth starts first it
#     auto-inits an empty datadir with the default mainnet genesis).
bench_init_geth() {
  echo "=== Initializing geth datadir + importing account keys ==="
  bash scripts/init-geth.sh
}

# [3] docker stack: postgres, ipfs, vault, geth
bench_start_stack() {
  echo "=== Starting geth stack (postgres, ipfs, vault, geth) ==="
  $COMPOSE up -d postgres ipfs vault geth
}

# [4] wait for geth JSON-RPC
bench_wait_geth() {
  echo "=== Waiting for geth RPC ==="
  READY=0
  for i in $(seq 1 60); do
    if curl -s -o /dev/null -m 2 \
         -H 'Content-Type: application/json' \
         --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
         http://localhost:8545 >/dev/null 2>&1; then
      echo "  geth ready after ${i}s"; READY=1; break
    fi
    sleep 1
  done
  if [ "$READY" != 1 ]; then echo "geth RPC did not come up"; exit 1; fi
}

# [5] apply DB schema (idempotent); the file is IF NOT EXISTS-everywhere, so a
#     stray duplicate run cannot error, but it must not race the entrypoint.
bench_apply_schema() {
  echo "=== Applying DB schema (idempotent) ==="
  PREADY=0
  for i in $(seq 1 30); do
    if docker exec postgres pg_isready -U postgres >/dev/null 2>&1; then
      PREADY=1; break
    fi
    sleep 1
  done
  if [ "$PREADY" != 1 ]; then echo "postgres did not come up"; exit 1; fi
  sleep 3
  timeout 120 docker exec postgres psql -U postgres -d well_db \
    -f /docker-entrypoint-initdb.d/init.sql >/dev/null 2>&1 || {
    echo "schema apply failed"; exit 1; }
}

# [6] compile + deploy (dev wiring), driven by the current networking env:
#     the loopback runner uses --network localhost (geth on 8545), the WAN runner
#     uses its own relays; both point at the same chain.
bench_deploy() {
  echo "=== Compiling + deploying contracts (dev wiring) ==="
  sleep 5
  npx hardhat compile --quiet
  # Wire the real Groth16Verifier into ZKAuthVerifier for genuine on-chain proof
  # verification unless explicitly overridden (BENCH_AUTH_VERIFIER=mock).
  export DEPLOY_AUTH_VERIFIER="${BENCH_AUTH_VERIFIER:-real}"
  timeout 420 npx hardhat run scripts/deploy_mock.js --network localhost || {
    echo "deploy_mock failed/timed out"; exit 1; }
}

# [7] build + start the backend, honouring whatever ETH_RPC / IPFS_API_URL /
#     VAULT_ADDR / AGGREGATOR_PORT are already exported (allows the WAN runner to
#     repoint the backend at relay ports ahead of the service legs). The caller
#     sources ./.env and applies its overrides before calling this.
bench_start_backend() {
  echo "=== Starting backend ==="
  go build -C well-backend -o "$LOGDIR/backend-bin" ./cmd
  "$LOGDIR/backend-bin" > "$LOGDIR/backend.log" 2>&1 < /dev/null &
  BACKEND_PID=$!

  READY=0
  for i in $(seq 1 90); do
    if curl -s -o /dev/null -m 2 http://localhost:8081/ >/dev/null 2>&1; then
      echo "  backend ready after ${i}s"; READY=1; break
    fi
    sleep 1
  done
  if [ "$READY" != 1 ]; then echo "backend did not come up (log: $LOGDIR/backend.log)"; exit 1; fi
}

# Build the reference PRE client once and export PRECLI_BIN. The seeder and
# the latency benchmark invoke it once per record; without a prebuilt binary
# each call falls back to `go run` (a compile per invocation).
bench_build_precli() {
  if [ -n "${PRECLI_BIN:-}" ]; then
    export PRECLI_BIN
    return
  fi
  go build -C well-backend -o "$LOGDIR/precli-bin" ./cmd/precli
  export PRECLI_BIN="$LOGDIR/precli-bin"
}

bench_stop_backend() {
  if [ -n "$BACKEND_PID" ]; then
    kill "$BACKEND_PID" 2>/dev/null || true
    BACKEND_PID=""
  fi
}

LOGVERIFY_BIN="${LOGVERIFY_BIN:-$LOGDIR/logverify-bin}"
LOGVERIFY_BUILT=""

bench_build_logverify() {
  if [ -n "$LOGVERIFY_BUILT" ]; then return; fi
  go build -C well-backend -o "$LOGVERIFY_BIN" ./cmd/logverify
  LOGVERIFY_BUILT=1
}

# scalar DB helper: returns one numeric value from the well_db.
ledger_q() {
  $DB_POLL_CMD -c "$1" 2>/dev/null | tr -d ' \t\r\n'
}

# bench_verify_ledger <resultsDir> — stops the backend (SIGTERM triggers the
# graceful ledger flush: every unbatched row is sealed and every batch anchored),
# then runs a full logverify audit and writes <resultsDir>/ledger.json with the
# row/batch/anchor counts and the verdict. Any ledger inconsistency fails the run.
bench_verify_ledger() {
  local results_dir="${1:-$RESULTS_DIR}"
  mkdir -p "$results_dir"

  local pid="$BACKEND_PID"
  bench_stop_backend
  local waited=0
  while [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; do
    sleep 1
    waited=$((waited + 1))
    if [ "$waited" -gt 30 ]; then
      echo "ERROR: backend did not exit after SIGTERM (ledger flush stuck)"; exit 1
    fi
  done
  [ -n "$pid" ] && echo "  backend exited after graceful LEDGER FLUSH (${waited}s)"

  bench_build_logverify

  echo "=== Ledger integrity check (logverify full) ==="
  local t0 t1 wall exit_code
  t0=$(date +%s)
  "$LOGVERIFY_BIN" full > "$results_dir/ledger-report.json" 2> "$results_dir/ledger-verify.log"
  exit_code=$?
  t1=$(date +%s)
  wall=$(( t1 - t0 ))
  if [ "$exit_code" != 0 ]; then
    echo "LEDGER INTEGRITY CHECK FAILED (exit $exit_code) — evidence: $results_dir/ledger-report.json"
    cat "$results_dir/ledger-verify.log"
    exit 1
  fi
  echo "  logverify full OK (${wall}s)"

  local chain_id=""
  chain_id=$(curl -s -m 5 -H 'Content-Type: application/json' \
    --data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' \
    "${CHAIN_HTTP_RPC:-http://localhost:8545}" 2>/dev/null \
    | node -e 'let s="";process.stdin.on("data",d=>s+=d);process.stdin.on("end",()=>{try{console.log(parseInt(JSON.parse(s).result,16))}catch(e){console.log("")}})')
  [ -n "$chain_id" ] || chain_id="unknown"

  local a_total a_hashed a_unbatched e_total e_hashed e_unbatched
  a_total=$(ledger_q "SELECT count(*) FROM access_logs")
  a_hashed=$(ledger_q "SELECT count(*) FROM access_logs WHERE row_hash IS NOT NULL")
  a_unbatched=$(ledger_q "SELECT count(*) FROM access_logs WHERE batch_id IS NULL")
  e_total=$(ledger_q "SELECT count(*) FROM event_logs")
  e_hashed=$(ledger_q "SELECT count(*) FROM event_logs WHERE row_hash IS NOT NULL")
  e_unbatched=$(ledger_q "SELECT count(*) FROM event_logs WHERE batch_id IS NULL")
  local batches anchored unanchored
  batches=$(ledger_q "SELECT count(*) FROM batched_logs")
  anchored=$(ledger_q "SELECT count(*) FROM batched_logs WHERE anchored = TRUE")
  unanchored=$(( batches - anchored ))

  local unbatched_total=$(( a_unbatched + e_unbatched ))
  local seal_logged anchor_failed
  seal_logged=$(grep -c "ledger: table=" "$LOGDIR/backend.log" 2>/dev/null || true)
  anchor_failed=$(grep -c "not anchored" "$LOGDIR/backend.log" 2>/dev/null || true)
  : "${seal_logged:=0}"; : "${anchor_failed:=0}"

  local verdict="ok"
  if [ "$unbatched_total" != 0 ] || [ "$unanchored" != 0 ]; then
    verdict="FAIL"
  fi
  [ "$verdict" = "ok" ] || { echo "Ledger not fully flushed: unbatched=$unbatched_total unanchored=$unanchored"; exit 1; }

  {
    printf '{\n'
    printf '  "date": "%s",\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
    printf '  "chainId": %s,\n' "$(printf '%s' "$chain_id" | grep -q '^[0-9]*$' && echo "$chain_id" || echo "\"$chain_id\"")"
    printf '  "environment": {\n'
    printf '    "nodeDesc": "%s",\n' "${BENCH_NODE_DESC:-single local geth Clique PoA node (2s period)}"
    printf '    "networkLatency": "%s"\n' "${BENCH_NET_LATENCY:-zero (loopback)}"
    printf '  },\n'
    printf '  "database": {\n'
    printf '    "accessLogs": { "total": %s, "hashed": %s, "unbatched": %s },\n' "$a_total" "$a_hashed" "$a_unbatched"
    printf '    "eventLogs":  { "total": %s, "hashed": %s, "unbatched": %s },\n' "$e_total" "$e_hashed" "$e_unbatched"
    printf '    "batches": %s, "anchoredBatches": %s, "unanchoredBatches": %s,\n' "$batches" "$anchored" "$unanchored"
    printf '    "unbatchedRows": %s\n' "$unbatched_total"
    printf '  },\n'
    printf '  "ledgerVerify": { "verdict": "%s", "wallSeconds": %s, "reportFile": "ledger-report.json" },\n' "$verdict" "$wall"
    printf '  "backendLog": { "sealedBatchesLogged": %s, "anchorFailuresLogged": %s, "logFile": "backend.log" }\n' "$seal_logged" "$anchor_failed"
    printf '}\n'
  } > "$results_dir/ledger.json"

  echo "  ledger.json written to $results_dir/ledger.json"
  echo "  rows: access=$a_total (hashed=$a_hashed) event=$e_total (hashed=$e_hashed)"
  echo "  batches=$batches anchored=$anchored unbatched=$unbatched_total"
}

bench_results_banner() {
  echo
  echo "=== Results written to $RESULTS_DIR/ ==="
  ls -1 "$RESULTS_DIR/" 2>/dev/null || true
  echo
  echo "Backend log: $LOGDIR/backend.log"
  echo "Note: benchmark JSON output is reproducible evidence for the thesis;"
  echo "      keep the geth stack running for interactive verification (docker compose down to stop)."
}