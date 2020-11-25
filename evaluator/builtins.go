package evaluator

import (
    "fmt"
    "myMonkey/object"
)

var Builtins = map[string]*object.Builtin {
    "len":     &object.Builtin{Fn: BuiltinFuncLen},
    "first":   &object.Builtin{Fn: BuiltinFuncFirst},
    "last":    &object.Builtin{Fn: BuiltinFuncLast},
    "rest":    &object.Builtin{Fn: BuiltinFuncRest},
    "push":    &object.Builtin{Fn: BuiltinFuncPush},
    "puts":    &object.Builtin{Fn: BuiltinFuncPuts},
}

func BuiltinFuncLen(args ...object.Object) object.Object {
    if len(args) != 1 {
        return newErrorObejct("wrong number of arguments. got=%d, want=1", len(args))
    }

    switch arg := args[0].(type) {
    case *object.String:
        return &object.Integer{Value: int64(len(arg.Value))}
    case *object.Array:
        return &object.Integer{Value: int64(len(arg.Elements))}
    default:
        return newErrorObejct("argument to `len` not supported, got %s", arg.Type())
    }
}

func BuiltinFuncFirst(args ...object.Object) object.Object {
    if len(args) != 1 {
        return newErrorObejct("wrong number of arguments. got=%d, want=1", len(args))
    }

    switch arg := args[0].(type) {
        case *object.Array:
            if len(arg.Elements) > 0 {
                return arg.Elements[0]
            } else {
                return NULL
            }
        default:
            return newErrorObejct("argument to `first` must be ARRAY, got %s", arg.Type())
    }
}

func BuiltinFuncLast(args ...object.Object) object.Object {
    if len(args) != 1 {
        return newErrorObejct("wrong number of arguments. got=%d, want=1", len(args))
    }

    switch arg := args[0].(type) {
        case *object.Array:
            var length int = len(arg.Elements)
            if length > 0 {
                return arg.Elements[length - 1]
            } else {
                return NULL
            }
        default:
            return newErrorObejct("argument to `last` must be ARRAY, got %s", arg.Type())
    }
}

func BuiltinFuncRest(args ...object.Object) object.Object {
    if len(args) != 1 {
        return newErrorObejct("wrong number of arguments. got=%d, want=1", len(args))
    }

    switch arg := args[0].(type) {
        case *object.Array:
            var length int = len(arg.Elements)
            if length > 0 {
                newElements := make([]object.Object, length - 1, length - 1)
                copy(newElements, arg.Elements[1:length])
                return &object.Array{Elements : newElements}
            } else {
                return NULL
            }
        default:
            return newErrorObejct("argument to `rest` must be ARRAY, got %s", arg.Type())
    }
}

func BuiltinFuncPush(args ...object.Object) object.Object {
    if len(args) != 2 {
        return newErrorObejct("wrong number of arguments. got=%d, want=2", len(args))
    }

    switch arg := args[0].(type) {
        case *object.Array:
            var length int = len(arg.Elements)
            newElements := make([]object.Object, length + 1, length + 1)
            copy(newElements, arg.Elements)
            newElements[length] = args[1]

            return &object.Array{Elements : newElements}
        default:
            return newErrorObejct("argument to `push` must be ARRAY, got %s", arg.Type())
    }
}

func BuiltinFuncPuts(args ...object.Object) object.Object {
    for _, arg := range args {
        fmt.Println(arg.Inspect())
    }

    return NULL
}


func newErrorObejct(format string, a ...interface{}) *object.Error {
    return &object.Error{Message : fmt.Sprintf(format, a...)}
}
