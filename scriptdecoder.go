package btctx

import "encoding/hex"

type scriptDecoder struct {
	*decoder
}

func newScriptDecoder(b []byte) *scriptDecoder {
	return &scriptDecoder{newDecoder(b)}
}

func (sd *scriptDecoder) decode() (*Script, error) {
	return &Script{
		Hex: hex.EncodeToString(sd.data),
		Asm: "", // TODO
	}, nil
}
