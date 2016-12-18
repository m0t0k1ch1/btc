package btctx

import (
	"bytes"
	"encoding/hex"
	"io"
)

const (
	DefaultTxVersion  int32  = 1
	DefaultTxLockTime uint32 = 0
)

type Tx struct {
	Version  int32    `json:"version"`
	TxIns    []*TxIn  `json:"txIns"`
	TxOuts   []*TxOut `json:"txOuts"`
	LockTime uint32   `json:"lockTime"`
}

func NewTx() *Tx {
	return &Tx{
		Version:  DefaultTxVersion,
		LockTime: DefaultTxLockTime,
	}
}

func NewTxFromHex(s string) (*Tx, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return NewTxFromBytes(b)
}

func NewTxFromBytes(b []byte) (*Tx, error) {
	txd := newTxDecoder(b)

	return txd.decode()
}

func (tx *Tx) AddTxIn(txIn *TxIn) {
	tx.TxIns = append(tx.TxIns, txIn)
}

func (tx *Tx) AddTxOut(txOut *TxOut) {
	tx.TxOuts = append(tx.TxOuts, txOut)
}

func (tx *Tx) ToBytes() ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := tx.WriteAll(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (tx *Tx) ToHex() (string, error) {
	b, err := tx.ToBytes()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func (tx *Tx) ToHash() (string, error) {
	txBytes, err := tx.ToBytes()
	if err != nil {
		return "", err
	}

	hashBytes, err := sha256Double(txBytes)
	if err != nil {
		return "", err
	}
	hashBytesLen := len(hashBytes)

	buf := &bytes.Buffer{}
	for i := 1; i <= hashBytesLen; i++ {
		buf.WriteByte(hashBytes[hashBytesLen-i])
	}

	return hex.EncodeToString(buf.Bytes()), nil
}

func (tx *Tx) WriteAll(w io.Writer) error {
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
	return writeVarInt(w, uint(len(tx.TxIns)))
}

func (tx *Tx) WriteTxIns(w io.Writer) error {
	for _, txIn := range tx.TxIns {
		if err := txIn.WriteAll(w); err != nil {
			return err
		}
	}

	return nil
}

func (tx *Tx) WriteTxOutCount(w io.Writer) error {
	return writeVarInt(w, uint(len(tx.TxOuts)))
}

func (tx *Tx) WriteTxOuts(w io.Writer) error {
	for _, txOut := range tx.TxOuts {
		if err := txOut.WriteAll(w); err != nil {
			return err
		}
	}

	return nil
}

func (tx *Tx) WriteLockTime(w io.Writer) error {
	return writeData(w, tx.LockTime)
}
