# Spire Based Stack Contributing Guide

## Disclaimer

- <b>Please keep in mind that we DO NOT have a token and we are not doing airdrops. Even if we do plan to launch a token in the future, contributing to this codebase will not be rewarded with an airdrop</b>
- <b>Please refrain from opening low effort PRs. PRs for typo fix and modifying markdown newlines will NOT be accepted as first time contributions.</b>

## What to Contribute

Welcome to the Spire Based stack Contributing Guide! If you're reading this then you might be interested in contributing to the Spire Based stack. Please keep in mind the following while contributing:

- Report issues in this repository. Great bug reports are detailed and give clear instructions for how a developer can reproduce the problem. Write good bug reports and developers will love you.
- If you are trying to solve a particular issue, mention it in your PR description as `Closes #<issue-number>`. Mention enough context for the developers to understand your thinking process.
- Help improve the [Spire Based stack Developer Docs](https://github.com/spire-labs/docs) by reporting issues, fixing typos, or adding missing sections.
- Make sure all the tests pass and the code is lint-free

## Code of Conduct

Interactions within this repository are subject to a [Code of Conduct](https://github.com/spire-labs/based-stack/tree/develop/.github/CODE_OF_CONDUCT.md) adapted from the [Contributor Covenant](https://www.contributor-covenant.org/version/1/4/code-of-conduct/).

## Development Quick Start

### Software Dependencies

| Dependency                                                    | Version  | Version Check Command    |
| ------------------------------------------------------------- | -------- | ------------------------ |
| [git](https://git-scm.com/)                                   | `^2`     | `git --version`          |
| [go](https://go.dev/)                                         | `^1.21`  | `go version`             |
| [node](https://nodejs.org/en/)                                | `^20`    | `node --version`         |
| [nvm](https://github.com/nvm-sh/nvm)                          | `^0.39`  | `nvm --version`          |
| [just](https://github.com/casey/just)                         | `^1.34.0`| `just --version`         |
| [foundry](https://github.com/foundry-rs/foundry#installation) | `^0.2.0` | `forge --version`        |
| [make](https://linux.die.net/man/1/make)                      | `^3`     | `make --version`         |
| [jq](https://github.com/jqlang/jq)                            | `^1.6`   | `jq --version`           |
| [direnv](https://direnv.net)                                  | `^2`     | `direnv --version`       |
| [docker](https://docs.docker.com/get-docker/)                 | `^24`    | `docker --version`       |
| [docker compose](https://docs.docker.com/compose/install/)    | `^2.23`  | `docker compose version` |

### Notes on Specific Dependencies

#### `node`

Make sure to use the version of `node` specified within [`.nvmrc`](./.nvmrc).
You can use [`nvm`](https://github.com/nvm-sh/nvm) to manage multiple versions of Node.js on your machine and automatically switch to the correct version when you enter this repository.

#### `foundry`

`foundry` is updated frequently and occasionally contains breaking changes.
This repository pins a specific version of `foundry` inside of [`versions.json`](./versions.json).
Use the command `just update-foundry` at the root of the repo to make sure that your version of `foundry` is the same as the one currently being used in CI.

#### `direnv`

[`direnv`](https://direnv.net) is a tool used to load environment variables from [`.envrc`](./.envrc) into your shell so you don't have to manually export variables every time you want to use them.
`direnv` only has access to files that you explicitly allow it to see.
After [installing `direnv`](https://direnv.net/docs/installation.html), you will need to **make sure that [`direnv` is hooked into your shell](https://direnv.net/docs/hook.html)**.
Make sure you've followed [the guide on the `direnv` website](https://direnv.net/docs/hook.html), then **close your terminal and reopen it** so that the changes take effect (or `source` your config file if you know how to do that).

#### `docker compose`

[Docker Desktop](https://docs.docker.com/get-docker/) should come with `docker compose` installed by default.
You'll have to install the `compose` plugin if you're not using Docker Desktop or you're on linux.

### Setting Up

Clone the repository and open it:

```bash
git clone git@github.com:ethereum-optimism/optimism.git
cd optimism
```

### Building the repo

Make sure that you've installed all of the required [Software Dependencies](#software-dependencies) before you continue.
You will need [foundry](https://github.com/foundry-rs/foundry) to build the smart contracts found within this repository.
Refer to the note on [foundry as a dependency](#foundry) for instructions.

Install dependencies and build all packages within the repo by running:

```bash
make build
```

Packages built on one branch may not be compatible with packages on a different branch.
**You should rebuild the repo whenever you move from one branch to another.**
Use the above command to rebuild the repo.

