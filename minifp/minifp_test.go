package minifp_test

import (
	"log"
	"strings"
	"testing"

	"github.com/grailbio/testutil/expect"
	"github.com/yasushi-saito/minifp/minifp"
)

func compileRun(t *testing.T, expr string) minifp.Literal {
	exprs := minifp.Parse(strings.NewReader(expr))
	km := &minifp.KMachine{}
	var val minifp.Literal
	for _, expr := range exprs {
		log.Printf("Run: %+v", expr)
		val = km.Run(minifp.Compile(expr))
	}
	return val
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

func TestAssign0(t *testing.T) {
	expect.EQ(t, compileRun(t, `x = 10; x+x`).String(), "20")
}
