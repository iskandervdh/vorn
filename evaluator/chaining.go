package evaluator

import (
	"strings"

	"github.com/iskandervdh/vorn/ast"
	"github.com/iskandervdh/vorn/object"
	"github.com/iskandervdh/vorn/token"
)

/**
 * Array chaining functions
 */

func getCallbackArgumentsCount(f object.Object, chainingExpression string) (int, *object.Error) {
	if f.Type() == object.BUILTIN_OBJ {
		return f.(*object.Builtin).ArgumentsCount[0], nil
	} else if f.Type() == object.FUNCTION_OBJ {
		return len(f.(*object.Function).Arguments), nil
	} else {
		return -1, object.NewError(f.Node(), "%s callback must be a function, got %s", chainingExpression, f.Type())
	}
}

func (e *Evaluator) arrayLength(left *object.Array, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(left.Node(), "Array.length() takes no arguments")
	}

	return object.NewInteger(left.Node(), int64(len(left.Elements)))
}

func (e *Evaluator) arrayPush(left *object.Array, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(left.Node(), "Array.push() takes exactly 1 argument")
	}

	left.Elements = append(left.Elements, args[0])

	return left
}

func (e *Evaluator) arrayPop(left *object.Array, args ...object.Object) object.Object {
	if len(args) > 1 {
		return object.NewError(left.Node(), "Array.pop() 0 or 1 argument")
	}

	if len(left.Elements) == 0 {
		return object.NewError(left.Node(), "Array.pop() called on empty array")
	}

	if len(args) == 1 {
		if _, ok := args[0].(*object.Integer); !ok {
			return object.NewError(left.Node(), "Array.pop() argument must be an integer")
		}

		index := args[0].(*object.Integer).Value

		if index < 0 || int(index) >= len(left.Elements) {
			return object.NewError(left.Node(), "Array.pop() index out of range")
		}

		popped := left.Elements[index]
		left.Elements = append(left.Elements[:index], left.Elements[index+1:]...)

		return popped
	}

	popped := left.Elements[len(left.Elements)-1]
	left.Elements = left.Elements[:len(left.Elements)-1]

	return popped
}

func (e *Evaluator) arrayMap(arr *object.Array, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(arr.Node(), "Array.map() takes exactly 1 argument")
	}

	f := args[0]
	callbackArgumentsCount, err := getCallbackArgumentsCount(f, "Array.map()")

	if err != nil {
		return err
	}

	if callbackArgumentsCount < 1 {
		return object.NewError(f.Node(), "Array.map() callback must take at least 1 argument")
	}

	newArray := object.NewArray(arr.Node(), make([]object.Object, len(arr.Elements)))

	// Create a call expression to pass to the applyFunction method to have the correct line and column numbers for errors
	callExpression := &ast.CallExpression{
		Token: token.Token{
			Type:    token.LPAREN,
			Literal: "(",
			Line:    arr.Node().Line(),
			Column:  arr.Node().Column(),
		},
		Function: &ast.Identifier{Token: token.Token{
			Type:    token.IDENT,
			Literal: "map",
			Line:    arr.Node().Line(),
			Column:  arr.Node().Column(),
		}, Value: "map"},
		Arguments: []ast.Expression{},
	}

	for i, el := range arr.Elements {
		arguments := []object.Object{el, object.NewInteger(arr.Node(), int64(i)), arr}

		// Handle functions that require only a certain amount of arguments
		value := e.applyFunction(callExpression, f, arguments[:callbackArgumentsCount])

		// If the value is an error, return the error
		if _, ok := value.(*object.Error); ok {
			return value
		}

		newArray.Elements[i] = value
	}

	return newArray
}

