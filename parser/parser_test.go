package parser

import (
	"blank/ast"
	"blank/lexer"
	"blank/token"
	"fmt"
	"testing"
)

func TestVarStatements(t *testing.T) {
	input := `
		var x = 5;
		var y = 10;
		var foobar = 838383;
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testVarStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testVarStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "var" {
		t.Errorf("s.TokenLiteral not 'var'. got=%q", s.TokenLiteral())
		return false
	}

	varStmt, exists := s.(*ast.VarStatement)
	if !exists {
		t.Errorf("s not *ast.VarStatement. got=%T", s)
		return false
	}

	if varStmt.Name.Value != name {
		t.Errorf("letvarStmt.Name.Value not '%s'. got=%s", name, varStmt.Name.Value)
		return false
	}

	if varStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, varStmt.Name)
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 993322;
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}
	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			ident.TokenLiteral())
	}
}

func TestIntegerExpression(t *testing.T) {
	input := "190;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IntegerLiteral. got=%T",
			program.Statements[0])
	}
	integer, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", integer.Value)
	}
	if integer.Value != 190 {
		t.Errorf("integer.Value not %s. got=%d", "190", integer.Value)
	}
	if integer.Token.Literal != "190" {
		t.Errorf("integer.TokenLiteral not %s. got=%s", "190",
			integer.TokenLiteral())
	}
}

func TestPrefixExpression(t *testing.T) {
	input := `
		!test;
		-30;
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 2 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	for i := 0; i < len(program.Statements); i++ {
		stmt, ok := program.Statements[i].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[i])
		}
		prefix, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("exp not *ast.PrefixExpression. got=%T", prefix.TokenLiteral())
		}
	}
	prefix := program.Statements[0].(*ast.ExpressionStatement).Expression.(*ast.PrefixExpression)
	if prefix.TokenLiteral() != token.BANG {
		t.Errorf("prefix.TokenLiteral not '!'. got=%s",
			prefix.TokenLiteral())
	}
	CheckIndentifierExpression(t, prefix.Right, "test")

	prefix = program.Statements[1].(*ast.ExpressionStatement).Expression.(*ast.PrefixExpression)
	if prefix.TokenLiteral() != token.MINUS {
		t.Errorf("prefix.TokenLiteral not '-'. got=%s",
			prefix.TokenLiteral())
	}
	CheckIntegerExpression(t, prefix.Right, 30)
}

func CheckIntegerExpression(t *testing.T, il ast.Expression, value int64) {
	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", integer)
	}
	if integer.Value != value {
		t.Errorf("integer.Value not %d. got=%d", value, integer.Value)
	}
	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral not %d. got=%s", value,
			integer.TokenLiteral())
	}
}

func CheckIndentifierExpression(t *testing.T, il ast.Expression, value string) {
	expression, ok := il.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", expression)
	}
	if expression.Value != value {
		t.Errorf("expression.Value not %s. got=%s", value, expression.Value)
	}
	if expression.TokenLiteral() != value {
		t.Errorf("expression.TokenLiteral not %s. got=%s", value,
			expression.TokenLiteral())
	}
}

func TestInfixExpressoin(t *testing.T) {
	input := `
		5 - 5;
		23 * 2;
		13 + 13;
		10 / 2;
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	type OutputStruct struct {
		left     int64
		operator string
		right    int64
	}
	output := []OutputStruct{
		{left: 5, operator: token.MINUS, right: 5},
		{left: 23, operator: token.ASTERISK, right: 2},
		{left: 13, operator: token.PLUS, right: 13},
		{left: 10, operator: token.SLASH, right: 2},
	}
	checkParserErrors(t, p)
	if len(program.Statements) != 4 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	for i := 0; i < len(program.Statements); i++ {
		stmt, ok := program.Statements[i].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[%d] is not ast.IntegerLiteral. got=%T", i,
				program.Statements[i])
		}
		infixExpression, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("infixExpression not *ast.infixExpression. got=%T", infixExpression.String())
		}
		if infixExpression.Left.(*ast.IntegerLiteral).Value != output[i].left {
			t.Errorf("infixExpression.Value not %d. got=%d", output[i].left, infixExpression.Left.(*ast.IntegerLiteral).Value)
		}
		if infixExpression.TokenLiteral() != output[i].operator {
			t.Errorf("infixExpression.Operator.Literal not %s. got=%s", output[i].operator,
				infixExpression.TokenLiteral())
		}
		if infixExpression.Right.(*ast.IntegerLiteral).Value != output[i].right {
			t.Errorf("infixExpression.Value not %d. got=%d", output[i].right, infixExpression.Right.(*ast.IntegerLiteral).Value)
		}
	}
}