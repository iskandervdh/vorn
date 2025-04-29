package parser

import (
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/iskandervdh/vorn/ast"
	"github.com/iskandervdh/vorn/lexer"
	"github.com/iskandervdh/vorn/token"
)

type (
	prefixParseFunction func() ast.Expression
	infixParseFunction  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	scope  ast.Scope
	errors []string
	trace  bool

	currentToken token.Token
	peekToken    token.Token

	prefixParseFunctions map[token.TokenType]prefixParseFunction
	infixParseFunctions  map[token.TokenType]infixParseFunction
}

/*
Precedence levels for the different operators
*/
const (
	_ int = iota
	LOWEST
	OR           // ||
	AND          // &&
	BITWISE_OR   // |
	BITWISE_XOR  // ^
	BITWISE_AND  // &
	EQUALS       // == or !=
	LESS_GREATER // >, <, >= or <=
	SHIFT        // << or >>
	SUM          // + or -
	PRODUCT      // * or /
	PREFIX       // -x, !x or ~x
	POSTFIX      // object.property, object.method(args), function(X), array[index]
)

/*
Precedence levels for the tokens based on their operator type
*/
var precedences = map[token.TokenType]int{
	token.OR:          OR,
	token.AND:         AND,
	token.BITWISE_OR:  BITWISE_OR,
	token.BITWISE_XOR: BITWISE_XOR,
	token.BITWISE_AND: BITWISE_AND,
	token.EQ:          EQUALS,
	token.NOT_EQ:      EQUALS,
	token.LT:          LESS_GREATER,
	token.GT:          LESS_GREATER,
	token.LTE:         LESS_GREATER,
	token.GTE:         LESS_GREATER,
	token.LEFT_SHIFT:  SHIFT,
	token.RIGHT_SHIFT: SHIFT,
	token.PLUS:        SUM,
	token.MINUS:       SUM,
	token.ASTERISK:    PRODUCT,
	token.SLASH:       PRODUCT,
	token.PERCENT:     PRODUCT,
	token.DOT:         POSTFIX,
	token.LPAREN:      POSTFIX,
	token.LBRACKET:    POSTFIX,
}

/*
Create a new parser with the given lexer and trace flag
*/
func New(l *lexer.Lexer, trace bool) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
		trace:  trace,
	}

	// Read two tokens, so currentToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	// Register the prefix and infix functions to be used
	// to determine what to do when we encounter a prefix or infix token
	p.prefixParseFunctions = make(map[token.TokenType]prefixParseFunction)
	p.registerPrefixFunctions()

	p.infixParseFunctions = make(map[token.TokenType]infixParseFunction)
	p.registerInfixFunctions()

	return p
}

func (p *Parser) registerPrefix(tokenType token.TokenType, function prefixParseFunction) {
	p.prefixParseFunctions[tokenType] = function
}

func (p *Parser) registerPrefixFunctions() {
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.NULL, p.parseNull)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.EXCLAMATION, p.parsePrefixExpression)
	p.registerPrefix(token.BITWISE_NOT, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.BREAK, p.parseBreakExpression)
	p.registerPrefix(token.CONTINUE, p.parseContinueExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)
}

func (p *Parser) registerInfix(tokenType token.TokenType, function infixParseFunction) {
	p.infixParseFunctions[tokenType] = function
}

func (p *Parser) registerInfixFunctions() {
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.PERCENT, p.parseInfixExpression)
	p.registerInfix(token.BITWISE_OR, p.parseInfixExpression)
	p.registerInfix(token.BITWISE_XOR, p.parseInfixExpression)
	p.registerInfix(token.BITWISE_AND, p.parseInfixExpression)
	p.registerInfix(token.LEFT_SHIFT, p.parseInfixExpression)
	p.registerInfix(token.RIGHT_SHIFT, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LTE, p.parseInfixExpression)
	p.registerInfix(token.GTE, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.DOT, p.parseChainingExpression)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func PrintErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Syntax errors:\n")

	for _, msg := range errors {
		io.WriteString(out, msg+"\n")
	}
}

/*
Get the precedence of the peek token
*/
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

/*
Get the precedence of the current token
*/
func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}

	return LOWEST
}

/*
Add an error to the parser
*/
func (p *Parser) addError(e string, peek bool) {
	var t token.Token

	if peek {
		t = p.peekToken
	} else {
		t = p.currentToken
	}

	e = fmt.Sprintf("[%d:%d]: %s", t.Line, t.Column, e)

	p.errors = append(p.errors, e)
}

/*
Add a peek error to the parser meaning the parser expected a certain token type but got another
*/
func (p *Parser) addPeekError(t token.TokenType) {
	p.addError(fmt.Sprintf("expected '%s', got %s instead", t, p.peekToken.Type), true)
}

