package evaluator

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/iskandervdh/vorn/ast"
	"github.com/iskandervdh/vorn/object"
)

// Common functions

func (e *Evaluator) builtinType(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	typeName := args[0].Type()

	return object.NewString(node, string(typeName))
}

func (e *Evaluator) builtinRange(node ast.Node, args ...object.Object) object.Object {
	if len(args) < 1 || len(args) > 2 {
		return newError(node, "wrong number of arguments. got %d, want 1 or 2", len(args))
	}

	if args[0].Type() != object.INTEGER_OBJ {
		return newError(node, "first argument to `range` must be INTEGER, got %s", args[0].Type())
	}

	firstArg := args[0].(*object.Integer)

	start := firstArg.Value
	end := start

	if len(args) == 1 {
		if start < 0 {
			return newError(node, "argument to `range` must be non-negative, got %d", start)
		}

		start = 0
	} else {
		if args[1].Type() != object.INTEGER_OBJ {
			return newError(node, "second argument to `range` must be INTEGER, got %s", args[1].Type())
		}

		end = args[1].(*object.Integer).Value
	}

	elementsLength := end - start

	if start > end {
		elementsLength = start - end
	}

	elements := make([]object.Object, elementsLength)

	if start > end {
		for i := start; i > end; i-- {
			elements[start-i] = object.NewInteger(firstArg.Node(), i)
		}
	} else {
		for i := start; i < end; i++ {
			elements[i-start] = object.NewInteger(firstArg.Node(), i)
		}
	}

	return object.NewArray(firstArg.Node(), elements)
}

// Conversion functions

func (e *Evaluator) builtinInt(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		return arg
	case *object.Float:
		return object.NewInteger(arg.Node(), int64(arg.Value))
	case *object.String:
		integer, err := strconv.ParseInt(arg.Value, 0, 64)

		if err != nil {
			return newError(node, "could not parse %q as INTEGER", arg.Value)
		}

		return object.NewInteger(arg.Node(), integer)
	default:
		return newError(node, "argument to `int` not supported, got %s", args[0].Type())
	}
}

func (e *Evaluator) builtinFloat(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		return object.NewFloat(arg.Node(), float64(arg.Value))
	case *object.Float:
		return arg
	case *object.String:
		float, err := strconv.ParseFloat(arg.Value, 64)

		if err != nil {
			return newError(node, "could not parse %q as FLOAT", arg.Value)
		}

		return object.NewFloat(arg.Node(), float)
	default:
		return newError(node, "argument to `float` not supported, got %s", args[0].Type())
	}
}

func (e *Evaluator) builtinString(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	return object.NewString(node, args[0].Inspect())
}

func (e *Evaluator) builtinBool(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Boolean:
		return arg
	case *object.Null:
		return FALSE
	case *object.Integer:
		return e.nativeBoolToBooleanObject(arg.Value != 0)
	case *object.Float:
		return e.nativeBoolToBooleanObject(arg.Value != 0)
	case *object.String:
		return e.nativeBoolToBooleanObject(arg.Value != "")
	case *object.Array:
		return e.nativeBoolToBooleanObject(len(arg.Elements) != 0)
	case *object.Hash:
		return e.nativeBoolToBooleanObject(len(arg.Pairs) != 0)
	default:
		return newError(node, "argument to `bool` not supported, got %s", args[0].Type())
	}
}

// String functions

func (e *Evaluator) builtinSplit(node ast.Node, args ...object.Object) object.Object {
	if len(args) == 0 || len(args) > 2 {
		return newError(node, "wrong number of arguments. got %d, want 1 or 2", len(args))
	}

	if args[0].Type() != object.STRING_OBJ {
		return newError(node, "first argument to `split` must be STRING, got %s", args[0].Type())
	}

	separator := " "

	if len(args) == 2 {
		if args[1].Type() != object.STRING_OBJ {
			return newError(node, "second argument to `split` must be STRING, got %s", args[1].Type())
		}

		separator = args[1].(*object.String).Value
	}

	str := args[0].(*object.String).Value

	parts := strings.Split(str, separator)
	elements := make([]object.Object, len(parts))

	for i, part := range parts {
		elements[i] = object.NewString(node, part)
	}

	return object.NewArray(node, elements)
}

