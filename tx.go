package btctx

type TxIn struct {
	Hash            string `json:"hash"`
	Index           uint32 `json:"index"`
	SignatureScript string `json:"signatureScript"`
	Sequence        uint32 `json:"sequence"`
}

type TxOut struct {
	Value    int64  `json:"value"` // satoshi
	PkScript string `json:"pkScript"`
}

type Tx struct {
	Version  int32    `json:"version"`
	TxIns    []*TxIn  `json:"txIns"`
	TxOuts   []*TxOut `json:"txOuts"`
	LockTime uint32   `json:"lockTime"`
}

func NewTxFromBytes(b []byte) (*Tx, error) {
	txr := NewTxReader(b)

	return txr.ReadTx()
}
