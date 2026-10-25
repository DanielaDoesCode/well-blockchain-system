// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.8.0 <0.9.0;

import "forge-std/Script.sol";
import "./Registry.sol";
import "./History.sol";
import "./Record.sol";
import "./Types.sol";

contract TestEvents is Script {
    function run() external {
        vm.startBroadcast();

        WellRegistry registry = WellRegistry(vm.envAddress("REGISTRY_ADDRESS"));

        bytes32 patientCommitment = keccak256("patient-commitment");
        bytes32 patientInternalID = keccak256("patient-internal-id");

        address historyAddress = registry.register(
            patientCommitment,
            patientInternalID,
            Types.Roles.Patient
        );
        console.log("History deployed at:", historyAddress);

        Types.RecordMetadata memory metadata = Types.RecordMetadata({
            recordType: Types.RecordType.LabResult,
            createdAt: 0,
            updatedAt: 0
        });

        History history = History(historyAddress);
        Record record = registry.createRecord(
            address(0x1234),
            metadata,
            keccak256("content-hash"),
            "ipfs://cid"
        );

        console.log("Record deployed at:", address(record));

        address doctor = address(0x1234);
        history.grantRecordAccess(record, doctor);

        console.log("Doctor granted access:", doctor);

        vm.stopBroadcast();
    }
}