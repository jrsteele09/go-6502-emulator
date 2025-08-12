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
	ByteDirectiveToken
	WordDirectiveToken
	TextDirectiveToken
	StringDirectiveToken
	StrDirectiveToken
	AscDirectiveToken
	AsciizDirectiveToken
	OrgDirectiveToken
	EquDirectiveToken
	DbDirectiveToken
	DwDirectiveToken
	DsDirectiveToken
	ProgramCounterToken
	ProgramCounterSetToken
)

// KeywordTokens defines keyword to token mappings
var KeywordTokens = map[string]lexer.TokenIdentifier{
	".BYTE":   ByteDirectiveToken,
	".WORD":   WordDirectiveToken,
	".TEXT":   TextDirectiveToken,
	".STRING": StringDirectiveToken,
	".STR":    StrDirectiveToken,
	".ASC":    AscDirectiveToken,
	".ASCIIZ": AsciizDirectiveToken,
	".ORG":    OrgDirectiveToken,
	".EQU":    EquDirectiveToken,
	".DB":     DbDirectiveToken,
	".DW":     DwDirectiveToken,
	".DS":     DsDirectiveToken,
	"*=":      ProgramCounterSetToken,
}

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
