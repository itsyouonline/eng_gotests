package encryption

import "crypto/md5"

// Hash returns the md5 hash
func Hash(data []byte) []byte {
	hash := md5.Sum(data)
	return hash[:]
}
