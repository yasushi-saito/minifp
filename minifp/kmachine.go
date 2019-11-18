package minifp

import (
	"fmt"
	"log"
	"strings"
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

type KLambda struct{ Body KCode }

func (k *KLambda) DebugString() string {
	return "ƛ"
}

type KAssign struct {
	Addr KAddr
	Expr KCode
}

func (k *KAssign) DebugString() string {
	return fmt.Sprintf("var:%v=%s", k.Addr, k.Expr.DebugString())
}

type KApply struct{ Head, Tail KCode }

func (k *KApply) DebugString() string {
	return fmt.Sprintf("(%s %s)", k.Head.DebugString(), k.Tail.DebugString())
}

type KLocalVar struct{ Addr KAddr }
type KGlobalVar struct{ Addr KAddr }

func (k *KLocalVar) DebugString() string {
	return fmt.Sprintf("localvar:%v", k.Addr)
}

func (k *KGlobalVar) DebugString() string {
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

type KEnv struct {
	kHeapMap
	Const *Literal
}

func (s KEnv) String() string {
	if s.Const != nil {
		return fmt.Sprintf("const:%+v", *s.Const)
	}
	return s.kHeapMap.String()
}

type kHeapEntry struct {
	KClosure
	link *kHeapEntry
}

type kHeapMap struct {
	head *kHeapEntry
}

func (s kHeapMap) Push(cl KClosure) (s2 kHeapMap) {
	s2.head = &kHeapEntry{KClosure: cl, link: s.head}
	return s2
}

func (s kHeapMap) Empty() bool { return s.head == nil }

func (s kHeapMap) Pop() (cl KClosure, s2 kHeapMap) {
	if s.head == nil {
		panic(s)
	}
	cl = s.head.KClosure
	s2.head = s.head.link
	return cl, s2
}

func (s kHeapMap) Read(i uint32) KClosure {
	cl := s.head
	for i > 0 {
		cl = cl.link
		i--
	}
	if cl == nil {
		panic(s)
	}
	return cl.KClosure
}

func (s kHeapMap) String() string {
	var (
		buf strings.Builder
		i   int
		cl  = s.head
	)
	buf.WriteRune('[')
	for cl != nil {
		if i > 0 {
			buf.WriteRune(' ')
		}
		i++
		buf.WriteString(cl.KClosure.String())
		cl = cl.link
	}
	buf.WriteRune(']')
	return buf.String()
}

type KClosure struct {
	Code KCode
	Env  KEnv
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

type KAddr struct {
	index  uint32
	global bool
}

type KMachine struct {
	Code    KCode
	Env     KEnv
	Globals kHeapMap
	Stack   kHeapMap
	step    int
}

func (k *KMachine) Run(code KCode) Literal {
	if v, ok := code.(*KAssign); ok {
		k.Globals = k.Globals.Push(KClosure{v.Expr, KEnv{}})
		log.Printf("Done: %v", k.Env.String())
		return kNil
	}
	k.Code = code
	for k.Step() {
	}
	if k.Code != kRet {
		panic(fmt.Sprintf("%d: %v %v %v", k.step, k.Code.DebugString(), k.Env.String(), k.Stack.String()))
	}
	return *k.Env.Const
}

func (k *KMachine) Step() bool {
	k.step++
	log.Printf("%d: %v %v %v", k.step, k.Code.DebugString(), k.Env.String(), k.Stack.String())
	switch v := k.Code.(type) {
	case *KApply:
		k.Code = v.Head
		k.Stack = k.Stack.Push(KClosure{Code: v.Tail, Env: k.Env})
	case *KLocalVar:
		cl := k.Env.Read(v.Addr.index)
		k.Code = cl.Code
		k.Env = cl.Env
	case *KGlobalVar:
		cl := k.Globals.Read(v.Addr.index)
		k.Code = cl.Code
		k.Env = cl.Env
	case *KLambda:
		k.Code = v.Body
		var arg KClosure
		arg, k.Stack = k.Stack.Pop()
		k.Env = KEnv{kHeapMap: k.Env.Push(arg)}
	case *KConst:
		k.Code = kRet
		k.Env = KEnv{Const: (*Literal)(v)}
	case *KRet:
		if k.Env.Const == nil {
			panic(k)
		}
		if k.Stack.Empty() {
			return false
		}
		val := k.Env
		top, stack := k.Stack.Pop()
		k.Code = top.Code
		k.Env = top.Env
		k.Stack = stack.Push(KClosure{Code: kRet, Env: val})
	case *KBuiltinOp:
		switch v.Op.NArg() {
		case 2:
			arg1, stack := k.Stack.Pop()
			arg0, stack := stack.Pop()
			k.Stack = stack
			v0, v1 := arg0.Literal(), arg1.Literal()
			k.Code = kRet

			var val Literal
			switch v.Op {
			case BuiltinOpAdd:
				val = NewLiteralInt(v0.Int() + v1.Int())
			case BuiltinOpMul:
				val = NewLiteralInt(v0.Int() * v1.Int())
			default:
				panic(v)
			}
			k.Env = KEnv{Const: &val}
		}
	case *KSwapStack:
		if v.N != 1 {
			panic(v)
		}
		arg0, stack := k.Stack.Pop()
		if arg0.Code != kRet {
			panic(arg0)
		}
		arg1, stack := stack.Pop()
		ret, stack := stack.Pop()
		k.Code = arg1.Code
		k.Env = arg1.Env
		k.Stack = stack.Push(arg0).Push(ret)
	default:
		return false
	}
	return true
}
