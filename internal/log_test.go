package internal_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tahardi/bearclave-foundry/internal"
)

//go:embed testdata/log.json
var logJSON []byte

func TestLog_JSON(t *testing.T) {
	t.Run("happy path - round trip", func(t *testing.T) {
		// given
		want := logJSON

		// when
		log := &internal.Log{}
		err := log.UnmarshalJSON(want)
		require.NoError(t, err)

		got, err := log.MarshalJSON()

		// then
		require.NoError(t, err)
		require.JSONEq(t, string(want), string(got))
	})
}
