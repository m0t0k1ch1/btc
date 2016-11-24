package btctx

type TxIn struct {
	TxId      string `json:"txid"`
	Index     uint32 `json:"index"`
	Signature string `json:"signature"`
	Sequence  uint32 `json:"sequence"`
}

type TxOut struct {
	Value  int64  `json:"value"`
	Script string `json:"script"`
}

type Tx struct {
	Version  int32    `json:"version"`
	TxIns    []*TxIn  `json:"txin"`
	TxOuts   []*TxOut `json:"txout"`
	LockTime uint32   `json:"locktime"`
}

func NewTxFromBytes(b []byte) (*Tx, error) {
	txr := NewTxReader(b)

	return txr.ReadTx()
}
