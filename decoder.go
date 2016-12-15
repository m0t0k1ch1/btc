package btctx

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type decoder struct {
	data   []byte
	reader *bytes.Reader
}

func newDecoder(b []byte) *decoder {
	return &decoder{
		data:   b,
		reader: bytes.NewReader(b),
	}
}

func (d *decoder) readByte() (byte, error) {
	return d.reader.ReadByte()
}

func (d *decoder) readBytes(size int) ([]byte, error) {
	data := make([]byte, size)

	for i := 0; i < size; i++ {
		b, err := d.readByte()
		if err != nil {
			return nil, err
		}

		data[i] = b
	}

	return data, nil
}

func (d *decoder) readBytesReverse(size int) ([]byte, error) {
	data := make([]byte, size)

	for i := size - 1; i >= 0; i-- {
		b, err := d.readByte()
		if err != nil {
			return nil, err
		}

		data[i] = b
	}

	return data, nil
}

func (d *decoder) readString(size int) (string, error) {
	b, err := d.readBytes(size)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (d *decoder) readStringReverse(size int) (string, error) {
	b, err := d.readBytesReverse(size)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (d *decoder) readHex(size int) (string, error) {
	s, err := d.readString(size)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", s), nil
}

func (d *decoder) readHexReverse(size int) (string, error) {
	s, err := d.readStringReverse(size)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", s), nil
}

func (d *decoder) readBinary(size int, data interface{}) error {
	b, err := d.readBytes(size)
	if err != nil {
		return err
	}

	buf := bytes.NewReader(b)
	if err := binary.Read(buf, binary.LittleEndian, data); err != nil {
		return err
	}

	return nil
}

func (d *decoder) readInt16() (int16, error) {
	var data int16
	if err := d.readBinary(2, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (d *decoder) readInt32() (int32, error) {
	var data int32
	if err := d.readBinary(4, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (d *decoder) readInt64() (int64, error) {
	var data int64
	if err := d.readBinary(8, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (d *decoder) readUint16() (uint16, error) {
	var data uint16
	if err := d.readBinary(2, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (d *decoder) readUint32() (uint32, error) {
	var data uint32
	if err := d.readBinary(4, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (d *decoder) readUint64() (uint64, error) {
	var data uint64
	if err := d.readBinary(8, &data); err != nil {
		return 0, err
	}

	return data, nil
}

// variable length integer
// ref. https://en.bitcoin.it/wiki/Protocol_documentation#Variable_length_integer
func (d *decoder) readVarInt() (uint, error) {
	head, err := d.readByte()
	if err != nil {
		return 0, err
	}

	switch head {
	case 0xff:
		val, err := d.readUint64()
		if err != nil {
			return 0, err
		}
		return uint(val), nil
	case 0xfe:
		val, err := d.readUint32()
		if err != nil {
			return 0, err
		}
		return uint(val), nil
	case 0xfd:
		val, err := d.readUint16()
		if err != nil {
			return 0, err
		}
		return uint(val), nil
	default:
		return uint(head), nil
	}
}
