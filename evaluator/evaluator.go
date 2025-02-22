package evaluator

import (
	"fmt"

	"github.com/iskandervdh/vorn/ast"
	"github.com/iskandervdh/vorn/object"
	"github.com/iskandervdh/vorn/token"
)

type StringChainingFunction func(left *object.String, args ...object.Object) object.Object
type ArrayChainingFunction func(left *object.Array, args ...object.Object) object.Object
type ObjectChainingFunction func(left *object.Hash, args ...object.Object) object.Object

type Evaluator struct {
	builtins map[string]*object.Builtin

	stringChainingFunctions map[string]StringChainingFunction
	arrayChainingFunctions  map[string]ArrayChainingFunction
	objectChainingFunctions map[string]ObjectChainingFunction
}

// Reusable objects for TRUE, FALSE and NULL
var (
	NULL = object.NewNull(&ast.NullLiteral{
		Token: token.Token{
			Type:    "null",
			Literal: "null",
		},
	})
	TRUE = object.NewBoolean(&ast.BooleanLiteral{
		Token: token.Token{
			Type:    "true",
			Literal: "true",
			Line:    1,
			Column:  1,
		},
		Value: true,
	}, true)
	FALSE = object.NewBoolean(&ast.BooleanLiteral{
		Token: token.Token{
			Type:    "false",
			Literal: "false",
			Line:    1,
			Column:  1,
		},
		Value: false,
	}, false)
)

