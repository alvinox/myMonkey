package token

type TokenType string

const (
    ILLEGAL     = "ILLEGAL"
    EOF         = "EOF"

    // identifiers + literals
    IDENT       = "IDENT"
    INT         = "INT"
    STRING      = "STRING"

    // operators
    ASSIGN      = "ASSIGN"
    PLUS        = "+"
    MINUS       = "-"
    BANG        = "!"
    ASTERISK    = "*"
    SLASH       = "/"

    LT          = "<"
    GT          = ">"

    EQ          = "=="
    NEQ         = "!="

    // delimiters
    COMMA       = ","
    SEMICOLON   = ";"
    COLON       = ":"

    LPAREN      = "("
    RPAREN      = ")"
    LBRACKET    = "["
    RBRACKET    = "]"
    LBRACE      = "{"
    RBRACE      = "}"

    // keywords
    FUNCTION    = "FUNCTION"
    LET         = "LET"
    TRUE        = "TRUE"
    FALSE       = "FALSE"
    IF          = "IF"
    ELSE        = "ELSE"
    RETURN      = "RETURN"
)

type Token struct {
    Type    TokenType
    Literal string
}

var keywords = map[string] TokenType {
    "fn":      FUNCTION,
    "let":     LET,
    "true":    TRUE,
    "false":   FALSE,
    "if":      IF,
    "else":    ELSE,
    "return":  RETURN,
}

func LookupIdent(ident string) TokenType {
    if tok, ok := keywords[ident]; ok {
        return tok
    }

    return IDENT
}