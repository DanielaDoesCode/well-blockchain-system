// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.20;

import "./Groth16Verifier.sol";

contract ZKAuthVerifier {
    Groth16Verifier public verifier;
    
    uint256 public sessionDuration;
    
    mapping(bytes32 => bool) public usedNullifiers;
    
    mapping(bytes32 => Session) public sessions;
    
    struct Session {
        address user;
        uint48 expiresAt;
        bytes32 commitment;
    }
    
    event Authenticated(
        address indexed user,
        bytes32 indexed commitment,
        bytes32 indexed sessionToken,
        uint48 expiresAt
    );
    
    event SessionRevoked(
        bytes32 indexed sessionToken,
        address indexed user
    );
    
    error InvalidProof();
    error NullifierAlreadyUsed(bytes32 nullifier);
    error SessionExpired();
    error InvalidSession();
    error Unauthorized();
    
    constructor(address _verifier, uint256 _sessionDuration) {
        verifier = Groth16Verifier(_verifier);
        sessionDuration = _sessionDuration;
    }
    
    function authenticate(
        address user,
        bytes32 commitment,
        uint256[2] calldata pA,
        uint256[2][2] calldata pB,
        uint256[2] calldata pC,
        bytes32 nullifier,
        bytes32 actionId
    ) external returns (bytes32 sessionToken) {
        if (usedNullifiers[nullifier]) {
            revert NullifierAlreadyUsed(nullifier);
        }
        
        
        uint256[3] memory pubSignals;
        pubSignals[0] = uint256(commitment);
        pubSignals[1] = uint256(nullifier);
        pubSignals[2] = uint256(actionId);
        
        bool isValid = verifier.verifyProof(pA, pB, pC, pubSignals);
        
        if (!isValid) {
            revert InvalidProof();
        }
        
        usedNullifiers[nullifier] = true;
        
        sessionToken = keccak256(abi.encodePacked(
            user,
            commitment,
            nullifier,
            actionId,
            block.timestamp,
            block.prevrandao
        ));
        
        uint48 expiresAt = uint48(block.timestamp + sessionDuration);
        
        sessions[sessionToken] = Session({
            user: user,
            expiresAt: expiresAt,
            commitment: commitment
        });
        
        emit Authenticated(user, commitment, sessionToken, expiresAt);
    }
    
    function verifySession(bytes32 sessionToken) external view returns (bool) {
        Session memory session = sessions[sessionToken];
        
        if (session.user == address(0)) {
            return false;
        }
        
        if (block.timestamp > session.expiresAt) {
            return false;
        }
        
        return true;
    }
    
    function getSession(bytes32 sessionToken) external view returns (
        address user,
        uint48 expiresAt,
        bytes32 commitment
    ) {
        Session memory session = sessions[sessionToken];
        
        if (session.user == address(0)) {
            revert InvalidSession();
        }
        
        return (session.user, session.expiresAt, session.commitment);
    }
    
    function revokeSession(bytes32 sessionToken) external {
        Session memory session = sessions[sessionToken];
        
        if (session.user == address(0)) {
            revert InvalidSession();
        }
        
        if (session.user != msg.sender) {
            revert Unauthorized();
        }
        
        delete sessions[sessionToken];
        
        emit SessionRevoked(sessionToken, msg.sender);
    }
    
    function revokeSessionByCommitter(bytes32 sessionToken, address committer) external {
        Session memory session = sessions[sessionToken];
        
        if (session.user != committer) {
            revert Unauthorized();
        }
        
        delete sessions[sessionToken];
        
        emit SessionRevoked(sessionToken, committer);
    }
}
