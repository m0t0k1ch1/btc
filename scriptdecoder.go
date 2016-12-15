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

func (sd *scriptDecoder) decodePushData(op OP) (string, error) {
	var len uint

	switch op {
	case OP_PUSHDATA1:
		len8, err := sd.readByte()
		if err != nil {
			return "", err
		}
		len = uint(len8)
	case OP_PUSHDATA2:
		len16, err := sd.readUint16()
		if err != nil {
			return "", err
		}
		len = uint(len16)
	case OP_PUSHDATA4:
		len32, err := sd.readUint32()
		if err != nil {
			return "", err
		}
		len = uint(len32)
	default:
		len = uint(op)
	}

	return sd.readHex(len)
}

func (sd *scriptDecoder) decodePart() (string, error) {
	op, err := sd.decodeOP()
	if err != nil {
		return "", err
	}

	if op.isPushData() {
		return sd.decodePushData(op)
	}

	opCode, ok := opCodeMap[op]
	if !ok {
		return "", fmt.Errorf("unknown operation code")
	}

	return opCode, nil
}

func (sd *scriptDecoder) decode() (*Script, error) {
	parts := []string{}

	for {
		part, err := sd.decodePart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		parts = append(parts, part)
	}

	return &Script{
		Hex: hex.EncodeToString(sd.data),
		Asm: strings.Join(parts, " "),
	}, nil
}
