#!/bin/bash
set -e

ORIG_DIR="$(pwd)"
cd "$(dirname "$0")/.."
PROJECT_ROOT="$(pwd)"

FILE=""
DEPARTMENT=""
INSTITUTION=""
KEEP=0
while [ $# -gt 0 ]; do
  case "$1" in
    --file)         FILE="$2"; shift 2 ;;
    --department)   DEPARTMENT="$2"; shift 2 ;;
    --institution)  INSTITUTION="$2"; shift 2 ;;
    --keep)         KEEP=1; shift ;;
    *) echo "Unknown argument: $1"; echo "Usage: $0 --file <path> [--department <name>] [--institution <name>] [--keep]"; exit 1 ;;
  esac
done

if [ -z "$FILE" ]; then
  echo "Usage: $0 --file <path> [--department <name>] [--institution <name>] [--keep]"
  echo "  e.g.: $0 --file ./report.pdf --department cardiology --institution general"
  echo "  --keep: leave hardhat node + backend running for manual curl afterwards"
  exit 1
fi

# Resolve the file path: absolute paths pass through; relative paths are tried
# against the original working directory first, then the project root.
case "$FILE" in
  /*) ;;
  *)
    if [ ! -f "$FILE" ] && [ -f "$ORIG_DIR/$FILE" ]; then
      FILE="$ORIG_DIR/$FILE"
    fi
    ;;
esac

if [ ! -f "$FILE" ]; then
  echo "File not found: $FILE"
  echo "  (searched in: $PROJECT_ROOT and $ORIG_DIR)"
  exit 1
fi

echo "=== Starting infrastructure (postgres/ipfs/vault) ==="
docker compose up -d postgres ipfs vault
sleep 5

echo "=== Applying DB schema (idempotent) ==="
docker exec -i postgres psql -U postgres -d well_db < db/init.sql 2>/dev/null || true

lsof -ti:8545 | xargs kill -9 2>/dev/null || true
lsof -ti:8081 | xargs kill -9 2>/dev/null || true

echo "=== Starting Hardhat node ==="
npx hardhat node --hostname 0.0.0.0 --port 8545 &>/tmp/hardhat_manual.log &
HHPID=$!
sleep 3

echo "=== Deploying mock contracts (real PolicyVerifier) ==="
npx hardhat run scripts/deploy_mock.js --network localhost

echo "=== Starting backend ==="
set -a && . ./.env && set +a
cd well-backend
go run ./cmd/ &>/tmp/backend_manual.log &
BEPID=$!
cd ..
sleep 3

echo "=== Running manual doc test with file: $FILE ==="
run_env=( "MANUAL_DOC_FILE=$FILE" )
if [ -n "$DEPARTMENT" ]; then run_env+=( "MANUAL_DOC_DEPARTMENT=$DEPARTMENT" ); fi
if [ -n "$INSTITUTION" ]; then run_env+=( "MANUAL_DOC_INSTITUTION=$INSTITUTION" ); fi
set +e
env "${run_env[@]}" npx hardhat run scripts/manual_doc_test.js --network localhost
E2E_EXIT=$?
set -e

if [ "$KEEP" = "1" ]; then
  echo ""
  echo "=== Stack left running (--keep) ==="
  echo "  Hardhat node: $HHPID (logs: /tmp/hardhat_manual.log, RPC :8545)"
  echo "  Backend:      $BEPID (logs: /tmp/backend_manual.log, API :8081)"
  echo "  Stop with: kill $BEPID $HHPID"
  exit $E2E_EXIT
fi

echo "=== Cleaning up ==="
kill $BEPID 2>/dev/null
kill $HHPID 2>/dev/null

exit $E2E_EXIT
