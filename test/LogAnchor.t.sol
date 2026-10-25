// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/LogAnchor.sol";

contract LogAnchorTest is Test {
    LogAnchor internal anchor;
    address internal operator = address(0xA);
    address internal stranger = address(0xB);

    function setUp() public {
        vm.prank(operator);
        anchor = new LogAnchor();
    }

    function test_DefaultState() public view {
        assertEq(anchor.operator(), operator);
        assertEq(anchor.batchCount(), 0);
    }

    function test_AnchorsStampBlockData() public {
        uint256 chainHash = 123456789;
        vm.prank(operator);
        uint256 batchId = anchor.anchor(chainHash);

        // batch ids start at 0
        assertEq(batchId, 0);
        assertEq(anchor.batchCount(), 1);

        (uint256 stored, uint256 blockNumber, uint256 timestamp) = anchor.anchors(0);
        assertEq(stored, chainHash);
        assertEq(blockNumber, block.number, "anchored block number must be recorded");
        assertEq(timestamp, block.timestamp, "anchored timestamp must be recorded");
    }

    function test_AnchorsAreAppendOnlyPerBatch() public {
        vm.prank(operator);
        anchor.anchor(111);
        vm.prank(operator);
        anchor.anchor(222);

        uint256[] memory expected = new uint256[](2);
        expected[0] = 111;
        expected[1] = 222;
        for (uint256 i = 0; i < 2; i++) {
            (uint256 stored,,) = anchor.anchors(i);
            assertEq(stored, expected[i], "each batch must keep its own chain hash");
        }
        assertEq(anchor.batchCount(), 2);
    }

    function test_anchorAtAllowsDeterministicBatchId() public {
        vm.prank(operator);
        anchor.anchorAt(9, 444);

        (uint256 stored,,) = anchor.anchors(9);
        assertEq(stored, 444);
        assertEq(anchor.batchCount(), 10, "gap counts towards batchCount");
    }

    function testRevert_OnlyOperatorAllowed() public {
        vm.prank(stranger);
        vm.expectRevert(LogAnchor.OnlyOperator.selector);
        anchor.anchor(1);
    }

    function testRevert_anchorAtRejectsReusedBatchId() public {
        vm.prank(operator);
        anchor.anchorAt(5, 100);

        vm.prank(operator);
        vm.expectRevert(abi.encodeWithSelector(AnchorExists.selector, uint256(5)));
        anchor.anchorAt(5, 200);
    }

    function testRevert_OperatorFixedAtDeployment() public {
        vm.prank(stranger);
        vm.expectRevert(LogAnchor.OnlyOperator.selector);
        anchor.anchor(0);
    }
}