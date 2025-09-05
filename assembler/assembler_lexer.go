package assembler

import (
	"bufio"
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
	'>': GreaterThanToken,
	'<': LessThanToken,
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

// NewAssemblerLexer creates a new preprocessor with the given file resolver
func NewAssemblerLexer(resolver utils.FileResolver) *AssemblerLexer {
	return &AssemblerLexer{
		fileResolver:    resolver,
		MaxIncludeDepth: 10, // Reasonable default for include depth
		includedFiles:   make(map[string]bool),
		includeCount:    make(map[string]int),
		importOnce:      make(map[string]bool),
	}
}

// Tokens tokenizes the input while expanding include directives, preserving
// SourceLine and SourceColumn as they appear in their original files.
// It tokenizes line-by-line so that line numbers remain accurate in the produced tokens.
func (p *AssemblerLexer) Tokens(cfg *lexer.LanguageConfig, input io.Reader, filename string) ([]*lexer.Token, error) {
	// Reset included files for each processing session
	p.includedFiles = make(map[string]bool)
	p.includeCount = make(map[string]int)
	p.importOnce = make(map[string]bool)

	tokens, err := p.readerTokens(cfg, input, filename, 0)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

// readerTokens recursively processes a reader, expanding includes and returning tokens
func (p *AssemblerLexer) readerTokens(cfg *lexer.LanguageConfig, input io.Reader, filename string, depth int) ([]*lexer.Token, error) {
	if depth > p.MaxIncludeDepth {
		return nil, fmt.Errorf("maximum include depth (%d) exceeded", p.MaxIncludeDepth)
	}

	var out []*lexer.Token
	scanner := bufio.NewScanner(input)
	lineNum := 0
	var sourceCode strings.Builder
	sourceLine := 1

	tokenizeSource := func(filename string) error {
		if sourceCode.Len() == 0 {
			return nil
		}
		lex := lexer.NewLexer(cfg)
		toks, err := lex.Tokenize(strings.NewReader(sourceCode.String()), filename)
		if err != nil {
			return fmt.Errorf("tokenize error around line %d: %w", sourceLine, err)
		}
		// Offset SourceLine to match original file lines and filter EOF
		lineOffset := uint(sourceLine - 1)
		for _, t := range toks {
			if t == nil || t.ID == lexer.EOFType {
				continue
			}
			if t.SourceLine > 0 {
				t.SourceLine += lineOffset
			} else {
				t.SourceLine = lineOffset + 1
			}
			out = append(out, t)
		}
		sourceCode.Reset()
		return nil
	}

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		// Check for preprocessor commands
		if strings.HasPrefix(trimmedLine, "#") {
			if skip, err := p.handlePreprocessorCommand(trimmedLine, filename); err != nil {
				return nil, err
			} else if skip {
				continue
			}
		}

		// Check for include directives
		includeFilePath := p.extractIncludePath(trimmedLine)
		if includeFilePath != "" {
			if p.importOnce[includeFilePath] && p.includeCount[includeFilePath] > 0 {
				// Skip this include as it has already been imported once
				continue
			}
			// Flush any accumulated non-include lines before processing include
			if err := tokenizeSource(filename); err != nil {
				return nil, err
			}

			// Prevent circular includes
			if p.includedFiles[includeFilePath] {
				return nil, fmt.Errorf("circular include detected: '%s' (line %d)", includeFilePath, lineNum)
			}

			p.includedFiles[includeFilePath] = true
			p.includeCount[includeFilePath]++

			includeFileReader, err := p.fileResolver.Resolve(includeFilePath)
			if err != nil {
				return nil, fmt.Errorf("line %d: %w", lineNum, err)
			}

			includedTokens, err := p.readerTokens(cfg, includeFileReader, includeFilePath, depth+1)
			if err != nil {
				return nil, fmt.Errorf("in file '%s': %w", includeFilePath, err)
			}

			out = append(out, includedTokens...)
			delete(p.includedFiles, includeFilePath)
			sourceLine = lineNum + 1
			continue
		}

		// Accumulate regular line into chunk so multi-line constructs tokenize correctly
		sourceCode.WriteString(line)
		sourceCode.WriteString("\n")
	}

	// Flush any remaining chunk
	if err := tokenizeSource(filename); err != nil {
		return nil, err
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	// Ensure a single EOF token terminates the stream
	if depth == 0 {
		out = append(out, lexer.NewToken(lexer.EOFType, "", 0))
	}

	return out, nil
}

func (p *AssemblerLexer) handlePreprocessorCommand(trimmedLine string, filename string) (skipLine bool, err error) {
	switch trimmedLine {
	case "#importonce":
		p.importOnce[filename] = true
		return true, nil
	}
	return false, nil
}

// extractIncludePath extracts the file path from include directives
// Supports both #include "file.asm" and .include "file.asm" formats
func (p *AssemblerLexer) extractIncludePath(line string) string {
	line = strings.TrimSpace(line)

	// Handle #include directive
	if strings.HasPrefix(line, "#include") {
		return p.extractQuotedPath(line[8:]) // Skip "#include"
	}

	// Handle .include directive
	if strings.HasPrefix(line, ".include") {
		return p.extractQuotedPath(line[8:]) // Skip ".include"
	}

	// Handle .INCLUDE directive (uppercase)
	if strings.HasPrefix(line, ".INCLUDE") {
		return p.extractQuotedPath(line[8:]) // Skip ".INCLUDE"
	}

	return ""
}

// extractQuotedPath extracts a quoted file path from the remaining part of an include line
func (p *AssemblerLexer) extractQuotedPath(remainder string) string {
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
