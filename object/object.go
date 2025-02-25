package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/iskandervdh/vorn/ast"
)

type ObjectType string

const (
	NULL_OBJ         = "NULL"
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	FLOAT_OBJ        = "FLOAT"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	BREAK_OBJ        = "BREAK"
	CONTINUE_OBJ     = "CONTINUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
)

type Object interface {
	Type() ObjectType
	Inspect() string
	Node() ast.Node
}

func IsNumber(obj Object) bool {
	return obj.Type() == INTEGER_OBJ || obj.Type() == FLOAT_OBJ
}

type Null struct {
	node ast.Node
}

func NewNull(node ast.Node) *Null {
	return &Null{node: node}
}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }
func (n *Null) Node() ast.Node   { return n.node }

type Integer struct {
	node  ast.Node
	Value int64
}

func NewInteger(node ast.Node, value int64) *Integer {
	return &Integer{node: node, Value: value}
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Node() ast.Node   { return i.node }

type Boolean struct {
	node  ast.Node
	Value bool
}

func NewBoolean(node ast.Node, value bool) *Boolean {
	return &Boolean{node: node, Value: value}
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Node() ast.Node   { return b.node }

type Float struct {
	node  ast.Node
	Value float64
}

func NewFloat(node ast.Node, value float64) *Float {
	return &Float{node: node, Value: value}
}

func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) Inspect() string  { return fmt.Sprintf("%g", f.Value) }
func (f *Float) Node() ast.Node   { return f.node }

type ReturnValue struct {
	node  ast.Node
	Value Object
}

func NewReturnValue(node ast.Node, value Object) *ReturnValue {
	return &ReturnValue{node: node, Value: value}
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Node() ast.Node   { return rv.node }

type Break struct {
	node ast.Node
}

func NewBreak(node ast.Node) *Break {
	return &Break{node: node}
}

func (b *Break) Type() ObjectType { return BREAK_OBJ }
func (b *Break) Inspect() string  { return "break" }
func (b *Break) Node() ast.Node   { return b.node }

type Continue struct {
	node ast.Node
}

func NewContinue(node ast.Node) *Continue {
	return &Continue{node: node}
}

func (c *Continue) Type() ObjectType { return CONTINUE_OBJ }
func (c *Continue) Inspect() string  { return "continue" }
func (c *Continue) Node() ast.Node   { return c.node }

type Error struct {
	node    ast.Node
	Message string
}

func NewError(node ast.Node, format string, a ...interface{}) *Error {
	location := fmt.Sprintf("[%d:%d]", node.Line(), node.Column())
	message := location + " " + fmt.Sprintf(format, a...)

	return &Error{node: node, Message: message}
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
func (e *Error) Node() ast.Node   { return e.node }

type Function struct {
	node      ast.Node
	Arguments []*ast.Identifier
	Body      *ast.BlockStatement
	Env       *Environment
}

func NewFunction(node ast.Node, args []*ast.Identifier, body *ast.BlockStatement, env *Environment) *Function {
	return &Function{node: node, Arguments: args, Body: body, Env: env}
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	args := []string{}

	for _, p := range f.Arguments {
		args = append(args, p.String())
	}

	out.WriteString("func")
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(") ")
	out.WriteString(f.Body.String())
	out.WriteString("\n")

	return out.String()
}
func (f *Function) Node() ast.Node { return f.node }

type String struct {
	node  ast.Node
	Value string
}

func NewString(node ast.Node, value string) *String {
	return &String{node: node, Value: value}
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) Node() ast.Node   { return s.node }

type BuiltinFunction func(node ast.Node, args ...Object) Object

type Builtin struct {
	node           ast.Node
	Function       BuiltinFunction
	ArgumentsCount int
}

func NewBuiltin(node ast.Node, function BuiltinFunction) *Builtin {
	return &Builtin{node: node, Function: function}
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) Node() ast.Node   { return b.node }

type Array struct {
	node     ast.Node
	Elements []Object
}

func NewArray(node ast.Node, elements []Object) *Array {
	return &Array{node: node, Elements: elements}
}

func (arr *Array) Type() ObjectType { return ARRAY_OBJ }
func (arr *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}

	for _, e := range arr.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
func (arr *Array) Node() ast.Node { return arr.node }

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	node  ast.Node
	Pairs map[HashKey]HashPair
}

func NewHash(node ast.Node, pairs map[HashKey]HashPair) *Hash {
	return &Hash{node: node, Pairs: pairs}
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}

	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
func (h *Hash) Node() ast.Node { return h.node }

type Hashable interface {
	HashKey() HashKey
}

func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

func Clone(node ast.Node, object Object) Object {
	switch object.Type() {
	case NULL_OBJ:
		return NewNull(node)
	case INTEGER_OBJ:
		return NewInteger(node, object.(*Integer).Value)
	case BOOLEAN_OBJ:
		return NewBoolean(node, object.(*Boolean).Value)
	case FLOAT_OBJ:
		return NewFloat(node, object.(*Float).Value)
	case RETURN_VALUE_OBJ:
		return NewReturnValue(node, object)
	case BREAK_OBJ:
		return NewBreak(node)
	case CONTINUE_OBJ:
		return NewContinue(node)
	case ERROR_OBJ:
		return &Error{node: node, Message: object.(*Error).Message}
	case FUNCTION_OBJ:
		return NewFunction(node, object.(*Function).Arguments, object.(*Function).Body, object.(*Function).Env)
	case STRING_OBJ:
		return NewString(node, object.(*String).Value)
	case BUILTIN_OBJ:
		return NewBuiltin(node, object.(*Builtin).Function)
	case ARRAY_OBJ:
		elements := make([]Object, len(object.(*Array).Elements))

		for i, el := range object.(*Array).Elements {
			elements[i] = Clone(node, el)
		}

		return NewArray(node, elements)
	case HASH_OBJ:
		pairs := make(map[HashKey]HashPair)

		for key, pair := range object.(*Hash).Pairs {
			pairCopy := HashPair{Key: Clone(node, pair.Key), Value: Clone(node, pair.Value)}

			pairs[key] = pairCopy
		}

		return NewHash(node, object.(*Hash).Pairs)
	default:
		return NewNull(node)
	}
}
