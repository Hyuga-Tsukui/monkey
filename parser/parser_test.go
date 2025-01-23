package parser

import (
	"reflect"
	"testing"

	"github.com/Hyuga-Tsukui/monkey/ast"
	"github.com/Hyuga-Tsukui/monkey/lexer"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		statementCount     int
		expectedIdentifier []string
		hasParseError      bool
		expectedErrors     []string
	}{
		{
			name: "let statement",
			input: `
			let x = 5;
			let y = 10;
			let foobar = 838383;
			`,
			statementCount:     3,
			expectedIdentifier: []string{"x", "y", "foobar"},
		},
		{
			name: "let statement with parse error",
			input: `
			let x = 5;
			let = 10;
			let foobar = 838383;
			`,
			hasParseError: true,
			expectedErrors: []string{
				"expected next token to be IDENT, got = instead",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			l := lexer.New(tt.input)
			p := New(l)

			program := p.ParseProgram()
			if tt.hasParseError {
				if !reflect.DeepEqual(p.errors, tt.expectedErrors) {
					t.Errorf("p.errors = %v, want %v", p.errors, tt.expectedErrors)
					t.Fatalf("ParseProgram() has parse error, but not expected error")
				}
				return
			}
			if program == nil {
				t.Fatalf("ParseProgram() returned nil")
			}
			if len(program.Statements) != tt.statementCount {
				t.Fatalf("program.Statements dose not contain %d statements. got=%d dump=%#v", tt.statementCount, len(program.Statements), program)
			}

			for i := 0; i < tt.statementCount; i++ {
				stmt := program.Statements[i]
				if !testLetStatement(t, stmt, tt.expectedIdentifier[i]) {
					return
				}
			}
		})
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStmt.TokenLiteral())
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, p *Parser, expectedErrors []string) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	if !reflect.DeepEqual(errors, expectedErrors) {
		t.Fail()
	}
}
