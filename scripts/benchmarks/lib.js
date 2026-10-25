import "dotenv/config";
import path from "path";
import os from "os";
import crypto from "crypto";
import { execFileSync } from "child_process";
import { fileURLToPath } from "url";
import { createRequire } from "module";

const require = createRequire(import.meta.url);
const __dirname = path.dirname(fileURLToPath(import.meta.url));

export const BENCH_ROOT = path.join(__dirname, "..", "..");
export const ZK_ROOT = path.join(BENCH_ROOT, "zk-circuits");
// BENCH_RESULTS_DIR lets the latency-injected suite write to
// results/wan/<situation>/<mode>/ without touching the loopback results.
export const RESULTS_DIR = process.env.BENCH_RESULTS_DIR || path.join(__dirname, "results");

export const SNAPSHOT = {
  date: new Date().toISOString(),
  node: `${process.version}`,
  platform: `${process.platform}-${process.arch}`,
  provider: undefined, // filled by bench scripts
  geth: undefined,     // filled by bench scripts
  environment: {
    singleNode: process.env.BENCH_NODE_DESC || "single local geth Clique PoA node (2s period)",
    networkLatency: process.env.BENCH_NET_LATENCY || "zero (loopback)",
    vault: "dev-grade Vault (transit engine), emergency-recovery escrow only",
    pre: "etclab/pre AFGH second attempt (BLS12-381), record delegation",
    ipfs: "single IPFS peer",
    zkAuthVerifier: zkAuthVerifierDesc(),
  },
};

export function zkAuthVerifierDesc() {
  const mode = process.env.ZK_AUTH_VERIFIER_MODE;
  if (mode === "real") {
    return "Groth16Verifier real on-chain verification (proofs genuinely verified)";
  }
  if (mode === "mock") {
    return "MockVerifier (dev wiring, accepts all proofs; real Groth16Verifier measured separately)";
  }
  // No env marker: infer from what deploy_mock was told (default mock for dev).
  return process.env.DEPLOY_AUTH_VERIFIER === "real"
    ? "Groth16Verifier real on-chain verification (proofs genuinely verified)"
    : "MockVerifier (dev wiring, accepts all proofs; real Groth16Verifier measured separately)";
}

const snarkjs = require(path.join(ZK_ROOT, "node_modules", "snarkjs"));
const circomlibjs = require(path.join(ZK_ROOT, "node_modules", "circomlibjs"));

const WTNS_DIR = (() => {
  const dir = process.env.BENCH_TMP_DIR || path.join(os.tmpdir(), "well-bench");
  return dir;
})();

export const BUILDS = {
  zkauth: {
    name: "zkauth",
    wasm: path.join(ZK_ROOT, "build", "zkauth_js", "zkauth.wasm"),
    zkey: path.join(ZK_ROOT, "build", "zkauth_final.zkey"),
  },
  policy: {
    name: "policy_match",
    wasm: path.join(ZK_ROOT, "build", "policy_match_js", "policy_match.wasm"),
    zkey: path.join(ZK_ROOT, "build", "policy_match_final.zkey"),
  },
  research: {
    name: "research_contribute",
    wasm: path.join(ZK_ROOT, "build", "research_contribute_js", "research_contribute.wasm"),
    zkey: path.join(ZK_ROOT, "build", "research_contribute_final.zkey"),
  },
};

// ---------------------------------------------------------------- statistics
export const pct = (sorted, p) => {
  if (sorted.length === 0) return null;
  const idx = Math.min(sorted.length - 1, Math.max(0, Math.ceil((p / 100) * sorted.length) - 1));
  return sorted[idx];
};

export const stat = (samples) => {
  const vals = samples.map(Number).sort((a, b) => a - b);
  if (vals.length === 0) return null;
  const sum = vals.reduce((a, b) => a + b, 0);
  return {
    n: vals.length,
    mean: sum / vals.length,
    min: vals[0],
    p50: pct(vals, 50),
    p95: pct(vals, 95),
    p99: pct(vals, 99),
    max: vals[vals.length - 1],
  };
};

