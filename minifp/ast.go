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

type compiler struct {
	globals []Symbol
	stack   []Symbol
}

func (c *compiler) lookup(pos scanner.Position, sym Symbol) KAddr {
	for i, name := range c.stack {
		if name == sym {
			return KAddr{index: uint32(i)}
		}
	}
	for i, name := range c.globals {
		if name == sym {
			return KAddr{index: uint32(i), global: true}
		}
	}
	panicf(pos, "variable %v not found in %+v %+v", sym, c.stack, c.globals)
	panic("")
}

func (c *compiler) compile(node ASTNode) KCode {
	switch v := node.(type) {
	case *ASTAssign:
		c.globals = append(c.globals, v.Sym)
		return &KAssign{Addr: c.lookup(v.pos, v.Sym), Expr: c.compile(v.Expr)}
	case *ASTConst:
		return (*KConst)(&v.Val)
	case *ASTLambda:
		c.stack = append(c.stack, v.Arg)
		defer func() { c.stack = c.stack[:len(c.stack)-1] }()
		return &KLambda{Body: c.compile(v.Body)}
	case *ASTVar:
		addr := c.lookup(v.pos, v.Sym)
		if addr.global {
			return &KGlobalVar{Addr: addr}
		} else {
			return &KLocalVar{Addr: addr}
		}
	case *ASTApply:
		return &KApply{Head: c.compile(v.Head), Tail: c.compile(v.Tail)}
	case *ASTApplyBuiltin:
		if len(v.Args) == 1 {
			return &KApply{Head: c.compile(v.Args[0]), Tail: &KBuiltinOp{v.Op}}
		}
		if len(v.Args) == 2 {
			c0 := &KApply{Head: c.compile(v.Args[0]), Tail: &KSwapStack{1}}
			c1 := &KApply{Head: c0, Tail: c.compile(v.Args[1])}
			return &KApply{Head: c1, Tail: &KBuiltinOp{v.Op}}
		}
		panic("blah")
	}
	panic(node)
}

func Compile(node ASTNode) KCode {
	var c compiler
	return c.compile(node)
}
