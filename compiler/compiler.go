package compiler

import (
    "fmt"
    "myMonkey/ast"
    "myMonkey/code"
    "myMonkey/object"
)

type EmittedInstruction struct {
    Opcode    code.Opcode
    Position  int
}

type Compiler struct {
    instructions code.Instructions
    constants    []object.Object

    lastInstruction     EmittedInstruction
    previousInstruction EmittedInstruction
}

func New() *Compiler {
    return &Compiler {
        instructions: code.Instructions{},
        constants:    []object.Object{},

        lastInstruction:     EmittedInstruction{},
        previousInstruction: EmittedInstruction{},
    }
}

func (c *Compiler) Compile(node ast.Node) error {
    switch node := node.(type) {
    case *ast.Program:
        for _, s := range node.Statements {
            err := c.Compile(s)
            if err != nil {
                return err
            }
        }

    case *ast.BlockStatement:
        for _, s := range node.Statements {
            err := c.Compile(s)
            if err != nil {
                return err
            }
        }

    case *ast.ExpressionStatement:
        err := c.Compile(node.Expression)
        if err != nil {
            return err
        }

        // the object is useless after ExpressionStatement
        c.emit(code.OpPop)

    case *ast.PrefixOpExpression:
        err := c.Compile(node.Right)
        if err != nil {
            return err
        }

        switch node.Operator {
        case "-":
            c.emit(code.OpMinus)
        case "!":
            c.emit(code.OpBang)
        default:
            return fmt.Errorf("unknown operator %s", node.Operator)
        }

    case *ast.InfixExpression:
        err := c.Compile(node.Left)
        if err != nil {
            return err
        }

        err = c.Compile(node.Right)
        if err != nil {
            return err
        }

        switch node.Operator {
        case "+":
            c.emit(code.OpAdd)
        case "-":
            c.emit(code.OpSub)
        case "*":
            c.emit(code.OpMul)
        case "/":
            c.emit(code.OpDiv)
        case "==":
            c.emit(code.OpEqual)
        case "!=":
            c.emit(code.OpNotEqual)
        case "<":
            c.emit(code.OpLessThan)
        case ">":
            c.emit(code.OpGreaterThan)
        default:
            return fmt.Errorf("unknown operator %s", node.Operator)
        }

    case *ast.IntegerLiteral:
        integer := &object.Integer{Value: node.Value}
        c.emit(code.OpConstant, c.addConstant(integer))

    case *ast.Boolean:
        if node.Value {
            c.emit(code.OpTrue)
        } else {
            c.emit(code.OpFalse)
        }
        
    case *ast.IfExpression:
        err := c.Compile(node.Condition)
        if err != nil {
            return err
        }

        jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

        err = c.Compile(node.Consequence)
        if err != nil {
            return err
        }

        if c.lastInstructionIsPop() {
            c.removeLastPop()
        }

        jumpPos := c.emit(code.OpJump, 9999)

        afterConsequencePos := len(c.instructions)
        c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

        if node.Alternative == nil {
            c.emit(code.OpNull)
        } else {
            err = c.Compile(node.Alternative)
            if err != nil {
                return err
            }
    
            if c.lastInstructionIsPop() {
                c.removeLastPop()
            }
        }

        afterAlternativePos := len(c.instructions)
        c.changeOperand(jumpPos, afterAlternativePos)
    }

    return nil
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
    ins := code.Make(op, operands...)
    pos := c.addInstruction(ins)

    c.setLastInstruction(op, pos)

    return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
    pos := len(c.instructions)
    c.instructions = append(c.instructions, ins...)
    return pos
}

func (c *Compiler) addConstant(object object.Object) int {
    c.constants = append(c.constants, object)
    return len(c.constants) - 1
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
    previous := c.lastInstruction
    last     := EmittedInstruction{Opcode: op, Position: pos}

    c.previousInstruction = previous
    c.lastInstruction     = last
}

func (c *Compiler) lastInstructionIsPop() bool {
    return c.lastInstruction.Opcode == code.OpPop
}

func (c *Compiler) removeLastPop() {
    c.instructions    = c.instructions[:c.lastInstruction.Position]
    c.lastInstruction = c.previousInstruction
}

func (c *Compiler) replaceInstruction(pos int, newInstructions []byte) {
    for i := 0; i < len(newInstructions); i++ {
        c.instructions[pos + i] = newInstructions[i]
    }
}

func (c *Compiler) changeOperand(opPos int, operand int) {
    op := code.Opcode(c.instructions[opPos])
    newInstructions := code.Make(op, operand)

    c.replaceInstruction(opPos, newInstructions)
}

type Bytecode struct {
    Instructions code.Instructions
    Constants    []object.Object
}

func (c *Compiler) Bytecode() *Bytecode {
    return &Bytecode {
        Instructions: c.instructions,
        Constants:    c.constants,
    }
}