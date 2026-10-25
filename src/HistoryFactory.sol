// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.20;

import "./Types.sol";
import "./History.sol";
import "./UserRegistry.sol";

contract HistoryFactory {

    UserRegistry public userRegistry;
    address public wellRegistry;
    address public policyVerifier;

    mapping(bytes32 => address) public commitmentToHistory;

    error NotRegistered();
    error Unauthorized();

    event HistoryCreated(
        address indexed patient,
        address indexed historyContract,
        bytes32 indexed commitment
    );

    modifier onlyWellRegistry() {
        if (msg.sender != wellRegistry) revert Unauthorized();
        _;
    }

    constructor(address _userRegistry) {
        userRegistry = UserRegistry(_userRegistry);
    }

    function setWellRegistry(address _wellRegistry) external {
        if (wellRegistry != address(0)) revert Unauthorized();
        wellRegistry = _wellRegistry;
    }

    /// One-time setter — must be called before any History is created so
    /// every Record gets the policy verifier at construction time.
    function setPolicyVerifier(address _policyVerifier) external {
        if (policyVerifier != address(0)) revert Unauthorized();
        policyVerifier = _policyVerifier;
    }

    function createHistory(bytes32 commitment, address patientAddress)
        external onlyWellRegistry returns (address)
    {
        History history = new History(patientAddress, wellRegistry, policyVerifier);
        commitmentToHistory[commitment] = address(history);

        emit HistoryCreated(patientAddress, address(history), commitment);
        return address(history);
    }

    function getHistoryContract(bytes32 internalID)
        external view returns (address)
    {
        bytes32 commitment = userRegistry.getCommitmentByInternalID(internalID);
        address h = commitmentToHistory[commitment];
        if (h == address(0)) revert NotRegistered();
        return h;
    }

    function getHistoryContractByAddress(address userAddress)
        external view returns (address)
    {
        return _getHistoryContractByAddress(userAddress);
    }

    function _getHistoryContractByAddress(address userAddress)
        internal view returns (address)
    {
        bytes32 commitment = userRegistry.getCommitmentByAddress(userAddress);
        address h = commitmentToHistory[commitment];
        if (h == address(0)) revert NotRegistered();
        return h;
    }
}
