package btctx

import (
	"encoding/hex"
	"errors"
	"io"
	"strings"
)

var (
	ErrUnknownOperationCode = errors.New("unknown operation code")
)

type scriptDecoder struct {
	*decoder
}

func newScriptDecoder(b []byte) *scriptDecoder {
	return &scriptDecoder{newDecoder(b)}
}

func (sd *scriptDecoder) decodeOP() (OP, error) {
	b, err := sd.readByte()
	if err != nil {
		return OP_0, err
	}

	return OP(b), nil
}

func (sd *scriptDecoder) decodePushData(op OP) ([]string, []byte, error) {
	parts := []string{}
	var len uint

	switch op {
	case OP_PUSHDATA1:
		parts = append(parts, opCodeMap[OP_PUSHDATA1])
		len8, err := sd.readByte()
		if err != nil {
			return nil, nil, err
		}
		len = uint(len8)
	case OP_PUSHDATA2:
		parts = append(parts, opCodeMap[OP_PUSHDATA2])
		len16, err := sd.readUint16()
		if err != nil {
			return nil, nil, err
		}
		len = uint(len16)
	case OP_PUSHDATA4:
		parts = append(parts, opCodeMap[OP_PUSHDATA4])
		len32, err := sd.readUint32()
		if err != nil {
			return nil, nil, err
		}
		len = uint(len32)
	default:
		len = uint(op)
	}

	dataBytes, err := sd.readBytes(len)
	if err != nil {
		return nil, nil, err
	}
	parts = append(parts, hex.EncodeToString(dataBytes))

	return parts, dataBytes, nil
}

func (sd *scriptDecoder) decodePart() ([]string, []byte, error) {
	op, err := sd.decodeOP()
	if err != nil {
		return nil, nil, err
	}

	if op.isPushData() {
		return sd.decodePushData(op)
	}

	opCode, ok := opCodeMap[op]
	if !ok {
		return nil, nil, ErrUnknownOperationCode
	}

	return []string{opCode}, nil, nil
}

func (sd *scriptDecoder) decode() (*Script, error) {
	sps := scriptParts{}
	data := [][]byte{}

	for {
		parts, dataBytes, err := sd.decodePart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		sps = append(sps, parts...)

		if dataBytes != nil {
			data = append(data, dataBytes)
		}
	}

	addresses, _ := sps.extractAddresses()

	return &Script{
		Hex:       hex.EncodeToString(sd.data),
		Asm:       strings.Join(sps, " "),
		Addresses: addresses,
		Data:      data,
	}, nil
}
