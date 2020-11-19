package ast

import (
    "bytes"
    "strings"
    "myMonkey/token"
)

type Node interface {
    TokenLiteral() string
    String()       string
}

type Statement interface {
    Node
    StatementNode()
}

type Expression interface {
    Node
    ExpressionNode()
}

type Program struct {
    Statements []Statement
}

func (p *Program) TokenLiteral() string {
    if len(p.Statements) > 0 {
        return p.Statements[0].TokenLiteral()
    } else {
        return ""
    }
}

func (p *Program) String() string {
    var out bytes.Buffer

    for _, s := range p.Statements {
        out.WriteString(s.String())
    }

    return out.String()
}

type LetStatement struct {
    Token   token.Token // the 'let' token
    Name    *Identifier
    Value   Expression
}

func (ls *LetStatement) TokenLiteral() string {
    return ls.Token.Literal
}

func (ls *LetStatement) String() string {
    var out bytes.Buffer

    out.WriteString(ls.TokenLiteral() + " ")
    out.WriteString(ls.Name.String())
    out.WriteString(" = ")

    if ls.Value != nil {
        out.WriteString(ls.Value.String())
    }

    out.WriteString(";")

    return out.String()
}

func (ls *LetStatement) StatementNode() {}

type ReturnStatement struct {
    Token       token.Token // the 'return' token
    ReturnValue Expression
}

func (rs *ReturnStatement) TokenLiteral() string {
    return rs.Token.Literal
}

func (rs *ReturnStatement) StatementNode() {}

func (rs *ReturnStatement) String() string {
    var out bytes.Buffer

    out.WriteString(rs.TokenLiteral() + " ")

    if rs.ReturnValue != nil {
        out.WriteString(rs.ReturnValue.String())
    }

    out.WriteString(";")

    return out.String()
}

type ExpressionStatement struct {
    Token       token.Token // the first token of the expression
    Expression  Expression
}

func (es *ExpressionStatement) TokenLiteral() string {
    return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
    var out bytes.Buffer

    if es.Expression != nil {
        out.WriteString(es.Expression.String())
    }

    return out.String()
}

func (es *ExpressionStatement) StatementNode() {}


type BlockStatement struct {
    Token      token.Token // the token.LBRACE token
    Statements []Statement
}

func (bs *BlockStatement) TokenLiteral() string {
    return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
    var out bytes.Buffer

    out.WriteString(" {")
    for _, stmt := range bs.Statements {
        out.WriteString(stmt.String())
    }
    out.WriteString("} ")

    return out.String()
}

func (bs *BlockStatement) StatementNode() {}

type Identifier struct {
    Token    token.Token // the token.IDENT token
    Value    string
}

func (i *Identifier) TokenLiteral() string {
    return i.Token.Literal
}

func (i *Identifier) String() string {
    return i.Value
}

func (i *Identifier) ExpressionNode() {}

type IntegerLiteral struct {
    Token    token.Token // the token.INT token
    Value    int64
}

func (i *IntegerLiteral) TokenLiteral() string {
    return i.Token.Literal
}

func (i *IntegerLiteral) String() string {
    return i.TokenLiteral()
}

func (i *IntegerLiteral) ExpressionNode() {}

type Boolean struct {
    Token    token.Token // the token.INT token
    Value    bool
}

func (b *Boolean) TokenLiteral() string {
    return b.Token.Literal
}

func (b *Boolean) String() string {
    return b.TokenLiteral()
}

func (b *Boolean) ExpressionNode() {}

type IfExpression struct {
    Token         token.Token // the token.IF token
    Condition     Expression
    Consequence   *BlockStatement
    Alternative   *BlockStatement
}

func (ie *IfExpression) TokenLiteral() string {
    return ie.Token.Literal
}

func (ie *IfExpression) String() string {
    var out bytes.Buffer

    out.WriteString("if")
    out.WriteString(ie.Condition.String())
    out.WriteString(" ")
    out.WriteString(ie.Consequence.String())

    if ie.Alternative != nil {
        out.WriteString("else")
        out.WriteString(ie.Alternative.String())
    }

    return out.String()
}

func (ie *IfExpression) ExpressionNode() {}

type FunctionLiteral struct {
    Token      token.Token // the token.FUNCTION token
    Parameters []*Identifier
    Body       *BlockStatement
}

func (fl *FunctionLiteral) TokenLiteral() string {
    return fl.Token.Literal
}

func (fl *FunctionLiteral) String() string {
    var out bytes.Buffer

    params := []string{}
    for _, p := range fl.Parameters {
        params = append(params, p.String())
    }

    out.WriteString("fn(")
    out.WriteString(strings.Join(params, ", "))
    out.WriteString(")")
    out.WriteString(fl.Body.String())

    return out.String()
}

func (fl *FunctionLiteral) ExpressionNode() {}

type CallExpression struct {
    Token      token.Token // the '(' token
    Function   Expression
    Arguments  []Expression
}

func (ce *CallExpression) TokenLiteral() string {
    return ce.Token.Literal
}

func (ce *CallExpression) String() string {
    var out bytes.Buffer

    args := []string{}
    for _, a := range ce.Arguments {
        args = append(args, a.String())
    }

    out.WriteString(ce.Function.String())
    out.WriteString("(")
    out.WriteString(strings.Join(args, ", "))
    out.WriteString(")")

    return out.String()
}

func (ce *CallExpression) ExpressionNode() {}

type PrefixOpExpression struct {
    Token    token.Token // the prefix token, e.g. !
    Operator string
    Right    Expression
}

func (pe *PrefixOpExpression) TokenLiteral() string {
    return pe.Token.Literal
}

func (pe *PrefixOpExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(pe.Operator)
    out.WriteString(pe.Right.String())
    out.WriteString(")")

    return out.String()
}

func (i *PrefixOpExpression) ExpressionNode() {}

type InfixExpression struct {
    Token    token.Token // the infix token, e.g. +
    Left     Expression
    Operator string
    Right    Expression
}

func (ie *InfixExpression) TokenLiteral() string {
    return ie.Token.Literal
}

func (ie *InfixExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(ie.Left.String())
    out.WriteString(" " + ie.Operator + " ")
    out.WriteString(ie.Right.String())
    out.WriteString(")")

    return out.String()
}

func (ie *InfixExpression) ExpressionNode() {}
