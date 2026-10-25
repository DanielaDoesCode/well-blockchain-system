// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.20;

import "./Types.sol";
import "./UserRegistry.sol";
import "./IncidentManager.sol";
import "./EmergencySummary.sol";
import "./EmergencySummaryFactory.sol";

contract RiskAssessment {

    IncidentManager public immutable incidentManager;
    UserRegistry public immutable userRegistry;
    EmergencySummaryFactory public immutable emergencySummaryFactory;

    uint256 public constant ACCESS_DURATION = 1 hours;

    error Unauthorized();
    error IncidentNotActive();
    error NotAssignedResponder();
    error NotFirstResponder();
    error LocationMismatch();
    error InvalidPatient();
    error PatientMismatch();

    event EmergencyGranted(
        address indexed responder,
        bytes32 indexed incidentId,
        uint256 expiresAt
    );

    // One-time binding from incident to the patient whose emergency data
    // that incident may unlock. Bound by the first successful grant.
    mapping(bytes32 => bytes32) private incidentPatient;

    constructor(
        address _incidentManager,
        address _userRegistry,
        address _emergencySummaryFactory
    ) {
        incidentManager = IncidentManager(_incidentManager);
        userRegistry = UserRegistry(_userRegistry);
        emergencySummaryFactory = EmergencySummaryFactory(_emergencySummaryFactory);
    }

    function emergencyGrant(
        address emergencySummary,
        bytes32 incidentId,
        bytes32 locationHash,
        bytes32 patientCommitment,
        bytes calldata encryptedKey
    ) external {
        if (!incidentManager.isIncidentActive(incidentId)) {
            revert IncidentNotActive();
        }

        if (userRegistry.getRoleByAddress(msg.sender) != Types.Roles.FirstResponder) {
            revert NotFirstResponder();
        }

        if (!incidentManager.isResponderInIncident(incidentId, msg.sender)) {
            revert NotAssignedResponder();
        }

        Types.Incident memory incident = incidentManager.getIncident(incidentId);
        if (incident.locationHash != locationHash) {
            revert LocationMismatch();
        }

        if (emergencySummaryFactory.getEmergencySummaryByCommitment(patientCommitment) != emergencySummary) {
            revert InvalidPatient();
        }

        bytes32 bound = incidentPatient[incidentId];
        if (bound == bytes32(0)) {
            incidentPatient[incidentId] = patientCommitment;
        } else if (bound != patientCommitment) {
            revert PatientMismatch();
        }

        uint256 expiry = block.timestamp + ACCESS_DURATION;

        EmergencySummary(emergencySummary).grantAccess(
            msg.sender,
            encryptedKey,
            ACCESS_DURATION
        );

        emit EmergencyGranted(msg.sender, incidentId, expiry);
    }

    function getIncidentPatient(bytes32 incidentId) external view returns (bytes32) {
        return incidentPatient[incidentId];
    }
}
