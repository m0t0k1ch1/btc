package btctx

import "encoding/hex"

type Script struct {
	Hex       string   `json:"hex"`
	Asm       string   `json:"asm"`
	Addresses []string `json:"addresses"`
	Data      [][]byte `json:"-"`
}

func NewScriptFromHex(s string) (*Script, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return NewScriptFromBytes(b)
}

func NewScriptFromBytes(b []byte) (*Script, error) {
	sd := newScriptDecoder(b)

	return sd.decode()
}

type scriptParts []string

func (sps scriptParts) likeP2PKH() bool {
	return len(sps) >= 5 &&
		sps[0] == opCodeMap[OP_DUP] &&
		sps[1] == opCodeMap[OP_HASH160] &&
		len(sps[2]) == 40 &&
		sps[3] == opCodeMap[OP_EQUALVERIFY] &&
		sps[4] == opCodeMap[OP_CHECKSIG]
}

func (sps scriptParts) extractAddresses() ([]string, error) {
	if sps.likeP2PKH() {
		pkh := PubKeyHash(sps[2])

		address, err := pkh.ToAddress()
		if err != nil {
			return nil, err
		}

		return []string{address}, nil
	}

	return []string{}, nil
}
