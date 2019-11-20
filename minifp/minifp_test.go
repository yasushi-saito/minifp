package minifp_test

import (
	"log"
	"strings"
	"testing"

	"github.com/grailbio/testutil/expect"
	"github.com/yasushi-saito/minifp/minifp"
)

func compileRun(t *testing.T, km *minifp.KMachine, expr string) minifp.Literal {
	exprs := minifp.Parse(strings.NewReader(expr))
	var val minifp.Literal
	for _, expr := range exprs {
		log.Printf("Run: %+v", expr)
		val = km.Run(km.Compile(expr))
	}
	return val
}

func TestConst(t *testing.T) {
	km := &minifp.KMachine{}
	val := compileRun(t, km, "10")
	expect.EQ(t, val.String(), "10")
}

func TestIDFunction(t *testing.T) {
	km := &minifp.KMachine{}
	val := compileRun(t, km, `(\x -> x) 10`)
	expect.EQ(t, val.String(), "10")
}

func TestFunction2(t *testing.T) {
	km := &minifp.KMachine{}
	val := compileRun(t, km, `(\x -> x*2) 10`)
	expect.EQ(t, val.String(), "20")
}

func TestBuiltinBinaryOp(t *testing.T) {
	km := &minifp.KMachine{}
	val := compileRun(t, km, `10+11`)
	expect.EQ(t, val.String(), "21")
	val = compileRun(t, km, `10*11`)
	expect.EQ(t, val.String(), "110")
}

func TestAssign0(t *testing.T) {
	km := &minifp.KMachine{}
	expect.EQ(t, compileRun(t, km, `x = 10; x+x`).String(), "20")
	expect.EQ(t, compileRun(t, km, `y = x; y+11`).String(), "21")
	expect.EQ(t, compileRun(t, km, `x = 20`).String(), "20")
}
