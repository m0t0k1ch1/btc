package btctx

import (
	"encoding/hex"
	"fmt"
	"io"
	"strings"
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

func (sd *scriptDecoder) decodePushData(op OP) ([]string, error) {
	parts := []string{}
	var len uint

	switch op {
	case OP_PUSHDATA1:
		parts = append(parts, opCodeMap[OP_PUSHDATA1])
		len8, err := sd.readByte()
		if err != nil {
			return nil, err
		}
		len = uint(len8)
	case OP_PUSHDATA2:
		parts = append(parts, opCodeMap[OP_PUSHDATA2])
		len16, err := sd.readUint16()
		if err != nil {
			return nil, err
		}
		len = uint(len16)
	case OP_PUSHDATA4:
		parts = append(parts, opCodeMap[OP_PUSHDATA4])
		len32, err := sd.readUint32()
		if err != nil {
			return nil, err
		}
		len = uint(len32)
	default:
		len = uint(op)
	}

	data, err := sd.readHex(len)
	if err != nil {
		return nil, err
	}

	return append(parts, data), nil
}

func (sd *scriptDecoder) decodePart() ([]string, error) {
	op, err := sd.decodeOP()
	if err != nil {
		return nil, err
	}

	if op.isPushData() {
		return sd.decodePushData(op)
	}

	opCode, ok := opCodeMap[op]
	if !ok {
		return nil, fmt.Errorf("unknown operation code")
	}

	return []string{opCode}, nil
}

func (sd *scriptDecoder) decode() (*Script, error) {
	sps := scriptParts{}

	for {
		parts, err := sd.decodePart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		sps = append(sps, parts...)
	}

	addresses, _ := sps.extractAddresses()

	return &Script{
		Hex:       hex.EncodeToString(sd.data),
		Asm:       strings.Join(sps, " "),
		Addresses: addresses,
	}, nil
}
