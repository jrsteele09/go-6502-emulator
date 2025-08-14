package assembler

import (
	"fmt"
	"io"
	"strings"

	"github.com/jrsteele09/go-6502-emulator/cpu"
	"github.com/jrsteele09/go-6502-emulator/utils"
	"github.com/jrsteele09/go-lexer/lexer"
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
	constants      map[string]any
	lexerConfig    *lexer.LanguageConfig
	programCounter uint16
	originAddress  uint16
}

type Directive struct {
	Type string
	Args []interface{}
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
		constants:      make(map[string]interface{}),
		programCounter: 0x0000,
		originAddress:  0x0000,
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

type AssembledData struct {
	StartAddress uint16
	Data         []byte
}

func (a *Assembler) Assemble(r io.Reader) ([]AssembledData, error) {
	// Reset assembler state for each assembly
	a.reset()

	tokens, err := lexer.NewLexer(a.lexerConfig).Tokenize(r)
	if err != nil {
		return nil, fmt.Errorf("[Assembler assemble] Tokenize [%w]", err)
	}

	// First pass: calculate memory layout and collect labels
	segments, err := a.preprocessor(tokens)
	if err != nil {
		return nil, fmt.Errorf("[Assembler assemble] preprocessor [%w]", err)
	}

	// Second pass: generate machine code
	err = a.generateCode(tokens, segments)
	if err != nil {
		return nil, fmt.Errorf("[Assembler assemble] generateCode [%w]", err)
	}

	return segments, nil
}

// AssembleWithPreprocessor assembles source code with include directive support
func (a *Assembler) AssembleWithPreprocessor(r io.Reader, fileResolver FileResolver) ([]AssembledData, error) {
	// Reset assembler state for each assembly
	a.reset()

	// Preprocess to handle includes
	preprocessor := NewPreprocessor(fileResolver)
	processedInput, err := preprocessor.Process(r)
	if err != nil {
		return nil, fmt.Errorf("[Assembler AssembleWithPreprocessor] Preprocess [%w]", err)
	}

	tokens, err := lexer.NewLexer(a.lexerConfig).Tokenize(processedInput)
	if err != nil {
		return nil, fmt.Errorf("[Assembler AssembleWithPreprocessor] Tokenize [%w]", err)
	}

	// First pass: calculate memory layout and collect labels
	segments, err := a.preprocessor(tokens)
	if err != nil {
		return nil, fmt.Errorf("[Assembler AssembleWithPreprocessor] preprocessor [%w]", err)
	}

	// Second pass: generate machine code
	err = a.generateCode(tokens, segments)
	if err != nil {
		return nil, fmt.Errorf("[Assembler AssembleWithPreprocessor] generateCode [%w]", err)
	}

	return segments, nil
}

// reset clears labels and variables for a fresh assembly
func (a *Assembler) reset() {
	a.labels = make(map[string]uint64)
	a.constants = make(map[string]interface{})
	a.programCounter = 0x0000
	a.originAddress = 0x0000
}

func (a *Assembler) generateCode(tokens []*lexer.Token, segments []AssembledData) error {
	// Create a map for quick segment lookup by address -> segment index
	segmentMap := make(map[uint16]int)
	for i, segment := range segments {
		segmentMap[segment.StartAddress] = i
	}

	a.programCounter = a.originAddress
	currentSegmentIndex := -1
	segmentOffset := 0

	// Find the initial segment
	if segmentIdx, exists := segmentMap[a.programCounter]; exists {
		currentSegmentIndex = segmentIdx
		segmentOffset = 0
	}

	insertIntoMemory := func(data []byte) {
		if currentSegmentIndex >= 0 && currentSegmentIndex < len(segments) {
			copy(segments[currentSegmentIndex].Data[segmentOffset:], data)
			segmentOffset += len(data)
			a.programCounter += uint16(len(data))
		}
	}

	// Function to update current segment when program counter changes
	updateCurrentSegment := func() {
		if segmentIdx, exists := segmentMap[a.programCounter]; exists {
			if currentSegmentIndex != segmentIdx {
				currentSegmentIndex = segmentIdx
				segmentOffset = 0
			}
		} else {
			currentSegmentIndex = -1 // No segment found for this address
		}
	}

	asmTokens := NewAssemblerTokens(tokens)

	for {
		t := asmTokens.Next()
		if t == nil || t.ID == lexer.EOFType {
			break
		}

		switch t.ID {
		case AsterixSymbolToken:
			err := a.checkForOrgAsterixDirective(asmTokens, updateCurrentSegment)
			if err != nil {
				return err
			}

		case PeriodToken:
			err := a.processAssemblerDirective(asmTokens, insertIntoMemory, updateCurrentSegment)
			if err != nil {
				return err
			}

		case LabelToken:
			// Labels already processed in first pass
			continue

		case MnemonicToken:
			err := a.generateInstructionCode(t, asmTokens, insertIntoMemory)
			if err != nil {
				return err
			}

		case IdentifierToken:
			// Check if this identifier would be a valid mnemonic if it were uppercase
			upperIdent := strings.ToUpper(t.Literal)
			if _, found := a.instructionSet[upperIdent]; found {
				return fmt.Errorf("[generateCode] invalid mnemonic case '%s' (did you mean '%s'?)", t.Literal, upperIdent)
			}
			return fmt.Errorf("[generateCode] unknown mnemonic or invalid identifier '%s'", t.Literal)

		default:
			// Skip other tokens (comments, whitespace, etc.)
			continue
		}
	}
	return nil
}

func (a *Assembler) addressForAsterixOrgDirective(asmTokens *AssemblerTokens, finalizeSegment func()) error {
	// Check if this is "*=" (program counter set)
	nextToken := asmTokens.Peek()
	if nextToken != nil && nextToken.ID == EqualsSymbolToken {
		// Consume the equals token
		asmTokens.Next()
		// Process as program counter set (same as .ORG)
		err := a.processOrgDirective(asmTokens, finalizeSegment)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Assembler) checkForOrgAsterixDirective(asmTokens *AssemblerTokens, updateCurrentSegment func()) error {
	// Check if this is "*=" (program counter set)
	nextToken := asmTokens.Peek()
	if nextToken != nil && nextToken.ID == EqualsSymbolToken {
		// Consume the equals token
		asmTokens.Next()
		// Process program counter change using the same logic as .ORG
		t := asmTokens.Next()
		if t != nil {
			newAddress, err := a.tokenAddressValue(t)
			if err != nil {
				return err
			}
			a.programCounter = newAddress
			// Update segment after program counter change
			updateCurrentSegment()
		}
		// Skip any remaining tokens on this line
		a.skipDirectiveTokens(asmTokens)
	}
	return nil
}

func (a *Assembler) processAssemblerDirective(asmTokens *AssemblerTokens, insertIntoMemory func([]byte), updateCurrentSegment func()) error {
	// Check if this is a directive (. followed by directive name)
	nextToken := asmTokens.Peek()
	if nextToken != nil && nextToken.ID == IdentifierToken {
		directiveName := "." + nextToken.Literal
		if tokenID, found := KeywordTokens[directiveName]; found {
			// Consume the directive name token
			asmTokens.Next()
			// Process the specific directive based on its token ID (second pass)
			switch tokenID {
			case ByteDirectiveToken, DbDirectiveToken:
				err := a.processByteDirective(asmTokens, insertIntoMemory)
				if err != nil {
					return err
				}
			case WordDirectiveToken, DwDirectiveToken:
				err := a.processWordDirective(asmTokens, insertIntoMemory)
				if err != nil {
					return err
				}
			case TextDirectiveToken, StringDirectiveToken, StrDirectiveToken, AscDirectiveToken:
				err := a.processTextDirective(asmTokens, insertIntoMemory)
				if err != nil {
					return err
				}
			case AsciizDirectiveToken:
				err := a.processAsciizDirective(asmTokens, insertIntoMemory)
				if err != nil {
					return err
				}
			case OrgDirectiveToken:
				err := a.generateCodeForOrgDirective(asmTokens, updateCurrentSegment)
				if err != nil {
					return err
				}
			case EquDirectiveToken:
				a.skipDirectiveTokens(asmTokens) // Skip in second pass, already processed
			case DsDirectiveToken:
				err := a.processDataSpaceDirective(asmTokens, insertIntoMemory)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (a *Assembler) generateCodeForOrgDirective(asmTokens *AssemblerTokens, updateCurrentSegment func()) error {
	// Process program counter change - same as *= but without equals sign
	t := asmTokens.Next()
	if t != nil {
		newAddress, err := a.tokenAddressValue(t)
		if err != nil {
			return err
		}
		a.programCounter = newAddress
		// Update segment after program counter change, just like *= does
		updateCurrentSegment()
	}
	// Skip any remaining tokens on this line
	a.skipDirectiveTokens(asmTokens)
	return nil
}

func (a *Assembler) generateInstructionCode(t *lexer.Token, asmTokens *AssemblerTokens, insertIntoMemory func([]byte)) error {
	addressingMode, err := a.parseAddressingMode(t.Literal, asmTokens, false)
	if err != nil {
		return err
	}
	instruction, ok := a.instructionSet[t.Literal][addressingMode.AddressingMode]
	if !ok {
		return fmt.Errorf("[Assembler Assemble] invalid addressing mode for instruction: %s %s", t.Literal, addressingMode)
	}

	data := []byte{byte(instruction.Opcode)}
	data = append(data, addressingMode.Operands...)
	insertIntoMemory(data)
	return nil
}

type AddressingMode struct {
	AddressingMode cpu.AddressingModeType
	Identifier     string
	Operands       []byte
}

func (a *Assembler) parseAddressingMode(mnemonic string, asmTokens *AssemblerTokens, preprocess bool) (AddressingMode, error) {
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
			operandSizeMask, v, err := parseOperandSize(false, t.Value)
			if err != nil {
				return AddressingMode{}, err
			}
			operandValues = append(operandValues, v)
			parsedAddressingMode += operandSizeMask

		case MinusToken:
			nextTokenID := utils.Value(asmTokens.Peek()).ID
			if nextTokenID == lexer.IntegerLiteral || nextTokenID == lexer.HexLiteral {
				nt := asmTokens.Next()
				operandSizeMask, v, err := parseOperandSize(true, nt.Value)
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
			if value, ok := a.constants[identifier]; ok {
				operandSizeMask, v, err := parseOperandSize(false, value)
				if err != nil {
					return AddressingMode{}, err
				}
				operandValues = append(operandValues, v)
				parsedAddressingMode += a.operandSizeForLabel(mnemonic, operandSizeMask)
				break
			}
			if address, ok := a.labels[identifier]; ok && !preprocess {
				operandSizeMask, v, err := a.parseLabelOffset(mnemonic, address)
				if err != nil {
					return AddressingMode{}, err
				}
				operandValues = append(operandValues, v)
				parsedAddressingMode += operandSizeMask
				break
			}
			if !preprocess {
				return AddressingMode{}, fmt.Errorf("[parseAddressingMode] identifier '%s' not found", identifier)
			}
			operandSizeMask, v, err := a.preprocessorLabelSizer(mnemonic)
			if err != nil {
				return AddressingMode{}, fmt.Errorf("[parseAddressingMode] preprocessor label sizing failed: %w", err)
			}
			operandValues = append(operandValues, v)
			parsedAddressingMode += operandSizeMask
			// If preprocesing then we need to make an assumption that this is referencing a labal
			// That hasn't been defined yet, so assumptions about the addressing mode needs to be made
			// Depending on the mnemonic

		case AsterixSymbolToken:
			// Check what comes after the asterisk
			signToken := asmTokens.Peek()

			if signToken != nil && signToken.ID == PlusToken || signToken.ID == MinusToken {

				asmTokens.Next() // Get passed the signToken
				valueToken := asmTokens.Next()
				if valueToken == nil || (valueToken.ID != lexer.HexLiteral && valueToken.ID != lexer.IntegerLiteral) {
					return AddressingMode{}, fmt.Errorf("[parseAddressingMode] expected value after %s", signToken.Literal)
				}

				operandSizeMask, v, err := parseOperandSize(signToken.ID == MinusToken, valueToken.Value)
				if err != nil {
					return AddressingMode{}, err
				}
				if len(operandSizeMask) != 2 {
					return AddressingMode{}, fmt.Errorf("[parseAddressingMode] parsing relative address expected 8 bit signed")
				}

				operandValues = append(operandValues, v)
				parsedAddressingMode += string(t.Literal)
				parsedAddressingMode += string(operandSizeMask)
			} else {
				operandSizeMask, v, err := parseOperandSize(false, a.programCounter)
				if err != nil {
					return AddressingMode{}, fmt.Errorf("[parseAddressingMode] parsing * addressing mode %w", err)
				}

				operandValues = append(operandValues, v)
				parsedAddressingMode += operandSizeMask
			}

		case ProgramCounterToken:
			// Program counter reference - check if followed by operator
			nextToken := asmTokens.Peek()
			var finalPcValue uint64

			if nextToken != nil && (nextToken.ID == PlusToken || nextToken.ID == MinusToken) {
				// Handle expressions like "* + 5" or "* - 3"
				operatorToken := asmTokens.Next()
				valueToken := asmTokens.Next()
				if valueToken == nil || (valueToken.ID != lexer.HexLiteral && valueToken.ID != lexer.IntegerLiteral) {
					return AddressingMode{}, fmt.Errorf("[parseAddressingMode] expected value after %s", operatorToken.Literal)
				}

				pcValue := int64(a.programCounter)
				offsetValue, err := toUint64(valueToken.Value)
				if err != nil {
					return AddressingMode{}, fmt.Errorf("[parseAddressingMode] invalid offset value: %w", err)
				}

				if operatorToken.ID == PlusToken {
					pcValue += int64(offsetValue)
				} else {
					pcValue -= int64(offsetValue)
				}
				finalPcValue = uint64(pcValue)
			} else {
				// Simple program counter reference
				finalPcValue = uint64(a.programCounter)
			}

			operandSizeMask, v, err := parseOperandSize(false, finalPcValue)
			if err != nil {
				return AddressingMode{}, err
			}
			operandValues = append(operandValues, v)
			parsedAddressingMode += operandSizeMask

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
			return "*nn"
		}
	}
	return currentSizeStr
}

// parseOperandSize converts any integer type to operand size mask and value
// Returns "nn" for 8-bit values, "nnnn" for 16-bit values, etc.
// Handles negative flag by promoting to larger size if needed

func (a *Assembler) parseLabelOffset(mnemonic string, address uint64) (string, any, error) {
	addressingModes, ok := a.instructionSet[mnemonic]
	if !ok {
		return "", nil, fmt.Errorf("[Assembler parseLabelOffset] unknown mnemonic '%s'", mnemonic)
	}
	if _, found := addressingModes[cpu.RelativeModeStr]; found {
		// Calculate relative displacement: target - (PC + 2)
		delta := int64(address) - int64(a.programCounter+2)
		if delta < -128 || delta > 127 {
			return "", nil, fmt.Errorf("[Assembler parseLabelOffset] relative address out of range: %d", delta)
		}
		return string(cpu.RelativeModeStr), ReduceBytes(delta, 1), nil
	}

	return fourByteOperand, ReduceBytes(address, 2), nil
}

// parseLabelOffset

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

func (a *Assembler) processOrgDirective(asmTokens *AssemblerTokens, finalizeSegment func()) error {
	finalizeSegment() // Close current segment

	t := asmTokens.Next()
	if t == nil {
		return fmt.Errorf("[processOrgDirective] expected address after .ORG")
	}

	address, err := a.tokenAddressValue(t)
	if err != nil {
		return fmt.Errorf("[processOrgDirective] %w", err)
	}

	a.programCounter = address
	a.originAddress = address

	return nil
}

// tokenAddressValue extracts an address value from a token
func (a *Assembler) tokenAddressValue(t *lexer.Token) (uint16, error) {
	switch t.ID {
	case lexer.HexLiteral, lexer.IntegerLiteral:
		value, err := toUint64(t.Value)
		if err != nil {
			return 0, fmt.Errorf("invalid address value: %w", err)
		}
		return uint16(value), nil
	case ProgramCounterToken:
		return a.programCounter, nil
	default:
		return 0, fmt.Errorf("expected address value, got %s", t.Literal)
	}
}

func (a *Assembler) processByteDirective(asmTokens *AssemblerTokens, insertIntoMemory func([]byte)) error {
	var bytes []byte

	for {
		t := asmTokens.Next()
		if t == nil || t.ID == lexer.EndOfLineType || t.ID == lexer.EOFType {
			break
		}

		if t.ID == CommaToken {
			continue
		}

		switch t.ID {
		case lexer.HexLiteral, lexer.IntegerLiteral:
			value, err := toUint64(t.Value)
			if err != nil {
				return fmt.Errorf("[processByteDirective] invalid byte value: %w", err)
			}
			if value > 255 {
				return fmt.Errorf("[processByteDirective] byte value %d exceeds 255", value)
			}
			bytes = append(bytes, byte(value))
		case IdentifierToken:
			if value, ok := a.constants[t.Literal]; ok {
				if byteVal, ok := value.(uint8); ok {
					bytes = append(bytes, byteVal)
				} else {
					return fmt.Errorf("[processByteDirective] variable %s is not a byte", t.Literal)
				}
			} else {
				return fmt.Errorf("[processByteDirective] undefined variable: %s", t.Literal)
			}
		default:
			return fmt.Errorf("[processByteDirective] unexpected token: %s", t.Literal)
		}
	}

	if len(bytes) > 0 {
		insertIntoMemory(bytes)
	}
	return nil
}

func (a *Assembler) processWordDirective(asmTokens *AssemblerTokens, insertIntoMemory func([]byte)) error {
	var bytes []byte

	for {
		t := asmTokens.Next()
		if t == nil || t.ID == lexer.EndOfLineType || t.ID == lexer.EOFType {
			break
		}

		if t.ID == CommaToken {
			continue
		}

		switch t.ID {
		case lexer.HexLiteral, lexer.IntegerLiteral:
			value, err := toUint64(t.Value)
			if err != nil {
				return fmt.Errorf("[processWordDirective] invalid word value: %w", err)
			}
			if value > 65535 {
				return fmt.Errorf("[processWordDirective] word value %d exceeds 65535", value)
			}
			// Store in little-endian format
			bytes = append(bytes, byte(value&0xFF), byte((value>>8)&0xFF))
		case IdentifierToken:
			if address, ok := a.labels[t.Literal]; ok {
				// Store address in little-endian format
				bytes = append(bytes, byte(address&0xFF), byte((address>>8)&0xFF))
			} else if value, ok := a.constants[t.Literal]; ok {
				if wordVal, ok := value.(uint16); ok {
					bytes = append(bytes, byte(wordVal&0xFF), byte((wordVal>>8)&0xFF))
				} else {
					return fmt.Errorf("[processWordDirective] variable %s is not a word", t.Literal)
				}
			} else {
				return fmt.Errorf("[processWordDirective] undefined label/variable: %s", t.Literal)
			}
		default:
			return fmt.Errorf("[processWordDirective] unexpected token: %s", t.Literal)
		}
	}

	if len(bytes) > 0 {
		insertIntoMemory(bytes)
	}
	return nil
}

func (a *Assembler) processTextDirective(asmTokens *AssemblerTokens, insertIntoMemory func([]byte)) error {
	t := asmTokens.Next()
	if t == nil || t.ID != lexer.StringLiteral {
		return fmt.Errorf("[processTextDirective] expected string after .TEXT")
	}

	str, ok := t.Value.(string)
	if !ok {
		return fmt.Errorf("[processTextDirective] invalid string value")
	}

	bytes := []byte(str)
	insertIntoMemory(bytes)
	return nil
}

func (a *Assembler) processAsciizDirective(asmTokens *AssemblerTokens, insertIntoMemory func([]byte)) error {
	t := asmTokens.Next()
	if t == nil || t.ID != lexer.StringLiteral {
		return fmt.Errorf("[processAsciizDirective] expected string after .ASCIIZ")
	}

	str, ok := t.Value.(string)
	if !ok {
		return fmt.Errorf("[processAsciizDirective] invalid string value")
	}

	// Add null terminator to the string
	bytes := []byte(str)
	bytes = append(bytes, 0) // Add null terminator
	insertIntoMemory(bytes)
	return nil
}

func (a *Assembler) processEquDirective(asmTokens *AssemblerTokens) error {
	// Get variable name
	nameToken := asmTokens.Next()
	if nameToken == nil || nameToken.ID != IdentifierToken {
		return fmt.Errorf("[processEquDirective] expected identifier after .EQU")
	}

	variableName := nameToken.Literal

	// Check for duplicate variable
	if _, exists := a.constants[variableName]; exists {
		return fmt.Errorf("[processEquDirective] duplicate variable '%s' already defined", variableName)
	}

	// Check if this variable name conflicts with an existing label
	if _, exists := a.labels[variableName]; exists {
		return fmt.Errorf("[processEquDirective] variable '%s' conflicts with existing label", variableName)
	}

	// Skip equals sign if present
	nextToken := asmTokens.Peek()
	if nextToken != nil && nextToken.ID == EqualsSymbolToken {
		asmTokens.Next()
	}

	// Get value
	valueToken := asmTokens.Next()
	if valueToken == nil {
		return fmt.Errorf("[processEquDirective] expected value after .EQU")
	}

	switch valueToken.ID {
	case lexer.HexLiteral, lexer.IntegerLiteral:
		a.constants[variableName] = valueToken.Value
	case ProgramCounterToken:
		a.constants[variableName] = uint16(a.programCounter)
	default:
		return fmt.Errorf("[processEquDirective] invalid value type for .EQU")
	}

	return nil
}

func (a *Assembler) processDataSpaceDirective(asmTokens *AssemblerTokens, insertIntoMemory func([]byte)) error {
	t := asmTokens.Next()
	if t == nil {
		return fmt.Errorf("[processDataSpaceDirective] expected size after .DS")
	}

	var size uint64
	switch v := t.Value.(type) {
	case uint8:
		size = uint64(v)
	case uint16:
		size = uint64(v)
	case uint64:
		size = v
	case int:
		size = uint64(v)
	case int64:
		size = uint64(v)
	default:
		return fmt.Errorf("[processDataSpaceDirective] invalid size value type: %T", v)
	}

	// Create zero-filled bytes
	bytes := make([]byte, size)
	insertIntoMemory(bytes)
	return nil
}

func (a *Assembler) skipDirectiveTokens(asmTokens *AssemblerTokens) {
	// Skip tokens until end of line
	for {
		t := asmTokens.Next()
		if t == nil || t.ID == lexer.EndOfLineType || t.ID == lexer.EOFType {
			break
		}
	}
}

func (a *Assembler) identifierTokenCreator(identifier string) *lexer.Token {
	// Program Counter
	if identifier == "*" {
		return lexer.NewToken(ProgramCounterToken, identifier, 0)
	}

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
