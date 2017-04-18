package btctx

import (
	"bytes"
	"encoding/hex"
	"io"
)

type TxOut struct {
	Amount Satoshi `json:"amount"`
	Script *Script `json:"script"`
}

func NewTxOut(amount int64, scriptHex string) (*TxOut, error) {
	script, err := NewScriptFromHex(scriptHex)
	if err != nil {
		return nil, err
	}

	return &TxOut{
		Amount: Satoshi(amount),
		Script: script,
	}, nil
}

func (txOut *TxOut) Bytes() ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := txOut.writeAll(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (txOut *TxOut) Hex() (string, error) {
	b, err := txOut.Bytes()
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
	return writeData(w, txOut.Amount)
}

func (txOut *TxOut) writePkScriptLength(w io.Writer) error {
	b, err := hex.DecodeString(txOut.Script.Hex)
	if err != nil {
		return err
	}

	return writeVarInt(w, uint(len(b)))
}

func (txOut *TxOut) writePkScript(w io.Writer) error {
	return writeHex(w, txOut.Script.Hex)
}
