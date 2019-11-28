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

type ASTApplyLeafFunction struct {
	pos  scanner.Position
	Op   *funcSpec
	Args []ASTNode
}

func (n ASTApplyLeafFunction) Pos() scanner.Position { return n.pos }

func (n ASTApplyLeafFunction) String() string {
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
