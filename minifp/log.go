package minifp

import (
	"fmt"
	"text/scanner"
)

func panicf(pos scanner.Position, format string, args ...interface{}) {
	panic(pos.String() + ": " + fmt.Sprintf(format, args...))
}
