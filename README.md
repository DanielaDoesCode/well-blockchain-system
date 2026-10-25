# well-blockchain
well-blockchain access control system - thesis development


##### Running the project

`make all` # Full setup (powers of tau + circuit + zkey + verifier)
`make verifier` # Export solidity verifier only



#### ZK-authentication

### What is a Zero Knowledge Proof?
You want to prove a secret password, without actually telling anyone what the password is. A zero knowledge proofs lets you prove to a verifier that you know the secret without actually revealing the secret itself.

### What is Groth16?
Groth16 is a specific zk-SNARK protocol (Zero-Knowledge Succint Non-Interactive Arguments of Knowledge). It's one of the most efficiente and widely-used systems for generating and verify proofs on blockchain

Why Groth16?
- Very fast verification (important for gast costs)
- Small proof size (~200 bytes)
- Used by many zkproof projects

Tradeoffs: requires a trusted setup ceremony (generating global setup keys)

### Trusted Setup
Ceremony to generate public parameters.
Powers of Tau generates random numbers.
Power 15 means the circuit can have up to 2¹5 = 32,768 constraints.
In real-life system we'd do multi-party ceremonies with many contributors, this is due to the decentralized nature of the problem. This gets simplified for academic purposes.


## Foundry Testing

### Run all existing tests:
    forge test
### compile contracts:
    forge build
### compile and show gas report:
    forge test --gas-report


#### To test with real proofs:
    - Compile the circuit and generate proving/verification keys

- cd zk-circuits
- npm installs
./node_modules/.bin/circom zkauth.circom --r1cs --wasm --sym -o build
./node_modules/.bin/snarkjs groth16 setup zkauth.r1cs pot20_final.ptau build/verification_key.json
./node_modules/.bin/snarkjs groth16 contribute ... # generate zkey
./node_modules/.bin/snarkjs groth16 export verificationkey build/verification_key.json

### how to use anvil:
    - Terminal 1: Start a local node
    # Default Anvil node (port 8545):
    anvil
    # or with custom chain config:
    anvil --host 0.0.0.0 --port 8545

    - Terminal 2 - Deploy contracts
    Using forge: 
        forge create src/Registry.sol:WellRegistry \
        --rpc-url http://localhost:8545 \ --private-key [private key]
    - Using Cast (for interactive calls):
    # Depoly
    cast send <contract address> "register(bytes32,bytes32,bytes32,uint8)" \
    --rpc-url http://localhost8545
    --private-key 0x