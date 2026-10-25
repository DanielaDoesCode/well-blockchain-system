// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/ZKAuthVerifier.sol";
import "../src/Groth16Verifier.sol";

contract ZKAuthVerifierTest is Test {
    ZKAuthVerifier public verifier;
    Groth16Verifier public groth16Verifier;
    
    address public user = address(0x1);
    bytes32 public commitment = keccak256("test-commitment");
    bytes32 public nullifier = keccak256("test-nullifier");
    bytes32 public actionId = keccak256("test-action-id");
    
    uint256 public constant SESSION_DURATION = 1 hours;
    
    function setUp() public {
        groth16Verifier = new Groth16Verifier();
        verifier = new ZKAuthVerifier(address(groth16Verifier), SESSION_DURATION);
    }
    
    function testInitialState() public view {
        assertEq(address(verifier.verifier()), address(groth16Verifier));
        assertEq(verifier.sessionDuration(), SESSION_DURATION);
    }
    
    function testNullifierTracking() public view {
        assertFalse(verifier.usedNullifiers(nullifier));
    }
    
    function testSessionVerificationFailsForInvalidToken() public view {
        bytes32 invalidToken = keccak256("invalid-token");
        assertFalse(verifier.verifySession(invalidToken));
    }
    
    function testGetSessionRevertsForInvalidToken() public {
        bytes32 invalidToken = keccak256("invalid-token");
        vm.expectRevert(ZKAuthVerifier.InvalidSession.selector);
        verifier.getSession(invalidToken);
    }
    
    function testAuthenticateFailsWithInvalidProof() public {
        uint256[2] memory pA = [uint256(1), uint256(2)];
        uint256[2][2] memory pB;
        pB[0] = [uint256(1), uint256(2)];
        pB[1] = [uint256(3), uint256(4)];
        uint256[2] memory pC = [uint256(1), uint256(2)];
        
        vm.expectRevert(ZKAuthVerifier.InvalidProof.selector);
        verifier.authenticate(
            user,
            commitment,
            pA,
            pB,
            pC,
            nullifier,
            actionId
        );
    }
    
    function testRevokeSessionByCommitterRevertsForInvalidSession() public {
        bytes32 sessionToken = keccak256("test-session");
        
        vm.expectRevert(ZKAuthVerifier.Unauthorized.selector);
        verifier.revokeSessionByCommitter(sessionToken, user);
    }
    
    function testRevokeSessionRevertsForUnauthorized() public {
        bytes32 sessionToken = keccak256("test-session");
        
        vm.prank(user);
        vm.expectRevert(ZKAuthVerifier.InvalidSession.selector);
        verifier.revokeSession(sessionToken);
    }
    
    function testAuthenticateFailsWithUsedNullifier() public {
        uint256[2] memory pA = [uint256(0), uint256(0)];
        uint256[2][2] memory pB;
        pB[0] = [uint256(0), uint256(0)];
        pB[1] = [uint256(0), uint256(0)];
        uint256[2] memory pC = [uint256(0), uint256(0)];
        
        bytes32 usedNullifier = keccak256("used-nullifier");
        
        vm.mockCall(
            address(groth16Verifier),
            abi.encodeWithSelector(
                Groth16Verifier.verifyProof.selector,
                pA,
                pB,
                pC,
                [uint256(commitment), uint256(usedNullifier), uint256(actionId)]
            ),
            abi.encode(true)
        );
        
        verifier.authenticate(
            user,
            commitment,
            pA,
            pB,
            pC,
            usedNullifier,
            actionId
        );
        
        assertTrue(verifier.usedNullifiers(usedNullifier));
        
        vm.expectRevert(abi.encodeWithSelector(
            ZKAuthVerifier.NullifierAlreadyUsed.selector,
            usedNullifier
        ));
        
        verifier.authenticate(
            user,
            commitment,
            pA,
            pB,
            pC,
            usedNullifier,
            actionId
        );
    }
    
    function testSessionTokenGeneration() public view {
        bytes32 sessionToken = keccak256(abi.encodePacked(
            user,
            commitment,
            nullifier,
            actionId,
            block.timestamp,
            block.prevrandao
        ));
        
        assertTrue(sessionToken != bytes32(0));
    }
    
    function testEventEmission() public {
        uint256[2] memory pA = [uint256(0), uint256(0)];
        uint256[2][2] memory pB;
        pB[0] = [uint256(0), uint256(0)];
        pB[1] = [uint256(0), uint256(0)];
        uint256[2] memory pC = [uint256(0), uint256(0)];
        
        vm.mockCall(
            address(groth16Verifier),
            abi.encodeWithSelector(
                Groth16Verifier.verifyProof.selector,
                pA,
                pB,
                pC,
                [uint256(commitment), uint256(nullifier), uint256(actionId)]
            ),
            abi.encode(true)
        );
        
        bytes32 sessionToken = verifier.authenticate(
            user,
            commitment,
            pA,
            pB,
            pC,
            nullifier,
            actionId
        );
        
        assertTrue(sessionToken != bytes32(0));
        assertEq(verifier.verifySession(sessionToken), true);
    }
}
