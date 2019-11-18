package minifp

type SourcePos struct {
	Path string
	Line int
}

type ASTNode interface {
	// SourcePos returns the source-code location of this node.
	SourcePos() SourcePos
}

type ASTConst struct {
	pos SourcePos
	Val Literal
}

func (n ASTConst) SourcePos() SourcePos { return n.pos }

type ASTVar struct {
	pos SourcePos
	Sym Symbol
}

func (n ASTVar) SourcePos() SourcePos { return n.pos }

type ASTApply struct {
	pos  SourcePos
	Head ASTNode
	Tail ASTNode
}

func (n ASTApply) SourcePos() SourcePos { return n.pos }

type ASTLambda struct {
	pos  SourcePos
	Arg  Symbol
	Body ASTNode
}

func (n ASTLambda) SourcePos() SourcePos { return n.pos }

type BuiltinOpType uint

const (
	BuiltinOpInvalid BuiltinOpType = iota
	BuiltinOpAdd
)

type ASTApplyBuiltin struct {
	pos  SourcePos
	Op   BuiltinOpType
	Args []ASTNode
}

func (n ASTApplyBuiltin) SourcePos() SourcePos { return n.pos }

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
	case ASTVar:
		for i, name := range c.stack {
			if name == v.Sym {
				return &KVar{Index: i}
			}
		}
		panic(v)
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
