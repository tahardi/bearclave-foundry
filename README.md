# Bearclave: Foundry

[Foundry](https://github.com/foundry-rs/foundry) is a set of CLI tools for
Ethereum smart contract development. They allow you to build, test, and deploy
smart contracts locally. The Ethereum [abigen](https://github.com/ethereum/go-ethereum) tool can then be used to
generate Go bindings for calling deployed contracts.

While smart contract unit tests are straightforward with Foundry, integration
tests with the Go bindings require orchestrating setup/teardown of a local
Ethereum node and waiting on transactions to confirm expected behaviors.
While you could handle this with scripts or tools like docker compose, I wanted
the ability to do all test setup and teardown from within the Go test framework.

This repository contains a Golang test harness for the Foundry CLI. It allows
you to start and stop a local Ethereum evm, deploy smart contracts, and read
contract events from within your Go tests.

## Getting Started

1. Install [Golang](https://golang.org/doc/install) (v1.25.5 or higher) to build
   and run the integration tests.
2. Install the [Foundry](https://github.com/foundry-rs/foundry) toolset for
   smart contract development.
3. Initialize the submodules.
```bash
git submodule update --init --recursive 
```
4. Install [Slither](https://github.com/crytic/slither) static analysis tool for
   auditing smart contracts.

The code in `contracts/` and `test/integration/` demonstrates how to use the
test harness to run integration tests against a local Foundry node.
