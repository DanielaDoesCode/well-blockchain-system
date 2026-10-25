#!/usr/bin/env python3
"""Plot the listener-ingestion lag distribution for a session-burst run.

Reads a run's `session_ingest.json` (written by seed_locust.js when
SESSION_CONCURRENCY > 1) and figures the time from a session's own block
confirmation to the first API response that accepts its token.

Two modes, chosen automatically by what the artifact contains:

  full   -- `ingestLagSortedMs` is present (every session's measured lag).
            Draws a histogram plus an ECDF. This is the mode to prefer.
  sample -- only the summary percentiles survive (older runs, recorded before
            the raw sample was captured). Draws a percentile curve instead. It
            shows the shape of the distribution and is NOT a histogram; do not
            present it as one.

Reference lines drawn on both:
  * the 2s Clique block period, which is what the ingestion tail clusters near
  * the 3s fixed sleep the seeder historically used, i.e. the provisioning
    assumption this measurement exists to test

Usage:
  plot_ingest_lag.py --ingest <run-dir>/session_ingest.json --out <png> [--title ...]
"""
import argparse
import json
import os

import matplotlib

matplotlib.use("Agg")
import matplotlib.pyplot as plt
import numpy as np

BLOCK_PERIOD_MS = 2000.0
LEGACY_SLEEP_MS = 3000.0
# Seeder knob INGEST_POLL_MS: lag is quantised to this interval, so a measured
# value is really an interval [lag - poll, lag].
POLL_MS = 200.0


def load(path):
    with open(path) as fh:
        return json.load(fh)


def annotate_refs(ax, ymax):
    ax.axhline(BLOCK_PERIOD_MS, color="#444444", lw=1.0, ls="--", zorder=1)
    ax.annotate("2 s block period", xy=(0.005, BLOCK_PERIOD_MS), xycoords=("axes fraction", "data"),
                va="bottom", ha="left", fontsize=7.5, color="#444444")
    if LEGACY_SLEEP_MS <= ymax * 1.12:
        ax.axhline(LEGACY_SLEEP_MS, color="#c44e52", lw=1.0, ls=":", zorder=1)
        ax.annotate("3 s historical seed sleep", xy=(0.005, LEGACY_SLEEP_MS),
                    xycoords=("axes fraction", "data"),
                    va="bottom", ha="left", fontsize=7.5, color="#c44e52")


def plot_full(lags, d, out, title):
    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(11, 4.6), constrained_layout=True)

    # Histogram. The two clusters are the finding, so bin on a boundary that
    # separates them rather than letting numpy pick a uniform width.
    edges = np.arange(0, max(lags) + 200, POLL_MS)
    ax1.hist(lags, bins=edges, color="#4c72b0", edgecolor="white", linewidth=0.4)
    ax1.axvline(BLOCK_PERIOD_MS, color="#444444", lw=1.0, ls="--")
    ax1.axvline(LEGACY_SLEEP_MS, color="#c44e52", lw=1.0, ls=":")
    ax1.set_xlabel("Ingestion lag from block confirmation (ms)")
    ax1.set_ylabel("Sessions")
    ax1.set_title(f"Distribution (n={len(lags)}, bin={POLL_MS:.0f} ms)", fontsize=10)
    ax1.annotate("fast cluster", xy=(0.05, 0.92), xycoords="axes fraction",
                 fontsize=8, color="#1f4e79")
    ax1.annotate("cluster near\nblock period", xy=(0.62, 0.92), xycoords="axes fraction",
                 fontsize=8, color="#1f4e79")

    x = np.arange(1, len(lags) + 1)
    ax2.plot(x, lags, color="#1f77b4", lw=1.4)
    ax2.axhline(BLOCK_PERIOD_MS, color="#444444", lw=1.0, ls="--")
    ax2.axhline(LEGACY_SLEEP_MS, color="#c44e52", lw=1.0, ls=":")
    for p in (50, 95, 99):
        v = d["ingestLagMs"][f"p{p}"]
        ax2.axhline(v, color="#8a8a8a", lw=0.6, ls="-.")
        ax2.annotate(f"p{p}={v:.0f} ms", xy=(0.62, v), fontsize=7, color="#555555",
                     va="bottom")
    ax2.set_xlabel("Sessions, ordered by lag (fastest first)")
    ax2.set_ylabel("Ingestion lag (ms)")
    ax2.set_title("ECDF", fontsize=10)
    ax2.set_xlim(1, len(lags))

    fig.suptitle(title, fontsize=11)
    fig.savefig(out, dpi=200)
    plt.close(fig)


