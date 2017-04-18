package btctx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type scriptMappingTestCase struct {
	hex string
	asm string
}

var (
	scriptMappingTestCases = []scriptMappingTestCase{
		{
			"",
			"",
		},
		{
			"76a914cbc222711a230ecdd9a5aa65b61ed39c24db2b3488ac",
			"OP_DUP OP_HASH160 cbc222711a230ecdd9a5aa65b61ed39c24db2b34 OP_EQUALVERIFY OP_CHECKSIG",
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
	}
}
