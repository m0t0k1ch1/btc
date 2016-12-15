package btctx

import (
	"bytes"
	"encoding/hex"
	"io"
)

const (
	DefaultTxInSequence uint32 = 4294967295
)

type TxIn struct {
	Hash      string  `json:"hash"`
	Index     uint32  `json:"index"`
	SigScript *Script `json:"sigScript"`
	Sequence  uint32  `json:"sequence"`
}

func NewTxIn() *TxIn {
	return &TxIn{
		Sequence: DefaultTxInSequence,
	}
}

func (txIn *TxIn) ToBytes() ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := txIn.WriteAll(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (txIn *TxIn) WriteAll(w io.Writer) error {
	if err := txIn.WriteHash(w); err != nil {
		return err
	}

	if err := txIn.WriteIndex(w); err != nil {
		return err
	}

	if err := txIn.WriteSigScriptLength(w); err != nil {
		return err
	}

	if err := txIn.WriteSigScript(w); err != nil {
		return err
	}

	if err := txIn.WriteSequence(w); err != nil {
		return err
	}

	return nil
}

func (txIn *TxIn) WriteHash(w io.Writer) error {
	return writeHexReverse(w, txIn.Hash)
}

func (txIn *TxIn) WriteIndex(w io.Writer) error {
	return writeData(w, txIn.Index)
}

func (txIn *TxIn) WriteSigScriptLength(w io.Writer) error {
	b, err := hex.DecodeString(txIn.SigScript.Hex)
	if err != nil {
		return err
	}

	return writeVarInt(w, len(b))
}

func (txIn *TxIn) WriteSigScript(w io.Writer) error {
	return writeHex(w, txIn.SigScript.Hex)
}

func (txIn *TxIn) WriteSequence(w io.Writer) error {
	return writeData(w, txIn.Sequence)
}