def plot_sample(d, out, title):
    pts = [("min", d["ingestLagMs"]["min"]), ("p50", d["ingestLagMs"]["p50"]),
           ("p95", d["ingestLagMs"]["p95"]), ("p99", d["ingestLagMs"]["p99"]),
           ("max", d["ingestLagMs"]["max"])]
    xs = np.arange(len(pts))
    ys = [p[1] for p in pts]
    fig, ax = plt.subplots(figsize=(8.2, 4.8), constrained_layout=True)
    ax.plot(xs, ys, color="#1f77b4", lw=1.6, marker="o", ms=5)
    for x, (lab, v) in zip(xs, pts):
        ax.annotate(f"{v:.0f}", (x, v), textcoords="offset points", xytext=(0, 7),
                    ha="center", fontsize=8)
    ax.set_xticks(xs)
    ax.set_xticklabels([p[0] for p in pts])
    annotate_refs(ax, max(ys))
    # Emphasise the p50 -> p95 jump: the bimodality lives in that segment.
    ax.annotate("", xy=(2, ys[2]), xytext=(1, ys[1]),
                arrowprops=dict(arrowstyle="->", color="#c44e52", lw=1.4))
    ratio = ys[2] / ys[1] if ys[1] else float("nan")
    ax.annotate(f"{ratio:.1f}x jump —\nfast cluster to\nblock-period cluster",
                xy=(1.5, (ys[1] + ys[2]) / 2), fontsize=8, color="#c44e52",
                ha="center", va="center")
    ax.set_xlabel("Percentile of session time-to-usable")
    ax.set_ylabel("Ingestion lag from block confirmation (ms)")
    ax.set_title(title, fontsize=11)
    # These are thesis figures, so this states the figure's scope as a fact about
    # the data and nothing else. No instructions to the author belong in here;
    # the caption belongs in the document text, not in the rendered PNG.
    fig.text(0.5, -0.02,
             "Percentile curve: this run recorded summary percentiles only "
             "(no per-session sample).",
             ha="center", fontsize=7.5, style="italic", color="#666666")
    fig.savefig(out, dpi=200, bbox_inches="tight")
    plt.close(fig)


def main():
    ap = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    ap.add_argument("--ingest", required=True, help="path to session_ingest.json")
    ap.add_argument("--out", required=True, help="output PNG path")
    ap.add_argument("--title", default=None)
    args = ap.parse_args()

    d = load(args.ingest)
    n = d["sessions"]
    conc = d["concurrency"]
    default_title = (f"Listener ingestion lag — {n} concurrent sessions "
                     f"(N={conc}), poll interval {d['ingestPollMs']} ms")
    title = args.title or default_title

    if d.get("ingestLagSortedMs"):
        plot_full(d["ingestLagSortedMs"], d, args.out, title)
        mode = "full (histogram + ECDF)"
    else:
        plot_sample(d, args.out, title)
        mode = "sample (percentile curve — summary only)"

    os.makedirs(os.path.dirname(os.path.abspath(args.out)), exist_ok=True)
    print(f"mode:   {mode}")
    print(f"wrote:  {args.out}")
    lg = d["ingestLagMs"]
    print(f"lag ms: min={lg['min']:.0f} p50={lg['p50']:.0f} p95={lg['p95']:.0f} "
          f"p99={lg['p99']:.0f} max={lg['max']:.0f}  (p95/p50 = {lg['p95']/lg['p50']:.1f}x)")
    print(f"headroom vs {LEGACY_SLEEP_MS:.0f} ms legacy sleep: "
          f"{LEGACY_SLEEP_MS - lg['max']:.0f} ms "
          f"({(LEGACY_SLEEP_MS - lg['max']) / LEGACY_SLEEP_MS * 100:.0f}%)")


if __name__ == "__main__":
    main()
