package parser

import (
	"ast"
	"lexer"
	"testing"
)

func TestLetStatments(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

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
