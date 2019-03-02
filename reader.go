package btc

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io"
	"strings"
)

type reader struct {
	*bytes.Buffer
}

func newReader(b []byte) *reader {
	return &reader{bytes.NewBuffer(b)}
}

func (r *reader) readData(size uint, data interface{}) error {
	b, err := r.readBytes(size)
	if err != nil {
		return err
	}

	return binary.Read(bytes.NewReader(b), binary.LittleEndian, data)
}

func (r *reader) readBytes(size uint) ([]byte, error) {
	data := make([]byte, size)

	var i uint
	for i = 0; i < size; i++ {
		c, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		data[i] = c
	}

	return data, nil
}

func (r *reader) readBytesReverse(size uint) ([]byte, error) {
	data := make([]byte, size)

	var i uint
	for i = 1; i <= size; i++ {
		c, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		data[size-i] = c
	}

	return data, nil
}

func (r *reader) readString(size uint) (string, error) {
	b, err := r.readBytes(size)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (r *reader) readStringReverse(size uint) (string, error) {
	b, err := r.readBytesReverse(size)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (r *reader) readHex(size uint) (string, error) {
	b, err := r.readBytes(size)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func (r *reader) readHexReverse(size uint) (string, error) {
	b, err := r.readBytesReverse(size)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func (r *reader) readInt16() (int16, error) {
	var data int16
	if err := r.readData(2, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (r *reader) readInt32() (int32, error) {
	var data int32
	if err := r.readData(4, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (r *reader) readInt64() (int64, error) {
	var data int64
	if err := r.readData(8, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (r *reader) readUint16() (uint16, error) {
	var data uint16
	if err := r.readData(2, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (r *reader) readUint32() (uint32, error) {
	var data uint32
	if err := r.readData(4, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (r *reader) readUint64() (uint64, error) {
	var data uint64
	if err := r.readData(8, &data); err != nil {
		return 0, err
	}

	return data, nil
}

// variable length integer
// ref. https://en.bitcoin.it/wiki/Protocol_documentation#Variable_length_integer
func (r *reader) readVarInt() (uint, error) {
	head, err := r.ReadByte()
	if err != nil {
		return 0, err
	}

	switch head {
	case 0xff:
		data, err := r.readUint64()
		if err != nil {
			return 0, err
		}
		return uint(data), nil
	case 0xfe:
		data, err := r.readUint32()
		if err != nil {
			return 0, err
		}
		return uint(data), nil
	case 0xfd:
		data, err := r.readUint16()
		if err != nil {
			return 0, err
		}
		return uint(data), nil
	default:
		return uint(head), nil
	}
}

func (r *reader) readOpCode() (OpCode, error) {
	b, err := r.ReadByte()
	if err != nil {
		return Op0, err
	}

	return OpCode(b), nil
}

func (r *reader) readPushedData(op OpCode) ([]byte, error) {
	var len uint

	switch op {
	case OpPushdata1:
		len8, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		len = uint(len8)
	case OpPushdata2:
		len16, err := r.readUint16()
		if err != nil {
			return nil, err
		}
		len = uint(len16)
	case OpPushdata4:
		len32, err := r.readUint32()
		if err != nil {
			return nil, err
		}
		len = uint(len32)
	default:
		// op: 0x01-0x4b (1-75)
		len = uint(op)
	}

	return r.readBytes(len)
}

func (r *reader) readScriptPart() ([]string, error) {
	op, err := r.readOpCode()
	if err != nil {
		return nil, err
	}

	if op.isPushData() {
		b, err := r.readPushedData(op)
		if err != nil {
			return nil, err
		}

		return []string{hex.EncodeToString(b)}, nil
	}

	return []string{op.Name()}, nil
}

func (r *reader) readScript() (*Script, error) {
	b := r.Bytes()
	asmParts := []string{}

	for {
		parts, err := r.readScriptPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		asmParts = append(asmParts, parts...)
	}

	return &Script{
		Hex: hex.EncodeToString(b),
		Asm: strings.Join(asmParts, " "),
	}, nil
}

func (r *reader) readTxVersion() (int32, error) {
	return r.readInt32()
}

func (r *reader) readTxIns() ([]*TxIn, error) {
	txInCnt, err := r.readVarInt()
	if err != nil {
		return nil, err
	}

	txIns := make([]*TxIn, txInCnt)
	for i := 0; i < int(txInCnt); i++ {
		txIn, err := r.readTxIn()
		if err != nil {
			return nil, err
		}

		txIns[i] = txIn
	}

	return txIns, nil
}

func (r *reader) readTxIn() (*TxIn, error) {
	txid, err := r.readHexReverse(32)
	if err != nil {
		return nil, err
	}

	index, err := r.readUint32()
	if err != nil {
		return nil, err
	}

	scriptLen, err := r.readVarInt()
	if err != nil {
		return nil, err
	}
	scriptHex, err := r.readHex(scriptLen)
	if err != nil {
		return nil, err
	}
	var script *Script
	if txid == CoinBaseTxid {
		script = &Script{
			Hex: scriptHex,
		}
	} else {
		script, err = NewScriptFromHex(scriptHex)
		if err != nil {
			return nil, err
		}
	}

	seq, err := r.readUint32()
	if err != nil {
		return nil, err
	}

	return &TxIn{
		Txid:     txid,
		Index:    index,
		Script:   script,
		Sequence: seq,
	}, nil
}

func (r *reader) readTxOuts() ([]*TxOut, error) {
	txOutCnt, err := r.readVarInt()
	if err != nil {
		return nil, err
	}

	txOuts := make([]*TxOut, txOutCnt)
	for i := 0; i < int(txOutCnt); i++ {
		txOut, err := r.readTxOut()
		if err != nil {
			return nil, err
		}

		txOuts[i] = txOut
	}

	return txOuts, nil
}

func (r *reader) readTxOut() (*TxOut, error) {
	amount, err := r.readInt64()
	if err != nil {
		return nil, err
	}

	scriptLen, err := r.readVarInt()
	if err != nil {
		return nil, err
	}
	scriptHex, err := r.readHex(scriptLen)
	if err != nil {
		return nil, err
	}
	script, err := NewScriptFromHex(scriptHex)
	if err != nil {
		return nil, err
	}

	return &TxOut{
		Amount: Satoshi(amount),
		Script: script,
	}, nil
}

func (r *reader) readLockTime() (uint32, error) {
	return r.readUint32()
}

func (r *reader) readTx() (*Tx, error) {
	version, err := r.readTxVersion()
	if err != nil {
		return nil, err
	}

	txIns, err := r.readTxIns()
	if err != nil {
		return nil, err
	}

	txOuts, err := r.readTxOuts()
	if err != nil {
		return nil, err
	}

	lockTime, err := r.readLockTime()
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

func (r *reader) readTxes() ([]*Tx, error) {
	cnt, err := r.readVarInt()
	if err != nil {
		return nil, err
	}

	txes := make([]*Tx, cnt)
	for i := uint(0); i < cnt; i++ {
		tx, err := r.readTx()
		if err != nil {
			return nil, err
		}

		txes[i] = tx
	}

	return txes, nil
}

func (r *reader) readBlockVersion() (int32, error) {
	return r.readInt32()
}

func (r *reader) readPrevBlock() (string, error) {
	return r.readHexReverse(32)
}

func (r *reader) readMerkleRoot() (string, error) {
	return r.readHexReverse(32)
}

func (r *reader) readTimestamp() (uint32, error) {
	return r.readUint32()
}

func (r *reader) readBits() (uint32, error) {
	return r.readUint32()
}

func (r *reader) readNonce() (uint32, error) {
	return r.readUint32()
}

func (r *reader) readBlock() (*Block, error) {
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

func (r *reader) readBlockHeader() (*BlockHeader, error) {
	version, err := r.readBlockVersion()
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
