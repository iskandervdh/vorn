package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/iskandervdh/vorn/ast"
	"github.com/iskandervdh/vorn/lexer"
)

func initializeParserTest(t *testing.T, input string, expectedStatementCount int) *ast.Program {
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if expectedStatementCount > 0 && len(program.Statements) != expectedStatementCount {
		t.Fatalf("program.Statements does not contain %d statements. got %d",
			expectedStatementCount,
			len(program.Statements),
		)
	}

	return program
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

func checkVariableStatement(t *testing.T, s ast.Statement, expectedLiteral string, name string) bool {
	if s.TokenLiteral() != expectedLiteral {
		t.Errorf("s.TokenLiteral not '%s'. got %q", expectedLiteral, s.TokenLiteral())
		return false
	}

	variableStatement, ok := s.(*ast.VariableStatement)

	if !ok {
		t.Errorf("s not *ast.VariableStatement. got %T", s)
		return false
	}

	if variableStatement.Name.Value != name {
		t.Errorf("variableStatement.Name.Value not '%s'. got '%s'", name, variableStatement.Name.Value)
		return false
	}

	if variableStatement.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got '%s'", name, variableStatement.Name)
		return false
	}

	if (expectedLiteral == "let" && !variableStatement.IsLet()) ||
		(expectedLiteral == "const" && !variableStatement.IsConst()) {
		t.Errorf("s not '%s'. got %q", expectedLiteral, variableStatement.TokenLiteral())
		return false
	}

	return true
}

func checkLetStatement(t *testing.T, s ast.Statement, name string) bool {
	return checkVariableStatement(t, s, "let", name)
}

func checkConstStatement(t *testing.T, s ast.Statement, name string) bool {
	return checkVariableStatement(t, s, "const", name)
}

func checkFunctionStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "func" {
		t.Errorf("s.TokenLiteral not '%s'. got %q", "func", s.TokenLiteral())
		return false
	}

	functionStatement, ok := s.(*ast.FunctionStatement)

	if !ok {
		t.Errorf("s not *ast.FunctionStatement. got %T", s)
		return false
	}

	if functionStatement.Name.Value != name {
		t.Errorf("functionStatement.Name.Value not '%s'. got '%s'", name, functionStatement.Name.Value)
		return false
	}

	if functionStatement.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got '%s'", name, functionStatement.Name)
		return false
	}

	return true
}

func checkIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got %T", il)
		return false
	}

	if integer.Value != value {
		t.Errorf("integer.Value not %d. got %d", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral not %d. got '%s'", value, integer.TokenLiteral())
		return false
	}

	return true
}

func checkFloatLiteral(t *testing.T, il ast.Expression, value float64) bool {
	float, ok := il.(*ast.FloatLiteral)

	if !ok {
		t.Errorf("il not *ast.FloatLiteral. got %T", il)
		return false
	}

	if float.Value != value {
		t.Errorf("float.Value not %f. got %f", value, float.Value)
		return false
	}

	floatStr := fmt.Sprintf("%g", value)

	if !strings.Contains(floatStr, ".") {
		floatStr += ".0"
	}

	if float.TokenLiteral() != floatStr {
		t.Errorf("float.TokenLiteral not %g. got '%s'", value, float.TokenLiteral())
		return false
	}

	return true
}

func checkIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)

	if !ok {
		t.Errorf("exp not *ast.Identifier. got %T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got '%s'", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got '%s'", value, ident.TokenLiteral())
		return false
	}

	return true
}

func checkBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	boolean, ok := exp.(*ast.BooleanLiteral)

	if !ok {
		t.Errorf("exp not *ast.Boolean. got %T", exp)
		return false
	}

	if boolean.Value != value {
		t.Errorf("bo.Value not %t. got %T", value, boolean.Value)
		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got '%s'", value, boolean.TokenLiteral())
		return false
	}

	return true
}

func checkLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return checkIntegerLiteral(t, exp, int64(v))
	case int64:
		return checkIntegerLiteral(t, exp, v)
	case float64:
		return checkFloatLiteral(t, exp, v)
	case string:
		return checkIdentifier(t, exp, v)
	case bool:
		return checkBooleanLiteral(t, exp, v)
	}

	t.Errorf("type of exp not handled. got %T", exp)
	return false
}

func checkInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	infixExpression, ok := exp.(*ast.InfixExpression)

	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got %T(%s)", exp, exp)
		return false
	}

	if !checkLiteralExpression(t, infixExpression.Left, left) {
		return false
	}

	if infixExpression.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got %q", operator, infixExpression.Operator)
		return false
	}

	if !checkLiteralExpression(t, infixExpression.Right, right) {
		return false
	}

	return true
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, test := range tests {
		program := initializeParserTest(t, test.input, 1)
		statement := program.Statements[0]

		if !checkLetStatement(t, statement, test.expectedIdentifier) {
			return
		}

		value := statement.(*ast.VariableStatement).Value

		if !checkLiteralExpression(t, value, test.expectedValue) {
			return
		}
	}
}

func TestConstStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"const x = 10;", "x", 10},
		{"const y = false;", "y", false},
		{"const foobar = x;", "foobar", "x"},
	}

	for _, test := range tests {
		program := initializeParserTest(t, test.input, 1)
		statement := program.Statements[0]

		if !checkConstStatement(t, statement, test.expectedIdentifier) {
			return
		}

		value := statement.(*ast.VariableStatement).Value

		if !checkLiteralExpression(t, value, test.expectedValue) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, test := range tests {
		program := initializeParserTest(t, test.input, 1)
		statement := program.Statements[0]

		returnStatement, ok := statement.(*ast.ReturnStatement)

		if !ok {
			t.Errorf("statement not *ast.returnStatement. got %T", statement)
			continue
		}

		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("returnStatement.TokenLiteral not 'return', got %q", returnStatement.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	program := initializeParserTest(t, "foobar;", 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T",
			program.Statements[0])
	}

	ident, ok := statement.Expression.(*ast.Identifier)

	if !ok {
		t.Fatalf("exp not *ast.Identifier. got %T", statement.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got '%s'", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got '%s'", "foobar",
			ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	program := initializeParserTest(t, "5;", 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T", program.Statements[0])
	}

	literal, ok := statement.Expression.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got %T", statement.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got %d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got '%s'", "5", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, prefixTest := range prefixTests {
		program := initializeParserTest(t, prefixTest.input, 1)

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T", program.Statements[0])
		}

		expression, ok := statement.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("statement is not ast.PrefixExpression. got %T", statement.Expression)
		}

		if expression.Operator != prefixTest.operator {
			t.Fatalf("expression.Operator is not '%s'. got '%s'", prefixTest.operator, expression.Operator)
		}

		if !checkLiteralExpression(t, expression.Right, prefixTest.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 <= 5;", 5, "<=", 5},
		{"5 >= 5;", 5, ">=", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"2.0 + 3.0;", 2.0, "+", 3.0},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar >= barfoo;", "foobar", ">=", "barfoo"},
		{"foobar <= barfoo;", "foobar", "<=", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}
	for _, infixTest := range infixTests {
		program := initializeParserTest(t, infixTest.input, 1)

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T", program.Statements[0])
		}

		if !checkInfixExpression(t, statement.Expression, infixTest.leftValue, infixTest.operator, infixTest.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"3 >= 5 == false",
			"((3 >= 5) == false)",
		},
		{
			"3 <= 5 == true",
			"((3 <= 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {
		program := initializeParserTest(t, tt.input, -1)
		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected %q, got %q", tt.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, test := range tests {
		program := initializeParserTest(t, test.input, 1)

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T", program.Statements[0])
		}

		boolean, ok := statement.Expression.(*ast.BooleanLiteral)

		if !ok {
			t.Fatalf("exp not *ast.Boolean. got %T", statement.Expression)
		}
		if boolean.Value != test.expectedBoolean {
			t.Errorf("boolean.Value not %t. got %T", test.expectedBoolean, boolean.Value)
		}
	}
}

func TestIfExpression(t *testing.T) {
	program := initializeParserTest(t, "if (x < y) { x }", 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T", program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got %T",
			stmt.Expression)
	}
	if !checkInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got %d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got %T",
			exp.Consequence.Statements[0])
	}
	if !checkIdentifier(t, consequence.Expression, "x") {
		return
	}
	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got %+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	program := initializeParserTest(t, "if (x < y) { x } else { y }", 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)

	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got %T", stmt.Expression)
	}

	if !checkInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got %d\n", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got %T", exp.Consequence.Statements[0])
	}

	if !checkIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements does not contain 1 statements. got %d\n",
			len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got %T",
			exp.Alternative.Statements[0])
	}

	if !checkIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	program := initializeParserTest(t, "func(x, y) { x + y; }", 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T", program.Statements[0])
	}

	function, ok := statement.Expression.(*ast.FunctionLiteral)

	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got %T", statement.Expression)
	}
	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got %d\n", len(function.Parameters))
	}

	checkLiteralExpression(t, function.Parameters[0], "x")
	checkLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got %d\n", len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got %T", function.Body.Statements[0])
	}

	checkInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "func() {};", expectedParams: []string{}},
		{input: "func(x) {};", expectedParams: []string{"x"}},
		{input: "func(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, test := range tests {
		program := initializeParserTest(t, test.input, -1)

		statement := program.Statements[0].(*ast.ExpressionStatement)
		function := statement.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(test.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got %d\n", len(test.expectedParams), len(function.Parameters))
		}

		for i, ident := range test.expectedParams {
			checkLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	program := initializeParserTest(t, "add(1, 2 * 3, 4 + 5);", 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got %T",
			program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.CallExpression)

	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got %T", statement.Expression)
	}

	if !checkIdentifier(t, expression.Function, "add") {
		return
	}

	if len(expression.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got %d", len(expression.Arguments))
	}

	checkLiteralExpression(t, expression.Arguments[0], 1)
	checkInfixExpression(t, expression.Arguments[1], 2, "*", 3)
	checkInfixExpression(t, expression.Arguments[2], 4, "+", 5)
}

func TestFunctionStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
	}{
		{"func test() { return 2; }", "test"},
		{"func add(x, y) { return x + y; }", "add"},
	}

	for _, test := range tests {
		program := initializeParserTest(t, test.input, 1)
		statement := program.Statements[0]

		if !checkFunctionStatement(t, statement, test.expectedIdentifier) {
			return
		}
	}
}

func TestStringLiteralExpression(t *testing.T) {
	program := initializeParserTest(t, `"hello world";`, 1)

	statement := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := statement.Expression.(*ast.StringLiteral)

	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got %T", statement.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got %q", "hello world", literal.Value)
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	program := initializeParserTest(t, "[1, 2 * 2, 3 + 3]", 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T", program.Statements[0])
	}

	array, ok := statement.Expression.(*ast.ArrayLiteral)

	if !ok {
		t.Fatalf("array not ast.ArrayLiteral. got %T", statement.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got %d", len(array.Elements))
	}

	checkIntegerLiteral(t, array.Elements[0], 1)
	checkInfixExpression(t, array.Elements[1], 2, "*", 2)
	checkInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	program := initializeParserTest(t, "myArray[1 + 1]", 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T", program.Statements[0])
	}

	indexExpression, ok := statement.Expression.(*ast.IndexExpression)

	if !ok {
		t.Fatalf("indexExpression not *ast.IndexExpression. got %T", statement.Expression)
	}

	if !checkIdentifier(t, indexExpression.Left, "myArray") {
		return
	}

	if !checkInfixExpression(t, indexExpression.Index, 1, "+", 1) {
		return
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	program := initializeParserTest(t, `{"one": 1, "two": 2, "three": 3}`, 1)

	statement := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := statement.Expression.(*ast.HashLiteral)

	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got %T", statement.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got %d", len(hash.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)

		if !ok {
			t.Errorf("key is not ast.StringLiteral. got %T", key)
		}

		expectedValue := expected[literal.String()]

		checkIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	program := initializeParserTest(t, "{}", 1)

	statement := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := statement.Expression.(*ast.HashLiteral)

	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got %T", statement.Expression)
	}
	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs has wrong length. got %d", len(hash.Pairs))
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	program := initializeParserTest(t, `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`, 1)

	statement := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := statement.Expression.(*ast.HashLiteral)

	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got %T", statement.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got %d", len(hash.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			checkInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			checkInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			checkInfixExpression(t, e, 15, "/", 5)
		},
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)

		if !ok {
			t.Errorf("key is not ast.StringLiteral. got %T", key)
			continue
		}

		testFunction, ok := tests[literal.String()]

		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}

		testFunction(value)
	}
}

func TestParsingReassignment(t *testing.T) {
	program := initializeParserTest(t, "let x = 4; x = 5;", 2)

	checkLetStatement(t, program.Statements[0], "x")

	assignment, ok := program.Statements[1].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[1] is not ast.ExpressionStatement. got %T", program.Statements[1])
	}

	expression, ok := assignment.Expression.(*ast.ReassignmentExpression)

	if !ok {
		t.Fatalf("expression not *ast.ReassignmentExpression. got %T", assignment.Expression)
	}

	if expression.Name.Value != "x" {
		t.Errorf("expression.Name.Value not %s. got %s", "x", expression.Value)
	}

	if expression.Token.Literal != "=" {
		t.Errorf("expression.Token.Literal not %s. got %s", "=", expression.Token.Literal)
	}
}

func TestParsingConstReassignmentError(t *testing.T) {
	input := `const NAME = "YOU";
NAME = "ME";`

	l := lexer.New(input)
	p := New(l)
	p.ParseProgram()

	errors := p.Errors()

	if len(errors) != 1 {
		t.Error("Expected a parser error")
	}

	expectedError := "[2:7] can not reassign constant NAME."

	if errors[0] != expectedError {
		t.Errorf("Expected error message to be %q, got %q", expectedError, errors[0])
	}
}
