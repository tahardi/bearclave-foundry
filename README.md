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

## Installation

Requires [Golang](https://golang.org/doc/install) v1.25.0 or higher.

```bash
go get github.com/tahardi/bearclave-foundry@v0.1.0
```

Install the [Foundry](https://github.com/foundry-rs/foundry) toolset. The
harness is tested with foundry `1.8.4` and warns if it finds another
version.

```bash
curl -L https://foundry.paradigm.xyz | bash && foundryup
```

## Usage

```go
f, err := foundry.NewFoundry(t.Context(), true, BroadcastDir, ScriptDir)
require.NoError(t, err)
defer f.Stop()

owner := f.Anvil().Account(0)
contract, err := foundry.DeployContract(t.Context(), f, owner, "KitchenSink", bindings.NewKitchenSink)
require.NoError(t, err)

call := func(opts *bind.TransactOpts) (*types.Transaction, error) {
	return contract.Burn(opts, big.NewInt(100))
}
receipt, err := foundry.CallContract(t.Context(), f, call, owner)
require.NoError(t, err)

event, err := foundry.GetEvent[*bindings.KitchenSinkBurn](receipt, contract.ParseBurn)
require.NoError(t, err)
```

`NewFoundry` starts a local anvil node. `DeployContract` runs a Forge script
that you provide, so it relies on a naming convention. For a contract named
`<Name>`, create `<scriptDir>/<Name>.s.sol` containing a contract named
`<Name>Script` that deploys the contract and broadcasts the deployment.

The code in `contracts/` and `test/integration/` shows a complete example.

## Development

1. Install [Foundry](https://github.com/foundry-rs/foundry).
2. Install [Slither](https://github.com/crytic/slither): `pipx install slither-analyzer`
3. Initialize the submodules: `git submodule update --init --recursive`
4. Run all checks: `make pre-pr`
