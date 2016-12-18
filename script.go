package btctx

import (
	"encoding/hex"

	"github.com/m0t0k1ch1/base58"
)

type Script struct {
	Hex       string   `json:"hex"`
	Asm       string   `json:"asm"`
	Addresses []string `json:"addresses"`
}

func NewScriptFromHex(s string) (*Script, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return NewScriptFromBytes(b)
}

func NewScriptFromBytes(b []byte) (*Script, error) {
	sd := newScriptDecoder(b)

	return sd.decode()
}

type scriptParts []string

func (sps scriptParts) isP2PKH() bool {
	if sps[0] == opCodeMap[OP_DUP] &&
		sps[1] == opCodeMap[OP_HASH160] &&
		len(sps[2]) == 40 &&
		sps[3] == opCodeMap[OP_EQUALVERIFY] &&
		sps[4] == opCodeMap[OP_CHECKSIG] {
		return true
	}

	return false
}

func (sps scriptParts) extractAddresses() ([]string, error) {
	if sps.isP2PKH() {
		pkHashBytes, err := hex.DecodeString(sps[2])
		if err != nil {
			return nil, err
		}

		if isTestNet() {
			pkHashBytes = append([]byte{0x6f}, pkHashBytes...)
		} else {
			pkHashBytes = append([]byte{0x00}, pkHashBytes...)
		}

		doubleHashBytes, err := sha256Double(pkHashBytes)
		if err != nil {
			return nil, err
		}

		checksumBytes := doubleHashBytes[0:4]
		addressBytes := append(pkHashBytes, checksumBytes...)

		b58 := base58.NewBitcoinBase58()
		address, err := b58.EncodeToString(addressBytes)
		if err != nil {
			return nil, err
		}

		return []string{address}, nil
	}

	return nil, nil
}
