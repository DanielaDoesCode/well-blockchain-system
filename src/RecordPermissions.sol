// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.8.0 <0.9.0;

contract RecordPermissions {

    address public immutable owner;
    address public immutable history; // History contract that can manage permissions
    address public immutable policyRecord; // the Record contract that deployed this
                                            // (only source of policy-based grants)

    // Mapping of addresses that have been granted access
    mapping(address => bool) private accessGranted;

    error Unauthorized();

    modifier onlyAuthorizedManager(address _caller) {
        if (_caller != owner && _caller != history) revert Unauthorized();
        _;
    }

    constructor(address _owner, address _history) {
        owner = _owner;
        history = _history;
        policyRecord = msg.sender;
    }

    function grantAccess(address _user, address _caller) external onlyAuthorizedManager(_caller) {
        accessGranted[_user] = true;
    }

    /// Policy-driven self-grant. Only the owning Record may call this,
    /// after it has verified a valid ZK policy proof.
    function grantAccessByPolicy(address _user) external {
        if (msg.sender != policyRecord) revert Unauthorized();
        accessGranted[_user] = true;
    }

    function revokeAccess(address _user, address _caller) external onlyAuthorizedManager(_caller) {
        accessGranted[_user] = false;
    }

    function hasAccess(address _user) external view returns (bool) {
        // Owner always has access
        return _user == owner || accessGranted[_user];
    }
}