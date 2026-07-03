package assembler

import "github.com/jrsteele09/go-lexer/lexer"

func isTerminatorToken(tokenID lexer.TokenIdentifier) bool {
	return tokenID == lexer.NullType || tokenID == lexer.EndOfLineType || tokenID == lexer.EOFType
}
