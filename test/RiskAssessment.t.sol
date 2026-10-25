// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Registry.sol";
import "../src/UserRegistry.sol";
import "../src/HistoryFactory.sol";
import "../src/EmergencySummaryFactory.sol";
import "../src/IncidentManager.sol";
import "../src/RiskAssessment.sol";
import "../src/EmergencySummary.sol";
import "../src/Types.sol";

contract RiskAssessmentTest is Test {
    WellRegistry internal registry;
    UserRegistry internal userRegistry;
    HistoryFactory internal historyFactory;
    EmergencySummaryFactory internal esFactory;
    IncidentManager internal incidentManager;
    RiskAssessment internal riskAssessment;

    address internal admin = address(0xA);
    address internal hospital = address(0x1);
    address internal responder = address(0x2);
    address internal otherResponder = address(0x3);
    address internal doctor = address(0x4);
    address internal patient = address(0x5);
    address internal patient2 = address(0x6);

    bytes32 internal patientCommitment = keccak256("patient-commitment");
    bytes32 internal patient2Commitment = keccak256("patient2-commitment");
    bytes32 internal locationHash = keccak256("emergency-location");
    bytes32 internal incidentId;

    EmergencySummary internal summary;
    EmergencySummary internal summary2;

    function setUp() public {
        userRegistry = new UserRegistry();
        historyFactory = new HistoryFactory(address(userRegistry));
        registry = new WellRegistry(address(userRegistry), address(historyFactory));
        historyFactory.setWellRegistry(address(registry));

        esFactory = new EmergencySummaryFactory(address(userRegistry));
        esFactory.setWellRegistry(address(registry));

        incidentManager = new IncidentManager(address(userRegistry));
        riskAssessment = new RiskAssessment(
            address(incidentManager),
            address(userRegistry),
            address(esFactory)
        );

        registry.setIncidentManager(address(incidentManager));
        registry.setEmergencySummaryFactory(address(esFactory));
        registry.setRiskAssessment(address(riskAssessment));
        esFactory.setRiskAssessment(address(riskAssessment));

        vm.prank(admin);
        registry.register(bytes32(uint256(111)), bytes32(uint256(222)), Types.Roles.Admin);

        vm.prank(admin);
        registry.registerByAdmin(hospital, bytes32(uint256(1)), bytes32(uint256(2)), Types.Roles.Hospital);
        vm.prank(admin);
        registry.registerByAdmin(responder, bytes32(uint256(3)), bytes32(uint256(4)), Types.Roles.FirstResponder);
        vm.prank(admin);
        registry.registerByAdmin(otherResponder, bytes32(uint256(5)), bytes32(uint256(6)), Types.Roles.FirstResponder);
        vm.prank(admin);
        registry.registerByAdmin(doctor, bytes32(uint256(7)), bytes32(uint256(8)), Types.Roles.Doctor);

        vm.prank(patient);
        registry.register(patientCommitment, bytes32(uint256(9)), Types.Roles.Patient);
        summary = EmergencySummary(registry.getEmergencySummaryByAddress(patient));

        vm.prank(patient2);
        registry.register(patient2Commitment, bytes32(uint256(10)), Types.Roles.Patient);
        summary2 = EmergencySummary(registry.getEmergencySummaryByAddress(patient2));
    }

    function _createActiveIncident() internal returns (bytes32 id) {
        vm.prank(hospital);
        id = incidentManager.createIncident(locationHash, uint48(block.timestamp + 1 hours));

        address[] memory responders = new address[](1);
        responders[0] = responder;
        vm.prank(hospital);
        incidentManager.addResponders(id, responders);
    }

    function testHappyPathGrantBindsPatient() public {
        incidentId = _createActiveIncident();

        vm.prank(responder);
        riskAssessment.emergencyGrant(
            address(summary),
            incidentId,
            locationHash,
            patientCommitment,
            abi.encodePacked("encrypted-key")
        );

        assertTrue(summary.hasAccess(responder));
        assertEq(riskAssessment.getIncidentPatient(incidentId), patientCommitment);
    }

    function testWrongLocationHashReverts() public {
        incidentId = _createActiveIncident();

        vm.prank(responder);
        vm.expectRevert(RiskAssessment.LocationMismatch.selector);
        riskAssessment.emergencyGrant(
            address(summary),
            incidentId,
            keccak256("other-location"),
            patientCommitment,
            abi.encodePacked("encrypted-key")
        );

        assertFalse(summary.hasAccess(responder));
        assertEq(riskAssessment.getIncidentPatient(incidentId), bytes32(0));
    }

    function testSummaryNotMatchingCommitmentReverts() public {
        incidentId = _createActiveIncident();

        vm.prank(responder);
        vm.expectRevert(RiskAssessment.InvalidPatient.selector);
        riskAssessment.emergencyGrant(
            address(summary),
            incidentId,
            locationHash,
            patient2Commitment,
            abi.encodePacked("encrypted-key")
        );

        assertFalse(summary.hasAccess(responder));
    }

    function testCrossPatientGrantOnSameIncidentReverts() public {
        incidentId = _createActiveIncident();

        vm.prank(responder);
        riskAssessment.emergencyGrant(
            address(summary),
            incidentId,
            locationHash,
            patientCommitment,
            abi.encodePacked("encrypted-key")
        );

        vm.prank(responder);
        vm.expectRevert(RiskAssessment.PatientMismatch.selector);
        riskAssessment.emergencyGrant(
            address(summary2),
            incidentId,
            locationHash,
            patient2Commitment,
            abi.encodePacked("encrypted-key")
        );

        assertTrue(summary.hasAccess(responder));
        assertFalse(summary2.hasAccess(responder));
    }

    function testUnassignedResponderReverts() public {
        incidentId = _createActiveIncident();

        vm.prank(otherResponder);
        vm.expectRevert(RiskAssessment.NotAssignedResponder.selector);
        riskAssessment.emergencyGrant(
            address(summary),
            incidentId,
            locationHash,
            patientCommitment,
            abi.encodePacked("encrypted-key")
        );
    }

    function testNonFirstResponderReverts() public {
        incidentId = _createActiveIncident();

        vm.prank(doctor);
        vm.expectRevert(RiskAssessment.NotFirstResponder.selector);
        riskAssessment.emergencyGrant(
            address(summary),
            incidentId,
            locationHash,
            patientCommitment,
            abi.encodePacked("encrypted-key")
        );
    }

    function testDeactivatedIncidentReverts() public {
        incidentId = _createActiveIncident();

        vm.prank(hospital);
        incidentManager.deactivateIncident(incidentId);

        vm.prank(responder);
        vm.expectRevert(RiskAssessment.IncidentNotActive.selector);
        riskAssessment.emergencyGrant(
            address(summary),
            incidentId,
            locationHash,
            patientCommitment,
            abi.encodePacked("encrypted-key")
        );

        assertFalse(summary.hasAccess(responder));
    }

    function testExpiredIncidentReverts() public {
        incidentId = _createActiveIncident();
        uint48 createdAt = uint48(block.timestamp);

        vm.warp(createdAt + 1 hours);

        vm.prank(responder);
        vm.expectRevert(RiskAssessment.IncidentNotActive.selector);
        riskAssessment.emergencyGrant(
            address(summary),
            incidentId,
            locationHash,
            patientCommitment,
            abi.encodePacked("encrypted-key")
        );
    }

    function testEmergencyAccessExpiresAfterDuration() public {
        incidentId = _createActiveIncident();

        vm.prank(responder);
        riskAssessment.emergencyGrant(
            address(summary),
            incidentId,
            locationHash,
            patientCommitment,
            abi.encodePacked("encrypted-key")
        );

        assertTrue(summary.hasAccess(responder));

        vm.warp(block.timestamp + riskAssessment.ACCESS_DURATION() + 1);
        assertFalse(summary.hasAccess(responder));

        vm.prank(responder);
        vm.expectRevert(EmergencySummary.AccessExpired.selector);
        summary.getData();
    }
}