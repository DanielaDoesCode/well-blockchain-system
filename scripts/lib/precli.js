import { execFileSync } from "node:child_process";
import { mkdtempSync } from "node:fs";
import { tmpdir } from "node:os";
import { join, dirname } from "node:path";
import { fileURLToPath } from "node:url";

const here = dirname(fileURLToPath(import.meta.url));
const repoRoot = join(here, "..", "..");
const backendDir = join(repoRoot, "well-backend");

// Runs the reference PRE client. Set PRECLI_BIN to a prebuilt binary; otherwise
// falls back to `go run`, which is slower but needs no build step.
function invoke(args) {
  const bin = process.env.PRECLI_BIN;
  const out = bin
    ? execFileSync(bin, args, { cwd: repoRoot, encoding: "utf8" })
    : execFileSync("go", ["run", "./cmd/precli", ...args], { cwd: backendDir, encoding: "utf8" });
  return JSON.parse(out);
}

// precliSetup generates a patient/doctor PRE key pair, seals a random lockbox
// to the patient, and returns the lockbox, salt, public key and re-encryption
// key. statePath holds the secret keys; it is a client-side temp file.
// Pass a statePath to reuse one file across many invocations (the file only
// matters for the `self` / `verify` commands); otherwise a fresh temp file is
// created each call.
export function precliSetup(statePath) {
  const state = statePath || join(mkdtempSync(join(tmpdir(), "well-pre-")), "state.json");
  const out = invoke(["setup", "--state", state]);
  return { ...out, statePath: state };
}

// precliBench runs the in-process PRE primitive benchmark (precli bench) and
// returns its JSON: per-operation latency distributions, wire sizes and
// round-trip correctness counts.
export function precliBench(reps = 1000, verify = 500) {
  return invoke(["bench", "--reps", String(reps), "--verify", String(verify)]);
}

export function precliSelf(statePath, lockboxHex) {
  return invoke(["self", "--state", statePath, "--lockbox", lockboxHex]).message;
}

export function precliVerify(statePath, ctHex) {
  return invoke(["verify", "--state", statePath, "--ct", ctHex]).message;
}
