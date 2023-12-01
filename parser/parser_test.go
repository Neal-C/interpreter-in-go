package parser

import (
	"fmt"
	"github.com/Neal-C/interpreter-in-go/ast"
	"github.com/Neal-C/interpreter-in-go/lexer"
	"log"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `let x = 5;
let y = 10;
let foobar = 838383;`

	var theLexer *lexer.Lexer = lexer.New(input)
	var parser *Parser = New(theLexer)

	var program *ast.Program = parser.ParseProgram()
	log.Println(program.String())
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

func TestReturnStatement(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`

	var theLexer *lexer.Lexer = lexer.New(input)
	var parser *Parser = New(theLexer)

	program := parser.ParseProgram()
	log.Println(program.String())
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

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`

	currentLexer := lexer.New(input)
	currentParser := New(currentLexer)
	currentProgram := currentParser.ParseProgram()
	log.Println(currentProgram.String())
	checkParserErrors(t, currentParser)

	if len(currentProgram.Statements) != 1 {
		t.Fatalf("program has not enough statements. got %d", len(currentProgram.Statements))
	}

	stmt, ok := currentProgram.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expression is not a *ast.ExpressStatement, got %T", currentProgram.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expression is not an *ast.Identifier, got %T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s, got %s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral() not %s, got : %s ", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteral(t *testing.T) {
	input := `5;`

	myLexer := lexer.New(input)
	parser := New(myLexer)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement, got %T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf("stmt.Expression is not a *ast.IntegerLiteral, got %T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value not %d, got %d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral() is not '%d', got %s ", 5, literal.TokenLiteral())
	}
}

func testIntegerLiteral(t *testing.T, intergerLiteral ast.Expression, value int64) bool {
	integer, ok := intergerLiteral.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("intergerLiteral is not a *ast.IntegerLiteral, got %T", intergerLiteral)
		return false
	}

	if integer.Value != value {
		t.Errorf("integer.Value not %d, got %d", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral() not %d, got %s", value, integer.TokenLiteral())
		return false
	}

	return true
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    any
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		myLexer := lexer.New(tt.input)
		parser := New(myLexer)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements, got %d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not a *ast.ExpressionStatement, got %T", program.Statements[0])
		}

		expression, ok := stmt.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("stmt is not a *ast.PrefixExpression, got %T", stmt.Expression)
		}

		if expression.Operator != tt.operator {
			t.Fatalf("expression operator is not %s, got %s", "!", expression.Operator)
		}

		if !testLiteralExpression(t, expression.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		myLexer := lexer.New(tt.input)
		parser := New(myLexer)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d, got %d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement, got %T", stmt)
		}

		expression, ok := stmt.Expression.(*ast.InfixExpression)

		if !ok {
			t.Fatalf("stmt.Expression is not a *ast.InfixExpression, got %T", stmt.Expression)
		}

		if !testLiteralExpression(t, expression.Left, tt.leftValue) {
			return
		}

		if expression.Operator != tt.operator {
			t.Fatalf("Operator is not %s , got=%s", tt.operator, expression.Operator)
		}

		if !testLiteralExpression(t, expression.Right, tt.rightValue) {
			return
		}

	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tableTest := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 > 4 != 3 > 4", "((5 > 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
	}

	for _, tt := range tableTest {
		myLexer := lexer.New(tt.input)
		parser := New(myLexer)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected = %q , got = %q", tt.expected, actual)
		}
	}
}

func TestIfExpressions(t *testing.T) {
	input := `if ( x < y ) { x }`

	myLexer := lexer.New(input)
	parser := New(myLexer)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program's body does not contain %d statement, got %d \n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement, got %T", program.Statements[0])
	}

	expression, ok := stmt.Expression.(*ast.IfExpression)

	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IfExpression, got %T", stmt.Expression)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("expression.Consequence.Statements is not 1 statements, got %d \n", len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("expression.Consequence.Statements[0] is not *ast.ExpressionStatement, got %T", expression.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if expression.Alternative != nil {
		t.Errorf("expression.Alternative was not nil, got %+v", expression.Alternative)
	}
}

func TestIfElseExpressions(t *testing.T) {
	input := `if ( x < y ) { x } else { y }`

	myLexer := lexer.New(input)
	parser := New(myLexer)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program's body does not contain %d statement, got %d \n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement, got %T", program.Statements[0])
	}

	expression, ok := stmt.Expression.(*ast.IfExpression)

	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IfExpression, got %T", stmt.Expression)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("expression.Consequence.Statements is not 1 statements, got %d \n", len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("expression.Consequence.Statements[0] is not *ast.ExpressionStatement, got %T", expression.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if expression.Alternative == nil {
		t.Errorf("expression.Alternative was  nil, got %+v", expression.Alternative)
	}

	if len(expression.Alternative.Statements) != 1 {
		t.Errorf("expression.Alternative.Statements is not 1 statements, got %d \n", len(expression.Alternative.Statements))
	}

	alternative, ok := expression.Alternative.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("expression.Alternative.Statements[0] is not *ast.ExpressionStatement, got %T", expression.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func testIdentifier(t *testing.T, expression ast.Expression, value string) bool {
	identifier, ok := expression.(*ast.Identifier)

	if !ok {
		t.Errorf("expression is not *ast.Identifier, got %T", expression)
		return false
	}

	if identifier.Value != value {
		t.Errorf("identifier.Value is not %s, got %s", value, identifier.Value)
		return false
	}

	if identifier.TokenLiteral() != value {
		t.Errorf("identifier.TokenLiteral() is not %s, got %s", value, identifier.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, expression ast.Expression, value bool) bool {
	boolean, ok := expression.(*ast.Boolean)

	if !ok {
		t.Errorf("expression is not a *ast.Boolean, got %T", expression)
		return false
	}

	if boolean.Value != value {
		t.Errorf("boolean.Value is not %t, got %t", value, boolean.Value)
		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("boolean.TokenLiteral() not %t, got %s", value, boolean.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, expression ast.Expression, expected any) bool {
	switch value := expected.(type) {
	case int:
		return testIntegerLiteral(t, expression, int64(value))
	case int64:
		return testIntegerLiteral(t, expression, value)
	case string:
		return testIdentifier(t, expression, value)
	case bool:
		return testBooleanLiteral(t, expression, value)
	}

	t.Errorf("type of the expression is not handled, got : %T", expression)
	return false
}

func testInfixExpression(t *testing.T, expression ast.Expression, leftHand any, operator string, rightHand any) bool {

	infixExpression, ok := expression.(*ast.InfixExpression)

	if !ok {
		t.Errorf("expression is not an *ast.InfixExpression, got = %T(%s)", expression, expression)
		return false
	}

	if !testLiteralExpression(t, infixExpression.Left, leftHand) {
		return false
	}

	if infixExpression.Operator != operator {
		t.Errorf("infixExpression.Operator is not '%s' , got = %q", operator, infixExpression.Operator)
		return false
	}

	if !testLiteralExpression(t, infixExpression.Right, rightHand) {
		return false
	}

	return true
}

// TODO: clean test suite by using the test utils