export const fmtStat = (s, unit = "ms", decimals = 1) => {
  if (!s) return "n/a";
  const d = (x) => (x == null ? "n/a" : x.toFixed(decimals) + unit);
  return `${d(s.mean)}  (p50 ${d(s.p50)} | p95 ${d(s.p95)} | min ${d(s.min)} | max ${d(s.max)}, n=${s.n})`;
};

export const fmtMs = (x) => (x == null ? "n/a" : `${x.toFixed(1)} ms`);
export const fmtGas = (x) => (x == null ? "n/a" : `${Number(x).toLocaleString("en-US")} gas`);

// ---------------------------------------------------------------- metrics sink
export const metrics = {};

export function addMetric(key, ms, extra = {}) {
  (metrics[key] = metrics[key] || []).push({ ms, ...extra });
  return metrics[key];
}

export function metricStats(key) {
  const entries = metrics[key] || [];
  return stat(entries.map((e) => e.ms));
}

export function gasStat(key) {
  const entries = metrics[key] || [];
  return stat(entries.filter((e) => e.gas != null).map((e) => e.gas));
}

export function collect(reports) {
  for (const [k, extra] of Object.entries(reports)) {
    if (extra != null) addMetric(k, extra.ms, extra);
  }
}

// ---------------------------------------------------------------- result output
function ensureDir(dir) {
  try {
    execFileSync("mkdir", ["-p", dir]);
  } catch { /* ignore */ }
}

export function writeJSON(filename, obj) {
  ensureDir(RESULTS_DIR);
  const outPath = path.join(RESULTS_DIR, filename);
  const fs = require("fs");
  fs.writeFileSync(outPath, JSON.stringify(obj, null, 2) + "\n");
  return outPath;
}

export function printTable(rows) {
  const widths = [];
  for (const row of rows) {
    row.forEach((c, i) => { widths[i] = Math.max(widths[i] || 0, String(c).length); });
  }
  const line = (row) => row.map((c, i) => String(c).padEnd(widths[i])).join("  ").trimEnd();
  console.log(line(rows[0]));
  console.log(rows.slice(1).map(line).join("\n"));
}

// ---------------------------------------------------------------- env / on-chain helpers
const required = (k) => {
  const v = process.env[k];
  if (!v) throw new Error(`missing ${k} in .env — run scripts/deploy_mock.js first`);
  return v;
};

export function readAddrs() {
  return {
    registry: required("REGISTRY_ADDRESS"),
    userRegistry: required("USER_REGISTRY_ADDRESS"),
    historyFactory: required("HISTORY_FACTORY_ADDRESS"),
    zkAuth: required("ZKAUTH_VERIFIER_ADDRESS"),
    incident: required("INCIDENT_MANAGER_ADDRESS"),
    riskAssessment: required("RISK_ASSESSMENT_ADDRESS"),
    emergencyFactory: required("EMERGENCY_SUMMARY_FACTORY_ADDRESS"),
    attestation: required("ATTENTATION_MANAGER_ADDRESS"),
    researchVerifier: required("RESEARCH_VERIFIER_ADDRESS"),
    researchContributer: required("RESEARCH_CONTRIBUTER_ADDRESS"),
    poseidonT3: required("POSEIDON_T3_ADDRESS"),
  };
}

// API_BASE_URL lets the latency-injected suite route REST calls through the API
// relay (18443) while the backend keeps listening on AGGREGATOR_PORT.
export const API_BASE = process.env.API_BASE_URL || `http://localhost:${process.env.AGGREGATOR_PORT || 8081}`;

export async function apiGet(path, { token, user } = {}) {
  const headers = {};
  if (user) headers["X-User-Address"] = user;
  if (token) headers["X-Session-Token"] = token;
  const res = await fetch(`${API_BASE}${path}`, { headers });
  let data = null;
  try { data = await res.json(); } catch { data = null; }
  return { status: res.status, data };
}

