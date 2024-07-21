package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT = "IDENT"
	INT   = "INT"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	ASTERISK = "*"
	MINUS    = "-"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"
	BANG     = "!"
	EQ       = "=="
	NOT_EQ   = "!="
	GTE      = ">="
	LTE      = "<="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	VAR      = "VAR"
	RETURN   = "RETURN"
	BREAK    = "BREAK"
	IF       = "IF"
	ELSE     = "ELSE"
	BLANKOUT = "BLANKOUT"
	BLANKIN  = "BLANKIN"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
)

var keywords = map[string]TokenType{
	"func":     FUNCTION,
	"var":      VAR,
	"return":   RETURN,
	"break":    BREAK,
	"if":       IF,
	"else":     ELSE,
	"blankout": BLANKOUT,
	"blankin":  BLANKIN,
	"true":     TRUE,
	"false":    FALSE,
}

func LookupIdent(ident string) TokenType {
	if tok, exists := keywords[ident]; exists {
		return tok
	}
	return IDENT
}
