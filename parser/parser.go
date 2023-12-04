package parser

import (
	"fmt"
	"github.com/Neal-C/interpreter-in-go/ast"
	"github.com/Neal-C/interpreter-in-go/lexer"
	"github.com/Neal-C/interpreter-in-go/token"
	"strconv"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	// // myFunction(x)
	CALL
)

type Parser struct {
	lexer          *lexer.Lexer
	currentToken   token.Token
	peekToken      token.Token
	errors         []string
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func (self *Parser) registerPrefix(tokenKey token.TokenType, associatedFn prefixParseFn) {
	self.prefixParseFns[tokenKey] = associatedFn
}

func (self *Parser) registerInfix(tokenKey token.TokenType, associatedFn infixParseFn) {
	self.infixParseFns[tokenKey] = associatedFn
}

func (self *Parser) nextToken() {
	self.currentToken = self.peekToken
	self.peekToken = self.lexer.NextToken()
}

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer:  lexer,
		errors: []string{},
	}

	parser.prefixParseFns = make(map[token.TokenType]prefixParseFn)

	parser.registerPrefix(token.IDENT, parser.parseIdentifier)
	parser.registerPrefix(token.INT, parser.parseIntegerLiteral)
	parser.registerPrefix(token.BANG, parser.parsePrefixExpression)
	parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)
	parser.registerPrefix(token.TRUE, parser.parseBoolean)
	parser.registerPrefix(token.FALSE, parser.parseBoolean)
	parser.registerPrefix(token.LPAREN, parser.parseGroupedExpression)
	parser.registerPrefix(token.IF, parser.parseIfExpression)
	parser.registerPrefix(token.FUNCTION, parser.parseFunctionLiteral)

	parser.infixParseFns = make(map[token.TokenType]infixParseFn)

	parser.registerInfix(token.PLUS, parser.parseInfixExpression)
	parser.registerInfix(token.MINUS, parser.parseInfixExpression)
	parser.registerInfix(token.SLASH, parser.parseInfixExpression)
	parser.registerInfix(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfix(token.EQ, parser.parseInfixExpression)
	parser.registerInfix(token.NOT_EQ, parser.parseInfixExpression)
	parser.registerInfix(token.LT, parser.parseInfixExpression)
	parser.registerInfix(token.GT, parser.parseInfixExpression)

	// Read two tokens, so curToken and peekToken are both set
	parser.nextToken()
	parser.nextToken()
	return parser
}

func (self *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !self.currentTokenIs(token.EOF) {
		stmt := self.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		self.nextToken()
	}

	return program
}

func (self *Parser) parseStatement() ast.Statement {
	switch self.currentToken.Type {
	case token.LET:
		return self.parseLetStatement()
	case token.RETURN:
		return self.parseReturnStatement()
	default:
		return self.parseExpressionStatement()
	}
}

func (self *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: self.currentToken, Value: self.currentToken.Literal}
}

func (self *Parser) parseExpression(precedence int) ast.Expression {

	prefixFn := self.prefixParseFns[self.currentToken.Type]

	if prefixFn == nil {
		self.noPrefixParseFnError(self.currentToken.Type)
		return nil
	}

	leftExpression := prefixFn()

	for !self.peekTokenIs(token.SEMICOLON) && precedence < self.peekPrecedence() {
		infixFn := self.infixParseFns[self.peekToken.Type]
		if infixFn == nil {
			return leftExpression
		}

		self.nextToken()

		leftExpression = infixFn(leftExpression)
	}

	return leftExpression
}

func (self *Parser) parseExpressionStatement() *ast.ExpressionStatement {

	stmt := &ast.ExpressionStatement{Token: self.currentToken}
	stmt.Expression = self.parseExpression(LOWEST)

	if self.peekTokenIs(token.SEMICOLON) {
		self.nextToken()
	}

	return stmt
}

func (self *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: self.currentToken}

	if !self.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: self.currentToken, Value: self.currentToken.Literal}

	if !self.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO : we're skipping until we encounter a semi-colon
	for !self.currentTokenIs(token.SEMICOLON) {
		self.nextToken()
	}

	return stmt

}

func (self *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: self.currentToken}
	self.nextToken()
	// TODO: We're skipping the expressions until we
	// encounter a semicolon
	for !self.currentTokenIs(token.SEMICOLON) {
		self.nextToken()
	}
	return stmt
}

