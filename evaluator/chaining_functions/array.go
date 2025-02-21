package chaining_functions

import "github.com/iskandervdh/vorn/object"

type ArrayChainingFunction func(left *object.Array, args ...object.Object) object.Object

var ArrayChainingFunctions = map[string]ArrayChainingFunction{
	"length": arrayLength,
}

func arrayLength(left *object.Array, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(left.Node(), "Array.length() takes no arguments")
	}

	return object.NewInteger(left.Node(), int64(len(left.Elements)))
}
