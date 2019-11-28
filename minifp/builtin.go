package minifp

type funcSpec struct {
	name string
	nArg int
	cb   func(args ...Literal) Literal
}

var funcs map[string]*funcSpec

func init() {
	funcs = map[string]*funcSpec{
		"builtin:+": &funcSpec{
			name: "builtin:+",
			nArg: 2,
			cb: func(args ...Literal) Literal {
				return NewLiteralInt(args[0].Int() + args[1].Int())
			},
		},
		"builtin:-": &funcSpec{
			name: "builtin:-",
			nArg: 2,
			cb: func(args ...Literal) Literal {
				return NewLiteralInt(args[0].Int() - args[1].Int())
			},
		},
		"builtin:*": &funcSpec{
			name: "builtin:*",
			nArg: 2,
			cb: func(args ...Literal) Literal {
				return NewLiteralInt(args[0].Int() * args[1].Int())
			},
		},
		"builtin:==": &funcSpec{
			name: "builtin:==",
			nArg: 2,
			cb: func(args ...Literal) (val Literal) {
				val = kFalse
				if args[0].Int() == args[1].Int() {
					val = kTrue
				}
				return
			},
		},
		"builtin:!=": &funcSpec{
			name: "builtin:!=",
			nArg: 2,
			cb: func(args ...Literal) (val Literal) {
				val = kFalse
				if args[0].Int() != args[1].Int() {
					val = kTrue
				}
				return
			},
		},
		"builtin:>=": &funcSpec{
			name: "builtin:>=",
			nArg: 2,
			cb: func(args ...Literal) (val Literal) {
				val = kFalse
				if args[0].Int() >= args[1].Int() {
					val = kTrue
				}
				return
			},
		},
		"builtin:<=": &funcSpec{
			name: "builtin:<=",
			nArg: 2,
			cb: func(args ...Literal) (val Literal) {
				val = kFalse
				if args[0].Int() <= args[1].Int() {
					val = kTrue
				}
				return
			},
		},
		"builtin:<": &funcSpec{
			name: "builtin:<",
			nArg: 2,
			cb: func(args ...Literal) (val Literal) {
				val = kFalse
				if args[0].Int() < args[1].Int() {
					val = kTrue
				}
				return
			},
		},
		"builtin:>": &funcSpec{
			name: "builtin:>",
			nArg: 2,
			cb: func(args ...Literal) (val Literal) {
				val = kFalse
				if args[0].Int() > args[1].Int() {
					val = kTrue
				}
				return
			},
		},
	}
}
