package btctx

import (
	"bytes"
	"io"
)

type Tx struct {
	Version  int32    `json:"version"`
	TxIns    []*TxIn  `json:"txIns"`
	TxOuts   []*TxOut `json:"txOuts"`
	LockTime uint32   `json:"lockTime"`
}

func NewTxFromBytes(b []byte) (*Tx, error) {
	txr := newTxReader(b)

	return txr.read()
}

func (tx *Tx) ToBytes() ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := tx.Write(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (tx *Tx) Write(w io.Writer) error {
	if err := tx.WriteVersion(w); err != nil {
		return err
	}

	if err := tx.WriteTxInCount(w); err != nil {
		return err
	}

	if err := tx.WriteTxIns(w); err != nil {
		return err
	}

	if err := tx.WriteTxOutCount(w); err != nil {
		return err
	}

	if err := tx.WriteTxOuts(w); err != nil {
		return err
	}

	if err := tx.WriteLockTime(w); err != nil {
		return err
	}

	return nil
}

func (tx *Tx) WriteVersion(w io.Writer) error {
	return writeData(w, tx.Version)
}

func (tx *Tx) WriteTxInCount(w io.Writer) error {
	return writeVarInt(w, len(tx.TxIns))
}

func (tx *Tx) WriteTxIns(w io.Writer) error {
	for _, txIn := range tx.TxIns {
		if err := txIn.Write(w); err != nil {
			return err
		}
	}

	return nil
}

func (tx *Tx) WriteTxOutCount(w io.Writer) error {
	return writeVarInt(w, len(tx.TxOuts))
}

func (tx *Tx) WriteTxOuts(w io.Writer) error {
	for _, txOut := range tx.TxOuts {
		if err := txOut.Write(w); err != nil {
			return err
		}
	}

	return nil
}

func (tx *Tx) WriteLockTime(w io.Writer) error {
	return writeData(w, tx.LockTime)
}
