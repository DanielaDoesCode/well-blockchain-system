// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/UserRegistry.sol";
import "../src/HistoryFactory.sol";
import "../src/AttestationManager.sol";
import "../src/Registry.sol";
import "../src/ResearchVerifier.sol";
import "../src/ResearchContributer.sol";
import "../src/Types.sol";

/// @notice End-to-end proof test for the research-contribution flow.
///         A real Groth16 proof was generated off-chain with gen_merkle_proof.js
///         using default inputs:
///           secret = 123456789, nullifierSeed = 987654321, conditionGroup = 42
///         and baked in below. The on-chain flow is:
///           patient registers commitment -> patient.optIn ->
///           doctor attest(commitment, conditionGroup=42) ->
///           ResearchContributer.contribute(proof...) verifies against the
///           freshly-produced Merkle root and marks the nullifier used.
contract ResearchProofTest is Test {
    UserRegistry userRegistry;
    HistoryFactory historyFactory;
    WellRegistry wellRegistry;
    ResearchVerifier researchVerifier;
    ResearchContributer researchContributer;
    AttestationManager attestationManager;

    address deployer = address(0x1234);
    address patient = address(0x5678);
    address doctor = address(0xDEAD);

    // Public signals from the reference proof (conditionGroup = 42).
    // Circuit signal order: [commitment, nullifier, merkleRoot].
    bytes32 internal constant COMMITMENT =
        0x15ff5f37f52326d1cea0983b154fecf00c2344e850a4620f5cdc58eb804e606c;
    bytes32 internal constant NULLIFIER =
        0x2246d113bbdac3f6fc85512f40dedde2f923ac941795749e0e7136c65e5f1d38;
    uint256 internal constant MERKLE_ROOT =
        uint256(0x10b9888d1ffa2566d9a0669d7acff53e68456c990401e18f4f4dc4bdfa996752);
    bytes32 internal constant CONDITION_GROUP = bytes32(uint256(42));

    function setUp() public {
        vm.startPrank(deployer);

        userRegistry = new UserRegistry();
        historyFactory = new HistoryFactory(address(userRegistry));
        wellRegistry = new WellRegistry(address(userRegistry), address(historyFactory));
        historyFactory.setWellRegistry(address(wellRegistry));

        researchVerifier = new ResearchVerifier();
        attestationManager = new AttestationManager(address(userRegistry));
        researchContributer = new ResearchContributer(
            address(researchVerifier),
            address(wellRegistry),
            address(attestationManager)
        );

        vm.stopPrank();

        // Register the patient (self-registration binds the patient address to the commitment)
        vm.prank(patient);
        wellRegistry.register(COMMITMENT, keccak256("patient-internal-id"), Types.Roles.Patient);

        // Register the doctor
        vm.prank(deployer);
        userRegistry.register(doctor, keccak256("doctor-commitment"), keccak256("doctor-internal-id"), Types.Roles.Doctor);
        vm.stopPrank();

        // Patient consents to research
        vm.prank(patient);
        attestationManager.optIn(COMMITMENT);
    }

    function validProof()
        internal
        pure
        returns (
            uint256[2] memory pA,
            uint256[2][2] memory pB,
            uint256[2] memory pC
        )
    {
        pA = [
            uint256(3831480553123679039632223734555113860696781073374334892069984429094464517553),
            uint256(5064119243712126451716180150283040556479763691375002318334320987464639863780)
        ];
        pB = [
            [
                uint256(800718596290355931826373255977184946302360546577400029621437089904378712441),
                uint256(973943999067726171253068796279385427611476090132579795576397804699965299650)
            ],
            [
                uint256(11637237966353197812932870152133779746903207593445162495078632338834711710260),
                uint256(4710663813818732215155745759899292364162280385164056187775128193426282510620)
            ]
        ];
        pC = [
            uint256(18029203776822555358723332613867591587777509080694974320934597779460776975974),
            uint256(512106847857850276623762304908117486311882938437293782216335942519040263462)
        ];
    }

    function test_Proof_VerifiesOnChain() public {
        (uint256[2] memory pA, uint256[2][2] memory pB, uint256[2] memory pC) = validProof();

        uint256[3] memory pubSignals;
        pubSignals[0] = uint256(COMMITMENT);
        pubSignals[1] = uint256(NULLIFIER);
        pubSignals[2] = MERKLE_ROOT;

        assertTrue(researchVerifier.verifyProof(pA, pB, pC, pubSignals));
    }

    function test_Contribution_FullFlow() public {
        // Doctor attests conditionGroup=42 for the patient commitment.
        // This inserts leaf = Poseidon(COMMITMENT, 42) at index 0, so the
        // resulting root must equal MERKLE_ROOT from the reference proof.
        vm.prank(doctor);
        attestationManager.attest(COMMITMENT, CONDITION_GROUP);

        assertEq(attestationManager.getRoot(), MERKLE_ROOT);
        assertTrue(attestationManager.isKnownRoot(bytes32(MERKLE_ROOT)));

        (uint256[2] memory pA, uint256[2][2] memory pB, uint256[2] memory pC) = validProof();

        researchContributer.contribute(pA, pB, pC, COMMITMENT, NULLIFIER, MERKLE_ROOT);

        assertTrue(researchContributer.usedNullifiers(NULLIFIER));
    }

    function test_Contribution_RevertOnNullifierReuse() public {
        vm.prank(doctor);
        attestationManager.attest(COMMITMENT, CONDITION_GROUP);

        (uint256[2] memory pA, uint256[2][2] memory pB, uint256[2] memory pC) = validProof();

        researchContributer.contribute(pA, pB, pC, COMMITMENT, NULLIFIER, MERKLE_ROOT);

        vm.expectRevert(ResearchContributer.NullifierAlreadyUsed.selector);
        researchContributer.contribute(pA, pB, pC, COMMITMENT, NULLIFIER, MERKLE_ROOT);
    }

    function test_Contribution_RevertOnTamperedProof() public {
        vm.prank(doctor);
        attestationManager.attest(COMMITMENT, CONDITION_GROUP);

        (uint256[2] memory pA, uint256[2][2] memory pB, uint256[2] memory pC) = validProof();
        pA[0] = pA[0] + 1; // tamper

        vm.expectRevert(ResearchContributer.InvalidProof.selector);
        researchContributer.contribute(pA, pB, pC, COMMITMENT, NULLIFIER, MERKLE_ROOT);
    }

    function test_Contribution_RevertOnStaleOrInvalidRoot() public {
        vm.prank(doctor);
        attestationManager.attest(COMMITMENT, CONDITION_GROUP);

        (uint256[2] memory pA, uint256[2][2] memory pB, uint256[2] memory pC) = validProof();

        vm.expectRevert(ResearchContributer.StaleOrInvalidRoot.selector);
        researchContributer.contribute(pA, pB, pC, COMMITMENT, NULLIFIER, uint256(1));
    }

    function test_Contribution_RevertOnCommitmentNotRegistered() public {
        vm.prank(doctor);
        attestationManager.attest(COMMITMENT, CONDITION_GROUP);

        (uint256[2] memory pA, uint256[2][2] memory pB, uint256[2] memory pC) = validProof();

        bytes32 otherCommitment = bytes32(uint256(COMMITMENT) ^ 1);

        vm.expectRevert(ResearchContributer.CommitmentNotRegistered.selector);
        researchContributer.contribute(pA, pB, pC, otherCommitment, NULLIFIER, MERKLE_ROOT);
    }
}
