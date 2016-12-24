package btctx

const (
	OP_0           OP = 0x00
	OP_FALSE       OP = 0x00
	OP_PUSHDATA1   OP = 0x4c
	OP_PUSHDATA2   OP = 0x4d
	OP_PUSHDATA4   OP = 0x4e
	OP_DROP        OP = 0x75
	OP_DUP         OP = 0x76
	OP_EQUALVERIFY OP = 0x88
	OP_HASH160     OP = 0xa9
	OP_CHECKSIG    OP = 0xac
)

var opCodeMap = map[OP]string{
	OP_0:           "OP_0",
	OP_PUSHDATA1:   "OP_PUSHDATA1",
	OP_PUSHDATA2:   "OP_PUSHDATA2",
	OP_PUSHDATA4:   "OP_PUSHDATA4",
	OP_DROP:        "OP_DROP",
	OP_DUP:         "OP_DUP",
	OP_EQUALVERIFY: "OP_EQUALVERIFY",
	OP_HASH160:     "OP_HASH160",
	OP_CHECKSIG:    "OP_CHECKSIG",
}

type OP byte

func (op OP) isPushData() bool {
	return (0x01 <= op && op <= 0x4b) ||
		op == OP_PUSHDATA1 ||
		op == OP_PUSHDATA2 ||
		op == OP_PUSHDATA4
}
