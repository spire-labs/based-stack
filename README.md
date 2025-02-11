<div align="center">
  <br />
  <br />
  <a href="https://optimism.io"><img alt="Optimism" src="https://raw.githubusercontent.com/ethereum-optimism/brand-kit/main/assets/svg/OPTIMISM-R.svg" width=600></a>
  <br />
  <h3>Derived from<a href="https://optimism.io">Optimism</a></h3>
  <br />
</div>

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [About Spire Labs?](#about-spire-labs)
- [Documentation](#documentation)
- [Spinning up local OP Stack](#spinning-up-local-op-stack)
- [Specification](#specification)
- [Community](#community)
- [Contributing](#contributing)
- [Directory Structure](#directory-structure)
- [License](#license)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## About Spire Labs?

[<i>Spire</i>](https://www.spire.dev/) is developing the first <b>based rollup framework</b> & network on Ethereum, enabling a wide range of app developers to create their own appchains on Ethereum. Spire's sequencing design enables appchains to interact with Ethereum's liquidity and protocols natively (think: cross chain swaps). Spire appchains support a wide range of execution environments, scalability upgrades, and decentralized preconfirmations for lightning -fast UX.

## Documentation

Coming Soon!

## Spinning up local OP Stack

Coming Soon!

## Specification

Coming Soon!

## Community

Coming Soon!

## Contributing

Coming Soon!

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

All other files within this repository are licensed under the [MIT License](https://github.com/ethereum-optimism/optimism/blob/master/LICENSE) unless stated otherwise.
