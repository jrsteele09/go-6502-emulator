package assembler

import "github.com/jrsteele09/go-lexer/lexer"

const (
	AsterixSymbolToken lexer.TokenIdentifier = lexer.LastStdLiteral + iota
	EqualsSymbolToken
	LeftParenthesis
	RightParenthesis
	CommaToken
	PeriodToken
	LabelToken
	IdentifierToken
	HashToken
	MnemonicToken
	MinusToken
	PlusToken
	DivideSymbolToken
	SemiColonToken
)

// KeywordTokens defines keyword to token mappings
var KeywordTokens = map[string]lexer.TokenIdentifier{}

// Custom tokenizers - On detection of the starting character, jump to a specific tokenizer.
var customTokenizers = map[string]lexer.TokenizerFunc{
	"$": lexer.HexTokenizer,
	"%": lexer.BinaryTokenizer,
}

var OperatorTokens = map[string]lexer.TokenIdentifier{}

// SymbolTokens defines single delimeter runes to token mappings
var SymbolTokens = map[rune]lexer.TokenIdentifier{
	'*': AsterixSymbolToken,
	'=': EqualsSymbolToken,
	'(': LeftParenthesis,
	')': RightParenthesis,
	',': CommaToken,
	'.': PeriodToken,
	'#': HashToken,
	'-': MinusToken,
	'+': PlusToken,
	'/': DivideSymbolToken,
	';': SemiColonToken,
}

// comments defines comment syntax mappings
var comments = map[string]string{
	";":  "\n",
	"//": "\n",
	"/*": "*/",
}
