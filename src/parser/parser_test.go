package parser

import (
	"../ast"
	"fmt"
	"../lexer"
	"testing"
)

func TestLetStatments(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;`

	program := CreateProgram(t, input)

	if program == nil {
		t.Fatalf("Expected ParseProgram to return an object, and not nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Expected program.Statements to contain 3 statements, but got %d", len(program.Statements))
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
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatments(t *testing.T) {
	input := `
		return 5;
		return 10;
		return 993322;`

	program := CreateProgram(t, input)

	if program == nil {
		t.Fatalf("Expected ParseProgram to return an object, and not nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Expected program.Statements to contain 3 statements, but got %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("Expected statment to be a return statment but was %T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("Expected token litaral to be 'return' but was '%q'", returnStmt.TokenLiteral())
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("Expected token literal to be 'let' but was %q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("Expected let statement but got %T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("Expected Name value identifer to be '%s' but got '%s'", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("Expected Name to be '%s' but got '%s'", name, letStmt.Name)
		return false
	}

	return true
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	program := CreateProgram(t, input)

	AssertProgramStatementLength(t, program, 1)

	stmt := GetExpressionStatement(t, program, 0)
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected Identifier but got %T", stmt.Expression)
	}

	AssertStringsAreEqual(t, ident.Value, "foobar", "Identifer Value")
	AssertStringsAreEqual(t, ident.TokenLiteral(), "foobar", "Identifier TokenLiteral")
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	program := CreateProgram(t, input)
	AssertProgramStatementLength(t, program, 1)

	stmt := GetExpressionStatement(t, program, 0)

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expected integer literal but got %T", stmt.Expression)
	}

	AssertIntegersAreEqual(t, literal.Value, 5, "Literal Value")
	AssertStringsAreEqual(t, literal.TokenLiteral(), "5", "Literal TokenLiteral")
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range prefixTests {
		program := CreateProgram(t, tt.input)
		AssertProgramStatementLength(t, program, 1)

		stmt := GetExpressionStatement(t, program, 0)

		prefixExp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("Expected prefix expression but got %T", stmt.Expression)
		}

		if prefixExp.Operator != tt.operator {
			t.Fatalf("Expected operator to be %s but got %s", tt.operator, prefixExp.Operator)
		}

		if !testIntegerLiteral(t, prefixExp.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"1 + 2", 1, "+", 2},
		{"1 - 2", 1, "-", 2},
		{"1 * 2", 1, "*", 2},
		{"4 / 2", 4, "/", 2},
		{"1 > 2", 1, ">", 2},
		{"1 < 2", 1, "<", 2},
		{"1 == 2", 1, "==", 2},
		{"1 != 2", 1, "!=", 2},
	}

	for _, tt := range infixTests {
		program := CreateProgram(t, tt.input)
		AssertProgramStatementLength(t, program, 1)

		stmt := GetExpressionStatement(t, program, 0)

		infixExp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("Expected infix expression but got %T", stmt.Expression)
		}

		if !testIntegerLiteral(t, infixExp.Left, tt.leftValue) {
			return
		}

		if infixExp.Operator != tt.operator {
			t.Fatalf("Expected infix expression operator to be %s but got %s", tt.operator, infixExp.Operator)
		}

		if !testIntegerLiteral(t, infixExp.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, tt := range tests {
		program := CreateProgram(t, tt.input)

		actual := program.String()
		AssertStringsAreEqual(t, actual, tt.expected, "program")
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("Expected expression to be integer literal but got %T", il)
		return false
	}

	if integer.Value != value {
		t.Errorf("Expected integer value to be %d but got %d", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("Expected integer literal token literal to be %d but got %s", value, integer.TokenLiteral())
		return false
	}

	return true
}

func CreateProgram(t *testing.T, input string) *ast.Program {
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	return program
}

func AssertProgramStatementLength(t *testing.T, program *ast.Program, expectedStatmemtCount int) {
	if len(program.Statements) != expectedStatmemtCount {
		t.Fatalf("Statments not parsed correctly. Expected one statement, but got %d", len(program.Statements))
	}
}

func GetExpressionStatement(t *testing.T, program *ast.Program, expectedStatementIndex int) *ast.ExpressionStatement {
	stmt, ok := program.Statements[expectedStatementIndex].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatment but got %T", program.Statements[0])
	}

	return stmt
}

func AssertStringsAreEqual(t *testing.T, actual string, expected string, actualDescription string) {
	if actual != expected {
		t.Errorf("Expected %s to be %s but found %s", actualDescription, expected, actual)
	}
}

func AssertIntegersAreEqual(t *testing.T, actual int64, expected int64, actualDescription string) {
	if actual != expected {
		t.Errorf("Expected %s to be %d but found %d", actualDescription, expected, actual)
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("Parser error: %q", msg)
	}

	t.FailNow()
}
