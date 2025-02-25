package object

import (
	"testing"

	"github.com/iskandervdh/vorn/ast"
	"github.com/iskandervdh/vorn/token"
)

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func TestHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is Jeff"}
	diff2 := &String{Value: "My name is Jeff"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have same keys")
	}

	true1 := &Boolean{Value: true}
	true2 := &Boolean{Value: true}
	false1 := &Boolean{Value: false}
	false2 := &Boolean{Value: false}

	if true1.HashKey() != true2.HashKey() {
		t.Errorf("booleans with same content have different keys")
	}

	if false1.HashKey() != false2.HashKey() {
		t.Errorf("booleans with same content have different keys")
	}

	if true1.HashKey() == false1.HashKey() {
		t.Errorf("booleans with different content have same keys")
	}

	integer1 := &Integer{Value: 1}
	integer2 := &Integer{Value: 1}
	integer3 := &Integer{Value: 2}
	integer4 := &Integer{Value: 2}

	if integer1.HashKey() != integer2.HashKey() {
		t.Errorf("integers with same content have different keys")
	}

	if integer3.HashKey() != integer4.HashKey() {
		t.Errorf("integers with same content have different keys")
	}

	if integer1.HashKey() == integer3.HashKey() {
		t.Errorf("integers with different content have same keys")
	}
}

type Nonexistent struct {
	node ast.Node
}

func (n *Nonexistent) Type() ObjectType { return "" }
func (n *Nonexistent) Inspect() string  { return "" }
func (n *Nonexistent) Node() ast.Node   { return n.node }

