package parser

import (
	"fmt"
	"github.com/Neal-C/interpreter-in-go/ast"
	"github.com/Neal-C/interpreter-in-go/lexer"
	"github.com/Neal-C/interpreter-in-go/token"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
	errors       []string
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
		return nil
	}
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
	if !self.currentTokenIs(token.SEMICOLON) {
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
