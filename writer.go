package btc

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
)

type writer struct {
	*bytes.Buffer
}

func newWriter() *writer {
	return &writer{&bytes.Buffer{}}
}

func (w *writer) writeData(data interface{}) error {
	return binary.Write(w, binary.LittleEndian, data)
}

func (w *writer) writeHex(data string) error {
	b, err := hex.DecodeString(data)
	if err != nil {
		return err
	}

	if _, err := w.Write(b); err != nil {
		return err
	}

	return nil
}

func (w *writer) writeHexReverse(data string) error {
	b, err := hex.DecodeString(data)
	if err != nil {
		return err
	}

	if _, err := w.Write(reverseBytes(b)); err != nil {
		return err
	}

	return nil
}

// variable length integer
// ref. https://en.bitcoin.it/wiki/Protocol_documentation#Variable_length_integer
func (w *writer) writeVarInt(data uint) error {
	if data < 0xfd {
		return w.writeData(byte(data))
	} else if data <= 0xffff {
		if err := w.WriteByte(0xfd); err != nil {
			return err
		}
		return w.writeData(uint16(data))
	} else if data <= 0xffffffff {
		if err := w.WriteByte(0xfe); err != nil {
			return err
		}
		return w.writeData(uint32(data))
	} else {
		if err := w.WriteByte(0xff); err != nil {
			return err
		}
		return w.writeData(uint64(data))
	}
}

func (w *writer) writeTxVersion(version int32) error {
	return w.writeData(version)
}

func (w *writer) writeTxIns(txIns []*TxIn) error {
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

func (w *writer) writeTxIn(txIn *TxIn) error {
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

func (w *writer) writeTxOuts(txOuts []*TxOut) error {
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

func (w *writer) writeTxOut(txOut *TxOut) error {
	if err := w.writeData(txOut.Amount); err != nil {
		return err
	}

	if err := w.writeScript(txOut.Script); err != nil {
		return err
	}

	return nil
}

func (w *writer) writeLockTime(lockTime uint32) error {
	return w.writeData(lockTime)
}

func (w *writer) writeScript(script *Script) error {
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

func (w *writer) writeTx(tx *Tx) error {
	if err := w.writeTxVersion(tx.Version); err != nil {
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

func (w *writer) writeTxes(txes []*Tx) error {
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

func (w *writer) writeBlockVersion(version int32) error {
	return w.writeData(version)
}

func (w *writer) writePrevBlockhash(prevBlockhash string) error {
	return w.writeHexReverse(prevBlockhash)
}

func (w *writer) writeMerkleRoot(merkleRoot string) error {
	return w.writeHexReverse(merkleRoot)
}

func (w *writer) writeTimestamp(timestamp uint32) error {
	return w.writeData(timestamp)
}

func (w *writer) writeBits(bits uint32) error {
	return w.writeData(bits)
}

func (w *writer) writeNonce(nonce uint32) error {
	return w.writeData(nonce)
}

func (w *writer) writeBlockHeader(bh *BlockHeader) error {
	if err := w.writeBlockVersion(bh.Version); err != nil {
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

func (w *writer) writeBlock(block *Block) error {
	if err := w.writeBlockHeader(block.BlockHeader); err != nil {
		return err
	}

	if err := w.writeTxes(block.Txes); err != nil {
		return err
	}

	return nil
}
