package btctx

import "encoding/hex"

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