func (e *Evaluator) builtinUppercase(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	if args[0].Type() != object.STRING_OBJ {
		return newError(node, "argument to `uppercase` must be STRING, got %s", args[0].Type())
	}

	str := args[0].(*object.String).Value

	return object.NewString(node, strings.ToUpper(str))
}

func (e *Evaluator) builtinLowercase(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	if args[0].Type() != object.STRING_OBJ {
		return newError(node, "argument to `lowercase` must be STRING, got %s", args[0].Type())
	}

	str := args[0].(*object.String).Value

	return object.NewString(node, strings.ToLower(str))
}

// String & Array functions

func (e *Evaluator) builtinLen(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Array:
		return object.NewInteger(arg.Node(), int64(len(arg.Elements)))
	case *object.String:
		return object.NewInteger(arg.Node(), int64(len(arg.Value)))
	default:
		return newError(node, "argument to `len` not supported, got %s", args[0].Type())
	}
}

func (e *Evaluator) builtinFirst(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Array:
		if len(arg.Elements) > 0 {
			return object.Clone(node, arg.Elements[0])
		}
	case *object.String:
		if len(arg.Value) > 0 {
			return object.NewString(node, string(arg.Value[0]))
		}
	default:
		return newError(node, "argument to `first` must be ARRAY or STRING, got %s", args[0].Type())
	}

	return NULL
}

func (e *Evaluator) builtinLast(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Array:
		length := len(arg.Elements)

		if length > 0 {
			return object.Clone(node, arg.Elements[length-1])
		}
	case *object.String:
		length := len(arg.Value)

		if length > 0 {
			return object.NewString(node, string(arg.Value[length-1]))
		}
	default:
		return newError(node, "argument to `last` must be ARRAY or STRING, got %s", args[0].Type())
	}

	return NULL
}

// Array functions

