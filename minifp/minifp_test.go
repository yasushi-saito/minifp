package minifp_test

import (
	"strings"
	"testing"

	"github.com/grailbio/testutil/assert"
	"github.com/grailbio/testutil/expect"
	"github.com/yasushi-saito/minifp/minifp"
)

func compileRun(t *testing.T, expr string) minifp.Literal {
	exprs := minifp.Parse(strings.NewReader(expr))
	assert.EQ(t, len(exprs), 1)
	km := &minifp.KMachine{Code: minifp.Compile(exprs[0])}
	return km.Run()
}

func TestConst(t *testing.T) {
	val := compileRun(t, "10")
	expect.EQ(t, val.String(), "10")
}

func TestIDFunction(t *testing.T) {
	val := compileRun(t, `(\x -> x) 10`)
	expect.EQ(t, val.String(), "10")
}

func TestFunction2(t *testing.T) {
	val := compileRun(t, `(\x -> x*2) 10`)
	expect.EQ(t, val.String(), "20")
}

func TestBuiltinBinaryOp(t *testing.T) {
	val := compileRun(t, `10+11`)
	expect.EQ(t, val.String(), "21")
	val = compileRun(t, `10*11`)
	expect.EQ(t, val.String(), "110")
}
