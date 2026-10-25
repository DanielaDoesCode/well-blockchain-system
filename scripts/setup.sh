#!/usr/bin/env bash
set -euo pipefail

# ──────────────────────────────────────────────────
# Full ZK setup: new → beacon → prepare phase2 → zkey → verifier → compile
# ──────────────────────────────────────────────────

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
PROJECT_DIR=$(dirname "$SCRIPT_DIR")
ZK_DIR="$PROJECT_DIR/zk-circuits"
BUILD_DIR="$ZK_DIR/build"
SRC_DIR="$PROJECT_DIR/src"

echo "=== 1/7: Removing old ptau files ==="
rm -f "$BUILD_DIR/pot15_beacon.ptau" "$BUILD_DIR/pot15_final.ptau"

echo "=== 2/7: Generating initial ptau ==="
cd "$ZK_DIR"
NODE_OPTIONS=--max-old-space-size=8192 npx snarkjs powersoftau new bn128 15 "$BUILD_DIR/pot15_0000.ptau"

echo "=== 3/8: Applying beacon to ptau ==="
NODE_OPTIONS=--max-old-space-size=8192 npx snarkjs powersoftau beacon \
  "$BUILD_DIR/pot15_0000.ptau" "$BUILD_DIR/pot15_beacon.ptau" \
  "77656c6c2d626c6f636b636861696e" 10 -n="Well Blockchain Beacon"

echo "=== 4/8: Preparing phase2 (takes ~2 min for power 15) ==="
NODE_OPTIONS=--max-old-space-size=8192 npx snarkjs powersoftau prepare phase2 \
  "$BUILD_DIR/pot15_beacon.ptau" "$BUILD_DIR/pot15_final.ptau"

echo "=== 5/8: Verifying ptau ==="
npx snarkjs powersoftau verify "$BUILD_DIR/pot15_final.ptau"

echo "=== 6/8: Generating zkey and contributing ==="
rm -f "$BUILD_DIR/zkauth_final.zkey" "$BUILD_DIR/verification_key.json" "$BUILD_DIR/Verifier.sol"
npx snarkjs groth16 setup "$BUILD_DIR/zkauth.r1cs" "$BUILD_DIR/pot15_final.ptau" "$BUILD_DIR/zkauth_final.zkey"
npx snarkjs zkey contribute "$BUILD_DIR/zkauth_final.zkey" "$BUILD_DIR/zkauth_final.zkey" \
  -e="well-blockchain-$(date +%s)"

echo "=== 7/8: Exporting verification key + Solidity verifier ==="
npx snarkjs zkey export verificationkey "$BUILD_DIR/zkauth_final.zkey" "$BUILD_DIR/verification_key.json"
npx snarkjs zkey export solidityverifier "$BUILD_DIR/zkauth_final.zkey" "$BUILD_DIR/Verifier.sol"
cp "$BUILD_DIR/Verifier.sol" "$SRC_DIR/Groth16Verifier.sol"

echo "=== 8/8: Compiling contracts ==="
cd "$PROJECT_DIR"
forge build

echo ""
echo "=== IC points (should be large random numbers, NOT [0,1]): ==="
jq '.IC' "$BUILD_DIR/verification_key.json"

echo ""
echo "=== Size check ==="
for c in Groth16Verifier ZKAuthVerifier WellRegistry; do
  size=$(forge inspect "$c" deployedBytecode | wc -c)
  echo "  $c: $size bytes ($([ "$size" -lt 24576 ] && echo 'OK' || echo '>24KB — needs refactor for testnet'))"
done

echo ""
echo "=== Setup complete ==="
echo "Deploy with: make deploy"
