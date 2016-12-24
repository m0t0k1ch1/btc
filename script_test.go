package btctx

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScriptMapping(t *testing.T) {
	UseTestnet()

	scriptHex := "76a914cbc222711a230ecdd9a5aa65b61ed39c24db2b3488ac"

	pubKeyHash := "cbc222711a230ecdd9a5aa65b61ed39c24db2b34"
	pubKeyHashBytes, err := hex.DecodeString(pubKeyHash)
	require.NoError(t, err)

	script, err := NewScriptFromHex(scriptHex)
	require.NoError(t, err)
	assert.Equal(t, script.Hex, scriptHex)
	assert.Equal(t, script.Asm, fmt.Sprintf("OP_DUP OP_HASH160 %s OP_EQUALVERIFY OP_CHECKSIG", pubKeyHash))
	assert.Equal(t, len(script.Addresses), 1)
	assert.Equal(t, script.Addresses, []string{"mz6L2hYM8jPR5nhH6kEsc3DQFiSDA1Jqpa"})
	assert.Equal(t, len(script.Data), 1)
	assert.Equal(t, script.Data, [][]byte{pubKeyHashBytes})
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
