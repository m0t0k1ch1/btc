package btc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScriptMapping(t *testing.T) {
	UseTestnet()

	testCases := []struct {
		hex string
		asm string
	}{
		{
			"",
			"",
		},
		{
			"76a914cbc222711a230ecdd9a5aa65b61ed39c24db2b3488ac",
			"OP_DUP OP_HASH160 cbc222711a230ecdd9a5aa65b61ed39c24db2b34 OP_EQUALVERIFY OP_CHECKSIG",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.hex, func(t *testing.T) {
			script, err := NewScriptFromHex(tc.hex)
			require.NoError(t, err)
			assert.Equal(t, script.Hex, tc.hex)
			assert.Equal(t, script.Asm, tc.asm)
		})
	}
}
