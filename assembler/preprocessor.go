package assembler

import (
	"fmt"
	"strings"

	"github.com/jrsteele09/go-6502-emulator/cpu"
	"github.com/jrsteele09/go-lexer/lexer"
)

// preprocessor performs the first pass of assembly: calculate memory layout and collect labels
func (a *Assembler) preprocessor(tokens []*lexer.Token) ([]AssembledData, error) {
	// Track memory segments that will be needed
	type SegmentInfo struct {
		StartAddress uint16
	}
	var segmentInfos []SegmentInfo
	var currentSegmentStart uint16

	a.programCounter = 0x0000
	currentSegmentStart = a.programCounter
	currentSegmentSize := 0

	finalizeCurrentSegment := func() {
		if currentSegmentSize > 0 {
			segmentInfos = append(segmentInfos, SegmentInfo{
				StartAddress: currentSegmentStart,
			})
		}
		currentSegmentSize = 0
	}

	advanceProgramCounter := func(size int) {
		a.programCounter += uint16(size)
		currentSegmentSize += size
	}

	asmTokens := NewAssemblerTokens(tokens)

	tokenPosition := 0

	for {
		tokenPosition++
		t := asmTokens.Next()
		if t == nil || t.ID == lexer.EOFType {
			break
		}

		switch t.ID {
		case AsterixSymbolToken:
			err := a.addressForAsterixOrgDirective(asmTokens, finalizeCurrentSegment)
			if err != nil {
				return nil, err
			}
			// Check if program counter changed (due to *=) and start new segment
			if a.programCounter != currentSegmentStart+uint16(currentSegmentSize) {
				currentSegmentStart = a.programCounter
				currentSegmentSize = 0
			}

		case PeriodToken:
			err := a.preprocessDirective(asmTokens, advanceProgramCounter, finalizeCurrentSegment)
			if err != nil {
				return nil, err
			}
			// Check if program counter changed (due to .ORG) and start new segment
			if a.programCounter != currentSegmentStart+uint16(currentSegmentSize) {
				currentSegmentStart = a.programCounter
				currentSegmentSize = 0
			}

		case LabelToken:
			err := a.recordLabelAddress(t)
			if err != nil {
				return nil, err
			}

		case MnemonicToken:
			// Calculate instruction size
			addressingMode, err := a.parseAddressingMode(t.Literal, asmTokens, true)
			if err != nil {
				return nil, err
			}
			instructionSize := 1 + len(addressingMode.Operands)
			advanceProgramCounter(instructionSize)
			tokenPosition = 0

		case IdentifierToken:
			// Check if this is a constant assignment (identifier = value)
			nextToken := asmTokens.Peek()
			if nextToken != nil && nextToken.ID == EqualsSymbolToken {
				err := a.processConstantAssignment(t, asmTokens)
				if err != nil {
					return nil, err
				}
			} else if tokenPosition == 1 {
				err := a.recordLabelAddress(t)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("[preprocessor] unexpected identifier '%s'", t.Literal)
			}
		case lexer.EndOfLineType:
			tokenPosition = 0 // Reset position on new line

		default:
			// Skip other tokens (comments, whitespace, etc.)
			continue
		}
	}

	finalizeCurrentSegment()

	// Create segements with start address
	segments := make([]AssembledData, len(segmentInfos))
	for i, info := range segmentInfos {
		segments[i] = AssembledData{
			StartAddress: info.StartAddress,
		}
	}

	return segments, nil
}

