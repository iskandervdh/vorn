package parser

import (
	"fmt"
	"io"
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

	currentToken token.Token
	peekToken    token.Token

	prefixParseFunctions map[token.TokenType]prefixParseFunction
	infixParseFunctions  map[token.TokenType]infixParseFunction
}

const (
	_ int = iota
	LOWEST
	EQUALS       // ==
	LESS_GREATER // > or <
	SUM          // +
	PRODUCT      // *
	PREFIX       // -X or !X
	CALL         // myFunction(X)
	INDEX        // array[index]
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESS_GREATER,
	token.GT:       LESS_GREATER,
	token.LTE:      LESS_GREATER,
	token.GTE:      LESS_GREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens, so currentToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	p.prefixParseFunctions = make(map[token.TokenType]prefixParseFunction)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.NULL, p.parseNull)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.EXCLAMATION, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.BREAK, p.parseBreakExpression)
	p.registerPrefix(token.CONTINUE, p.parseContinueExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)
	p.registerPrefix(token.ASSIGN, p.parseReassignLiteral)

	p.infixParseFunctions = make(map[token.TokenType]infixParseFunction)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LTE, p.parseInfixExpression)
	p.registerInfix(token.GTE, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	return p
}

func (p *Parser) registerPrefix(tokenType token.TokenType, function prefixParseFunction) {
	p.prefixParseFunctions[tokenType] = function
}

func (p *Parser) registerInfix(tokenType token.TokenType, function infixParseFunction) {
	p.infixParseFunctions[tokenType] = function
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

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}

	return LOWEST
}

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

func (p *Parser) peekError(t token.TokenType) {
	p.addError(fmt.Sprintf("expected '%s', got %s instead", t, p.peekToken.Type), true)
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()

		return true
	} else {
		p.peekError(t)

		return false
	}
}

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

func (p *Parser) parseFunction() ast.Statement {
	if p.peekTokenIs(token.IDENT) {
		return p.parseFunctionStatement()
	}

	return p.parseExpressionStatement()
}

func (p *Parser) parseFunctionStatement() *ast.FunctionStatement {
	statement := &ast.FunctionStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	statement.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	statement.Body = p.parseBlockStatement()

	return statement
}

func (p *Parser) parseVariableStatement() *ast.VariableStatement {
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

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()

	statement.ReturnValue = p.parseExpression(LOWEST)

	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
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

func (p *Parser) parseForStatement() *ast.ForStatement {
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

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFunctions[p.currentToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.currentToken.Type)

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

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: p.currentToken}
	statement.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseReassignmentExpression() ast.Expression {
	statement := &ast.ReassignmentExpression{}

	statement.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
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

func (p *Parser) parseIdentifier() ast.Expression {
	identifier := &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	if p.currentTokenIs(token.IDENT) && p.peekTokenIs(token.ASSIGN) {
		return p.parseReassignmentExpression()
	}

	return identifier
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: p.currentToken}
	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)

	if err != nil {
		p.addError(fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal), false)

		return nil
	}

	literal.Value = value

	return literal
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.BooleanLiteral{Token: p.currentToken, Value: p.currentTokenIs(token.TRUE)}
}

func (p *Parser) parseNull() ast.Expression {
	return &ast.NullLiteral{Token: p.currentToken}
}

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

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	p.addError(fmt.Sprintf("unexpected token %s", t), false)
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
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

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.currentToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)

	return exp
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	expression := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return expression
}

func (p *Parser) parseIfExpression() ast.Expression {
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

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBreakExpression() ast.Expression {
	return &ast.BreakExpression{Token: p.currentToken}
}

func (p *Parser) parseContinueExpression() ast.Expression {
	return &ast.ContinueExpression{Token: p.currentToken}
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token:      p.currentToken,
		Statements: []ast.Statement{},
	}

	// Create a new scope for the block
	block.Parent = p.scope
	p.scope = block

	p.nextToken()

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

func (p *Parser) parseFunctionLiteral() ast.Expression {
	fl := &ast.FunctionLiteral{Token: p.currentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	fl.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	fl.Body = p.parseBlockStatement()

	return fl
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
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

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.currentToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)

	return array
}

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

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expression := &ast.IndexExpression{Token: p.currentToken, Left: left}

	p.nextToken()
	expression.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return expression
}

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

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return hash
}

func (p *Parser) parseReassignLiteral() ast.Expression {
	p.nextToken()

	return p.parseExpression(LOWEST)
}
