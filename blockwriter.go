package btc

type blockWriter struct {
	*writer
}

func newBlockWriter() *blockWriter {
	return &blockWriter{newWriter()}
}

func (w *blockWriter) writeVersion(version int32) error {
	return w.writeData(version)
}

func (w *blockWriter) writePrevBlockhash(prevBlockhash string) error {
	return w.writeHexReverse(prevBlockhash)
}

func (w *blockWriter) writeMerkleRoot(merkleRoot string) error {
	return w.writeHex(merkleRoot)
}

func (w *blockWriter) writeTimestamp(timestamp uint32) error {
	return w.writeData(timestamp)
}

func (w *blockWriter) writeBits(bits uint32) error {
	return w.writeData(bits)
}

func (w *blockWriter) writeNonce(nonce uint32) error {
	return w.writeData(nonce)
}

func (w *blockWriter) writeTxes(txes []*Tx) error {
	if err := w.writeVarInt(uint(len(txes))); err != nil {
		return nil
	}

	for _, tx := range txes {
		if err := w.writeTx(tx); err != nil {
			return err
		}
	}

	return nil
}

func (w *blockWriter) writeTx(tx *Tx) error {
	txWriter := newTxWriter()
	if err := txWriter.writeTx(tx); err != nil {
		return err
	}

	if _, err := w.Write(txWriter.Bytes()); err != nil {
		return err
	}

	return nil
}

func (w *blockWriter) writeBlockHeaderBase(bh *BlockHeader) error {
	if err := w.writeVersion(bh.Version); err != nil {
		return err
	}

	if err := w.writePrevBlockhash(bh.PrevBlockhash); err != nil {
		return err
	}

	if err := w.writeMerkleRoot(bh.MerkleRoot); err != nil {
		return err
	}

	if err := w.writeTimestamp(bh.Timestamp); err != nil {
		return err
	}

	if err := w.writeBits(bh.Bits); err != nil {
		return err
	}

	if err := w.writeNonce(bh.Nonce); err != nil {
		return err
	}

	return nil
}

func (w *blockWriter) writeBlockHeader(bh *BlockHeader) error {
	if err := w.writeBlockHeaderBase(bh); err != nil {
		return err
	}

	// number of transaction entries, this value is always 0
	// ref. https://en.bitcoin.it/wiki/Protocol_documentation#Block_Headers
	if err := w.WriteByte(0); err != nil {
		return err
	}

	return nil
}

func (w *blockWriter) writeBlock(block *Block) error {
	if err := w.writeBlockHeaderBase(block.BlockHeader); err != nil {
		return err
	}

	if err := w.writeTxes(block.Txes); err != nil {
		return err
	}

	return nil
}
