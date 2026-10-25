// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.20;

import "./UserRegistry.sol";
import "./lib/PoseidonT3.sol";

/// @notice Incremental Merkle tree of doctor-attested (commitment, conditionGroup) pairs.
///         Tree depth 20 matches the ResearchContribute(20) circuit parameter exactly.
///         Leaf value = PoseidonT3.hash([commitment, conditionGroup]) -- identical to the
///         circuit's `leaf = Poseidon(commitment, conditionGroup)` at lines 28-30.
///         Empty leaf value = 0 (field element zero).
///         The precomputed zeros array holds the root of an all-zero subtree at each level,
///         so no on-chain initialisation loop is needed.
contract AttestationManager {

    // -- Constants ------------------------------------------------------------

    uint256 public constant TREE_DEPTH = 20;
    uint256 public constant MAX_LEAVES = 1 << TREE_DEPTH; // 2^20 = 1,048,576

    /// @dev zeros[i] = Poseidon root of a subtree of depth i where every leaf = 0.
    ///      Computed off-chain once and hardcoded here.
    ///      zeros[0] = 0  (the empty leaf itself)
    ///      zeros[i] = PoseidonT3.hash([zeros[i-1], zeros[i-1]])  for i >= 1
    ///
    ///      To regenerate:
    ///        const { buildPoseidon } = require("circomlibjs");
    ///        const p = await buildPoseidon();
    ///        let z = 0n;
    ///        for (let i = 0; i < 20; i++) {
    ///          console.log(z.toString(16));
    ///          z = p.F.toObject(p([z, z]));
    ///        }
    function _zeros(uint256 i) internal pure returns (uint256) {
        if (i == 0)  return 0;
        if (i == 1)  return 0x2098f5fb9e239eab3ceac3f27b81e481dc3124d55ffed523a839ee8446b64864;
        if (i == 2)  return 0x1069673dcdb12263df301a6ff584a7ec261a44cb9dc68df067a4774460b1f1e1;
        if (i == 3)  return 0x18f43331537ee2af2e3d758d50f72106467c6eea50371dd528d57eb2b856d238;
        if (i == 4)  return 0x7f9d837cb17b0d36320ffe93ba52345f1b728571a568265caac97559dbc952a;
        if (i == 5)  return 0x2b94cf5e8746b3f5c9631f4c5df32907a699c58c94b2ad4d7b5cec1639183f55;
        if (i == 6)  return 0x2dee93c5a666459646ea7d22cca9e1bcfed71e6951b953611d11dda32ea09d78;
        if (i == 7)  return 0x78295e5a22b84e982cf601eb639597b8b0515a88cb5ac7fa8a4aabe3c87349d;
        if (i == 8)  return 0x2fa5e5f18f6027a6501bec864564472a616b2e274a41211a444cbe3a99f3cc61;
        if (i == 9)  return 0xe884376d0d8fd21ecb780389e941f66e45e7acce3e228ab3e2156a614fcd747;
        if (i == 10) return 0x1b7201da72494f1e28717ad1a52eb469f95892f957713533de6175e5da190af2;
        if (i == 11) return 0x1f8d8822725e36385200c0b201249819a6e6e1e4650808b5bebc6bface7d7636;
        if (i == 12) return 0x2c5d82f66c914bafb9701589ba8cfcfb6162b0a12acf88a8d0879a0471b5f85a;
        if (i == 13) return 0x14c54148a0940bb820957f5adf3fa1134ef5c4aaa113f4646458f270e0bfbfd0;
        if (i == 14) return 0x190d33b12f986f961e10c0ee44d8b9af11be25588cad89d416118e4bf4ebe80c;
        if (i == 15) return 0x22f98aa9ce704152ac17354914ad73ed1167ae6596af510aa5b3649325e06c92;
        if (i == 16) return 0x2a7c7c9b6ce5880b9f6f228d72bf6a575a526f29c66ecceef8b753d38bba7323;
        if (i == 17) return 0x2e8186e558698ec1c67af9c14d463ffc470043c9c2988b954d75dd643f36b992;
        if (i == 18) return 0x0f57c5571e9a4eab49e2c8cf050dae948aef6ead647392273546249d1c1ff10f;
        if (i == 19) return 0x1830ee67b5fb554ad5f63d4388800e1cfe78e310697d46e43c9ce36134f72cca;
        if (i == 20) return 0x2134e76ac5d21aab186c2be1dd8f84ee880a1e46eaf712f9d371b6df22191f3e;
        revert("zeros: index out of range");
    }

    // -- State ----------------------------------------------------------------

    UserRegistry public userRegistry;

    /// @dev Standard incremental tree state: filledSubtrees[i] is the last
    ///      known non-zero node at level i on the rightmost populated path.
    uint256[20] public filledSubtrees;

    /// @dev Index of the next free leaf slot.
    uint256 public nextIndex;

    /// @dev Current Merkle root -- updated on every insertion.
    uint256 public currentRoot;

    /// @notice Every root that has ever been the canonical tree root.
    mapping(bytes32 => bool) public knownRoots;

    /// @notice GDPR-aligned consent: patient opted in to research for a commitment.
    ///         Only enforced at attest time — the incremental tree cannot delete
    ///         already-inserted leaves, so withdrawing consent does not remove
    ///         prior attestations.
    mapping(bytes32 => bool) public optedIn;

    // -- Events ---------------------------------------------------------------

    event Attested(
        address indexed doctor,
        bytes32 indexed commitment,
        bytes32 indexed conditionGroup,
        uint256 leafIndex,
        uint256 merkleRoot
    );

    event ConsentGiven(
        address indexed patient,
        bytes32 indexed commitment
    );

    // -- Constructor ----------------------------------------------------------

    constructor(address _userRegistry) {
        userRegistry = UserRegistry(_userRegistry);

        // Initialise filledSubtrees to the zero values for each level.
        for (uint256 i = 0; i < TREE_DEPTH; i++) {
            filledSubtrees[i] = _zeros(i);
        }

        // The root of a completely empty tree of depth 20.
        currentRoot = _zeros(TREE_DEPTH);
        knownRoots[bytes32(currentRoot)] = true;
    }


    /// @notice Called by a patient (the owner of `commitment`) to consent to
    ///         having it attested by doctors. Consent is required by `attest`.
    function optIn(bytes32 commitment) external {
        if (userRegistry.getUserByAddress(msg.sender).commitment != commitment) {
            revert("not the commitment owner");
        }
        optedIn[commitment] = true;
        emit ConsentGiven(msg.sender, commitment);
    }

    /// @notice Called by a registered doctor to attest that a patient
    ///         (identified by their research commitment) has a given conditionGroup.
    ///         Inserts leaf = PoseidonT3.hash([commitment, conditionGroup]) into the
    ///         tree and updates currentRoot.
    function attest(
        bytes32 commitment,
        bytes32 conditionGroup
    ) external {
        if (!userRegistry.isDoctor(msg.sender)) revert("not a doctor");
        if (!optedIn[commitment]) revert("consent required");
        if (nextIndex >= MAX_LEAVES) revert("tree is full");

        uint256 leaf = PoseidonT3.hash(
            [uint256(commitment), uint256(conditionGroup)]
        );

        uint256 index = nextIndex;
        uint256 currentLevelHash = leaf;

        for (uint256 i = 0; i < TREE_DEPTH; i++) {
            uint256 left;
            uint256 right;

            if (index % 2 == 0) {
                // current node is a left child; right sibling is empty at this level
                left  = currentLevelHash;
                right = _zeros(i);
                filledSubtrees[i] = currentLevelHash;
            } else {
                // current node is a right child; left sibling was stored on prior insert
                left  = filledSubtrees[i];
                right = currentLevelHash;
            }

            currentLevelHash = PoseidonT3.hash([left, right]);
            index >>= 1;
        }

        currentRoot = currentLevelHash;
        knownRoots[bytes32(currentRoot)] = true;
        nextIndex++;

        emit Attested(
            msg.sender,
            commitment,
            conditionGroup,
            nextIndex - 1,
            currentRoot
        );
    }

    function getRoot() external view returns (uint256) {
        return currentRoot;
    }

    function isKnownRoot(bytes32 root) external view returns (bool) {
        return knownRoots[root];
    }
}
