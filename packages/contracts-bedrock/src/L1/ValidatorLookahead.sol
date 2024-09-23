// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

contract  ValidatorLookahead {
    struct ValidatorData {
        bytes pubkey;
        uint256 validatorIndex;
        uint256 slot;
    }

    ValidatorData[32] private validatorLookahead;

    event ValidatorLookaheadUpdated(uint256 epoch);

    function storeValidatorLookahead(
        ValidatorData[32] calldata newValidatorData,
        uint256 epoch
    ) external {
        for (uint256 i = 0; i < 32; i++) {
            validatorLookahead[i] = newValidatorData[i];
        }

        emit ValidatorLookaheadUpdated(epoch);
    }

    function getValidatorLookahead()
        external
        view
        returns (ValidatorData[32] memory)
    {
        return validatorLookahead;
    }
}
