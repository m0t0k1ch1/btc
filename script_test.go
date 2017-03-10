package btctx

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type scriptMappingTestCase struct {
	hex     string
	asm     string
	pkh     string
	address string
}

var (
	scriptMappingTestCases = []scriptMappingTestCase{
		{
			"76a914cbc222711a230ecdd9a5aa65b61ed39c24db2b3488ac",
			"OP_DUP OP_HASH160 cbc222711a230ecdd9a5aa65b61ed39c24db2b34 OP_EQUALVERIFY OP_CHECKSIG",
			"cbc222711a230ecdd9a5aa65b61ed39c24db2b34",
			"mz6L2hYM8jPR5nhH6kEsc3DQFiSDA1Jqpa",
		},
	}
)

func TestScriptMapping(t *testing.T) {
	UseTestnet()

	for _, testCase := range scriptMappingTestCases {
		pubKeyHashBytes, err := hex.DecodeString(testCase.pkh)
		require.NoError(t, err)

		script, err := NewScriptFromHex(testCase.hex)
		require.NoError(t, err)
		assert.Equal(t, script.Hex, testCase.hex)
		assert.Equal(t, script.Asm, testCase.asm)
		assert.Equal(t, len(script.Addresses), 1)
		assert.Equal(t, script.Addresses, []string{testCase.address})
		assert.Equal(t, len(script.Data), 1)
		assert.Equal(t, script.Data, [][]byte{pubKeyHashBytes})
	}
}

// ==================================================
// target script in testnet
// ==================================================
// hex:
// 76a914cbc222711a230ecdd9a5aa65b61ed39c24db2b3488ac
// --------------------------------------------------
// asm:
// OP_DUP OP_HASH160 cbc222711a230ecdd9a5aa65b61ed39c24db2b34 OP_EQUALVERIFY OP_CHECKSIG
// --------------------------------------------------
