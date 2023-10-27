package parser

import (
	"github.com/Neal-C/interpreter-in-go/ast"
	"github.com/Neal-C/interpreter-in-go/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `let x = 5;
let y = 10;
let foobar = 838383`

	var theLexer *lexer.Lexer = lexer.New(input)
	var parser *Parser = New(theLexer)

	var program *ast.Program = parser.ParseProgram()

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements do not contain 3 statements, got : %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifer string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for index, tabletest := range tests {
		stmt := program.Statements[index]
		if !testLetStatement(t, stmt, tabletest.expectedIdentifer) {
			return
		}
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("stmt.TœokenLiteral not 'let' , got = %q", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("stmt is not a let statement, got %T", letStmt)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s', got = %s ", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("stmt.Name not '%s' , got=%s", name, letStmt.Name)
		return false
	}

	return true
}
