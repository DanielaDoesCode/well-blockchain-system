// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.20;

import "./Types.sol";
import "./UserRegistry.sol";

contract IncidentManager {

    UserRegistry public immutable userRegistry;

    uint256 public constant MAX_INCIDENT_DURATION = 1 hours;

    mapping(bytes32 => Types.Incident) private incidents;
    mapping(bytes32 => mapping(address => bool)) private responderInIncident;
    bytes32[] private incidentIds;

    error Unauthorized();
    error InvalidIncident();
    error IncidentExpired();
    error IncidentNotActive();

    
    event IncidentCreated(
        bytes32 indexed id,
        address indexed hospital,
        uint48 expiresAt
    );

    event RespondersAdded(
        bytes32 indexed id,
        address[] responders
    );

    event IncidentDeactivated(bytes32 indexed id);

    modifier onlyHospital() {
        if (userRegistry.getRoleByAddress(msg.sender) != Types.Roles.Hospital) {
            revert Unauthorized();
        }
        _;
    }

    modifier onlyIncidentHospital(bytes32 incidentId) {
        Types.Incident storage incident = incidents[incidentId];
        if (incident.hospital == address(0)) revert InvalidIncident();
        if (msg.sender != incident.hospital) revert Unauthorized();
        _;
    }

    constructor(address _userRegistry) {
        userRegistry = UserRegistry(_userRegistry);
    }

    function createIncident(
        bytes32 locationHash,
        uint48 expiresAt
    ) external onlyHospital returns (bytes32) {
        if (expiresAt <= uint48(block.timestamp)) revert InvalidIncident();
        if (expiresAt > uint48(block.timestamp) + uint48(MAX_INCIDENT_DURATION)) {
            revert InvalidIncident();
        }

        bytes32 id = keccak256(abi.encodePacked(
            block.timestamp,
            msg.sender,
            locationHash,
            block.prevrandao
        ));

        incidents[id] = Types.Incident({
            id: id,
            hospital: msg.sender,
            locationHash: locationHash,
            createdAt: uint48(block.timestamp),
            expiresAt: expiresAt,
            active: true
        });

        incidentIds.push(id);

        emit IncidentCreated(id, msg.sender, expiresAt);
        return id;
    }

    function addResponders(
        bytes32 incidentId,
        address[] calldata responders
    ) external onlyIncidentHospital(incidentId) {
        Types.Incident storage incident = incidents[incidentId];
        if (!incident.active) revert IncidentNotActive();
        if (block.timestamp >= incident.expiresAt) revert IncidentExpired();

        for (uint256 i = 0; i < responders.length; i++) {
            responderInIncident[incidentId][responders[i]] = true;
        }

        emit RespondersAdded(incidentId, responders);
    }

    function deactivateIncident(bytes32 incidentId)
        external onlyIncidentHospital(incidentId)
    {
        incidents[incidentId].active = false;
        emit IncidentDeactivated(incidentId);
    }


    function isIncidentActive(bytes32 incidentId) external view returns (bool) {
        Types.Incident storage incident = incidents[incidentId];
        if (incident.hospital == address(0)) return false;
        return incident.active && block.timestamp < incident.expiresAt;
    }

    function isResponderInIncident(
        bytes32 incidentId,
        address responder
    ) external view returns (bool) {
        return responderInIncident[incidentId][responder];
    }

    function getIncident(bytes32 incidentId)
        external view returns (Types.Incident memory)
    {
        return incidents[incidentId];
    }

    function getIncidentCount() external view returns (uint256) {
        return incidentIds.length;
    }
}
