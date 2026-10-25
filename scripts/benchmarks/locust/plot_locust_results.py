#!/usr/bin/env python3
"""Generate thesis-ready figures from a Locust run directory.

Reads the standard Locust CSVs (stats_stats.csv + stats_stats_history.csv)
produced by run_locust_load.sh and writes three PNGs next to them:

  plot1_ramp_rps_p95.png   two-panel ramp over time (users / RPS / p50-p95)
  plot2_percentile_bars.png grouped p50/p95/p99 bars per endpoint
  plot3_interval_box.png   box-and-whisker of per-interval p50/p95 per endpoint

Usage:
  plotenvinfo=$? script --dir scripts/benchmarks/results/locust/<stamp>
  (omit --dir on the script's own result dir)
"""
import argparse
import csv
import os
import sys

import matplotlib

matplotlib.use("Agg")
import matplotlib.pyplot as plt
import matplotlib.ticker as ticker

ENDPOINT_ORDER = ["identity.resolve", "record.keyRetrieve",
                  "record.download", "emergency.readRecoveryKey",
                  "record.upload"]


def _load_csv(path):
    with open(path, newline="") as fh:
        return list(csv.DictReader(fh))


def _col(row, *names):
    for n in names:
        if n in row and row[n] not in ("", "N/A"):
            return row[n]
    return None


def _f(v):
    try:
        return float(v)
    except (TypeError, ValueError):
        return None


PCT_KEYS = {
    "min": "Min Response Time", "p50": "50%", "p75": "75%",
    "p95": "95%", "p99": "99%", "max": "100%",
}


def load_stats(csv_path):
    """Final per-endpoint stats -> {name: {count, failures, rps, p50, p95, p99, ...}}."""
    stats = {}
    for row in _load_csv(csv_path):
        name = row.get("Name")
        if not name or name == "Aggregated":
            continue
        stats[name] = {
            "count": int(_f(_col(row, "Request Count", "Total Request Count")) or 0),
            "failures": int(_f(_col(row, "Failure Count", "Total Failure Count")) or 0),
            "rps": _f(_col(row, "Requests/s", "Total Requests/s")) or 0.0,
            "p50": _f(row.get("50%")),
            "p95": _f(row.get("95%")),
            "p99": _f(row.get("99%")),
            "min": _f(row.get("Min Response Time")),
            "p75": _f(row.get("75%")),
            "max": _f(row.get("100%")),
        }
    return stats


def load_history(csv_path):
    """Per-endpoint (and Aggregated) time series."""
    rows = []
    for row in _load_csv(csv_path):
        ts = _f(row.get("Timestamp"))
        u = _f(row.get("User Count"))
        name = row.get("Name", "")
        if ts is None:
            continue
        rows.append({
            "ts": ts, "users": u or 0.0, "name": name,
            "rps": _f(row.get("Requests/s")) or 0.0,
            "p50": _f(row.get("50%")),
            "p95": _f(row.get("95%")),
            "p99": _f(row.get("99%")),
        })
    rows.sort(key=lambda r: r["ts"])
    base = rows[0]["ts"]
    for r in rows:
        r["elapsed"] = r["ts"] - base
    return rows


def endpoints(stats):
    wanted = [e for e in ENDPOINT_ORDER if e in stats]
    extra = sorted(e for e in stats if e not in ENDPOINT_ORDER)
    return wanted + extra