func (e *Evaluator) arrayFilter(arr *object.Array, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(arr.Node(), "Array.filter() takes exactly 1 argument")
	}

	f := args[0]
	callbackArgumentsCount, err := getCallbackArgumentsCount(f, "Array.filter()")

	if err != nil {
		return err
	}

	if callbackArgumentsCount < 1 {
		return object.NewError(f.Node(), "Array.filter() callback must take at least 1 argument")
	}

	newArray := object.NewArray(arr.Node(), []object.Object{})

	// Create a call expression to pass to the applyFunction method to have the correct line and column numbers for errors
	callExpression := &ast.CallExpression{
		Token: token.Token{
			Type:    token.LPAREN,
			Literal: "(",
			Line:    arr.Node().Line(),
			Column:  arr.Node().Column(),
		},
		Function: &ast.Identifier{Token: token.Token{
			Type:    token.IDENT,
			Literal: "filter",
			Line:    arr.Node().Line(),
			Column:  arr.Node().Column(),
		}, Value: "filter"},
		Arguments: []ast.Expression{},
	}

	for i, el := range arr.Elements {
		arguments := []object.Object{el, object.NewInteger(arr.Node(), int64(i)), arr}

		// Handle functions that require only a certain amount of arguments
		value := e.applyFunction(callExpression, f, arguments[:callbackArgumentsCount])

		if _, ok := value.(*object.Error); ok {
			return value
		}

		if value == TRUE {
			newArray.Elements = append(newArray.Elements, el)
			// Check if value cast to boolean is true
		} else if value != FALSE && e.builtinBool(callExpression, value) == TRUE {
			newArray.Elements = append(newArray.Elements, el)
		}
	}

	return newArray
}

func (e *Evaluator) arrayReduce(arr *object.Array, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewError(arr.Node(), "Array.reduce() takes exactly 2 arguments, got %d", len(args))
	}

	f := args[0]
	accumulator := args[1]
	callbackArgumentsCount, err := getCallbackArgumentsCount(f, "Array.reduce()")

	if err != nil {
		return err
	}

	if callbackArgumentsCount < 2 {
		return object.NewError(f.Node(), "Array.reduce() callback must take at least 2 arguments")
	}

	// Create a call expression to pass to the applyFunction method to have the correct line and column numbers for errors
	callExpression := &ast.CallExpression{
		Token: token.Token{
			Type:    token.LPAREN,
			Literal: "(",
			Line:    arr.Node().Line(),
			Column:  arr.Node().Column(),
		},
		Function: &ast.Identifier{Token: token.Token{
			Type:    token.IDENT,
			Literal: "reduce",
			Line:    arr.Node().Line(),
			Column:  arr.Node().Column(),
		}, Value: "reduce"},
		Arguments: []ast.Expression{},
	}

	for i, el := range arr.Elements {
		arguments := []object.Object{accumulator, el, object.NewInteger(arr.Node(), int64(i)), arr}

		// Handle functions that require only a certain amount of arguments
		accumulator = e.applyFunction(callExpression, f, arguments[:callbackArgumentsCount])

		// If the value is an error, return the error
		if _, ok := accumulator.(*object.Error); ok {
			return accumulator
		}
	}

	return accumulator
}

/**
 * String chaining functions
 */

func (e *Evaluator) stringLength(left *object.String, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(left.Node(), "String.length() takes no arguments")
	}

	return object.NewInteger(left.Node(), int64(len(left.Value)))
}

func (e *Evaluator) stringUpper(left *object.String, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(left.Node(), "String.upper() takes no arguments")
	}

	return object.NewString(left.Node(), strings.ToUpper(left.Value))
}

func (e *Evaluator) stringLower(left *object.String, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(left.Node(), "String.lower() takes no arguments")
	}

	return object.NewString(left.Node(), strings.ToLower(left.Value))
}

func (e *Evaluator) stringSplit(left *object.String, args ...object.Object) object.Object {
	if len(args) > 1 {
		return object.NewError(left.Node(), "String.split() takes at most 1 argument, got %d", len(args))
	}

	separator := " "

	if len(args) == 1 {
		if args[0].Type() != object.STRING_OBJ {
			return object.NewError(left.Node(), "argument to `split` must be STRING, got %s", args[0].Type())
		}

		separator = args[0].(*object.String).Value
	}

	parts := strings.Split(left.Value, separator)
	elements := make([]object.Object, len(parts))

	for i, part := range parts {
		elements[i] = object.NewString(left.Node(), part)
	}

	return object.NewArray(left.Node(), elements)
}
