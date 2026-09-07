package assembler

import (
	"fmt"
	"io"
	"strings"

	"github.com/jrsteele09/go-6502-emulator/utils"
	"github.com/jrsteele09/go-lexer/lexer"
)

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
	GreaterThanToken
	LessThanToken
	AmpersandToken
	PipeToken
	CaretToken
	TildeToken
	BangToken
	ShiftLeftToken
	ShiftRightToken
	NotEqualToken
	EqualEqualToken
	LessEqualToken
	GreaterEqualToken
)

// KeywordTokens defines keyword to token mappings
var KeywordTokens = map[string]lexer.TokenIdentifier{}

// Custom tokenizers - On detection of the starting character, jump to a specific tokenizer.
var prefixTokenizers = map[string]lexer.TokenizerFunc{
	"$": lexer.HexTokenizer,
	"%": lexer.BinaryTokenizer,
}

var OperatorTokens = map[string]lexer.TokenIdentifier{
	"<<": ShiftLeftToken,
	">>": ShiftRightToken,
	"!=": NotEqualToken,
	"==": EqualEqualToken,
	"<=": LessEqualToken,
	">=": GreaterEqualToken,
}

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
	'>': GreaterThanToken,
	'<': LessThanToken,
	'&': AmpersandToken,
	'|': PipeToken,
	'^': CaretToken,
	'~': TildeToken,
	'!': BangToken,
}

// comments defines comment syntax mappings
var comments = map[string]string{
	";":  "\n",
	"//": "\n",
	"/*": "*/",
}

// AssemblerLexer converts file(s) to a continuous stream of tokens
type AssemblerLexer struct {
	fileResolver    utils.FileResolver
	MaxIncludeDepth int
	includedFiles   map[string]bool // Track included files to prevent circular includes
	includeCount    map[string]int
	importOnce      map[string]bool
}

// NewAssemblerLexer creates a lexer with the given file resolver.
func NewAssemblerLexer(resolver utils.FileResolver) *AssemblerLexer {
	return &AssemblerLexer{
		fileResolver:    resolver,
		MaxIncludeDepth: 10, // Reasonable default for include depth
		includedFiles:   make(map[string]bool),
		includeCount:    make(map[string]int),
		importOnce:      make(map[string]bool),
	}
}

// Tokens preprocesses includes in source order so the entire translation unit
// shares one macro and constant context.
func (p *AssemblerLexer) Tokens(cfg *lexer.LanguageConfig, input io.Reader, filename string) ([]lexer.Token, error) {
	p.includedFiles = make(map[string]bool)
	p.includeCount = make(map[string]int)
	p.importOnce = make(map[string]bool)

	processed, err := p.preprocessReader(input, filename, 0, make(map[string]sourceMacro), make(map[string]int64))
	if err != nil {
		return nil, fmt.Errorf("source preprocessing: %w", err)
	}
	tokens, err := lexer.NewLexer(cfg).Tokenize(strings.NewReader(processed), filename)
	if err != nil {
		return nil, fmt.Errorf("tokenize error: %w", err)
	}
	return tokens, nil
}

func (p *AssemblerLexer) preprocessReader(input io.Reader, filename string, depth int, macros map[string]sourceMacro, constants map[string]int64) (string, error) {
	if depth > p.MaxIncludeDepth {
		return "", fmt.Errorf("maximum include depth (%d) exceeded", p.MaxIncludeDepth)
	}

	data, err := io.ReadAll(input)
	if err != nil {
		return "", err
	}

	include := func(includeFilePath string, lineNumber int) (string, error) {
		if p.importOnce[includeFilePath] && p.includeCount[includeFilePath] > 0 {
			return "", nil
		}
		if p.includedFiles[includeFilePath] {
			return "", fmt.Errorf("circular include detected: '%s' (line %d)", includeFilePath, lineNumber+1)
		}

		includeFileReader, err := p.fileResolver.Resolve(includeFilePath)
		if err != nil {
			return "", fmt.Errorf("line %d: %w", lineNumber+1, err)
		}
		p.includedFiles[includeFilePath] = true
		p.includeCount[includeFilePath]++
		includedSource, err := p.preprocessReader(includeFileReader, includeFilePath, depth+1, macros, constants)
		delete(p.includedFiles, includeFilePath)
		if err != nil {
			return "", fmt.Errorf("in file '%s': %w", includeFilePath, err)
		}
		return includedSource, nil
	}

	directive := func(name string, _ int) error {
		if name == "#importonce" {
			p.importOnce[filename] = true
		}
		return nil
	}

	return processSourceLines(strings.Split(string(data), "\n"), filename, macros, constants, 0, include, directive)
}

// extractIncludePath supports both #include and .include syntax.
func extractIncludePath(line string) string {
	line = strings.TrimSpace(line)
	words := strings.Fields(line)
	if len(words) < 2 || (!strings.EqualFold(words[0], "#include") && !strings.EqualFold(words[0], ".include")) {
		return ""
	}
	return extractQuotedPath(strings.TrimSpace(strings.TrimPrefix(line, words[0])))
}

func extractQuotedPath(remainder string) string {
	remainder = strings.TrimSpace(remainder)

	// Handle both single and double quotes
	if len(remainder) < 2 {
		return ""
	}

	// Check for double quotes
	if remainder[0] == '"' {
		if endPos := strings.Index(remainder[1:], `"`); endPos != -1 {
			return remainder[1 : endPos+1]
		}
	}

	// Check for single quotes
	if remainder[0] == '\'' {
		if endPos := strings.Index(remainder[1:], `'`); endPos != -1 {
			return remainder[1 : endPos+1]
		}
	}

	// Handle unquoted paths (space-delimited)
	parts := strings.Fields(remainder)
	if len(parts) > 0 {
		return parts[0]
	}

	return ""
}
