package evaluator

import (
	"sort"
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
		return f.(*object.Builtin).ArgumentsCount, nil
	} else if f.Type() == object.FUNCTION_OBJ {
		return len(f.(*object.Function).Arguments), nil
	} else {
		return -1, object.NewError(f.Node(), "%s callback must be a function, got %s", chainingExpression, f.Type())
	}
}

func (e *Evaluator) sortObjects(elements []object.Object, reverse bool) []object.Object {
	// Copy the elements to a new slice to avoid modifying the original array
	sorted := make([]object.Object, len(elements))

	copy(sorted, elements)

	sort.Slice(sorted, func(i, j int) bool {
		if reverse {
			return sorted[i].Inspect() > sorted[j].Inspect()
		}

		return sorted[i].Inspect() < sorted[j].Inspect()
	})

	return sorted
}

func (e *Evaluator) arrayLength(arr *object.Array, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(arr.Node(), "Array.length() takes no arguments")
	}

	return object.NewInteger(arr.Node(), int64(len(arr.Elements)))
}

func (e *Evaluator) arrayPrepend(arr *object.Array, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(arr.Node(), "Array.prepend() takes exactly 1 argument")
	}

	arr.Elements = append([]object.Object{args[0]}, arr.Elements...)

	return arr
}

func (e *Evaluator) arrayAppend(arr *object.Array, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(arr.Node(), "Array.append() takes exactly 1 argument")
	}

	arr.Elements = append(arr.Elements, args[0])

	return arr
}

func (e *Evaluator) arrayShift(arr *object.Array, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(arr.Node(), "Array.shift() takes no arguments")
	}

	if len(arr.Elements) == 0 {
		return object.NewError(arr.Node(), "Array.shift() called on empty array")
	}

	shifted := arr.Elements[0]
	arr.Elements = arr.Elements[1:]

	return shifted
}

func (e *Evaluator) arrayPop(arr *object.Array, args ...object.Object) object.Object {
	if len(args) > 1 {
		return object.NewError(arr.Node(), "Array.pop() 0 or 1 argument")
	}

	if len(arr.Elements) == 0 {
		return object.NewError(arr.Node(), "Array.pop() called on empty array")
	}

	if len(args) == 1 {
		if _, ok := args[0].(*object.Integer); !ok {
			return object.NewError(arr.Node(), "Array.pop() argument must be an integer")
		}

		index := args[0].(*object.Integer).Value

		if index < 0 || int(index) >= len(arr.Elements) {
			return object.NewError(arr.Node(), "Array.pop() index out of range")
		}

		popped := arr.Elements[index]
		arr.Elements = append(arr.Elements[:index], arr.Elements[index+1:]...)

		return popped
	}

	popped := arr.Elements[len(arr.Elements)-1]
	arr.Elements = arr.Elements[:len(arr.Elements)-1]

	return popped
}

func (e *Evaluator) arrayConcat(arr *object.Array, args ...object.Object) object.Object {
	if len(args) == 0 {
		return object.NewError(arr.Node(), "Array.concat() takes at least 1 argument")
	}

	concatenated := make([]object.Object, len(arr.Elements))

	copy(concatenated, arr.Elements)

	for _, arg := range args {
		if arg.Type() != object.ARRAY_OBJ {
			return object.NewError(arr.Node(), "argument to `concat` must be ARRAY, got %s", arg.Type())
		}

		concatenated = append(concatenated, arg.(*object.Array).Elements...)
	}

	return object.NewArray(arr.Node(), concatenated)
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
		args := []object.Object{el, object.NewInteger(arr.Node(), int64(i)), arr}

		// Handle functions that require only a certain amount of arguments
		value := e.applyFunction(callExpression, f, args[:callbackArgumentsCount])

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
		args := []object.Object{el, object.NewInteger(arr.Node(), int64(i)), arr}

		// Handle functions that require only a certain amount of arguments
		value := e.applyFunction(callExpression, f, args[:callbackArgumentsCount])

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
		args := []object.Object{accumulator, el, object.NewInteger(arr.Node(), int64(i)), arr}

		// Handle functions that require only a certain amount of arguments
		accumulator = e.applyFunction(callExpression, f, args[:callbackArgumentsCount])

		// If the value is an error, return the error
		if _, ok := accumulator.(*object.Error); ok {
			return accumulator
		}
	}

	return accumulator
}

