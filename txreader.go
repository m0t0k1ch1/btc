package btctx

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type txReader struct {
	r *bytes.Reader
}

func newTxReader(b []byte) *txReader {
	return &txReader{
		r: bytes.NewReader(b),
	}
}

func (txr *txReader) readByte() (byte, error) {
	return txr.r.ReadByte()
}

func (txr *txReader) readBytes(size int) ([]byte, error) {
	data := make([]byte, size)

	for i := 0; i < size; i++ {
		b, err := txr.readByte()
		if err != nil {
			return nil, err
		}

		data[i] = b
	}

	return data, nil
}

func (txr *txReader) readBytesReverse(size int) ([]byte, error) {
	data := make([]byte, size)

	for i := size - 1; i >= 0; i-- {
		b, err := txr.readByte()
		if err != nil {
			return nil, err
		}

		data[i] = b
	}

	return data, nil
}

func (txr *txReader) readString(size int) (string, error) {
	b, err := txr.readBytes(size)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (txr *txReader) readStringReverse(size int) (string, error) {
	b, err := txr.readBytesReverse(size)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (txr *txReader) readHex(size int) (string, error) {
	s, err := txr.readString(size)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", s), nil
}

func (txr *txReader) readHexReverse(size int) (string, error) {
	s, err := txr.readStringReverse(size)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", s), nil
}

func (txr *txReader) readBinary(size int, data interface{}) error {
	b, err := txr.readBytes(size)
	if err != nil {
		return err
	}

	buf := bytes.NewReader(b)
	if err := binary.Read(buf, binary.LittleEndian, data); err != nil {
		return err
	}

	return nil
}

func (txr *txReader) readInt16() (int16, error) {
	var data int16
	if err := txr.readBinary(2, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (txr *txReader) readInt32() (int32, error) {
	var data int32
	if err := txr.readBinary(4, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (txr *txReader) readInt64() (int64, error) {
	var data int64
	if err := txr.readBinary(8, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (txr *txReader) readUint16() (uint16, error) {
	var data uint16
	if err := txr.readBinary(2, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (txr *txReader) readUint32() (uint32, error) {
	var data uint32
	if err := txr.readBinary(4, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (txr *txReader) readUint64() (uint64, error) {
	var data uint64
	if err := txr.readBinary(8, &data); err != nil {
		return 0, err
	}

	return data, nil
}

// variable length integer
// ref. https://en.bitcoin.it/wiki/Protocol_documentation#Variable_length_integer
func (txr *txReader) readVarInt() (uint, error) {
	head, err := txr.readByte()
	if err != nil {
		return 0, err
	}

	switch head {
	case 0xff:
		val, err := txr.readUint64()
		if err != nil {
			return 0, err
		}
		return uint(val), nil
	case 0xfe:
		val, err := txr.readUint32()
		if err != nil {
			return 0, err
		}
		return uint(val), nil
	case 0xfd:
		val, err := txr.readUint16()
		if err != nil {
			return 0, err
		}
		return uint(val), nil
	default:
		return uint(head), nil
	}
}

func (txr *txReader) readVersion() (int32, error) {
	return txr.readInt32()
}

func (txr *txReader) readTxIns() ([]*TxIn, error) {
	txInCnt, err := txr.readVarInt()
	if err != nil {
		return nil, err
	}

	txIns := make([]*TxIn, txInCnt)
	for i := 0; i < int(txInCnt); i++ {
		txIn, err := txr.readTxIn()
		if err != nil {
			return nil, err
		}

		txIns[i] = txIn
	}

	return txIns, nil
}

func (txr *txReader) readTxIn() (*TxIn, error) {
	hash, err := txr.readHexReverse(32)
	if err != nil {
		return nil, err
	}

	index, err := txr.readUint32()
	if err != nil {
		return nil, err
	}

	sigScriptLen, err := txr.readVarInt()
	if err != nil {
		return nil, err
	}

	sigScript, err := txr.readHex(int(sigScriptLen))
	if err != nil {
		return nil, err
	}

	seq, err := txr.readUint32()
	if err != nil {
		return nil, err
	}

	txin := &TxIn{
		Hash:            hash,
		Index:           index,
		SignatureScript: sigScript,
		Sequence:        seq,
	}

	return txin, nil
}

func (txr *txReader) readTxOuts() ([]*TxOut, error) {
	txOutCnt, err := txr.readVarInt()
	if err != nil {
		return nil, err
	}

	txOuts := make([]*TxOut, txOutCnt)
	for i := 0; i < int(txOutCnt); i++ {
		txOut, err := txr.readTxOut()
		if err != nil {
			return nil, err
		}

		txOuts[i] = txOut
	}

	return txOuts, nil
}

func (txr *txReader) readTxOut() (*TxOut, error) {
	value, err := txr.readInt64()
	if err != nil {
		return nil, err
	}

	pkScriptLen, err := txr.readVarInt()
	if err != nil {
		return nil, err
	}

	pkScript, err := txr.readHex(int(pkScriptLen))
	if err != nil {
		return nil, err
	}

	txout := &TxOut{
		Value:    value,
		PkScript: pkScript,
	}

	return txout, nil
}

func (txr *txReader) readLockTime() (uint32, error) {
	return txr.readUint32()
}

func (txr *txReader) readAll() (*Tx, error) {
	version, err := txr.readVersion()
	if err != nil {
		return nil, err
	}

	txIns, err := txr.readTxIns()
	if err != nil {
		return nil, err
	}

	txOuts, err := txr.readTxOuts()
	if err != nil {
		return nil, err
	}

	lockTime, err := txr.readLockTime()
	if err != nil {
		return nil, err
	}

	tx := &Tx{
		Version:  version,
		TxIns:    txIns,
		TxOuts:   txOuts,
		LockTime: lockTime,
	}

	return tx, nil
}
