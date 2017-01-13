package btctx

import (
	"bytes"
	"encoding/hex"
	"errors"

	"github.com/m0t0k1ch1/base58"
)

type PubKeyHash string

const (
	AddressVersionMain byte = 0x00
	AddressVersionTest byte = 0x6f
)

var (
	ErrInvalidAddressVersion = errors.New("invalid address version")
	ErrInvalidChecksum       = errors.New("invalid checksum")
)

func IsValidAddressVersion(version byte) bool {
	if isTestNet() {
		if version == AddressVersionTest {
			return true
		}
	} else {
		if version == AddressVersionMain {
			return true
		}
	}

	return false
}

func NewPkhFromAddress(address string) (PubKeyHash, error) {
	b58 := base58.NewBitcoinBase58()

	pkhBytes, err := b58.DecodeString(address)
	if err != nil {
		return "", err
	}

	// validate address version byte
	if !IsValidAddressVersion(pkhBytes[0]) {
		return "", ErrInvalidAddressVersion
	}

	checksumBytes := pkhBytes[21:]
	pkhBytes = pkhBytes[0:21]

	doubleHashBytes, err := sha256Double(pkhBytes)
	if err != nil {
		return "", err
	}

	// validate checksum bytes
	checksumBytesValid := doubleHashBytes[0:4]
	if bytes.Compare(checksumBytes, checksumBytesValid) != 0 {
		return "", ErrInvalidChecksum
	}

	return PubKeyHash(hex.EncodeToString(pkhBytes[1:21])), nil
}

func (pkh PubKeyHash) ToString() string {
	return string(pkh)
}

func (pkh PubKeyHash) ToAddress() (string, error) {
	pkhBytes, err := hex.DecodeString(pkh.ToString())
	if err != nil {
		return "", err
	}

	if isTestNet() {
		pkhBytes = append([]byte{AddressVersionTest}, pkhBytes...)
	} else {
		pkhBytes = append([]byte{AddressVersionMain}, pkhBytes...)
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
