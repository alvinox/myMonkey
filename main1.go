package main

import (
    "fmt"
    "myMonkey/lexer"
    "myMonkey/parser"
    "myMonkey/object"
    "myMonkey/evaluator"
)

func main() {
    line    := `"Hello" + " " + "World!"`
    l       := lexer.New(line)
    p       := parser.New(l)
    program := p.ParseProgram()
    env     := object.NewEnvironment()

    evaluated := evaluator.Eval(program, env)
    fmt.Printf(evaluated.Inspect())
}