import "dotenv/config";
import { execFileSync } from "child_process";
import path from "path";
import os from "os";
import { fileURLToPath } from "url";
import { stat, printTable, writeJSON, SNAPSHOT } from "./lib.js";
import { precliBench, precliSetup } from "../lib/precli.js";

// AFGH primitive benchmark: the reference client (`precli bench`) times each
// operation in-process on BLS12-381 — no HTTP, no Postgres, no chain. This is
// the only place the cryptographic cost is measured on its own; the HTTP
// benchmark and the Locust runs fold it together with everything else and can
// only show its effect.
//
// Knobs: BENCH_PRE_REPS (timed iterations), BENCH_PRE_VERIFY (round-trip
// correctness checks), BENCH_PRE_SETUP_REPS (end-to-end client setup samples).

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const REPO_ROOT = path.join(__dirname, "..", "..");

const REPS = Number(process.env.BENCH_PRE_REPS || 1000);
const VERIFY = Number(process.env.BENCH_PRE_VERIFY || 500);
const SETUP_REPS = Number(process.env.BENCH_PRE_SETUP_REPS || 25);

const OPS = [
  ["keygen", "KeyGen (once per identity)"],
  ["seal", "Seal lockbox (patient, per record)"],
  ["deriveRk", "Derive re-encryption key (patient, per delegation)"],
  ["reencrypt", "Re-encrypt transform (backend, pure)"],
  ["reencrypt_endpoint", "Re-encrypt request path (parse + transform + marshal)"],
  ["selfDecrypt", "Patient self-recovery decrypt (paper D2)"],
  ["delegateeDecrypt", "Doctor decrypt of C_doctor (paper D1)"],
];

function buildPrecli() {
  if (process.env.PRECLI_BIN) return process.env.PRECLI_BIN;
  const bin = path.join(os.tmpdir(), `well-precli-${process.pid}`);
  execFileSync("go", ["build", "-o", bin, "./cmd/precli"], { cwd: path.join(REPO_ROOT, "well-backend") });
  process.env.PRECLI_BIN = bin;
  return bin;
}

async function main() {
  const goVersion = (() => {
    try { return execFileSync("go", ["version"], { encoding: "utf8" }).trim(); } catch { return "n/a"; }
  })();
  SNAPSHOT.provider = "not used (no chain, no HTTP, no database)";
  SNAPSHOT.geth = goVersion;
  SNAPSHOT.environment.pre = "etclab/pre AFGH second attempt (BLS12-381), CIRCL";

  console.log("=== PRE primitive benchmark (AFGH, BLS12-381) ===");
  buildPrecli();
  console.log(`binary: ${process.env.PRECLI_BIN}`);
  console.log(`reps per operation: ${REPS}, round-trip checks: ${VERIFY}\n`);

  const bench = precliBench(REPS, VERIFY);
  if (bench.correctness?.failures) {
    throw new Error(`${bench.correctness.failures} round-trips failed`);
  }

  const rows = OPS.map(([op, desc]) => {
    const s = bench.ops[op];
    return s
      ? [op, `${s.p50Ms.toFixed(3)} ms`, `${s.p95Ms.toFixed(3)} ms`,
        `${s.p99Ms.toFixed(3)} ms`, `${s.meanMs.toFixed(3)} ms`]
      : [op, "n/a", "n/a", "n/a", "n/a"];
  });
  console.log("Latency per operation:");
  printTable([["operation", "p50", "p95", "p99", "mean"], ...rows]);

  console.log("\nWire sizes:");
  printTable([
    ["value", "bytes"],
    ...Object.entries(bench.sizes).map(([k, v]) => [k, String(v)]),
  ]);

  // What a real client actually pays to prepare one record: one process
  // round-trip (or its in-process equivalent) that runs KeyGen for both
  // parties, seals the lockbox and derives the re-encryption key.
  const setupMs = [];
  const statePath = path.join(os.tmpdir(), `well-pre-client-setup-${process.pid}.json`);
  for (let i = 0; i < SETUP_REPS; i++) {
    const t0 = performance.now();
    precliSetup(statePath);
    setupMs.push(performance.now() - t0);
  }
  const setup = stat(setupMs);
  console.log(`\nEnd-to-end client setup (keygen x2 + seal + derive rk + spawn): `
    + `p50 ${setup.p50.toFixed(2)} ms, mean ${setup.mean.toFixed(2)} ms over ${SETUP_REPS} samples`);

  const sum = (a, b) => (bench.ops[a]?.p50Ms ?? 0) + (bench.ops[b]?.p50Ms ?? 0);
  const derived = {
    patientPerRecordMs: sum("seal", "deriveRk"),
    backendPerReencryptMs: bench.ops.reencrypt_endpoint?.p50Ms ?? null,
    patientSelfRecoverMs: bench.ops.selfDecrypt?.p50Ms ?? null,
    doctorReadMs: bench.ops.delegateeDecrypt?.p50Ms ?? null,
    identityKeygenMs: bench.ops.keygen?.p50Ms ?? null,
  };
  console.log("\nDerived per-unit costs (sum of p50s):");
  console.log(`  patient prepares one record: ${derived.patientPerRecordMs.toFixed(3)} ms (seal + derive rk)`);
  console.log(`  backend re-encrypts one record: ${derived.backendPerReencryptMs.toFixed(3)} ms`);
  console.log(`  patient self-recovery: ${derived.patientSelfRecoverMs.toFixed(3)} ms`);
  console.log(`  doctor opens a delegated record: ${derived.doctorReadMs.toFixed(3)} ms`);

  const outPath = writeJSON("pre.json", {
    snapshot: SNAPSHOT,
    unit: "milliseconds",
    reps: REPS,
    verify: VERIFY,
    ops: bench.ops,
    sizes: bench.sizes,
    correctness: bench.correctness,
    clientSetup: { samples: SETUP_REPS, ...setup },
    derived,
    notes: [
      "reencrypt is the pure AFGH transform; reencrypt_endpoint additionally parses and marshals the wire form, which is what handleReEncrypt spends.",
      "clientSetup includes the reference client's process spawn, so it is an upper bound on an in-process client.",
      "SealRandom generates its own message, so `seal` includes message generation.",
    ],
  });
  console.log(`\nWrote ${outPath}`);
}

main().then(() => process.exit(0)).catch((err) => {
  console.error("bench_pre failed:", err);
  process.exit(1);
});
