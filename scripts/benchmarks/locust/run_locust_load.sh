#!/usr/bin/env bash
# Locust load-test runner for the Well EHR API.
#
# One command does the full deterministic cycle:
#   reset chain+DB -> init geth -> start stack -> schema -> deploy -> backend
#   -> seed users/records/sessions -> seed CSV -> run Locust headless
#   -> stop backend (geth stack kept running, like the benchmark runner)
#
# Knobs (env): LOCUST_PATIENTS, LOCUST_DOCTORS, LOCUST_RECORDS_PER_PATIENT,
#   LOCUST_USERS (users spawned), LOCUST_RUN_TIME, LOCUST_SPAWN_RATE,
#   LOCUST_HOST, LOCUST_CSV, LOCUST_WEB=1 (interactive UI instead of headless),
#   LOCUST_SKIP_RESET=1 (reuse running stack), LOCUST_EXTRA (extra locust args),
#   LOCUST_SHAPE=soak (staged ramp/hold/ramp-down LoadTestShape; needs
#     LOCUST_RUN_TIME=6m because the shape self-terminates at 325s and a
#     shorter --run-time truncates it mid-ramp),
#   LOCUST_WRITE_TASKS=1 (enable the /records/upload task; drives IPFS rather
#     than the API, so its numbers are storage-bound),
#   LOCUST_PRE_TASKS=1 (enable the proxy re-encryption tasks: POST /records/
#     {addr}/reencrypt runs the AFGH transform on the backend, POST /records/
#     {addr}/delegate re-stores the patient's rk. See docs/pre-evaluation.md),
#   LOCUST_MIN_WAIT / LOCUST_MAX_WAIT (think-time bounds, read by the locustfile),
#   SESSION_CONCURRENCY (seeder: 1=serial, >1=mints the zkauth sessions as a
#     burst and measures listener ingestion lag into $OUT_DIR/session_ingest.json).
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
cd "$REPO_ROOT"

source scripts/benchmarks/lib-run.sh

export RESULTS_DIR="${RESULTS_DIR:-scripts/benchmarks/results}"
export DB_POLL_CMD="${DB_POLL_CMD:-docker exec postgres psql -U postgres -d well_db -t -A}"

# ---- knobs ---------------------------------------------------------------
LOCUST_PATIENTS="${LOCUST_PATIENTS:-25}"
LOCUST_DOCTORS="${LOCUST_DOCTORS:-5}"
LOCUST_RECORDS_PER_PATIENT="${LOCUST_RECORDS_PER_PATIENT:-2}"
LOCUST_USERS="${LOCUST_USERS:-$((LOCUST_PATIENTS + LOCUST_DOCTORS))}"
LOCUST_SPAWN_RATE="${LOCUST_SPAWN_RATE:-4}"
LOCUST_RUN_TIME="${LOCUST_RUN_TIME:-90s}"
LOCUST_HOST="${LOCUST_HOST:-http://localhost:8081}"
LOCUST_CSV="${LOCUST_CSV:-scripts/benchmarks/results/locust/seed.csv}"
LOCUST_WEB="${LOCUST_WEB:-0}"
SKIP_RESET="${LOCUST_SKIP_RESET:-0}"
LOCUST_KEEP_BACKEND="${LOCUST_KEEP_BACKEND:-0}"
LOCUST_SKIP_BACKEND="${LOCUST_SKIP_BACKEND:-0}"
LOCUST_SKIP_SEED="${LOCUST_SKIP_SEED:-0}"
LOCUST_EXTRA="${LOCUST_EXTRA:-}"
SESSION_CONCURRENCY="${SESSION_CONCURRENCY:-1}"
LOCUST_PRE_TASKS="${LOCUST_PRE_TASKS:-0}"

STAMP="$(date -u +%Y%m%dT%H%M%SZ)"
OUT_DIR="$RESULTS_DIR/locust/$STAMP"
mkdir -p "$OUT_DIR"

if [ -x .venv/bin/locust ]; then
  LOCUST_BIN="$(pwd)/.venv/bin/locust"
elif [ -n "${VIRTUAL_ENV:-}" ] && [ -x "$VIRTUAL_ENV/bin/locust" ]; then
  LOCUST_BIN="$VIRTUAL_ENV/bin/locust"
elif command -v locust >/dev/null 2>&1; then
  LOCUST_BIN="$(command -v locust)"
else
  echo "locust not installed. Create a venv:"
  echo "  python3 -m venv .venv && ./.venv/bin/pip install -U pip locust"
  exit 1
fi

if [ "$LOCUST_KEEP_BACKEND" != "1" ]; then
  trap bench_stop_backend EXIT
