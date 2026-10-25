#!/bin/bash
set -e
cd "$(dirname "$0")/.."

echo "=== Running all e2e suites against a fresh stack each ==="

./scripts/run_e2e.sh
echo ""

./scripts/run_e2e_policy.sh
echo ""

./scripts/run_e2e_research.sh
echo ""

./scripts/run_e2e_emergency.sh

echo ""
echo "=== All e2e suites passed ==="
