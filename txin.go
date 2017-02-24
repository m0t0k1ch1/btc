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

func NewTxIn(hash string, index uint32, sigScriptHex string) (*TxIn, error) {
	sigScript, err := NewScriptFromHex(sigScriptHex)
	if err != nil {
		return nil, err
	}

	return &TxIn{
		Hash:      hash,
		Index:     index,
		SigScript: sigScript,
		Sequence:  DefaultTxInSequence,
	}, nil
}

func (txIn *TxIn) ToBytes() ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := txIn.writeAll(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (txIn *TxIn) ToHex() (string, error) {
	b, err := txIn.ToBytes()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func (txIn *TxIn) writeAll(w io.Writer) error {
	if err := txIn.writeHash(w); err != nil {
		return err
	}

	if err := txIn.writeIndex(w); err != nil {
		return err
	}

	if err := txIn.writeSigScriptLength(w); err != nil {
		return err
	}

	if err := txIn.writeSigScript(w); err != nil {
		return err
	}

	if err := txIn.writeSequence(w); err != nil {
		return err
	}

	return nil
}

func (txIn *TxIn) writeHash(w io.Writer) error {
	return writeHexReverse(w, txIn.Hash)
}

func (txIn *TxIn) writeIndex(w io.Writer) error {
	return writeData(w, txIn.Index)
}

func (txIn *TxIn) writeSigScriptLength(w io.Writer) error {
	b, err := hex.DecodeString(txIn.SigScript.Hex)
	if err != nil {
		return err
	}

	return writeVarInt(w, uint(len(b)))
}

func (txIn *TxIn) writeSigScript(w io.Writer) error {
	return writeHex(w, txIn.SigScript.Hex)
}

func (txIn *TxIn) writeSequence(w io.Writer) error {
	return writeData(w, txIn.Sequence)
}
