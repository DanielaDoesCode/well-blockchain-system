#!/usr/bin/env bash
set -euo pipefail

# Initialize the Geth (Clique PoA) data directory and import the 10 well-known
# Hardhat account keys so the node can seal blocks and sign transactions for
# all pre-funded accounts.
#
# The AWS plan's docker-compose passes --unlock to the Geth node, but a freshly
# `geth init` datadir contains no keystore for those accounts, so the node would
# fail to seal blocks and eth_sendTransaction would reject the hardhat scripts.
# This script fixes that by importing the well-known Hardhat mnemonic keys
# (same keys forge/hardhat deploy with), creating the keystore entries Geth
# needs. Addresses must stay in sync with genesis.json's alloc.
#
# Usage: scripts/init-geth.sh
# Idempotent: safe to run any time, even while the Geth node is up (the init
# step is skipped if the datadir already has chaindata; keystore imports do not
# touch the chaindata lock).

COMPOSE_FILE="docker-compose.geth.yml"
DATADIR="/data"
GENESIS="/genesis.json"
PASSWORD="/data/password.txt"
GETH_IMAGE="ethereum/client-go:v1.13.15"
GETH_VOLUME="well-blockchain_geth_data"

# Hardhat mnemonic "test ... junk" accounts (private key / address, no 0x).
ACCOUNTS=(
  "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80|f39fd6e51aad88f6f4ce6ab8827279cfffb92266"
  "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d|70997970c51812dc3a010c7d01b50e0d17dc79c8"
  "5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a|3c44cdddb6a900fa2b585dd299e03d12fa4293bc"
  "7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6|90f79bf6eb2c4f870365e785982e1f101e93b906"
  "47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a|15d34aaf54267db7d7c367839aaf71a00a2c6a65"
  "8b3a350cf5c34c9194ca85829a2df0ec3153be0318b5e2d3348e872092edffba|9965507d1a55bcc2695c58ba16fb37d819b0a4dc"
  "92db14e403b83dfe3df233f83dfa3a0d7096f21ca9b0d6d6b8d88b2b4ec1564e|976ea74026e726554db657fa54763abd0c3a0aa9"
  "4bbbf85ce3377467afe5d46f804f221813b2bb87f24d81f60f1fcdbf7cbf4356|14dc79964da2c08b23698b3d3cc7ca32193d9955"
  "dbda1821b80551c9d65939329250298aa3472ba22feea921c0cf5d620ea67b97|23618e81e3f5cdf7f54c3d65f7fbc0abf5b21e8f"
  "2a871d0798f97d79848a013d4936a73bf4cc922c825d33c1cf7073dff6d409c6|a0ee7a142d267c1f36714e4a8f75612f20a79720"
)

cd "$(dirname "$0")/.."

run_geth() {
  docker compose -f "$COMPOSE_FILE" run --rm geth "$@"
}

# Import a single key into the geth keystore using the plain docker CLI with
# explicit mounts (compose "run" does not take "docker run -v" style flags).
run_import() {
  docker run --rm \
    -v "$GETH_VOLUME":/data \
    -v "$(pwd)/password.txt":/data/password.txt:ro \
    -v /tmp/geth-import.key:/data/import.key:ro \
    "$GETH_IMAGE" \
    account import --datadir "$DATADIR" --password "$PASSWORD" /data/import.key
}

echo "=== [1/2] Initializing Geth datadir with genesis (idempotent) ==="
if docker compose -f "$COMPOSE_FILE" ps --status running --services 2>/dev/null | grep -qx geth; then
  echo "  Geth node is running — datadir already initialized, skipping init."
elif run_geth sh -c 'test -d /data/geth/chaindata' 2>/dev/null; then
  echo "  Datadir already initialized, skipping init."
else
  run_geth init --datadir "$DATADIR" "$GENESIS"
fi

echo "=== [2/2] Importing Hardhat account keys (idempotent) ==="
LISTED="$(run_geth geth account list --datadir "$DATADIR" 2>/dev/null || true)"
for entry in "${ACCOUNTS[@]}"; do
  KEY="${entry%%|*}"
  ADDR="${entry##*|}"
  if echo "$LISTED" | grep -qi "$ADDR"; then
    echo "  $ADDR already in keystore — skipping."
  else
    printf '%s' "$KEY" > /tmp/geth-import.key
    if run_import >/dev/null 2>&1; then
      echo "  Imported $ADDR"
    else
      echo "  FAILED to import $ADDR" >&2
      exit 1
    fi
    rm -f /tmp/geth-import.key
  fi
done

echo "=== Geth datadir ready ==="
run_geth geth account list --datadir "$DATADIR" 2>/dev/null | grep 'Account #' || true
