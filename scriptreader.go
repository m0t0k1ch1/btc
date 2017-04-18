package btctx

import (
	"encoding/hex"
	"io"
	"strings"
)

type scriptReader struct {
	*reader
}

func newScriptReader(b []byte) *scriptReader {
	return &scriptReader{newReader(b)}
}

func (r *scriptReader) readOpCode() (OpCode, error) {
	b, err := r.ReadByte()
	if err != nil {
		return OpZero, err
	}

	return OpCode(b), nil
}

func (r *scriptReader) readPushedData(op OpCode) ([]byte, error) {
	var len uint

	switch op {
	case OpPushdata1:
		len8, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		len = uint(len8)
	case OpPushdata2:
		len16, err := r.readUint16()
		if err != nil {
			return nil, err
		}
		len = uint(len16)
	case OpPushdata4:
		len32, err := r.readUint32()
		if err != nil {
			return nil, err
		}
		len = uint(len32)
	default:
		// op: 0x01-0x4b (1-75)
		len = uint(op)
	}

	return r.readBytes(len)
}

func (r *scriptReader) readPart() ([]string, error) {
	op, err := r.readOpCode()
	if err != nil {
		return nil, err
	}

	if op.isPushData() {
		b, err := r.readPushedData(op)
		if err != nil {
			return nil, err
		}

		return []string{hex.EncodeToString(b)}, nil
	}

	return []string{op.Name()}, nil
}

func (r *scriptReader) readScript() (*Script, error) {
	b := r.Bytes()
	asmParts := []string{}

	for {
		parts, err := r.readPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		asmParts = append(asmParts, parts...)
	}

	return &Script{
		Hex: hex.EncodeToString(b),
		Asm: strings.Join(asmParts, " "),
	}, nil
}
