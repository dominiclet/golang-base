package randgenerate

import (
	cryptoRand "crypto/rand"
	"encoding/hex"
	"math/rand"
)

var runes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

func GenerateAlphaNumericString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}

func GenerateSecureToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := cryptoRand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
