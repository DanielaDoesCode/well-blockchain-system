import "dotenv/config";
import hre from "hardhat";
import {
  BUILDS,
  buildPoseidon, zkauthInputs, policyInputs, researchInputs,
  genProofSplit, stat, printTable, writeJSON, SNAPSHOT,
} from "./lib.js";

const { ethers } = hre;

import crypto from "crypto";

const REPS = Number(process.env.BENCH_ZK_REPS || 5);

const rand = () => BigInt("0x" + crypto.randomBytes(31).toString("hex")) + 1n;

async function runCircuit(poseidon, key) {
  const samples = [];
  for (let i = 0; i < REPS; i++) {
    let inputs;
    if (key === "zkauth") {
      inputs = zkauthInputs(poseidon, {
        secret: rand(), nullifierSeed: rand(), actionId: rand(),
      }).inputs;
    } else if (key === "policy") {
      const dept = rand();
      inputs = policyInputs({
        department: dept, institution: dept,
        doctorAddr: BigInt("0x" + crypto.randomBytes(20).toString("hex")),
        validUntil: 2000000000n, evaluatedAt: 1700000000n,
      });
    } else if (key === "research") {
      inputs = researchInputs(poseidon, {
        secret: rand(), nullifierSeed: rand(), conditionGroup: 42n,
      }).inputs;
    }
    const res = await genProofSplit({ wasm: BUILDS[key].wasm, zkey: BUILDS[key].zkey, inputs });
    samples.push(res);
  }
  return {
    key,
    total: stat(samples.map((s) => s.totalMs)),
    witness: stat(samples.map((s) => s.witnessMs)),
    prove: stat(samples.map((s) => s.proveMs)),
    samples,
  };
}

async function main() {
  const provider = ethers.provider;
  SNAPSHOT.provider = `${provider._network?.name || "jsonrpc"}@${provider._connection?.url || "localhost:8545"}`;
  SNAPSHOT.geth = await provider.send("web3_clientVersion").catch(() => "n/a");

  console.log("=== ZK proof generation benchmark (Groth16, real circuits) ===");
  console.log(`reps per circuit: ${REPS}\n`);

  const poseidon = await buildPoseidon();
  const meanMs = (s) => `${(s?.mean ?? 0).toFixed(2)} ms`;
  const rows = [];
  const results = {};
  for (const key of ["zkauth", "policy", "research"]) {
    console.log(`Proving ${BUILDS[key].name}...`);
    const r = await runCircuit(poseidon, key);
    results[key] = r;
    rows.push([
      BUILDS[key].name,
      meanMs(r.total),
      meanMs(r.witness),
      meanMs(r.prove),
      `${(r.prove.p95 ?? 0).toFixed(2)} ms`,
    ]);
  }

  console.log("\nPer-circuit (per sample is witness+Poseidon + Groth16 prove, all in-proc):");
  printTable([
    ["circuit", "total mean", "witness mean", "prove mean", "prove p95"],
    ...rows,
  ]);

  const summary = {};
  for (const [k, r] of Object.entries(results)) {
    summary[k] = {
      circuit: BUILDS[k].name,
      totalMeanMs: r.total.mean, totalP50Ms: r.total.p50, totalP95Ms: r.total.p95,
      witnessMeanMs: r.witness.mean, witnessP95Ms: r.witness.p95,
      proveMeanMs: r.prove.mean, proveP95Ms: r.prove.p95,
      reps: REPS,
    };
  }
  const outPath = writeJSON("zk.json", {
    snapshot: SNAPSHOT,
    unit: "milliseconds",
    repsPerCircuit: REPS,
    circuits: summary,
  });
  console.log(`\nWrote ${outPath}`);
}

main().then(() => process.exit(0)).catch((err) => {
  console.error("bench_zk failed:", err);
  process.exit(1);
});