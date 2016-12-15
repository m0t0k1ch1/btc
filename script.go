package btctx

import "encoding/hex"

const (
	OP_PUSHDATA1   byte = 76
	OP_PUSHDATA2   byte = 77
	OP_PUSHDATA4   byte = 78
	OP_DUP         byte = 118
	OP_EQUALVERIFY byte = 136
	OP_HASH160     byte = 169
	OP_CHECKSIG    byte = 172
)

var opCodeMap = map[byte]string{
	OP_DUP:         "OP_DUP",
	OP_EQUALVERIFY: "OP_EQUALVERIFY",
	OP_HASH160:     "OP_HASH160",
	OP_CHECKSIG:    "OP_CHECKSIG",
}

type Script struct {
	Hex string `json:"hex"`
	Asm string `json:"asm"`
}

func NewScriptFromHex(s string) (*Script, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return NewScriptFromBytes(b)
}

func NewScriptFromBytes(b []byte) (*Script, error) {
	sd := newScriptDecoder(b)

	return sd.decode()
}