func (e *Evaluator) arrayContains(arr *object.Array, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(arr.Node(), "Array.contains() takes exactly 1 argument")
	}

	for _, el := range arr.Elements {
		if el.Inspect() == args[0].Inspect() {
			return TRUE
		}
	}

	return FALSE
}

func (e *Evaluator) arrayIndexOf(arr *object.Array, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(arr.Node(), "Array.indexOf() takes exactly 1 argument")
	}

	for i, el := range arr.Elements {
		if el.Inspect() == args[0].Inspect() {
			return object.NewInteger(arr.Node(), int64(i))
		}
	}

	return object.NewInteger(arr.Node(), -1)
}

func (e *Evaluator) arrayFind(arr *object.Array, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(arr.Node(), "Array.find() takes exactly 1 argument")
	}

	f := args[0]

	callbackArgumentsCount, err := getCallbackArgumentsCount(f, "Array.find()")

	if err != nil {
		return err
	}

	if callbackArgumentsCount < 1 {
		return object.NewError(f.Node(), "Array.find() callback must take at least 1 argument")
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
			Literal: "find",
			Line:    arr.Node().Line(),
			Column:  arr.Node().Column(),
		}, Value: "find"},
		Arguments: []ast.Expression{},
	}

	for i, el := range arr.Elements {
		args := []object.Object{el, object.NewInteger(arr.Node(), int64(i)), arr}

		// Handle functions that require only a certain amount of arguments
		value := e.applyFunction(callExpression, f, args[:callbackArgumentsCount])

		// If the value is an error, return the error
		if _, ok := value.(*object.Error); ok {
			return value
		}

		// Check if value cast to boolean is true
		if value == TRUE {
			return el
		} else if value != FALSE && e.builtinBool(callExpression, value) == TRUE {
			return el
		}
	}

	return NULL
}

func (e *Evaluator) arrayJoin(arr *object.Array, args ...object.Object) object.Object {
	if len(args) > 1 {
		return object.NewError(arr.Node(), "Array.join() takes at most 1 argument, got %d", len(args))
	}

	separator := ", "

	if len(args) == 1 {
		if args[0].Type() != object.STRING_OBJ {
			return object.NewError(arr.Node(), "argument to `join` must be STRING, got %s", args[0].Type())
		}

		separator = args[0].(*object.String).Value
	}

	elements := make([]string, len(arr.Elements))

	for i, el := range arr.Elements {
		elements[i] = el.Inspect()
	}

	return object.NewString(arr.Node(), strings.Join(elements, separator))
}

func (e *Evaluator) arrayReverse(arr *object.Array, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(arr.Node(), "Array.reverse() takes no arguments")
	}

	reversed := make([]object.Object, len(arr.Elements))

	for i, j := 0, len(arr.Elements)-1; i < len(arr.Elements); i, j = i+1, j-1 {
		reversed[i] = arr.Elements[j]
	}

	arr.Elements = reversed

	return arr
}

