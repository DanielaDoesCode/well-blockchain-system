#!/bin/bash
set -e
cd "$(dirname "$0")/.."

echo "=== Starting infrastructure (postgres/ipfs/vault) ==="
docker compose up -d postgres ipfs vault
sleep 5

lsof -ti:8545 | xargs kill -9 2>/dev/null || true
lsof -ti:8081 | xargs kill -9 2>/dev/null || true

echo "=== Starting Hardhat node ==="
npx hardhat node --hostname 0.0.0.0 --port 8545 &>/tmp/hardhat_emergency.log &
HHPID=$!
sleep 3

echo "=== Deploying mock contracts ==="
npx hardhat run scripts/deploy_mock.js --network localhost

echo "=== Starting backend (listener) ==="
set -a && . ./.env && set +a
cd well-backend
go run ./cmd/ &>/tmp/backend_emergency.log &
BEPID=$!
cd ..
sleep 3

echo "=== Running e2e emergency test ==="
set +e
npx hardhat run scripts/e2e_emergency.js --network localhost
E2E_EXIT=$?
set -e

echo "=== Verifying listener captured per-patient EmergencySummary events ==="
FOUND=0
for i in $(seq 1 15); do
  if grep -q "EmergencySummary.AccessGranted" /tmp/backend_emergency.log; then
    FOUND=1
    break
  fi
  sleep 1
done
if [ "$FOUND" -ne 1 ]; then
  echo "FAIL: listener did not capture EmergencySummary events"
  tail -30 /tmp/backend_emergency.log
  kill $BEPID 2>/dev/null
  kill $HHPID 2>/dev/null
  exit 1
fi
grep -q "subscribed: EmergencySummary(" /tmp/backend_emergency.log \
  && echo "  [PASS] listener subscribed to per-patient EmergencySummary" \
  || echo "  [WARN] dynamic EmergencySummary subscription not seen in log"
echo "  [PASS] listener captured EmergencySummary.AccessGranted"

echo "=== Cleaning up ==="
kill $BEPID 2>/dev/null
kill $HHPID 2>/dev/null

exit $E2E_EXIT
