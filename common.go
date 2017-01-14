package btctx

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"io"
	"os"
)

const (
	SatoshiPerBtc = 100000000

	NetworkEnvKey  = "BTCTX_NETWORK"
	NetworkEnvMain = "mainnet"
	NetworkEnvTest = "testnet"
)

type Satoshi int64

func (satoshi Satoshi) ToInt64() int64 {
	return int64(satoshi)
}

func (satoshi Satoshi) ToBtc() float64 {
	return float64(satoshi) / SatoshiPerBtc
}

func UseTestnet() error {
	return os.Setenv(NetworkEnvKey, NetworkEnvTest)
}

func isTestNet() bool {
	if os.Getenv(NetworkEnvKey) == NetworkEnvTest {
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

func writeData(w io.Writer, data interface{}) error {
	return binary.Write(w, binary.LittleEndian, data)
}

// variable length integer
// ref. https://en.bitcoin.it/wiki/Protocol_documentation#Variable_length_integer
func writeVarInt(w io.Writer, data uint) error {
	if data < 0xfd {
		return writeData(w, byte(data))
	} else if data <= 0xffff {
		if err := writeData(w, byte(0xfd)); err != nil {
			return err
		}
		return writeData(w, uint16(data))
	} else if data <= 0xffffffff {
		if err := writeData(w, byte(0xfe)); err != nil {
			return err
		}
		return writeData(w, uint32(data))
	} else {
		if err := writeData(w, byte(0xff)); err != nil {
			return err
		}
		return writeData(w, uint64(data))
	}
}

func writeHex(w io.Writer, data string) error {
	b, err := hex.DecodeString(data)
	if err != nil {
		return err
	}

	return writeData(w, b)
}

func writeHexReverse(w io.Writer, data string) error {
	b, err := hex.DecodeString(data)
	if err != nil {
		return err
	}
	size := len(b)

	reversed := make([]byte, size)

	for i, j := size-1, 0; i >= 0; i-- {
		reversed[j] = b[i]
		j++
	}

	return writeData(w, reversed)
}