func (e *Evaluator) builtinRest(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError(node, "argument to `rest` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)

	if length > 0 {
		elements := make([]object.Object, length-1)

		for i := 1; i < length; i++ {
			elements[i-1] = object.Clone(node, arr.Elements[i])
		}

		return object.NewArray(node, elements)
	}

	return NULL
}

func (e *Evaluator) builtinPush(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError(node, "wrong number of arguments. got %d, want 2", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError(node, "first argument to `push` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	elements := make([]object.Object, length+1)

	for i := 0; i < length; i++ {
		elements[i] = object.Clone(node, arr.Elements[i])
	}

	elements[length] = object.Clone(node, args[1])

	return object.NewArray(node, elements)
}

func (e *Evaluator) builtinPop(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError(node, "first argument to `pop` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)

	if length > 0 {
		elements := make([]object.Object, length-1)

		for i := 0; i < length-1; i++ {
			elements[i] = object.Clone(node, arr.Elements[i])
		}

		return object.NewArray(node, elements)
	}

	return NULL
}

func (e *Evaluator) builtinIterMap(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 3 {
		return newError(node, "wrong number of arguments. got %d, want 2", len(args))
	}

	arr := args[0].(*object.Array)
	accumulated := args[1]
	f := args[2]

	if e.builtinLen(node, arr).Inspect() == "0" {
		return accumulated
	}

	return e.builtinIterMap(
		node,
		e.builtinRest(node, arr),
		e.builtinPush(node, accumulated, e.applyFunction(nil, f, []object.Object{e.builtinFirst(node, arr)})),
		f,
	)
}

func (e *Evaluator) builtinMap(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError(node, "wrong number of arguments. got %d, want 2", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError(node, "first argument to `map` must be ARRAY, got %s", args[0].Type())
	}

	if args[1].Type() != object.FUNCTION_OBJ && args[1].Type() != object.BUILTIN_OBJ {
		return newError(node, "second argument to `map` must be FUNCTION or BUILTIN, got %s", args[1].Type())
	}

	arr := args[0].(*object.Array)
	f := args[1]

	return e.builtinIterMap(node, arr, &object.Array{Elements: []object.Object{}}, f)
}

func (e *Evaluator) builtinIterReduce(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 3 {
		return newError(node, "wrong number of arguments. got %d, want 3", len(args))
	}

	arr := args[0].(*object.Array)
	result := args[1]
	f := args[2]

	if e.builtinLen(node, arr).Inspect() == "0" {
		return result
	}

	return e.builtinIterReduce(
		node,
		e.builtinRest(node, arr),
		e.applyFunction(nil, f, []object.Object{result, e.builtinFirst(node, arr)}),
		f,
	)
}

func (e *Evaluator) builtinReduce(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 3 {
		return newError(node, "wrong number of arguments. got %d, want 3", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError(node, "first argument to `reduce` must be ARRAY, got %s", args[0].Type())
	}

	if args[2].Type() != object.FUNCTION_OBJ && args[2].Type() != object.BUILTIN_OBJ {
		return newError(node, "third argument to `reduce` must be FUNCTION or BUILTIN, got %s", args[2].Type())
	}

	arr := args[0].(*object.Array)
	initial := args[1]
	f := args[2]

	return e.builtinIterReduce(node, arr, initial, f)
}

// IO functions

func (e *Evaluator) builtinPrint(node ast.Node, args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}

	return NULL
}

// Math functions

func (e *Evaluator) builtinAbs(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		if arg.Value < 0 {
			return object.NewInteger(arg.Node(), -arg.Value)
		}

		return arg
	case *object.Float:
		if arg.Value < 0 {
			return object.NewFloat(arg.Node(), -arg.Value)
		}

		return arg
	default:
		return newError(node, "argument to `abs` must be INTEGER or FLOAT, got %s", args[0].Type())
	}
}

func powFloat(node ast.Node, x float64, y float64) object.Object {
	pow := math.Pow(x, y)

	return object.NewFloat(node, pow)
}

func powInt(node ast.Node, x int64, y int64) object.Object {
	result := int64(1)

	if y < 0 {
		return powFloat(node, float64(x), float64(y))
	}

	for i := int64(0); i < y; i++ {
		result *= x
	}

	return object.NewInteger(node, result)
}

func (e *Evaluator) builtinPow(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError(node, "wrong number of arguments. got %d, want 2", len(args))
	}

	switch {
	case args[0].Type() == object.INTEGER_OBJ && args[1].Type() == object.INTEGER_OBJ:
		x := args[0].(*object.Integer).Value
		y := args[1].(*object.Integer).Value

		return powInt(node, x, y)

	case args[0].Type() == object.FLOAT_OBJ && args[1].Type() == object.FLOAT_OBJ:
		x := args[0].(*object.Float).Value
		y := args[1].(*object.Float).Value

		return powFloat(node, x, y)

	case args[0].Type() == object.INTEGER_OBJ && args[1].Type() == object.FLOAT_OBJ:
		x := args[0].(*object.Integer).Value
		y := args[1].(*object.Float).Value

		return powFloat(node, float64(x), y)

	case args[0].Type() == object.FLOAT_OBJ && args[1].Type() == object.INTEGER_OBJ:
		x := args[0].(*object.Float).Value
		y := args[1].(*object.Integer).Value

		return powFloat(node, x, float64(y))

	default:
		return newError(node, "arguments to `pow` must be INTEGER or FLOAT, got %s and %s", args[0].Type(), args[1].Type())
	}
}

func (e *Evaluator) builtinSqrt(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	if !object.IsNumber(args[0]) {
		return newError(node, "argument to `sqrt` must be INTEGER or FLOAT, got %s", args[0].Type())
	}

	var x float64

	if args[0].Type() == object.FLOAT_OBJ {
		x = args[0].(*object.Float).Value
	} else if args[0].Type() == object.INTEGER_OBJ {
		x = float64(args[0].(*object.Integer).Value)
	}

	if x < 0 {
		return newError(node, "argument to `sqrt` must be non-negative, got %g", x)
	}

	sqrt := math.Sqrt(x)

	return &object.Float{Value: sqrt}
}
