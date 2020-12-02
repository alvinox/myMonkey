package evaluator

import (
    "fmt"
    "myMonkey/ast"
    "myMonkey/object"
)

var (
    NULL  = &object.Null{}
    TRUE  = &object.Boolean{Value: true}
    FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
    switch node := node.(type) {
    case *ast.Program:
        return evalProgram(node, env)

    case *ast.BlockStatement:
        return evalBlockStatement(node, env)
    
    case *ast.ExpressionStatement:
        return Eval(node.Expression, env)

    case *ast.ReturnStatement:
        val := Eval(node.ReturnValue, env)
        if isError(val) {
            return val
        }
        return &object.ReturnValue{Value: val}

    case *ast.LetStatement:
        val := Eval(node.Value, env)
        if isError(val) {
            return val
        }
        env.Set(node.Name.Value, val)

    case *ast.Identifier:
        return evalIdentifier(node, env)

    case *ast.IntegerLiteral:
        return &object.Integer{Value: node.Value}

    case *ast.StringLiteral:
        return &object.String{Value: node.Value}

    case *ast.Boolean:
        return nativeBoolToBooleanObject(node.Value)

    case *ast.ArrayLiteral:
        elements := evalExpressions(node.Elements, env)
        if len(elements) == 1 && isError(elements[0]) {
            return elements[0]
        }
        return &object.Array{Elements: elements}

    case *ast.HashLiteral:
        return evalHashLiteral(node, env)

    case *ast.PrefixOpExpression:
        right := Eval(node.Right, env)
        if isError(right) {
            return right
        }
        return evalPrefixOpExpression(node.Operator, right)

    case *ast.InfixExpression:
        left  := Eval(node.Left, env)
        if isError(left) {
            return left
        }

        right := Eval(node.Right, env)
        if isError(right) {
            return right
        }
        return evalInfixExpression(node.Operator, left, right)

    case *ast.IndexExpression:
        left  := Eval(node.Left, env)
        if isError(left) {
            return left
        }

        index := Eval(node.Index, env)
        if isError(index) {
            return index
        }

        return evalIndexExpression(left, index)

    case *ast.IfExpression:
        return evalIfExpression(node, env)

    case *ast.FunctionLiteral:
        params := node.Parameters
        body   := node.Body
        return &object.Function{
            Parameters: params,
            Body:       body,
            Env:        env,
        }
    case *ast.CallExpression:
        function := Eval(node.Function, env)
        if isError(function) {
            return function
        }

        args := evalExpressions(node.Arguments, env)
        if len(args) == 1 && isError(args[0]) {
            return args[0]
        }

        return applyFunction(function, args)
    }

    return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
    var result object.Object

    for _, stmt := range program.Statements {
        result = Eval(stmt, env)

        switch result := result.(type) {
        case *object.ReturnValue:
            return result.Value
        case *object.Error:
            return result
        }
    }

    return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
    var result object.Object

    for _, stmt := range block.Statements {
        result = Eval(stmt, env)

        if result != nil {
            rt := result.Type()
            if rt == object.RETURN_VALUE_OBJ ||
               rt == object.ERROR_OBJ {
                return result
            }
        }
    }

    return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
    if val, ok := env.Get(node.Value); ok {
        return val
    }

    if builtin, ok := Builtins[node.Value]; ok {
        return builtin
    }

    return newError("identifier not found: " + node.Value)
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
    pairs := make(map[object.HashKey]object.HashPair)

    for keyNode, valueNode := range node.Pairs {
        key := Eval(keyNode, env)
        if isError(key) {
            return key
        }

        hashKey, ok := key.(object.Hashable)
        if !ok {
            newError("unusable as hash key: %s", key.Type())
        }

        value := Eval(valueNode, env)
        if isError(value) {
            return value
        }

        hashed := hashKey.HashKey()
        pairs[hashed] = object.HashPair{Key: key, Value: value}
    }

    return &object.Hash{Pairs: pairs}
}

func evalPrefixOpExpression(operator string, right object.Object) object.Object {
    switch operator {
    case "!":
        return evalBangOperatorExpression(right)
    case "-":
        return evalMinusPrefixOperatorExpression(right)
    default:
        return newError("unknown operator %s%s", operator, right.Type())
    }
}

func evalInfixExpression(operator string, 
    left, right object.Object) object.Object {

    lType := left.Type()
    rType := right.Type()

    switch {
    case lType == object.INTEGER_OBJ && rType == object.INTEGER_OBJ:
        return evalIntegerInfixExpression(operator, left, right)
    case lType == object.STRING_OBJ && rType == object.STRING_OBJ:
        return evalStringInfixExpression(operator, left, right)
    case lType == object.BOOLEAN_OBJ && rType == object.BOOLEAN_OBJ:
        return evalBooleanInfixExpression(operator, left, right)
    case lType != rType:
        return newError("type mismatch: %s %s %s",
            lType, operator, rType)
    default:
        return newError("unknown operator: %s %s %s",
            lType, operator, rType)
    }
}

