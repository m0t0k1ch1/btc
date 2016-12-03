package btctx

import (
	"bytes"
	"encoding/hex"
	"io"
)

type TxOut struct {
	Value    int64  `json:"value"`    // satoshi
	PkScript string `json:"pkScript"` // hex
}

func (txOut *TxOut) ToBytes() ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := txOut.write(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (txOut *TxOut) write(w io.Writer) error {
	if err := txOut.writeValue(w); err != nil {
		return err
	}

	if err := txOut.writePkScript(w); err != nil {
		return err
	}

	return nil
}

func (txOut *TxOut) writeValue(w io.Writer) error {
	return writeData(w, txOut.Value)
}

func (txOut *TxOut) writePkScript(w io.Writer) error {
	s := txOut.PkScript

	b, err := hex.DecodeString(s)
	if err != nil {
		return err
	}

	if err := writeVarInt(w, len(b)); err != nil {
		return err
	}
	if err := writeHex(w, s); err != nil {
		return err
	}

	return nil
}
