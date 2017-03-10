package btctx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type scriptMappingTestCase struct {
	hex       string
	asm       string
	addresses []string
	data      [][]byte
}

var (
	scriptMappingTestCases = []scriptMappingTestCase{
		{
			"",
			"",
			[]string{},
			[][]byte{},
		},
		{
			"76a914cbc222711a230ecdd9a5aa65b61ed39c24db2b3488ac",
			"OP_DUP OP_HASH160 cbc222711a230ecdd9a5aa65b61ed39c24db2b34 OP_EQUALVERIFY OP_CHECKSIG",
			[]string{"mz6L2hYM8jPR5nhH6kEsc3DQFiSDA1Jqpa"},
			[][]byte{[]byte{203, 194, 34, 113, 26, 35, 14, 205, 217, 165, 170, 101, 182, 30, 211, 156, 36, 219, 43, 52}},
		},
	}
)

func TestScriptMapping(t *testing.T) {
	UseTestnet()

	for _, testCase := range scriptMappingTestCases {
		script, err := NewScriptFromHex(testCase.hex)
		require.NoError(t, err)
		assert.Equal(t, script.Hex, testCase.hex)
		assert.Equal(t, script.Asm, testCase.asm)
		assert.Equal(t, len(script.Addresses), len(testCase.addresses))
		assert.Equal(t, script.Addresses, testCase.addresses)
		assert.Equal(t, len(script.Data), len(testCase.data))
		assert.Equal(t, script.Data, testCase.data)
	}
}
