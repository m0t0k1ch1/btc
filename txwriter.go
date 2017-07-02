package btc

type txWriter struct {
	*writer
}

func newTxWriter() *txWriter {
	return &txWriter{newWriter()}
}

func (w *txWriter) writeVersion(version int32) error {
	return w.writeData(version)
}

func (w *txWriter) writeTxIns(txIns []*TxIn) error {
	if err := w.writeVarInt(uint(len(txIns))); err != nil {
		return err
	}

	for _, txIn := range txIns {
		if err := w.writeTxIn(txIn); err != nil {
			return err
		}
	}

	return nil
}

func (w *txWriter) writeTxIn(txIn *TxIn) error {
	if err := w.writeHexReverse(txIn.Txid); err != nil {
		return err
	}

	if err := w.writeData(txIn.Index); err != nil {
		return err
	}

	if err := w.writeScript(txIn.Script); err != nil {
		return err
	}

	if err := w.writeData(txIn.Sequence); err != nil {
		return err
	}

	return nil
}

func (w *txWriter) writeTxOuts(txOuts []*TxOut) error {
	if err := w.writeVarInt(uint(len(txOuts))); err != nil {
		return err
	}

	for _, txOut := range txOuts {
		if err := w.writeTxOut(txOut); err != nil {
			return err
		}
	}

	return nil
}

func (w *txWriter) writeTxOut(txOut *TxOut) error {
	if err := w.writeData(txOut.Amount); err != nil {
		return err
	}

	if err := w.writeScript(txOut.Script); err != nil {
		return err
	}

	return nil
}

func (w *txWriter) writeLockTime(lockTime uint32) error {
	return w.writeData(lockTime)
}

func (w *txWriter) writeScript(script *Script) error {
	b, err := script.Bytes()
	if err != nil {
		return err
	}

	if err := w.writeVarInt(uint(len(b))); err != nil {
		return err
	}

	if _, err := w.Write(b); err != nil {
		return err
	}

	return nil
}

func (w *txWriter) writeTx(tx *Tx) error {
	if err := w.writeVersion(tx.Version); err != nil {
		return err
	}

	if err := w.writeTxIns(tx.TxIns); err != nil {
		return err
	}

	if err := w.writeTxOuts(tx.TxOuts); err != nil {
		return err
	}

	if err := w.writeLockTime(tx.LockTime); err != nil {
		return err
	}

	return nil
}
