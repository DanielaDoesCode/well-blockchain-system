#!/usr/bin/env python3
"""Figure how a concurrent session burst maps onto block gas capacity.

Reads a run's `burst_blocks.json` (block composition of the burst, plus the
per-session ingestion lag) and plots, for each block the burst occupied:

  * gas consumed against the block gas limit, with the limit drawn as a
    reference so the reader can see immediately whether the burst saturated
    the block or merely fitted inside it
  * the count of sessions in that block
  * the ingestion lag of the sessions in that block

The point of the figure is that the *distribution* of lag is not a property of
the workload alone: it inherits the block structure of the burst. Sessions that
land in an earlier block and sessions that land in a later one are ingested in
different polls, so the lag distribution inherits the block split. That is only
visible when the two are plotted together.

Usage:
  plot_burst_blocks.py --blocks <run-dir>/burst_blocks.json --out <png>
"""
import argparse
import json
import os

import matplotlib

matplotlib.use("Agg")
import matplotlib.pyplot as plt


def load(path):
    with open(path) as fh:
        return json.load(fh)


def main():
    ap = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    ap.add_argument("--blocks", required=True, help="path to burst_blocks.json")
    ap.add_argument("--ingest", default=None,
                    help="path to session_ingest.json (default: alongside --blocks)")
    ap.add_argument("--out", required=True, help="output PNG path")
    ap.add_argument("--title", default=None)
    args = ap.parse_args()

    d = load(args.blocks)
    ingest_path = args.ingest or os.path.join(
        os.path.dirname(os.path.abspath(args.blocks)), "session_ingest.json")
    lags = load(ingest_path)["ingestLagSortedMs"]

    limit = d["gasLimit"]
    blocks = d["blocks"]
    labels = [f"block {b['number']}" for b in blocks]
    used = [b["gasUsed"] for b in blocks]
    txs = [b["txs"] for b in blocks]

    fig, (ax1, ax2) = plt.subplots(
        1, 2, figsize=(11.5, 4.8), constrained_layout=True, width_ratios=[1, 1.15])

    # -- gas occupancy per block, against the limit ------------------------
    x = range(len(blocks))
    ax1.bar(x, used, width=0.55, color="#4c72b0", edgecolor="white",
            label="gas consumed by burst sessions")
    ax1.axhline(limit, color="#c44e52", lw=1.4, ls="--",
                label=f"block gas limit ({limit / 1e6:.0f}M)")
    for i, (u, t) in enumerate(zip(used, txs)):
        ax1.annotate(f"{u:,}\n{t} sessions\n{u / limit * 100:.1f}%",
                     (i, u), ha="center", va="bottom", fontsize=8,
                     xytext=(0, 3), textcoords="offset points")
    ax1.set_xticks(list(x))
    ax1.set_xticklabels(labels)
    ax1.set_ylim(0, limit * 1.18)
    ax1.set_ylabel("Gas used in block")
    ax1.set_title("Burst gas consumption vs the block gas limit", fontsize=10)
    ax1.legend(fontsize=7.5, loc="upper right", framealpha=0.9)
    ax1.grid(axis="y", alpha=0.25, lw=0.6)
    ax1.set_axisbelow(True)

    # -- ingestion lag: the two clusters -----------------------------------
    # NOTE: the seeder does not record which block each individual transaction
    # was mined in, only how many. The cluster sizes happen to equal the
    # per-block counts, but which specific sessions belong to which cluster is
    # an inference, so this panel deliberately does NOT colour sessions by
    # block. It shows the measured distribution and reports the block split as
    # context.
    ax2.scatter(range(len(lags)), lags, s=14, color="#4c72b0", alpha=0.8,
                edgecolor="white", linewidth=0.3, zorder=3)
    for p in (50, 95):
        ax2.axhline(d["ingestLagMs"][f"p{p}"], color="#555555", lw=0.7, ls="-.")
    ax2.set_xlabel("Sessions, ordered by ingestion lag (fastest first)")
    ax2.set_ylabel("Ingestion lag from block confirmation (ms)")
    split = " + ".join(str(b["txs"]) for b in blocks)
    ax2.set_title(
        f"Ingestion lag is bimodal; burst occupied {len(blocks)} blocks ({split})",
        fontsize=10)
    ax2.grid(axis="y", alpha=0.25, lw=0.6)
    ax2.set_axisbelow(True)
    ax2.set_xlim(-4, len(lags) + 4)
    ax2.annotate(
        f"clusters of {blocks[0]['txs']} and {blocks[-1]['txs']} sessions,\n"
        f"separated by ~1 block period",
        xy=(0.5, 0.06), xycoords="axes fraction", fontsize=8, ha="center",
        color="#555555")

    title = args.title or (
        f"Concurrent session burst of {d['sessions']} sessions mapped onto "
        f"block gas capacity ({d['blockPeriodSec']}s blocks, "
        f"{d['gasPerTx']:,} gas per authenticate transaction)")
    fig.suptitle(title, fontsize=11)
    fig.savefig(args.out, dpi=200)
    plt.close(fig)

    os.makedirs(os.path.dirname(os.path.abspath(args.out)), exist_ok=True)
    print(f"wrote:  {args.out}")
    print(f"ceiling: {d['sessionsPerBlockCeiling']} sessions/block "
          f"at {limit:,} gas limit")
    for b in blocks:
        print(f"  block {b['number']}: {b['txs']} txs, {b['gasUsed']:,} gas "
              f"({b['utilisation'] * 100:.1f}%)")


if __name__ == "__main__":
    main()