/*
Move to the next token in the lexer by setting the current token to the peek token
and the peek token to the next token in the lexer
*/
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

/*
Check if the current token is of a certain type
*/
func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

/*
Check if the peek token is of a certain type
*/
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

/*
Check if the peek token is an assignment operator
*/
func (p *Parser) peekTokenIsAssignmentOperator() bool {
	return slices.Contains(token.AssignmentOperators, p.peekToken.Type)
}

/*
Expect the next token to be of a certain type.

If the next token is of the expected type, the parser will move to the next token and return true.
If the next token is not of the expected type, the parser will add an error and return false.
*/
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()

		return true
	}

	p.addPeekError(t)

	return false

}

/*
Expect the next token to be an assignment operator.

If the next token is an assignment operator, the parser will move to the next token and return true.
If the next token is not an assignment operator, the parser will add an error and return false.
*/
func (p *Parser) expectReassignmentOperator() bool {
	if p.peekTokenIsAssignmentOperator() {
		p.nextToken()

		return true
	}

	p.addError(fmt.Sprintf("expected assignment operator, got %s instead", p.peekToken.Type), true)

	return false
}

/*
Check if the current token is a reassignment of a constant.

Add an error if it is
*/
func (p *Parser) checkConstReassignment(scope ast.Scope, newStatement ast.Statement) {
	expressionStatement, ok := newStatement.(*ast.ExpressionStatement)

	if !ok || expressionStatement == nil {
		return
	}

	reassignmentExpression, ok := expressionStatement.Expression.(*ast.ReassignmentExpression)

	if !ok || reassignmentExpression == nil {
		return
	}

	definedAsLetInScope := false

	for _, statement := range scope.GetScopeStatements() {
		variableStatement, ok := statement.(*ast.VariableStatement)

		if !ok {
			continue
		}

		if !variableStatement.IsConst() {
			definedAsLetInScope = true
			continue
		}

		if variableStatement.Name.Value == reassignmentExpression.Name.Value {
			e := fmt.Sprintf("[%d:%d] can not reassign constant %s.",
				reassignmentExpression.Token.Line,
				reassignmentExpression.Token.Column,
				reassignmentExpression.Name.Value,
			)

			p.errors = append(p.errors, e)
			return
		}
	}

	if scope.GetParentScope() != nil && !definedAsLetInScope {
		p.checkConstReassignment(scope.GetParentScope(), newStatement)
	}
}

func (p *Parser) checkVariableRedefinition(statements []ast.Statement, newStatement ast.Statement) {
	expressionStatement, ok := newStatement.(*ast.VariableStatement)

	if !ok || expressionStatement == nil {
		return
	}

	for _, statement := range statements {
		variableStatement, ok := statement.(*ast.VariableStatement)

		if !ok {
			continue
		}

		if variableStatement.Name.Value == expressionStatement.Name.Value {
			e := fmt.Sprintf("[%d:%d] can not redefine variable %s.",
				expressionStatement.Token.Line,
				expressionStatement.Token.Column,
				expressionStatement.Name.Value,
			)

			p.errors = append(p.errors, e)
		}
	}
}

/*
Parse the program by parsing all the statements in the program.

Check for reassignments of constants and redefinitions of variables after each statement.

Return the parsed program
*/
func (p *Parser) ParseProgram() *ast.Program {
	program := ast.NewProgram()
	p.scope = program

	for p.currentToken.Type != token.EOF {
		statement := p.parseStatement()

		p.checkConstReassignment(program, statement)
		p.checkVariableRedefinition(program.Statements, statement)

		program.Statements = append(program.Statements, statement)

		p.nextToken()
	}

	return program
}

/*
Parse a statement based on the current token type
*/
func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseVariableStatement()
	case token.CONST:
		return p.parseVariableStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.FUNCTION:
		return p.parseFunction()
	case token.FOR:
		return p.parseForStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	default:
		return p.parseExpressionStatement()
	}
}

/*
Parse a named or anonymous function based on if the peek token is an identifier
*/
func (p *Parser) parseFunction() ast.Statement {
	if p.peekTokenIs(token.IDENT) {
		return p.parseFunctionStatement()
	}

	return p.parseExpressionStatement()
}

/*
Parse a function statement, including the function name, arguments and body
*/
func (p *Parser) parseFunctionStatement() *ast.FunctionStatement {
	statement := &ast.FunctionStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDENT) { // coverage-ignore
		return nil
	}

	statement.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	statement.Arguments = p.parseFunctionArguments()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	statement.Body = p.parseBlockStatement()

	return statement
}

