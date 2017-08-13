package btc

import "encoding/hex"

type Block struct {
	*BlockHeader
	Txes []*Tx
}

func NewBlockFromHex(s string) (*Block, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return NewBlockFromBytes(b)
}

func (block *Block) Bytes() ([]byte, error) {
	w := newBlockWriter()
	if err := w.writeBlock(block); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func NewBlockFromBytes(b []byte) (*Block, error) {
	return newBlockReader(b).readBlock()
}

type BlockHeader struct {
	Version       int32
	PrevBlockhash string
	MerkleRoot    string
	Timestamp     uint32
	Bits          uint32
	Nonce         uint32
}

func NewBlockHeaderFromHex(s string) (*BlockHeader, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return NewBlockHeaderFromBytes(b)
}

func NewBlockHeaderFromBytes(b []byte) (*BlockHeader, error) {
	return newBlockReader(b).readBlockHeader()
}

func (bh *BlockHeader) Bytes() ([]byte, error) {
	w := newBlockWriter()
	if err := w.writeBlockHeader(bh); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}
