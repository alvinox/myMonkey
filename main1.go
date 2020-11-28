package main

import (
    "fmt"
    "myMonkey/lexer"
    "myMonkey/parser"
    "myMonkey/object"
    "myMonkey/evaluator"
    "myMonkey/compiler"
    "myMonkey/vm"
)

func main() {
    // evaluate()
    machine()
}

func evaluate() {
    line    := `"Hello" + " " + "World!"`
    l       := lexer.New(line)
    p       := parser.New(l)
    program := p.ParseProgram()
    env     := object.NewEnvironment()

    evaluated := evaluator.Eval(program, env)
    fmt.Printf(evaluated.Inspect())
}

func machine() {
    line     := `if (true) { 10 }; 3333;`
    l        := lexer.New(line)
    p        := parser.New(l)
    program  := p.ParseProgram()

    compiler := compiler.New()
    compiler.Compile(program)

    machine := vm.New(compiler.Bytecode())
    machine.Run()

    lastPopped := machine.LastPoppedStackElem()
    fmt.Printf(lastPopped.Inspect())
}