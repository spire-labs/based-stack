// SPDX-License-Identifier: MIT
pragma solidity 0.8.25;

import { SSZ } from "src/libraries/SSZ.sol";

contract ValidatorRegistry {
    /// @notice Error emitted when a beacon block root is not found
    error BeaconRootNotFound();

    error InvalidValidatorProof();

    error NoWithdrawalAddress();

    event ValidatorRegistered(uint256 indexed validatorIndex, address indexed delegatedAddress);

    address internal constant BEACON_ROOTS_CONTRACT = 0x000F3df6D732807Ef1319fB7B8bB8522d0Beac02;

    uint256 internal constant ETH2_GENESIS_TIMESTAMP = 1606824023;

    uint256 internal constant GENESIS_INDEX = 798245441765376;

    // validatorIndex -> address that is delegated to for the validator
    mapping(uint256 => address) public validatorIndexToDelegatedAddress;

    // validatorIndex -> validator
    mapping(uint256 => SSZ.Validator) public validatorIndexToValidator;

    function registerValidator(
        SSZ.Validator calldata _validator,
        uint256 _validatorIndex,
        bytes32[] calldata _validatorInclusionProof
    )
        external
    {
        if (bytes1(_validator.withdrawalCredentials) == bytes1(0x00)) revert NoWithdrawalAddress();

        bytes32 _validatorTreeRoot = SSZ._hashValidatorTreeRoot(_validator);
        uint256 _gi = GENESIS_INDEX + _validatorIndex;
        bytes32 _blockRoot = _getBeaconParentBlockRoot();

        bytes32 _calculatedBlockRoot = SSZ._restoreMerkleRoot(_validatorTreeRoot, _gi, _validatorInclusionProof);

        if (_calculatedBlockRoot != _blockRoot) revert InvalidValidatorProof();

        validatorIndexToValidator[_validatorIndex] = _validator;

        address _withdrawalAddress = _extractWithdrawalCredentials(_validator.withdrawalCredentials);

        validatorIndexToDelegatedAddress[_validatorIndex] = _withdrawalAddress;

        emit ValidatorRegistered(_validatorIndex, _withdrawalAddress);
    }

    function _extractWithdrawalCredentials(bytes32 _withdrawalCredentials)
        internal
        pure
        returns (address withdrawalAddress_)
    {
        withdrawalAddress_ = address(uint160(uint256(_withdrawalCredentials)));
    }

    function _getBeaconParentBlockRoot() internal view returns (bytes32 parentBlockRoot_) {
        (bool _success, bytes memory _data) = BEACON_ROOTS_CONTRACT.staticcall(abi.encode(block.timestamp));

        if (!_success || _data.length == 0) {
            revert BeaconRootNotFound();
        }

        parentBlockRoot_ = abi.decode(_data, (bytes32));
    }
}
