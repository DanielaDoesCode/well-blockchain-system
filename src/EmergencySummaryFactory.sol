// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.20;

import "./Types.sol";
import "./EmergencySummary.sol";
import "./UserRegistry.sol";

contract EmergencySummaryFactory {

    UserRegistry public userRegistry;
    address public wellRegistry;
    address public riskAssessment;

    mapping(address => address) public patientToEmergencySummary;
    mapping(bytes32 => address) public commitmentToEmergencySummary;

    error Unauthorized();

    event EmergencySummaryCreated(
        bytes32 indexed commitment,
        address indexed patient,
        address emergencySummary
    );

    modifier onlyWellRegistry() {
        if (msg.sender != wellRegistry) revert Unauthorized();
        _;
    }

    constructor(address _userRegistry) {
        userRegistry = UserRegistry(_userRegistry);
    }

    function setWellRegistry(address _wellRegistry) external {
        if (wellRegistry != address(0)) revert Unauthorized();
        wellRegistry = _wellRegistry;
    }

    function setRiskAssessment(address _riskAssessment) external {
        if (riskAssessment != address(0)) revert Unauthorized();
        riskAssessment = _riskAssessment;
    }

    function createEmergencySummary(bytes32 commitment, address patientAddress)
        external onlyWellRegistry returns (address)
    {
        EmergencySummary summary = new EmergencySummary(patientAddress, wellRegistry, riskAssessment);
        commitmentToEmergencySummary[commitment] = address(summary);
        patientToEmergencySummary[patientAddress] = address(summary);

        emit EmergencySummaryCreated(commitment, patientAddress, address(summary));
        return address(summary);
    }

    function getEmergencySummaryByAddress(address patientAddress)
        external view returns (address)
    {
        return patientToEmergencySummary[patientAddress];
    }

    function getEmergencySummaryByCommitment(bytes32 commitment)
        external view returns (address)
    {
        return commitmentToEmergencySummary[commitment];
    }
}
