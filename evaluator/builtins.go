package evaluator

import (
	"fmt"
	"math"
	"strconv"

	"github.com/iskandervdh/vorn/ast"
	"github.com/iskandervdh/vorn/object"
)

// Common functions

func (e *Evaluator) builtinType(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	typeName := args[0].Type()

	return object.NewString(node, string(typeName))
}

func (e *Evaluator) builtinRange(node ast.Node, args ...object.Object) object.Object {
	if len(args) < 1 || len(args) > 2 {
		return object.NewError(node, "wrong number of arguments. got %d, want 1 or 2", len(args))
	}

	if args[0].Type() != object.INTEGER_OBJ {
		return object.NewError(node, "first argument to `range` must be INTEGER, got %s", args[0].Type())
	}

	firstArg := args[0].(*object.Integer)

	start := firstArg.Value
	end := start

	if len(args) == 1 {
		if start < 0 {
			return object.NewError(node, "argument to `range` must be non-negative, got %d", start)
		}

		start = 0
	} else {
		if args[1].Type() != object.INTEGER_OBJ {
			return object.NewError(node, "second argument to `range` must be INTEGER, got %s", args[1].Type())
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
		return object.NewError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		return arg
	case *object.Float:
		return object.NewInteger(arg.Node(), int64(arg.Value))
	case *object.String:
		integer, err := strconv.ParseInt(arg.Value, 0, 64)

		if err != nil {
			return object.NewError(node, "could not parse %q as INTEGER", arg.Value)
		}

		return object.NewInteger(arg.Node(), integer)
	default:
		return object.NewError(node, "argument to `int` not supported, got %s", args[0].Type())
	}
}

func (e *Evaluator) builtinFloat(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		return object.NewFloat(arg.Node(), float64(arg.Value))
	case *object.Float:
		return arg
	case *object.String:
		float, err := strconv.ParseFloat(arg.Value, 64)

		if err != nil {
			return object.NewError(node, "could not parse %q as FLOAT", arg.Value)
		}

		return object.NewFloat(arg.Node(), float)
	default:
		return object.NewError(node, "argument to `float` not supported, got %s", args[0].Type())
	}
}

func (e *Evaluator) builtinString(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	return object.NewString(node, args[0].Inspect())
}

func (e *Evaluator) builtinBool(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(node, "wrong number of arguments. got %d, want 1", len(args))
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
		return object.NewError(node, "argument to `bool` not supported, got %s", args[0].Type())
	}
}

// String & Array functions

func (e *Evaluator) builtinLen(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Array:
		return object.NewInteger(arg.Node(), int64(len(arg.Elements)))
	case *object.String:
		return object.NewInteger(arg.Node(), int64(len(arg.Value)))
	default:
		return object.NewError(node, "argument to `len` not supported, got %s", args[0].Type())
	}
}

func (e *Evaluator) builtinFirst(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(node, "wrong number of arguments. got %d, want 1", len(args))
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
		return object.NewError(node, "argument to `first` must be ARRAY or STRING, got %s", args[0].Type())
	}

	return NULL
}

func (e *Evaluator) builtinLast(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(node, "wrong number of arguments. got %d, want 1", len(args))
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
		return object.NewError(node, "argument to `last` must be ARRAY or STRING, got %s", args[0].Type())
	}

	return NULL
}

// Array functions

func (e *Evaluator) builtinRest(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return object.NewError(node, "argument to `rest` must be ARRAY, got %s", args[0].Type())
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
		return object.NewError(node, "wrong number of arguments. got %d, want 1", len(args))
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
		return object.NewError(node, "argument to `abs` must be INTEGER or FLOAT, got %s", args[0].Type())
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
		return object.NewError(node, "wrong number of arguments. got %d, want 2", len(args))
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
		return object.NewError(node, "arguments to `pow` must be INTEGER or FLOAT, got %s and %s", args[0].Type(), args[1].Type())
	}
}

