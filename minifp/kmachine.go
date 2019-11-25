package minifp

import (
	"fmt"
	"log"
	"strings"
	"text/scanner"
)

type KCode interface {
	DebugString() string
}

type LiteralType uint

const (
	LiteralInvalid LiteralType = iota
	LiteralInt
	LiteralNil
)

type Literal struct {
	typ    LiteralType
	intVal int64
}

var kNil = Literal{typ: LiteralNil}

func NewLiteralInt(v int64) Literal { return Literal{typ: LiteralInt, intVal: int64(v)} }

func (l Literal) String() string {
	switch l.typ {
	case LiteralInt:
		return fmt.Sprint(l.intVal)
	}
	panic(l)
}

func (l Literal) Int() int64 {
	if l.typ != LiteralInt {
		panic(l)
	}
	return l.intVal
}

type KLambda struct {
	Arg  Symbol
	Body KCode
}

func (k *KLambda) DebugString() string {
	return "ƛ"
}

type KLetrec struct {
	VarNames []Symbol
	VarExprs []KCode
	Body     KCode
}

func (k *KLetrec) DebugString() string {
	return fmt.Sprintf("letrec %v %v %v", k.VarNames, k.VarExprs, k.Body.DebugString())
}

type KApply struct{ Head, Tail KCode }

func (k *KApply) DebugString() string {
	return fmt.Sprintf("(%s %s)", k.Head.DebugString(), k.Tail.DebugString())
}

type KVar struct{ Addr KAddr }

func (k *KVar) DebugString() string {
	return fmt.Sprintf("localvar:%v", k.Addr)
}

type KRet struct{}

var kRet = &KRet{}

func (k *KRet) DebugString() string {
	return "ret"
}

type KConst Literal

func (k *KConst) DebugString() string {
	return fmt.Sprintf("const:%+v", (*Literal)(k).String())
}

type KBuiltinOp struct{ Op BuiltinOpType }

func (k *KBuiltinOp) DebugString() string {
	return fmt.Sprintf("builtin:%+v", k.Op)
}

type KSwapStack struct{ N int }

func (k *KSwapStack) DebugString() string {
	return fmt.Sprintf("swapstack:%+v", k.N)
}

type kStackEntry struct {
	cl      KClosure
	pointer *KClosure
}

func (s kStackEntry) String() string {
	if s.pointer != nil {
		return "*"
	}
	return s.cl.String()
}

type kVarEntry struct {
	// sym is the name of the variable that's stored in this slot.  It's mainly
	// for debugging; runtime access is done using index.
	sym Symbol
	cl  KClosure
}

type kEnvFrame struct {
	Const *Literal
	vars  []kVarEntry
	next  *kEnvFrame
}

// func (s kEnvFrame) Empty() bool { return s.head == nil }

// func (s kEnvFrame) Pop() (cl KClosure, s2 kEnvFrame) {
// 	if s.head == nil {
// 		panic(s)
// 	}
// 	cl = s.head.KClosure
// 	s2.head = s.head.link
// 	return cl, s2
// }

func NewMachine() *KMachine {
	return &KMachine{}
}

func (k *KMachine) Read(addr KAddr) *KClosure {
	if addr.frameIndex == kGlobalFrame {
		return &k.Globals[addr.varIndex].cl
	}
	f := k.Locals
	for addr.frameIndex > 0 {
		f = f.next
		addr.frameIndex--
	}
	if f == nil {
		panic(k)
	}
	if f.Const != nil {
		panic(f)
	}
	return &f.vars[addr.varIndex].cl
}