func New() *Evaluator {
	e := &Evaluator{}

	e.builtins = map[string]*object.Builtin{
		// Common
		"type":  {Function: e.builtinType, ArgumentsCount: 1},
		"range": {Function: e.builtinRange, ArgumentsCount: -1}, // Variable amount of arguments

		// Conversions
		"int":    {Function: e.builtinInt, ArgumentsCount: 1},
		"float":  {Function: e.builtinFloat, ArgumentsCount: 1},
		"string": {Function: e.builtinString, ArgumentsCount: 1},
		"bool":   {Function: e.builtinBool, ArgumentsCount: 1},

		// Strings & Arrays
		"len":   {Function: e.builtinLen, ArgumentsCount: 1},
		"first": {Function: e.builtinFirst, ArgumentsCount: 1},
		"last":  {Function: e.builtinLast, ArgumentsCount: 1},

		// Arrays
		"rest": {Function: e.builtinRest, ArgumentsCount: 1},

		// IO
		"print": {Function: e.builtinPrint, ArgumentsCount: -1}, // Variable amount of arguments

		// Math
		"abs":  {Function: e.builtinAbs, ArgumentsCount: 1},
		"pow":  {Function: e.builtinPow, ArgumentsCount: 2},
		"sqrt": {Function: e.builtinSqrt, ArgumentsCount: 1},
		"sin":  {Function: e.builtinSin, ArgumentsCount: 1},
		"cos":  {Function: e.builtinCos, ArgumentsCount: 1},
		"tan":  {Function: e.builtinTan, ArgumentsCount: 1},
		"sum":  {Function: e.builtinSum, ArgumentsCount: 1},
		"mean": {Function: e.builtinMean, ArgumentsCount: 1},
	}

	e.stringChainingFunctions = map[string]StringChainingFunction{
		"length":     e.stringLength,
		"upper":      e.stringUpper,
		"lower":      e.stringLower,
		"split":      e.stringSplit,
		"contains":   e.stringContains,
		"replace":    e.stringReplace,
		"trim":       e.stringTrim,
		"trimStart":  e.stringTrimStart,
		"trimEnd":    e.stringTrimEnd,
		"repeat":     e.stringRepeat,
		"reverse":    e.stringReverse,
		"slice":      e.stringSlice,
		"startsWith": e.stringStartsWith,
		"endsWith":   e.stringEndsWith,
	}

	e.arrayChainingFunctions = map[string]ArrayChainingFunction{
		"length":   e.arrayLength,
		"prepend":  e.arrayPrepend,
		"append":   e.arrayAppend,
		"shift":    e.arrayShift,
		"pop":      e.arrayPop,
		"concat":   e.arrayConcat,
		"map":      e.arrayMap,
		"filter":   e.arrayFilter,
		"reduce":   e.arrayReduce,
		"contains": e.arrayContains,
		"indexOf":  e.arrayIndexOf,
		"find":     e.arrayFind,
		"join":     e.arrayJoin,
		"reverse":  e.arrayReverse,
		"slice":    e.arraySlice,
		"sort":     e.arraySort,
		"any":      e.arrayAny,
		"every":    e.arrayAll,
	}

	e.objectChainingFunctions = map[string]ObjectChainingFunction{
		"keys":   e.objectKeys,
		"values": e.objectValues,
		"items":  e.objectItems,
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

func (e *Evaluator) evalBlockStatement(block *ast.BlockStatement, parentEnv *object.Environment) object.Object {
	var result object.Object
	env := object.NewEnclosedEnvironment(parentEnv)

	for _, statement := range block.Statements {
		result = e.Eval(statement, env)

		if result != nil {
			rt := result.Type()

			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ || rt == object.BREAK_OBJ || rt == object.CONTINUE_OBJ {
				return result
			}
		}
	}

	return NULL
}

func (e *Evaluator) evalWhileStatement(we *ast.WhileStatement, env *object.Environment) object.Object {
	for {
		condition := e.Eval(we.Condition, env)

		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result := e.evalBlockStatement(we.Consequence, env)

		if result != nil {
			resultType := result.Type()

			if resultType == object.CONTINUE_OBJ {
				continue
			} else if resultType == object.BREAK_OBJ {
				break
			} else if resultType == object.RETURN_VALUE_OBJ || resultType == object.ERROR_OBJ {
				return result
			}
		}
	}

	return NULL
}

func (e *Evaluator) evalForStatement(fs *ast.ForStatement, env *object.Environment) object.Object {
	if fs.Init != nil {
		e.Eval(fs.Init, env)
	}

	for {
		condition := e.Eval(fs.Condition, env)

		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result := e.evalBlockStatement(fs.Body, env)

		if result != nil {
			resultType := result.Type()

			if resultType == object.CONTINUE_OBJ {
				if fs.Update != nil {
					e.Eval(fs.Update, env)
				}

				continue
			} else if resultType == object.BREAK_OBJ {
				break
			} else if resultType == object.RETURN_VALUE_OBJ || resultType == object.ERROR_OBJ {
				return result
			}
		}

		if fs.Update != nil {
			e.Eval(fs.Update, env)
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
		return object.NewInteger(right.Node(), -value)

	case object.FLOAT_OBJ:
		value := right.(*object.Float).Value
		return object.NewFloat(right.Node(), -value)

	default:
		return object.NewError(right.Node(), "unknown operator: -%s", right.Type())
	}
}

func (e *Evaluator) evalPrefixExpression(node *ast.PrefixExpression, right object.Object) object.Object {
	switch node.Operator {
	case "!":
		return e.evalExclamationOperatorExpression(right)
	case "-":
		return e.evalMinusPrefixOperatorExpression(right)
	default:
		return object.NewError(node, "unknown operator: %s%s", node.Operator, right.Type())
	}
}

func (e *Evaluator) evalIntegerInfixExpression(node *ast.InfixExpression, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch node.Operator {
	case "+":
		return object.NewInteger(node, leftVal+rightVal)
	case "-":
		return object.NewInteger(node, leftVal-rightVal)
	case "*":
		return object.NewInteger(node, leftVal*rightVal)
	case "/":
		return object.NewFloat(node, float64(leftVal)/float64(rightVal))
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
		return object.NewError(node, "unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
	}
}

func (e *Evaluator) evalFloatInfixExpression(node *ast.InfixExpression, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch node.Operator {
	case "+":
		return object.NewFloat(node, leftVal+rightVal)
	case "-":
		return object.NewFloat(node, leftVal-rightVal)
	case "*":
		return object.NewFloat(node, leftVal*rightVal)
	case "/":
		return object.NewFloat(node, leftVal/rightVal)
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
		return object.NewError(node, "unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
	}
}

func (e *Evaluator) evalNumberInfixExpression(node *ast.InfixExpression, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return e.evalIntegerInfixExpression(node, left, right)

	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return e.evalFloatInfixExpression(node, left, right)

	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		leftValue := float64(left.(*object.Integer).Value)

		return e.evalFloatInfixExpression(node, object.NewFloat(left.Node(), leftValue), right)

	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		rightValue := float64(right.(*object.Integer).Value)

		return e.evalFloatInfixExpression(node, left, object.NewFloat(right.Node(), rightValue))

	default:
		return object.NewError(node, "type mismatch: %s %s %s", left.Type(), node.Operator, right.Type())
	}
}

func (e *Evaluator) evalStringInfixExpression(node *ast.InfixExpression, left, right object.Object) object.Object {
	if node.Operator != "+" {
		return object.NewError(node, "unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	return object.NewString(node, leftVal+rightVal)
}

func (e *Evaluator) evalInfixExpression(node *ast.InfixExpression, left, right object.Object) object.Object {
	switch {
	case object.IsNumber(left) && object.IsNumber(right):
		return e.evalNumberInfixExpression(node, left, right)

	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return e.evalStringInfixExpression(node, left, right)

	case node.Operator == "==":
		return e.nativeBoolToBooleanObject(left == right)

	case node.Operator == "!=":
		return e.nativeBoolToBooleanObject(left != right)

	case left.Type() != right.Type():
		return object.NewError(node, "type mismatch: %s %s %s", left.Type(), node.Operator, right.Type())

	default:
		return object.NewError(node, "unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
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

	return object.NewError(node, "identifier not found: %s", node.Value)
}

func (e *Evaluator) evalExpressions(expressions []ast.Expression, env *object.Environment) ([]object.Object, *object.Error) {
	var result []object.Object

	for _, expression := range expressions {
		evaluated := e.Eval(expression, env)

		if isError(evaluated) {
			return []object.Object{}, evaluated.(*object.Error)
		}

		result = append(result, evaluated)
	}

	return result, nil
}

func (e *Evaluator) extendFunctionEnv(function *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(function.Env)

	for paramIdx, param := range function.Arguments {
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

func (e *Evaluator) applyFunction(node *ast.CallExpression, function object.Object, args []object.Object) object.Object {
	switch function := function.(type) {
	case *object.Function:
		extendedEnv := e.extendFunctionEnv(function, args)
		evaluated := e.Eval(function.Body, extendedEnv)

		return e.unwrapReturnValue(evaluated)
	case *object.Builtin:
		return function.Function(node, args...)
	}

	return object.NewError(node, "not a function: %s", function.Type())
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
		fmt.Println(index.Type())
		return object.NewError(index.Node(), "unusable as object key: %s", index.Type())
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
		return object.NewError(index.Node(), "index operator not supported: %s", left.Type())
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
			return object.NewError(node, "unusable as object key: %s", key.Type())
		}

		value := e.Eval(valueNode, env)

		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return object.NewHash(node, pairs)
}

func (e *Evaluator) evalChainingCallExpression(left ast.Node, rightCallExpression *ast.CallExpression, env *object.Environment) object.Object {
	leftValue := e.Eval(left, env)
	args, err := e.evalExpressions(rightCallExpression.Arguments, env)

	if err != nil {
		return err
	}

	switch leftValue := leftValue.(type) {
	case *object.String:
		chainingFunction, ok := e.stringChainingFunctions[rightCallExpression.Function.TokenLiteral()]

		if !ok {
			return object.NewError(rightCallExpression.Function, "String has no method %s", rightCallExpression.Function.TokenLiteral())
		}

		return chainingFunction(leftValue, args...)
	case *object.Array:
		chainingFunction, ok := e.arrayChainingFunctions[rightCallExpression.Function.TokenLiteral()]

		if !ok {
			return object.NewError(rightCallExpression.Function, "Array has no method %s", rightCallExpression.Function.TokenLiteral())
		}

		return chainingFunction(leftValue, args...)
	case *object.Hash:
		chainingFunction, ok := e.objectChainingFunctions[rightCallExpression.Function.TokenLiteral()]

		if !ok {
			return object.NewError(rightCallExpression.Function, "Object has no method %s", rightCallExpression.Function.TokenLiteral())
		}

		return chainingFunction(leftValue, args...)
	}

	return object.NewError(rightCallExpression.Function, "chaining operator not supported: %s.%s", leftValue.Type(), rightCallExpression.Function.TokenLiteral())
}

func (e *Evaluator) evalChainingExpression(left ast.Node, right ast.Node, env *object.Environment) object.Object {
	switch right := right.(type) {
	case *ast.CallExpression:
		return e.evalChainingCallExpression(left, right, env)
	case *ast.ExpressionStatement:
		switch rightExpression := right.Expression.(type) {
		case *ast.CallExpression:
			return e.evalChainingCallExpression(left, rightExpression, env)

		case *ast.Identifier:
			// TODO: implement chaining ident operator
			return object.NewError(right, "chaining ident operator not supported: %s.%s", left.TokenLiteral(), rightExpression.Value)
		}
	}

	leftValue := e.Eval(left, env)

	return object.NewError(right, "chaining operator not supported: %s.%s", leftValue.Type(), right.String())
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
		value := e.Eval(node.ReturnValue, env)

		if isError(value) {
			return value
		}

		return object.NewReturnValue(node, value)

	case *ast.VariableStatement:
		_, defined := env.GetFromCurrent(node.Name.Value)

		if defined {
			return object.NewError(node, "identifier already defined: %s", node.Name.Value)
		}

		value := e.Eval(node.Value, env)

		if isError(value) {
			return value
		}

		env.Set(node.Name.Value, value)

	case *ast.FunctionStatement:
		_, defined := env.GetFromCurrent(node.Name.Value)

		if defined {
			return object.NewError(node, "identifier already defined: %s", node.Name.Value)
		}

		params := node.Parameters
		body := node.Body

		function := object.NewFunction(node, params, body, env)

		env.Set(node.Name.Value, function)

	case *ast.WhileStatement:
		return e.evalWhileStatement(node, env)

	case *ast.ForStatement:
		return e.evalForStatement(node, env)

	// Expressions
	case *ast.IntegerLiteral:
		return object.NewInteger(node, node.Value)

	case *ast.BooleanLiteral:
		return e.nativeBoolToBooleanObject(node.Value)

	case *ast.FloatLiteral:
		return object.NewFloat(node, node.Value)

	case *ast.NullLiteral:
		return object.NewNull(node)

	case *ast.StringLiteral:
		return object.NewString(node, node.Value)

	case *ast.ArrayLiteral:
		elements, err := e.evalExpressions(node.Elements, env)

		if err != nil {
			return err
		}

		return object.NewArray(node, elements)

	case *ast.HashLiteral:
		return e.evalHashLiteral(node, env)

	case *ast.PrefixExpression:
		right := e.Eval(node.Right, env)

		if isError(right) {
			return right
		}

		return e.evalPrefixExpression(node, right)

	case *ast.InfixExpression:
		left := e.Eval(node.Left, env)

		if isError(left) {
			return left
		}

		right := e.Eval(node.Right, env)

		if isError(right) {
			return right
		}

		return e.evalInfixExpression(node, left, right)

	case *ast.IfExpression:
		return e.evalIfExpression(node, env)

	case *ast.BreakExpression:
		return object.NewBreak(node)

	case *ast.ContinueExpression:
		return object.NewContinue(node)

	case *ast.Identifier:
		return e.evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body

		return object.NewFunction(node, params, body, env)

	case *ast.CallExpression:
		function := e.Eval(node.Function, env)

		if isError(function) {
			return function
		}

		args, err := e.evalExpressions(node.Arguments, env)

		if err != nil {
			return err
		}

		return e.applyFunction(node, function, args)

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
			return object.NewError(node, "variable %s has not been initialized.", node.Name.Value)
		}

		value := e.Eval(node.Value, env)

		if isError(value) {
			return value
		}

		environment.Set(node.Name.Value, value)
	case *ast.ChainingExpression:
		return e.evalChainingExpression(node.Left, node.Right, env)
	}

	return nil
}
