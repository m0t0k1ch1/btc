package btc

import "encoding/hex"

type TxIn struct {
	Txid     string  `json:"txid"`
	Index    uint32  `json:"index"`
	Script   *Script `json:"script"`
	Sequence uint32  `json:"sequence"`
}

func NewTxIn(txid string, index uint32, script *Script) *TxIn {
	return &TxIn{
		Txid:     txid,
		Index:    index,
		Script:   script,
		Sequence: TxInSequence,
	}
}

type TxOut struct {
	Amount Satoshi `json:"amount"`
	Script *Script `json:"script"`
}

func NewTxOut(amount int64, script *Script) *TxOut {
	return &TxOut{
		Amount: Satoshi(amount),
		Script: script,
	}
}

type Tx struct {
	Version  int32    `json:"version"`
	TxIns    []*TxIn  `json:"txIns"`
	TxOuts   []*TxOut `json:"txOuts"`
	LockTime uint32   `json:"lockTime"`
}

func NewTx() *Tx {
	return &Tx{
		Version:  TxVersion,
		TxIns:    []*TxIn{},
		TxOuts:   []*TxOut{},
		LockTime: TxLockTime,
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
	return newReader(b).readTx()
}

func (tx *Tx) AddTxIn(txIn *TxIn) {
	tx.TxIns = append(tx.TxIns, txIn)
}

func (tx *Tx) AddTxOut(txOut *TxOut) {
	tx.TxOuts = append(tx.TxOuts, txOut)
}

func (tx *Tx) Bytes() ([]byte, error) {
	w := newWriter()
	if err := w.writeTx(tx); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func (tx *Tx) Hex() (string, error) {
	b, err := tx.Bytes()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func (tx *Tx) Txid() (string, error) {
	txBytes, err := tx.Bytes()
	if err != nil {
		return "", err
	}

	hashBytes, err := sha256Double(txBytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(reverseBytes(hashBytes)), nil
}
