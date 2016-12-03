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
	if err := tx.write(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (tx *Tx) write(w io.Writer) error {
	if err := tx.writeVersion(w); err != nil {
		return err
	}

	if err := tx.writeTxIns(w); err != nil {
		return err
	}

	if err := tx.writeTxOuts(w); err != nil {
		return err
	}

	if err := tx.writeLockTime(w); err != nil {
		return err
	}

	return nil
}

func (tx *Tx) writeVersion(w io.Writer) error {
	return writeData(w, tx.Version)
}

func (tx *Tx) writeTxIns(w io.Writer) error {
	if err := writeVarInt(w, len(tx.TxIns)); err != nil {
		return err
	}

	for _, txIn := range tx.TxIns {
		if err := txIn.write(w); err != nil {
			return err
		}
	}

	return nil
}

func (tx *Tx) writeTxOuts(w io.Writer) error {
	if err := writeVarInt(w, len(tx.TxOuts)); err != nil {
		return err
	}

	for _, txOut := range tx.TxOuts {
		if err := txOut.write(w); err != nil {
			return err
		}
	}

	return nil
}

func (tx *Tx) writeLockTime(w io.Writer) error {
	return writeData(w, tx.LockTime)
}
