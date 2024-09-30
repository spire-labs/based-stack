import { Test } from "forge-std/Test.sol";
import { EightBitPacker, BitPackerLib } from "src/libraries/EightBitPacker.sol";

contract TestEightBitPacker is Test {
    using BitPackerLib for EightBitPacker;

    EightBitPacker public packedValue;
    uint8 constant BIT_LENGTH = 4;
    uint8 constant SHIFT_AMOUNT = 2;

    function testPack() public {
        uint8 value = 0x12;
        uint8 scalar = 2;

        packedValue = packedValue.pack(value, BIT_LENGTH, SHIFT_AMOUNT, scalar);

        assertEq(packedValue.unpack(BIT_LENGTH, SHIFT_AMOUNT, scalar), value);
    }

    function testPackMultipleDifferentValues() public {
        uint8 value1 = 0x12; // 18 in decimal
        uint8 value2 = 0x34; // 52 in decimal

        // Scalars for each value to scale them down to fit in 4 bits
        uint8 scalar1 = 2; // Scale value1 down by 2 (to fit within 4 bits)
        uint8 scalar2 = 4; // Scale value2 down by 4 (to fit within 4 bits)

        // Shift amounts for each value
        uint8 shiftAmount1 = 0; // No shift for value1 (store in the least significant 4 bits)
        uint8 shiftAmount2 = 4; // Shift value2 by 4 bits to store in the next 4 bits

        // Pack the values into the 8-bit storage
        packedValue = packedValue.pack(value1, 4, shiftAmount1, scalar1); // Pack value1 into the first 4 bits
        packedValue = packedValue.pack(value2, 4, shiftAmount2, scalar2); // Pack value2 into the next 4 bits

        // Unpack and check the values, scaling them back to the original values
        assertEq(packedValue.unpack(4, shiftAmount1, scalar1), value1);
        assertEq(packedValue.unpack(4, shiftAmount2, scalar2), value2);
    }

    function testPackMultipleDifferentValuesWithDifferentBitLengths() public {
        uint8 value1 = 4; // Example small value that will be scaled down
        uint8 value2 = 8; // Another small value
        uint8 value3 = 12; // Larger value
        uint8 value4 = 15; // Largest value fitting in 2 bits after scaling

        // Scalars to reduce each value so it fits in 2 bits (values between 0 and 3)
        uint8 scalar1 = 2; // Scaling down value1 by 2
        uint8 scalar2 = 4; // Scaling down value2 by 4
        uint8 scalar3 = 6; // Scaling down value3 by 6
        uint8 scalar4 = 5; // Scaling down value4 by 7

        // Since we are storing 4 values, and 8 bits in total are available, each value gets 2 bits
        uint8 bitLength = 2;

        // Shift amounts for each value in the 8-bit storage
        uint8 shiftAmount1 = 0; // No shift for value1 (store in the least significant bits)
        uint8 shiftAmount2 = 2; // Shift value2 by 2 bits
        uint8 shiftAmount3 = 4; // Shift value3 by 4 bits
        uint8 shiftAmount4 = 6; // Shift value4 by 6 bits

        // Pack the values into the 8-bit storage
        packedValue = packedValue.pack(value1, bitLength, shiftAmount1, scalar1); // Pack value1 into the first 2 bits
        packedValue = packedValue.pack(value2, bitLength, shiftAmount2, scalar2); // Pack value2 into the next 2 bits
        packedValue = packedValue.pack(value3, bitLength, shiftAmount3, scalar3); // Pack value3 into the next 2 bits
        packedValue = packedValue.pack(value4, bitLength, shiftAmount4, scalar4); // Pack value4 into the last 2 bits

        // Unpack and check the values, scaling them back to the original values
        assertEq(packedValue.unpack(bitLength, shiftAmount1, scalar1), value1);
        assertEq(packedValue.unpack(bitLength, shiftAmount2, scalar2), value2);
        assertEq(packedValue.unpack(bitLength, shiftAmount3, scalar3), value3);
        assertEq(packedValue.unpack(bitLength, shiftAmount4, scalar4), value4);
    }

    function testPackWithDynamicScalar() public {
        // Let's assume the range of values is between 1 and 99
        uint8 value1 = 45; // Example value within the range 1-99
        uint8 value2 = 77; // Another example value within the range

        uint8 maxRangeValue = 99; // Maximum value in the range is 99
        uint8 bitLength = 4; // We're using 4 bits for each value, giving us a max storable value of 15

        // Dynamically calculate the scalar based on the range and bit length
        uint8 scalar1 = packedValue.calculateScalar(value1, bitLength);
        uint8 scalar2 = packedValue.calculateScalar(value2, bitLength);

        // Shift amounts for each value in the 8-bit storage
        uint8 shiftAmount1 = 0; // No shift for value1
        uint8 shiftAmount2 = 4; // Shift value2 by 4 bits

        // Pack the values into the 8-bit storage (we can only pack 2 values into 8 bits with 4 bits each)
        packedValue = packedValue.pack(value1, bitLength, shiftAmount1, scalar1); // Pack value1
        packedValue = packedValue.pack(value2, bitLength, shiftAmount2, scalar2); // Pack value2

        // Unpack and check the values, scaling them back to the original values
        assertEq(packedValue.unpack(bitLength, shiftAmount1, scalar1), value1);
        assertEq(packedValue.unpack(bitLength, shiftAmount2, scalar2), value2);
    }
}
