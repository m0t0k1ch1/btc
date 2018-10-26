package btc

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPkhAddressConversion(t *testing.T) {
	testCases := []struct {
		pkh     string
		address string
	}{
		{
			// ref. https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses
			"010966776006953d5567439e5e39f86a0d273bee",
			"16UwLL9Risc3QfPqBUvKofHmBQ7wMtjvM",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.pkh, func(t *testing.T) {
			b, err := hex.DecodeString(tc.pkh)
			require.NoError(t, err)

			// address -> pkh
			address, err := Pkh(b).Address()
			require.NoError(t, err)
			assert.Equal(t, address.String(), tc.address)

			// pkh -> address
			pkh, err := address.Pkh()
			require.NoError(t, err)
			assert.Equal(t, hex.EncodeToString(pkh.Bytes()), tc.pkh)
		})
	}
}

func TestAddressValidation(t *testing.T) {
	testCases := []struct {
		address Address
		isValid bool
	}{
		{
			Address("16UwLL9Risc3QfPqBUvKofHmBQ7wMtjvM"),
			true,
		},
		{
			Address("000000000000000000000000000000000"),
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.address.String(), func(t *testing.T) {
			ok, err := tc.address.IsValid()
			require.NoError(t, err)
			assert.Equal(t, ok, tc.isValid)
		})
	}
}
