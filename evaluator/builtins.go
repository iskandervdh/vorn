package evaluator

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/iskandervdh/vorn/object"
)

// Common functions

func (e *Evaluator) builtinType(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got %d, want 1", len(args))
	}

	typeName := args[0].Type()

	return &object.String{Value: string(typeName)}
}

func (e *Evaluator) builtinRange(args ...object.Object) object.Object {
	if len(args) < 1 || len(args) > 2 {
		return newError("wrong number of arguments. got %d, want 1 or 2", len(args))
	}

	if args[0].Type() != object.INTEGER_OBJ {
		return newError("first argument to `range` must be INTEGER, got %s", args[0].Type())
	}

	start := args[0].(*object.Integer).Value
	end := start

	if len(args) == 2 {
		if args[1].Type() != object.INTEGER_OBJ {
			return newError("second argument to `range` must be INTEGER, got %s", args[1].Type())
		}

		end = args[1].(*object.Integer).Value
	} else {
		start = 0
	}

	elementsLength := end - start

	if start > end {
		elementsLength = start - end
	}

	elements := make([]object.Object, elementsLength)

	if start > end {
		for i := start; i > end; i-- {
			elements[start-i] = &object.Integer{Value: i}
		}
	} else {
		for i := start; i < end; i++ {
			elements[i-start] = &object.Integer{Value: i}
		}
	}

	return &object.Array{Elements: elements}
}

// Conversion functions

func (e *Evaluator) builtinInt(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got %d, want 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		return arg
	case *object.Float:
		return &object.Integer{Value: int64(arg.Value)}
	case *object.String:
		integer, err := strconv.ParseInt(arg.Value, 0, 64)

		if err != nil {
			return newError("could not parse %q as INTEGER", arg.Value)
		}

		return &object.Integer{Value: integer}
	default:
		return newError("argument to `int` not supported, got %s", args[0].Type())
	}
}

func (e *Evaluator) builtinFloat(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got %d, want 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		return &object.Float{Value: float64(arg.Value)}
	case *object.Float:
		return arg
	case *object.String:
		float, err := strconv.ParseFloat(arg.Value, 64)

		if err != nil {
			return newError("could not parse %q as FLOAT", arg.Value)
		}

		return &object.Float{Value: float}
	default:
		return newError("argument to `float` not supported, got %s", args[0].Type())
	}
}

func (e *Evaluator) builtinString(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got %d, want 1", len(args))
	}

	return &object.String{Value: args[0].Inspect()}
}

func (e *Evaluator) builtinBool(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got %d, want 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Boolean:
		return arg
	case *object.Integer:
		return &object.Boolean{Value: arg.Value != 0}
	case *object.Float:
		return &object.Boolean{Value: arg.Value != 0}
	case *object.String:
		return &object.Boolean{Value: arg.Value != ""}
	case *object.Array:
		return &object.Boolean{Value: len(arg.Elements) != 0}
	default:
		return newError("argument to `bool` not supported, got %s", args[0].Type())
	}
}

// String functions

func (e *Evaluator) builtinSplit(args ...object.Object) object.Object {
	if len(args) == 0 || len(args) > 2 {
		return newError("wrong number of arguments. got %d, want 1 or 2", len(args))
	}

	if args[0].Type() != object.STRING_OBJ {
		return newError("first argument to `split` must be STRING, got %s", args[0].Type())
	}

	separator := " "

	if len(args) == 2 {
		if args[1].Type() != object.STRING_OBJ {
			return newError("second argument to `split` must be STRING, got %s", args[1].Type())
		}

		separator = args[1].(*object.String).Value
	}

	str := args[0].(*object.String).Value

	parts := strings.Split(str, separator)
	elements := make([]object.Object, len(parts))

	for i, part := range parts {
		elements[i] = &object.String{Value: part}
	}
	return &object.Array{Elements: elements}
}

func (e *Evaluator) builtinUppercase(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got %d, want 1", len(args))
	}

	if args[0].Type() != object.STRING_OBJ {
		return newError("argument to `uppercase` must be STRING, got %s", args[0].Type())
	}

	str := args[0].(*object.String).Value

	return &object.String{Value: strings.ToUpper(str)}
}

func (e *Evaluator) builtinLowercase(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got %d, want 1", len(args))
	}

	if args[0].Type() != object.STRING_OBJ {
		return newError("argument to `lowercase` must be STRING, got %s", args[0].Type())
	}

	str := args[0].(*object.String).Value

	return &object.String{Value: strings.ToLower(str)}
}

// String & Array functions

func (e *Evaluator) builtinLen(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got %d, want 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	default:
		return newError("argument to `len` not supported, got %s",
			args[0].Type())
	}
}

func (e *Evaluator) builtinFirst(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got %d, want 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Array:
		if len(arg.Elements) > 0 {
			return arg.Elements[0]
		}
	case *object.String:
		if len(arg.Value) > 0 {
			return &object.String{Value: string(arg.Value[0])}
		}
	default:
		return newError("argument to `first` must be ARRAY or STRING, got %s", args[0].Type())
	}

	return NULL
}

func (e *Evaluator) builtinLast(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got %d, want 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Array:
		length := len(arg.Elements)

		if length > 0 {
			return arg.Elements[length-1]
		}
	case *object.String:
		length := len(arg.Value)

		if length > 0 {
			return &object.String{Value: string(arg.Value[length-1])}
		}
	default:
		return newError("argument to `last` must be ARRAY or STRING, got %s", args[0].Type())
	}

	return NULL
}

// Array functions

