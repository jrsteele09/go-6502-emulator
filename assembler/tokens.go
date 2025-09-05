package assembler

import (
	"github.com/jrsteele09/go-lexer/lexer"
)

// Tokens wraps a slice of tokens and provides methods for sequential processing
type Tokens struct {
	tokens    []*lexer.Token
	tokenIdx  int
	currToken *lexer.Token
	nextToken *lexer.Token
}

// NewAssemblerTokens creates a new AssemblerTokens instance with the given tokens
func NewAssemblerTokens(tokens []*lexer.Token) *Tokens {
	at := &Tokens{
		tokens: tokens,
	}

	at.nextToken = tokens[0]
	at.tokenIdx = 0
	return at
}

// Next advances to the next token and returns the current token
func (at *Tokens) Next() *lexer.Token {
	at.currToken = at.nextToken
	if at.tokenIdx >= len(at.tokens)-1 {
		at.nextToken = nil
		return nil
	}
	at.tokenIdx++
	at.nextToken = at.tokens[at.tokenIdx]
	return at.currToken
}

// Peek returns the next token without advancing the position
func (at *Tokens) Peek() *lexer.Token {
	return at.nextToken
}

func (at *Tokens) Current() *lexer.Token {
	return at.currToken
}
