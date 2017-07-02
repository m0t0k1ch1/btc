package btc

const (
	OpZero        OpCode = 0x00
	OpFalse       OpCode = 0x00
	OpDataLenMin  OpCode = 0x01
	OpDataLenMax  OpCode = 0x4b
	OpPushdata1   OpCode = 0x4c
	OpPushdata2   OpCode = 0x4d
	OpPushdata4   OpCode = 0x4e
	OpReturn      OpCode = 0x6a
	OpDrop        OpCode = 0x75
	OpDup         OpCode = 0x76
	OpEqualVerify OpCode = 0x88
	OpHash160     OpCode = 0xa9
	OpCheckSig    OpCode = 0xac
)

var opCodeNameMap = map[OpCode]string{
	OpZero:        "OP_0",
	OpPushdata1:   "OP_PUSHDATA1",
	OpPushdata2:   "OP_PUSHDATA2",
	OpPushdata4:   "OP_PUSHDATA4",
	OpReturn:      "OP_RETURN",
	OpDrop:        "OP_DROP",
	OpDup:         "OP_DUP",
	OpEqualVerify: "OP_EQUALVERIFY",
	OpHash160:     "OP_HASH160",
	OpCheckSig:    "OP_CHECKSIG",
}

type OpCode byte

func (op OpCode) Name() string {
	return opCodeNameMap[op]
}

func (op OpCode) Byte() byte {
	return byte(op)
}

func (op OpCode) isDataLen() bool {
	return OpDataLenMin <= op && op <= OpDataLenMax
}

func (op OpCode) isPushData() bool {
	return op.isDataLen() ||
		op == OpPushdata1 ||
		op == OpPushdata2 ||
		op == OpPushdata4
}