func (e *Evaluator) builtinRest(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got %d, want 1", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)

	if length > 0 {
		newElements := make([]object.Object, length-1)
		copy(newElements, arr.Elements[1:length])

		return &object.Array{Elements: newElements}
	}

	return NULL
}

func (e *Evaluator) builtinPush(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got %d, want 2", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("first argument to `push` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	newElements := make([]object.Object, length+1)

	copy(newElements, arr.Elements)
	newElements[length] = args[1]

	return &object.Array{Elements: newElements}
}

func (e *Evaluator) builtinPop(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got %d, want 1", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("first argument to `pop` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)

	if length > 0 {
		newElements := make([]object.Object, length-1)
		copy(newElements, arr.Elements[:length-1])

		return &object.Array{Elements: newElements}
	}

	return NULL
}

func (e *Evaluator) builtinIterMap(args ...object.Object) object.Object {
	if len(args) != 3 {
		return newError("wrong number of arguments. got %d, want 2", len(args))
	}

	arr := args[0].(*object.Array)
	accumulated := args[1]
	f := args[2]

	if e.builtinLen(arr).Inspect() == "0" {
		return accumulated
	}

	return e.builtinIterMap(e.builtinRest(arr), e.builtinPush(accumulated, e.applyFunction(f, []object.Object{e.builtinFirst(arr)})), f)
}

func (e *Evaluator) builtinMap(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got %d, want 2", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("first argument to `map` must be ARRAY, got %s", args[0].Type())
	}

	if args[1].Type() != object.FUNCTION_OBJ && args[1].Type() != object.BUILTIN_OBJ {
		return newError("second argument to `map` must be FUNCTION or BUILTIN, got %s", args[1].Type())
	}

	arr := args[0].(*object.Array)
	f := args[1]

	return e.builtinIterMap(arr, &object.Array{Elements: []object.Object{}}, f)
}

func (e *Evaluator) builtinIterReduce(args ...object.Object) object.Object {
	if len(args) != 3 {
		return newError("wrong number of arguments. got %d, want 3", len(args))
	}

	arr := args[0].(*object.Array)
	result := args[1]
	f := args[2]

	if e.builtinLen(arr).Inspect() == "0" {
		return result
	}

	return e.builtinIterReduce(e.builtinRest(arr), e.applyFunction(f, []object.Object{result, e.builtinFirst(arr)}), f)
}

func (e *Evaluator) builtinReduce(args ...object.Object) object.Object {
	if len(args) != 3 {
		return newError("wrong number of arguments. got %d, want 3", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("first argument to `reduce` must be ARRAY, got %s", args[0].Type())
	}

	if args[2].Type() != object.FUNCTION_OBJ && args[2].Type() != object.BUILTIN_OBJ {
		return newError("third argument to `reduce` must be FUNCTION or BUILTIN, got %s", args[2].Type())
	}

	arr := args[0].(*object.Array)
	initial := args[1]
	f := args[2]

	return e.builtinIterReduce(arr, initial, f)
}

// IO functions

func (e *Evaluator) builtinPrint(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}

	return NULL
}

// Math functions

func (e *Evaluator) builtinAbs(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got %d, want 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		if arg.Value < 0 {
			return &object.Integer{Value: -arg.Value}
		}

		return arg
	case *object.Float:
		if arg.Value < 0 {
			return &object.Float{Value: -arg.Value}
		}

		return arg
	default:
		return newError("argument to `abs` must be INTEGER or FLOAT, got %s", args[0].Type())
	}
}

func powFloat(x float64, y float64) object.Object {
	pow := math.Pow(x, y)

	return &object.Float{Value: pow}
}

func powInt(x int64, y int64) object.Object {
	result := int64(1)

	if y < 0 {
		return powFloat(float64(x), float64(y))
	}

	for i := int64(0); i < y; i++ {
		result *= x
	}

	return &object.Integer{Value: result}
}

func (e *Evaluator) builtinPow(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got %d, want 2", len(args))
	}

	switch {
	case args[0].Type() == object.INTEGER_OBJ && args[1].Type() == object.INTEGER_OBJ:
		x := args[0].(*object.Integer).Value
		y := args[1].(*object.Integer).Value

		return powInt(x, y)

	case args[0].Type() == object.FLOAT_OBJ && args[1].Type() == object.FLOAT_OBJ:
		x := args[0].(*object.Float).Value
		y := args[1].(*object.Float).Value

		return powFloat(x, y)

	case args[0].Type() == object.INTEGER_OBJ && args[1].Type() == object.FLOAT_OBJ:
		x := args[0].(*object.Integer).Value
		y := args[1].(*object.Float).Value

		return powFloat(float64(x), y)

	case args[0].Type() == object.FLOAT_OBJ && args[1].Type() == object.INTEGER_OBJ:
		x := args[0].(*object.Float).Value
		y := args[1].(*object.Integer).Value

		return powFloat(x, float64(y))

	default:
		return newError("arguments to `pow` must be INTEGER or FLOAT, got %s and %s", args[0].Type(), args[1].Type())
	}
}

func (e *Evaluator) builtinSqrt(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got %d, want 1", len(args))
	}

	if !object.IsNumber(args[0]) {
		return newError("argument to `sqrt` must be INTEGER or FLOAT, got %s", args[0].Type())
	}

	var x float64

	if args[0].Type() == object.FLOAT_OBJ {
		x = args[0].(*object.Float).Value
	} else if args[0].Type() == object.INTEGER_OBJ {
		x = float64(args[0].(*object.Integer).Value)
	}

	if x < 0 {
		return newError("argument to `sqrt` must be non-negative, got %g", x)
	}

	sqrt := math.Sqrt(x)

	return &object.Float{Value: sqrt}
}
