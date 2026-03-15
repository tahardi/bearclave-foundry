package foundry

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	ErrFoundry = errors.New("foundry")
)

type Foundry struct {
	anvil *Anvil
	forge *Forge
}

func NewFoundry(
	ctx context.Context,
	silent bool,
	broadcastDir string,
	scriptDir string,
) (*Foundry, error) {
	anvil, err := NewAnvil(ctx)
	if err != nil {
		return nil, foundryError("creating anvil", err)
	}

	err = anvil.Start(ctx, silent)
	if err != nil {
		return nil, foundryError("starting anvil", err)
	}

	forge, err := NewForge(ctx, broadcastDir, scriptDir, anvil.URL(), anvil.ChainID())
	if err != nil {
		return nil, foundryError("creating forge", err)
	}

	return &Foundry{
		anvil: anvil,
		forge: forge,
	}, nil
}

func (f *Foundry) Stop() {
	_ = f.anvil.Stop()
}

func (f *Foundry) Anvil() *Anvil {
	return f.anvil
}

func (f *Foundry) Forge() *Forge {
	return f.forge
}

func DeployContract[T any](
	ctx context.Context,
	f *Foundry,
	owner *Account,
	contractName string,
	newContract func(address common.Address, backend bind.ContractBackend) (*T, error),
) (*T, error) {
	contractAddress, err := f.Forge().DeployContract(ctx, contractName, owner)
	if err != nil {
		return nil, foundryError("deploying contract", err)
	}

	client, err := f.Anvil().Client(ctx)
	if err != nil {
		return nil, foundryError("getting anvil client", err)
	}

	contract, err := newContract(*contractAddress, client)
	if err != nil {
		return nil, foundryError("creating contract", err)
	}
	return contract, nil
}

func CallContract(
	ctx context.Context,
	f *Foundry,
	contractCall func(opts *bind.TransactOpts) (*types.Transaction, error),
	sender *Account,
) (*types.Receipt, error) {
	opts, err := bind.NewKeyedTransactorWithChainID(
		sender.PrivateKey(),
		f.Anvil().ChainID(),
	)
	if err != nil {
		return nil, foundryError("creating transactor", err)
	}

	tx, err := contractCall(opts)
	if err != nil {
		return nil, foundryError("calling contract", err)
	}

	client, err := f.Anvil().Client(ctx)
	if err != nil {
		return nil, foundryError("getting anvil client", err)
	}

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return nil, foundryError("waiting for transaction", err)
	}
	return receipt, nil
}

func GetEvent[T any](
	receipt *types.Receipt,
	eventParser func(log types.Log) (T, error),
) (T, error) {
	for _, log := range receipt.Logs {
		event, err := eventParser(*log)
		if err != nil {
			continue
		}
		return event, nil
	}

	var zero T
	return zero, foundryError("no event found", nil)
}
