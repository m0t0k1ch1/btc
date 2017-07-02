package btc

type txReader struct {
	*reader
}

func newTxReader(b []byte) *txReader {
	return &txReader{newReader(b)}
}

func (r *txReader) readVersion() (int32, error) {
	return r.readInt32()
}

func (r *txReader) readTxIns() ([]*TxIn, error) {
	txInCnt, err := r.readVarInt()
	if err != nil {
		return nil, err
	}

	txIns := make([]*TxIn, txInCnt)
	for i := 0; i < int(txInCnt); i++ {
		txIn, err := r.readTxIn()
		if err != nil {
			return nil, err
		}

		txIns[i] = txIn
	}

	return txIns, nil
}

func (r *txReader) readTxIn() (*TxIn, error) {
	txid, err := r.readHexReverse(32)
	if err != nil {
		return nil, err
	}

	index, err := r.readUint32()
	if err != nil {
		return nil, err
	}

	scriptLen, err := r.readVarInt()
	if err != nil {
		return nil, err
	}
	scriptHex, err := r.readHex(scriptLen)
	if err != nil {
		return nil, err
	}
	var script *Script
	if txid == CoinBaseTxid {
		script = &Script{
			Hex: scriptHex,
		}
	} else {
		script, err = NewScriptFromHex(scriptHex)
		if err != nil {
			return nil, err
		}
	}

	seq, err := r.readUint32()
	if err != nil {
		return nil, err
	}

	return &TxIn{
		Txid:     txid,
		Index:    index,
		Script:   script,
		Sequence: seq,
	}, nil
}

func (r *txReader) readTxOuts() ([]*TxOut, error) {
	txOutCnt, err := r.readVarInt()
	if err != nil {
		return nil, err
	}

	txOuts := make([]*TxOut, txOutCnt)
	for i := 0; i < int(txOutCnt); i++ {
		txOut, err := r.readTxOut()
		if err != nil {
			return nil, err
		}

		txOuts[i] = txOut
	}

	return txOuts, nil
}

func (r *txReader) readTxOut() (*TxOut, error) {
	amount, err := r.readInt64()
	if err != nil {
		return nil, err
	}

	scriptLen, err := r.readVarInt()
	if err != nil {
		return nil, err
	}
	scriptHex, err := r.readHex(scriptLen)
	if err != nil {
		return nil, err
	}
	script, err := NewScriptFromHex(scriptHex)
	if err != nil {
		return nil, err
	}

	return &TxOut{
		Amount: Satoshi(amount),
		Script: script,
	}, nil
}

func (r *txReader) readLockTime() (uint32, error) {
	return r.readUint32()
}

func (r *txReader) readTx() (*Tx, error) {
	version, err := r.readVersion()
	if err != nil {
		return nil, err
	}

	txIns, err := r.readTxIns()
	if err != nil {
		return nil, err
	}

	txOuts, err := r.readTxOuts()
	if err != nil {
		return nil, err
	}

	lockTime, err := r.readLockTime()
	if err != nil {
		return nil, err
	}

	return &Tx{
		Version:  version,
		TxIns:    txIns,
		TxOuts:   txOuts,
		LockTime: lockTime,
	}, nil
}
