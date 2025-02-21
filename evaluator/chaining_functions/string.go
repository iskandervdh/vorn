package chaining_functions

import (
	"strings"

	"github.com/iskandervdh/vorn/object"
)

type StringChainingFunction func(left *object.String, args ...object.Object) object.Object

var StringChainingFunctions = map[string]StringChainingFunction{
	"length": stringLength,
	"upper":  stringUpper,
	"lower":  stringLower,
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
