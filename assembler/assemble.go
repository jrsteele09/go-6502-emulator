package assembler

import (
	"fmt"
	"io"
	"strings"

	"github.com/jrsteele09/go-6502-emulator/cpu"
	"github.com/jrsteele09/go-6502-emulator/utils"
	"github.com/jrsteele09/go-lexer/lexer"
	"golang.org/x/exp/constraints"
)

const (
	twoByteOperand  = "nn"
	fourByteOperand = "nnnn"
)

type Instruction struct {
	Opcode     int
	Definition *cpu.OpCodeDef
}

type Assembler struct {
	instructionSet map[string]map[cpu.AddressingModeType]Instruction
	labels         map[string]uint64 // Labels address
	variables      map[string]interface{}
	lexerConfig    *lexer.LanguageConfig
}

type AssemblerTokens struct {
	tokens    []*lexer.Token
	tokenIdx  int
	currToken *lexer.Token
	nextToken *lexer.Token
}

func NewAssemblerTokens(tokens []*lexer.Token) *AssemblerTokens {
	at := &AssemblerTokens{
		tokens: tokens,
	}

	at.nextToken = tokens[0]
	at.tokenIdx = 0
	return at
}

func (at *AssemblerTokens) Next() *lexer.Token {
	at.currToken = at.nextToken
	if at.tokenIdx >= len(at.tokens)-1 {
		at.nextToken = nil
		return nil
	}
	at.tokenIdx++
	at.nextToken = at.tokens[at.tokenIdx]
	return at.currToken
}

func (at *AssemblerTokens) Peek() *lexer.Token {
	return at.nextToken
}

func New(opcodes []*cpu.OpCodeDef) *Assembler {
	instructionSet := make(map[string]map[cpu.AddressingModeType]Instruction)
	for opcode, opCodeDef := range opcodes {
		if opCodeDef == nil {
			continue
		}
		if _, found := instructionSet[opCodeDef.Mnemonic]; !found {
			instructionSet[opCodeDef.Mnemonic] = make(map[cpu.AddressingModeType]Instruction)
		}
		instructionSet[opCodeDef.Mnemonic][opCodeDef.AddressingModeType] = Instruction{
			Opcode:     opcode,
			Definition: opCodeDef,
		}
	}

	assembler := &Assembler{
		instructionSet: instructionSet,
		labels:         make(map[string]uint64),
		variables:      make(map[string]interface{}),
	}

	assembler.lexerConfig = lexer.NewLexerLanguage(
		lexer.WithKeywords(KeywordTokens),
		lexer.WithCustomTokenizers(customTokenizers),
		lexer.WithOperators(OperatorTokens),
		lexer.WithSymbols(SymbolTokens),
		lexer.WithCommentMap(comments),
		lexer.WithSpecializationCreators(assembler.identifierTokenCreator),
		lexer.WithExtendendedIdentifierRunes("_", ":"), // Allow underscores in identifiers, but when parsing an identifier, stop at a colon (Enables things like Labels)
	)

	return assembler
}

type ParsedInstruction struct {
	instruction Instruction
	data        []byte
}

func (a *Assembler) Assemble(r io.Reader) error {
	tokens, err := lexer.NewLexer(a.lexerConfig).Tokenize(r)
	if err != nil {
		return fmt.Errorf("[Assembler assemble] Tokenize [%w]", err)
	}

	asmTokens := NewAssemblerTokens(tokens)
	var parsedInstructions = []ParsedInstruction{}
	for {
		t := asmTokens.Next()
		if t == nil || t.ID == lexer.EOFType {
			break
		}

		switch t.ID {
		case MnemonicToken:
			addressingMode, err := a.parseAddressingMode(t.Literal, asmTokens)
			fmt.Printf("%v\n", addressingMode.AddressingMode)
			if err != nil {
				return err
			}
			instruction, ok := a.instructionSet[t.Literal][addressingMode.AddressingMode]
			if !ok {
				return fmt.Errorf("[Assembler Assemble] invalid addressing mode for instruction: %s %s", t.Literal, addressingMode)
			}

			parsedInstructions = append(parsedInstructions, ParsedInstruction{
				instruction: instruction,
				data:        addressingMode.Operands,
			})
		}
	}
	return nil
}

type AddressingMode struct {
	AddressingMode cpu.AddressingModeType
	Identifier     string
	Operands       []byte
}

