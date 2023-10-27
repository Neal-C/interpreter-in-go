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
	return nil
}
