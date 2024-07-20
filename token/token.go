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
