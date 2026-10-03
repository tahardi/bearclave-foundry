package foundry

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"math/big"
	"os/exec"
	"strings"
)

const (
	CastCommand    = "cast"
	CastVersion    = "1.8.4"
	BalanceCommand = "balance"
	SendCommand    = "send"

	ValueFlag = "--value"
	Base10    = 10
)

var (
	ErrCast = errors.New("cast")
)

type Denomination string

const (
	Ether Denomination = "ether"
	Wei   Denomination = "gwei"
	Raw   Denomination = "raw"
)

type Cast struct {
	url     string
	chainID *big.Int
}

func NewCast(
	ctx context.Context,
	url string,
	chainID *big.Int,
) (*Cast, error) {
	_, err := exec.LookPath(CastCommand)
	if err != nil {
		msg := CastCommand + " command not found in system path"
		return nil, anvilError(msg, err)
	}

	cmd := exec.CommandContext(ctx, CastCommand, "--version")
	out, err := cmd.Output()
	if err != nil {
		msg := "running " + CastCommand + " --version: "
		return nil, anvilError(msg, err)
	}

	version := strings.TrimSpace(string(bytes.TrimSpace(out)))
	if !strings.Contains(version, CastVersion) {
		slog.Warn(
			"version mismatch",
			"command", CastCommand,
			"expected", CastVersion,
			"got", version,
		)
	}

	return &Cast{
		url:     url,
		chainID: chainID,
	}, nil
}

func (c *Cast) URL() string       { return c.url }
func (c *Cast) ChainID() *big.Int { return c.chainID }

func (c *Cast) Balance(ctx context.Context, account *Account) (*big.Int, error) {
	args := []string{
		BalanceCommand, account.Address().Hex(),
		RPCFlag, c.url,
	}

	getBalance := exec.CommandContext(ctx, CastCommand, args...)
	out, err := getBalance.CombinedOutput()
	if err != nil {
		msg := "getting balance: " + string(out)
		return nil, castError(msg, err)
	}

	balanceStr := strings.TrimSpace(string(out))
	balance := new(big.Int)
	if _, ok := balance.SetString(balanceStr, Base10); !ok {
		return nil, castError("parsing balance", nil)
	}
	return balance, nil
}

func (c *Cast) SendEther(
	ctx context.Context,
	from *Account,
	to *Account,
	amount *big.Int,
	denomination Denomination,
) error {
	amountString := amount.String()
	switch denomination {
	case Ether:
		amountString += string(Ether)
	case Wei:
		amountString += string(Wei)
	case Raw:
		break
	default:
		msg := "invalid denomination: " + string(denomination)
		return castError(msg, nil)
	}

	args := []string{
		SendCommand, to.Address().Hex(),
		ValueFlag, amountString,
		PrivateKeyFlag, from.PrivateKeyHex(),
		RPCFlag, c.url,
	}

	send := exec.CommandContext(ctx, CastCommand, args...)
	out, err := send.CombinedOutput()
	if err != nil {
		msg := "sending ether: " + string(out)
		return castError(msg, err)
	}
	return nil
}
