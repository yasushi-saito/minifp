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
	BuiltinOpSub                   = (2 << 16) | 2
	BuiltinOpMul                   = (3 << 16) | 2
	BuiltinOpGE                    = (4 << 16) | 2
	BuiltinOpLE                    = (5 << 16) | 2
	BuiltinOpGT                    = (6 << 16) | 2
	BuiltinOpLT                    = (7 << 16) | 2
	BuiltinOpEQ                    = (8 << 16) | 2
	BuiltinOpNEQ                   = (9 << 16) | 2
)

func (op BuiltinOpType) String() string {
	switch op {
	case BuiltinOpAdd:
		return "builtin:+"
	case BuiltinOpSub:
		return "builtin:-"
	case BuiltinOpMul:
		return "builtin:*"
	case BuiltinOpGE:
		return "builtin:>="
	case BuiltinOpLE:
		return "builtin:<="
	case BuiltinOpGT:
		return "builtin:>"
	case BuiltinOpLT:
		return "builtin:<"
	case BuiltinOpEQ:
		return "builtin:=="
	case BuiltinOpNEQ:
		return "builtin:!="
	}
	return fmt.Sprintf("builtin:%d", op)
}
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
	buf.WriteRune('(')
	buf.WriteString(fmt.Sprint(n.Op))
	for _, arg := range n.Args {
		buf.WriteRune(' ')
		buf.WriteString(arg.String())
	}
	buf.WriteRune(')')
	return buf.String()
}

type ASTLetrec struct {
	pos      scanner.Position
	Bindings []*ASTAssign
	Body     ASTNode
}

func (n ASTLetrec) Pos() scanner.Position { return n.pos }
func (n ASTLetrec) String() string {
	return fmt.Sprintf("letrec %+v in %v", n.Bindings, n.Body)
}

type ASTIf struct {
	pos              scanner.Position
	Cond, Then, Else ASTNode
}

func (n ASTIf) Pos() scanner.Position { return n.pos }
func (n ASTIf) String() string {
	return fmt.Sprintf("if %v %v %v", n.Cond, n.Then, n.Else)
}
