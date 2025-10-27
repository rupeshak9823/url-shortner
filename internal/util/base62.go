package util

import (
	"crypto/sha256"
	"math/big"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func padLeft(s string, length int) string {
	for len(s) < length {
		s = "0" + s
	}
	return s
}

func StringToBase62(s string, length int) string {
	hash := sha256.Sum256([]byte(s))
	num := new(big.Int).SetBytes(hash[:])

	encoded := toBase62(num)

	// Pad or truncate to desired length
	if len(encoded) < 8 {
		encoded = padLeft(encoded, length)
	} else if len(encoded) > length {
		encoded = encoded[:length]
	}

	return encoded
}

func toBase62(num *big.Int) string {
	if num.Sign() == 0 {
		return "0"
	}

	base := big.NewInt(62)
	mod := new(big.Int)
	result := ""

	for num.Sign() > 0 {
		num.DivMod(num, base, mod)
		result = string(base62Chars[mod.Int64()]) + result
	}

	return result
}