func TestObjectClone(t *testing.T) {
	testCases := []struct {
		name           string
		node           ast.Node
		object         Object
		modifyFunction func(Object)
	}{
		{
			name: "Boolean",
			node: &ast.BooleanLiteral{
				Value: true,
				Token: token.Token{
					Type:    token.TRUE,
					Literal: "true",
					Line:    1,
					Column:  1,
				},
			},
			object: newBoolean(&ast.BooleanLiteral{
				Value: true,
				Token: token.Token{
					Type:    token.TRUE,
					Literal: "true",
					Line:    1,
					Column:  1,
				},
			}, true),
			modifyFunction: func(obj Object) {
				obj.(*Boolean).Value = false
			},
		},
		{
			name: "String",
			node: &ast.StringLiteral{
				Token: token.Token{
					Type:    token.STRING,
					Literal: "hello",
					Line:    1,
					Column:  1,
				},
				Value: "hello",
			},
			object: NewString(&ast.StringLiteral{
				Token: token.Token{
					Type:    token.STRING,
					Literal: "hello",
					Line:    1,
					Column:  1,
				},
				Value: "hello",
			}, "hello"),
			modifyFunction: func(obj Object) {
				obj.(*String).Value = "world"
			},
		},
		{
			name: "Integer",
			node: &ast.IntegerLiteral{
				Token: token.Token{
					Type:    token.INT,
					Literal: "1",
					Line:    1,
					Column:  1,
				},
				Value: 1,
			},
			object: NewInteger(&ast.IntegerLiteral{
				Token: token.Token{
					Type:    token.INT,
					Literal: "1",
					Line:    1,
					Column:  1,
				},
				Value: 1,
			}, 1),
			modifyFunction: func(obj Object) {
				obj.(*Integer).Value = 2
			},
		},
		{
			name: "Float",
			node: &ast.FloatLiteral{
				Token: token.Token{
					Type:    token.FLOAT,
					Literal: "1.5",
					Line:    1,
					Column:  1,
				},
				Value: 1.5,
			},
			object: NewFloat(&ast.FloatLiteral{
				Token: token.Token{
					Type:    token.FLOAT,
					Literal: "1.5",
					Line:    1,
					Column:  1,
				},
				Value: 1.5,
			}, 1.5),
			modifyFunction: func(obj Object) {
				obj.(*Float).Value = 2.5
			},
		},
		{
			name: "Return",
			node: &ast.ReturnStatement{
				Token: token.Token{
					Type:    token.RETURN,
					Literal: "return",
					Line:    1,
					Column:  1,
				},
				ReturnValue: &ast.IntegerLiteral{
					Token: token.Token{
						Type:    token.INT,
						Literal: "1",
						Line:    1,
						Column:  1,
					},
					Value: 1,
				},
			},
			object: NewReturnValue(&ast.ReturnStatement{
				Token: token.Token{
					Type:    token.RETURN,
					Literal: "return",
					Line:    1,
					Column:  1,
				},
				ReturnValue: &ast.IntegerLiteral{
					Token: token.Token{
						Type:    token.INT,
						Literal: "1",
						Line:    1,
						Column:  1,
					},
					Value: 1,
				},
			}, NewInteger(&ast.IntegerLiteral{
				Token: token.Token{
					Type:    token.INT,
					Literal: "1",
					Line:    1,
					Column:  1,
				},
				Value: 1,
			}, 1)),
			modifyFunction: func(obj Object) {
				obj.(*ReturnValue).Value = NewInteger(&ast.IntegerLiteral{
					Token: token.Token{
						Type:    token.INT,
						Literal: "2",
						Line:    1,
						Column:  1,
					},
					Value: 2,
				}, 2)
			},
		},
		{
			name: "Break",
			node: &ast.BreakExpression{
				Token: token.Token{
					Type:    token.BREAK,
					Literal: "break",
					Line:    1,
					Column:  1,
				},
			},
			object: NewBreak(&ast.BreakExpression{
				Token: token.Token{
					Type:    token.BREAK,
					Literal: "break",
					Line:    1,
					Column:  1,
				},
			}),
			modifyFunction: func(obj Object) {
				// nothing to modify
			},
		},
		{
			name: "Continue",
			node: &ast.ContinueExpression{
				Token: token.Token{
					Type:    token.CONTINUE,
					Literal: "continue",
					Line:    1,
					Column:  1,
				},
			},
			object: NewContinue(&ast.ContinueExpression{
				Token: token.Token{
					Type:    token.CONTINUE,
					Literal: "continue",
					Line:    1,
					Column:  1,
				},
			}),
			modifyFunction: func(obj Object) {
				// nothing to modify
			},
		},
		{
			name: "Null",
			node: &ast.NullLiteral{
				Token: token.Token{
					Type:    token.NULL,
					Literal: "null",
					Line:    1,
					Column:  1,
				},
			},
			object: newNull(&ast.NullLiteral{
				Token: token.Token{
					Type:    token.NULL,
					Literal: "null",
					Line:    1,
					Column:  1,
				},
			}),
			modifyFunction: func(obj Object) {
				// nothing to modify
			},
		},
		{
			name: "Error",
			node: &ast.StringLiteral{
				Value: "error",
				Token: token.Token{
					Type:    token.STRING,
					Literal: "error",
					Line:    2,
					Column:  5,
				},
			},
			object: NewError(&ast.StringLiteral{
				Value: "error",
				Token: token.Token{
					Type:    token.STRING,
					Literal: "error",
					Line:    2,
					Column:  5,
				},
			}, "error"),
			modifyFunction: func(obj Object) {
				obj.(*Error).Message = "new error"
			},
		},
		{
			name: "Function",
			node: &ast.FunctionLiteral{
				Token: token.Token{
					Type:    token.FUNCTION,
					Literal: "func",
					Line:    1,
					Column:  1,
				},
				Arguments: []*ast.Identifier{
					{
						Value: "a",
						Token: token.Token{
							Type:    token.IDENT,
							Literal: "a",
							Line:    1,
							Column:  6,
						},
					},
				},
				Body: &ast.BlockStatement{
					Token: token.Token{
						Type:    token.LBRACE,
						Literal: "{",
						Line:    1,
						Column:  8,
					},
					Statements: []ast.Statement{},
				},
			},
			object: NewFunction(
				&ast.FunctionLiteral{
					Token: token.Token{
						Type:    token.FUNCTION,
						Literal: "func",
						Line:    1,
						Column:  1,
					},
					Arguments: []*ast.Identifier{
						{
							Value: "a",
							Token: token.Token{
								Type:    token.IDENT,
								Literal: "a",
								Line:    1,
								Column:  6,
							},
						},
					},
					Body: &ast.BlockStatement{
						Token: token.Token{
							Type:    token.LBRACE,
							Literal: "{",
							Line:    1,
							Column:  8,
						},
						Statements: []ast.Statement{},
					},
				},
				[]*ast.Identifier{
					{
						Value: "a",
						Token: token.Token{
							Type:    token.IDENT,
							Literal: "a",
							Line:    1,
							Column:  6,
						},
					},
				},
				&ast.BlockStatement{
					Token: token.Token{
						Type:    token.LBRACE,
						Literal: "{",
						Line:    1,
						Column:  8,
					},
					Statements: []ast.Statement{},
				},
				NewEnvironment(),
			),
			modifyFunction: func(obj Object) {
				obj.(*Function).Arguments = []*ast.Identifier{
					{
						Value: "b",
						Token: token.Token{
							Type:    token.IDENT,
							Literal: "b",
							Line:    1,
							Column:  6,
						},
					},
				}
			},
		},
		{
			name: "Builtin",
			node: &ast.Identifier{
				Value: "len",
				Token: token.Token{
					Type:    token.IDENT,
					Literal: "len",
					Line:    1,
					Column:  1,
				},
			},
			object: NewBuiltin(&ast.Identifier{
				Value: "len",
				Token: token.Token{
					Type:    token.IDENT,
					Literal: "len",
					Line:    1,
					Column:  1,
				},
			}, nil),
			modifyFunction: func(obj Object) {
				// nothing to modify
			},
		},
		{
			name: "Array",
			node: &ast.ArrayLiteral{
				Token: token.Token{
					Type:    token.LBRACKET,
					Literal: "[",
					Line:    1,
					Column:  1,
				},
				Elements: []ast.Expression{
					&ast.IntegerLiteral{
						Token: token.Token{
							Type:    token.INT,
							Literal: "1",
						},
						Value: 1,
					},
					&ast.IntegerLiteral{
						Token: token.Token{
							Type:    token.INT,
							Literal: "2",
						},
						Value: 2,
					},
				},
			},
			object: NewArray(&ast.ArrayLiteral{
				Token: token.Token{
					Type:    token.LBRACKET,
					Literal: "[",
					Line:    1,
					Column:  1,
				},
				Elements: []ast.Expression{
					&ast.IntegerLiteral{
						Token: token.Token{
							Type:    token.INT,
							Literal: "1",
						},
						Value: 1,
					},
					&ast.IntegerLiteral{
						Token: token.Token{
							Type:    token.INT,
							Literal: "2",
						},
						Value: 2,
					},
				},
			}, []Object{
				NewInteger(&ast.IntegerLiteral{
					Token: token.Token{
						Type:    token.INT,
						Literal: "1",
					},
					Value: 1,
				}, 1),
				NewInteger(&ast.IntegerLiteral{
					Token: token.Token{
						Type:    token.INT,
						Literal: "2",
					},
					Value: 2,
				}, 2),
			}),
			modifyFunction: func(obj Object) {
				obj.(*Array).Elements = []Object{
					NewInteger(&ast.IntegerLiteral{
						Token: token.Token{
							Type:    token.INT,
							Literal: "3",
						},
						Value: 3,
					}, 3),
				}
			},
		},
		{
			name: "Hash",
			node: &ast.HashLiteral{
				Token: token.Token{
					Type:    token.LBRACE,
					Literal: "{",
					Line:    1,
					Column:  1,
				},
				Pairs: map[ast.Expression]ast.Expression{
					&ast.StringLiteral{
						Token: token.Token{
							Type:    token.STRING,
							Literal: "a",
						},
						Value: "a",
					}: &ast.IntegerLiteral{
						Token: token.Token{
							Type:    token.INT,
							Literal: "1",
						},
						Value: 1,
					},
				},
			},
			object: NewHash(&ast.HashLiteral{
				Token: token.Token{
					Type:    token.LBRACE,
					Literal: "{",
					Line:    1,
					Column:  1,
				},
				Pairs: map[ast.Expression]ast.Expression{
					&ast.StringLiteral{
						Token: token.Token{
							Type:    token.STRING,
							Literal: "a",
						},
						Value: "a",
					}: &ast.IntegerLiteral{
						Token: token.Token{
							Type:    token.INT,
							Literal: "1",
						},
						Value: 1,
					},
				},
			}, map[HashKey]HashPair{
				NewString(&ast.StringLiteral{
					Token: token.Token{
						Type:    token.STRING,
						Literal: "a",
					},
					Value: "a",
				}, "a").HashKey(): {
					Key: NewString(&ast.StringLiteral{
						Token: token.Token{
							Type:    token.STRING,
							Literal: "a",
						},
						Value: "a",
					}, "a"),
					Value: NewInteger(&ast.IntegerLiteral{
						Token: token.Token{
							Type:    token.INT,
							Literal: "1",
						},
						Value: 1,
					}, 1),
				},
			}),
			modifyFunction: func(obj Object) {
				obj.(*Hash).Pairs = map[HashKey]HashPair{
					NewString(&ast.StringLiteral{
						Token: token.Token{
							Type:    token.STRING,
							Literal: "b",
						},
						Value: "b",
					}, "b").HashKey(): {
						Key: NewString(&ast.StringLiteral{
							Token: token.Token{
								Type:    token.STRING,
								Literal: "b",
							},
							Value: "b",
						}, "b"),
						Value: NewInteger(&ast.IntegerLiteral{
							Token: token.Token{
								Type:    token.INT,
								Literal: "2",
							},
							Value: 2,
						}, 2),
					},
				}
			},
		},
		{
			name: "Nonexistent",
			node: &ast.StringLiteral{
				Value: "nonexistent",
				Token: token.Token{
					Type:    token.STRING,
					Literal: "nonexistent",
				},
			},
			object: &Nonexistent{
				node: &ast.StringLiteral{
					Value: "nonexistent",
					Token: token.Token{
						Type:    token.STRING,
						Literal: "nonexistent",
					},
				},
			},
			modifyFunction: func(obj Object) {
				// nothing to modify
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clone := Clone(tc.node, tc.object)

			if tc.name != "Nonexistent" && clone.Type() != tc.object.Type() {
				t.Errorf("cloned object has different type")
			}

			if tc.name != "Nonexistent" && clone.Inspect() != tc.object.Inspect() {
				t.Errorf("cloned object has different value")
			}

			if clone.Node().TokenLiteral() != tc.object.Node().TokenLiteral() {
				t.Errorf("cloned object has different token literal")
			}

			if clone.Node().Line() != tc.object.Node().Line() {
				t.Errorf("cloned object has different line")
			}

			if clone.Node().Column() != tc.object.Node().Column() {
				t.Errorf("cloned object has different column")
			}

			if clone.Node().String() != tc.object.Node().String() {
				t.Errorf("cloned object has different string")
			}

			if clone == tc.object {
				t.Errorf("cloned object is the same object")
			}

			tc.modifyFunction(clone)

			if !contains([]string{"Break", "Continue", "Null", "Builtin"}, tc.name) && clone.Inspect() == tc.object.Inspect() {
				t.Errorf("changed cloned object has same value")
			}
		})
	}
}

