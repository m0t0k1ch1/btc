package btctx

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
)

type TxWriter struct {
	tx  *Tx
	buf *bytes.Buffer
}

func NewTxWriter(tx *Tx) *TxWriter {
	return &TxWriter{
		tx:  tx,
		buf: &bytes.Buffer{},
	}
}

func (txw *TxWriter) writeData(data interface{}) error {
	return binary.Write(txw.buf, binary.LittleEndian, data)
}

func (txw *TxWriter) writeHex(data string) error {
	b, err := hex.DecodeString(data)
	if err != nil {
		return err
	}

	return txw.writeData(b)
}

func (txw *TxWriter) writeHexReverse(data string) error {
	b, err := hex.DecodeString(data)
	if err != nil {
		return err
	}
	size := len(b)

	reversed := make([]byte, size)

	for i, j := size-1, 0; i >= 0; i-- {
		reversed[j] = b[i]
		j++
	}

	return txw.writeData(reversed)
}

// variable length integer
// ref. https://en.bitcoin.it/wiki/Protocol_documentation#Variable_length_integer
func (txw *TxWriter) writeVarInt(data int) error {
	if data < 0xfd {
		return txw.writeData(uint8(data))
	} else if data <= 0xffff {
		return txw.writeData(uint16(data))
	} else if data <= 0xffffffff {
		return txw.writeData(uint32(data))
	} else {
		return txw.writeData(uint64(data))
	}
}

func (txw *TxWriter) writeVersion() error {
	return txw.writeData(txw.tx.Version)
}

func (txw *TxWriter) writeTxIns() error {
	txIns := txw.tx.TxIns

	if err := txw.writeVarInt(len(txIns)); err != nil {
		return err
	}

	for _, txIn := range txIns {
		if err := txw.writeHexReverse(txIn.Hash); err != nil {
			return err
		}

		if err := txw.writeData(txIn.Index); err != nil {
			return err
		}

		sigScript := txIn.SignatureScript
		sigScriptBytes, err := hex.DecodeString(sigScript)
		if err != nil {
			return err
		}
		if err := txw.writeVarInt(len(sigScriptBytes)); err != nil {
			return err
		}
		if err := txw.writeHex(sigScript); err != nil {
			return err
		}

		if err := txw.writeData(txIn.Sequence); err != nil {
			return err
		}
	}

	return nil
}

func (txw *TxWriter) writeLockTime() error {
	return txw.writeData(txw.tx.LockTime)
}

func (txw *TxWriter) writeTxOuts() error {
	txOuts := txw.tx.TxOuts

	if err := txw.writeVarInt(len(txOuts)); err != nil {
		return err
	}

	for _, txOut := range txOuts {
		if err := txw.writeData(txOut.Value); err != nil {
			return err
		}

		pkScript := txOut.PkScript
		pkScriptBytes, err := hex.DecodeString(pkScript)
		if err != nil {
			return err
		}
		if err := txw.writeVarInt(len(pkScriptBytes)); err != nil {
			return err
		}
		if err := txw.writeHex(pkScript); err != nil {
			return err
		}
	}

	return nil
}

func (txw *TxWriter) WriteTx() ([]byte, error) {
	if err := txw.writeVersion(); err != nil {
		return nil, err
	}

	if err := txw.writeTxIns(); err != nil {
		return nil, err
	}

	if err := txw.writeTxOuts(); err != nil {
		return nil, err
	}

	if err := txw.writeLockTime(); err != nil {
		return nil, err
	}

	return txw.buf.Bytes(), nil
}
