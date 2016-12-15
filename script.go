package btctx

import "encoding/hex"

const (
	OP_0           OP = 0x00
	OP_FALSE       OP = 0x00
	OP_PUSHDATA1   OP = 0x4c
	OP_PUSHDATA2   OP = 0x4d
	OP_PUSHDATA4   OP = 0x4e
	OP_DUP         OP = 0x76
	OP_EQUALVERIFY OP = 0x88
	OP_HASH160     OP = 0xa9
	OP_CHECKSIG    OP = 0xac
)

var opCodeMap = map[OP]string{
	OP_DUP:         "OP_DUP",
	OP_EQUALVERIFY: "OP_EQUALVERIFY",
	OP_HASH160:     "OP_HASH160",
	OP_CHECKSIG:    "OP_CHECKSIG",
}

type OP byte

func (op OP) isPushData() bool {
	return (0x01 <= op && op <= 0x75) ||
		op == OP_PUSHDATA1 ||
		op == OP_PUSHDATA2 ||
		op == OP_PUSHDATA4
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