func (e *Evaluator) arraySlice(arr *object.Array, args ...object.Object) object.Object {
	if len(args) == 0 || len(args) > 2 {
		return object.NewError(arr.Node(), "Array.slice() takes 1 or 2 arguments")
	}

	if args[0].Type() != object.INTEGER_OBJ {
		return object.NewError(arr.Node(), "first argument to `slice` must be INTEGER, got %s", args[0].Type())
	}

	start := int(args[0].(*object.Integer).Value)
	end := len(arr.Elements)

	if len(args) == 2 {
		if args[1].Type() != object.INTEGER_OBJ {
			return object.NewError(arr.Node(), "second argument to `slice` must be INTEGER, got %s", args[1].Type())
		}

		end = int(args[1].(*object.Integer).Value)
	}

	if start < 0 || start > len(arr.Elements) {
		return object.NewError(arr.Node(), "first argument to `slice` out of range")
	}

	if end < 0 || end > len(arr.Elements) {
		return object.NewError(arr.Node(), "second argument to `slice` out of range")
	}

	return object.NewArray(arr.Node(), arr.Elements[start:end])
}

func (e *Evaluator) arraySort(arr *object.Array, args ...object.Object) object.Object {
	if len(args) > 1 {
		return object.NewError(arr.Node(), "Array.sort() takes at most 1 argument, got %d", len(args))
	}

	// If no arguments are passed, sort the array in ascending order
	if len(args) == 0 {
		arr.Elements = e.sortObjects(arr.Elements, false)

		return arr
	}

	if len(args) == 1 {
		if args[0].Type() == object.BOOLEAN_OBJ {
			arr.Elements = e.sortObjects(arr.Elements, args[0].(*object.Boolean).Value)

			return arr
		} else if args[0].Type() != object.FUNCTION_OBJ && args[0].Type() != object.BUILTIN_OBJ {
			return object.NewError(arr.Node(), "argument to `sort` must be FUNCTION or BUILTIN, got %s", args[0].Type())
		}
	}

	f := args[0]
	callbackArgumentsCount, err := getCallbackArgumentsCount(f, "Array.sort()")

	if err != nil {
		return err
	}

	if callbackArgumentsCount < 2 {
		return object.NewError(f.Node(), "Array.sort() callback must take at least 2 arguments")
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
			Literal: "sort",
			Line:    arr.Node().Line(),
			Column:  arr.Node().Column(),
		}, Value: "sort"},
		Arguments: []ast.Expression{},
	}

	sorted := make([]object.Object, len(arr.Elements))
	copy(sorted, arr.Elements)

	sort.Slice(arr.Elements, func(i, j int) bool {
		args := []object.Object{arr.Elements[i], arr.Elements[j], object.NewInteger(arr.Node(), int64(i)), object.NewInteger(arr.Node(), int64(j)), arr}

		// Handle functions that require only a certain amount of arguments
		value := e.applyFunction(callExpression, f, args[:callbackArgumentsCount])

		// If the value is an error, return the error
		if _, ok := value.(*object.Error); ok {
			return false
		}

		return e.builtinBool(callExpression, value) == TRUE
	})

	return arr
}

func (e *Evaluator) arrayAny(arr *object.Array, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(arr.Node(), "Array.any() takes exactly 1 argument")
	}

	f := args[0]
	callbackArgumentsCount, err := getCallbackArgumentsCount(f, "Array.any()")

	if err != nil {
		return err
	}

	if callbackArgumentsCount < 1 {
		return object.NewError(f.Node(), "Array.any() callback must take at least 1 argument")
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
			Literal: "any",
			Line:    arr.Node().Line(),
			Column:  arr.Node().Column(),
		}, Value: "any"},
		Arguments: []ast.Expression{},
	}

	for i, el := range arr.Elements {
		args := []object.Object{el, object.NewInteger(arr.Node(), int64(i)), arr}

		// Handle functions that require only a certain amount of arguments
		value := e.applyFunction(callExpression, f, args[:callbackArgumentsCount])

		// If the value is an error, return the error
		if _, ok := value.(*object.Error); ok {
			return value
		}

		// Check if value cast to boolean is true
		if value == TRUE {
			return TRUE
		} else if value != FALSE && e.builtinBool(callExpression, value) == TRUE {
			return TRUE
		}
	}

	return FALSE
}

