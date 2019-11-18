package minifp

//go:generate goyacc -l -o parser_generated.go parser.y

import (
	"errors"
	"io"
	"strconv"
	"text/scanner"
)

func Parse(in io.Reader) []ASTNode {
	p := parser{
		sc: &scanner.Scanner{},
	}
	p.sc.Error = func(_ *scanner.Scanner, msg string) {
		if p.err == nil {
			p.err = errors.New(msg)
		}
	}
	p.sc.Mode = scanner.GoTokens
	p.sc.Init(in)

	yyParse(&p)
	if p.err != nil {
		panic(p.err)
	}
	return p.result
}

type parser struct {
	err    error
	result []ASTNode
	sc     *scanner.Scanner
}

// Lex implements yyLexer
func (p *parser) Error(msg string) {
	if p.err == nil {
		p.err = errors.New(msg)
	}
}

// Lex implements yyLexer
func (p *parser) Lex(y *yySymType) int {
	if p.err != nil {
		return 0
	}
	ch := p.sc.Scan()
	if p.err != nil {
		return 0
	}
	if ch == scanner.EOF {
		return 0
	}
	if ch == scanner.Int {
		val, err := strconv.ParseInt(p.sc.TokenText(), 0, 64)
		if err != nil {
			panicf(p.sc.Pos(), "parse int %s: %s", p.sc.TokenText(), err)
		}
		y.ast = &ASTConst{Val: Literal{typ: LiteralInt, intVal: val}}
		return tokLiteral
	}
	if ch == scanner.Ident {
		y.ident = p.sc.TokenText()
		return tokIdent
	}
	if ch == '=' || ch == '(' || ch == ')' || ch == '+' ||
		ch == '*' || ch == '\\' || ch == ';' {
		return int(ch)
	}
	if ch == '-' {
		ch2 := p.sc.Peek()
		if ch2 == '>' {
			p.sc.Next()
			return tokArrow
		}
		return int(ch)
	}
	panic(ch)
	return 1
}

func newLambda(pos scanner.Position, args []string, expr ASTNode) ASTNode {
	arg := InternSymbol(args[0])
	if len(args) == 1 {
		return &ASTLambda{Arg: arg, Body: expr}
	}
	return &ASTLambda{pos: pos, Arg: arg, Body: newLambda(pos, args[1:], expr)}
}
