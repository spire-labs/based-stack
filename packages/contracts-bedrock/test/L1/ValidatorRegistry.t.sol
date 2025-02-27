// SPDX-License-Identifier: MIT
pragma solidity 0.8.25;

import { Test } from "forge-std/Test.sol";
import { ValidatorRegistry } from "src/L1/ValidatorRegistry.sol";
import { SSZ } from "src/libraries/SSZ.sol";

contract ValidatorRegistry_Test is Test {
    ValidatorRegistry public registry;
    bytes32[] public proof = [
        bytes32(0xcf6a41338cd11cde6f9367651e0934a3be310388eb0d306d7e175d5b1f4810e9),
        bytes32(0x532fb6e2b256b136bd83d4d641e6f04e03a838e40ae8c0cd41efb28b29484f15),
        bytes32(0x6fd2159bf85d39214268277c02c64dca3db4b1a098955c5bb834a7bd483ceabd),
        bytes32(0xd6ae0614839a25c593c1e9a545037fa147357e46c9b9ee8e7572b0eaf7508242),
        bytes32(0x97ad2f340a9e35d08303bed926cdfc7b1cd7974446c43e077223fda471701f10),
        bytes32(0x7f88dc656cadd5b529c3bc6f5de3d2b9942c691efdb5f43f49effa81cac89ef5),
        bytes32(0x99d928f79357f4bb58af5b653769fda94bffded07810b5493ce79779478b5b62),
        bytes32(0x8bb1f6f710e4b270d0f588eda78756c51c296e106fee89d8244b4babdd624ed8),
        bytes32(0x41d8cd46f2c703c33375843631ccff00ec62968b2ba9b3c538c43fba29491b72),
        bytes32(0x24d9fc75c7a511b6f576f12db7259c406e82c0d30ac584edc7b9de1bc7b6c59d),
        bytes32(0xb4d07f7688de8304de2c8bcf620f1b6b21fa2cbfbe79a4fb42a41bac0ad29bc7),
        bytes32(0x891054d1e18a46c10665aa08262acac31764e2e7b7473ce02c0b8ec6def4e601),
        bytes32(0xfc3bedb43644ae52f635473d6ef70d3eaa285cb991a31dada307a9565d40a9a0),
        bytes32(0x799a334fafd2bd458b885c99cb89ee811e6612d56eb70188103bd820f0cdffde),
        bytes32(0x577b4edddb8c1deea223e00ce8897628382d0c0aaeb36eea1f9338e1d2470b5b),
        bytes32(0x159a64a7290879f113c6208f9e3fb39500f7fa827b207f8d51e016b5642d3cd9),
        bytes32(0x40b6f71a6269c8985364ca110596b9c173cf127df1d2b030d019d9885aa4cf6c),
        bytes32(0xb0bd51b36e76f42cb3f739cf932030911e8451ef8516e2aaca525b53fed6b99c),
        bytes32(0x908898cd0e59e23bff58328991c6c6fd4b341e5f178a9b9e055380a5647f24e5),
        bytes32(0xdd8ebdd2af2b2607a2e2a7c9cdd28593fb7c5e5192982b3b21d3e92a527b8c60),
        bytes32(0xb9019b833bb16c0b40cc2f15a749a4be385ec775b8f04589ae30724dbc091f36),
        bytes32(0x8a8d7fe3af8caa085a7639a832001457dfb9128a8061142ad0335629ff23ff9c),
        bytes32(0xfeb3c337d7a51a6fbf00b9e34c52e1c9195c969bd4e7a0bfd51d5c5bed9c1167),
        bytes32(0xe71f0aa83cc32edfbefa9f4d3e0174ca85182eec9f3a09f6a6c0df6377a510d7),
        bytes32(0x31206fa80a50bb6abe29085058f16212212a60eec8f049fecb92d8c8e0a84bc0),
        bytes32(0x21352bfecbeddde993839f614c3dac0a3ee37543f9b412b16199dc158e23b544),
        bytes32(0x619e312724bb6d7c3153ed9de791d764a366b389af13c58bf8a8d90481a46765),
        bytes32(0x7cdd2986268250628d0c10e385c58c6191e6fbe05191bcc04f133f2cea72c1c4),
        bytes32(0x848930bd7ba8cac54661072113fb278869e07bb8587f91392933374d017bcbe1),
        bytes32(0x8869ff2c22b28cc10510d9853292803328be4fb0e80495e8bb8d271f5b889636),
        bytes32(0xb5fe28e79f1b850f8658246ce9b6a1e7b49fc06db7143e8fe0b4f2b0c5523a5c),
        bytes32(0x985e929f70af28d0bdd1a90a808f977f597c7c778c489e98d3bd8910d31ac0f7),
        bytes32(0xc6f67e02e6e4e1bdefb994c6098953f34636ba2b6ca20a4721d2b26a886722ff),
        bytes32(0x1c9a7e5ff1cf48b4ad1582d3f4e4a1004f3b20d8c5a2b71387a4254ad933ebc5),
        bytes32(0x2f075ae229646b6f6aed19a5e372cf295081401eb893ff599b3f9acc0c0d3e7d),
        bytes32(0x328921deb59612076801e8cd61592107b5c67c79b846595cc6320c395b46362c),
        bytes32(0xbfb909fdb236ad2411b4e4883810a074b840464689986c3f8a8091827e17c327),
        bytes32(0x55d8fb3687ba3ba49f342c77f5a1f89bec83d811446e1a467139213d640b6a74),
        bytes32(0xf7210d4f8e7e1039790e7bf4efa207555a10a6db1dd4b95da313aaa88b88fe76),
        bytes32(0xad21b516cbc645ffe34ab5de1c8aef8cd4e7f8d2b51e8e1456adc7563cda206f),
        bytes32(0x3437180000000000000000000000000000000000000000000000000000000000),
        bytes32(0x25c6190000000000000000000000000000000000000000000000000000000000),
        bytes32(0x3dcb9b4645a0a4a65cc1f72b8b8c4e644392da51c51b9ea28839b3dde1a1198c),
        bytes32(0x3788ecd7dc04418be021f25c1115b5e2ed1a38834b9c4b252ef7f618c4080e0c),
        bytes32(0x60d14d895b2b9417026ab4838fd533ec2dfb079eca6c2e46bc5bad53dda730d9),
        bytes32(0x1245953df09f1013ec92045421475b93a7effda1109853f57e564b77ab373940),
        bytes32(0xd3ad9b21b16725992b033bb5fc5a367ceb000e5b79766ee91e4259db0b4398ac),
        bytes32(0xc0cf9dc3e6b0ec473c95186bf71c9819a79888d7445c835956c7abc04590c504),
        bytes32(0x049f2496c622abed3ca813aa87dc419085a1ba824a8c439af841cdfb4f986183)
    ];
    uint256 internal constant _FORK_BLOCK = 20814319;

    function setUp() public {
      // random free tier public infura key cause archive node needed for this block fork
      // remove for env variable in production
        vm.createSelectFork(vm.rpcUrl("https://mainnet.infura.io/v3/1a2bbfa3895b4335b74a963eda5aadc5"), _FORK_BLOCK);
        registry = new ValidatorRegistry();
    }

    event ValidatorRegistered(uint256 indexed validatorIndex, address indexed delegatedAddress);

    function test_registerValidator() public {
        SSZ.Validator memory _validator = SSZ.Validator({
            pubkey: hex"8fb78536e82bcec34e98fff85c907f0a8e6f4b1ccdbf1e8ace26b59eb5a06d16f34e50837f6c490e2ad6a255db8d543b",
            //hex"8fb78536e82bcec34e98fff85c907f0a8e6f4b1ccdbf1e8ace26b59eb5a06d16f34e50837f6c490e2ad6a255db8d543b"
            withdrawalCredentials: bytes32(0x0100000000000000000000005e0ce4d7a65e4de934b75b7e48a447e299ff23d3),
            effectiveBalance: 32000000000,
            slashed: false,
            activationEligibilityEpoch: 2593,
            activationEpoch: 5890,
            exitEpoch: type(uint64).max,
            withdrawableEpoch: type(uint64).max
        });

        uint256 _validatorIndex = 44444;
        address _expectedDelegatedAddress = 0x5e0cE4d7a65e4dE934b75b7e48A447e299ff23D3;

        vm.expectEmit(true, true, true, true);
        emit ValidatorRegistered(_validatorIndex, _expectedDelegatedAddress);
        registry.registerValidator(_validator, _validatorIndex, proof);

        (
            bytes memory _expectedPubkey,
            bytes32 _expectedWithdrawalCredentials,
            uint64 _expectedEffectiveBalance,
            bool _expectedSlashed,
            uint64 _expectedActivationEligibilityEpoch,
            uint64 _expectedActivationEpoch,
            uint64 _expectedExitEpoch,
            uint64 _expectedWithdrawableEpoch
        ) = registry.validatorIndexToValidator(_validatorIndex);

        assertEq(_expectedPubkey, _validator.pubkey);
        assertEq(_expectedWithdrawalCredentials, _validator.withdrawalCredentials);
        assertEq(_expectedEffectiveBalance, _validator.effectiveBalance);
        assertEq(_expectedSlashed, _validator.slashed);
        assertEq(_expectedActivationEligibilityEpoch, _validator.activationEligibilityEpoch);
        assertEq(_expectedActivationEpoch, _validator.activationEpoch);
        assertEq(_expectedExitEpoch, _validator.exitEpoch);
        assertEq(_expectedWithdrawableEpoch, _validator.withdrawableEpoch);
        assertEq(registry.validatorIndexToDelegatedAddress(_validatorIndex), _expectedDelegatedAddress);
    }
}
