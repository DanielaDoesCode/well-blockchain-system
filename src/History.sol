// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.8.0 <0.9.0;

import "./Types.sol";
import "./RecordPermissions.sol";
import "./Record.sol";

/**
 * History Contract
 * Stores references to all medical records for a specific patient
 * and redirects operations to the Record contracts.
 *
 * Deployed once per patient by WellRegistry at registration time.
 * The access control backend address is set at deploy time and is
 * the only caller allowed to create records on the patient's behalf.
 */
contract History {

    address public immutable owner;         // patient address
    address public immutable accessControl; // trusted backend — only caller for createRecord
    address public immutable policyVerifier; // Groth16 verifier for policy_match proofs

    Record[] private records;
    mapping(Record => bool) private isOwnedRecord;

    error Unauthorized();
    error RecordNotFound();

    event RecordCreated(
        address indexed record,
        address indexed uploader,
        Types.RecordType recordType
    );

    event AccessGranted(
        address indexed record,
        address indexed user
    );

    event AccessRevoked(
        address indexed record,
        address indexed user
    );

    modifier onlyOwner() {
        if (msg.sender != owner) revert Unauthorized();
        _;
    }

    modifier onlyAccessControl() {
        if (msg.sender != accessControl) revert Unauthorized();
        _;
    }

    modifier validRecord(Record record) {
        if (!isOwnedRecord[record]) revert RecordNotFound();
        _;
    }

    constructor(address _owner, address _accessControl, address _policyVerifier) {
        owner           = _owner;
        accessControl   = _accessControl;
        policyVerifier  = _policyVerifier;
    }


    function createRecord(
        address                    _uploader,
        Types.RecordMetadata memory _metadata,
        bytes32                    _contentHash,
        string memory              _storageReference
    ) external onlyAccessControl returns (Record) {
        Record newRecord = new Record(
            owner,
            _uploader,
            _metadata,
            _contentHash,
            _storageReference,
            policyVerifier
        );

        records.push(newRecord);
        isOwnedRecord[newRecord] = true;

        emit RecordCreated(address(newRecord), _uploader, _metadata.recordType);

        return newRecord;
    }

    function grantRecordAccess(Record record, address user)
        external
        onlyOwner
        validRecord(record)
    {
        record.grantAccess(user);
        emit AccessGranted(address(record), user);
    }

    function revokeRecordAccess(Record record, address user)
        external
        onlyOwner
        validRecord(record)
    {
        record.revokeAccess(user);
        emit AccessRevoked(address(record), user);
    }

    function getRecord(Record record)
        external
        view
        validRecord(record)
        returns (
            Types.RecordMetadata memory,
            bytes32,
            string memory,
            address,
            address
        )
    {
        return record.getRecord();
    }

    function verifyRecordIntegrity(Record record, bytes32 hash)
        external
        view
        validRecord(record)
        returns (bool)
    {
        return record.verifyIntegrity(hash);
    }

    function hasAccessToRecord(Record record, address user)
        external
        view
        validRecord(record)
        returns (bool)
    {
        return record.permissions().hasAccess(user);
    }

    function getAllRecords() external view returns (Record[] memory) {
        return records;
    }

    function getRecordCount() external view returns (uint256) {
        return records.length;
    }
}