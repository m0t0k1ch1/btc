package btc

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
)

type reader struct {
	*bytes.Buffer
}

func newReader(b []byte) *reader {
	return &reader{bytes.NewBuffer(b)}
}

func (r *reader) readData(size uint, data interface{}) error {
	b, err := r.readBytes(size)
	if err != nil {
		return err
	}

	return binary.Read(bytes.NewReader(b), binary.LittleEndian, data)
}

func (r *reader) readBytes(size uint) ([]byte, error) {
	data := make([]byte, size)

	var i uint
	for i = 0; i < size; i++ {
		c, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		data[i] = c
	}

	return data, nil
}

func (r *reader) readBytesReverse(size uint) ([]byte, error) {
	data := make([]byte, size)

	var i uint
	for i = 1; i <= size; i++ {
		c, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		data[size-i] = c
	}

	return data, nil
}

func (r *reader) readString(size uint) (string, error) {
	b, err := r.readBytes(size)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (r *reader) readStringReverse(size uint) (string, error) {
	b, err := r.readBytesReverse(size)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (r *reader) readHex(size uint) (string, error) {
	b, err := r.readBytes(size)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func (r *reader) readHexReverse(size uint) (string, error) {
	b, err := r.readBytesReverse(size)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func (r *reader) readInt16() (int16, error) {
	var data int16
	if err := r.readData(2, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (r *reader) readInt32() (int32, error) {
	var data int32
	if err := r.readData(4, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (r *reader) readInt64() (int64, error) {
	var data int64
	if err := r.readData(8, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (r *reader) readUint16() (uint16, error) {
	var data uint16
	if err := r.readData(2, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (r *reader) readUint32() (uint32, error) {
	var data uint32
	if err := r.readData(4, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (r *reader) readUint64() (uint64, error) {
	var data uint64
	if err := r.readData(8, &data); err != nil {
		return 0, err
	}

	return data, nil
}

// variable length integer
// ref. https://en.bitcoin.it/wiki/Protocol_documentation#Variable_length_integer
func (r *reader) readVarInt() (uint, error) {
	head, err := r.ReadByte()
	if err != nil {
		return 0, err
	}

	switch head {
	case 0xff:
		data, err := r.readUint64()
		if err != nil {
			return 0, err
		}
		return uint(data), nil
	case 0xfe:
		data, err := r.readUint32()
		if err != nil {
			return 0, err
		}
		return uint(data), nil
	case 0xfd:
		data, err := r.readUint16()
		if err != nil {
			return 0, err
		}
		return uint(data), nil
	default:
		return uint(head), nil
	}
}
