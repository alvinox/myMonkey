package object

import "fmt"

var Builtins = []struct {
    Name    string
    Builtin *Builtin
} {
    {"len",    &Builtin{Fn: BuiltinFuncLen},},
    {"first",  &Builtin{Fn: BuiltinFuncFirst},},
    {"last",   &Builtin{Fn: BuiltinFuncLast},},
    {"rest",   &Builtin{Fn: BuiltinFuncRest},},
    {"push",   &Builtin{Fn: BuiltinFuncPush},},
    {"puts",   &Builtin{Fn: BuiltinFuncPuts},},
}

func GetBuiltinByName(name string) *Builtin {
    for _, def := range Builtins {
        if def.Name == name {
            return def.Builtin
        }
    }

    return nil
}

func BuiltinFuncLen(args ...Object) Object {
    if len(args) != 1 {
        return newErrorObejct("wrong number of arguments. got=%d, want=1", len(args))
    }

    switch arg := args[0].(type) {
    case *String:
        return &Integer{Value: int64(len(arg.Value))}
    case *Array:
        return &Integer{Value: int64(len(arg.Elements))}
    default:
        return newErrorObejct("argument to `len` not supported, got %s", arg.Type())
    }
}

func BuiltinFuncFirst(args ...Object) Object {
    if len(args) != 1 {
        return newErrorObejct("wrong number of arguments. got=%d, want=1", len(args))
    }

    switch arg := args[0].(type) {
        case *Array:
            if len(arg.Elements) > 0 {
                return arg.Elements[0]
            } else {
                return nil
            }
        default:
            return newErrorObejct("argument to `first` must be ARRAY, got %s", arg.Type())
    }
}

func BuiltinFuncLast(args ...Object) Object {
    if len(args) != 1 {
        return newErrorObejct("wrong number of arguments. got=%d, want=1", len(args))
    }

    switch arg := args[0].(type) {
        case *Array:
            var length int = len(arg.Elements)
            if length > 0 {
                return arg.Elements[length - 1]
            } else {
                return nil
            }
        default:
            return newErrorObejct("argument to `last` must be ARRAY, got %s", arg.Type())
    }
}

func BuiltinFuncRest(args ...Object) Object {
    if len(args) != 1 {
        return newErrorObejct("wrong number of arguments. got=%d, want=1", len(args))
    }

    switch arg := args[0].(type) {
        case *Array:
            var length int = len(arg.Elements)
            if length > 0 {
                newElements := make([]Object, length - 1, length - 1)
                copy(newElements, arg.Elements[1:length])
                return &Array{Elements : newElements}
            } else {
                return nil
            }
        default:
            return newErrorObejct("argument to `rest` must be ARRAY, got %s", arg.Type())
    }
}

func BuiltinFuncPush(args ...Object) Object {
    if len(args) != 2 {
        return newErrorObejct("wrong number of arguments. got=%d, want=2", len(args))
    }

    switch arg := args[0].(type) {
        case *Array:
            var length int = len(arg.Elements)
            newElements := make([]Object, length + 1, length + 1)
            copy(newElements, arg.Elements)
            newElements[length] = args[1]

            return &Array{Elements : newElements}
        default:
            return newErrorObejct("argument to `push` must be ARRAY, got %s", arg.Type())
    }
}

func BuiltinFuncPuts(args ...Object) Object {
    for _, arg := range args {
        fmt.Println(arg.Inspect())
    }

    return nil 
}

func newErrorObejct(format string, a ...interface{}) *Error {
    return &Error{Message : fmt.Sprintf(format, a...)}
}