/*
Parse a variable statement, including the variable name and value
*/
func (p *Parser) parseVariableStatement() *ast.VariableStatement {
	if p.trace { // coverage-ignore
		defer untrace(trace("VariableStatement"))
	}

	statement := &ast.VariableStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	statement.Value = p.parseExpression(LOWEST)

	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

/*
Parse a return statement, including the return value
*/
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()

	statement.ReturnValue = p.parseExpression(LOWEST)

	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

/*
Parse a while statement, including the condition and body
*/
func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	if p.trace { // coverage-ignore
		defer untrace(trace("WhileStatement"))
	}

	whileStatement := &ast.WhileStatement{Token: p.currentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	whileStatement.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	whileStatement.Consequence = p.parseBlockStatement()

	return whileStatement
}

/*
Parse a for statement, including the initialization, condition, update and body
*/
func (p *Parser) parseForStatement() *ast.ForStatement {
	if p.trace { // coverage-ignore
		defer untrace(trace("ForStatement"))
	}

	// Initialize for loop scope
	forStatement := &ast.ForStatement{Token: p.currentToken}
	forStatement.Parent = p.scope
	p.scope = forStatement

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	if p.currentTokenIs(token.LET) {
		forStatement.Init = p.parseVariableStatement()
	} else if p.currentTokenIs(token.CONST) {
		forStatement.Init = p.parseVariableStatement()
	} else if p.currentTokenIs(token.SEMICOLON) {
		forStatement.Init = nil
	} else {
		forStatement.Init = p.parseExpressionStatement()
	}

	p.nextToken()
	forStatement.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	p.nextToken()

	if p.currentTokenIs(token.SEMICOLON) {
		forStatement.Update = nil
	} else {
		forStatement.Update = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	forStatement.Body = p.parseBlockStatement()

	// Restore the parent scope
	p.scope = p.scope.GetParentScope()

	return forStatement
}

/*
Parse an expression based on the current token type

The precedence is used to determine the order of operations.
The parser will keep parsing expressions until it encounters a semicolon or the precedence of the next token is lower than the current precedence

The parser will then return the parsed expression
*/
func (p *Parser) parseExpression(precedence int) ast.Expression {
	if p.trace { // coverage-ignore
		defer untrace(trace("Expression"))
	}

	prefix := p.prefixParseFunctions[p.currentToken.Type]

	if prefix == nil {
		p.noPrefixParseFunctionError(p.currentToken.Type)

		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFunctions[p.peekToken.Type]

		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

/*
Parse an expression statement, including the expression and a semicolon
*/
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	if p.trace { // coverage-ignore
		defer untrace(trace("ExpressionStatement"))
	}

	statement := &ast.ExpressionStatement{Token: p.currentToken}
	statement.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

/*
Parse a reassignment expression, including the name, value and a semicolon
*/
func (p *Parser) parseReassignmentExpression() ast.Expression {
	if p.trace { // coverage-ignore
		defer untrace(trace("ReassignmentExpression"))
	}

	statement := &ast.ReassignmentExpression{}

	statement.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectReassignmentOperator() { // coverage-ignore
		return nil
	}

	statement.Token = p.currentToken

	p.nextToken()

	statement.Value = p.parseExpression(LOWEST)

	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

/*
Parse an identifier expression based on the current token

If the current token is an identifier, the parser will check if the next token is an assignment operator.
If the next token is an assignment operator, the parser will parse a reassignment expression.
If the next token is not an assignment operator, the parser will parse an identifier expression.

The parser will then return the parsed expression
*/
func (p *Parser) parseIdentifier() ast.Expression {
	identifier := &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	if p.currentTokenIs(token.IDENT) && p.peekTokenIsAssignmentOperator() {
		return p.parseReassignmentExpression()
	}

	return identifier
}

/*
Parse an integer literal expression based on the current token
*/
func (p *Parser) parseIntegerLiteral() ast.Expression {
	if p.trace { // coverage-ignore
		defer untrace(trace("IntegerLiteral"))
	}

	literal := &ast.IntegerLiteral{Token: p.currentToken}
	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)

	if err != nil {
		p.addError(fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal), false)

		return nil
	}

	literal.Value = value

	return literal
}

/*
Parse a boolean expression based on the current token
*/
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.BooleanLiteral{Token: p.currentToken, Value: p.currentTokenIs(token.TRUE)}
}

/*
Parse a null expression based on the current token
*/
func (p *Parser) parseNull() ast.Expression {
	return &ast.NullLiteral{Token: p.currentToken}
}

/*
Parse a float literal expression based on the current token
*/
func (p *Parser) parseFloatLiteral() ast.Expression {
	literal := &ast.FloatLiteral{Token: p.currentToken}
	value, err := strconv.ParseFloat(p.currentToken.Literal, 64)

	if err != nil {
		p.addError(fmt.Sprintf("could not parse %q as float", p.currentToken.Literal), false)

		return nil
	}

	literal.Value = value

	return literal
}

/*
Add an error to the parser if there is no prefix parse function for the current token
*/
func (p *Parser) noPrefixParseFunctionError(t token.TokenType) {
	p.addError(fmt.Sprintf("unexpected token %s", t), false)
}

/*
Parse a prefix expression, including the operator and right-hand side
*/
func (p *Parser) parsePrefixExpression() ast.Expression {
	if p.trace { // coverage-ignore
		defer untrace(trace("PrefixExpression"))
	}

	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

/*
Parse an infix expression, including the left-hand side, operator and right-hand side
*/
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	if p.trace { // coverage-ignore
		defer untrace(trace("InfixExpression"))
	}

	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

/*
Parse a call expression, including the function name and the call arguments
*/
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	if p.trace { // coverage-ignore
		defer untrace(trace("CallExpression"))
	}

	exp := &ast.CallExpression{Token: p.currentToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)

	return exp
}

/*
Parse a grouped expression, including the expression inside the parentheses
*/
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	expression := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return expression
}

/*
Parse an if expression, including the condition, consequence and alternative
*/
func (p *Parser) parseIfExpression() ast.Expression {
	if p.trace { // coverage-ignore
		defer untrace(trace("IfExpression"))
	}

	expression := &ast.IfExpression{Token: p.currentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		// TODO: Add support for else if

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

/*
Parse a break expression
*/
func (p *Parser) parseBreakExpression() ast.Expression {
	return &ast.BreakExpression{Token: p.currentToken}
}

/*
Parse a continue expression
*/
func (p *Parser) parseContinueExpression() ast.Expression {
	return &ast.ContinueExpression{Token: p.currentToken}
}

/*
Parse a block statement, including all the statements inside the block.

The parser will create a new scope for the block and restore the parent scope after parsing the block.
It will also check for reassignments of constants and redefinitions of variables after each statement.

The parser will then return the parsed block statement
*/
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	if p.trace { // coverage-ignore
		defer untrace(trace("BlockStatement"))
	}

	block := &ast.BlockStatement{
		Token:      p.currentToken,
		Statements: []ast.Statement{},
	}

	// Create a new scope for the block
	block.Parent = p.scope
	p.scope = block

	p.nextToken()

	// Parse all the statements inside the block
	for !p.currentTokenIs(token.RBRACE) && !p.currentTokenIs(token.EOF) {
		statement := p.parseStatement()

		p.checkConstReassignment(block, statement)
		p.checkVariableRedefinition(block.Statements, statement)

		block.Statements = append(block.Statements, statement)

		p.nextToken()
	}

	// Restore the parent scope
	p.scope = p.scope.GetParentScope()

	return block
}

/*
Parse a function literal, including the arguments and body
*/
func (p *Parser) parseFunctionLiteral() ast.Expression {
	fl := &ast.FunctionLiteral{Token: p.currentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	fl.Arguments = p.parseFunctionArguments()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	fl.Body = p.parseBlockStatement()

	return fl
}

/*
Parse a the arguments of a function until a right parenthesis is encountered
*/
func (p *Parser) parseFunctionArguments() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		ident := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

/*
Parse a string literal
*/
func (p *Parser) parseStringLiteral() ast.Expression {
	if p.trace { // coverage-ignore
		defer untrace(trace("StringLiteral"))
	}

	return &ast.StringLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
}

/*
Parse an array literal, including all the elements inside the array
*/
func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.currentToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)

	return array
}

/*
Parse a expression list until a certain token type is encountered
*/
func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

/*
Parse an index expression, including the left-hand side and the index
*/
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expression := &ast.IndexExpression{Token: p.currentToken, Left: left}

	p.nextToken()
	expression.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return expression
}

/*
Parse a chaining expression, including the left-hand side and the right-hand side

The parser will continue parsing chained expressions if there are any
*/
func (p *Parser) parseChainingExpression(left ast.Expression) ast.Expression {
	if p.trace { // coverage-ignore
		defer untrace(trace("ChainingExpression"))
	}

	expression := &ast.ChainingExpression{Token: p.currentToken, Left: left}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// Parse the right-hand side
	identifier := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	// Check if the next token is a left parenthesis indicating a function call
	if p.peekTokenIs(token.LPAREN) {
		p.nextToken()
		expression.Right = p.parseCallExpression(identifier)
	} else {
		expression.Right = identifier
	}

	// Continue parsing chained expressions if there are any
	for p.peekTokenIs(token.DOT) {
		p.nextToken()
		return p.parseChainingExpression(expression)
	}

	return expression
}

/*
Parse a object literal, including all the key-value pairs inside the object
*/
func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{
		Token: p.currentToken,
		Pairs: make(map[ast.Expression]ast.Expression),
	}

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()

		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()

		value := p.parseExpression(LOWEST)
		hash.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) { // coverage-ignore
		return nil
	}

	return hash
}
