package utils

import (
	"math/rand"
	"time"
)

var chars string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateName() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	mn := ""
	for i := 0; i < 10; i++ {
		index := r.Intn(len(chars))
		mn += string(chars[index])
	}
	return mn
}
