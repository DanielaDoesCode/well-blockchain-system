pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/EmergencySummary.sol";

contract EmergencyRecoveryKeyTest is Test {
    event RecoveryKeyStored(uint256 indexed version, uint256 timestamp);

    EmergencySummary summary;
    address owner = address(0x1111);
    address attacker = address(0x2222);

    bytes wrappedKey = abi.encodePacked("wrapped-content-key");

    function setUp() public {
        // registry/riskAssessment are unused for the recovery-key surface.
        summary = new EmergencySummary(owner, address(0), address(0));
    }

    function testOwnerStoresRecoveryKey() public {
        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit RecoveryKeyStored(1, block.timestamp);
        summary.storeRecoveryKey(wrappedKey);

        assertEq(summary.recoveryVersion(), 1);
        vm.prank(owner);
        assertEq(summary.getRecoveryKey(), wrappedKey);
    }

    function testRecoveryKeyVersionIncrementsOnRotation() public {
        vm.prank(owner);
        summary.storeRecoveryKey(wrappedKey);

        bytes memory rotated = abi.encodePacked("rotated-key");
        vm.prank(owner);
        summary.storeRecoveryKey(rotated);

        assertEq(summary.recoveryVersion(), 2);
        vm.prank(owner);
        assertEq(summary.getRecoveryKey(), rotated);
    }

    function testNonOwnerCannotStoreRecoveryKey() public {
        vm.prank(attacker);
        vm.expectRevert(EmergencySummary.Unauthorized.selector);
        summary.storeRecoveryKey(wrappedKey);
    }

    function testNonOwnerCannotReadRecoveryKey() public {
        vm.prank(owner);
        summary.storeRecoveryKey(wrappedKey);

        vm.prank(attacker);
        vm.expectRevert(EmergencySummary.Unauthorized.selector);
        summary.getRecoveryKey();
    }

    function testRecoveryKeyIsIndependentOfGrantedUsers() public {
        // Store a recovery key, then verify other access paths are unaffected.
        vm.prank(owner);
        summary.storeRecoveryKey(wrappedKey);

        // A responder can request and be granted access without touching recoveryKey.
        vm.prank(attacker);
        uint256 requestId = summary.requestAccess(
            abi.encodePacked("pubkey"), "loc", "desc", "pat-1"
        );

        vm.prank(owner);
        summary.approveAccessRequest(requestId, abi.encodePacked("granted-key"));

        assertTrue(summary.hasAccess(attacker));
        assertEq(summary.recoveryVersion(), 1);
        vm.prank(owner);
        assertEq(summary.getRecoveryKey(), wrappedKey);
    }
}