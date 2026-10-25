#!/usr/bin/env bash
set -euo pipefail

# Provision + start a local OP Stack L2 + its L1.
#
# The checkout lives OUTSIDE the thesis repo (../optimism) so the repo stays
# clean; only this helper and the eval scripts live in l2/.
#
# Two devnet flavors exist:
#   - `just devnet` (Kurtosis)  — the current optimizer monorepo default
#   - `make devnet` (docker-compose) — the legacy ops-bedrock layout:
#       L1 geth (chainid 900): http://127.0.0.1:8545   ws://127.0.0.1:8546
#       L2 op-geth (chainid 901): http://127.0.0.1:9545 ws://127.0.0.1:9546
#
# Prereqs: git, golang, Docker (compose plugin); for Kurtosis also `just`.
# The first start is heavy (downloads + build, ~20-40 min for make devnet).
#
# After startup the script probes common endpoints and prints which chain ids
# it found; export L1_RPC_HTTP/L1_RPC_WS/L2_RPC_HTTP/L2_RPC_WS (and optionally
# L1_CHAIN_ID/L2_CHAIN_ID) to tell the eval runners where the layers live.

cd "$(dirname "$0")/.."

OPTIMISM_DIR="${OPTIMISM_DIR:-$(cd . && pwd)/../optimism}"

if [ ! -d "$OPTIMISM_DIR/.git" ]; then
  echo "Cloning ethereum-optimism/optimism -> $OPTIMISM_DIR ..."
  git clone https://github.com/ethereum-optimism/optimism.git "$OPTIMISM_DIR"
fi

cd "$OPTIMISM_DIR"
echo "=== optimism @ $(git rev-parse --short HEAD 2>/dev/null || echo unknown) on $(git branch --show-current 2>/dev/null || echo '(detached)') ==="

for cmd in docker go git; do
  command -v "$cmd" >/dev/null || { echo "ERROR: '$cmd' required"; exit 1; }
done
docker compose version >/dev/null 2>&1 || { echo "ERROR: docker compose plugin required"; exit 1; }

# Newer optimism monorepos drive the devnet via `just` (Kurtosis, the current
# default) or the legacy docker-compose `make devnet`. Try both; set
# DEVNET_ARGS to pass arguments to just (e.g. DEVNET_ARGS="simple.yaml").
if command -v just >/dev/null 2>&1 && [ -f justfile ]; then
  echo "=== just devnet ${DEVNET_ARGS:-} (first run downloads/boots the Kurtosis devnet, several minutes) ==="
  just devnet ${DEVNET_ARGS:-}
elif [ -f Makefile ]; then
  echo "=== make devnet (legacy docker-compose devnet; first run builds op-node/op-geth/bedrock, ~20-40 min) ==="
  make devnet
else
  echo "ERROR: no devnet entrypoint found in $OPTIMISM_DIR" >&2
  exit 1
fi

PROBE() { curl -s -m 2 -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' "$1" \
  | node -e 'let s="";process.stdin.on("data",d=>s+=d);process.stdin.on("end",()=>{try{console.log(parseInt(JSON.parse(s).result,16))}catch(e){console.log("")}})'; }

echo
echo "=== probing common devnet endpoints ==="
# The ops-bedrock compose devnet serves L1 on 8545 and L2 on 9545 (ids 900/901);
# Kurtosis devnets use other ports/ids. Anything reachable with a chain id is
# reported; export the right endpoints + ids for the eval runners.
for port in 8545 8546 8547 9545 9546 9222 8551; do
  c=$(PROBE "http://127.0.0.1:$port" 2>/dev/null)
  [ -n "$c" ] && echo "  :$port  chainId=$c"
done

c1=$(PROBE http://127.0.0.1:8545)
c2=$(PROBE http://127.0.0.1:9545)
if [ -n "$c1" ] && [ -n "$c2" ] && [ "$c1" != "$c2" ]; then
  echo "  -> found L1 (chainid $c1) on :8545 and L2 (chainid $c2) on :9545 (the legacy-compose layout)"
  echo
  echo "=== devnet ready (legacy-compose layout) ==="
  echo "  Next: bash l2/deploy-layers.sh && bash l2/run-eval.sh"
else
  echo "  -> devnet is up but not on the assumed :8545/:9545 ports."
  echo "     Set the endpoints explicitly and re-check, e.g.:"
  echo "       export L1_RPC_HTTP=http://127.0.0.1:<l1>  L1_RPC_WS=ws://127.0.0.1:<l1ws>"
  echo "       export L2_RPC_HTTP=http://127.0.0.1:<l2>  L2_RPC_WS=ws://127.0.0.1:<l2ws>"
  echo "       bash l2/run-eval.sh"
fi