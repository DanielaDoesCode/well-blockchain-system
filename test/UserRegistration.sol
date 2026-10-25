pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Registry.sol";
import "../src/UserRegistry.sol";
import "../src/HistoryFactory.sol";
import "../src/Types.sol";
import "../src/Groth16Verifier.sol";
import "../src/ZKAuthVerifier.sol";

contract RegistryContract is Test {
    WellRegistry register;
    UserRegistry userRegistry;
    HistoryFactory historyFactory;
    ZKAuthVerifier zkVerifier;
    address patient = address(0x1);
    bytes32 patientCommitment = keccak256("patient-commitment");
    bytes32 patientInternalID = keccak256("patient-internal-id");

    address doctor = address(0x2);
    bytes32 doctorCommitment = keccak256("doctor-commitment");
    bytes32 doctorInternalID = keccak256("doctor-internal-id");

    function setUp() public {
        Groth16Verifier groth16Verifier = new Groth16Verifier();
        zkVerifier = new ZKAuthVerifier(address(groth16Verifier), 1 hours);

        userRegistry = new UserRegistry();
        historyFactory = new HistoryFactory(address(userRegistry));
        register = new WellRegistry(address(userRegistry), address(historyFactory));
        historyFactory.setWellRegistry(address(register));
    }


    function testRegisterPatient() public {
        vm.prank(patient);
        address history = register.register(
            patientCommitment,
            patientInternalID,
            Types.Roles.Patient
        );

        assertTrue(register.isRegistered(patientCommitment));
        assertTrue(history != address(0));
        assertEq(register.getHistoryContract(patientInternalID), history);
        assertEq(register.getCommitmentByAddress(patient), patientCommitment);
    }

    function testCannotRegisterTwice() public {
        vm.prank(patient);
        register.register(patientCommitment, patientInternalID, Types.Roles.Patient);

        vm.prank(patient);
        vm.expectRevert(UserRegistry.AlreadyRegistered.selector);
        register.register(patientCommitment, patientInternalID, Types.Roles.Patient);
    }

    function testCannotRegisterSameAddressTwice() public {
        vm.prank(patient);
        register.register(patientCommitment, patientInternalID, Types.Roles.Patient);

        bytes32 newCommitment = keccak256("new-commitment");
        bytes32 newInternalID = keccak256("new-internal-id");
        
        vm.prank(patient);
        vm.expectRevert(UserRegistry.AddressAlreadyRegistered.selector);
        register.register(newCommitment, newInternalID, Types.Roles.Patient);
    }
}