export async function apiPost(path, body, { token, user } = {}) {
  const headers = { "Content-Type": "application/json" };
  if (user) headers["X-User-Address"] = user;
  if (token) headers["X-Session-Token"] = token;
  const res = await fetch(`${API_BASE}${path}`, {
    method: "POST", headers, body: JSON.stringify(body),
  });
  let data = null;
  try { data = await res.json(); } catch { data = null; }
  return { status: res.status, data };
}

export async function apiUpload(fileData, { token, user } = {}) {
  const form = new FormData();
  form.append("file", new Blob([fileData]));
  const headers = {};
  if (user) headers["X-User-Address"] = user;
  if (token) headers["X-Session-Token"] = token;
  const res = await fetch(`${API_BASE}/api/v1/records/upload`, {
    method: "POST", headers, body: form,
  });
  let data = null;
  try { data = await res.json(); } catch { data = null; }
  return { status: res.status, data };
}

// Measure a single HTTP op (includes full request/response round trip).
export async function measureApi(label, fn) {
  const t0 = Date.now();
  const value = await fn();
  const ms = Date.now() - t0;
  addMetric(label, ms);
  return { ms, value, status: value?.status };
}

// Measure a transaction: send (RPC round trip) vs mine (block inclusion) + gas.
export async function measureTx(label, fn) {
  const t0 = Date.now();
  const tx = await fn();
  const sendMs = Date.now() - t0;
  const t1 = Date.now();
  const receipt = await tx.wait();
  const mineMs = Date.now() - t1;
  const gas = receipt && receipt.gasUsed != null ? receipt.gasUsed.toString() : null;
  addMetric(`${label}.send`, sendMs);
  addMetric(`${label}.mine`, mineMs);
  addMetric(`${label}.total`, sendMs + mineMs, { gas });
  return { sendMs, mineMs, totalMs: sendMs + mineMs, gas, receipt };
}

// ---------------------------------------------------------------- DB polling (listener propagation)
const DB_POLL_CMD = () =>
  (process.env.DB_POLL_CMD || "docker exec postgres psql -U postgres -d well_db -t -A").split(/\s+/);

export function dbQuery(sql) {
  const [cmd, ...args] = DB_POLL_CMD();
  const out = execFileSync(cmd, [...args, "-c", sql], { encoding: "utf8", stdio: ["ignore", "pipe", "pipe"] });
  return out.trim();
}

export function dbRow(sql) {
  const rows = dbQuery(sql).split("\n").map((r) => r.trim()).filter(Boolean);
  return rows.length ? rows[0] : null;
}

// Poll the DB until a condition row appears; returns latency from start.
export async function pollDb(sql, timeoutMs = 15000, pollEveryMs = 100) {
  const t0 = Date.now();
  for (;;) {
    if (dbRow(sql)) return Date.now() - t0;
    if (Date.now() - t0 > timeoutMs) return null;
    await new Promise((r) => setTimeout(r, pollEveryMs));
  }
}

// Listener latency: on-chain event -> event_logs row.
export async function waitForEvent(txHash, eventType, timeoutMs = 15000) {
  return pollDb(
    `SELECT event_type FROM event_logs WHERE tx_hash='${txHash}' AND event_type='${eventType}' LIMIT 1`,
    timeoutMs
  );
}

// Listener latency: on-chain access event -> access_logs row.
export async function waitForAccessLog(txHash, accessType, timeoutMs = 15000) {
  return pollDb(
    `SELECT access_type FROM access_logs WHERE tx_hash='${txHash}' AND access_type='${accessType}' LIMIT 1`,
    timeoutMs
  );
}

// ---------------------------------------------------------------- identity helpers
export async function provisionIdentity(adminSession, admin, identity) {
  const res = await apiPost("/api/v1/identities", identity, { token: adminSession, user: admin });
  if (res.status !== 201) throw new Error(`provision identity failed: ${JSON.stringify(res)}`);
  return res.data;
}

