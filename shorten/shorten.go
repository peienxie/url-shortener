package shorten

import (
	"crypto/sha256"

	"github.com/jxskiss/base62"
)

const intMaxWidth = 20

// ShortenByHash encodes given long URL string into a shorten URL string by
// apply SHA256 hashing and base62 encoding with given long URL
func ShortenByHash(longURL string) string {
	digest := sha256.Sum256([]byte(longURL))
	return base62.EncodeToString(digest[:])
}

func itoaLeadingZeros(num uint64, width int) []byte {
	b := []byte{}
	for num > 0 {
		b = append(b, byte((num%10)+'0'))
		num /= 10
	}
	for width > len(b) {
		b = append(b, byte('0'))
	}

	// reverse the output
	for i := 0; i < len(b)/2; i++ {
		b[i], b[len(b)-1-i] = b[len(b)-1-i], b[i]
	}

	return b
}

// ShortenByInt encodes given integer into a base62 string,
// the given integer will be converted to a numeric string padded with leading zeros first
func ShortenByInt(num uint64) string {
	b := itoaLeadingZeros(num, intMaxWidth)
	return base62.EncodeToString(b)
}
