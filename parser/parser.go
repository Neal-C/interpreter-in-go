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
	//// myFunction(x)
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

	// Read two tokens, so curToken and peekToken are both set
	parser.nextToken()
	parser.nextToken()
	return parser
}

func (self *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !self.currentTokenIs(token.EOF) {
		stmt := self.ParseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		self.nextToken()
	}

	return program
}

func (self *Parser) ParseStatement() ast.Statement {
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
	prefix := self.prefixParseFns[self.currentToken.Type]

	if prefix == nil {
		self.noPrefixParseFnError(self.currentToken.Type)
		return nil
	}

	leftExpression := prefix()

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

func (self *Parser) noPrefixParseFnError(tok token.TokenType) {
	msg := fmt.Sprintf("not prefix parse function found for %s found", tok)
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