func (a *Assembler) parseAddressingMode(mnemonic string, asmTokens *AssemblerTokens) (AddressingMode, error) {
	parsedAddressingMode := ""
	var operandValues []any
	var identifier string
parseLoop:
	for {
		t := asmTokens.Next()
		switch t.ID {
		case lexer.EndOfLineType, lexer.EOFType:
			break parseLoop
		case lexer.HexLiteral, lexer.IntegerLiteral:
			operandSizeMask, v, err := a.parseOperandSizeOfValue(false, t.Value)
			if err != nil {
				return AddressingMode{}, err
			}
			operandValues = append(operandValues, v)
			parsedAddressingMode += operandSizeMask

		case MinusToken:
			nextTokenID := utils.Value(asmTokens.Peek()).ID
			if nextTokenID == lexer.IntegerLiteral || nextTokenID == lexer.HexLiteral {
				negative := t.ID == MinusToken
				nt := asmTokens.Next()
				operandSizeMask, v, err := a.parseOperandSizeOfValue(negative, nt.Value)
				if err != nil {
					return AddressingMode{}, err
				}
				operandValues = append(operandValues, v)
				parsedAddressingMode += operandSizeMask
			} else {
				parsedAddressingMode += t.Literal
			}

		case IdentifierToken:
			identifier = t.Literal
			if value, ok := a.variables[identifier]; ok {
				operandSizeMask, v, err := a.parseOperandSizeOfValue(false, value)
				if err != nil {
					return AddressingMode{}, err
				}
				operandValues = append(operandValues, v)
				parsedAddressingMode += a.operandSizeForLabel(mnemonic, operandSizeMask)
				break
			}
			if address, ok := a.labels[identifier]; ok {
				operandSizeMask, v, err := a.parseOperandSizeOfValue(false, address)
				if err != nil {
					return AddressingMode{}, err
				}
				operandValues = append(operandValues, v)
				parsedAddressingMode += operandSizeMask
			} else {
				parsedAddressingMode += a.operandSizeForLabel(mnemonic, fourByteOperand) // Unresolved, assume a four byte label
			}

		default:
			parsedAddressingMode += t.Literal
		}
		if t.ID == lexer.EndOfLineType {
			break
		}
	}
	bytes, err := ValuesToLittleEndianBytes(operandValues)
	if err != nil {
		return AddressingMode{}, err
	}
	return AddressingMode{
		AddressingMode: cpu.AddressingModeType(parsedAddressingMode),
		Identifier:     identifier,
		Operands:       bytes,
	}, nil
}

func (a *Assembler) operandSizeForLabel(mnemonic, currentSizeStr string) string {
	if addressingModeTable, foundMnemonic := a.instructionSet[mnemonic]; foundMnemonic {
		if _, found := addressingModeTable[cpu.RelativeModeStr]; found {
			return twoByteOperand
		}
	}
	return currentSizeStr
}

func (a *Assembler) parseOperandSizeOfValue(negative bool, value any) (string, any, error) {
	switch v := value.(type) {
	case int8:
		return parseIntOperand(negative, v)
	case uint8:
		return parseIntOperand(negative, v)
	case int16:
		return parseIntOperand(negative, v)
	case uint16:
		return parseIntOperand(negative, v)
	case int32:
		return parseIntOperand(negative, v)
	case uint32:
		return parseIntOperand(negative, v)
	case int64:
		return parseIntOperand(negative, v)
	case uint64:
		return parseIntOperand(negative, v)
	default:
		return "", nil, fmt.Errorf("[Assembler parseSizeOfValue] expected valid operand %v", v)
	}
}

func parseIntOperand[T constraints.Integer](negative bool, value T) (string, any, error) {
	var signedValue T
	if value < 0 {
		signedValue = value
	} else {
		intValue := int64(value)
		if negative {
			intValue *= -1
		}
		if int64(T(intValue)) != intValue {
			return "", value, fmt.Errorf("[Assembler parseSizeOfValue] number too large %v", intValue)
		}
		signedValue = T(intValue)
	}

	operandSizeMask := strings.Replace(fmt.Sprintf("%x", signedValue), "-", "", 1)

	if len(operandSizeMask) > 4 {
		return "", signedValue, fmt.Errorf("[Assembler parseSizeOfValue] number too large %v", signedValue)
	}

	noOfBytes := 0
	if len(operandSizeMask) <= 2 {
		operandSizeMask = twoByteOperand
		noOfBytes = 1
	} else {
		operandSizeMask = fourByteOperand
		noOfBytes = 2
	}

	return operandSizeMask, ReduceBytes(signedValue, noOfBytes), nil
}

func (a *Assembler) mnemonicTokenCreator(identifier string) *lexer.Token {
	identifier = strings.ToUpper(identifier)
	if _, found := a.instructionSet[identifier]; !found {
		return nil
	}
	return lexer.NewToken(MnemonicToken, identifier, 0)
}

func labelTokenCreator(identifier string) *lexer.Token {
	if len(identifier) < 2 && !strings.HasSuffix(identifier, ":") {
		return nil
	}
	return lexer.NewToken(LabelToken, identifier, 0)
}

func (a *Assembler) identifierTokenCreator(identifier string) *lexer.Token {
	// Label
	if strings.HasSuffix(identifier, ":") {
		return labelTokenCreator(identifier)
	}

	// Mnemonic
	t := a.mnemonicTokenCreator(identifier)
	if t != nil {
		return t
	}

	// Identifier
	return lexer.NewToken(IdentifierToken, identifier, 0)
}
