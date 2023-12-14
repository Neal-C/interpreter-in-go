package evaluator

import (
	"github.com/Neal-C/interpreter-in-go/lexer"
	"github.com/Neal-C/interpreter-in-go/object"
	"github.com/Neal-C/interpreter-in-go/parser"
	"testing"
)

func TestEvalIntegerExpress(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}

}

func testEval(input string) object.Object {
	myLexer := lexer.New(input)
	myParser := parser.New(myLexer)
	program := myParser.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, evaluated object.Object, expected int64) bool {
	result, ok := evaluated.(*object.Integer)

	if !ok {
		t.Errorf("object is not an Integer, got %T (%+v)", evaluated, evaluated)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value, got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tableTests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tableTests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}

}

func testBooleanObject(t *testing.T, evaluated object.Object, expected bool) bool {
	result, ok := evaluated.(*object.Boolean)

	if !ok {
		t.Errorf("evaluated is not an *object.Boolean, got %T, (%+v) ", evaluated, evaluated)
		return false
	}

	if result.Value != expected {
		t.Errorf("evaluated has wrong value, got = %t , want = %t", result.Value, expected)
		return false
	}

	return true

}

func TestBangOperator(t *testing.T) {
	tableTests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tableTests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}
