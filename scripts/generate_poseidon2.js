#!/usr/bin/env node
// Generates Poseidon2.sol with a complete Poseidon(2) implementation for BN254

import { buildPoseidon } from "circomlibjs";
import { writeFileSync } from "fs";

async function main() {
  const poseidon = await buildPoseidon();
  const F = poseidon.F;

  // Define the nRoundsF and nRoundsP for t=3
  const nRoundsF = 8;
  const nRoundsP = 57;
  const t = 3;

  // The buildPoseidon function internally uses the constants.
  // We need to get them from the internal source.
  // Let's use the poseidon reference implementation to extract constants.

  // The userData (from evmasm or the wasm) has the constants.
  // But we can also compute them from the poseidon function itself.
  
  // Actually, let's use the pre-computed constants from the circom source.
  // We'll read them from the poseidon_constants.js file.
  const poseidonConstants = (await import("circomlibjs/src/poseidon_constants.js")).default;
  const { utils } = await import("ffjavascript");
  const { unstringifyBigInts } = utils;
  const { C, M } = unstringifyBigInts(poseidonConstants);

  const cT3 = C[t - 2]; // C array for t=3 (81 elements)
  const mT3 = M[t - 2]; // MDS matrix for t=3 (3x3)

  // For P matrix (the partial first layer matrix), it's the same as the M matrix 
  // but with the last row replaced. Actually, in the circom code, P is separate.
  // Let me check…
  // From the circom poseidon_constants.circom, POSEIDON_P(t) returns a different matrix.
  // But circomlibjs doesn't seem to expose it. Let me check the evmasm.

  // Actually, the buildPoseidon function's evmasm code uses only M and C.
  // The P matrix and S matrix are not used in the WASM/optimized implementations.
  // Instead, the optimized impl uses a different structure.
  
  // Let's check what the optimized implementation looks like:
  const poseidonConstantsOpt = (await import("circomlibjs/src/poseidon_constants_opt.js")).default;
  const { C: Copt, S: Sopt } = unstringifyBigInts(poseidonConstantsOpt);

  const cOptT3 = Copt[t - 2];
  const sOptT3 = Sopt[t - 2];

  // Total constants in S for t=3: nRoundsP * (t*2 - 1) = 57 * 5 = 285
  console.log(`C[t-2] length: ${cT3.length} (expected ${t * nRoundsF + nRoundsP})`);
  console.log(`Copt[t-2] length: ${cOptT3.length} (expected ${t * nRoundsF + nRoundsP})`);
  console.log(`Sopt[t-2] length: ${sOptT3.length} (expected ${nRoundsP * (t*2 - 1)})`);
  console.log(`M[t-2] rows: ${mT3.length}, cols: ${mT3[0].length}`);

  function toHex256(n) {
    let s = n.toString(16);
    while (s.length < 64) s = "0" + s;
    return "0x" + s;
  }

  function formatArray(arr, indent) {
    const values = arr.map(v => toHex256(v));
    const lines = [];
    let line = indent;
    for (let i = 0; i < values.length; i++) {
      const val = values[i];
      const isLast = (i === values.length - 1);
      if (line.length + val.length + 2 > 120 && line !== indent) {
        lines.push(line);
        line = indent;
      }
      line += val;
      if (!isLast) line += ", ";
    }
    if (line !== indent) lines.push(line);
    return lines.join("\n");
  }

  function formatMatrix(mat, indent) {
    const rows = [];
    for (let i = 0; i < mat.length; i++) {
      const rowVals = mat[i].map(v => toHex256(v));
      rows.push(indent + "[" + rowVals.join(", ") + "]");
    }
    return rows.join(",\n");
  }

  const sol = `// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.20;

// Poseidon(2) hash function for BN254 (t=3)
// Generated from circomlibjs constants
library Poseidon2 {
    uint256 internal constant P = 21888242871839275222246405745257275088548364400416034343698204186575808495617;

    uint256 internal constant N_ROUNDS_F = 8;
    uint256 internal constant N_ROUNDS_P = 57;

    function pow5(uint256 x) internal pure returns (uint256) {
        uint256 x2 = mulmod(x, x, P);
        uint256 x4 = mulmod(x2, x2, P);
        return mulmod(x4, x, P);
    }

    function hash(uint256[2] memory inputs) internal pure returns (uint256) {
        uint256[3] memory st;
        st[0] = inputs[0];
        st[1] = inputs[1];
        st[2] = 0;

        // Round constants C[t][i] for t=3
        uint256[81] memory C = [
${formatArray(cOptT3, "            ")}
        ];

        // MDS matrix M[3][3]
        uint256[3][3] memory M;
        M[0] = [${mT3[0].map(toHex256).join(", ")}];
        M[1] = [${mT3[1].map(toHex256).join(", ")}];
        M[2] = [${mT3[2].map(toHex256).join(", ")}];

        // Sparse matrix S[t*2-1][nRoundsP] flattened
        // S[0..4] for round 0, S[5..9] for round 1, etc.
        uint256[285] memory S = [
${formatArray(sOptT3, "            ")}
        ];

        // Helper: apply full round (Sigma on all elements)
        function fullRound(uint256[3] memory state, uint256 roundIdx) pure returns (uint256[3] memory) {
            for (uint256 i = 0; i < 3; i++) {
                state[i] = pow5(state[i]);
                state[i] = addmod(state[i], C[roundIdx * 3 + i], P);
            }
            return matMul(state, M);
        }

        // Helper: S-box on element 0 only (partial round)
        function partialRound(uint256[3] memory state, uint256 roundIdx) pure returns (uint256[3] memory) {
            state[0] = pow5(state[0]);
            state[0] = addmod(state[0], C[roundIdx * 3], P);
            return sparseMul(state, roundIdx);
        }

        // Helper: matrix multiplication (3x3)
        function matMul(uint256[3] memory state, uint256[3][3] memory matrix) pure returns (uint256[3] memory) {
            uint256[3] memory result;
            for (uint256 i = 0; i < 3; i++) {
                result[i] = addmod(
                    addmod(
                        mulmod(matrix[i][0], state[0], P),
                        mulmod(matrix[i][1], state[1], P),
                        P
                    ),
                    mulmod(matrix[i][2], state[2], P),
                    P
                );
            }
            return result;
        }

        // Helper: sparse matrix multiplication for MixS
        function sparseMul(uint256[3] memory state, uint256 roundIdx) pure returns (uint256[3] memory) {
            uint256[3] memory result;
            uint256 base = roundIdx * 5;
            // First row: dense
            result[0] = addmod(
                addmod(
                    mulmod(S[base], state[0], P),
                    mulmod(S[base + 1], state[1], P),
                    P
                ),
                mulmod(S[base + 2], state[2], P),
                P
            );
            // Remaining rows: identity + column scaling
            result[1] = addmod(state[1], mulmod(S[base + 3], state[0], P), P);
            result[2] = addmod(state[2], mulmod(S[base + 4], state[0], P), P);
            return result;
        }

        // ── Poseidon permutation ──

        // Round 0: initial Ark
        for (uint256 i = 0; i < 3; i++) {
            st[i] = addmod(st[i], C[i], P);
        }

        // First half: 3 full rounds (nRoundsF/2 - 1 = 3)
        for (uint256 r = 0; r < 3; r++) {
            st = fullRound(st, r + 1);
        }

        // Middle full round (using P matrix = same as M for the final row mix)
        // The circom code uses P matrix here, which differs from M only for t>3.
        // For t=3, P has the same first two rows as M but different third row.
        // Source: poseidon_constants.circom POSEIDON_P(3):
        // [0x109b7f..., 0x1e6f20..., 0x1bd8c5...]
        // [0x16ed41..., 0x2d51ba..., 0x1b07d6...]
        // [0x2b90bb..., 0x11e12a..., 0x221c17...]
        // vs M(3):
        // [0x109b7f..., 0x2969f2..., 0x143021...]
        // [0x16ed41..., 0x2e2419..., 0x176cc0...]
        // [0x2b90bb..., 0x101071..., 0x19a3fc...]
        uint256[3][3] memory Pmat;
        Pmat[0] = [${[
            "0x109b7f411ba0e4c9b2b70caf5c36a7b194be7c11ad24378bfedb68592ba8118b",
            "0x1e6f20a11d1e31e43f83dcedddb9a0236203f5f24ae72c925a8a79a66831f51d",
            "0x1bd8c528472e57bdc722a141f8785694484f426725403ae24084e3027e782467"
        ].join(", ")}];
        Pmat[1] = [${[
            "0x16ed41e13bb9c0c66ae119424fddbcbc9314dc9fdbdeea55d6c64543dc4903e0",
            "0x2d51ba82c8073c6d6bacf1ad5e56655b7143625b0a9e9c3190527a1a5f05079a",
            "0x1b07d6d51e6f7e97e0ab10fc2e51ea83ce0611f940ff0731b5f927fe8d6a77c9"
        ].join(", ")}];
        Pmat[2] = [${[
            "0x2b90bba00fca0589f617e7dcbfe82e0df706ab640ceb247b791a93b74e36736d",
            "0x11e12a40d262ae88e8376f62d19edf43093cdef1ccf34d985a3e53f0bc5765a0",
            "0x221c170e4d02a2479c6f3e47b5ff55781574f980d89038308a3ef37cce8463bd"
        ].join(", ")}];

        // Middle: Sigma on all, then Ark, then Mix with P
        for (uint256 i = 0; i < 3; i++) {
            st[i] = pow5(st[i]);
            st[i] = addmod(st[i], C[4 * 3 + i], P);
        }
        st = matMul(st, Pmat);

        // 57 partial rounds
        for (uint256 r = 0; r < N_ROUNDS_P; r++) {
            st = partialRound(st, 5 * 3 + r);
        }

        // Last half: 3 full rounds (nRoundsF/2 - 1 = 3)
        for (uint256 r = 0; r < 3; r++) {
            st = fullRound(st, 5 * 3 + N_ROUNDS_P + r * 3);
        }

        // Final: Sigma on all, then output = M[0] * state (MixLast)
        for (uint256 i = 0; i < 3; i++) {
            st[i] = pow5(st[i]);
        }

        uint256 result;
        result = addmod(
            addmod(
                mulmod(M[0][0], st[0], P),
                mulmod(M[0][1], st[1], P),
                P
            ),
            mulmod(M[0][2], st[2], P),
            P
        );

        return result;
    }
}
`;

  writeFileSync("/home/daniela/Documents/Thesis/well-blockchain/src/lib/Poseidon2.sol", sol);
  console.log("Poseidon2.sol generated");

  // Verify against circomlibjs
  const testInputs = [BigInt(1), BigInt(2)];
  const expected = poseidon(testInputs);
  console.log("Test vector Poseidon(1,2) =", F.toString(expected));
}

main().catch(console.error);
