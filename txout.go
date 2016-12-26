package btctx

import (
	"bytes"
	"encoding/hex"
	"io"
)

type TxOut struct {
	Value    int64   `json:"value"`
	PkScript *Script `json:"pkScript"`
}

func NewTxOut(value int64, pkScriptHex string) (*TxOut, error) {
	pkScript, err := NewScriptFromHex(pkScriptHex)
	if err != nil {
		return nil, err
	}

	return &TxOut{
		Value:    value,
		PkScript: pkScript,
	}, nil
}

func (txOut *TxOut) ToBytes() ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := txOut.WriteAll(buf); err != nil {
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

func (txOut *TxOut) WriteAll(w io.Writer) error {
	if err := txOut.WriteValue(w); err != nil {
		return err
	}

	if err := txOut.WritePkScriptLength(w); err != nil {
		return err
	}

	if err := txOut.WritePkScript(w); err != nil {
		return err
	}

	return nil
}

func (txOut *TxOut) WriteValue(w io.Writer) error {
	return writeData(w, txOut.Value)
}

func (txOut *TxOut) WritePkScriptLength(w io.Writer) error {
	b, err := hex.DecodeString(txOut.PkScript.Hex)
	if err != nil {
		return err
	}

	return writeVarInt(w, uint(len(b)))
}

func (txOut *TxOut) WritePkScript(w io.Writer) error {
	return writeHex(w, txOut.PkScript.Hex)
}
