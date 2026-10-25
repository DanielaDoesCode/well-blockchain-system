// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.20;

import "./Types.sol";
import "./UserRegistry.sol";
import "./HistoryFactory.sol";
import "./EmergencySummaryFactory.sol";
import "./History.sol";
import "./Record.sol";

contract WellRegistry {

    UserRegistry public userRegistry;
    HistoryFactory public historyFactory;
    EmergencySummaryFactory public emergencySummaryFactory;
    address public incidentManager;
    address public riskAssessment;

    struct User {
        bytes32      commitment;
        bytes32      internalID;
        Types.Roles  role;
        address      userAddress;
        address      historyContract;
    }

    error NotRegistered();
    error Unauthorized();
    error AlreadySet();
    error NoSelfRegistrationDoctor();
    error AdminAlreadyExists();

    bool public firstAdminRegistered;

    modifier onlyAdmin() {
        if (userRegistry.getRoleByAddress(msg.sender) != Types.Roles.Admin) {
            revert Unauthorized();
        }
        _;
    }

    constructor(address _userRegistry, address _historyFactory) {
        userRegistry = UserRegistry(_userRegistry);
        historyFactory = HistoryFactory(_historyFactory);
    }

    function register(
        bytes32     commitment,
        bytes32     internalID,
        Types.Roles role
    ) external returns (address historyContract) {
        if (role == Types.Roles.Patient) {
            userRegistry.register(msg.sender, commitment, internalID, role);

            historyContract = historyFactory.createHistory(commitment, msg.sender);

            if (address(emergencySummaryFactory) != address(0)) {
            emergencySummaryFactory.createEmergencySummary(commitment, msg.sender);
            }

            return historyContract;
        }

        if (role == Types.Roles.Admin) {
            if (firstAdminRegistered) revert AdminAlreadyExists();

            userRegistry.register(msg.sender, commitment, internalID, role);
            firstAdminRegistered = true;

            return historyContract;
        }

        revert NoSelfRegistrationDoctor();
    }

    function registerByAdmin(
        address     userAddress,
        bytes32     commitment,
        bytes32     internalID,
        Types.Roles role
    ) external onlyAdmin returns (address historyContract) {
        
        userRegistry.register(userAddress, commitment, internalID, role);

        if(role == Types.Roles.Patient){
            historyContract = historyFactory.createHistory(commitment, userAddress);

        if (address(emergencySummaryFactory) != address(0)) {
            emergencySummaryFactory.createEmergencySummary(commitment, userAddress);
        }

            return historyContract;

        }

        return historyContract;

    }

    function revoke(bytes32 commitment) external {
        userRegistry.revoke(commitment);
    }

    function createRecord(
        address                    uploader,
        Types.RecordMetadata memory metadata,
        bytes32                    contentHash,
        string memory              storageReference
    ) external returns (Record) {
        address historyAddress = historyFactory.getHistoryContractByAddress(msg.sender);
        return History(historyAddress).createRecord(
            uploader,
            metadata,
            contentHash,
            storageReference
        );
    }

    // ─────────────────────────────────────────────────────────────
    // One-time setters for emergency infrastructure
    // ─────────────────────────────────────────────────────────────

    function setEmergencySummaryFactory(address _emergencySummaryFactory) external {
        if (address(emergencySummaryFactory) != address(0)) revert AlreadySet();
        emergencySummaryFactory = EmergencySummaryFactory(_emergencySummaryFactory);
    }

    function setIncidentManager(address _incidentManager) external {
        if (incidentManager != address(0)) revert AlreadySet();
        incidentManager = _incidentManager;
    }

    function setRiskAssessment(address _riskAssessment) external {
        if (riskAssessment != address(0)) revert AlreadySet();
        riskAssessment = _riskAssessment;
    }

    // ─────────────────────────────────────────────────────────────
    // Read Functions
    // ─────────────────────────────────────────────────────────────

    function isRegistered(bytes32 commitment) external view returns (bool) {
        return userRegistry.isRegistered(commitment);
    }

    function getCommitmentByInternalID(bytes32 internalID)
        external view returns (bytes32 commitment)
    {
        return userRegistry.getCommitmentByInternalID(internalID);
    }

    function getCommitmentByAddress(address userAddress)
        external view returns (bytes32 commitment)
    {
        return userRegistry.getCommitmentByAddress(userAddress);
    }

    function getUserByAddress(address userAddress)
        external view returns (User memory)
    {
        bytes32 commitment = userRegistry.getCommitmentByAddress(userAddress);
        UserRegistry.User memory u = userRegistry.getUser(commitment);
        address h = historyFactory.commitmentToHistory(commitment);
        return User({
            commitment:      u.commitment,
            internalID:      u.internalID,
            role:            u.role,
            userAddress:     u.userAddress,
            historyContract: h
        });
    }

    function getUser(bytes32 commitment)
        external view returns (User memory)
    {
        UserRegistry.User memory u = userRegistry.getUser(commitment);
        address h = historyFactory.commitmentToHistory(commitment);
        return User({
            commitment:      u.commitment,
            internalID:      u.internalID,
            role:            u.role,
            userAddress:     u.userAddress,
            historyContract: h
        });
    }

    function getRole(bytes32 commitment)
        external view returns (Types.Roles)
    {
        return userRegistry.getRole(commitment);
    }

    function getRoleByAddress(address userAddress)
        external view returns (Types.Roles)
    {
        return userRegistry.getRoleByAddress(userAddress);
    }

    function getHistoryContract(bytes32 internalID)
        external view returns (address)
    {
        return historyFactory.getHistoryContract(internalID);
    }

    function getHistoryContractByAddress(address userAddress)
        external view returns (address)
    {
        return historyFactory.getHistoryContractByAddress(userAddress);
    }

    function getEmergencySummaryByAddress(address patientAddress)
        external view returns (address)
    {
        return emergencySummaryFactory.getEmergencySummaryByAddress(patientAddress);
    }

    function getIncidentManager() external view returns (address) {
        return incidentManager;
    }

    function getRiskAssessment() external view returns (address) {
        return riskAssessment;
    }
}
