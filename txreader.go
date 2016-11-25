package btctx

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type TxReader struct {
	r *bytes.Reader
}

func NewTxReader(b []byte) *TxReader {
	return &TxReader{
		r: bytes.NewReader(b),
	}
}

func (txr *TxReader) readByte() (byte, error) {
	return txr.r.ReadByte()
}

func (txr *TxReader) readBytes(size int) ([]byte, error) {
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

func (txr *TxReader) readBytesReverse(size int) ([]byte, error) {
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

func (txr *TxReader) readString(size int) (string, error) {
	b, err := txr.readBytes(size)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (txr *TxReader) readStringReverse(size int) (string, error) {
	b, err := txr.readBytesReverse(size)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (txr *TxReader) readBinary(size int, data interface{}) error {
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

func (txr *TxReader) readInt16() (int16, error) {
	var data int16
	if err := txr.readBinary(2, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (txr *TxReader) readInt32() (int32, error) {
	var data int32
	if err := txr.readBinary(4, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (txr *TxReader) readInt64() (int64, error) {
	var data int64
	if err := txr.readBinary(8, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (txr *TxReader) readUint16() (uint16, error) {
	var data uint16
	if err := txr.readBinary(2, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (txr *TxReader) readUint32() (uint32, error) {
	var data uint32
	if err := txr.readBinary(4, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (txr *TxReader) readUint64() (uint64, error) {
	var data uint64
	if err := txr.readBinary(8, &data); err != nil {
		return 0, err
	}

	return data, nil
}

// variable length integer
// ref. https://en.bitcoin.it/wiki/Protocol_documentation#Variable_length_integer
func (txr *TxReader) readVarInt() (uint, error) {
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

func (txr *TxReader) readVersion() (int32, error) {
	return txr.readInt32()
}

func (txr *TxReader) readTxIns() ([]*TxIn, error) {
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

func (txr *TxReader) readTxIn() (*TxIn, error) {
	hash, err := txr.readStringReverse(32)
	if err != nil {
		return nil, err
	}
	hash = fmt.Sprintf("%x", hash)

	index, err := txr.readUint32()
	if err != nil {
		return nil, err
	}

	signatureLen, err := txr.readVarInt()
	if err != nil {
		return nil, err
	}

	signature, err := txr.readString(int(signatureLen))
	if err != nil {
		return nil, err
	}
	signature = fmt.Sprintf("%x", signature)

	sequence, err := txr.readUint32()
	if err != nil {
		return nil, err
	}

	txin := &TxIn{
		Hash:            hash,
		Index:           index,
		SignatureScript: signature,
		Sequence:        sequence,
	}

	return txin, nil
}

func (txr *TxReader) readTxOuts() ([]*TxOut, error) {
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

func (txr *TxReader) readTxOut() (*TxOut, error) {
	value, err := txr.readInt64()
	if err != nil {
		return nil, err
	}

	scriptLen, err := txr.readVarInt()
	if err != nil {
		return nil, err
	}

	script, err := txr.readString(int(scriptLen))
	if err != nil {
		return nil, err
	}
	script = fmt.Sprintf("%x", script)

	txout := &TxOut{
		Value:    value,
		PkScript: script,
	}

	return txout, nil
}

func (txr *TxReader) readLockTime() (uint32, error) {
	return txr.readUint32()
}

func (txr *TxReader) ReadTx() (*Tx, error) {
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
