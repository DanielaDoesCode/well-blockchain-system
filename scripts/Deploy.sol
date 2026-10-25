// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.8.0 <0.9.0;

import "forge-std/Script.sol";
import "./Groth16Verifier.sol";
import "./ZKAuthVerifier.sol";
import "./UserRegistry.sol";
import "./HistoryFactory.sol";
import "./Registry.sol";
import "./IncidentManager.sol";
import "./EmergencySummaryFactory.sol";
import "./RiskAssessment.sol";
import "./AttestationManager.sol";
import "./ResearchVerifier.sol";
import "./ResearchContributer.sol";
import "./PolicyVerifier.sol";

contract Deploy is Script {

    function run() external {
        vm.startBroadcast();

        // -- Auth layer -------------------------------------------------------
        Groth16Verifier groth16Verifier = new Groth16Verifier();
        ZKAuthVerifier zkVerifier = new ZKAuthVerifier(
            address(groth16Verifier),
            1 hours
        );

        // -- Identity / registry layer ----------------------------------------
        UserRegistry userRegistry = new UserRegistry();
        PolicyVerifier policyVerifier = new PolicyVerifier();
        HistoryFactory historyFactory = new HistoryFactory(address(userRegistry));
        WellRegistry wellRegistry = new WellRegistry(
            address(userRegistry),
            address(historyFactory)
        );
        historyFactory.setWellRegistry(address(wellRegistry));
        historyFactory.setPolicyVerifier(address(policyVerifier));

        // -- Emergency infrastructure -----------------------------------------
        IncidentManager incidentManager = new IncidentManager(
            address(userRegistry)
        );
        EmergencySummaryFactory emergencySummaryFactory = new EmergencySummaryFactory(
            address(userRegistry)
        );
        RiskAssessment riskAssessment = new RiskAssessment(
            address(incidentManager),
            address(userRegistry),
            address(emergencySummaryFactory)
        );

        // Wire emergency contracts into WellRegistry
        wellRegistry.setIncidentManager(address(incidentManager));
        wellRegistry.setEmergencySummaryFactory(address(emergencySummaryFactory));
        wellRegistry.setRiskAssessment(address(riskAssessment));

        // Wire EmergencySummaryFactory to WellRegistry
        emergencySummaryFactory.setWellRegistry(address(wellRegistry));
        emergencySummaryFactory.setRiskAssessment(address(riskAssessment));

        // -- Attestation layer ------------------------------------------------
        // AttestationManager depends on UserRegistry to verify isDoctor().
        // Must be deployed before ResearchContributer.
        AttestationManager attestationManager = new AttestationManager(
            address(userRegistry)
        );

        // -- Research layer ---------------------------------------------------
        ResearchVerifier researchVerifier = new ResearchVerifier();
        ResearchContributer researchContributer = new ResearchContributer(
            address(researchVerifier),
            address(wellRegistry),
            address(attestationManager)
        );

        // -- Log addresses ----------------------------------------------------
        console.log("Groth16Verifier deployed at:    ", address(groth16Verifier));
        console.log("ZKAuthVerifier deployed at:     ", address(zkVerifier));
        console.log("UserRegistry deployed at:       ", address(userRegistry));
        console.log("PolicyVerifier deployed at:     ", address(policyVerifier));
        console.log("HistoryFactory deployed at:     ", address(historyFactory));
        console.log("WellRegistry deployed at:       ", address(wellRegistry));
        console.log("IncidentManager deployed at:    ", address(incidentManager));
        console.log("EmergencySummaryFactory at:     ", address(emergencySummaryFactory));
        console.log("RiskAssessment deployed at:     ", address(riskAssessment));
        console.log("AttestationManager deployed at: ", address(attestationManager));
        console.log("ResearchVerifier deployed at:   ", address(researchVerifier));
        console.log("ResearchContributer deployed at:", address(researchContributer));

        vm.stopBroadcast();
    }
}
