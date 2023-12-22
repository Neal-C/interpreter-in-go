package evaluator

import (
	"fmt"
	"github.com/Neal-C/interpreter-in-go/ast"
	"github.com/Neal-C/interpreter-in-go/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeNodeToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		rightHandSign := Eval(node.Right, env)
		if isError(rightHandSign) {
			return rightHandSign
		}
		return evalPrefixExpression(node.Operator, rightHandSign)
	case *ast.InfixExpression:
		leftHandSign := Eval(node.Left, env)
		if isError(leftHandSign) {
			return leftHandSign
		}
		rightHandSign := Eval(node.Right, env)
		if isError(rightHandSign) {
			return rightHandSign
		}
		return evalInfixExpression(node.Operator, leftHandSign, rightHandSign)
	case *ast.BlockStatement:
		return evalBlockStatements(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		value := Eval(node.ReturnValue, env)
		if isError(value) {
			return value
		}
		return &object.ReturnValue{Value: value}
	case *ast.LetStatement:
		value := Eval(node.Value, env)
		if isError(value) {
			return value
		}
		env.Set(node.Name.Value, value)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}
	case *ast.CallExpression:
		fnCall := Eval(node.Function, env)

		if isError(fnCall) {
			return fnCall
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(fnCall, args)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	}

	return nil
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
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
		return newError("unknown operator: %s%s", operator, rightHandSign.Type())
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
		return newError("unknown operator: -%s", rightHandSign.Type())
	}
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
	case leftHandSign.Type() != rightHandSign.Type():
		return newError("type mismatch: %s %s %s", leftHandSign.Type(), operator, rightHandSign.Type())
	case leftHandSign.Type() == object.STRING_OBJ && rightHandSign.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, leftHandSign, rightHandSign)
	default:
		return newError("unknown operator: %s %s %s", leftHandSign.Type(), operator, rightHandSign.Type())
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
		return newError("unknown operator: %s %s %s", leftHandSign.Type(), operator, rightHandSign.Type())
	}
}

func evalIfExpression(ifExpr *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ifExpr.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ifExpr.Consequence, env)
	} else if ifExpr.Alternative != nil {
		return Eval(ifExpr.Alternative, env)
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

func evalBlockStatements(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			resultType := result.Type()
			if resultType == object.RETURN_VALUE_OBJ || resultType == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func newError(format string, others ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, others...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {

	if value, ok := env.Get(node.Value); ok {
		return value
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("identifier not found: " + node.Value)
}

func evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, expr := range expressions {
		evaluated := Eval(expr, env)

		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fnCall object.Object, args []object.Object) object.Object {
	fn, ok := fnCall.(*object.Function)

	if !ok {
		return newError("fn is not a function: %s ", fn.Type())
	}

	extendedEnv := extendFunctionEnv(fn, args)
	evaluated := Eval(fn.Body, extendedEnv)

	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIndex, param := range fn.Parameters {
		env.Set(param.Value, args[paramIndex])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalStringInfixExpression(operator string, leftHandSign object.Object, rightHandSign object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", leftHandSign.Type(), operator, rightHandSign.Type())
	}

	leftValue := leftHandSign.(*object.String).Value
	rightValue := rightHandSign.(*object.String).Value
	return &object.String{Value: leftValue + rightValue}
}
