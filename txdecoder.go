package btctx

type txDecoder struct {
	*decoder
}

func newTxDecoder(b []byte) *txDecoder {
	return &txDecoder{newDecoder(b)}
}

func (txd *txDecoder) decodeVersion() (int32, error) {
	return txd.readInt32()
}

func (txd *txDecoder) decodeTxIns() ([]*TxIn, error) {
	txInCnt, err := txd.readVarInt()
	if err != nil {
		return nil, err
	}

	txIns := make([]*TxIn, txInCnt)
	for i := 0; i < int(txInCnt); i++ {
		txIn, err := txd.decodeTxIn()
		if err != nil {
			return nil, err
		}

		txIns[i] = txIn
	}

	return txIns, nil
}

func (txd *txDecoder) decodeTxIn() (*TxIn, error) {
	hash, err := txd.readHexReverse(32)
	if err != nil {
		return nil, err
	}

	index, err := txd.readUint32()
	if err != nil {
		return nil, err
	}

	sigScriptLen, err := txd.readVarInt()
	if err != nil {
		return nil, err
	}

	sigScriptHex, err := txd.readHex(sigScriptLen)
	if err != nil {
		return nil, err
	}
	sigScript, err := NewScriptFromHex(sigScriptHex)
	if err != nil {
		return nil, err
	}

	seq, err := txd.readUint32()
	if err != nil {
		return nil, err
	}

	return &TxIn{
		Hash:      hash,
		Index:     index,
		SigScript: sigScript,
		Sequence:  seq,
	}, nil
}

func (txd *txDecoder) decodeTxOuts() ([]*TxOut, error) {
	txOutCnt, err := txd.readVarInt()
	if err != nil {
		return nil, err
	}

	txOuts := make([]*TxOut, txOutCnt)
	for i := 0; i < int(txOutCnt); i++ {
		txOut, err := txd.decodeTxOut()
		if err != nil {
			return nil, err
		}

		txOuts[i] = txOut
	}

	return txOuts, nil
}

func (txd *txDecoder) decodeTxOut() (*TxOut, error) {
	value, err := txd.readInt64()
	if err != nil {
		return nil, err
	}

	pkScriptLen, err := txd.readVarInt()
	if err != nil {
		return nil, err
	}

	pkScriptHex, err := txd.readHex(pkScriptLen)
	if err != nil {
		return nil, err
	}
	pkScript, err := NewScriptFromHex(pkScriptHex)
	if err != nil {
		return nil, err
	}

	return &TxOut{
		Value:    Satoshi(value),
		PkScript: pkScript,
	}, nil
}

func (txd *txDecoder) decodeLockTime() (uint32, error) {
	return txd.readUint32()
}

func (txd *txDecoder) decode() (*Tx, error) {
	version, err := txd.decodeVersion()
	if err != nil {
		return nil, err
	}

	txIns, err := txd.decodeTxIns()
	if err != nil {
		return nil, err
	}

	txOuts, err := txd.decodeTxOuts()
	if err != nil {
		return nil, err
	}

	lockTime, err := txd.decodeLockTime()
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
