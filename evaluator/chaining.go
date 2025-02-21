package evaluator

import (
	"github.com/iskandervdh/vorn/ast"
	"github.com/iskandervdh/vorn/evaluator/chaining_functions"
	"github.com/iskandervdh/vorn/object"
)

func (e *Evaluator) evalChainingCallExpression(left ast.Node, rightCallExpression *ast.CallExpression, env *object.Environment) object.Object {
	leftValue := e.Eval(left, env)
	args, err := e.evalExpressions(rightCallExpression.Arguments, env)

	if err != nil {
		return err
	}

	switch leftValue := leftValue.(type) {
	case *object.String:
		chainingFunction, ok := chaining_functions.StringChainingFunctions[rightCallExpression.Function.TokenLiteral()]

		if !ok {
			return object.NewError(rightCallExpression.Function, "String has no method %s", rightCallExpression.Function.TokenLiteral())
		}

		return chainingFunction(leftValue, args...)
	case *object.Array:
		chainingFunction, ok := chaining_functions.ArrayChainingFunctions[rightCallExpression.Function.TokenLiteral()]

		if !ok {
			return object.NewError(rightCallExpression.Function, "Array has no method %s", rightCallExpression.Function.TokenLiteral())
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

	return object.NewError(right, "chaining operator not supported: %s.%s", left.TokenLiteral(), right.String())
}
