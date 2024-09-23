// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.15;

import { Test } from "forge-std/Test.sol";
import "src/L1/ValidatorLookahead.sol";

contract ValidatorLookaheadTest is Test {
    ValidatorLookahead lookahead;

    function setUp() public {
        lookahead = new ValidatorLookahead();
    }

    function testWriteAndReadEpoch() public {
       ValidatorLookahead.ValidatorData[32] memory newValidatorData;
        newValidatorData[0] = ValidatorLookahead.ValidatorData(
            hex"b21473292db9acaac18f102de23410c4f13eed817c5c3c904c0c37f5a20c9ad4701a6424e71c9d45a5bb47ed2186bc1e",
            1447741,
            10013248
        );
        newValidatorData[1] = ValidatorLookahead.ValidatorData(
            hex"93a9460486fd119973cf209c5022c04b7e7cbf597644c216e2748ed3049310d1fc068bbe4eebbd1b0d194d6062ef6cea",
            1537318,
            10013249
        );
        newValidatorData[2] = ValidatorLookahead.ValidatorData(
            hex"a9e2860885798e355b330466f76324629cb907123d8568100d9bf3701578cee6a4d799c7a22b9c539ebc91547f0afdbc",
            411164,
            10013250
        );
        newValidatorData[3] = ValidatorLookahead.ValidatorData(
            hex"a576d7de216c471118823f7d77cbb0bd2f0c2386db094c4fea562218e66c1bed4cf4a463b934967847b235e647bbb641",
            1274803,
            10013251
        );
        newValidatorData[4] = ValidatorLookahead.ValidatorData(
            hex"8b2827266d74ee284640388b989888c533ed5d39a5de6b4bcd02ec3e2201965a65530b8ac4160d77ae8df76d955604e2",
            1044927,
            10013252
        );
        newValidatorData[5] = ValidatorLookahead.ValidatorData(
            hex"b7be504318baa32c3d1af9d86de8cba25e693a677d89629ed89cfa52dc28e58c4ce2521cf90b4add6e8f9d6901cdbb7c",
            900626,
            10013253
        );
        newValidatorData[6] = ValidatorLookahead.ValidatorData(
            hex"a4d2ca6130528436a3cf8cdabc4d10ba03bee7f046a712e76551439044fa8eeee4e646a68eb72052c3a7528beb8339d9",
            556089,
            10013254
        );
        newValidatorData[7] = ValidatorLookahead.ValidatorData(
            hex"b3e7ec69d0bdf0cdfc4c3be75cb9ae4d375fe2edbfc37b511bd90c81ccc5de6f2a6c0b4b32f3e34fbbaf90db1b90e5a5",
            1259185,
            10013255
        );
        newValidatorData[8] = ValidatorLookahead.ValidatorData(
            hex"b0025e9c3fb038ef3317feb42bc3d7ad0e5387786080aa3bc95da55be26e9fa76cd504bf09f92787d2b5bf056cd07562",
            806981,
            10013256
        );
        newValidatorData[9] = ValidatorLookahead.ValidatorData(
            hex"ac94f1c9cc6c9facf81703db2609bb13faef9f69c818090c441096c9d665df1748c6cae2cead23053d5858ff16f91cac",
            1543449,
            10013257
        );
        newValidatorData[10] = ValidatorLookahead.ValidatorData(
            hex"93457fa367f234f5bf4c0aee6f0c301fe7145c06d741e563f9175324638d974f35f92fa7a0894e5b70cfde399f56344f",
            2711,
            10013258
        );
        newValidatorData[11] = ValidatorLookahead.ValidatorData(
            hex"aa00e178f6d5499d635504f8da396152d3786bdad1c97e31a747b00d95545dc9bf4bcec458039e54002aced86b71a327",
            1042689,
            10013259
        );
        newValidatorData[12] = ValidatorLookahead.ValidatorData(
            hex"95c302129fd03ba825972c8ea4950e5279845ac995caca053c2dabfeedbb28458cd6796cc646d74a52d1e5353d7a44c5",
            989502,
            10013260
        );
        newValidatorData[13] = ValidatorLookahead.ValidatorData(
            hex"8a8d31bb4bac741f0f2ab97d95322ae5bf3310b9949052913c3443f9edcd26d2d9d5408ae30ab17217c03a41e477a8bc",
            548815,
            10013261
        );
        newValidatorData[14] = ValidatorLookahead.ValidatorData(
            hex"a78cd7e0c7aa3c50f7e0032a6ea86bb42a99e7c3d4cdff33a16d42fa5b1f69affb0b441ab7e0a922485793a2362f155b",
            964048,
            10013262
        );
        newValidatorData[15] = ValidatorLookahead.ValidatorData(
            hex"b3f8b0ecb687ec0ca3d234b469b53869181bff01bff6524e64b5437432efa317f7e7cb8dca7506f9c2f2dedc35d2debc",
            960007,
            10013263
        );
        newValidatorData[16] = ValidatorLookahead.ValidatorData(
            hex"8c9d3b00111b12c08e72acfb5132f36f5e825e412066860cca3014f3f66b2a4d465b96c8f2aa3543198fba8f1b4d6171",
            526992,
            10013264
        );
        newValidatorData[17] = ValidatorLookahead.ValidatorData(
            hex"910e07f5e35f54d3cd7c8c2c0d5d505bad5d4abc288eb0841faf2df9f61887f7647a8cec0654c3e9ae767f54bc7cdf68",
            639376,
            10013265
        );
        newValidatorData[18] = ValidatorLookahead.ValidatorData(
            hex"aa572317455945bd44f427bed0edd5440db8b2b6aa0cf8634e614849a7902d11d7595ec94d4ca52ac364fd13db13c213",
            1184892,
            10013266
        );
        newValidatorData[19] = ValidatorLookahead.ValidatorData(
            hex"94e4df12b200baf64e04826492fe0d833bdb567918d79f0e64beefb29f5d3135b493995a08e3db8b87f14ece6aaf3839",
            856510,
            10013267
        );
        newValidatorData[20] = ValidatorLookahead.ValidatorData(
            hex"b5ff28b4a5aa0f6dc030dbdf148abd2b0e2178ba4bf90630eb2bbadd6a992ffdbaaeeff799272096e7d93b4b04c03825",
            888531,
            10013268
        );
        newValidatorData[21] = ValidatorLookahead.ValidatorData(
            hex"b74369e0df8014a40550a2c3ccb26d7faa598e1c9d8f04191a9a80f389eab5ed93efb9d6bdad5c1ff9405e314cc9774e",
            733388,
            10013269
        );
        newValidatorData[22] = ValidatorLookahead.ValidatorData(
            hex"aa330fe2777b8ad76cb7a0840fb804c89b8cfd34454ad7c8d8454a74077e04bb429ec46044ecd54cf49ac3ee2e4cb2e3",
            778253,
            10013270
        );
        newValidatorData[23] = ValidatorLookahead.ValidatorData(
            hex"b5f02745e5e17e210eab2cbc07dc9f4e40164db09e83d1a1eda45c27918368e5859b19e27ebf6a8461d1095419475dd9",
            1012094,
            10013271
        );
        newValidatorData[24] = ValidatorLookahead.ValidatorData(
            hex"a2540d16552fcc6e205d3c28e0e19d48033136d6415fdf2de6d0d9592467c7d66e0c8edc8821c99e5f41857636c34821",
            1545329,
            10013272
        );
        newValidatorData[25] = ValidatorLookahead.ValidatorData(
            hex"8ab5146e1091580f4bfb922670cf85f380e8c084363313544af4756013da7c5b11b86a79d0641bf23f2f991452e233ab",
            855198,
            10013273
        );
        newValidatorData[26] = ValidatorLookahead.ValidatorData(
            hex"a083146b950c3c929772e48b872263bcae4a8e2a530becbbe8e5c4b1b0febd46ab32706384cd086a4730f510a50b2a1b",
            1138383,
            10013274
        );
        newValidatorData[27] = ValidatorLookahead.ValidatorData(
            hex"a4e878fc8769b33b759fc0fb07c4349f7469719624d869fa81a403aa56d90b4500da9365973a0326b89e24435b759059",
            1498898,
            10013275
        );
        newValidatorData[28] = ValidatorLookahead.ValidatorData(
            hex"87e2993111ba3a1447f2af1baf1943ed3e450bc2479af9626519e80f6bbfb344a79c9fabd6a3fc47e1b614e02fa1338e",
            299227,
            10013276
        );
        newValidatorData[29] = ValidatorLookahead.ValidatorData(
            hex"99c7fb88e4d6bcd68b243acd6e7dce6020f8bca560106475db102a5c849a4750b82390801e7c73780cc5f7bd03458c6e",
            1391411,
            10013277
        );
        newValidatorData[30] = ValidatorLookahead.ValidatorData(
            hex"b8809e7f51c1ffc7f2162f40d7ad3b8a8a4814b2fbd128b790040b3f09b373a0be48ec08e1769063e43b2646fcd8a266",
            238552,
            10013278
        );
        newValidatorData[31] = ValidatorLookahead.ValidatorData(
            hex"b615ad419badac7afb07a37c3b793436aa15537ccd823a0888332cdae72d1e486dc31b8f524db784ef80580704d204dd",
            822661,
            10013279
        );


        lookahead.storeValidatorLookahead(newValidatorData, 312924);

        // read back to ensure data was stored correctly
        ValidatorLookahead.ValidatorData[32] memory storedData = lookahead.getValidatorLookahead();

        // Verify that the stored data matches the original data
        for (uint256 i = 0; i < 32; i++) {
            assertEq(storedData[i].validatorIndex, newValidatorData[i].validatorIndex);
            assertEq(storedData[i].slot, newValidatorData[i].slot);
            assertEq(storedData[i].pubkey, newValidatorData[i].pubkey);
        }
    }
}