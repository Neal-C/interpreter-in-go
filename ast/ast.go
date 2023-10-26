package ast

import "github.com/Neal-C/interpreter-in-go/token"

type Node interface {
	TokenLiteral() string
}

// marker interface$

type Statement interface {
	Node
	statementNode()
}

// marker interface

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (self *Program) TokenLiteral() string {
	if len(self.Statements) > 0 {
		return self.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type Identifier struct {
	Token token.Token // The token.IDENT token
	Value string
}

func (self *Identifier) expressionNode() {}

func (self *Identifier) TokenLiteral() string {
	return self.Token.Literal
}

type LetStatement struct {
	Token token.Token // token.LET token
	Name  *Identifier
	Value Expression
}

func (self *LetStatement) statementNode() {}

func (self *LetStatement) TokenLiteral() string {
	return self.Token.Literal
}
