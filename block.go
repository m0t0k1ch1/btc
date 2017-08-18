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

func NewBlockFromBytes(b []byte) (*Block, error) {
	return newReader(b).readBlock()
}

func (block *Block) Bytes() ([]byte, error) {
	w := newWriter()
	if err := w.writeBlock(block); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func (block *Block) Hex() (string, error) {
	b, err := block.Bytes()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func (block *Block) Blockhash() (string, error) {
	return block.BlockHeader.Blockhash()
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
	return newReader(b).readBlockHeader()
}

func (bh *BlockHeader) Bytes() ([]byte, error) {
	w := newWriter()
	if err := w.writeBlockHeader(bh); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func (bh *BlockHeader) Hex() (string, error) {
	b, err := bh.Bytes()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func (bh *BlockHeader) Blockhash() (string, error) {
	blockBytes, err := bh.Bytes()
	if err != nil {
		return "", err
	}

	hashBytes, err := Sha256Double(blockBytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(reverseBytes(hashBytes)), nil
}
