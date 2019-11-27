package minifp_test

import (
	"log"
	"strings"
	"testing"

	"github.com/grailbio/testutil/expect"
	"github.com/yasushi-saito/minifp/minifp"
)

func run(t *testing.T, km *minifp.KMachine, expr string) minifp.Literal {
	exprs := minifp.Parse(strings.NewReader(expr))
	var val minifp.Literal
	for _, expr := range exprs {
		log.Printf("Run: %+v", expr)
		val = km.Run(km.Compile(expr))
	}
	return val
}

func TestConst(t *testing.T) {
	km := minifp.NewMachine()
	val := run(t, km, "10")
	expect.EQ(t, val.String(), "10")
}

func TestIDFunction(t *testing.T) {
	km := minifp.NewMachine()
	val := run(t, km, `(\x -> x) 10`)
	expect.EQ(t, val.String(), "10")
}

func TestFunction2(t *testing.T) {
	km := minifp.NewMachine()
	val := run(t, km, `(\x -> x*2) 10`)
	expect.EQ(t, val.String(), "20")
}

func TestBuiltinBinaryOp(t *testing.T) {
	km := minifp.NewMachine()
	expect.EQ(t, run(t, km, `10+11`).String(), "21")
	expect.EQ(t, run(t, km, `11-10`).String(), "1")
	expect.EQ(t, run(t, km, `11-10-1`).String(), "0")
	expect.EQ(t, run(t, km, `10*11`).String(), "110")
}

func TestBuiltinPred(t *testing.T) {
	km := minifp.NewMachine()
	expect.EQ(t, run(t, km, `10>11`).String(), "false")
	expect.EQ(t, run(t, km, `10>=11`).String(), "false")
	expect.EQ(t, run(t, km, `10<11`).String(), "true")
	expect.EQ(t, run(t, km, `10<=11`).String(), "true")
	expect.EQ(t, run(t, km, `10<=10`).String(), "true")
	expect.EQ(t, run(t, km, `10>=10`).String(), "true")
	expect.EQ(t, run(t, km, `10>10`).String(), "false")
	expect.EQ(t, run(t, km, `10<10`).String(), "false")
	expect.EQ(t, run(t, km, `10==10`).String(), "true")
	expect.EQ(t, run(t, km, `10!=10`).String(), "false")
}

func TestIf(t *testing.T) {
	km := minifp.NewMachine()
	expect.EQ(t, run(t, km, `if (10==10) 1 2`).String(), "1")
	expect.EQ(t, run(t, km, `if (10!=10) 1 2`).String(), "2")
}

func TestAssign0(t *testing.T) {
	km := minifp.NewMachine()
	expect.EQ(t, run(t, km, `x = 10; x+x`).String(), "20")
	expect.EQ(t, run(t, km, `y = x; y+11`).String(), "21")
	expect.EQ(t, run(t, km, `x = 20`).String(), "20")
}

func TestLetrec(t *testing.T) {
	km := minifp.NewMachine()
	expect.EQ(t, run(t, km, `letrec x=10 in x*x`).String(), "100")
	expect.EQ(t, run(t, km, `letrec x=10; y=x+1 in x*y`).String(), "110")
	expect.EQ(t, run(t, km, `letrec x=10 in (letrec y=12 in x*y)`).String(), "120")
}
