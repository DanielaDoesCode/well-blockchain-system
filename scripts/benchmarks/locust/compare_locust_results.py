#!/usr/bin/env python3
"""Comparison figure across Locust runs (e.g. mock vs real-loopback vs relayed).

Reads each run's stats_stats.csv and writes a grouped bar chart of p50/p95
per endpoint, plus a table to stdout:

  compare_plot.py --dir RUN:LABEL [--dir RUN2:LABEL2 ...] --out out.png
"""
import argparse
import csv
import os
import sys

import matplotlib

matplotlib.use("Agg")
import matplotlib.pyplot as plt

ENDPOINT_ORDER = ["identity.resolve", "record.keyRetrieve",
                  "record.download", "emergency.readRecoveryKey",
                  "record.upload"]


def load(path):
    out = {}
    with open(path, newline="") as fh:
        for row in csv.DictReader(fh):
            name = row.get("Name")
            if not name or name == "Aggregated":
                continue
            out[name] = {
                "count": int(float(row.get("Request Count") or 0)),
                "failures": int(float(row.get("Failure Count") or 0)),
                "rps": float(row.get("Requests/s") or 0),
                "p50": float(row.get("50%") or 0),
                "p95": float(row.get("95%") or 0),
            }
    return out


def main():
    ap = argparse.ArgumentParser(description="Compare Locust runs side by side.")
    ap.add_argument("--dir", action="append", required=True,
                    metavar="RUN:LABEL", help="result dir + label (repeatable)")
    ap.add_argument("--out", required=True, help="output PNG path")
    args = ap.parse_args()

    runs = []
    for spec in args.dir:
        d, label = (spec.split(":", 1) + [spec])[:2] if ":" in spec else (spec, spec)
        runs.append((label, load(os.path.join(d, "stats_stats.csv"))))

    eps = [e for e in ENDPOINT_ORDER if all(e in s for _, s in runs)]
    eps += sorted({e for _, s in runs for e in s if e not in ENDPOINT_ORDER})

    import numpy as np
    fig, ax = plt.subplots(figsize=(11, 5.2), constrained_layout=True)
    n_runs = len(runs)
    bar_w = 0.32
    colors = ["#4c72b0", "#55a868", "#c44e52"]
    for ri, (label, s) in enumerate(runs):
        x = np.arange(len(eps)) + (ri - (n_runs - 1) / 2) * bar_w * 1.5
        p50 = [s[e]["p50"] for e in eps]
        p95 = [s[e]["p95"] for e in eps]
        width = bar_w / 2
        b1 = ax.bar(x - width / 2, p50, width, color=colors[ri], alpha=0.75,
                    label=f"{label} p50")
        b2 = ax.bar(x + width / 2, p95, width, color=colors[ri], alpha=0.45,
                    hatch="//", label=f"{label} p95")
        for b, v in zip(list(b1) + list(b2), p50 + p95):
            ax.text(b.get_x() + b.get_width() / 2, b.get_height(),
                    f"{v:g}", ha="center", va="bottom", fontsize=6)
    ax.set_xticks(np.arange(len(eps)))
    ax.set_xticklabels(eps, rotation=20, ha="right")
    ax.set_ylabel("Response time (ms)")
    ax.set_title("Locust comparison — p50 (solid) / p95 (hatched) per endpoint")
    ax.legend(frameon=False, ncol=min(n_runs * 2, 6), loc="upper left", fontsize=8)
    fig.savefig(args.out, dpi=200)
    plt.close(fig)

    print(f"{'endpoint':<28}" + "".join(
        f"{' | ' + l + ' p50':>16}{'p95':>7}{'rps':>7}{'fail':>6}" for l, _ in runs))
    for e in eps:
        row = f"{e:<28}"
        for _, s in runs:
            st = s.get(e)
            row += f"{st['p50']:>16.0f}{st['p95']:>7.0f}{st['rps']:>7.1f}{st['failures']:>6}"
        print(row)
    print(f"\nFigure written to {args.out}")


if __name__ == "__main__":
    sys.exit(main())