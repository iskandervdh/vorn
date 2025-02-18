package evaluator

import (
	"fmt"

	"github.com/iskandervdh/vorn/ast"
	"github.com/iskandervdh/vorn/object"
)

type Evaluator struct {
	builtins map[string]*object.Builtin
}

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func New() *Evaluator {
	e := &Evaluator{}

	e.builtins = map[string]*object.Builtin{
		"len":    {Function: e.builtinLen},
		"first":  {Function: e.builtinFirst},
		"last":   {Function: e.builtinLast},
		"rest":   {Function: e.builtinRest},
		"push":   {Function: e.builtinPush},
		"pop":    {Function: e.builtinPop},
		"map":    {Function: e.builtinMap},
		"reduce": {Function: e.builtinReduce},
		"print":  {Function: e.builtinPrint},
		"pow":    {Function: e.builtinPow},
		"sqrt":   {Function: e.builtinSqrt},
	}

	return e
}

func (e *Evaluator) evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = e.Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func (e *Evaluator) evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = e.Eval(statement, env)

		if result != nil {
			rt := result.Type()

			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return NULL
}

func (e *Evaluator) nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}

func (e *Evaluator) evalExclamationOperatorExpression(right object.Object) object.Object {
	switch right {
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

func (e *Evaluator) evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	switch right.Type() {
	case object.INTEGER_OBJ:
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}

	case object.FLOAT_OBJ:
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}

	default:
		return newError("unknown operator: -%s", right.Type())
	}
}

func (e *Evaluator) evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return e.evalExclamationOperatorExpression(right)
	case "-":
		return e.evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func (e *Evaluator) evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: float64(leftVal) / float64(rightVal)}
	case "<":
		return e.nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return e.nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return e.nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return e.nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return e.nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return e.nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func (e *Evaluator) evalFloatInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "<":
		return e.nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return e.nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return e.nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return e.nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return e.nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return e.nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func (e *Evaluator) evalNumberInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return e.evalIntegerInfixExpression(operator, left, right)

	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return e.evalFloatInfixExpression(operator, left, right)

	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		left := float64(left.(*object.Integer).Value)
		return e.evalFloatInfixExpression(operator, &object.Float{Value: left}, right)

	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		right := float64(right.(*object.Integer).Value)
		return e.evalFloatInfixExpression(operator, left, &object.Float{Value: right})

	default:
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}
}

func (e *Evaluator) evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	return &object.String{Value: leftVal + rightVal}
}

func (e *Evaluator) evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case object.IsNumber(left) && object.IsNumber(right):
		return e.evalNumberInfixExpression(operator, left, right)

	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return e.evalStringInfixExpression(operator, left, right)

	case operator == "==":
		return e.nativeBoolToBooleanObject(left == right)

	case operator == "!=":
		return e.nativeBoolToBooleanObject(left != right)

	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func (e *Evaluator) evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := e.Eval(ie.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return e.evalBlockStatement(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return e.evalBlockStatement(ie.Alternative, env)
	}

	return NULL

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

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}

func (e *Evaluator) evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, _, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := e.builtins[node.Value]; ok {
		return builtin
	}

	return newError("[%d:%d] identifier not found: %s", node.Token.Line, node.Token.Column, node.Value)
}

func (e *Evaluator) evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, expression := range expressions {
		evaluated := e.Eval(expression, env)

		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func (e *Evaluator) extendFunctionEnv(function *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(function.Env)

	for paramIdx, param := range function.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func (e *Evaluator) unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func (e *Evaluator) applyFunction(function object.Object, args []object.Object) object.Object {
	switch function := function.(type) {
	case *object.Function:
		extendedEnv := e.extendFunctionEnv(function, args)
		evaluated := e.Eval(function.Body, extendedEnv)

		return e.unwrapReturnValue(evaluated)
	case *object.Builtin:
		return function.Function(args...)
	}

	return newError("not a function: %s", function.Type())
}

func (e *Evaluator) evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	// Negative index support from end of array
	if idx < 0 {
		idx = max + idx + 1
	}

	if idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

func (e *Evaluator) evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)
	key, ok := index.(object.Hashable)

	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]

	if !ok {
		return NULL
	}

	return pair.Value
}

func (e *Evaluator) evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return e.evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return e.evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func (e *Evaluator) evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := e.Eval(keyNode, env)

		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)

		if !ok {
			return newError("[%d:%d] unusable as hash key: %s", node.Token.Line, node.Token.Column, key.Type())
		}

		value := e.Eval(valueNode, env)

		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func (e *Evaluator) Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return e.evalProgram(node, env)

	case *ast.ExpressionStatement:
		return e.Eval(node.Expression, env)

	case *ast.BlockStatement:
		return e.evalBlockStatement(node, env)

	case *ast.ReturnStatement:
		val := e.Eval(node.ReturnValue, env)

		if isError(val) {
			return val
		}

		return &object.ReturnValue{Value: val}

	case *ast.VariableStatement:
		_, defined := env.GetFromCurrent(node.Name.Value)

		if defined {
			return newError("[%d:%d] variable already defined: %s", node.Token.Line, node.Token.Column, node.Name.Value)
		}

		value := e.Eval(node.Value, env)

		if isError(value) {
			return value
		}

		env.Set(node.Name.Value, value)

	case *ast.FunctionStatement:
		params := node.Parameters
		body := node.Body

		function := &object.Function{Parameters: params, Env: env, Body: body}

		env.Set(node.Name.Value, function)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.BooleanLiteral:
		return e.nativeBoolToBooleanObject(node.Value)

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.ArrayLiteral:
		elements := e.evalExpressions(node.Elements, env)

		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}

		return &object.Array{Elements: elements}

	case *ast.HashLiteral:
		return e.evalHashLiteral(node, env)

	case *ast.PrefixExpression:
		right := e.Eval(node.Right, env)

		if isError(right) {
			return right
		}

		return e.evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := e.Eval(node.Left, env)

		if isError(left) {
			return left
		}

		right := e.Eval(node.Right, env)

		if isError(right) {
			return right
		}

		return e.evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return e.evalIfExpression(node, env)

	case *ast.Identifier:
		return e.evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body

		return &object.Function{Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		function := e.Eval(node.Function, env)

		if isError(function) {
			return function
		}

		args := e.evalExpressions(node.Arguments, env)

		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return e.applyFunction(function, args)

	case *ast.IndexExpression:
		left := e.Eval(node.Left, env)

		if isError(left) {
			return left
		}

		index := e.Eval(node.Index, env)

		if isError(index) {
			return index
		}

		return e.evalIndexExpression(left, index)

	case *ast.ReassignmentExpression:
		_, environment, defined := env.Get(node.Name.Value)

		if !defined {
			return newError("[%d:%d] variable %s has not been initialized.", node.Token.Line, node.Token.Column, node.Name.Value)
		}

		value := e.Eval(node.Value, env)

		if isError(value) {
			return value
		}

		environment.Set(node.Name.Value, value)
	}

	return nil
}