export async function resolveIdentity(token, user, query) {
  const [k, v] = Object.entries(query)[0];
  return apiGet(`/api/v1/identities?${k}=${encodeURIComponent(v)}`, { token, user });
}

// ---------------------------------------------------------------- poseidon + circuit helpers
export async function buildPoseidon() {
  return circomlibjs.buildPoseidon();
}

const F = (poseidon, x) => poseidon.F.toObject(x);

export function toHex(n) {
  return "0x" + BigInt(n).toString(16).padStart(64, "0");
}

// zkauth.circom: commitment=Poseidon(secret); nullifier=Poseidon(secret, seed, actionId)
export function zkauthInputs(poseidon, { secret, nullifierSeed, actionId }) {
  const commitment = F(poseidon, poseidon([secret]));
  const nullifier = F(poseidon, poseidon([secret, nullifierSeed, actionId]));
  return {
    commitment,
    nullifier,
    inputs: {
      secret: secret.toString(),
      nullifierSeed: nullifierSeed.toString(),
      commitment: commitment.toString(),
      nullifier: nullifier.toString(),
      actionId: actionId.toString(),
    },
  };
}

// policy_match.circom inputs (requiredDept/requiredInst = 0 is a wildcard).
export function policyInputs({ department, institution, requiredDept = 0n, requiredInst = 0n, doctorAddr, validUntil, evaluatedAt }) {
  return {
    department: department.toString(),
    institution: institution.toString(),
    requiredDept: requiredDept.toString(),
    requiredInst: requiredInst.toString(),
    doctorAddr: doctorAddr.toString(),
    validUntil: validUntil.toString(),
    evaluatedAt: evaluatedAt.toString(),
  };
}

// research_contribute.circom: commitment=Poseidon(secret,1); leaf=Poseidon(commitment,conditionGroup).
// Merkle tree semantics match zk-circuits/gen_merkle_proof.js and AttestationManager.
export function researchZeros(poseidon, treeDepth = 20) {
  const zeros = new Array(treeDepth + 1);
  zeros[0] = 0n;
  for (let i = 1; i <= treeDepth; i++) {
    zeros[i] = F(poseidon, poseidon([zeros[i - 1], zeros[i - 1]]));
  }
  return zeros;
}

export function researchLeafInsert(poseidon, leaf, treeDepth = 20) {
  const zeros = researchZeros(poseidon, treeDepth);
  const filledSubtrees = zeros.slice(0, treeDepth);
  const proofElements = new Array(treeDepth);
  const proofIndices = new Array(treeDepth);
  let index = 0;
  let currentHash = leaf;
  for (let i = 0; i < treeDepth; i++) {
    const isLeft = index % 2 === 0;
    if (isLeft) {
      proofElements[i] = zeros[i];
      proofIndices[i] = 0;
      filledSubtrees[i] = currentHash;
    } else {
      proofElements[i] = filledSubtrees[i];
      proofIndices[i] = 1;
    }
    const left = isLeft ? currentHash : filledSubtrees[i];
    const right = isLeft ? zeros[i] : currentHash;
    currentHash = F(poseidon, poseidon([left, right]));
    index >>= 1;
  }
  return { merkleRoot: currentHash, merkleProofElements: proofElements, merkleProofIndices: proofIndices };
}

export function researchInputs(poseidon, { secret, nullifierSeed, conditionGroup }) {
  const commitment = F(poseidon, poseidon([secret, 1n]));
  const leaf = F(poseidon, poseidon([commitment, conditionGroup]));
  const { merkleRoot, merkleProofElements, merkleProofIndices } = researchLeafInsert(poseidon, leaf);
  return {
    commitment,
    merkleRoot,
    inputs: {
      secret: secret.toString(),
      nullifierSeed: nullifierSeed.toString(),
      conditionGroup: conditionGroup.toString(),
      merkleProofElements: merkleProofElements.map((e) => e.toString()),
      merkleProofIndices: merkleProofIndices.map((i) => i.toString()),
      merkleRoot: merkleRoot.toString(),
    },
  };
}

