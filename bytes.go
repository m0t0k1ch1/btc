package btc

func reverseBytes(b []byte) []byte {
	l := len(b)

	rb := make([]byte, l)
	for i, j := 0, l-1; i < l; i++ {
		rb[i] = b[j]
		j--
	}

	return rb
}
