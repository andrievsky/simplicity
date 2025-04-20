package genid

import (
	"github.com/google/uuid"
	"math/rand"
)

func GenerateUUID() string {
	return uuid.New().String()
}

const dict = "abcdefghijklmnopqrstuvwxyz0123456789"

func GeneratePartialID(n int) string {
	s := make([]byte, n)
	for i := range s {
		s[i] = dict[rand.Intn(len(dict))]
	}
	return string(s)
}