func evalIndexExpression(left, index object.Object) object.Object {
    lType := left.Type()
    iType := index.Type()

    switch {
    case lType == object.ARRAY_OBJ && iType == object.INTEGER_OBJ:
        return evalArrayIndexExpression(left, index)
    case lType == object.HASH_OBJ:
        return evalHashIndexExpression(left, index)
    default:
        return newError("index operator not supported: %s", left.Type())
    }
}

func evalArrayIndexExpression(left, index object.Object) object.Object {
    array := left.(*object.Array)
    idx   := index.(*object.Integer).Value
    max   := int64(len(array.Elements) - 1)

    if idx < 0 || idx > max {
        return NULL
    }

    return array.Elements[idx]
}

func evalHashIndexExpression(left, index object.Object) object.Object {
    hashObject := left.(*object.Hash)

    key, ok := index.(object.Hashable)
    if !ok {
        return newError("unusable as hash key: %s", index.Type())
    }

    pair, ok := hashObject.Pairs[key.HashKey()]
    if !ok {
        return NULL
    }

    return pair.Value
}

func evalBangOperatorExpression(right object.Object) object.Object {
    switch right {
    case TRUE:
        return FALSE
    case FALSE:
        return TRUE
    case NULL:
        return TRUE
    default:
        return FALSE
    }
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
    if right.Type() != object.INTEGER_OBJ {
        return newError("unknown operator: -%s", right.Type())
    }

    value := right.(*object.Integer).Value
    return &object.Integer{Value: -value}
}

func evalIntegerInfixExpression(operator string,
    left, right object.Object) object.Object {
    
    leftVal  := left.(*object.Integer).Value
    rightVal := right.(*object.Integer).Value
    switch operator {
    case "+":
        return &object.Integer{Value: leftVal + rightVal}
    case "-":
        return &object.Integer{Value: leftVal - rightVal}
    case "*":
        return &object.Integer{Value: leftVal * rightVal}
    case "/":
        return &object.Integer{Value: leftVal / rightVal}
    case "<":
        return nativeBoolToBooleanObject(leftVal < rightVal)
    case ">":
        return nativeBoolToBooleanObject(leftVal > rightVal)
    case "==":
        return nativeBoolToBooleanObject(leftVal == rightVal)
    case "!=":
        return nativeBoolToBooleanObject(leftVal != rightVal)
    default:
        return newError("unknown operator: %s %s %s",
            left.Type(), operator, right.Type())
    }
}

func evalStringInfixExpression(operator string,
    left, right object.Object) object.Object {
    
    leftVal  := left.(*object.String).Value
    rightVal := right.(*object.String).Value

    switch operator {
    case "+":
        return &object.String{Value: leftVal + rightVal}
    default:
        return newError("unknown operator: %s %s %s",
            left.Type(), operator, right.Type())
    }
}

func evalBooleanInfixExpression(operator string,
    left, right object.Object) object.Object {

    switch operator {
    case "==":
        return nativeBoolToBooleanObject(left == right)
    case "!=":
        return nativeBoolToBooleanObject(left != right)
    default:
        return newError("unknown operator: %s %s %s",
            left.Type(), operator, right.Type())
    }
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
    condition := Eval(ie.Condition, env)
    if isError(condition) {
        return condition
    }

    if isTruthy(condition) {
        return Eval(ie.Consequence, env)
    } else if ie.Alternative != nil {
        return Eval(ie.Alternative, env)
    } else {
        return NULL
    }
}

func evalExpressions(exps []ast.Expression, 
    env *object.Environment) []object.Object {
    var result []object.Object

    for _, e := range exps {
        evaluated := Eval(e, env)
        if isError(evaluated) {
            return []object.Object{evaluated}
        }

        result = append(result, evaluated)
    }

    return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
    switch fn := fn.(type) {
    case *object.Function:
        extendedEnv := extendFunctionEnv(fn, args)
        evaluated   := Eval(fn.Body, extendedEnv)
        return unwrapReturnValue(evaluated)
    case *object.Builtin:
        if result := fn.Fn(args...); result != nil {
            return result
        }
        return NULL
    default:
        return newError("not a function: %s", fn.Type())
    }
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
    env := object.NewEnclosedEnvironment(fn.Env)

    for idx, param := range fn.Parameters {
        env.Set(param.Value, args[idx])
    }

    return env
}

func unwrapReturnValue(obj object.Object) object.Object {
    if returnValue, ok := obj.(*object.ReturnValue); ok {
        return returnValue.Value
    }

    return obj
}

func newError(format string, a ...interface {}) *object.Error {
    return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func nativeBoolToBooleanObject(b bool) *object.Boolean {
    if b {
        return TRUE
    }

    return FALSE
}

func isTruthy(o object.Object) bool {
    switch o {
    case NULL:
        return false
    case TRUE:
        return true
    case FALSE:
        return false
    default:
        return true
    }
}

func isError(obj object.Object) bool {
    if obj != nil {
        return obj.Type() == object.ERROR_OBJ
    }
    return false
}
