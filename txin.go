package btctx

import (
	"bytes"
	"encoding/hex"
	"io"
)

type TxIn struct {
	Hash            string `json:"hash"`
	Index           uint32 `json:"index"`
	SignatureScript string `json:"signatureScript"` // hex
	Sequence        uint32 `json:"sequence"`
}

func (txIn *TxIn) ToBytes() ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := txIn.write(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (txIn *TxIn) write(w io.Writer) error {
	if err := txIn.writeHash(w); err != nil {
		return err
	}

	if err := txIn.writeIndex(w); err != nil {
		return err
	}

	if err := txIn.writeSignatureScript(w); err != nil {
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

func (txIn *TxIn) writeSignatureScript(w io.Writer) error {
	s := txIn.SignatureScript

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

func (txIn *TxIn) writeSequence(w io.Writer) error {
	return writeData(w, txIn.Sequence)
}
