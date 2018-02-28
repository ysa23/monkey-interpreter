package lexer

import "src/token"

type Lexer struct {
	input           string
	currentPosition int
	nextPosition    int
	ch              byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.nextPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.nextPosition]
	}

	l.currentPosition = l.nextPosition
	l.nextPosition++
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespaces()

	switch l.ch {
	case '=':
		tok = l.extractEqualityOperatorsToken(l.ch, token.EQ, token.ASSIGN)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '!':
		tok = l.extractEqualityOperatorsToken(l.ch, token.NOT_EQ, token.BANG)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)

	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)

	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
			return tok
		}
		if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		}

		tok = newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()
	return tok
}

func (l *Lexer) peekChar() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	}

	return l.input[l.nextPosition]
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) extractEqualityOperatorsToken(ch byte, tokenType token.TokenType, defaultTokenType token.TokenType) token.Token {
	if l.peekChar() == '=' {
		ch := l.ch
		l.readChar()
		return token.Token{Type: tokenType, Literal: string(ch) + string(l.ch)}
	}

	return newToken(defaultTokenType, l.ch)
}

func (l *Lexer) readIdentifier() string {
	position := l.currentPosition
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.currentPosition]
}

func isLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string {
	position := l.currentPosition
	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.currentPosition]
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) skipWhitespaces() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}
