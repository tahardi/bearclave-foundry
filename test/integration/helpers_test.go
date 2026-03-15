package integration_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tahardi/bearclave-foundry/contracts/bindings"
	"github.com/tahardi/bearclave-foundry/foundry"
)

func assertAddressesEqual(
	t *testing.T,
	address1 common.Address,
	address2 common.Address,
) {
	t.Helper()
	assert.Equal(t, 0, address1.Cmp(address2))
}

func approve(
	t *testing.T,
	f *foundry.Foundry,
	contract *bindings.KitchenSink,
	principal *foundry.Account,
	proxy *foundry.Account,
	amount *big.Int,
) (*types.Receipt, error) {
	t.Helper()
	call := func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return contract.Approve(opts, proxy.Address(), amount)
	}
	return foundry.CallContract(t.Context(), f, call, principal)
}

func burn(
	t *testing.T,
	f *foundry.Foundry,
	contract *bindings.KitchenSink,
	account *foundry.Account,
	amount *big.Int,
) (*types.Receipt, error) {
	t.Helper()
	call := func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return contract.Burn(opts, amount)
	}
	return foundry.CallContract(t.Context(), f, call, account)
}

func mint(
	t *testing.T,
	f *foundry.Foundry,
	contract *bindings.KitchenSink,
	owner *foundry.Account,
	to *foundry.Account,
	amount *big.Int,
) (*types.Receipt, error) {
	t.Helper()
	call := func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return contract.Mint(opts, to.Address(), amount)
	}
	return foundry.CallContract(t.Context(), f, call, owner)
}

func requireAllowance(
	t *testing.T,
	contract *bindings.KitchenSink,
	principal *foundry.Account,
	proxy *foundry.Account,
	want *big.Int,
) {
	t.Helper()
	got, err := contract.Allowance(nil, principal.Address(), proxy.Address())
	require.NoError(t, err)
	if want == nil {
		require.Equal(t, 0, got.Cmp(big.NewInt(0)))
	} else {
		require.Equal(t, want, got)
	}
}

func requireBalance(
	t *testing.T,
	contract *bindings.KitchenSink,
	account *foundry.Account,
	want *big.Int,
) {
	t.Helper()
	got, err := contract.BalanceOf(nil, account.Address())
	require.NoError(t, err)
	if want == nil {
		require.Equal(t, 0, got.Cmp(big.NewInt(0)))
	} else {
		require.Equal(t, want, got)
	}
}

func requireDeployKitchenSink(
	t *testing.T,
	f *foundry.Foundry,
	owner *foundry.Account,
) *bindings.KitchenSink {
	t.Helper()
	contract, err := foundry.DeployContract[bindings.KitchenSink](
		t.Context(),
		f,
		owner,
		ContractName,
		bindings.NewKitchenSink,
	)
	require.NoError(t, err)
	return contract
}

func requireMaxUint256(t *testing.T) *big.Int {
	t.Helper()
	maxUint256, ok := new(big.Int).
		SetString(
			"115792089237316195423570985008687907853269984665640564039457584007913129639935",
			10,
		)
	require.True(t, ok)
	return maxUint256
}

func requireBurnEvent(
	t *testing.T,
	contract *bindings.KitchenSink,
	receipt *types.Receipt,
	from common.Address,
	amount *big.Int,
) {
	t.Helper()
	event, err := foundry.GetEvent[*bindings.KitchenSinkBurn](receipt, contract.ParseBurn)
	require.NoError(t, err)

	assertAddressesEqual(t, from, event.From)
	if amount == nil {
		require.Equal(t, 0, event.Amount.Cmp(big.NewInt(0)))
	} else {
		require.Equal(t, amount, event.Amount)
	}
}

func requireMintEvent(
	t *testing.T,
	contract *bindings.KitchenSink,
	receipt *types.Receipt,
	to common.Address,
	amount *big.Int,
) {
	t.Helper()
	event, err := foundry.GetEvent[*bindings.KitchenSinkMint](receipt, contract.ParseMint)
	require.NoError(t, err)

	assertAddressesEqual(t, to, event.To)
	if amount == nil {
		require.Equal(t, 0, event.Amount.Cmp(big.NewInt(0)))
	} else {
		require.Equal(t, amount, event.Amount)
	}
}

func totalSupply() *big.Int {
	decimals := big.NewInt(Decimals)
	base := big.NewInt(Base)
	ten := big.NewInt(10)
	pow := big.NewInt(0).Exp(ten, decimals, nil)
	return big.NewInt(0).Mul(pow, base)
}

func transfer(
	t *testing.T,
	f *foundry.Foundry,
	contract *bindings.KitchenSink,
	from *foundry.Account,
	to *foundry.Account,
	amount *big.Int,
) (*types.Receipt, error) {
	t.Helper()
	call := func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return contract.Transfer(opts, to.Address(), amount)
	}
	return foundry.CallContract(t.Context(), f, call, from)
}

func transferFrom(
	t *testing.T,
	f *foundry.Foundry,
	contract *bindings.KitchenSink,
	principal *foundry.Account,
	proxy *foundry.Account,
	to *foundry.Account,
	amount *big.Int,
) (*types.Receipt, error) {
	t.Helper()
	call := func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return contract.TransferFrom(opts, principal.Address(), to.Address(), amount)
	}
	return foundry.CallContract(t.Context(), f, call, proxy)
}
