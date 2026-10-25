// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Registry.sol";
import "../src/UserRegistry.sol";
import "../src/HistoryFactory.sol";
import "../src/Types.sol";
import "../src/History.sol";
import "../src/Record.sol";
import "../src/RecordPermissions.sol";
import "../src/Groth16Verifier.sol";
import "../src/ZKAuthVerifier.sol";

contract HistoryRecordPermissionsTest is Test {
    WellRegistry internal registry;
    UserRegistry internal userRegistry;
    HistoryFactory internal historyFactory;
    ZKAuthVerifier internal zkVerifier;

    address internal patient = address(0x1);
    address internal doctor = address(0x2);

    bytes32 internal patientCommitment = keccak256("patient-commitment");
    bytes32 internal patientInternalID = keccak256("patient-internal-id");

    function setUp() public {
        Groth16Verifier groth16Verifier = new Groth16Verifier();
        zkVerifier = new ZKAuthVerifier(address(groth16Verifier), 1 hours);

        userRegistry = new UserRegistry();
        historyFactory = new HistoryFactory(address(userRegistry));
        registry = new WellRegistry(address(userRegistry), address(historyFactory));
        historyFactory.setWellRegistry(address(registry));
    }

    function testHistoryContractIsCreatedOnRegistration() public {
        vm.prank(patient);
        address historyAddress = registry.register(
            patientCommitment,
            patientInternalID,
            Types.Roles.Patient
        );

        assertTrue(registry.isRegistered(patientCommitment));
        assertEq(registry.getHistoryContract(patientInternalID), historyAddress);
        assertTrue(historyAddress != address(0));
        assertGt(historyAddress.code.length, 0, "History contract not deployed");
    }

    function testCreateRecordDeploysRecordAndPermissions() public {
        vm.prank(patient);
        address historyAddress = registry.register(
            patientCommitment,
            patientInternalID,
            Types.Roles.Patient
        );

        History history = History(historyAddress);
        Types.RecordMetadata memory metadata = Types.RecordMetadata({
            recordType: Types.RecordType.LabResult,
            createdAt: 0,
            updatedAt: 0
        });

        vm.prank(patient);
        Record record = registry.createRecord(
            doctor,
            metadata,
            keccak256("content-hash"),
            "ipfs://cid"
        );

        assertGt(address(record).code.length, 0, "Record not deployed");
        assertGt(address(record.permissions()).code.length, 0, "RecordPermissions not deployed");
        assertEq(history.getRecordCount(), 1);
    }

    function testPatientCanGrantAndRevokeAccessThroughHistory() public {
        vm.prank(patient);
        address historyAddress = registry.register(
            patientCommitment,
            patientInternalID,
            Types.Roles.Patient
        );

        History history = History(historyAddress);
        Types.RecordMetadata memory metadata = Types.RecordMetadata({
            recordType: Types.RecordType.LabResult,
            createdAt: 0,
            updatedAt: 0
        });

        vm.prank(patient);
        Record record = registry.createRecord(
            doctor,
            metadata,
            keccak256("content-hash"),
            "ipfs://cid"
        );

        RecordPermissions permissions = record.permissions();

        vm.prank(patient);
        history.grantRecordAccess(record, doctor);
        assertTrue(permissions.hasAccess(doctor));

        vm.prank(patient);
        history.revokeRecordAccess(record, doctor);
        assertFalse(permissions.hasAccess(doctor));
    }
}