fi

# ---- pipeline (mirror bench runner; loopback results untouched) ---------
if [ "$SKIP_RESET" != "1" ]; then
  bench_kill_stale_backend
  bench_reset
  bench_init_geth
  bench_start_stack
  bench_wait_geth
  bench_apply_schema
  bench_deploy
fi

# shellcheck disable=SC1091
set -a; source ./.env; set +a
# the deploy script writes the HTTP port into `.env`; geth serves WS on 8546
if [ -z "${ETH_RPC:-}" ] || [ "$ETH_RPC" = "ws://localhost:8545" ]; then
  export ETH_RPC="ws://localhost:8546"
fi
if [ "$LOCUST_SKIP_BACKEND" != "1" ]; then
  bench_start_backend
fi

if [ "$LOCUST_SKIP_SEED" != "1" ]; then
echo "=== Seeding locust users/records/sessions ==="
# precli is invoked once per record (patient lockbox + rk); build it once.
bench_build_precli
LOCUST_PATIENTS="$LOCUST_PATIENTS" LOCUST_DOCTORS="$LOCUST_DOCTORS" \
LOCUST_RECORDS_PER_PATIENT="$LOCUST_RECORDS_PER_PATIENT" LOCUST_CSV="$LOCUST_CSV" \
SESSION_CONCURRENCY="$SESSION_CONCURRENCY" \
SESSION_INGEST_JSON="$OUT_DIR/session_ingest.json" \
DELEGATION_INGEST_JSON="$OUT_DIR/delegation_ingest.json" \
  npx hardhat run scripts/benchmarks/locust/seed_locust.js --network localhost
fi

echo
echo "=== Running Locust: $LOCUST_USERS users @ $LOCUST_SPAWN_RATE/s for $LOCUST_RUN_TIME against $LOCUST_HOST ==="
# LOCUST_USERS is passed explicitly: the locustfile's StagedSoak shape reads it
# from the environment, and this runner may have *computed* the value above (from
# LOCUST_PATIENTS + LOCUST_DOCTORS), in which case it is not an exported shell
# variable and the shape would silently fall back to its 30-user default.
if [ "$LOCUST_WEB" = "1" ]; then
  LOCUST_SEED_CSV="$LOCUST_CSV" LOCUST_HOST="$LOCUST_HOST" \
  LOCUST_USERS="$LOCUST_USERS" LOCUST_SHAPE="${LOCUST_SHAPE:-}" \
  LOCUST_WRITE_TASKS="${LOCUST_WRITE_TASKS:-0}" \
  LOCUST_PRE_TASKS="$LOCUST_PRE_TASKS" \
    "$LOCUST_BIN" -f "$SCRIPT_DIR/locustfile.py" $LOCUST_EXTRA
else
  LOCUST_SEED_CSV="$LOCUST_CSV" LOCUST_HOST="$LOCUST_HOST" \
  LOCUST_USERS="$LOCUST_USERS" LOCUST_SHAPE="${LOCUST_SHAPE:-}" \
  LOCUST_WRITE_TASKS="${LOCUST_WRITE_TASKS:-0}" \
  LOCUST_PRE_TASKS="$LOCUST_PRE_TASKS" \
    "$LOCUST_BIN" -f "$SCRIPT_DIR/locustfile.py" --headless \
      -u "$LOCUST_USERS" --spawn-rate "$LOCUST_SPAWN_RATE" \
      --run-time "$LOCUST_RUN_TIME" \
      --html "$OUT_DIR/report.html" \
      --csv "$OUT_DIR/stats" \
      --csv-full-history \
      $LOCUST_EXTRA
fi

# Locust exits 2 when a LoadTestShape ends the run (returning None), which is the
# *expected* way a shape-driven soak terminates — not a failure. Under `set -e`
# that non-zero status would abort the script here, skipping the results summary
# and making a completed soak look broken. Anything else non-zero is still fatal.
locust_status=$?
if [ "$locust_status" -ne 0 ] && ! { [ "$locust_status" -eq 2 ] && [ -n "${LOCUST_SHAPE:-}" ]; }; then
  echo "locust exited $locust_status" >&2
  exit "$locust_status"
fi
if [ "$locust_status" -eq 2 ]; then
  echo "(locust exit 2 = LoadTestShape ended the run as designed)"
fi

echo
echo "=== Locust results in $OUT_DIR/ ==="
ls -1 "$OUT_DIR"
echo "Seed CSV: $LOCUST_CSV"
echo "Backend log: $LOGDIR/backend.log"
echo "Keep geth stack running for further runs (docker compose down to stop)."