.PHONY: up down db node compile deploy backend stop-db stop-node stop-backend clean ptau zkey research-zkey policy-zkey zk-setup test-ledger test-ledger-unit test-ledger-contract

# ── Full stack ──────────────────────────────────────
up: compile db
	@echo "Starting Hardhat node..."
	@npx hardhat node --hostname 0.0.0.0 --port 8545 &>/tmp/hardhat.log & sleep 3
	@echo "Hardhat node ready on :8545"
	npx hardhat run scripts/deploy.js --network localhost
	@set -a && . ./.env && set +a && cd well-backend && go run ./cmd/ & sleep 2
	@echo "All services started. Showing chain logs (Ctrl+C to hide, make down to stop):"
	@tail -f /tmp/hardhat.log

# ── Compile ────────────────────────────────────────
compile:
	npx hardhat compile

# ── Database ───────────────────────────────────────
db:
	@echo "Starting PostgreSQL..."
	docker compose up -d postgres
	@sleep 2
	@echo "PostgreSQL ready on :5432"

stop-db:
	docker compose stop postgres

# ── Hardhat node ───────────────────────────────────
node:
	@echo "Starting Hardhat node..."
	@npx hardhat node --hostname 0.0.0.0 --port 8545 &>/tmp/hardhat.log &
	@sleep 3
	@echo "Hardhat node ready on :8545"

stop-node:
	@-pkill -f 'hardhat node' 2>/dev/null; echo "Hardhat node stopped."

# ── Deploy contracts ───────────────────────────────
deploy:
	npx hardhat run scripts/deploy.js --network localhost

# ── Backend ────────────────────────────────────────
backend:
	@echo "Starting backend..."
	@set -a && . ./.env && set +a && cd well-backend && go run ./cmd/ &

stop-backend:
	@-pkill -f 'go run ./cmd/' 2>/dev/null; echo "Backend stopped."

# ── Stop everything ────────────────────────────────
down: stop-backend stop-node stop-db
	@-pkill -f 'hardhat node' 2>/dev/null; echo "Hardhat (local) stopped."

# ── ZK circuit setup ───────────────────────────────
# Apply beacon to ptau, prepare phase2 (needed once; takes ~2 min)
ptau:
	rm -f zk-circuits/build/pot15_beacon.ptau zk-circuits/build/pot15_final.ptau
	cd zk-circuits && \
	NODE_OPTIONS=--max-old-space-size=8192 npx snarkjs powersoftau new bn128 15 build/pot15_0000.ptau && \
	NODE_OPTIONS=--max-old-space-size=8192 npx snarkjs powersoftau beacon \
		build/pot15_0000.ptau build/pot15_beacon.ptau \
		"77656c6c2d626c6f636b636861696e" 10 -n="Well Blockchain Beacon" && \
	NODE_OPTIONS=--max-old-space-size=8192 npx snarkjs powersoftau prepare phase2 \
		build/pot15_beacon.ptau build/pot15_final.ptau && \
	npx snarkjs powersoftau verify build/pot15_final.ptau

# Regenerate zkey for zkauth circuit, export verifier
zkey:
	cd zk-circuits && \
	rm -f build/zkauth_final.zkey build/verification_key.json build/Verifier.sol && \
	npx snarkjs groth16 setup build/zkauth.r1cs build/pot15_final.ptau build/zkauth_final.zkey && \
	npx snarkjs zkey contribute build/zkauth_final.zkey build/zkauth_final.zkey \
		-e="well-blockchain-$$(date +%s)" && \
	npx snarkjs zkey export verificationkey build/zkauth_final.zkey build/verification_key.json && \
	npx snarkjs zkey export solidityverifier build/zkauth_final.zkey build/Verifier.sol && \
	cp build/Verifier.sol ../src/Groth16Verifier.sol && \
	echo "IC points:" && jq '.IC' build/verification_key.json

# Regenerate zkey for research_contribute circuit
research-zkey:
	cd zk-circuits && \
	rm -f build/research_contribute_0000.zkey build/research_contribute_final.zkey build/verification_key_research.json && \
	npx snarkjs groth16 setup build/research_contribute.r1cs build/pot15_final.ptau build/research_contribute_0000.zkey && \
	npx snarkjs zkey contribute build/research_contribute_0000.zkey build/research_contribute_final.zkey \
		-e="well-blockchain-research-$$(date +%s)" && \
	npx snarkjs zkey export verificationkey build/research_contribute_final.zkey build/verification_key_research.json && \
	npx snarkjs zkey export solidityverifier build/research_contribute_final.zkey build/ResearchVerifier.sol && \
	sed -i 's/contract Groth16Verifier/contract ResearchVerifier/' build/ResearchVerifier.sol && \
	cp build/ResearchVerifier.sol ../src/ResearchVerifier.sol && \
	echo "nPublic (must be 3):" && jq '.nPublic' build/verification_key_research.json && \
	echo "IC points:" && jq '.IC' build/verification_key_research.json

