package object

import (
    "fmt"
    "bytes"
    "strings"
    "hash/fnv"
    "myMonkey/ast"
)

type ObjectType string

const (
    NULL_OBJ         = "NULL"
    ERROR_OBJ        = "ERROR"
    INTEGER_OBJ      = "INTEGER"
    STRING_OBJ       = "STRING"
    BOOLEAN_OBJ      = "BOOLEAN"
    ARRAY_OBJ        = "ARRAY"
    HASH_OBJ         = "HASH"
    RETURN_VALUE_OBJ = "RETURN_VALU"
    FUNCTION_OBJ     = "FUNCTION"
    BUILTIN_OBJ      = "BUILTIN"
)

type Object interface {
    Type() ObjectType
    Inspect() string
}

type Null struct { }

func (n *Null) Type() ObjectType { return NULL_OBJ }

func (n *Null) Inspect() string { return "null" }

type Error struct {
    Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }

func (e *Error) Inspect() string { return "ERROR:" + e.Message }

type Integer struct {
    Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

func (i *Integer) HashKey() HashKey {
    return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type String struct {
    Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }

func (s *String) Inspect() string { return s.Value }

func (s *String) HashKey() HashKey {
    h := fnv.New64a()
    h.Write([]byte(s.Value))
    return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type Boolean struct {
    Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

func (b *Boolean) HashKey() HashKey {
    var val uint64

    if b.Value {
        val = 1
    } else {
        val = 0
    }

    return HashKey{Type: b.Type(), Value: val}
}

type Array struct {
    Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }

func (a *Array) Inspect() string { 
    var out bytes.Buffer

    elements := []string{}
    for _, e := range a.Elements {
        elements = append(elements, e.Inspect())
    }

    out.WriteString("[")
    out.WriteString(strings.Join(elements, ", "))
    out.WriteString("]")

    return out.String()
}

type Hash struct {
    Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }

func (h *Hash) Inspect() string {
    var out bytes.Buffer

    pairs := []string{}
    for _, pair := range h.Pairs {
        pairs = append(pairs, pair.Key.Inspect() + ":" + pair.Value.Inspect())
    }

    out.WriteString("{")
    out.WriteString(strings.Join(pairs, ", "))
    out.WriteString("}")

    return out.String()
}

type ReturnValue struct {
    Value Object
}

func (r *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

func (r *ReturnValue) Inspect() string { return r.Value.Inspect() }

type Function struct {
    Parameters []*ast.Identifier
    Body       *ast.BlockStatement
    Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }

func (f *Function) Inspect() string {
    var out bytes.Buffer

    params := []string{}
    for _, p := range f.Parameters {
        params = append(params, p.String())
    }

    out.WriteString("fn(")
    out.WriteString(strings.Join(params, ", "))
    out.WriteString(")")
    out.WriteString(f.Body.String())

    return out.String()
}

type BuiltinFunction func(args ...Object) Object
type Builtin struct {
    Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }

func (b *Builtin) Inspect() string { return "builtin function" }

type Hashable interface {
    HashKey() HashKey
}

type HashKey struct {
    Type   ObjectType
    Value  uint64
}

type HashPair struct {
    Key    Object
    Value  Object
}