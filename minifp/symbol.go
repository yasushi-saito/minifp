package minifp

import (
	"sync"
)

type Symbol struct{ *string }

var (
	symMu    sync.Mutex
	symTable = map[string]Symbol{}
)

func (s Symbol) String() string {
	return *s.string
}

func InternSymbol(name string) Symbol {
	symMu.Lock()
	s, ok := symTable[name]
	if !ok {
		s := Symbol{&name}
		symTable[name] = s
	}
	symMu.Unlock()
	return s
}
