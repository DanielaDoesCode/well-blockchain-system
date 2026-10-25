// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Registry.sol";
import "../src/UserRegistry.sol";
import "../src/HistoryFactory.sol";
import "../src/ZKAuthVerifier.sol";
import "../src/Groth16Verifier.sol";
import "../src/Types.sol";

contract ZKAuthIntegrationTest is Test {
    WellRegistry public registry;
    UserRegistry public userRegistry;
    HistoryFactory public historyFactory;
    ZKAuthVerifier public zkVerifier;
    Groth16Verifier public groth16Verifier;

    address public patient = address(0x1);
    bytes32 public patientCommitment = keccak256("patient-commitment");
    bytes32 public patientInternalID = keccak256("patient-internal-id");

    function setUp() public {
        groth16Verifier = new Groth16Verifier();
        zkVerifier = new ZKAuthVerifier(address(groth16Verifier), 1 hours);

        userRegistry = new UserRegistry();
        historyFactory = new HistoryFactory(address(userRegistry));
        registry = new WellRegistry(address(userRegistry), address(historyFactory));
        historyFactory.setWellRegistry(address(registry));
    }

    function testRegisterPatientAndGetCommitmentByAddress() public {
        vm.prank(patient);
        registry.register(patientCommitment, patientInternalID, Types.Roles.Patient);

        assertEq(registry.getCommitmentByAddress(patient), patientCommitment);
        assertEq(registry.getCommitmentByInternalID(patientInternalID), patientCommitment);
    }

    function testGetUserByAddress() public {
        vm.prank(patient);
        registry.register(patientCommitment, patientInternalID, Types.Roles.Patient);

        WellRegistry.User memory user = registry.getUserByAddress(patient);
        assertEq(user.commitment, patientCommitment);
        assertEq(user.internalID, patientInternalID);
        assertEq(uint256(user.role), uint256(Types.Roles.Patient));
        assertTrue(user.historyContract != address(0));
    }

    function testGetRoleByAddress() public {
        vm.prank(patient);
        registry.register(patientCommitment, patientInternalID, Types.Roles.Patient);

        assertEq(uint256(registry.getRoleByAddress(patient)), uint256(Types.Roles.Patient));
    }

    function testGetHistoryContractByAddress() public {
        vm.prank(patient);
        address history = registry.register(patientCommitment, patientInternalID, Types.Roles.Patient);

        assertEq(registry.getHistoryContractByAddress(patient), history);
    }

    function testCannotRegisterTwice() public {
        vm.prank(patient);
        registry.register(patientCommitment, patientInternalID, Types.Roles.Patient);

        vm.prank(patient);
        vm.expectRevert(UserRegistry.AlreadyRegistered.selector);
        registry.register(patientCommitment, patientInternalID, Types.Roles.Patient);
    }
}