func (self *Parser) parseIntegerLiteral() ast.Expression {

	literal := &ast.IntegerLiteral{Token: self.currentToken}

	value, err := strconv.ParseInt(self.currentToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", self.currentToken.Literal)
		self.errors = append(self.errors, msg)
		return nil
	}

	literal.Value = value

	return literal
}

func (self *Parser) parsePrefixExpression() ast.Expression {

	expression := &ast.PrefixExpression{
		Token:    self.currentToken,
		Operator: self.currentToken.Literal,
	}

	self.nextToken()

	expression.Right = self.parseExpression(PREFIX)

	return expression
}

func (self *Parser) parseInfixExpression(left ast.Expression) ast.Expression {

	infixExpression := &ast.InfixExpression{
		Token:    self.currentToken,
		Operator: self.currentToken.Literal,
		Left:     left,
	}

	precedence := self.currentPrecedence()
	self.nextToken()
	infixExpression.Right = self.parseExpression(precedence)

	return infixExpression
}

func (self *Parser) parseGroupedExpression() ast.Expression {
	self.nextToken()

	expression := self.parseExpression(LOWEST)

	if !self.expectPeek(token.RPAREN) {
		return nil
	}

	return expression
}

func (self *Parser) parseBlockStatement() *ast.BlockStatement {

	block := &ast.BlockStatement{Token: self.currentToken}

	block.Statements = []ast.Statement{}

	self.nextToken()

	for !self.currentTokenIs(token.RBRACE) && !self.currentTokenIs(token.EOF) {
		stmt := self.parseStatement()

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		self.nextToken()
	}

	return block
}

func (self *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: self.currentToken}

	if !self.expectPeek(token.LPAREN) {
		return nil
	}

	self.nextToken()

	expression.Condition = self.parseExpression(LOWEST)

	if !self.expectPeek(token.RPAREN) {
		return nil
	}

	if !self.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = self.parseBlockStatement()

	if self.peekTokenIs(token.ELSE) {

		self.nextToken()

		if !self.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = self.parseBlockStatement()

	}

	return expression
}

func (self *Parser) parseFunctionLiteral() ast.Expression {
	functionLiteral := &ast.FunctionLiteral{Token: self.currentToken}

	if !self.expectPeek(token.LPAREN) {
		return nil
	}

	functionLiteral.Parameters = self.parseFunctionParameters()

	if !self.expectPeek(token.LBRACE) {
		return nil
	}

	functionLiteral.Body = self.parseBlockStatement()

	return functionLiteral
}

func (self *Parser) parseFunctionParameters() []*ast.Identifier {
	var identifiers []*ast.Identifier

	// no param function literall
	if self.peekTokenIs(token.RPAREN) {
		self.nextToken()
		return identifiers
	}

	self.nextToken()

	ident := &ast.Identifier{Token: self.currentToken, Value: self.currentToken.Literal}
	identifiers = append(identifiers, ident)

	for self.peekTokenIs(token.COMMA) {

		self.nextToken()
		self.nextToken()

		nextIdent := &ast.Identifier{Token: self.currentToken, Value: self.currentToken.Literal}
		identifiers = append(identifiers, nextIdent)

	}

	if !self.expectPeek(token.RPAREN) {
		return nil
	}
	return identifiers
}

func (self *Parser) noPrefixParseFnError(tok token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function found for %s found", tok)
	self.errors = append(self.errors, msg)
}

func (self *Parser) currentTokenIs(t token.TokenType) bool {
	return self.currentToken.Type == t
}

func (self *Parser) peekTokenIs(t token.TokenType) bool {
	return self.peekToken.Type == t
}

func (self *Parser) expectPeek(t token.TokenType) bool {
	if self.peekTokenIs(t) {
		self.nextToken()
		return true
	} else {
		self.peekErrors(t)
		return false
	}
}

func (self *Parser) Errors() []string {
	return self.errors
}

func (self *Parser) peekErrors(t token.TokenType) {
	message := fmt.Sprintf("expected next token to be %s, got %s instead", t, self.peekToken.Type)
	self.errors = append(self.errors, message)
}

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

func (self *Parser) peekPrecedence() int {
	if precedence, ok := precedences[self.peekToken.Type]; ok {
		return precedence
	}

	return LOWEST
}

func (self *Parser) currentPrecedence() int {
	if precedence, ok := precedences[self.currentToken.Type]; ok {
		return precedence
	}

	return LOWEST
}

func (self *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: self.currentToken, Value: self.currentTokenIs(token.TRUE)}
}
