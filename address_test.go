package btc

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type pkhAddressConversionTestCase struct {
	pkh     string
	address string
}

var (
	pkhAddressConversionTestCases = []pkhAddressConversionTestCase{
		{ // ref. https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses
			"010966776006953d5567439e5e39f86a0d273bee",
			"16UwLL9Risc3QfPqBUvKofHmBQ7wMtjvM",
		},
	}
)

func TestPkhAddressConversion(t *testing.T) {
	for _, testCase := range pkhAddressConversionTestCases {
		b, err := hex.DecodeString(testCase.pkh)
		require.NoError(t, err)

		// address -> pkh
		address, err := Pkh(b).Address()
		require.NoError(t, err)
		assert.Equal(t, address.String(), testCase.address)

		// pkh -> address
		pkh, err := address.Pkh()
		require.NoError(t, err)
		assert.Equal(t, hex.EncodeToString(pkh.Bytes()), testCase.pkh)
	}
}
