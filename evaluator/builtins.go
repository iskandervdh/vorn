package evaluator

import (
	"fmt"
	"math"

	"github.com/iskandervdh/vorn/object"
)

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

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)

	if len(arr.Elements) > 0 {
		return arr.Elements[0]
	}

	return NULL
}

func (e *Evaluator) builtinLast(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got %d, want 1", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `last` must be ARRAY, got '%s'", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)

	if length > 0 {
		return arr.Elements[length-1]
	}

	return NULL
}

func (e *Evaluator) builtinRest(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got %d, want 1", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `rest` must be ARRAY, got %s",
			args[0].Type())
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
		return newError("argument to `push` must be ARRAY, got '%s'", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	newElements := make([]object.Object, length+1)

	copy(newElements, arr.Elements)
	newElements[length] = args[1]

	return &object.Array{Elements: newElements}
}

func (e *Evaluator) builtinIterMap(args ...object.Object) object.Object {
	if len(args) != 3 {
		return newError("wrong number of arguments. got %d, want 2", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("first argument to `iter` must be ARRAY, got %s", args[0].Type())
	}

	if args[2].Type() != object.FUNCTION_OBJ {
		return newError("third argument to `iter` must be FUNCTION, got %s", args[2].Type())
	}

	arr := args[0].(*object.Array)
	accumulated := args[1]
	f := args[2].(*object.Function)

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
		return newError("first argument to `map` must be ARRAY, got '%s'", args[0].Type())
	}

	if args[1].Type() != object.FUNCTION_OBJ {
		return newError("second argument to `map` must be FUNCTION, got '%s'", args[1].Type())
	}

	arr := args[0].(*object.Array)
	f := args[1].(*object.Function)

	return e.builtinIterMap(arr, &object.Array{Elements: []object.Object{}}, f)
}

func (e *Evaluator) builtinIterReduce(args ...object.Object) object.Object {
	if len(args) != 3 {
		return newError("wrong number of arguments. got %d, want 3", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("first argument to `iter` must be ARRAY, got '%s'", args[0].Type())
	}

	if args[2].Type() != object.FUNCTION_OBJ {
		return newError("third argument to `iter` must be FUNCTION, got '%s'", args[2].Type())
	}

	arr := args[0].(*object.Array)
	result := args[1]
	f := args[2].(*object.Function)

	if e.builtinLen(arr).Inspect() == "0" {
		return result
	}

	return e.builtinIterReduce(e.builtinRest(arr), e.applyFunction(f, []object.Object{result, e.builtinFirst(arr)}), f)
}

func (e *Evaluator) builtinReduce(args ...object.Object) object.Object {
	if len(args) != 3 {
		return newError("wrong number of arguments. got %d, want 2", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("first argument to `reduce` must be ARRAY, got '%s'", args[0].Type())
	}

	if args[2].Type() != object.FUNCTION_OBJ {
		return newError("third argument to `reduce` must be FUNCTION, got '%s'", args[1].Type())
	}

	arr := args[0].(*object.Array)
	initial := args[1]
	f := args[2].(*object.Function)

	return e.builtinIterReduce(arr, initial, f)
}

func (e *Evaluator) builtinPrint(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}

	return NULL
}

func powFloat(x float64, y float64) object.Object {
	if x < 0 {
		return newError("first argument to `pow` must be non-negative, got %f", x)
	}

	pow := math.Pow(x, y)

	return &object.Float{Value: pow}
}

func powInt(x int64, y int64) object.Object {
	result := int64(1)

	if x < 0 {
		return newError("first argument to `pow` must be non-negative, got %d", x)
	}

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

	if !object.IsNumber(args[0]) {
		return newError("first argument to `pow` must be INTEGER OR FLOAT, got %s", args[0].Type())
	}

	if !object.IsNumber(args[1]) {
		return newError("second argument to `pow` must be INTEGER OR FLOAT, got %s", args[1].Type())
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
		return newError("arguments to `pow` must be INTEGER OR FLOAT, got %s and %s", args[0].Type(), args[1].Type())
	}
}

func (e *Evaluator) builtinSqrt(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got %d, want 1", len(args))
	}

	if !object.IsNumber(args[0]) {
		return newError("argument to `sqrt` must be INTEGER OR FLOAT, got %s", args[0].Type())
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
