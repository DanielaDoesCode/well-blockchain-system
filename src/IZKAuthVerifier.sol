// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.20;

interface IZKAuthVerifier {
    function authenticate(
        address user,
        bytes32 commitment,
        uint256[2] calldata pA,
        uint256[2][2] calldata pB,
        uint256[2] calldata pC,
        bytes32 nullifier,
        bytes32 actionId
    ) external returns (bytes32 sessionToken);
    
    function verifySession(bytes32 sessionToken) external view returns (bool);
    
    function getSession(bytes32 sessionToken) external view returns (
        address user,
        uint48 expiresAt,
        bytes32 commitment
    );
    
    function revokeSession(bytes32 sessionToken) external;
    
    function usedNullifiers(bytes32 nullifier) external view returns (bool);
}
