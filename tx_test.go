package btctx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type txMappingTestCase struct {
	txid     string
	hex      string
	version  int32
	txIns    []*TxIn
	txOuts   []*TxOut
	lockTime uint32
}

var (
	txMappingTestCases = []txMappingTestCase{
		{
			"d7a4684b71776c8c96edd670a9d0c61d03c293f4c6266b70ff5030b2c4f0bdfe",
			"0100000001ce3cf2e2b334e7e9fa84619469d9edc49368c2f752ea30fb48b080fc794f6d56010000006a473044022065fe1ea4e94a9b44fb62c2b874b63a947504273a60b99b8f7bbf77b4db9331b002205559d8ee93cf341d75866f9eb912af05904fb6eed7372a837308c4e37f3ab58f012103bae5f04799c40862358560e42e441c3080b997a3dec161dd40395e992362bfc9feffffff0200f2052a010000001976a914cbc222711a230ecdd9a5aa65b61ed39c24db2b3488acc08d931a1d0000001976a914426c1ad9fa94f9ea3e6f9248b8bff6768e3ac8c488ac951a1000",
			1,
			[]*TxIn{&TxIn{
				Hash:  "566d4f79fc80b048fb30ea52f7c26893c4edd969946184fae9e734b3e2f23cce",
				Index: 1,
				SigScript: &Script{
					Hex: "473044022065fe1ea4e94a9b44fb62c2b874b63a947504273a60b99b8f7bbf77b4db9331b002205559d8ee93cf341d75866f9eb912af05904fb6eed7372a837308c4e37f3ab58f012103bae5f04799c40862358560e42e441c3080b997a3dec161dd40395e992362bfc9",
				},
				Sequence: 4294967294,
			}},
			[]*TxOut{
				&TxOut{
					Value: Satoshi(5000000000),
					PkScript: &Script{
						Hex: "76a914cbc222711a230ecdd9a5aa65b61ed39c24db2b3488ac",
					},
				},
				&TxOut{
					Value: Satoshi(124999929280),
					PkScript: &Script{
						Hex: "76a914426c1ad9fa94f9ea3e6f9248b8bff6768e3ac8c488ac",
					},
				},
			},
			1055381,
		},
	}
)

func TestTxMapping(t *testing.T) {
	UseTestnet()

	for _, testCase := range txMappingTestCases {
		tx, err := NewTxFromHex(testCase.hex)
		require.NoError(t, err)
		assert.Equal(t, tx.Version, testCase.version)
		require.Equal(t, len(tx.TxIns), len(testCase.txIns))
		require.Equal(t, len(tx.TxOuts), len(testCase.txOuts))
		assert.Equal(t, tx.LockTime, testCase.lockTime)

		for idx, txIn := range tx.TxIns {
			assert.Equal(t, txIn.Hash, testCase.txIns[idx].Hash)
			assert.Equal(t, txIn.Index, testCase.txIns[idx].Index)
			assert.Equal(t, txIn.SigScript.Hex, testCase.txIns[idx].SigScript.Hex)
			assert.Equal(t, txIn.Sequence, testCase.txIns[idx].Sequence)
		}

		for idx, txOut := range tx.TxOuts {
			assert.Equal(t, txOut.Value, testCase.txOuts[idx].Value)
			assert.Equal(t, txOut.PkScript.Hex, testCase.txOuts[idx].PkScript.Hex)
		}

		txid, err := tx.ToHash()
		require.NoError(t, err)
		assert.Equal(t, txid, testCase.txid)

		txHex, err := tx.ToHex()
		require.NoError(t, err)
		assert.Equal(t, txHex, testCase.hex)
	}
}

// ==================================================
// target transaction in testnet
// ==================================================
// txid:
// d7a4684b71776c8c96edd670a9d0c61d03c293f4c6266b70ff5030b2c4f0bdfe
// --------------------------------------------------
// raw:
// 0100000001ce3cf2e2b334e7e9fa84619469d9edc49368c2f752ea30fb48b080fc794f6d56010000006a473044022065fe1ea4e94a9b44fb62c2b874b63a947504273a60b99b8f7bbf77b4db9331b002205559d8ee93cf341d75866f9eb912af05904fb6eed7372a837308c4e37f3ab58f012103bae5f04799c40862358560e42e441c3080b997a3dec161dd40395e992362bfc9feffffff0200f2052a010000001976a914cbc222711a230ecdd9a5aa65b61ed39c24db2b3488acc08d931a1d0000001976a914426c1ad9fa94f9ea3e6f9248b8bff6768e3ac8c488ac951a1000
// --------------------------------------------------
// json:
// {
//   "result": {
//     "txid": "d7a4684b71776c8c96edd670a9d0c61d03c293f4c6266b70ff5030b2c4f0bdfe",
//     "hash": "d7a4684b71776c8c96edd670a9d0c61d03c293f4c6266b70ff5030b2c4f0bdfe",
//     "size": 225,
//     "vsize": 225,
//     "version": 1,
//     "locktime": 1055381,
//     "vin": [
//       {
//         "txid": "566d4f79fc80b048fb30ea52f7c26893c4edd969946184fae9e734b3e2f23cce",
//         "vout": 1,
//         "scriptSig": {
//           "asm": "3044022065fe1ea4e94a9b44fb62c2b874b63a947504273a60b99b8f7bbf77b4db9331b002205559d8ee93cf341d75866f9eb912af05904fb6eed7372a837308c4e37f3ab58f[ALL] 03bae5f04799c40862358560e42e441c3080b997a3dec161dd40395e992362bfc9",
//           "hex": "473044022065fe1ea4e94a9b44fb62c2b874b63a947504273a60b99b8f7bbf77b4db9331b002205559d8ee93cf341d75866f9eb912af05904fb6eed7372a837308c4e37f3ab58f012103bae5f04799c40862358560e42e441c3080b997a3dec161dd40395e992362bfc9"
//         },
//         "sequence": 4294967294
//       }
//     ],
//     "vout": [
//       {
//         "value": 50,
//         "n": 0,
//         "scriptPubKey": {
//           "asm": "OP_DUP OP_HASH160 cbc222711a230ecdd9a5aa65b61ed39c24db2b34 OP_EQUALVERIFY OP_CHECKSIG",
//           "hex": "76a914cbc222711a230ecdd9a5aa65b61ed39c24db2b3488ac",
//           "reqSigs": 1,
//           "type": "pubkeyhash",
//           "addresses": [
//             "mz6L2hYM8jPR5nhH6kEsc3DQFiSDA1Jqpa"
//           ]
//         }
//       },
//       {
//         "value": 1249.9992928,
//         "n": 1,
//         "scriptPubKey": {
//           "asm": "OP_DUP OP_HASH160 426c1ad9fa94f9ea3e6f9248b8bff6768e3ac8c4 OP_EQUALVERIFY OP_CHECKSIG",
//           "hex": "76a914426c1ad9fa94f9ea3e6f9248b8bff6768e3ac8c488ac",
//           "reqSigs": 1,
//           "type": "pubkeyhash",
//           "addresses": [
//             "mmaAPyTMoK3p1K7qCYytuSRRQMPCbrUAMN"
//           ]
//         }
//       }
//     ]
//   },
//   "error": null,
//   "id": 0
// }
// --------------------------------------------------
