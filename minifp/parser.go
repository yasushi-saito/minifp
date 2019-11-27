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
		sc:  &scanner.Scanner{},
		ops: map[byte]*opTrieNode{},
	}
	p.addOp("-", '-')
	p.addOp("+", '+')
	p.addOp("*", '*')
	p.addOp("(", '(')
	p.addOp(")", ')')
	p.addOp(";", ';')
	p.addOp("\\", '\\')
	p.addOp("=", '=')
	p.addOp("==", tokEQ)
	p.addOp("!=", tokNEQ)
	p.addOp("->", tokArrow)
	p.addOp(">", '>')
	p.addOp("<", '<')
	p.addOp(">=", tokGE)
	p.addOp("<=", tokLE)
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

type opTrieNode struct {
	tok int
	ch2 map[rune]int
}

type parser struct {
	err    error
	ops    map[byte]*opTrieNode
	result []ASTNode
	sc     *scanner.Scanner
}

func (p *parser) addOp(op string, tok int) {
	n := p.ops[op[0]]
	if n == nil {
		n = &opTrieNode{}
		p.ops[op[0]] = n
	}
	if len(op) == 1 {
		n.tok = tok
		return
	}
	if len(op) == 2 {
		if n.ch2 == nil {
			n.ch2 = map[rune]int{}
		}
		n.ch2[rune(op[1])] = tok
		return
	}
	panic(op)
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
		switch p.sc.TokenText() {
		case "letrec":
			return tokLetrec
		case "in":
			return tokIn
		case "if":
			return tokIf
		default:
			return tokIdent
		}
	}
	if e, ok := p.ops[byte(ch)]; ok {
		if e.ch2 != nil {
			ch2 := p.sc.Peek()
			if tok, ok := e.ch2[ch2]; ok {
				p.sc.Next()
				return tok
			}
		}
		if e.tok == 0 {
			panicf(p.sc.Pos(), "invalid char '%c'", ch)
		}
		return e.tok
	}
	// if ch == '-' {
	// 	ch2 := p.sc.Peek()
	// 	if ch2 == '>' {
	// 		p.sc.Next()
	// 		return tokArrow
	// 	}
	// 	return int(ch)
	// }
	// if ch == '>' {
	// 	ch2 := p.sc.Peek()
	// 	if ch2 == '=' {
	// 		p.sc.Next()
	// 		return tokGE
	// 	}
	// 	return int(ch)
	// }
	// if ch == '=' || ch == '(' || ch == ')' || ch == '+' ||
	// 	ch == '*' || ch == '\\' || ch == ';' {
	// 	return int(ch)
	// }
	panicf(p.sc.Pos(), "invalid char '%c'", ch)
	return 1
}

func newLambda(pos scanner.Position, args []string, expr ASTNode) ASTNode {
	arg := InternSymbol(args[0])
	if len(args) == 1 {
		return &ASTLambda{Arg: arg, Body: expr}
	}
	return &ASTLambda{pos: pos, Arg: arg, Body: newLambda(pos, args[1:], expr)}
}
