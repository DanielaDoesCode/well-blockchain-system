#!/bin/bash
set -e

# run_geth_doc_test.sh
#
# Same end-to-end ZK-policy document test as run_manual_doc_test.sh, but backed
# by Geth (Clique PoA, 2s block period) instead of the Hardhat node — this is
# the "Phase 0: local Geth validation" flow from THESIS-AWS-SETUP.md.
#
#   ./scripts/run_geth_doc_test.sh --file ./report.pdf [--department name] [--institution name] [--keep]
#
# The stack it manages (all via docker-compose.geth.yml):
#   geth (8545 HTTP / 8546 WS), postgres, ipfs, vault + the Go backend on :8081
#
# Flow: init Geth datadir (if needed) -> start services -> deploy_mock.js
# against Geth RPC -> fix .env ETH_RPC to Geth's WS port -> start backend ->
# run manual_doc_test.js -> cleanup.

ORIG_DIR="$(pwd)"
cd "$(dirname "$0")/.."
PROJECT_ROOT="$(pwd)"
COMPOSE="docker-compose.geth.yml"

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
  echo "  --keep: leave geth + backend running for manual curl afterwards"
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

echo "=== Ensuring Geth datadir is initialized (init-geth.sh) ==="
./scripts/init-geth.sh

echo "=== Starting infrastructure (geth/postgres/ipfs/vault) ==="
docker compose -f "$COMPOSE" up -d geth postgres ipfs vault
sleep 5

echo "=== Applying DB schema (idempotent) ==="
docker exec -i postgres psql -U postgres -d well_db < db/init.sql 2>/dev/null || true

lsof -ti:8081 | xargs kill -9 2>/dev/null || true

echo "=== Deploying contracts to Geth (deploy_mock.js, real PolicyVerifier) ==="
npx hardhat run scripts/deploy_mock.js --network localhost

echo "=== Pointing backend at Geth WS (:8546) ==="
sed -i 's|^ETH_RPC=ws://localhost:8545|ETH_RPC=ws://localhost:8546|' .env

echo "=== Starting backend ==="
set -a && . ./.env && set +a
cd well-backend
go build -o /tmp/well-backend-bin ./cmd/
/tmp/well-backend-bin &>/tmp/backend_geth.log &
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
  echo "  Geth:    :8545 HTTP / :8546 WS (logs: docker compose -f $COMPOSE logs geth)"
  echo "  Backend: $BEPID (logs: /tmp/backend_geth.log, API :8081)"
  echo "  Stop with: kill $BEPID; docker compose -f $COMPOSE stop geth"
  exit $E2E_EXIT
fi

echo "=== Cleaning up ==="
kill $BEPID 2>/dev/null
rm -f /tmp/well-backend-bin

exit $E2E_EXIT
