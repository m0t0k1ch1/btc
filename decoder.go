package btctx

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
)

type decoder struct {
	src    []byte
	reader *bytes.Reader
}

func newDecoder(b []byte) *decoder {
	return &decoder{
		src:    b,
		reader: bytes.NewReader(b),
	}
}

func (d *decoder) readByte() (byte, error) {
	return d.reader.ReadByte()
}

func (d *decoder) readBytes(size uint) ([]byte, error) {
	data := make([]byte, size)

	var i uint
	for i = 0; i < size; i++ {
		b, err := d.readByte()
		if err != nil {
			return nil, err
		}

		data[i] = b
	}

	return data, nil
}

func (d *decoder) readBytesReverse(size uint) ([]byte, error) {
	data := make([]byte, size)

	var i uint
	for i = 1; i <= size; i++ {
		b, err := d.readByte()
		if err != nil {
			return nil, err
		}

		data[size-i] = b
	}

	return data, nil
}

func (d *decoder) readString(size uint) (string, error) {
	b, err := d.readBytes(size)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (d *decoder) readStringReverse(size uint) (string, error) {
	b, err := d.readBytesReverse(size)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (d *decoder) readHex(size uint) (string, error) {
	b, err := d.readBytes(size)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func (d *decoder) readHexReverse(size uint) (string, error) {
	b, err := d.readBytesReverse(size)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func (d *decoder) readBinary(size uint, data interface{}) error {
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
		data, err := d.readUint64()
		if err != nil {
			return 0, err
		}
		return uint(data), nil
	case 0xfe:
		data, err := d.readUint32()
		if err != nil {
			return 0, err
		}
		return uint(data), nil
	case 0xfd:
		data, err := d.readUint16()
		if err != nil {
			return 0, err
		}
		return uint(data), nil
	default:
		return uint(head), nil
	}
}