// preprocessDirective handles directive processing during the first pass
func (a *Assembler) preprocessDirective(asmTokens *Tokens, advanceProgramCounter func(int), finalizeSegment func()) error {
	// Check if this is a directive (. followed by directive name)
	nextToken := asmTokens.Peek()
	if nextToken != nil && nextToken.ID == IdentifierToken {
		directiveName := "." + nextToken.Literal
		if directiveID, found := a.directives[strings.ToUpper(directiveName)]; found {
			// Consume the directive name token
			asmTokens.Next()
			// Process the specific directive based on its token ID
			switch directiveID {
			case ByteDirective, DbDirective:
				size, err := a.calculateByteDirectiveSize(asmTokens)
				if err != nil {
					return err
				}
				advanceProgramCounter(size)
			case WordDirective, DwDirective:
				size, err := a.calculateWordDirectiveSize(asmTokens)
				if err != nil {
					return err
				}
				advanceProgramCounter(size)
			case TextDirective, StringDirective, StrDirective, AscDirective:
				size, err := a.calculateTextDirectiveSize(asmTokens)
				if err != nil {
					return err
				}
				advanceProgramCounter(size)
			case AsciizDirective:
				size, err := a.calculateAsciizDirectiveSize(asmTokens)
				if err != nil {
					return err
				}
				advanceProgramCounter(size)
			case OrgDirective:
				err := a.processOrgDirective(asmTokens, finalizeSegment)
				if err != nil {
					return err
				}
			case DsDirective:
				size, err := a.calculateDataSpaceDirectiveSize(asmTokens)
				if err != nil {
					return err
				}
				advanceProgramCounter(size)
			case VarDirective:
				err := a.processVarDirective(asmTokens)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// recordLabelAddress records the current program counter as the address for a label
func (a *Assembler) recordLabelAddress(t *lexer.Token) error {
	labelName := strings.TrimSuffix(t.Literal, ":")

	// Check for duplicate label
	if _, exists := a.labels[labelName]; exists {
		return fmt.Errorf("[recordLabelAddress] duplicate label '%s' already defined", labelName)
	}

	// Check if this label name conflicts with an existing variable
	if _, exists := a.constants[labelName]; exists {
		return fmt.Errorf("[recordLabelAddress] label '%s' conflicts with existing variable", labelName)
	}

	a.labels[labelName] = uint64(a.programCounter)
	return nil
}

// Size calculation functions for first pass
func (a *Assembler) calculateByteDirectiveSize(asmTokens *Tokens) (int, error) {
	size := 0
	for {
		t := asmTokens.Peek()
		if t.ID == CommaToken {
			asmTokens.Next() // Skip comma
			continue
		}
		if t.ID != lexer.HexLiteral && t.ID != lexer.IntegerLiteral && t.ID != IdentifierToken {
			break
		}
		size++
		asmTokens.Next() // Consume the token
	}
	return size, nil
}

func (a *Assembler) calculateWordDirectiveSize(asmTokens *Tokens) (int, error) {
	size := 0
	for {
		t := asmTokens.Peek()
		if t == nil || t.ID == lexer.EndOfLineType || t.ID == lexer.EOFType {
			break
		}
		asmTokens.Next() // Consume the token
		if t.ID == CommaToken {
			continue
		}
		if t.ID == lexer.HexLiteral || t.ID == lexer.IntegerLiteral || t.ID == IdentifierToken {
			size += 2 // Words are 2 bytes
		}
	}
	return size, nil
}

func (a *Assembler) calculateTextDirectiveSize(asmTokens *Tokens) (int, error) {
	t := asmTokens.Next()
	if t == nil || t.ID != lexer.StringLiteral {
		return 0, fmt.Errorf("[calculateTextDirectiveSize] expected string after .TEXT")
	}
	str, ok := t.Value.(string)
	if !ok {
		return 0, fmt.Errorf("[calculateTextDirectiveSize] invalid string value")
	}
	return len(str), nil
}

func (a *Assembler) calculateAsciizDirectiveSize(asmTokens *Tokens) (int, error) {
	t := asmTokens.Next()
	if t == nil || t.ID != lexer.StringLiteral {
		return 0, fmt.Errorf("[calculateAsciizDirectiveSize] expected string after .ASCIIZ")
	}
	str, ok := t.Value.(string)
	if !ok {
		return 0, fmt.Errorf("[calculateAsciizDirectiveSize] invalid string value")
	}
	return len(str) + 1, nil // +1 for null terminator
}

func (a *Assembler) calculateDataSpaceDirectiveSize(asmTokens *Tokens) (int, error) {
	t := asmTokens.Next()
	if t == nil {
		return 0, fmt.Errorf("[calculateDataSpaceDirectiveSize] expected size after .DS")
	}

	size, err := toUint64(t.Value)
	if err != nil {
		return 0, fmt.Errorf("[calculateDataSpaceDirectiveSize] %w", err)
	}

	return int(size), nil
}

// preprocessorLabelSizer determines operand size for forward-referenced labels during preprocessing
// It makes assumptions based on the mnemonic and whether it has Relative addressing modes
func (a *Assembler) preprocessorLabelSizer(mnemonic string) (string, any, error) {
	// If Non Mnemonic, assume a label address and not a relative address calculation
	if mnemonic == "" {
		return twoByteOperand, ReduceBytes(0x0, 2), nil
	}
	addressingModes, ok := a.instructionSet[mnemonic]
	if !ok {
		return "", nil, fmt.Errorf("[Assembler preprocessorLabelSizer] unknown mnemonic '%s'", mnemonic)
	}
	if _, found := addressingModes[cpu.RelativeModeStr]; found {
		return oneByteOperand, ReduceBytes(0x0, 1), nil
	}

	return twoByteOperand, ReduceBytes(0x0, 2), nil
}

// processConstantAssignment handles identifier = value assignments during preprocessing
func (a *Assembler) processConstantAssignment(identifierToken *lexer.Token, asmTokens *Tokens) error {
	variableName := identifierToken.Literal

	// Check for duplicate variable
	if _, exists := a.constants[variableName]; exists {
		return fmt.Errorf("[processConstantAssignment] duplicate variable '%s' already defined", variableName)
	}

	// Check if this variable name conflicts with an existing label
	if _, exists := a.labels[variableName]; exists {
		return fmt.Errorf("[processConstantAssignment] variable '%s' conflicts with existing label", variableName)
	}

	// Consume the equals token
	equalsToken := asmTokens.Next()
	if equalsToken == nil || equalsToken.ID != EqualsSymbolToken {
		return fmt.Errorf("[processConstantAssignment] expected '=' after identifier '%s'", variableName)
	}

	// Get value
	valueToken := asmTokens.Next()
	if valueToken == nil {
		return fmt.Errorf("[processConstantAssignment] expected value after %s '='", variableName)
	}

	evaluatedValue, err := a.EvaluateExpression(asmTokens, "", true)
	if err != nil {
		return fmt.Errorf("[processConstantAssignment] failed to evaluate expression for %s: %w", variableName, err)
	}

	a.constants[variableName] = evaluatedValue
	return nil
}
