// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.8.0 <0.9.0;

import "./Types.sol";

contract EmergencySummary {

    address public immutable owner;
    address public immutable registry;
    address public immutable riskAssessement;

    bytes public encryptedData;
    
    bytes32 public dataHash;
    
    uint256 public dataVersion;

    uint256 public lastUpdated;

    // Recovery escrow: the data key, wrapped by the operator's Vault transit
    // machinery, persisted on-chain as a tamper-evident, recoverable copy.
    bytes public recoveryKey;

    uint256 public recoveryVersion;

    struct AccessGrant {
        bool isGranted;
        uint256 grantedAt;
        uint256 expiresAt;
    }
    
    mapping(address => AccessGrant) private accessGrants;
    mapping(address => bytes) private encryptedKeys;
    mapping(address => bytes) private userPublicKeys;
    address[] private authorizedUsers;

    struct AccessRequest {
        address responder;
        bytes publicKey;            // first responder public key
        string location;
        string incidentDescription;
        string patientIdentifier;
        uint256 timestamp;
        AccessStatus status;
    }
    
    enum AccessStatus { Pending, Granted, Denied, Expired }
    
    mapping(uint256 => AccessRequest) private accessRequests;
    uint256 private requestCounter;
    uint256[] private accessLog;

    uint256 public constant ACCESS_DURATION = 2 hours;
    uint256 public constant REQUEST_EXPIRY = 15 minutes;

    error Unauthorized();
    error InvalidRequest();
    error RequestNotFound();
    error RequestExpired();
    error NoAccessKey();
    error AccessExpired();
    error PublicKeyRequired();

    event DataUpdated(uint256 version, bytes32 dataHash, uint256 timestamp);
    event AccessRequested(uint256 indexed requestId, address indexed responder, uint256 timestamp);
    event AccessGranted(address indexed user, uint256 expiresAt);
    event AccessRevoked(address indexed user);
    event PublicKeyRegistered(address indexed user);
    event RecoveryKeyStored(uint256 indexed version, uint256 timestamp);


    modifier onlyOwnerOrRiskAssessement() {
        if (msg.sender != owner && msg.sender != riskAssessement) revert Unauthorized();
        _;
    }

    modifier onlyOwner() {
        if (msg.sender != owner) revert Unauthorized();
        _;
    }

    modifier hasValidAccess() {
        if (msg.sender != owner) {
            AccessGrant storage grant = accessGrants[msg.sender];
            if (!grant.isGranted) revert Unauthorized();
            if (grant.expiresAt != 0 && block.timestamp > grant.expiresAt) {
                revert AccessExpired();
            }
        }
        _;
    }


    constructor(address _owner, address _registry, address _riskAssessement) {
        owner = _owner;
        registry = _registry;
        riskAssessement = _riskAssessement;
    }


    function updateData(
        bytes calldata _encryptedData,
        bytes32 _dataHash
    ) external onlyOwner {
        encryptedData = _encryptedData;
        dataHash = _dataHash;
        dataVersion++;
        lastUpdated = block.timestamp;
        
        emit DataUpdated(dataVersion, _dataHash, block.timestamp);
    }

    function grantAccess(
        address _user,
        bytes calldata _encryptedKey,
        uint256 _duration
    ) external onlyOwnerOrRiskAssessement {
        uint256 expiry = _duration == 0 ? 0 : block.timestamp + _duration;
        
        encryptedKeys[_user] = _encryptedKey;
        accessGrants[_user] = AccessGrant({
            isGranted: true,
            grantedAt: block.timestamp,
            expiresAt: expiry
        });
        
        // Track authorized users for key rotation
        bool exists = false;
        for (uint256 i = 0; i < authorizedUsers.length; i++) {
            if (authorizedUsers[i] == _user) {
                exists = true;
                break;
            }
        }
        if (!exists) {
            authorizedUsers.push(_user);
        }
        
        emit AccessGranted(_user, expiry);
    }


    function updateKeys(
        address[] calldata _users,
        bytes[] calldata _encryptedKeys
    ) external onlyOwner {
        if (_users.length != _encryptedKeys.length) revert InvalidRequest();
        
        for (uint256 i = 0; i < _users.length; i++) {
            if (accessGrants[_users[i]].isGranted) {
                encryptedKeys[_users[i]] = _encryptedKeys[i];
            }
        }
    }

    /**
     * @notice Escrow a wrapped copy of the data key for operator-assisted recovery.
     * @dev The blob is opaque to the contract. Storage is owner-only and increments
     *      recoveryVersion so rotation of the escrow is auditable.
     */
    function storeRecoveryKey(bytes calldata _recoveryKey) external onlyOwner {
        recoveryKey = _recoveryKey;
        recoveryVersion++;

        emit RecoveryKeyStored(recoveryVersion, block.timestamp);
    }

    /**
     * @notice Read the escrowed, wrapped data key (owner only).
     */
    function getRecoveryKey() external view onlyOwner returns (bytes memory) {
        return recoveryKey;
    }


    function registerPublicKey(bytes calldata _publicKey) external {
        userPublicKeys[msg.sender] = _publicKey;
        emit PublicKeyRegistered(msg.sender);
    }

    function requestAccess(
        bytes calldata _publicKey,
        string calldata _location,
        string calldata _incidentDescription,
        string calldata _patientIdentifier
    ) external returns (uint256 requestId) {
        if (_publicKey.length == 0) revert PublicKeyRequired();
        

        userPublicKeys[msg.sender] = _publicKey;
        
        requestId = requestCounter++;
        accessRequests[requestId] = AccessRequest({
            responder: msg.sender,
            publicKey: _publicKey,
            location: _location,
            incidentDescription: _incidentDescription,
            patientIdentifier: _patientIdentifier,
            timestamp: block.timestamp,
            status: AccessStatus.Pending
        });
        
        emit AccessRequested(requestId, msg.sender, block.timestamp);
        return requestId;
    }


    function approveAccessRequest(
        uint256 _requestId,
        bytes calldata _encryptedKey
    ) external onlyOwner {
        AccessRequest storage request = accessRequests[_requestId];
        
        if (request.responder == address(0)) revert RequestNotFound();
        if (request.status != AccessStatus.Pending) revert InvalidRequest();
        if (block.timestamp > request.timestamp + REQUEST_EXPIRY) {
            request.status = AccessStatus.Expired;
            revert RequestExpired();
        }
        
        request.status = AccessStatus.Granted;
        
        uint256 expiry = block.timestamp + ACCESS_DURATION;
        encryptedKeys[request.responder] = _encryptedKey;
        accessGrants[request.responder] = AccessGrant({
            isGranted: true,
            grantedAt: block.timestamp,
            expiresAt: expiry
        });
        
        accessLog.push(_requestId);
        emit AccessGranted(request.responder, expiry);
    }

    /**
     * @notice Deny access request
     */
    function denyAccessRequest(uint256 _requestId) external onlyOwner {
        AccessRequest storage request = accessRequests[_requestId];
        if (request.responder == address(0)) revert RequestNotFound();
        request.status = AccessStatus.Denied;
    }

    // ─────────────────────────────────────────────────────────────
    // Read Functions
    // ─────────────────────────────────────────────────────────────

    /**
     * @notice Get encrypted data and your decryption key
     * @return _encryptedData The encrypted summary
     * @return _encryptedKey Your copy of the AES key (encrypted with your public key)
     * @return _version Data version
     * @return _hash Hash of plaintext for verification
     */
    function getData() external view hasValidAccess returns (
        bytes memory _encryptedData,
        bytes memory _encryptedKey,
        uint256 _version,
        bytes32 _hash
    ) {
        bytes storage key = encryptedKeys[msg.sender];
        if (msg.sender != owner && key.length == 0) revert NoAccessKey();
        
        return (
            encryptedData,
            key,
            dataVersion,
            dataHash
        );
    }

    /**
     * @notice Get just the encrypted data (if you already have the key)
     */
    function getEncryptedData() external view hasValidAccess returns (bytes memory) {
        return encryptedData;
    }

    /**
     * @notice Get your encrypted key
     */
    function getEncryptedKey() external view hasValidAccess returns (bytes memory) {
        return encryptedKeys[msg.sender];
    }

    /**
     * @notice Check access status
     */
    function hasAccess(address _user) external view returns (bool) {
        if (_user == owner) return true;
        AccessGrant storage grant = accessGrants[_user];
        if (!grant.isGranted) return false;
        if (grant.expiresAt != 0 && block.timestamp > grant.expiresAt) return false;
        return true;
    }

    /**
     * @notice Get access details
     */
    function getAccessInfo(address _user) external view returns (
        bool isGranted,
        uint256 grantedAt,
        uint256 expiresAt
    ) {
        AccessGrant storage grant = accessGrants[_user];
        return (grant.isGranted, grant.grantedAt, grant.expiresAt);
    }

    /**
     * @notice Get request details (for responder or owner)
     */
    function getAccessRequest(uint256 _requestId) external view returns (
        address responder,
        string memory location,
        string memory incidentDescription,
        uint256 timestamp,
        AccessStatus status
    ) {
        AccessRequest storage request = accessRequests[_requestId];
        if (msg.sender != request.responder && msg.sender != owner) {
            revert Unauthorized();
        }
        return (
            request.responder,
            request.location,
            request.incidentDescription,
            request.timestamp,
            request.status
        );
    }

    /**
     * @notice Get all authorized users (for key rotation)
     */
    function getAuthorizedUsers() external view onlyOwner returns (address[] memory) {
        return authorizedUsers;
    }

    /**
     * @notice Get access log
     */
    function getAccessLog() external view onlyOwner returns (uint256[] memory) {
        return accessLog;
    }

    /**
     * @notice Get data metadata (public)
     */
    function getMetadata() external view returns (
        uint256 version,
        uint256 updated,
        uint256 dataSize
    ) {
        return (dataVersion, lastUpdated, encryptedData.length);
    }
}