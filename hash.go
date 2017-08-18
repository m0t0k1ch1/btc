package btc

import (
	"crypto/sha256"

	"golang.org/x/crypto/ripemd160"
)

func Hash160(b []byte) ([]byte, error) {
	h256 := sha256.New()
	if _, err := h256.Write(b); err != nil {
		return nil, err
	}

	tmp := h256.Sum(nil)

	h160 := ripemd160.New()
	if _, err := h160.Write(tmp); err != nil {
		return nil, err
	}

	return h160.Sum(nil), nil
}

func Sha256Double(b []byte) ([]byte, error) {
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
