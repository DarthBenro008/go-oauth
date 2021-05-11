package utils

import (
	"crypto/rand"
	"fmt"
)

func GenerateRandomString() (string, error) {
	n := 5
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	s := fmt.Sprintf("%X", b)
	return s, nil
}