func (s *kEnvFrame) String() string {
	var (
		buf strings.Builder
		i   int
		f   = s
	)
	buf.WriteRune('[')
	for f != nil {
		if i > 0 {
			buf.WriteRune(' ')
		}
		if s.Const != nil {
			buf.WriteString(fmt.Sprintf("const:%+v", *s.Const))
		} else {
			buf.WriteRune('[')
			for j, e := range f.vars {
				if j > 0 {
					buf.WriteRune(' ')
				}
				buf.WriteString(*e.sym.string)
				buf.WriteString("=")
				buf.WriteString(e.cl.String())
			}
			buf.WriteRune(']')
		}
		f = f.next
		i++
	}
	buf.WriteRune(']')
	return buf.String()
}

type KClosure struct {
	Code KCode
	Env  *kEnvFrame
}

func (cl KClosure) String() string {
	var buf strings.Builder
	if cl.Code == kRet {
		buf.WriteString("ret:")
		buf.WriteString(cl.Env.Const.String())
	} else {
		buf.WriteString(cl.Code.DebugString())
	}
	return buf.String()
}

func (cl KClosure) Literal() Literal {
	if cl.Code != kRet {
		panic(cl)
	}
	if cl.Env.Const == nil {
		panic(cl)
	}
	return *cl.Env.Const
}

const kGlobalFrame = uint32(0xffffffff)

type KAddr struct {
	frameIndex uint32
	varIndex   uint32
}

type KMachine struct {
	Code    KCode
	Globals []kVarEntry
	Locals  *kEnvFrame
	Stack   []kStackEntry
	step    int
}

func (k *KMachine) Run(code KCode) Literal {
	k.Code = code
	for k.Step() {
	}
	if k.Code != kRet {
		panic(fmt.Sprintf("%d: %v %v %v", k.step, k.Code.DebugString(), k.Locals.String(), k.Stack))
	}
	return *k.Locals.Const
}

func (k *KMachine) pushStack(e kStackEntry) {
	k.Stack = append(k.Stack, e)
}

func (k *KMachine) popStack() kStackEntry {
	n := len(k.Stack)
	v := k.Stack[n-1]
	k.Stack = k.Stack[:n-1]
	return v
}

func (k *KMachine) Step() bool {
	k.step++
	log.Printf("%d: %v %v %v", k.step, k.Code.DebugString(), k.Locals.String(), k.Stack)
	switch v := k.Code.(type) {
	case *KApply:
		k.Code = v.Head
		k.Stack = append(k.Stack, kStackEntry{cl: KClosure{Code: v.Tail, Env: k.Locals}})
	case *KVar:
		cl := k.Read(v.Addr)
		k.Code = cl.Code
		k.Locals = cl.Env
		k.pushStack(kStackEntry{pointer: cl})
	case *KLambda:
		k.Code = v.Body
		arg := k.popStack()
		if arg.pointer != nil {
			panic(arg)
		}
		k.Locals = &kEnvFrame{
			vars: []kVarEntry{{sym: v.Arg, cl: arg.cl}},
			next: k.Locals}
	case *KLetrec:
		frame := &kEnvFrame{vars: make([]kVarEntry, len(v.VarExprs)), next: k.Locals}
		for i, b := range v.VarExprs {
			frame.vars[i] = kVarEntry{
				sym: v.VarNames[i],
				cl:  KClosure{Code: b, Env: frame},
			}
		}
		k.Code = v.Body
		k.Locals = frame
	case *KConst:
		k.Code = kRet
		k.Locals = &kEnvFrame{Const: (*Literal)(v)}
	case *KRet:
		if k.Locals.Const == nil {
			panic(k)
		}
		if len(k.Stack) == 0 {
			return false
		}
		val := k.Locals
		top := k.popStack()
		for top.pointer != nil {
			*top.pointer = KClosure{Code: kRet, Env: k.Locals}
			if len(k.Stack) == 0 {
				return false
			}
			top = k.popStack()
		}
		k.Code = top.cl.Code
		k.Locals = top.cl.Env
		k.pushStack(kStackEntry{cl: KClosure{Code: kRet, Env: val}})
	case *KBuiltinOp:
		switch v.Op.NArg() {
		case 2:
			arg1 := k.popStack()
			arg0 := k.popStack()
			v0, v1 := arg0.cl.Literal(), arg1.cl.Literal()
			k.Code = kRet

			var val Literal
			switch v.Op {
			case BuiltinOpAdd:
				val = NewLiteralInt(v0.Int() + v1.Int())
			case BuiltinOpSub:
				val = NewLiteralInt(v0.Int() - v1.Int())
			case BuiltinOpMul:
				val = NewLiteralInt(v0.Int() * v1.Int())
			default:
				panic(v)
			}
			k.Locals = &kEnvFrame{Const: &val}
		}
	case *KSwapStack:
		if v.N != 1 {
			panic(v)
		}
		arg0, arg1, ret := k.popStack(), k.popStack(), k.popStack()
		if arg0.pointer != nil || arg0.cl.Code != kRet {
			panic(arg0)
		}
		if arg1.pointer != nil {
			panic(arg1)
		}
		k.Code = arg1.cl.Code
		k.Locals = arg1.cl.Env
		k.pushStack(arg0)
		k.pushStack(ret)
	default:
		return false
	}
	return true
}

