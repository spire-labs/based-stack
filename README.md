<div align="center">
  <br />
  <br />
  <a href="https://www.spire.dev"><img alt="Spire" src="https://github.com/user-attachments/assets/09a7fe72-1db2-4cce-a3c0-f34a16de4937" width=600></a>
  <br />
  <h3>Derived from<a href="https://optimism.io"> Optimism</a></h3>
  <br />
</div>

### ⚠️ This codebase is currently un-audited and not production ready

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
## Table of Contents

- [Table of Contents](#table-of-contents)
- [About Spire Labs?](#about-spire-labs)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [Directory Structure](#directory-structure)
- [License](#license)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## About Spire Labs?

[<i>**Spire Labs**</i>](https://www.spire.dev/) is building a based rollup stack called “Based Stack” that enables developers to spin up app-specific **based rollups**. Spire's sequencing design enables appchains to interact with Ethereum's liquidity and protocols natively (e.g. cross chain contract calls). Also, based rollups/appchains built the Based Stack can retain MEV, customize their execution environment, and integrate natively with decentralized preconfirmations for a lightning-fast UX.

## Documentation

Visit https://docs.spire.dev/ for documentation about this project.

## Contributing

Check out our [CONTRIBUTING.md](https://github.com/spire-labs/based-stack/blob/develop/CONTRIBUTING.md) file for a detailed explanation of the contributing process for this repository. Make sure to use the [Developer Quick Start](https://github.com/spire-labs/based-stack/blob/develop/CONTRIBUTING.md#development-quick-start) to properly set up your development environment.

## Directory Structure

<pre>
├── <a href="./docs">docs</a>: A collection of documents including audits and post-mortems
├── <a href="./op-batcher">op-batcher</a>: L2-Batch Submitter, submits bundles of batches to L1
├── <a href="./op-bootnode">op-bootnode</a>: Standalone op-node discovery bootnode
├── <a href="./op-chain-ops">op-chain-ops</a>: State surgery utilities
├── <a href="./op-challenger">op-challenger</a>: Dispute game challenge agent
├── <a href="./op-e2e">op-e2e</a>: End-to-End testing of all bedrock components in Go
├── <a href="./op-node">op-node</a>: rollup consensus-layer client
├── <a href="./op-preimage">op-preimage</a>: Go bindings for Preimage Oracle
├── <a href="./op-program">op-program</a>: Fault proof program
├── <a href="./op-proposer">op-proposer</a>: L2-Output Submitter, submits proposals to L1
├── <a href="./op-service">op-service</a>: Common codebase utilities
├── <a href="./op-ufm">op-ufm</a>: Simulations for monitoring end-to-end transaction latency
├── <a href="./op-wheel">op-wheel</a>: Database utilities
├── <a href="./ops">ops</a>: Various operational packages
├── <a href="./ops-bedrock">ops-bedrock</a>: Bedrock devnet work
├── <a href="./packages">packages</a>
│   ├── <a href="./packages/contracts-bedrock">contracts-bedrock</a>: OP Stack smart contracts
├── <a href="./proxyd">proxyd</a>: Configurable RPC request router and proxy
├── <a href="./specs">specs</a>: Specs of the rollup starting at the Bedrock upgrade
</pre>

## License

All other files within this repository are licensed under the [MIT License](https://github.com/spire-labs/based-stack/blob/develop/LICENSE) unless stated otherwise.
