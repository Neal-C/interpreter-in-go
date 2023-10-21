package lexer

import "github.com/Neal-C/interpreter-in-go/token"

type Lexer struct {
	input        string
	position     int  // current position in the input (points to current character) // and where we last read
	readPosition int  // current reading position in the input (after current character)
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	lexer := &Lexer{
		input: input,
	}

	lexer.readChar()

	return lexer
}

func (lexer *Lexer) readChar() {
	//  is to check whether we have reached the end of input.
	if lexer.readPosition >= len(lexer.input) {
		// 0 is the ASCII code for 'NUL' character
		lexer.ch = 0
		// signifies either
		// “we haven’t read anything yet” or “end of file” for us.
	} else {
		// But if we haven’t reached the end of input yet,
		// it sets l.ch to the next character by accessing l.input[l.readPosition].
		lexer.ch = lexer.input[lexer.readPosition]
	}

	lexer.position = lexer.readPosition
	lexer.readPosition++
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

func (lexer *Lexer) NextToken() token.Token {
	var tok token.Token

	switch lexer.ch {
	case '=':
		tok = newToken(token.ASSIGN, lexer.ch)
	case ';':
		tok = newToken(token.SEMICOLON, lexer.ch)
	case '(':
		tok = newToken(token.LPAREN, lexer.ch)
	case ')':
		tok = newToken(token.RPAREN, lexer.ch)
	case ',':
		tok = newToken(token.COMMA, lexer.ch)
	case '+':
		tok = newToken(token.PLUS, lexer.ch)
	case '{':
		tok = newToken(token.LBRACE, lexer.ch)
	case '}':
		tok = newToken(token.RBRACE, lexer.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(lexer.ch) {
			tok.Literal = lexer.readIdentifier()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, lexer.ch)
		}
	}

	lexer.readChar()
	return tok
}
