// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.20;

import "./Registry.sol";
import "./AttestationManager.sol";

interface IResearchVerifier {
    function verifyProof(
        uint256[2] calldata _pA,
        uint256[2][2] calldata _pB,
        uint256[2] calldata _pC,
        uint256[3] calldata _pubSignals
    ) external view returns (bool);
}

contract ResearchContributer {
    IResearchVerifier    public researchVerifier;
    WellRegistry         public registry;
    AttestationManager   public attestationManager;

    mapping(bytes32 => bool) public usedNullifiers;

    error CommitmentNotRegistered();
    error NullifierAlreadyUsed();
    error InvalidProof();
    error StaleOrInvalidRoot();

    event ContributionSubmitted(
        bytes32 indexed commitment,
        bytes32 indexed nullifier
    );

    constructor(
        address _researchVerifier,
        address _registry,
        address _attestationManager
    ) {
        researchVerifier   = IResearchVerifier(_researchVerifier);
        registry           = WellRegistry(_registry);
        attestationManager = AttestationManager(_attestationManager);
    }

    /// @notice Submit a ZK proof of research contribution.
    /// @param pA          Groth16 proof element A
    /// @param pB          Groth16 proof element B
    /// @param pC          Groth16 proof element C
    /// @param commitment  Public output: patient research commitment (Poseidon(secret, 1))
    /// @param nullifier   Public output: per-group nullifier (Poseidon(secret, nullifierSeed, conditionGroup))
    /// @param merkleRoot  Public input: the Merkle root the proof was built against.
    ///                    Must match AttestationManager.isKnownRoot() -- i.e. any historical root.
    function contribute(
        uint256[2] calldata pA,
        uint256[2][2] calldata pB,
        uint256[2] calldata pC,
        bytes32 commitment,
        bytes32 nullifier,
        uint256 merkleRoot
    ) external {
        // 1. Commitment must be a registered patient.
        if (!registry.isRegistered(commitment)) {
            revert CommitmentNotRegistered();
        }

        // 2. Nullifier must not have been used before (replay protection).
        if (usedNullifiers[nullifier]) {
            revert NullifierAlreadyUsed();
        }

        // 3. The claimed Merkle root must match a known historical root.
        //    This prevents proofs built against fabricated trees.
        if (!attestationManager.isKnownRoot(bytes32(merkleRoot))) {
            revert StaleOrInvalidRoot();
        }

        // 4. Verify the Groth16 proof.
        //    Public signal order matches the circuit's signal declaration order:
        //      [0] commitment  (signal output -- public output)
        //      [1] nullifier   (signal output -- public output)
        //      [2] merkleRoot  (signal input  -- public input)
        uint256[3] memory pubSignals;
        pubSignals[0] = uint256(commitment);
        pubSignals[1] = uint256(nullifier);
        pubSignals[2] = merkleRoot;

        if (!researchVerifier.verifyProof(pA, pB, pC, pubSignals)) {
            revert InvalidProof();
        }

        usedNullifiers[nullifier] = true;

        emit ContributionSubmitted(commitment, nullifier);
    }
}