func TestIsNumber(t *testing.T) {
	integer := &Integer{Value: 1}
	float := &Float{Value: 1.5}

	if !IsNumber(integer) {
		t.Errorf("integer is not a number")
	}

	if !IsNumber(float) {
		t.Errorf("float is not a number")
	}
}

func TestNewFunctions(t *testing.T) {
	tests := []struct {
		name      string
		newObject Object
		expected  interface{}
	}{
		{BOOLEAN_OBJ, newBoolean(&ast.BooleanLiteral{Value: true}, true), "true"},
		{INTEGER_OBJ, NewInteger(&ast.IntegerLiteral{Value: 1}, 1), "1"},
		{FLOAT_OBJ, NewFloat(&ast.FloatLiteral{Value: 1.5}, 1.5), "1.5"},
		{STRING_OBJ, NewString(&ast.StringLiteral{Value: "hello"}, "hello"), "hello"},
		{ARRAY_OBJ, NewArray(&ast.ArrayLiteral{}, []Object{
			NewInteger(&ast.IntegerLiteral{Value: 1}, 1),
			NewInteger(&ast.IntegerLiteral{Value: 2}, 2),
		}), "[1, 2]"},
		{HASH_OBJ, NewHash(&ast.HashLiteral{}, map[HashKey]HashPair{
			NewString(&ast.StringLiteral{Value: "a"}, "a").HashKey(): {
				Key:   NewString(&ast.StringLiteral{Value: "a"}, "a"),
				Value: NewInteger(&ast.IntegerLiteral{Value: 1}, 1),
			},
		}), "{a: 1}"},
		{NULL_OBJ, newNull(&ast.NullLiteral{}), "null"},
		{RETURN_VALUE_OBJ, NewReturnValue(&ast.ReturnStatement{}, NewInteger(&ast.IntegerLiteral{Value: 2}, 2)), "2"},
		{BREAK_OBJ, NewBreak(&ast.BreakExpression{}), "break"},
		{CONTINUE_OBJ, NewContinue(&ast.ContinueExpression{}), "continue"},
		{ERROR_OBJ, NewError(&ast.StringLiteral{
			Value: "error",
			Token: token.Token{
				Type:    token.STRING,
				Literal: "error",
				Line:    2,
				Column:  5,
			},
		}, "error"), "ERROR: [2:5] error"},
		{FUNCTION_OBJ, NewFunction(
			&ast.FunctionLiteral{},
			[]*ast.Identifier{
				{
					Value: "a",
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "a",
						Line:    1,
						Column:  1,
					},
				},
			},
			&ast.BlockStatement{},
			NewEnvironment(),
		), `func(a) {
}
`},
		{BUILTIN_OBJ, NewBuiltin(&ast.Identifier{}, nil), "builtin function"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if string(test.newObject.Type()) != test.name {
				t.Errorf("object type is not %s, got '%s'", test.name, string(test.newObject.Type()))
			}

			if test.newObject.Inspect() != test.expected {
				t.Errorf("object value is not %v, got '%v'", test.expected, test.newObject.Inspect())
			}

			if test.newObject.Node() == nil {
				t.Errorf("object node is nil")
			}
		})
	}
}
