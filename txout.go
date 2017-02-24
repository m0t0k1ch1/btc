package btctx

import (
	"bytes"
	"encoding/hex"
	"io"
)

type TxOut struct {
	Value    Satoshi `json:"value"`
	PkScript *Script `json:"pkScript"`
}

func NewTxOut(value int64, pkScriptHex string) (*TxOut, error) {
	pkScript, err := NewScriptFromHex(pkScriptHex)
	if err != nil {
		return nil, err
	}

	return &TxOut{
		Value:    Satoshi(value),
		PkScript: pkScript,
	}, nil
}

func (txOut *TxOut) ToBytes() ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := txOut.writeAll(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (txOut *TxOut) ToHex() (string, error) {
	b, err := txOut.ToBytes()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func (txOut *TxOut) writeAll(w io.Writer) error {
	if err := txOut.writeValue(w); err != nil {
		return err
	}

	if err := txOut.writePkScriptLength(w); err != nil {
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

func (txOut *TxOut) writePkScriptLength(w io.Writer) error {
	b, err := hex.DecodeString(txOut.PkScript.Hex)
	if err != nil {
		return err
	}

	return writeVarInt(w, uint(len(b)))
}

func (txOut *TxOut) writePkScript(w io.Writer) error {
	return writeHex(w, txOut.PkScript.Hex)
}
