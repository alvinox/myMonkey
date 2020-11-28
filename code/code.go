package code

import (
    "fmt"
    "bytes"
    "encoding/binary"
)

type Instructions []byte

type Opcode byte

const (
    _ Opcode = iota
    OpConstant       // Constant
    OpPop            // pop the topmost element off the stack
    OpAdd            // Infix +   operator
    OpSub            // Infix -   operator
    OpMul            // Infix *   operator
    OpDiv            // Infix /   operator
    OpEqual          // Infix ==  operator
    OpNotEqual       // Infix !=  operator
    OpLessThan       // Infix <   operator
    OpGreaterThan    // Infix >   operator
    OpMinus          // Prefix -  operator
    OpBang           // Prefix !  operator
    OpTrue           // true
    OpFalse          // false
    OpNull           // null
    OpJumpNotTruthy
    OpJump
)

const (
    ConstWidth         int = 2 // const pool max size is 65536
    GlobalWidth        int = 2 // max number of global variables is 65536
    LocalWidth         int = 1 // max number of local variables is 256
    FreeWidth          int = 1 // max number of free variables is 256
    BuiltinWidth       int = 1 // max number of builtin functions is 256
    InstructionWidth   int = 2 // max number of instructions is 65536
    CallParamWidth     int = 1 // max number of parameters of each function is 256
)

type Definition struct {
    Name         string
    OperandWidths []int
}

var definitions = map[Opcode]*Definition {
    OpConstant:      {"OpConstant",      []int{2}},
    OpPop:           {"OpPop",           []int{}},
    OpAdd:           {"OpAdd",           []int{}},
    OpSub:           {"OpSub",           []int{}},
    OpMul:           {"OpMul",           []int{}},
    OpDiv:           {"OpDiv",           []int{}},
    OpEqual:         {"OpEqual",         []int{}},
    OpNotEqual:      {"OpNotEqual",      []int{}},
    OpLessThan:      {"OpLessThan",      []int{}},
    OpGreaterThan:   {"OpGreaterThan",   []int{}},
    OpMinus:         {"OpMinus",         []int{}},
    OpBang:          {"OpBang",          []int{}},
    OpTrue:          {"OpTrue",          []int{}},
    OpFalse:         {"OpFalse",         []int{}},
    OpNull:          {"OpNull",          []int{}},
    OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}},
    OpJump:          {"OpJump",          []int{2}},
}

func (ins Instructions) String() string {
    var out bytes.Buffer

    i := 0
    for i < len(ins) {
        def, err := Lookup(ins[i])
        if err != nil {
            fmt.Fprintf(&out, "ERROR: %s\n", err)
            continue
        }

        operands, n := ReadOperands(def, ins[i+1:])

        fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))

        i += 1 + n
    }

    return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
    operandCount := len(def.OperandWidths)

    if len(operands) != operandCount {
        return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n",
            len(operands), operandCount)
    }

    switch operandCount {
    case 0:
        return def.Name
    case 1:
        return fmt.Sprintf("%s %d", def.Name, operands[0])
    }

    return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

func Lookup(op byte) (*Definition, error) {
    def, ok := definitions[Opcode(op)]
    if !ok {
        return nil, fmt.Errorf("opcode %d undefined", op)
    }

    return def, nil
}

func Make(op Opcode, operands ...int) []byte {
    def, ok := definitions[op]
    if !ok {
        return []byte{}
    }

    instructionLen := 1
    for _, w := range def.OperandWidths {
        instructionLen += w
    }

    instructions := make([]byte, instructionLen)
    instructions[0] = byte(op)

    offset := 1
    for i, o := range operands {
        width := def.OperandWidths[i]
        switch width {
        case 2:
            binary.BigEndian.PutUint16(instructions[offset:], uint16(o))
        }
        offset += width
    }

    return instructions
}

func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
    operands := make([]int, len(def.OperandWidths))
    offset   := 0

    for i, width := range def.OperandWidths {
        switch width {
        case 2:
            operands[i] = int(ReadUint16(ins[offset:]))
        }

        offset += width
    }

    return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
    return binary.BigEndian.Uint16(ins)
}