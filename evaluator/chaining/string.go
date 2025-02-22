package chaining

import (
	"strings"

	"github.com/iskandervdh/vorn/object"
)

type StringChainingFunction func(left *object.String, args ...object.Object) object.Object

var StringChainingFunctions = map[string]StringChainingFunction{
	"length": stringLength,
	"upper":  stringUpper,
	"lower":  stringLower,
	"split":  stringSplit,
}

func stringLength(left *object.String, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(left.Node(), "String.length() takes no arguments")
	}

	return object.NewInteger(left.Node(), int64(len(left.Value)))
}

func stringUpper(left *object.String, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(left.Node(), "String.upper() takes no arguments")
	}

	return object.NewString(left.Node(), strings.ToUpper(left.Value))
}

func stringLower(left *object.String, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(left.Node(), "String.lower() takes no arguments")
	}

	return object.NewString(left.Node(), strings.ToLower(left.Value))
}

func stringSplit(left *object.String, args ...object.Object) object.Object {
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
