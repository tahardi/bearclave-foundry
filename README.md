# Bearclave: Foundry

[Foundry](https://github.com/foundry-rs/foundry) is a set of CLI tools for
Ethereum smart contract development. They allow you to build, test, and deploy
smart contracts locally. To interact with deployed contracts, I use the Ethereum
[abigen](https://github.com/ethereum/go-ethereum) tool to generate Go bindings
for calling the deployed contracts.

While smart contract unit tests are straightforward with Foundry, integration
tests required orchestrating setup/teardown of a local Ethereum node and
polling for contract events to confirm expected behaviors. Instead of doing
this through a script, I wanted the ability to do all test setup and teardown
within the Go test framework.

This repository contains a Golang test harness for the Foundry CLI. It allows
you to start and stop a local Ethereum node, deploy smart contracts, and read
contract events all from within your Go tests.

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