func (e *Evaluator) arrayAll(arr *object.Array, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(arr.Node(), "Array.all() takes exactly 1 argument")
	}

	f := args[0]
	callbackArgumentsCount, err := getCallbackArgumentsCount(f, "Array.all()")

	if err != nil {
		return err
	}

	if callbackArgumentsCount < 1 {
		return object.NewError(f.Node(), "Array.all() callback must take at least 1 argument")
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
			Literal: "all",
			Line:    arr.Node().Line(),
			Column:  arr.Node().Column(),
		}, Value: "all"},
		Arguments: []ast.Expression{},
	}

	for i, el := range arr.Elements {
		args := []object.Object{el, object.NewInteger(arr.Node(), int64(i)), arr}

		// Handle functions that require only a certain amount of arguments
		value := e.applyFunction(callExpression, f, args[:callbackArgumentsCount])

		// If the value is an error, return the error
		if _, ok := value.(*object.Error); ok {
			return value
		}

		// Check if value cast to boolean is true
		if value == FALSE {
			return FALSE
		} else if value != TRUE && e.builtinBool(callExpression, value) == FALSE {
			return FALSE
		}
	}

	return TRUE
}

/**
 * String chaining functions
 */

func (e *Evaluator) stringLength(str *object.String, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(str.Node(), "String.length() takes no arguments")
	}

	return object.NewInteger(str.Node(), int64(len(str.Value)))
}

func (e *Evaluator) stringUpper(str *object.String, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(str.Node(), "String.upper() takes no arguments")
	}

	return object.NewString(str.Node(), strings.ToUpper(str.Value))
}

func (e *Evaluator) stringLower(str *object.String, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(str.Node(), "String.lower() takes no arguments")
	}

	return object.NewString(str.Node(), strings.ToLower(str.Value))
}

func (e *Evaluator) stringSplit(str *object.String, args ...object.Object) object.Object {
	if len(args) > 1 {
		return object.NewError(str.Node(), "String.split() takes at most 1 argument, got %d", len(args))
	}

	separator := " "

	if len(args) == 1 {
		if args[0].Type() != object.STRING_OBJ {
			return object.NewError(str.Node(), "argument to `split` must be STRING, got %s", args[0].Type())
		}

		separator = args[0].(*object.String).Value
	}

	parts := strings.Split(str.Value, separator)
	elements := make([]object.Object, len(parts))

	for i, part := range parts {
		elements[i] = object.NewString(str.Node(), part)
	}

	return object.NewArray(str.Node(), elements)
}

func (e *Evaluator) stringContains(str *object.String, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(str.Node(), "String.contains() takes exactly 1 argument")
	}

	if args[0].Type() != object.STRING_OBJ {
		return object.NewError(str.Node(), "argument to `contains` must be STRING, got %s", args[0].Type())
	}

	return e.nativeBoolToBooleanObject(strings.Contains(str.Value, args[0].(*object.String).Value))
}

func (e *Evaluator) stringReplace(str *object.String, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewError(str.Node(), "String.replace() takes exactly 2 arguments")
	}

	if args[0].Type() != object.STRING_OBJ {
		return object.NewError(str.Node(), "first argument to `replace` must be STRING, got %s", args[0].Type())
	}

	if args[1].Type() != object.STRING_OBJ {
		return object.NewError(str.Node(), "second argument to `replace` must be STRING, got %s", args[1].Type())
	}

	return object.NewString(str.Node(), strings.ReplaceAll(str.Value, args[0].(*object.String).Value, args[1].(*object.String).Value))
}

func (e *Evaluator) stringTrim(str *object.String, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(str.Node(), "String.trim() takes no arguments")
	}

	return object.NewString(str.Node(), strings.TrimSpace(str.Value))
}

func (e *Evaluator) stringTrimStart(str *object.String, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(str.Node(), "String.trimStart() takes no arguments")
	}

	return object.NewString(str.Node(), strings.TrimLeft(str.Value, " \t\n\r"))
}