# Regenerate zkey for policy_match circuit
policy-zkey:
	cd zk-circuits && \
	rm -rf build/policy_match.r1cs build/policy_match.sym build/policy_match_js build/policy_match_0000.zkey build/policy_match_final.zkey build/verification_key_policy.json && \
	circom policy_match.circom --r1cs --wasm --sym -o build -l node_modules && \
	npx snarkjs groth16 setup build/policy_match.r1cs build/pot15_final.ptau build/policy_match_0000.zkey && \
	npx snarkjs zkey contribute build/policy_match_0000.zkey build/policy_match_final.zkey \
		-e="well-blockchain-policy-$$(date +%s)" && \
	npx snarkjs zkey export verificationkey build/policy_match_final.zkey build/verification_key_policy.json && \
	npx snarkjs zkey export solidityverifier build/policy_match_final.zkey build/PolicyVerifier.sol && \
	sed -i 's/contract Groth16Verifier/contract PolicyVerifier/' build/PolicyVerifier.sol && \
	cp build/PolicyVerifier.sol ../src/PolicyVerifier.sol && \
	echo "nPublic (must be 1):" && jq '.nPublic' build/verification_key_policy.json && \
	echo "IC points:" && jq '.IC' build/verification_key_policy.json

# Full ZK pipeline: ptau → zkey → compile
zk-setup: ptau zkey research-zkey policy-zkey compile

# ── Deploy contracts with mock verifier (for e2e) ──
deploy-mock:
	npx hardhat run scripts/deploy_mock.js --network localhost

# ── E2E test ───────────────────────────────────────
test-e2e: compile
	@echo "Starting docker services (postgres, ipfs, vault)..."
	docker compose up -d postgres ipfs vault 2>/dev/null
	@sleep 5
	@echo "Applying DB schema..."
	@docker exec -i postgres psql -U postgres -d well_db < db/init.sql 2>/dev/null || true
	@echo "Starting local Hardhat node..."
	@npx hardhat node --hostname 0.0.0.0 --port 8545 &>/tmp/hardhat.log & echo "$$!" > /tmp/hardhat.pid
	@sleep 3
	$(MAKE) deploy-mock
	@echo "Starting backend..."
	@set -a && . ./.env && set +a && cd well-backend && go run ./cmd/ &>/tmp/backend.log & echo "$$!" > /tmp/backend.pid
	@sleep 3
	@echo "Running e2e test..."
	@npx hardhat run scripts/e2e_test.js --network localhost || (kill $$(cat /tmp/backend.pid) 2>/dev/null; kill $$(cat /tmp/hardhat.pid) 2>/dev/null; exit 1)
	@kill $$(cat /tmp/backend.pid) 2>/dev/null; kill $$(cat /tmp/hardhat.pid) 2>/dev/null; echo "Cleaned up."
	@echo "E2E test complete."

# ── Single service helpers for manual e2e ──────────
start-infra:
	docker compose up -d postgres ipfs vault
	@echo "Then start Hardhat locally:  npx hardhat node --hostname 0.0.0.0 --port 8545"

start-backend:
	@set -a && . ./.env && set +a && cd well-backend && go run ./cmd/

run-e2e:
	npx hardhat run scripts/e2e_test.js --network localhost

# ── Ledger (tamper-evident audit log) tests ─────────
# Pure-unit: Go Poseidon + canonical row encoding golden tests.
test-ledger-unit:
	go test ./well-backend/internal/ledger/hash/...

# Contract: LogAnchor.sol behaviour via foundry.
test-ledger-contract:
	forge test --match-path test/LogAnchor.t.sol

# Full-stack: hardhat node + postgres, sealed+anchored batches, tamper detection.
test-ledger: compile
	@echo "Bringing up postgres + applying schema..."
	docker compose up -d postgres
	@sleep 3
	@docker exec -i postgres psql -U postgres -d well_db < db/init.sql 2>/dev/null || true
	@echo "Freeing :8545 (stopping geth/hardhat containers)..."
	@-docker stop geth 2>/dev/null || true
	@-docker stop hardhat 2>/dev/null || true
	@echo "Starting Hardhat node on :8545..."
	@npx hardhat node --hostname 0.0.0.0 --port 8545 &>/tmp/hardhat-ledger.log & echo "$$!" > /tmp/hardhat-ledger.pid
	@sleep 3
	@cd well-backend && go test -tags integration ./internal/ledger/ -count=1 -v || (kill $$(cat /tmp/hardhat-ledger.pid) 2>/dev/null; exit 1)
	@kill $$(cat /tmp/hardhat-ledger.pid) 2>/dev/null; echo "Cleaned up hardhat node."
	@echo "Ledger integration test complete."

# ── Clean ──────────────────────────────────────────
clean:
	@rm -rf out/ cache/ artifacts/
	@echo "Cleaned build artifacts."
