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
	Txid     string  `json:"txid"`
	Index    uint32  `json:"index"`
	Script   *Script `json:"script"`
	Sequence uint32  `json:"sequence"`
}

func NewTxIn(txid string, index uint32, scriptHex string) (*TxIn, error) {
	script, err := NewScriptFromHex(scriptHex)
	if err != nil {
		return nil, err
	}

	return &TxIn{
		Txid:     txid,
		Index:    index,
		Script:   script,
		Sequence: DefaultTxInSequence,
	}, nil
}

func (txIn *TxIn) Bytes() ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := txIn.writeAll(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (txIn *TxIn) Hex() (string, error) {
	b, err := txIn.Bytes()
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
	return writeHexReverse(w, txIn.Txid)
}

func (txIn *TxIn) writeIndex(w io.Writer) error {
	return writeData(w, txIn.Index)
}

func (txIn *TxIn) writeSigScriptLength(w io.Writer) error {
	b, err := hex.DecodeString(txIn.Script.Hex)
	if err != nil {
		return err
	}

	return writeVarInt(w, uint(len(b)))
}

func (txIn *TxIn) writeSigScript(w io.Writer) error {
	return writeHex(w, txIn.Script.Hex)
}

func (txIn *TxIn) writeSequence(w io.Writer) error {
	return writeData(w, txIn.Sequence)
}
