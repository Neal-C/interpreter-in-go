package evaluator

import (
	"github.com/Neal-C/interpreter-in-go/ast"
	"github.com/Neal-C/interpreter-in-go/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeNodeToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		leftHandSign := Eval(node.Left)
		rightHandSign := Eval(node.Right)
		return evalInfixExpression(node.Operator, leftHandSign, rightHandSign)
	case *ast.BlockStatement:
		return evalStatements(node.Statements)
	case *ast.IfExpression:
		return evalIfExpression(node)

	}

	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
	}

	return result
}

func nativeNodeToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}

func evalPrefixExpression(operator string, rightHandSign object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(rightHandSign)
	case "-":
		return evalMinusPrefixOperatorExpression(rightHandSign)
	default:
		return NULL
	}
}

func evalBangOperatorExpression(rightHandSign object.Object) object.Object {
	switch rightHandSign {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(rightHandSign object.Object) object.Object {
	if rightHandSign.Type() != object.INTEGER_OBJ {
		return NULL
	}

	value := rightHandSign.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, leftHandSign object.Object, rightHandSign object.Object) object.Object {
	switch {
	case leftHandSign.Type() == object.INTEGER_OBJ && rightHandSign.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, leftHandSign, rightHandSign)
	case operator == "==":
		return nativeNodeToBooleanObject(leftHandSign == rightHandSign)
	case operator == "!=":
		return nativeNodeToBooleanObject(leftHandSign != rightHandSign)
	default:
		return NULL
	}

}

func evalIntegerInfixExpression(operator string, leftHandSign object.Object, rightHandSign object.Object) object.Object {
	leftValue := leftHandSign.(*object.Integer).Value
	rightValue := rightHandSign.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case "<":
		return nativeNodeToBooleanObject(leftValue < rightValue)
	case ">":
		return nativeNodeToBooleanObject(leftValue > rightValue)
	case "==":
		return nativeNodeToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeNodeToBooleanObject(leftValue != rightValue)
	default:
		return NULL
	}
}

func evalIfExpression(ifExpr *ast.IfExpression) object.Object {
	condition := Eval(ifExpr.Condition)

	if isTruthy(condition) {
		return Eval(ifExpr.Consequence)
	} else if ifExpr.Alternative != nil {
		return Eval(ifExpr.Alternative)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}
