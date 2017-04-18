package btctx

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
)

type writer struct {
	*bytes.Buffer
}

func newWriter() *writer {
	return &writer{&bytes.Buffer{}}
}

func (w *writer) writeData(data interface{}) error {
	return binary.Write(w, binary.LittleEndian, data)
}

func (w *writer) writeHex(data string) error {
	b, err := hex.DecodeString(data)
	if err != nil {
		return err
	}

	if _, err := w.Write(b); err != nil {
		return err
	}

	return nil
}

func (w *writer) writeHexReverse(data string) error {
	b, err := hex.DecodeString(data)
	if err != nil {
		return err
	}

	if _, err := w.Write(reverseBytes(b)); err != nil {
		return err
	}

	return nil
}

// variable length integer
// ref. https://en.bitcoin.it/wiki/Protocol_documentation#Variable_length_integer
func (w *writer) writeVarInt(data uint) error {
	if data < 0xfd {
		return w.writeData(byte(data))
	} else if data <= 0xffff {
		if err := w.WriteByte(0xfd); err != nil {
			return err
		}
		return w.writeData(uint16(data))
	} else if data <= 0xffffffff {
		if err := w.WriteByte(0xfe); err != nil {
			return err
		}
		return w.writeData(uint32(data))
	} else {
		if err := w.WriteByte(0xff); err != nil {
			return err
		}
		return w.writeData(uint64(data))
	}
}
