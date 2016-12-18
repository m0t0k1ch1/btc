package btctx

import (
	"encoding/hex"

	"github.com/m0t0k1ch1/base58"
)

type PubKeyHash string

func (pkh PubKeyHash) ToAddress() (string, error) {
	pkhBytes, err := hex.DecodeString(string(pkh))
	if err != nil {
		return "", err
	}

	if isTestNet() {
		pkhBytes = append([]byte{0x6f}, pkhBytes...)
	} else {
		pkhBytes = append([]byte{0x00}, pkhBytes...)
	}

	doubleHashBytes, err := sha256Double(pkhBytes)
	if err != nil {
		return "", err
	}

	checksumBytes := doubleHashBytes[0:4]
	addressBytes := append(pkhBytes, checksumBytes...)

	b58 := base58.NewBitcoinBase58()

	return b58.EncodeToString(addressBytes)
}
