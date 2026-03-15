package foundry

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"math/big"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	AnvilCommand = "anvil"
	AnvilVersion = "1.5.1-stable"

	BaseFee          = 1_000_000_000
	ChainID          = 31337
	GasLimit         = 30_000_000
	GenesisTimestamp = 1769011998
	GenesisNumber    = 0
	StartingBalance  = 10_000
	URL              = "http://127.0.0.1:8545"
	PollInterval     = 25 * time.Millisecond
)

var (
	ErrAnvil = errors.New("anvil")
)

type Anvil struct {
	accounts         []*Account
	baseFee          uint64
	chainID          *big.Int
	gasLimit         uint64
	genesisTimestamp uint64
	genesisNumber    uint64
	url              string
	client           *ethclient.Client
	server           *exec.Cmd
}

func NewAnvil(ctx context.Context) (*Anvil, error) {
	_, err := exec.LookPath(AnvilCommand)
	if err != nil {
		msg := AnvilCommand + " command not found in system path"
		return nil, anvilError(msg, err)
	}

	cmd := exec.CommandContext(ctx, AnvilCommand, "--version")
	out, err := cmd.Output()
	if err != nil {
		msg := "running " + AnvilCommand + " --version: "
		return nil, anvilError(msg, err)
	}

	version := strings.TrimSpace(string(bytes.TrimSpace(out)))
	if !strings.Contains(version, AnvilVersion) {
		slog.Warn(
			"version mismatch",
			"command", AnvilCommand,
			"expected", AnvilVersion,
			"got", version,
		)
	}

	accounts, err := NewDefaultAnvilAccounts()
	if err != nil {
		return nil, anvilError("creating accounts", err)
	}

	return &Anvil{
		accounts:         accounts,
		baseFee:          BaseFee,
		chainID:          big.NewInt(ChainID),
		gasLimit:         GasLimit,
		genesisTimestamp: GenesisTimestamp,
		genesisNumber:    GenesisNumber,
		url:              URL,
		client:           nil,
		server:           nil,
	}, nil
}

func (a *Anvil) Accounts() []*Account     { return a.accounts }
func (a *Anvil) Account(i int) *Account   { return a.accounts[i] }
func (a *Anvil) BaseFee() uint64          { return a.baseFee }
func (a *Anvil) ChainID() *big.Int        { return a.chainID }
func (a *Anvil) GasLimit() uint64         { return a.gasLimit }
func (a *Anvil) GenesisTimestamp() uint64 { return a.genesisTimestamp }
func (a *Anvil) GenesisNumber() uint64    { return a.genesisNumber }
func (a *Anvil) URL() string              { return a.url }

func (a *Anvil) Start(ctx context.Context, silent bool) error {
	a.server = exec.CommandContext(ctx, AnvilCommand)
	if silent {
		a.server.Stdout = io.Discard
		a.server.Stderr = io.Discard
	} else {
		a.server.Stdout = os.Stdout
		a.server.Stderr = os.Stderr
	}

	err := a.server.Start()
	if err != nil {
		return anvilError("starting anvil process", err)
	}

	err = a.waitUntilReady(ctx)
	if err != nil {
		return anvilError("waiting for anvil readiness", err)
	}
	return nil
}

func (a *Anvil) Stop() error {
	if a.client != nil {
		a.client.Close()
	}
	if a.server != nil {
		return a.server.Process.Kill()
	}
	return nil
}

func (a *Anvil) Client(ctx context.Context) (*ethclient.Client, error) {
	if a.client != nil {
		return a.client, nil
	}

	client, err := ethclient.DialContext(ctx, a.url)
	if err != nil {
		return nil, anvilError("dialing client", err)
	}
	a.client = client
	return a.client, nil
}

func (a *Anvil) waitUntilReady(ctx context.Context) error {
	ticker := time.NewTicker(PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			client, err := a.Client(ctx)
			if err != nil {
				continue
			}

			_, err = client.ChainID(ctx)
			if err != nil {
				continue
			}
			return nil
		}
	}
}
