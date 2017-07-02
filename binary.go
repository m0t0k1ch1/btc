package btc

import "crypto/sha256"

func reverseBytes(b []byte) []byte {
	l := len(b)

	rb := make([]byte, l)
	for i, j := 0, l-1; i < l; i++ {
		rb[i] = b[j]
		j--
	}

	return rb
}

func sha256Double(b []byte) ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write(b); err != nil {
		return nil, err
	}

	tmp := h.Sum(nil)

	h.Reset()
	if _, err := h.Write(tmp); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
