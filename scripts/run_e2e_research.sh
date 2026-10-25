#!/bin/bash
set -e
cd "$(dirname "$0")/.."

lsof -ti:8545 | xargs kill -9 2>/dev/null || true

echo "=== Starting Hardhat node ==="
npx hardhat node --hostname 0.0.0.0 --port 8545 &>/tmp/hardhat_research.log &
HHPID=$!
sleep 3

echo "=== Deploying mock contracts (real ResearchVerifier) ==="
npx hardhat run scripts/deploy_mock.js --network localhost

echo "=== Running e2e research test ==="
set +e
npx hardhat run scripts/e2e_research.js --network localhost
E2E_EXIT=$?
set -e

echo "=== Cleaning up ==="
kill $HHPID 2>/dev/null

exit $E2E_EXIT
