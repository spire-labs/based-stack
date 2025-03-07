# Based Stack Smart Contracts

This package contains the L1 and L2 smart contracts for the Based Stack.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
## Table of Contents

- [Based Stack Smart Contracts](#based-stack-smart-contracts)
  - [Table of Contents](#table-of-contents)
  - [Deployment](#deployment)
  - [Generating L2 Genesis Allocs](#generating-l2-genesis-allocs)
    - [Configuration](#configuration)
      - [Custom Gas Token](#custom-gas-token)
    - [Execution](#execution)
    - [Deploying a single contract](#deploying-a-single-contract)
  - [Testing](#testing)
    - [Test Setup](#test-setup)
    - [Static Analysis](#static-analysis)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Deployment

The smart contracts are deployed using `foundry`. The `DEPLOYMENT_OUTFILE` env var will determine the filepath that the
deployment artifact is written to on disk after the deployment. It comes in the form of a JSON file where keys are
the names of the contracts and the values are the addresses the contract was deployed to.

The `DEPLOY_CONFIG_PATH` is a filepath to a deploy config file, see the `deploy-config` directory for examples and the
[DeployConfig](https://github.com/spire-labs/based-stack/blob/develop/op-chain-ops/genesis/config.go) definition for
descriptions of the values.
If you are following the official deployment tutorial, please make sure to use the `getting-started.json` file.

```bash
DEPLOYMENT_OUTFILE=deployments/artifact.json \
DEPLOY_CONFIG_PATH=<PATH_TO_MY_DEPLOY_CONFIG> \
  forge script scripts/deploy/Deploy.s.sol:Deploy \
  --broadcast --private-key $PRIVATE_KEY \
  --rpc-url $ETH_RPC_URL
```

The `IMPL_SALT` env var can be used to set the `create2` salt for deploying the implementation
contracts.

This will deploy an entire new system of L1 smart contracts including a new `SuperchainConfig`.

## Generating L2 Genesis Allocs

A foundry script is used to generate the L2 genesis allocs. This is a JSON file that represents the L2 genesis state.
The `CONTRACT_ADDRESSES_PATH` env var represents the deployment artifact that was generated during a contract deployment.
The same deploy config JSON file should be used for L1 contracts deployment as when generating the L2 genesis allocs.
The `STATE_DUMP_PATH` env var represents the filepath at which the allocs will be written to on disk.

```bash
CONTRACT_ADDRESSES_PATH=deployments/artifact.json \
DEPLOY_CONFIG_PATH=<PATH_TO_MY_DEPLOY_CONFIG> \
STATE_DUMP_PATH=<PATH_TO_WRITE_L2_ALLOCS> \
  forge script scripts/L2Genesis.s.sol:L2Genesis \
  --sig 'runWithStateDump()'
```

### Configuration

Create or modify a file `<network-name>.json` inside of the [`deploy-config`](./deploy-config/) folder.
Use the env var `DEPLOY_CONFIG_PATH` to use a particular deploy config file at runtime.

The script will read the latest active fork from the deploy config and the L2 genesis allocs generated will be
compatible with this fork. The automatically detected fork can be overwritten by setting the environment variable
`FORK` either to the lower-case fork name (currently `delta`, `ecotone`, `fjord`, or `granite`) or to `latest`, which
will select the latest fork available (currently `granite`).

By default, the script will dump the L2 genesis allocs of the detected or selected fork only, to the file at `STATE_DUMP_PATH`.
The optional environment variable `OUTPUT_MODE` allows to modify this behavior by setting it to one of the following values:
* `latest` (default) - only dump the selected fork's allocs.
* `all` - also dump all intermediary fork's allocs. This only works if `STATE_DUMP_PATH` is _not_ set. In this case, all allocs
          will be written to files `/state-dump-<fork>.json`. Another path cannot currently be specified for this use case.
* `none` - won't dump any allocs. Only makes sense for internal test usage.

#### Custom Gas Token

The Custom Gas Token feature is a Beta feature of the MIT licensed Based Stack.
While it has received initial review from core contributors, it is still undergoing testing, and may have bugs or other issues.

### Execution

Before deploying the contracts, you can verify the state diff produced by the deploy script using the `runWithStateDiff()` function signature which produces the outputs inside [`snapshots/state-diff/`](./snapshots/state-diff).
Run the deployment with state diffs by executing: `forge script -vvv scripts/deploy/Deploy.s.sol:Deploy --sig 'runWithStateDiff()' --rpc-url $ETH_RPC_URL --broadcast --private-key $PRIVATE_KEY`.

1. Set the env vars `ETH_RPC_URL`, `PRIVATE_KEY` and `ETHERSCAN_API_KEY` if contract verification is desired.
1. Set the `DEPLOY_CONFIG_PATH` env var to a path on the filesystem that points to a deploy config.
1. Deploy the contracts with `forge script -vvv scripts/deploy/Deploy.s.sol:Deploy --rpc-url $ETH_RPC_URL --broadcast --private-key $PRIVATE_KEY`
   Pass the `--verify` flag to verify the deployments automatically with Etherscan.

### Deploying a single contract

All of the functions for deploying a single contract are `public` meaning that the `--sig` argument to `forge script` can be used to
target the deployment of a single contract.

## Testing

### Test Setup

The Solidity unit tests use the same codepaths to set up state that are used in production. The same L1 deploy script is used to deploy the L1 contracts for the in memory tests
and the L2 state is set up using the same L2 genesis generation code that is used for production and then loaded into foundry via the `vm.loadAllocs` cheatcode. This helps
to reduce the overhead of maintaining multiple ways to set up the state as well as give additional coverage to the "actual" way that the contracts are deployed.

The L1 contract addresses are held in `deployments/hardhat/.deploy` and the L2 test state is held in a `.testdata` directory. The L1 addresses are used to create the L2 state
and it is possible for stale addresses to be pulled into the L2 state, causing tests to fail. Stale addresses may happen if the order of the L1 deployments happen differently
since some contracts are deployed using `CREATE`. Run `just clean` and rerun the tests if they are failing for an unknown reason.

### Static Analysis

`contracts-bedrock` uses [slither](https://github.com/crytic/slither) as its primary static analysis tool.
Slither will be run against PRs as part of CI, and new findings will be reported as a comment on the PR.
CI will fail if there are any new findings of medium or higher severity, as configured in the repo's Settings > Code Security and Analysis > Code Scanning > Protection rules setting.

There are two corresponding jobs in CI: one calls "Slither Analysis" and one called "Code scanning results / Slither".
The former will always pass if Slither runs successfully, and the latter will fail if there are any new findings of medium or higher severity.

Existing findings can be found in the repo's Security tab > [Code Scanning](https://github.com/ethereum-optimism/optimism/security/code-scanning) section.
You can view findings for a specific PR using the `pr:{number}` filter, such [`pr:9405`](https://github.com/ethereum-optimism/optimism/security/code-scanning?query=is:open+pr:9405).

For each finding, either fix it locally and push a new commit, or dismiss it through the PR comment's UI.

Note that you can run slither locally by running `slither .`, but because it does not contain the triaged results from GitHub, it will be noisy.
Instead, you should run `slither ./path/to/contract.sol` to run it against a specific file.