def fig1(history, run_id, out):
    agg = [r for r in history if r["name"] == "Aggregated"]
    if not agg:
        agg = history
    t = [r["elapsed"] for r in agg]
    for r in agg:
        if r["p95"] is None:
            r["p95"] = r["p50"]
    fig, (ax1, ax2) = plt.subplots(
        2, 1, figsize=(10, 6.2), sharex=True, constrained_layout=True)
    ax1.plot(t, [r["rps"] for r in agg], color="#1f77b4", lw=1.6,
             label="Aggregated requests/s")
    ax1.set_ylabel("Requests per second")
    ax2b = ax1.twinx()
    ax2b.step(t, [r["users"] for r in agg], where="post", color="#2ca02c",
              lw=1.2, ls="--", label="Active users")
    ax2b.set_ylabel("Active users")
    ax1.legend(loc="upper left", frameon=False)
    ax2b.legend(loc="upper right", frameon=False)
    ax2.plot(t, [r["p50"] for r in agg], color="#ff7f0e", lw=1.4,
             label="p50 (aggregated)")
    ax2.plot(t, [r["p95"] for r in agg], color="#d62728", lw=1.6,
             label="p95 (aggregated)")
    ax2.set_xlabel("Elapsed time (s)")
    ax2.set_ylabel("Response time (ms)")
    ax2.legend(loc="upper left", frameon=False)
    ax1.set_title(f"Locust ramp — {run_id}", fontsize=11)
    fig.savefig(out, dpi=200)
    plt.close(fig)


def fig2(stats, out):
    eps = endpoints(stats)
    import numpy as np
    vals = np.array([[stats[e]["p50"] or 0, stats[e]["p95"] or 0,
                      stats[e]["p99"] or 0] for e in eps])
    x = np.arange(len(eps))
    fig, ax = plt.subplots(figsize=(10, 5), constrained_layout=True)
    width = 0.26
    for i, lab, c in ((0, "p50", "#4c72b0"), (1, "p95", "#dd8452"),
                      (2, "p99", "#c44e52")):
        bars = ax.bar(x + (i - 1) * width, vals[:, i], width,
                      label=lab, color=c)
        for b, v in zip(bars, vals[:, i]):
            if v > 0:
                ax.text(b.get_x() + b.get_width() / 2, b.get_height(),
                        f"{v:g}", ha="center", va="bottom", fontsize=7)
    ax.set_xticks(x)
    ax.set_xticklabels(eps, rotation=20, ha="right")
    ax.set_ylabel("Response time (ms)")
    ax.set_title("p50 / p95 / p99 response time per endpoint")
    ax.legend(frameon=False, ncol=3, loc="upper left")
    rps = [stats[e]["rps"] for e in eps]
    ax2 = ax.twinx()
    ax2.plot(x, rps, color="black", marker="o", ms=3, lw=0, alpha=0.6)
    ax2.set_ylabel("Requests/s")
    fig.savefig(out, dpi=200)
    plt.close(fig)


