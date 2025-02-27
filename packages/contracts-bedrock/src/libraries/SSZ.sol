// SPDX-License-Identifier: MIT
pragma solidity 0.8.25;

library SSZ {
    error InvalidPubkey();

    // As defined in phase0/beacon-chain.md:356
    struct Validator {
        bytes pubkey;
        bytes32 withdrawalCredentials;
        uint64 effectiveBalance;
        bool slashed;
        uint64 activationEligibilityEpoch;
        uint64 activationEpoch;
        uint64 exitEpoch;
        uint64 withdrawableEpoch;
    }

    function _restoreMerkleRoot(
        bytes32 _leaf,
        uint256 _gi,
        bytes32[] memory _branch
    )
        internal
        pure
        returns (bytes32 root_)
    {
        require(2 ** (_branch.length + 1) > _gi);
        root_ = _leaf;
        uint256 i = 0;
        while (_gi != 1) {
            if (_gi % 2 == 1) {
                root_ = sha256(bytes.concat(_branch[i], root_));
            } else {
                root_ = sha256(bytes.concat(root_, _branch[i]));
            }
            _gi /= 2;

            unchecked {
                ++i;
            }
        }
    }

    function _hashValidatorTreeRoot(Validator calldata validator) internal pure returns (bytes32 validatorRoot_) {
        bytes32[8] memory _chunks = [
            _hashPubkey(validator.pubkey),
            validator.withdrawalCredentials,
            _toLittleEndian(validator.effectiveBalance),
            _toLittleEndian(validator.slashed),
            _toLittleEndian(validator.activationEligibilityEpoch),
            _toLittleEndian(validator.activationEpoch),
            _toLittleEndian(validator.exitEpoch),
            _toLittleEndian(validator.withdrawableEpoch)
        ];

        validatorRoot_ = _merkleize(_chunks);
    }

    function _merkleize(bytes32[8] memory _chunks) internal pure returns (bytes32 root_) {
        // Level 1
        bytes32[4] memory _level1;
        _level1[0] = _hashNode(_chunks[0], _chunks[1]);
        _level1[1] = _hashNode(_chunks[2], _chunks[3]);
        _level1[2] = _hashNode(_chunks[4], _chunks[5]);
        _level1[3] = _hashNode(_chunks[6], _chunks[7]);

        // Level 2
        bytes32[2] memory _level2;
        _level2[0] = _hashNode(_level1[0], _level1[1]);
        _level2[1] = _hashNode(_level1[2], _level1[3]);

        // Level 3 (root)
        root_ = _hashNode(_level2[0], _level2[1]);
    }

    // Helper function to compute the hash tree root of the pubkey
    function _hashPubkey(bytes memory _pubkey) internal pure returns (bytes32 pubkeyRoot_) {
        if (_pubkey.length != 48) revert InvalidPubkey();

        bytes32 _chunk0;
        bytes32 _chunk1;

        assembly {
            _chunk0 := mload(add(_pubkey, 32)) // First 32 bytes
            _chunk1 := mload(add(_pubkey, 64)) // Next 16 bytes (and 16 bytes of padding)
        }

        // Mask _chunk1 to keep only the first 16 bytes
        _chunk1 = bytes32(uint256(_chunk1) & 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000000000000000000000000000);

        // Compute the hash of the two chunks
        pubkeyRoot_ = _hashNode(_chunk0, _chunk1);
    }

    function _hashNode(bytes32 _left, bytes32 _right) internal pure returns (bytes32 node_) {
        return sha256(bytes.concat(_left, _right));
    }

    /// @notice Converts a value to a little-endian byte array
    /// @dev from https://github.com/succinctlabs/telepathy-contracts/blob/main/src/libraries/SimpleSerialize.sol
    function _toLittleEndian(uint256 _v) internal pure returns (bytes32 result_) {
        _v = ((_v & 0xFF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00) >> 8)
            | ((_v & 0x00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF) << 8);
        _v = ((_v & 0xFFFF0000FFFF0000FFFF0000FFFF0000FFFF0000FFFF0000FFFF0000FFFF0000) >> 16)
            | ((_v & 0x0000FFFF0000FFFF0000FFFF0000FFFF0000FFFF0000FFFF0000FFFF0000FFFF) << 16);
        _v = ((_v & 0xFFFFFFFF00000000FFFFFFFF00000000FFFFFFFF00000000FFFFFFFF00000000) >> 32)
            | ((_v & 0x00000000FFFFFFFF00000000FFFFFFFF00000000FFFFFFFF00000000FFFFFFFF) << 32);
        _v = ((_v & 0xFFFFFFFFFFFFFFFF0000000000000000FFFFFFFFFFFFFFFF0000000000000000) >> 64)
            | ((_v & 0x0000000000000000FFFFFFFFFFFFFFFF0000000000000000FFFFFFFFFFFFFFFF) << 64);
        _v = (_v >> 128) | (_v << 128);

        result_ = bytes32(_v);
    }

    /// @notice Converts a boolean to a little-endian byte array
    function _toLittleEndian(bool _v) internal pure returns (bytes32 result_) {
        result_ = bytes32(_v ? 1 << 248 : 0);
    }
}
