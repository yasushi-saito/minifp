package minifp

import "text/scanner"

type ASTNode interface {
	// scanner.Position returns the source-code location of this node.
	Pos() scanner.Position
}

type ASTConst struct {
	pos scanner.Position
	Val Literal
}

func (n ASTConst) Pos() scanner.Position { return n.pos }

type ASTVar struct {
	pos scanner.Position
	Sym Symbol
}

func (n ASTVar) Pos() scanner.Position { return n.pos }

type ASTApply struct {
	pos  scanner.Position
	Head ASTNode
	Tail ASTNode
}

func (n ASTApply) Pos() scanner.Position { return n.pos }

type ASTLambda struct {
	pos  scanner.Position
	Arg  Symbol
	Body ASTNode
}

func (n ASTLambda) Pos() scanner.Position { return n.pos }

type ASTAssign struct {
	pos  scanner.Position
	Sym  Symbol
	Expr ASTNode
}

func (n ASTAssign) Pos() scanner.Position { return n.pos }

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

type compiler struct {
	stack []Symbol
}

func (c *compiler) compile(node ASTNode) KCode {
	switch v := node.(type) {
	case *ASTConst:
		return (*KConst)(&v.Val)
	case *ASTLambda:
		c.stack = append(c.stack, v.Arg)
		return &KLambda{Body: c.compile(v.Body)}
		c.stack = c.stack[:len(c.stack)-1]
	case *ASTVar:
		for i, name := range c.stack {
			if name == v.Sym {
				return &KVar{Index: i}
			}
		}
		panicf(node.Pos(), "variable %v not found in %+v", v.Sym, c.stack)
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
