package minifp_test

import (
	"testing"

	"github.com/grailbio/testutil/expect"
	"github.com/yasushi-saito/minifp/minifp"
)

func compileRun(t *testing.T, node minifp.ASTNode) minifp.Literal {
	km := &minifp.KMachine{Code: minifp.Compile(node)}
	return km.Run()
}

func TestConst(t *testing.T) {
	x := minifp.InternSymbol("x")
	val := compileRun(t,
		&minifp.ASTApply{
			Head: &minifp.ASTLambda{
				Arg:  x,
				Body: minifp.ASTVar{Sym: x},
			},
			Tail: &minifp.ASTConst{Val: minifp.NewLiteralInt(10)}})
	expect.EQ(t, val.String(), "10")
}

func TestAdd(t *testing.T) {
	val := compileRun(t,
		&minifp.ASTApplyBuiltin{
			Op: minifp.BuiltinOpAdd,
			Args: []minifp.ASTNode{
				&minifp.ASTConst{Val: minifp.NewLiteralInt(10)},
				&minifp.ASTConst{Val: minifp.NewLiteralInt(11)}}})
	expect.EQ(t, val.String(), "21")
}
