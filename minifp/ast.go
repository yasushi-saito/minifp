package minifp

import (
	"fmt"
	"strings"
	"text/scanner"
)

type ASTNode interface {
	// scanner.Position returns the source-code location of this node.
	Pos() scanner.Position
	String() string
}

type ASTConst struct {
	pos scanner.Position
	Val Literal
}

func (n ASTConst) Pos() scanner.Position { return n.pos }
func (n ASTConst) String() string        { return n.Val.String() }

type ASTVar struct {
	pos scanner.Position
	Sym Symbol
}

func (n ASTVar) Pos() scanner.Position { return n.pos }
func (n ASTVar) String() string        { return n.Sym.String() }

type ASTApply struct {
	pos  scanner.Position
	Head ASTNode
	Tail ASTNode
}

func (n ASTApply) Pos() scanner.Position { return n.pos }
func (n ASTApply) String() string        { return "(" + n.Head.String() + " " + n.Tail.String() + ")" }

type ASTLambda struct {
	pos  scanner.Position
	Arg  Symbol
	Body ASTNode
}

func (n ASTLambda) Pos() scanner.Position { return n.pos }
func (n ASTLambda) String() string        { return "ƛ" + n.Arg.String() + "." + n.Body.String() }

type ASTAssign struct {
	pos  scanner.Position
	Sym  Symbol
	Expr ASTNode
}

func (n ASTAssign) Pos() scanner.Position { return n.pos }
func (n ASTAssign) String() string {
	return n.Sym.String() + "=" + n.Expr.String()
}

type BuiltinOpType uint

const (
	BuiltinOpInvalid BuiltinOpType = iota
	BuiltinOpAdd                   = (1 << 16) | 2
	BuiltinOpMul                   = (2 << 16) | 2
)

func (o BuiltinOpType) NArg() int {
	return int(o & 0xffff)
}

type ASTApplyBuiltin struct {
	pos  scanner.Position
	Op   BuiltinOpType
	Args []ASTNode
}

func (n ASTApplyBuiltin) Pos() scanner.Position { return n.pos }

func (n ASTApplyBuiltin) String() string {
	var buf strings.Builder
	buf.WriteString("(builtin:")
	buf.WriteString(fmt.Sprint(n.Op))
	for _, arg := range n.Args {
		buf.WriteRune(' ')
		buf.WriteString(arg.String())
	}
	buf.WriteRune(')')
	return buf.String()
}
