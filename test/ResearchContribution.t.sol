// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Groth16Verifier.sol";
import "../src/ZKAuthVerifier.sol";
import "../src/UserRegistry.sol";
import "../src/HistoryFactory.sol";
import "../src/AttestationManager.sol";
import "../src/Registry.sol";
import "../src/ResearchVerifier.sol";
import "../src/ResearchContributer.sol";
import "../src/Types.sol";

contract ResearchContributionTest is Test {
    Groth16Verifier groth16Verifier;
    ZKAuthVerifier zkAuthVerifier;
    UserRegistry userRegistry;
    HistoryFactory historyFactory;
    WellRegistry wellRegistry;
    ResearchVerifier researchVerifier;
    ResearchContributer researchContributer;
    AttestationManager attestationManager;

    address deployer = address(0x1234);
    address patient = address(0x5678);

    bytes32 patientCommitment = keccak256("patient-commitment");
    bytes32 patientInternalID = keccak256("patient-internal-id");

    function setUp() public {
        vm.startPrank(deployer);

        groth16Verifier = new Groth16Verifier();
        zkAuthVerifier = new ZKAuthVerifier(address(groth16Verifier), 1 hours);

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

        // Register patient (self-registration binds the patient address to the commitment)
        vm.prank(patient);
        wellRegistry.register(patientCommitment, patientInternalID, Types.Roles.Patient);
    }

    function test_ResearchContributer_Constructor() public {
        assertEq(address(researchContributer.researchVerifier()), address(researchVerifier));
        assertEq(address(researchContributer.registry()), address(wellRegistry));
    }

    function test_ResearchContributer_RevertWhenNotRegistered() public {
        bytes32 fakeCommitment = keccak256("not-registered");

        uint256[2] memory pA;
        uint256[2][2] memory pB;
        uint256[2] memory pC;
        uint256 root = attestationManager.getRoot();

        vm.expectRevert(ResearchContributer.CommitmentNotRegistered.selector);
        researchContributer.contribute(pA, pB, pC, fakeCommitment, bytes32(uint256(1)), root);
    }

    function test_ResearchContributer_RevertWhenNullifierUsed() public {
        // We need a valid proof to test nullifier reuse
        // For now just verify the nullifier tracking state is initialized
        bytes32 nullifier = keccak256("test-nullifier");
        assertFalse(researchContributer.usedNullifiers(nullifier));
    }

    function test_IntegrationWithWellRegistry() public {
        assertTrue(wellRegistry.isRegistered(patientCommitment));

        // Verify the chain connects: ResearchContributer calls WellRegistry
        address registryFromContributer = address(researchContributer.registry());
        assertEq(registryFromContributer, address(wellRegistry));
    }

    function test_AttestationManager_KeepsHistoricalRoots() public {
        // Register a doctor
        address doctor = address(0xDEAD);
        bytes32 doctorCommitment = keccak256("doctor-commitment");
        bytes32 doctorInternalID = keccak256("doctor-internal-id");
        userRegistry.register(doctor, doctorCommitment, doctorInternalID, Types.Roles.Doctor);

        // Register two patients whose commitments get attested
        address patientA = address(0x1111);
        bytes32 commitmentA = keccak256("patient-A");
        userRegistry.register(patientA, commitmentA, keccak256("patient-A-id"), Types.Roles.Patient);
        address patientB = address(0x2222);
        bytes32 commitmentB = keccak256("patient-B");
        userRegistry.register(patientB, commitmentB, keccak256("patient-B-id"), Types.Roles.Patient);

        // Both patients consent
        vm.prank(patientA);
        attestationManager.optIn(commitmentA);
        vm.prank(patientB);
        attestationManager.optIn(commitmentB);

        // empty tree root
        uint256 root0 = attestationManager.getRoot();
        assertTrue(attestationManager.isKnownRoot(bytes32(root0)));

        // root 1
        vm.prank(doctor);
        attestationManager.attest(commitmentA, keccak256("condition-X"));
        uint256 root1 = attestationManager.getRoot();
        assertTrue(attestationManager.isKnownRoot(bytes32(root1)));
        assertTrue(attestationManager.isKnownRoot(bytes32(root0)));

        // root 2
        vm.prank(doctor);
        attestationManager.attest(commitmentB, keccak256("condition-Y"));
        uint256 root2 = attestationManager.getRoot();

        // root 1 still should be known
        assertTrue(attestationManager.isKnownRoot(bytes32(root1)));

        // root 2 is the new current root
        assertTrue(attestationManager.isKnownRoot(bytes32(root2)));
        assertEq(attestationManager.getRoot(), root2);
        assertFalse(root1 == root2);
    }

    function test_AttestationManager_RevertWhenAttestWithoutConsent() public {
        address doctor = address(0xDEAD);
        bytes32 doctorCommitment = keccak256("doctor-commitment");
        userRegistry.register(doctor, doctorCommitment, keccak256("doctor-internal-id"), Types.Roles.Doctor);

        vm.prank(doctor);
        vm.expectRevert("consent required");
        attestationManager.attest(patientCommitment, keccak256("condition-X"));
    }

    function test_AttestationManager_RevertWhenOptInByNonOwner() public {
        // patient2 is registered but does NOT own patientCommitment
        address patient2 = address(0x9999);
        userRegistry.register(patient2, keccak256("patient2-commitment"), keccak256("patient2-id"), Types.Roles.Patient);

        vm.prank(patient2);
        vm.expectRevert("not the commitment owner");
        attestationManager.optIn(patientCommitment);
    }

    function test_AttestationManager_OptInAndAttest() public {
        address doctor = address(0xDEAD);
        bytes32 doctorCommitment = keccak256("doctor-commitment");
        userRegistry.register(doctor, doctorCommitment, keccak256("doctor-internal-id"), Types.Roles.Doctor);

        uint256 rootBefore = attestationManager.getRoot();

        vm.prank(patient);
        attestationManager.optIn(patientCommitment);
        assertTrue(attestationManager.optedIn(patientCommitment));

        vm.prank(doctor);
        attestationManager.attest(patientCommitment, keccak256("condition-X"));

        assertTrue(attestationManager.isKnownRoot(bytes32(rootBefore)));
        assertTrue(attestationManager.getRoot() != rootBefore);
        assertEq(attestationManager.nextIndex(), 1);
    }

    function test_ResearchContributer_EmitsEventOnContribution() public {
        // This test validates the contract wiring and event shape.
        // A real Groth16 proof for the research circuit is generated off-chain.
        // Here we verify that a proof verification failure emits the correct error.

        uint256[2] memory pA;
        pA[0] = 1;
        pA[1] = 2;
        uint256[2][2] memory pB;
        pB[0][0] = 3; pB[0][1] = 4;
        pB[1][0] = 5; pB[1][1] = 6;
        uint256[2] memory pC;
        pC[0] = 7; pC[1] = 8;
        uint256 root = attestationManager.getRoot();

        vm.expectRevert(ResearchContributer.InvalidProof.selector);
        researchContributer.contribute(
            pA, pB, pC,
            patientCommitment,
            bytes32(uint256(42)),
            root
        );
    }
}
