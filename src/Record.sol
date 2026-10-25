// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.8.0 <0.9.0;

import "./Types.sol";
import "./RecordPermissions.sol";
import "./PolicyVerifier.sol";
import "./lib/PoseidonT3.sol";

contract Record {

    /**
     * ---------------------------------------------------------
     * Assumptions:
     * Records, once created, cannot have their core data altered.
     * Only the integrity hash can be updated (e.g. after key rotation).
     * Access policy constraints (department/institution commitments)
     * can be updated by the owner at any time.
     * ---------------------------------------------------------
     */

    address public immutable owner;    // patient the data relates to
    address public immutable uploader; // doctor who created the record
    address public immutable history;  // parent History contract
    PolicyVerifier public immutable policyVerifier;

    Types.RecordMetadata public metadata;

    // Integrity
    bytes32 public contentHash;
    string  public storageReference; // IPFS CID of the encrypted EHR

    // ZK-obfuscated policy. Only Poseidon commitments of the required
    // department/institution UUIDs are stored — never the plaintext values.
    // bytes32(0) means unconstrained (wildcard).
    bytes32 public requiredDeptCommitment;
    bytes32 public requiredInstCommitment;

    // Freshness window for policy proofs (seconds).
    uint256 public constant POLICY_PROOF_TTL = 300;

    // Permissions contract — one per record
    RecordPermissions public permissions;

    // -----------------------------------------------------------
    //  Errors
    // -----------------------------------------------------------
    error Unauthorized();
    error InvalidProof();
    error StaleProof();
    error PolicyMismatch();

    // -----------------------------------------------------------
    //  Events
    //  These are the canonical source for access events.
    //  The off-chain backend listens to AccessGranted as the
    //  PRE trigger — do not emit duplicate events in History.
    // -----------------------------------------------------------
    event AccessGranted(
        address indexed doctor,
        address indexed grantedBy,
        uint48          grantedAt
    );

    event AccessRevoked(
        address indexed doctor,
        uint48          revokedAt
    );

    event PolicyUpdated(
        bytes32 requiredDeptCommitment,
        bytes32 requiredInstCommitment
    );

    event IntegrityHashUpdated(bytes32 newHash);

    // -----------------------------------------------------------
    //  Modifiers
    // -----------------------------------------------------------
    modifier onlyOwner() {
        if (msg.sender != owner) revert Unauthorized();
        _;
    }

    // History contract acts on the patient's behalf for grant/revoke
    modifier onlyOwnerOrHistory() {
        if (msg.sender != owner && msg.sender != history)
            revert Unauthorized();
        _;
    }

    modifier onlyAuthorized() {
        if (!permissions.hasAccess(msg.sender)) revert Unauthorized();
        _;
    }

    // -----------------------------------------------------------
    //  Constructor
    // -----------------------------------------------------------
    constructor(
        address                    _owner,
        address                    _uploader,
        Types.RecordMetadata memory _metadata,
        bytes32                    _contentHash,
        string memory              _storageReference,
        address                    _policyVerifier
    ) {
        owner            = _owner;
        uploader         = _uploader;
        history          = msg.sender; // History deploys Record, so msg.sender is History
        metadata         = _metadata;
        metadata.createdAt = uint48(block.timestamp);
        contentHash      = _contentHash;
        storageReference = _storageReference;
        policyVerifier   = PolicyVerifier(_policyVerifier);

        // Deploy the permissions contract.
        // Owner and History can both manage access.
        permissions = new RecordPermissions(_owner, history);
    }

    // -----------------------------------------------------------
    //  Access control
    //  Callable by owner directly or via History (onlyOwnerOrHistory).
    //  Emitting here makes Record the canonical event source.
    // -----------------------------------------------------------

    /// Grant a doctor access to this record.
    /// Emits AccessGranted — the off-chain backend treats this as
    /// the signal to trigger proxy re-encryption in the PRE network.
    function grantAccess(address doctor) external onlyOwnerOrHistory {
        permissions.grantAccess(doctor, msg.sender);
        emit AccessGranted(doctor, msg.sender, uint48(block.timestamp));
    }

    /// Policy-driven self-grant. The doctor proves — in zero-knowledge —
    /// that their hidden attributes satisfy the record's committed policy.
    ///
    /// Public signals (verifier order):
    ///   [deptCommitment, instCommitment, doctorCommitment, evaluatedAt]
    ///
    /// Checks (in order):
    ///   1. doctorCommitment == PoseidonT3.hash([msg.sender, 3])  — binds
    ///      the proof to the caller (anti-replay).
    ///   2. commitments match the stored policy (bytes32(0) = wildcard).
    ///   3. evaluatedAt is fresh (within POLICY_PROOF_TTL).
    ///   4. Groth16 proof verifies.
    function grantAccessWithProof(
        uint256[2]   calldata pA,
        uint256[2][2] calldata pB,
        uint256[2]   calldata pC,
        bytes32      deptCommitment,
        bytes32      instCommitment,
        bytes32      doctorCommitment,
        uint48       evaluatedAt
    ) external {
        // 1) doctor binding
        uint256 expected = PoseidonT3.hash([uint256(uint160(msg.sender)), 3]);
        if (doctorCommitment != bytes32(expected)) revert PolicyMismatch();

        // 2) policy match
        if (requiredDeptCommitment != bytes32(0) &&
            deptCommitment != requiredDeptCommitment) revert PolicyMismatch();
        if (requiredInstCommitment != bytes32(0) &&
            instCommitment != requiredInstCommitment) revert PolicyMismatch();

        // 3) freshness
        if (block.timestamp < evaluatedAt ||
            block.timestamp - evaluatedAt > POLICY_PROOF_TTL) revert StaleProof();

        // 4) proof verification
        if (!policyVerifier.verifyProof(
            pA, pB, pC,
            [uint256(deptCommitment), uint256(instCommitment),
             uint256(doctorCommitment), uint256(evaluatedAt)]
        )) revert InvalidProof();

        // 5) self-grant
        permissions.grantAccessByPolicy(msg.sender);
        emit AccessGranted(msg.sender, msg.sender, uint48(block.timestamp));
    }

    /// Revoke a doctor's access to this record.
    function revokeAccess(address doctor) external onlyOwnerOrHistory {
        permissions.revokeAccess(doctor, msg.sender);
        emit AccessRevoked(doctor, uint48(block.timestamp));
    }

    /// Check whether an address currently has access.
    /// Called by the backend as the single on-chain gate before
    /// triggering PRE — no attribute logic lives here.
    function isAuthorized(address doctor) external view returns (bool) {
        return permissions.hasAccess(doctor);
    }

    // -----------------------------------------------------------
    //  ZK-obfuscated static policy constraints (optional, owner-only)
    //  Only Poseidon commitments are stored; the plaintext department/
    //  institution UUIDs never touch the chain. The backend gates on
    //  isAuthorized() after a doctor has self-granted via
    //  grantAccessWithProof() — keeping the hot read path gas-cheap.
    //  Pass bytes32(0) to leave a field unconstrained (wildcard).
    // -----------------------------------------------------------
    function setPolicy(
        bytes32 _requiredDeptCommitment,
        bytes32 _requiredInstCommitment
    ) external onlyOwner {
        requiredDeptCommitment  = _requiredDeptCommitment;
        requiredInstCommitment  = _requiredInstCommitment;
        emit PolicyUpdated(_requiredDeptCommitment, _requiredInstCommitment);
    }

    function getPolicy()
        external view
        returns (bytes32 deptCommitment, bytes32 instCommitment)
    {
        return (requiredDeptCommitment, requiredInstCommitment);
    }

    // -----------------------------------------------------------
    //  Integrity
    //  Updated after key rotation or EHR amendment.
    //  Only the uploader (doctor of record) or owner can do this.
    // -----------------------------------------------------------
    function updateIntegrityHash(bytes32 newHash)
        external
    {
        if (msg.sender != owner && msg.sender != uploader)
            revert Unauthorized();
        contentHash = newHash;
        emit IntegrityHashUpdated(newHash);
    }

    function verifyIntegrity(bytes32 hash) external view returns (bool) {
        return contentHash == hash;
    }

    // -----------------------------------------------------------
    //  Reads
    // -----------------------------------------------------------

    /// Returns the full record bundle.
    /// Caller must be the owner, uploader, or an authorised doctor.
    function getRecord()
        external
        view
        onlyAuthorized
        returns (
            Types.RecordMetadata memory,
            bytes32  _contentHash,
            string memory _storageReference,
            address  _owner,
            address  _uploader
        )
    {
        return (metadata, contentHash, storageReference, owner, uploader);
    }
}