// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

/// @title LogAnchor
/// @notice Anchors off-chain log-ledger chain hashes onto the blockchain so the
///         integrity of the PostgreSQL audit tables is independently verifiable.
/// @dev Each call to anchor() records the ledger's current chain hash together
///      with the block number and timestamp in which it was sealed. Values are
///      opinionated: only the delegated operator may anchor.
contract LogAnchor {
    struct Anchor {
        uint256 chainHash;
        uint256 blockNumber;
        uint256 blockTimestamp;
    }

    event Anchored(uint256 indexed batchId, uint256 chainHash, uint256 blockNumber);

    error OnlyOperator();

    address public immutable operator;

    /// Batch (anchor) sequence -> recorded digest.
    mapping(uint256 => Anchor) public anchors;

    /// Total number of anchors recorded so far; also the batchId for the next anchor.
    uint256 public batchCount;

    constructor() {
        operator = msg.sender;
    }

    modifier onlyOperator() {
        if (msg.sender != operator) revert OnlyOperator();
        _;
    }

    /// Record a new chain hash and seal it with the current block context.
    function anchor(uint256 chainHash) external onlyOperator returns (uint256 batchId) {
        return _anchor(batchCount, chainHash);
    }

    /// Backfill a missed/malformed anchor for an exact batchId.
    function anchorAt(uint256 batchId, uint256 chainHash) external onlyOperator returns (uint256) {
        if (anchors[batchId].blockNumber != 0) revert AnchorExists(batchId);
        return _anchor(batchId, chainHash);
    }

    function _anchor(uint256 batchId, uint256 chainHash) internal returns (uint256) {
        anchors[batchId] = Anchor({
            chainHash: chainHash,
            blockNumber: block.number,
            blockTimestamp: block.timestamp
        });
        if (batchId >= batchCount) {
            batchCount = batchId + 1;
        }
        emit Anchored(batchId, chainHash, block.number);
        return batchId;
    }
}

error AnchorExists(uint256 batchId);