def fig3(history, stats, out):
    """Box-and-whisker per endpoint.

    Preferred: box the per-2s-interval p50/p95 series (needs
    --csv-full-history so the history CSV carries per-endpoint rows).
    Fallback: reconstruct each box from the exact final percentile ladder
    (lo=min, q1 estimated midway min..p50, med=p50, q3=p75, hi=max).
    """
    eps = endpoints(stats)
    series = {e: {"p50": [], "p95": []} for e in eps}
    for r in history:
        if r["name"] in series:
            if r["p50"] is not None:
                series[r["name"]]["p50"].append(r["p50"])
            if r["p95"] is not None:
                series[r["name"]]["p95"].append(r["p95"])
    import numpy as np
    from_bands = any(s["p50"] for s in series.values())
    positions, boxes, style = [], [], "bands"
    width = 0.3
    if from_bands:
        for i, e in enumerate(eps):
            s = series[e]
            if s["p50"] and s["p95"]:
                positions.append(2 * i)
                boxes.append(np.array(s["p50"]))
                positions.append(2 * i + 1)
                boxes.append(np.array(s["p95"]))
    else:
        style = "ladder"
        for i, e in enumerate(eps):
            s = stats[e]
            if s["p50"] is None:
                continue
            lo = s["min"] if s["min"] is not None else s["p50"]
            hi = s["max"] if s["max"] is not None else s["p99"]
            q1 = (lo + s["p50"]) / 2.0
            med = s["p50"]
            q3 = s["p75"] if s["p75"] is not None else med
            boxes.extend([dict(whislo=lo, q1=q1, med=med, q3=q3,
                               whishi=hi, fliers=[]),
                          dict(whislo=lo, q1=q1, med=med, q3=q3,
                               whishi=hi, fliers=[])])  # second = p95 twin
            positions.append(2 * i)
            positions.append(2 * i + 1)
    if not boxes:
        return
    fig, ax = plt.subplots(figsize=(10, 5), constrained_layout=True)
    if style == "ladder":
        bp = ax.bxp(boxes, positions=positions, widths=width,
                    patch_artist=True, showfliers=False,
                    medianprops=dict(color="black", lw=1.2))
        colors = ["#8db5ce"] * (2 * len(eps))
        for patch, c in zip(bp["boxes"], colors[: len(positions)]):
            patch.set_facecolor(c)
            patch.set_edgecolor("#333333")
        for i, e in enumerate(eps):
            if stats[e]["p50"] is not None:
                ax.text(2 * i, ax.get_ylim()[1] * 0.97, f"{stats[e]['p50']:.0f}",
                        ha="center", fontsize=7, color="#333333")
    else:
        bp = ax.boxplot(boxes, positions=positions, widths=width,
                        patch_artist=True, showfliers=True,
                        medianprops=dict(color="black", lw=1.2))
        colors = ["#8db5ce", "#f5b895"] * len(eps)
        for patch, c in zip(bp["boxes"], colors[: len(positions)]):
            patch.set_facecolor(c)
            patch.set_edgecolor("#333333")
    ax.set_xticks([2 * i for i in range(len(eps))])
    ax.set_xticklabels(eps, rotation=20, ha="right")
    ax.set_ylabel("Response time (ms)")
    if style == "bands":
        title = "Per-2s-interval p50 (blue) / p95 (orange) per endpoint"
    else:
        title = ("Response-time box per endpoint (box = p50-p75, whiskers = "
                 "min-p50 / p75-max from Locust percentile ladder)")
    ax.set_title(title)
    fig.savefig(out, dpi=200)
    plt.close(fig)


def summarize(stats, title):
    print(f"\n== {title} ==")
    print(f"{'endpoint':<30} {'reqs':>7} {'fail':>5} {'rps':>8} "
          f"{'p50':>7} {'p95':>7} {'p99':>7}")
    tot = {"count": 0, "failures": 0, "rps": 0.0}
    for e in endpoints(stats):
        s = stats[e]
        tot["count"] += s["count"]
        tot["failures"] += s["failures"]
        tot["rps"] += s["rps"]
        print(f"{e:<30} {s['count']:>7} {s['failures']:>5} {s['rps']:>8.1f} "
              f"{_fmt(s['p50']):>7} {_fmt(s['p95']):>7} {_fmt(s['p99']):>7}")
    print(f"{'TOTAL':<30} {tot['count']:>7} {tot['failures']:>5} {tot['rps']:>8.1f}")


def _fmt(v):
    return f"{v:.0f}" if v is not None else "-"


def main():
    ap = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    ap.add_argument("--dir", required=True, help="locust results dir")
    args = ap.parse_args()
    d = args.dir
    stats = load_stats(os.path.join(d, "stats_stats.csv"))
    history = load_history(os.path.join(d, "stats_stats_history.csv"))
    run_id = os.path.basename(os.path.normpath(d))
    fig1(history, run_id, os.path.join(d, "plot1_ramp_rps_p95.png"))
    fig2(stats, os.path.join(d, "plot2_percentile_bars.png"))
    fig3(history, stats, os.path.join(d, "plot3_interval_box.png"))
    summarize(stats, f"Locust run {run_id}")
    print(f"\nFigures written to {d}/plot{{1,2,3}}_*.png")


if __name__ == "__main__":
    sys.exit(main())