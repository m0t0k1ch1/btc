package btctx

import (
	"crypto/sha256"
	"errors"
	"os"
)

const (
	NetworkTypeEnvKey = "BTCTX_NETWORK"
	NetworkTypeMain   = "mainnet"
	NetworkTypeTest   = "testnet"

	AddressVersionMain byte = 0x00
	AddressVersionTest byte = 0x6f

	SatoshiPerBtc = 100000000

	TxVersion    int32  = 1
	TxLockTime   uint32 = 0
	TxInSequence uint32 = 4294967295

	CoinBaseTxid = "0000000000000000000000000000000000000000000000000000000000000000"
)

var (
	ErrInvalidAddressVersion = errors.New("invalid address version")
	ErrInvalidChecksum       = errors.New("invalid checksum")
)

type Satoshi int64

func NewSatoshiFromBtc(btc float64) Satoshi {
	return Satoshi(btc * SatoshiPerBtc)
}

func (satoshi Satoshi) ToInt64() int64 {
	return int64(satoshi)
}

func (satoshi Satoshi) ToBtc() float64 {
	return float64(satoshi / SatoshiPerBtc)
}

func UseTestnet() error {
	return os.Setenv(NetworkTypeEnvKey, NetworkTypeTest)
}

func isTestNet() bool {
	if os.Getenv(NetworkTypeEnvKey) == NetworkTypeTest {
		return true
	}

	return false
}

func sha256Double(b []byte) ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write(b); err != nil {
		return nil, err
	}

	tmp := h.Sum(nil)

	h.Reset()
	if _, err := h.Write(tmp); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func reverseBytes(b []byte) []byte {
	l := len(b)

	rb := make([]byte, l)
	for i, j := 0, l-1; i < l; i++ {
		rb[i] = b[j]
		j--
	}

	return rb
}
