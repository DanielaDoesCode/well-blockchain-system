#!/bin/bash
set -e
cd "$(dirname "$0")/.."

lsof -ti:8081 | xargs kill -9 2>/dev/null || true
lsof -ti:8545 | xargs kill -9 2>/dev/null || true
sleep 1

echo "=== Starting Hardhat node ==="
npx hardhat node --hostname 0.0.0.0 --port 8545 &>/tmp/hardhat_final.log &
HHPID=$!
sleep 3

echo "=== Deploying mock contracts ==="
npx hardhat run scripts/deploy_mock.js --network localhost

echo "=== Starting backend ==="
set -a && . ./.env && set +a
cd well-backend
go run ./cmd/ &>/tmp/backend_final.log &
BEPID=$!
cd ..
sleep 3

echo "=== Running e2e test ==="
set +e
npx hardhat run scripts/e2e_test.js --network localhost
E2E_EXIT=$?
set -e

echo "=== Cleaning up ==="
kill $BEPID 2>/dev/null
kill $HHPID 2>/dev/null

exit $E2E_EXIT
