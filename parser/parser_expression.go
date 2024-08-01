package parser

import (
	"blank/ast"
	"blank/token"
	"fmt"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

// operatorsPrecendence associates the operator with its precedence order
var operatorsPrecendence = map[token.TokenType]int{
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
}

// parseExpressionStatement parses a whole expression. Example: 2 - 2 * 5 + 4;
// Returns a reference to a expression statement node
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// parseExpression parses each token of the expression assigning it to an ast node evaluating their operators precedence.
// Returns a precedence of operations tree. Example: 2 + 2 * 4 - 5
//
//				expression
//				/        \
//			infix 		prefix
//			/	\
//		infix	prefix
//		/	\
//	prefix	 prefix
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		msg := fmt.Sprintf("no prefix parse function for %s found", p.curToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

// parsePrefixExpression creates prefix expression node and returns it reference.
func (p *Parser) parsePrefixExpression() ast.Expression {
	stmtPrefix := &ast.PrefixExpression{Token: p.curToken}

	p.nextToken()

	stmtPrefix.Right = p.parseExpression(PREFIX)

	return stmtPrefix
}

// parseIdentifier creates identifier expression node and returns it reference.
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// parseIntegerLiteral creates integer literal expression node and returns it reference.
func (p *Parser) parseIntegerLiteral() ast.Expression {
	stmtInt := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("Error, %s is not a number!!", p.curToken.Literal)
		p.errors = append(p.errors, msg)
	}
	stmtInt.Value = value
	return stmtInt
}

// parseInfixExpression creates infix expression node and returns it reference.
func (p *Parser) parseInfixExpression(leftExp ast.Expression) ast.Expression {
	stmtInfix := &ast.InfixExpression{
		Left:     leftExp,
		Operator: p.curToken,
	}
	p.nextToken()
	stmtInfix.Right = p.parseExpression(operatorsPrecendence[stmtInfix.Operator.Type])
	return stmtInfix
}

// parseBoolean creates boolean expression node and returns it reference.
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.BooleanExpression{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

// HELPERS

// peekPrecedence returns, if exists, the operator precedence of the next token, otherwise,
// returns LOWEST
func (p *Parser) peekPrecedence() int {
	if p, ok := operatorsPrecendence[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// curPrecedence returns, if exists, the operator precedence of the current token, otherwise,
// returns LOWEST
func (p *Parser) curPrecedence() int {
	if p, ok := operatorsPrecendence[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}
