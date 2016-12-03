package btctx

import (
	"encoding/binary"
	"encoding/hex"
	"io"
)

func writeData(w io.Writer, data interface{}) error {
	return binary.Write(w, binary.LittleEndian, data)
}

// variable length integer
// ref. https://en.bitcoin.it/wiki/Protocol_documentation#Variable_length_integer
func writeVarInt(w io.Writer, data int) error {
	if data < 0xfd {
		return writeData(w, uint8(data))
	} else if data <= 0xffff {
		return writeData(w, uint16(data))
	} else if data <= 0xffffffff {
		return writeData(w, uint32(data))
	} else {
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