func (k *KMachine) Compile(node ASTNode) KCode {
	var (
		c = compiler{globals: &k.Globals}
	)
	return c.compile(node)
}

type compiler struct {
	// Points to KMachine.Globals
	globals *[]kVarEntry
	locals  [][]Symbol
}

func (c *compiler) lookup(pos scanner.Position, sym Symbol) (addr KAddr, ok bool) {
	for i, frame := range c.locals {
		for j, name := range frame {
			if sym == name {
				return KAddr{frameIndex: uint32(i), varIndex: uint32(j)}, true
			}
		}
	}
	for j, e := range *c.globals {
		if sym == e.sym {
			return KAddr{frameIndex: kGlobalFrame, varIndex: uint32(j)}, true
		}
	}
	return KAddr{}, false
}

func (c *compiler) compile(node ASTNode) KCode {
	switch v := node.(type) {
	case *ASTAssign:
		addr, ok := c.lookup(v.pos, v.Sym)
		cl := KClosure{Code: c.compile(v.Expr), Env: nil}
		if !ok {
			*c.globals = append(
				*c.globals,
				kVarEntry{sym: v.Sym, cl: cl})
		} else {
			if addr.frameIndex != kGlobalFrame {
				panicf(v.pos, "local variable found where global is expected:%v", v.Sym)
			}
			(*c.globals)[addr.varIndex].cl = cl
		}
		return cl.Code
	case *ASTConst:
		return (*KConst)(&v.Val)
	case *ASTLambda:
		mustf(v.pos, v.Arg.string != nil, "v:%v", v)
		c.locals = append(c.locals, []Symbol{v.Arg})
		defer func() { c.locals = c.locals[:len(c.locals)-1] }()
		return &KLambda{Arg: v.Arg, Body: c.compile(v.Body)}
	case *ASTVar:
		addr, ok := c.lookup(v.pos, v.Sym)
		if !ok {
			panicf(v.pos, "variable %v not found in %+v", v.Sym, c.locals)
		}
		return &KVar{Addr: addr}
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
	case *ASTLetrec:
		n := len(v.Bindings)
		var frame []Symbol
		for _, b := range v.Bindings {
			frame = append(frame, b.Sym)
		}
		c.locals = append(c.locals, frame)
		var (
			varNames = make([]Symbol, 0, n)
			varExprs = make([]KCode, 0, n)
		)
		for _, b := range v.Bindings {
			varNames = append(varNames, b.Sym)
			varExprs = append(varExprs, c.compile(b.Expr))
		}
		defer func() { c.locals = c.locals[:len(c.locals)-1] }()
		return &KLetrec{VarNames: varNames, VarExprs: varExprs, Body: c.compile(v.Body)}

	}
	panic(node)
}
