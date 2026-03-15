package internal_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tahardi/bearclave-foundry/internal"
)

//go:embed testdata/inner.json
var innerJSON []byte

//go:embed testdata/transaction.json
var transactionJSON []byte

func TestInner_JSON(t *testing.T) {
	t.Run("happy path - round trip", func(t *testing.T) {
		// given
		want := innerJSON

		// when
		inner := &internal.InnerTransaction{}
		err := inner.UnmarshalJSON(want)
		require.NoError(t, err)

		got, err := inner.MarshalJSON()

		// then
		require.NoError(t, err)
		require.JSONEq(t, string(want), string(got))
	})
}

func TestTransaction_JSON(t *testing.T) {
	t.Run("happy path - round trip", func(t *testing.T) {
		// given
		want := transactionJSON

		// when
		transaction := &internal.Transaction{}
		err := transaction.UnmarshalJSON(want)
		require.NoError(t, err)

		got, err := transaction.MarshalJSON()

		// then
		require.NoError(t, err)
		require.JSONEq(t, string(want), string(got))
	})
}
