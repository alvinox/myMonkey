package compiler

type SymbolScope string

const (
    GlobalScope SymbolScope = "GLOBAL"
)

// holds all the necessary information about a symbol
type Symbol struct {
    Name    string
    Scope   SymbolScope
    Index   int
}

// associates string s with Symbol s in its store 
// and keeps track of the numDefinitions it has
type SymbolTable struct {
    store          map[string]Symbol
    numDefinitions int
}

func NewSymbolTable() *SymbolTable {
    s := make(map[string]Symbol)
    return &SymbolTable{store: s}
}

func (s *SymbolTable) Define(name string) Symbol {
    symbol := Symbol {
        Name:   name,
        Scope:  GlobalScope,
        Index:  s.numDefinitions,
    }

    s.store[name] = symbol
    s.numDefinitions++

    return symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
    obj, ok := s.store[name]
    return obj, ok
}