func (e *Evaluator) stringTrimEnd(str *object.String, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(str.Node(), "String.trimEnd() takes no arguments")
	}

	return object.NewString(str.Node(), strings.TrimRight(str.Value, " \t\n\r"))
}

func (e *Evaluator) stringRepeat(str *object.String, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(str.Node(), "String.repeat() takes exactly 1 argument")
	}

	if args[0].Type() != object.INTEGER_OBJ {
		return object.NewError(str.Node(), "argument to `repeat` must be INTEGER, got %s", args[0].Type())
	}

	return object.NewString(str.Node(), strings.Repeat(str.Value, int(args[0].(*object.Integer).Value)))
}

func (e *Evaluator) stringReverse(str *object.String, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(str.Node(), "String.reverse() takes no arguments")
	}

	runes := []rune(str.Value)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return object.NewString(str.Node(), string(runes))
}

func (e *Evaluator) stringSlice(str *object.String, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewError(str.Node(), "String.slice() takes exactly 2 arguments")
	}

	if args[0].Type() != object.INTEGER_OBJ {
		return object.NewError(str.Node(), "first argument to `slice` must be INTEGER, got %s", args[0].Type())
	}

	if args[1].Type() != object.INTEGER_OBJ {
		return object.NewError(str.Node(), "second argument to `slice` must be INTEGER, got %s", args[1].Type())
	}

	start := int(args[0].(*object.Integer).Value)
	end := int(args[1].(*object.Integer).Value)

	if start < 0 || start > len(str.Value) {
		return object.NewError(str.Node(), "first argument to `slice` out of range")
	}

	if end < 0 || end > len(str.Value) {
		return object.NewError(str.Node(), "second argument to `slice` out of range")
	}

	return object.NewString(str.Node(), str.Value[start:end])
}

func (e *Evaluator) stringStartsWith(str *object.String, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(str.Node(), "String.startsWith() takes exactly 1 argument")
	}

	if args[0].Type() != object.STRING_OBJ {
		return object.NewError(str.Node(), "argument to `startsWith` must be STRING, got %s", args[0].Type())
	}

	return e.nativeBoolToBooleanObject(strings.HasPrefix(str.Value, args[0].(*object.String).Value))
}

func (e *Evaluator) stringEndsWith(str *object.String, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(str.Node(), "String.endsWith() takes exactly 1 argument")
	}

	if args[0].Type() != object.STRING_OBJ {
		return object.NewError(str.Node(), "argument to `endsWith` must be STRING, got %s", args[0].Type())
	}

	return e.nativeBoolToBooleanObject(strings.HasSuffix(str.Value, args[0].(*object.String).Value))
}

/**
 * Object chaining functions
 */

func (e *Evaluator) objectKeys(hash *object.Hash, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(hash.Node(), "Object.keys() takes no arguments")
	}

	keys := make([]object.Object, len(hash.Pairs))

	i := 0

	for key := range hash.Pairs {
		keys[i] = object.NewString(hash.Node(), hash.Pairs[key].Key.Inspect())
		i++
	}

	return object.NewArray(hash.Node(), keys)
}

func (e *Evaluator) objectValues(hash *object.Hash, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(hash.Node(), "Object.values() takes no arguments")
	}

	values := make([]object.Object, len(hash.Pairs))

	i := 0

	for _, pair := range hash.Pairs {
		values[i] = pair.Value

		i++
	}

	return object.NewArray(hash.Node(), values)
}

func (e *Evaluator) objectItems(hash *object.Hash, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(hash.Node(), "Object.items() takes no arguments")
	}

	items := make([]object.Object, len(hash.Pairs))

	i := 0

	for key, pair := range hash.Pairs {
		items[i] = object.NewArray(hash.Node(), []object.Object{object.NewString(hash.Node(), hash.Pairs[key].Key.Inspect()), pair.Value})

		i++
	}

	return object.NewArray(hash.Node(), items)
}