// ---------------------------------------------------------------- witness + prove (split timing)
function wtnsFile() {
  ensureDir(WTNS_DIR);
  return path.join(WTNS_DIR, `wtns-${crypto.randomBytes(6).toString("hex")}.wtns`);
}

/**
 * Compute the witness and the Groth16 proof separately so each stage is timed.
 * Returns the proof in the on-chain (Solidity) argument layout.
 */
export async function genProofSplit({ wasm, zkey, inputs }) {
  const wtns = wtnsFile();
  let witnessMs, proveMs;
  let t0 = Date.now();
  await snarkjs.wtns.calculate(inputs, wasm, wtns);
  witnessMs = Date.now() - t0;

  t0 = Date.now();
  const { proof, publicSignals } = await snarkjs.groth16.prove(zkey, wtns);
  proveMs = Date.now() - t0;

  return {
    witnessMs,
    proveMs,
    totalMs: witnessMs + proveMs,
    proof,
    publicSignals,
    formatted: formatProof(proof, publicSignals),
  };
}

export function formatProof(proof, publicSignals) {
  return {
    pA: proof.pi_a.slice(0, 2).map((x) => x.toString()),
    pB: [
      proof.pi_b[0].slice(0, 2).reverse().map((x) => x.toString()),
      proof.pi_b[1].slice(0, 2).reverse().map((x) => x.toString()),
    ],
    pC: proof.pi_c.slice(0, 2).map((x) => x.toString()),
    publicSignals: publicSignals.map(toHex),
  };
}

// Yields a proof for a circuit with its own pre-built inputs.
export async function proveFor(buildKey, inputs) {
  const build = BUILDS[buildKey];
  return genProofSplit({ wasm: build.wasm, zkey: build.zkey, inputs });
}

// ---------------------------------------------------------------- sessions
async function parseEvent(receipt, iface, name) {
  const log = receipt.logs
    .map((l) => { try { return iface.parseLog({ topics: [...l.topics], data: l.data }); } catch { return null; } })
    .find((e) => e?.name === name);
  return log;
}

// Mock-wired dev authenticate (fast path for session provisioning in throughput runs).
export async function devZkAuthSession(zkAuth, signer) {
  const tx = await zkAuth.connect(signer).authenticate(
    signer.address,
    ethersHex(32),
    ["0", "0"], [["0", "0"], ["0", "0"]], ["0", "0"],
    ethersHex(32),
    ethersHex(32)
  );
  const receipt = await tx.wait();
  const log = await parseEvent(receipt, zkAuth.interface, "Authenticated");
  if (!log) throw new Error(`no Authenticated event for ${signer.address}`);
  return { token: log.args.sessionToken, receipt, mineMs: 0 };
}

function ethersHex(bytes) {
  const b = crypto.randomBytes(bytes);
  return "0x" + b.toString("hex");
}

// Genuine ZK session: real zkauth proof + on-chain authenticate.
export async function realZkAuthSession(zkAuth, signer, poseidon, { secret, nullifierSeed, actionId }) {
  const { commitment, nullifier, inputs } = zkauthInputs(poseidon, { secret, nullifierSeed, actionId });
  const proof = await proveFor("zkauth", inputs);
  const receipt = await (async () => {
    const tx = await zkAuth.connect(signer).authenticate(
      signer.address,
      toHex(commitment),
      proof.formatted.pA, proof.formatted.pB, proof.formatted.pC,
      toHex(nullifier),
      toHex(actionId)
    );
    return tx.wait();
  })();
  const log = await parseEvent(receipt, zkAuth.interface, "Authenticated");
  if (!log) throw new Error(`genuine authenticate failed for ${signer.address}`);
  return { token: log.args.sessionToken, commitment, nullifier, receipt, proof };
}

export function zeroAddress() {
  return "0x0000000000000000000000000000000000000000";
}