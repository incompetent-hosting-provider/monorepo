package util

import (
	"flag"
	"os"
	"strings"
)

func IsTestRun() bool {
	if flag.Lookup("test.v") != nil || strings.HasSuffix(os.Args[0], ".test") {
		return true
	}
	return false
}
