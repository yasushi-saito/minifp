package minifp

import (
	"fmt"
	"text/scanner"
)

var kUnknownPos scanner.Position

func panicf(pos scanner.Position, format string, args ...interface{}) {
	panic(pos.String() + ": " + fmt.Sprintf(format, args...))
}

func mustf(pos scanner.Position, cond bool, format string, args ...interface{}) {
	if !cond {
		panicf(pos, format, args...)
	}
}
