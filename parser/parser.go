package parser

import (
	"github.com/Neal-C/interpreter-in-go/ast"
	"github.com/Neal-C/interpreter-in-go/lexer"
	"github.com/Neal-C/interpreter-in-go/token"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
}

func (self *Parser) nextToken() {
	self.currentToken = self.peekToken
	self.peekToken = self.lexer.NextToken()
}

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{lexer: lexer}

	return parser
}

func (self *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for self.currentToken.Type != token.EOF {
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
		return self.ParseLetStatement()
	default:
		return nil
	}
}

func (self *Parser) ParseLetStatement() *ast.LetStatement {
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
