package chaining

import "github.com/iskandervdh/vorn/object"

type ArrayChainingFunction func(left *object.Array, args ...object.Object) object.Object

var ArrayChainingFunctions = map[string]ArrayChainingFunction{
	"length": arrayLength,
	"push":   arrayPush,
	"pop":    arrayPop,
}

func arrayLength(left *object.Array, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(left.Node(), "Array.length() takes no arguments")
	}

	return object.NewInteger(left.Node(), int64(len(left.Elements)))
}

func arrayPush(left *object.Array, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(left.Node(), "Array.push() takes exactly 1 argument")
	}

	left.Elements = append(left.Elements, args[0])

	return left
}

func arrayPop(left *object.Array, args ...object.Object) object.Object {
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
