// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

contract MockV3Aggregator {
    uint8 public decimals;
    uint80 private _roundId;
    int256 private _answer;
    uint256 private _updatedAt;

    constructor(uint8 decimals_, int256 answer_) {
        decimals = decimals_;
        _roundId = 1;
        _answer = answer_;
        _updatedAt = block.timestamp;
    }

    function updateAnswer(int256 answer_) external {
        _roundId++;
        _answer = answer_;
        _updatedAt = block.timestamp;
    }

    function updateAnswerWithTimestamp(int256 answer_, uint256 updatedAt_) external {
        _roundId++;
        _answer = answer_;
        _updatedAt = updatedAt_;
    }

    function latestRoundData()
        external
        view
        returns (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
    {
        return (_roundId, _answer, _updatedAt, _updatedAt, _roundId);
    }
}
