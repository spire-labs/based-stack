// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

/// @title EightBitPacker
/// @notice A custom type for packing multiple values into a single byte to optimize storage use

type EightBitPacker is uint8;

library BitPackerLib {
    // Pack a single value into the byte at the specified bit position, modifying storage directly
    function pack(
        EightBitPacker packedValue,
        uint8 value,
        uint8 bitLength,
        uint8 shiftAmount,
        uint8 scalar
    )
        internal
        returns (EightBitPacker)
    {
        uint8 scaledValue = (value + (scalar / 2)) / scalar;

        require(scaledValue < (1 << bitLength), "scaledValue exceeds bit length");

        // Unwrap the existing packedValue to get the raw uint8 data
        uint8 packed = EightBitPacker.unwrap(packedValue);

        // Clear the bits where the new value will be placed
        uint8 mask = uint8(((1 << bitLength) - 1) << shiftAmount);
        packed = (packed & ~mask) | (scaledValue << shiftAmount); // Pack the scaled value into the correct position

        // Return the new packed value wrapped as EightBitPacker
        return EightBitPacker.wrap(packed);
    }

    // Unpack a single value from the byte at the specified bit position, reading from storage
    function unpack(
        EightBitPacker packedValue,
        uint8 bitLength,
        uint8 shiftAmount,
        uint8 scalar
    )
        internal
        view
        returns (uint8)
    {
        uint8 packed = EightBitPacker.unwrap(packedValue); // Unwrap the packed storage value

        // Extract the value at the shiftAmount position using bitmasking and shifting
        uint8 extracted = uint8((packed >> shiftAmount) & ((1 << bitLength) - 1));

        // Multiply by scalar to return the original value
        return extracted * scalar;
    }

    function calculateScalar(EightBitPacker, uint8 value, uint8 bitLength) internal pure returns (uint8) {
            uint8 maxStorableValue = uint8((1 << bitLength) - 1);  // 2^bitLength - 1 (e.g., for 4 bits, max value is 15)

                // Start with the smallest scalar (1) and increment until a valid one is found
    for (uint8 scalar = 1; scalar <= value; scalar++) {
        if (value % scalar == 0 && value / scalar <= maxStorableValue) {
            return scalar;
        }
    }

    // If no valid scalar is found, revert
    revert("No valid scalar found");
        }
}
