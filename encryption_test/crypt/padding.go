package crypt

// AddPadding adds padding to the input according to the PKCS7 specification.
// Padding is always added, even if the last block is also a full block.
func AddPadding(input []byte, blockSize int) []byte {
	paddingSize := blockSize - len(input)%blockSize
	padding := make([]byte, paddingSize)
	for i := range padding {
		padding[i] = byte(paddingSize)
	}
	return append(input, padding...)
}

// RemovePadding removes padding from the input. Padding is expected to be according
// ot the PKCS7 specification. Padding is always removed (so the output is alwasys
// shorter than the input)
func RemovePadding(input []byte) []byte {
	paddingSize := int(input[len(input)-1])
	return input[:len(input)-paddingSize]
}
