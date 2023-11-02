package parser

import (
	"github.com/Neal-C/interpreter-in-go/ast"
	"github.com/Neal-C/interpreter-in-go/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `let x = 5;
let y = 10;
let foobar = 838383;`

	var theLexer *lexer.Lexer = lexer.New(input)
	var parser *Parser = New(theLexer)

	var program *ast.Program = parser.ParseProgram()
	checkParserErrors(t, parser)

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

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser had %d erros", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error %q", msg)
	}
	t.FailNow()
}

func testReturnStatement(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`

	var theLexer *lexer.Lexer = lexer.New(input)
	var parser *Parser = New(theLexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not a *ast.ReturnStatement, got %T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("retrunStmt literal is not 'return', got %q", returnStmt.TokenLiteral())
		}
	}
}
