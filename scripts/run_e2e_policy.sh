#!/bin/bash
set -e
cd "$(dirname "$0")/.."

echo "=== Starting infrastructure (postgres/ipfs/vault) ==="
docker compose up -d postgres ipfs vault
sleep 5

lsof -ti:8545 | xargs kill -9 2>/dev/null || true
lsof -ti:8081 | xargs kill -9 2>/dev/null || true

echo "=== Starting Hardhat node ==="
npx hardhat node --hostname 0.0.0.0 --port 8545 &>/tmp/hardhat_policy.log &
HHPID=$!
sleep 3

echo "=== Deploying mock contracts (real PolicyVerifier) ==="
npx hardhat run scripts/deploy_mock.js --network localhost

echo "=== Starting backend ==="
set -a && . ./.env && set +a
cd well-backend
go run ./cmd/ &>/tmp/backend_policy.log &
BEPID=$!
cd ..
sleep 3

echo "=== Running e2e policy test ==="
set +e
npx hardhat run scripts/e2e_policy.js --network localhost
E2E_EXIT=$?
set -e

echo "=== Cleaning up ==="
kill $BEPID 2>/dev/null
kill $HHPID 2>/dev/null

exit $E2E_EXIT