func (e *Evaluator) builtinSqrt(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	if !object.IsNumber(args[0]) {
		return object.NewError(node, "argument to `sqrt` must be INTEGER or FLOAT, got %s", args[0].Type())
	}

	var x float64

	if args[0].Type() == object.FLOAT_OBJ {
		x = args[0].(*object.Float).Value
	} else if args[0].Type() == object.INTEGER_OBJ {
		x = float64(args[0].(*object.Integer).Value)
	}

	if x < 0 {
		return object.NewError(node, "argument to `sqrt` must be non-negative, got %g", x)
	}

	sqrt := math.Sqrt(x)

	return &object.Float{Value: sqrt}
}

func (e *Evaluator) builtinSin(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	if !object.IsNumber(args[0]) {
		return object.NewError(node, "argument to `sin` must be INTEGER or FLOAT, got %s", args[0].Type())
	}

	var x float64

	if args[0].Type() == object.FLOAT_OBJ {
		x = args[0].(*object.Float).Value
	} else if args[0].Type() == object.INTEGER_OBJ {
		x = float64(args[0].(*object.Integer).Value)
	}

	sin := math.Sin(x)

	return &object.Float{Value: sin}
}

func (e *Evaluator) builtinCos(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	if !object.IsNumber(args[0]) {
		return object.NewError(node, "argument to `cos` must be INTEGER or FLOAT, got %s", args[0].Type())
	}

	var x float64

	if args[0].Type() == object.FLOAT_OBJ {
		x = args[0].(*object.Float).Value
	} else if args[0].Type() == object.INTEGER_OBJ {
		x = float64(args[0].(*object.Integer).Value)
	}

	cos := math.Cos(x)

	return &object.Float{Value: cos}
}

func (e *Evaluator) builtinTan(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	if !object.IsNumber(args[0]) {
		return object.NewError(node, "argument to `tan` must be INTEGER or FLOAT, got %s", args[0].Type())
	}

	var x float64

	if args[0].Type() == object.FLOAT_OBJ {
		x = args[0].(*object.Float).Value
	}

	if args[0].Type() == object.INTEGER_OBJ {
		x = float64(args[0].(*object.Integer).Value)
	}

	tan := math.Tan(x)

	return &object.Float{Value: tan}
}

func (e *Evaluator) builtinSum(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return object.NewError(node, "argument to `sum` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)

	var sum float64

	for _, element := range arr.Elements {
		if !object.IsNumber(element) {
			return object.NewError(node, "elements in array must be INTEGER or FLOAT, got %s", element.Type())
		}

		if element.Type() == object.FLOAT_OBJ {
			sum += element.(*object.Float).Value
		} else if element.Type() == object.INTEGER_OBJ {
			sum += float64(element.(*object.Integer).Value)
		}
	}

	// Check if the sum is an integer
	if sum == float64(int64(sum)) {
		return &object.Integer{Value: int64(sum)}
	}

	return &object.Float{Value: sum}
}

func (e *Evaluator) builtinMean(node ast.Node, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(node, "wrong number of arguments. got %d, want 1", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return object.NewError(node, "argument to `mean` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)

	if len(arr.Elements) == 0 {
		return &object.Integer{Value: 0}
	}

	var sum float64

	for _, element := range arr.Elements {
		if !object.IsNumber(element) {
			return object.NewError(node, "elements in array must be INTEGER or FLOAT, got %s", element.Type())
		}

		if element.Type() == object.FLOAT_OBJ {
			sum += element.(*object.Float).Value
		} else if element.Type() == object.INTEGER_OBJ {
			sum += float64(element.(*object.Integer).Value)
		}
	}

	mean := sum / float64(len(arr.Elements))

	// Check if the mean is an integer
	if mean == float64(int64(mean)) {
		return &object.Integer{Value: int64(mean)}
	}

	return &object.Float{Value: mean}
}
