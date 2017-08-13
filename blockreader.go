package btc

type blockReader struct {
	*reader
}

func newBlockReader(b []byte) *blockReader {
	return &blockReader{newReader(b)}
}

func (r *blockReader) readVersion() (int32, error) {
	return r.readInt32()
}

func (r *blockReader) readPrevBlock() (string, error) {
	return r.readHexReverse(32)
}

func (r *blockReader) readMerkleRoot() (string, error) {
	return r.readHexReverse(32)
}

func (r *blockReader) readTimestamp() (uint32, error) {
	return r.readUint32()
}

func (r *blockReader) readBits() (uint32, error) {
	return r.readUint32()
}

func (r *blockReader) readNonce() (uint32, error) {
	return r.readUint32()
}

func (r *blockReader) readTxes() ([]*Tx, error) {
	cnt, err := r.readVarInt()
	if err != nil {
		return nil, err
	}

	tr := newTxReader(r.Bytes())

	txes := make([]*Tx, cnt)
	for i := uint(0); i < cnt; i++ {
		tx, err := tr.readTx()
		if err != nil {
			return nil, err
		}

		txes[i] = tx
	}

	return txes, nil
}

func (r *blockReader) readBlock() (*Block, error) {
	bh, err := r.readBlockHeader()
	if err != nil {
		return nil, err
	}

	txes, err := r.readTxes()
	if err != nil {
		return nil, err
	}

	return &Block{
		BlockHeader: bh,
		Txes:        txes,
	}, nil
}

func (r *blockReader) readBlockHeader() (*BlockHeader, error) {
	version, err := r.readVersion()
	if err != nil {
		return nil, err
	}

	prevBlock, err := r.readPrevBlock()
	if err != nil {
		return nil, err
	}

	merkleRoot, err := r.readMerkleRoot()
	if err != nil {
		return nil, err
	}

	timestamp, err := r.readTimestamp()
	if err != nil {
		return nil, err
	}

	bits, err := r.readBits()
	if err != nil {
		return nil, err
	}

	nonce, err := r.readNonce()
	if err != nil {
		return nil, err
	}

	return &BlockHeader{
		Version:       version,
		PrevBlockhash: prevBlock,
		MerkleRoot:    merkleRoot,
		Timestamp:     timestamp,
		Bits:          bits,
		Nonce:         nonce,
	}, nil
}
