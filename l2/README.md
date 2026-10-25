# L2 evaluation suite (Optimism devnet)

Everything needed to run the thesis benchmark suite against a real local Layer 2 —
an [Optimism OP Stack](https://github.com/ethereum-optimism/optimism) devnet
provisioned with `ops-bedrock`'s own L1 — so the thesis can show *"first
everything on L1, then add L2 and here are the data"*.

Everything added for L2 lives in this `l2/` folder. The repo root stays clean:
the Optimism checkout lives **outside** the repo at `../optimism`, and all
generated state (`.env-l1`, `.env-l2`, `results/l2/`) is gitignored.

## Architecture

```
                     ops-bedrock devnet (`make devnet`)
  ┌──────────────────────────────────────────────────────────┐
  │  L1 geth (devnet, chainid 900)          L2 op-geth (901) │
  │  http://127.0.0.1:8545   ws://8546     http://127.0.0.1:9545   ws://9546 │
  │        │  ▲                                     │  ▲                  │
  │        │  └── OptimismPortal  ←──── deposit ────┘  │ (op sequencer)   │
  │        └──────┐  L2ToL1MessagePasser ──┐           │                  │
  └───────────────┼─────────────────────────┼───────────┼──────────────────┘
                  │                         │           │
        .env-l1 (L1 addresses)     .env-l2 (L2 addresses)
                  └──── well-backend (env-driven) ←───┘
                     bench_zk / bench_latency / bench_throughput
```

The backend is fully env-driven (`well-backend/internal/config`), so switching a
benchmark run between L1 baseline and L2 is purely: source `.env`/`.env-l2`,
point `ETH_RPC` at the right WebSocket, run the *same* bench scripts.

## Quickstart

```bash
# 1. provision the devnet (Kurtosis `just devnet` on current optimism; the
#    legacy `make devnet` uses L1 8545 + L2 9545, chain ids 900/901)
bash l2/setup-devnet.sh

# 2. deploy the Well contracts to BOTH layers (.env-l1, .env-l2)
bash l2/deploy-layers.sh

# 3. top up hardhat signers on L2 via the real bridge (idempotent)
node l2/fund-accounts.js

# 4. loopback evaluation on the L2
bash l2/run-eval.sh

# 5. latency-injected (WAN) evaluation on the L2 (~2h full matrix)
bash l2/run-wan-eval.sh

# 6. cross-layer bridge + app-metric demo
node l2/bridge-demo.js
```

## Endpoint configuration

Defaults target the legacy docker-compose devnet layout
(L1 `http://127.0.0.1:8545`/`ws://127.0.0.1:8546`, chain id 900;
L2 `http://127.0.0.1:9545`/`ws://127.0.0.1:9546`, chain id 901).
`setup-devnet.sh` probes what is actually reachable and prints the chain ids.
If your devnet differs, export the four endpoint variables (+ optionally
`L1_CHAIN_ID`/`L2_CHAIN_ID`) before running the eval scripts; `l2_check()`
only requires both layers to be reachable and to be *different* chains.

`DEV_KEY` (private key) is honored when the devnet funds a key other than the
well-known hardhat mnemonic accounts: it becomes the deployer/funder on L1,
while the bench signers remain the mnemonic accounts (topped up via the bridge).

## Files

| file | purpose |
|------|---------|
| `setup-devnet.sh` | clone `ethereum-optimism/optimism` to `../optimism`, run `just devnet` (Kurtosis) or `make devnet`, probe what came up |
| `deploy-layers.sh` | deploy Well contracts to devnet L1 + L2 via `scripts/deploy_mock.js`; writes `.env-l1` / `.env-l2` (gitignored) |
| `env.sh` | shared config: endpoints, chain-id sanity check, postgres-only reset, L2 deploy/fund helpers |
| `fund-accounts.js` | top up the 10 hardhat signers on L2 through the L1→L2 bridge (real cross-chain balance) |
| `bridge-demo.js` | full bridge round-trip: L1→L2 deposit, L2→L1 withdrawal initiation, plus same-call L1-vs-L2 metrics |
| `lib-bridge.js` | shared bridge helpers (portal/bridge detection, legacy `depositTransaction` probes) |
| `run-eval.sh` | loopback suite on L2 → `scripts/benchmarks/results/l2/` |
| `run-wan-eval.sh` | latency-injected suite on L2 → `scripts/benchmarks/results/l2/wan/` |

## What the results mean (honest framing)

- Bench scripts are byte-identical to the L1 baseline; all differences in the
  data are real property differences of the chain being exercised.
- `gasUsed` reported by op-geth is **pure EVM gas** (the Well contracts are
  identical artifacts on both layers). The L2 *total fee* additionally includes
  the L1 data fee, which op-geth bills separately and is **not** part of
  `gasUsed` — compare gas units / latency / throughput, do not compare fees
  between layers.
- Block times and the sequencer pipeline make L2 confirmation latency different
  from the Clique baseline; expect that in the numbers.
- This is a *local* devnet: absolute fees are synthetic (gas price
  `--miner.gasprice 0` on L1, devnet L1 base-fee oracle on L2). Scope cross-layer
  claims to mechanics and gas units, not fee economics.

## Adding an L2 run to the thesis

See `docs/benchmark-methodology-l2.md` for the methodology write-up and the
paragraphs to slot into `docs/benchmarks.md` for the L1-then-L2 narrative.