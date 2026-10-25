// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.20;

import "./Types.sol";

contract UserRegistry {

    struct User {
        bytes32      commitment;
        bytes32      internalID;
        Types.Roles  role;
        address      userAddress;
    }

    mapping(bytes32 => bool)  public commitments;
    mapping(bytes32 => User)  public commitmentToUser;
    mapping(bytes32 => bytes32) public internalIDToCommitment;
    mapping(address => bytes32) public addressToCommitment;

    error AlreadyRegistered();
    error NotRegistered();
    error InternalIDAlreadyUsed();
    error AddressAlreadyRegistered();

    event Registered(
        bytes32 indexed commitment,
        bytes32 indexed internalID,
        Types.Roles     role
    );

    event Revoked(bytes32 indexed commitment, bytes32 indexed internalID);

    function register(
        address     userAddress,
        bytes32     commitment,
        bytes32     internalID,
        Types.Roles role
    ) external {
        if (commitments[commitment]) {
            revert AlreadyRegistered();
        }

        if (addressToCommitment[userAddress] != bytes32(0)) {
            revert AddressAlreadyRegistered();
        }

        if (internalIDToCommitment[internalID] != bytes32(0)) {
            revert InternalIDAlreadyUsed();
        }

        commitments[commitment] = true;
        addressToCommitment[userAddress] = commitment;

        commitmentToUser[commitment] = User({
            commitment: commitment,
            internalID: internalID,
            role:       role,
            userAddress: userAddress
        });

        internalIDToCommitment[internalID] = commitment;
        emit Registered(commitment, internalID, role);
    }

    function revoke(bytes32 commitment) external {
        if (!commitments[commitment]) {
            revert NotRegistered();
        }

        bytes32 internalID = commitmentToUser[commitment].internalID;
        address userAddr = commitmentToUser[commitment].userAddress;

        delete internalIDToCommitment[internalID];
        delete commitmentToUser[commitment];
        delete addressToCommitment[userAddr];
        commitments[commitment] = false;

        emit Revoked(commitment, internalID);
    }

    function isRegistered(bytes32 commitment) external view returns (bool) {
        return commitments[commitment];
    }

    function getCommitmentByInternalID(bytes32 internalID)
        external view returns (bytes32 commitment)
    {
        commitment = internalIDToCommitment[internalID];
        if (commitment == bytes32(0)) revert NotRegistered();
    }

    function getCommitmentByAddress(address userAddress)
        external view returns (bytes32 commitment)
    {
        commitment = addressToCommitment[userAddress];
        if (commitment == bytes32(0)) revert NotRegistered();
    }

    function getUser(bytes32 commitment)
        external view returns (User memory)
    {
        if (!commitments[commitment]) revert NotRegistered();
        return commitmentToUser[commitment];
    }

    function getUserByAddress(address userAddress)
        external view returns (User memory)
    {
        bytes32 commitment = addressToCommitment[userAddress];
        if (commitment == bytes32(0)) revert NotRegistered();
        return commitmentToUser[commitment];
    }

    function getRole(bytes32 commitment)
        external view returns (Types.Roles)
    {
        if (!commitments[commitment]) revert NotRegistered();
        return commitmentToUser[commitment].role;
    }

    function getRoleByAddress(address userAddress)
        external view returns (Types.Roles)
    {
        bytes32 commitment = addressToCommitment[userAddress];
        if (commitment == bytes32(0)) revert NotRegistered();
        return commitmentToUser[commitment].role;
    }

    function isDoctor(address userAddress)
        external view returns (bool)
    {
        bytes32 commitment = addressToCommitment[userAddress];
        if(commitment == bytes32(0)) return false;
        return commitmentToUser[commitment].role == Types.Roles.Doctor;
    }

    function isFirstResponder(address userAddress)
        external view returns (bool)
    {
        bytes32 commitment = addressToCommitment[userAddress];
        if (commitment == bytes32(0)) return false;
        return commitmentToUser[commitment].role == Types.Roles.FirstResponder;
    }

    function isHospital(address userAddress)
        external view returns (bool)
    {
        bytes32 commitment = addressToCommitment[userAddress];
        if (commitment == bytes32(0)) return false;
        return commitmentToUser[commitment].role == Types.Roles.Hospital;
    }
}
