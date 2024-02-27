package util

import (
	"flag"
	"math/rand"
	"os"
	"strings"
)

func IsTestRun() bool {
	if flag.Lookup("test.v") != nil || strings.HasSuffix(os.Args[0], ".test") {
		return true
	}
	return false
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
