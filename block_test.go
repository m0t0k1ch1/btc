package btc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlockMapping(t *testing.T) {
	testCases := []struct {
		blockhash     string
		hex           string
		version       int32
		prevBlockHash string
		merkleRoot    string
		timestamp     uint32
		bits          uint32
		nonce         uint32
		txids         []string
	}{
		{
			// genesis block of Bitcoin main network
			"000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f",
			"0100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4a29ab5f49ffff001d1dac2b7c0101000000010000000000000000000000000000000000000000000000000000000000000000ffffffff4d04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73ffffffff0100f2052a01000000434104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac00000000",
			1,
			"0000000000000000000000000000000000000000000000000000000000000000",
			"4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b",
			1231006505,
			486604799,
			2083236893,
			[]string{
				"4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.blockhash, func(t *testing.T) {
			block, err := NewBlockFromHex(tc.hex)
			require.NoError(t, err)
			assert.Equal(t, block.Version, tc.version)
			assert.Equal(t, block.PrevBlockhash, tc.prevBlockHash)
			assert.Equal(t, block.MerkleRoot, tc.merkleRoot)
			assert.Equal(t, block.Timestamp, tc.timestamp)
			assert.Equal(t, block.Bits, tc.bits)
			assert.Equal(t, block.Nonce, tc.nonce)

			for idx, tx := range block.Txes {
				txid, err := tx.Txid()
				require.NoError(t, err)
				assert.Equal(t, txid, tc.txids[idx])
			}

			blockHex, err := block.Hex()
			require.NoError(t, err)
			assert.Equal(t, blockHex, tc.hex)
		})
	}
}
