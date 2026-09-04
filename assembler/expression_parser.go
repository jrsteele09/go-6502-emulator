package assembler

import (
	"fmt"

	"github.com/jrsteele09/go-lexer/lexer"
)

// Operator precedence levels for Pratt parser
const (
	PRECEDENCE_LOWEST = iota
	PRECEDENCE_COMPARE
	PRECEDENCE_BIT_OR
	PRECEDENCE_BIT_XOR
	PRECEDENCE_BIT_AND
	PRECEDENCE_SHIFT
	PRECEDENCE_SUM     // +, -
	PRECEDENCE_PRODUCT // *, /
	PRECEDENCE_PREFIX  // -x, ~x, <x, >x
)

// getPrecedence returns the precedence of an operator token
func (a *Assembler) getPrecedence(tokenID lexer.TokenIdentifier) int {
	switch tokenID {
	case EqualsSymbolToken, EqualEqualToken, NotEqualToken, LessThanToken, LessEqualToken, GreaterThanToken, GreaterEqualToken:
		return PRECEDENCE_COMPARE
	case PipeToken:
		return PRECEDENCE_BIT_OR
	case CaretToken:
		return PRECEDENCE_BIT_XOR
	case AmpersandToken:
		return PRECEDENCE_BIT_AND
	case ShiftLeftToken, ShiftRightToken:
		return PRECEDENCE_SHIFT
	case PlusToken, MinusToken:
		return PRECEDENCE_SUM
	case AsterixSymbolToken, DivideSymbolToken:
		return PRECEDENCE_PRODUCT
	default:
		return PRECEDENCE_LOWEST
	}
}

func (a *Assembler) parseNextExpression(asmTokens *Tokens, mnemonic string, precedence int, preprocess bool) (int64, error) {
	asmTokens.Next() // Advance to next token
	return a.parseCurrentExpression(asmTokens, mnemonic, precedence, preprocess)
}

// parseExpression implements a Pratt parser for mathematical expressions
func (a *Assembler) parseCurrentExpression(asmTokens *Tokens, mnemonic string, precedence int, preprocess bool) (int64, error) {
	// Parse prefix expression (primary)
	left, err := a.parsePrimary(asmTokens, mnemonic, preprocess)
	if err != nil {
		return 0, err
	}

	// Parse infix expressions based on precedence
	for {
		nextToken := asmTokens.Peek()
		if isTerminatorToken(nextToken.ID) {
			break
		}

		tokenPrecedence := a.getPrecedence(nextToken.ID)
		if tokenPrecedence <= precedence {
			break
		}

		// Consume the operator token
		operatorToken := asmTokens.Next()

		// Parse the right operand
		right, err := a.parseNextExpression(asmTokens, mnemonic, tokenPrecedence, preprocess)
		if err != nil {
			return 0, err
		}

		// Apply the operator
		switch operatorToken.ID {
		case PlusToken:
			left = left + right
		case MinusToken:
			left = left - right
		case AsterixSymbolToken:
			left = left * right
		case DivideSymbolToken:
			if right == 0 {
				return 0, fmt.Errorf("[parseExpression] division by zero")
			}
			left = left / right
		case ShiftLeftToken:
			if right < 0 {
				return 0, fmt.Errorf("[parseExpression] negative shift count")
			}
			left = left << uint(right)
		case ShiftRightToken:
			if right < 0 {
				return 0, fmt.Errorf("[parseExpression] negative shift count")
			}
			left = left >> uint(right)
		case AmpersandToken:
			left = left & right
		case PipeToken:
			left = left | right
		case CaretToken:
			left = left ^ right
		case EqualsSymbolToken, EqualEqualToken:
			left = expressionBool(left == right)
		case NotEqualToken:
			left = expressionBool(left != right)
		case LessThanToken:
			left = expressionBool(left < right)
		case LessEqualToken:
			left = expressionBool(left <= right)
		case GreaterThanToken:
			left = expressionBool(left > right)
		case GreaterEqualToken:
			left = expressionBool(left >= right)
		default:
			return 0, fmt.Errorf("[parseExpression] unknown operator: %s", operatorToken.Literal)
		}
	}

	return left, nil
}

// parsePrimary parses primary expressions (literals, identifiers, parentheses, unary minus)
func (a *Assembler) parsePrimary(asmTokens *Tokens, mnemonic string, preprocess bool) (int64, error) {
	token := asmTokens.Current()
	if isTerminatorToken(token.ID) {
		return 0, fmt.Errorf("[parsePrimary] unexpected end of expression")
	}

	switch token.ID {
	case lexer.HexLiteral, lexer.IntegerLiteral:
		// Convert literal to int64
		value, err := toInt64(token.Value)
		if err != nil {
			return 0, fmt.Errorf("invalid literal: %w", err)
		}
		return value, nil
	case lexer.StringLiteral:
		value, ok := token.Value.(string)
		if !ok {
			return 0, fmt.Errorf("invalid string literal")
		}
		if len([]rune(value)) != 1 {
			return 0, fmt.Errorf("expected character literal, got %q", value)
		}
		return int64([]rune(value)[0]), nil

	case IdentifierToken:
		if isExpressionByteFunction(token.Literal) {
			return a.parseExpressionByteFunction(asmTokens, mnemonic, token.Literal, preprocess)
		}
		value, err := a.expressionIdentifierValue(mnemonic, token.Literal, preprocess)
		if err != nil {
			return 0, err
		}
		return value, nil

	case MinusToken:
		// Unary minus
		right, err := a.parseNextExpression(asmTokens, mnemonic, PRECEDENCE_PREFIX, preprocess)
		if err != nil {
			return 0, err
		}
		return -right, nil
	case TildeToken:
		right, err := a.parseNextExpression(asmTokens, mnemonic, PRECEDENCE_PREFIX, preprocess)
		if err != nil {
			return 0, err
		}
		return ^right, nil

	case LessThanToken:
		right, err := a.parseNextExpression(asmTokens, mnemonic, PRECEDENCE_PREFIX, preprocess)
		if err != nil {
			return 0, err
		}
		return right & 0xff, nil

	case GreaterThanToken:
		right, err := a.parseNextExpression(asmTokens, mnemonic, PRECEDENCE_PREFIX, preprocess)
		if err != nil {
			return 0, err
		}
		return (right >> 8) & 0xff, nil

	case AsterixSymbolToken:
		return int64(a.programCounter), nil

	case LeftParenthesis:
		// Parenthesized expression
		result, err := a.parseNextExpression(asmTokens, mnemonic, PRECEDENCE_LOWEST, preprocess)
		if err != nil {
			return 0, err
		}

		// Expect closing parenthesis
		closeParen := asmTokens.Next()
		if closeParen.ID != RightParenthesis {
			return 0, fmt.Errorf("expected closing parenthesis")
		}

		return result, nil

	default:
		return 0, fmt.Errorf("unexpected token: %s", token.Literal)
	}
}

func (a *Assembler) parseExpressionByteFunction(asmTokens *Tokens, mnemonic string, name string, preprocess bool) (int64, error) {
	if asmTokens.Peek().ID == LeftParenthesis {
		asmTokens.Next()
		value, err := a.parseNextExpression(asmTokens, mnemonic, PRECEDENCE_LOWEST, preprocess)
		if err != nil {
			return 0, err
		}
		closeParen := asmTokens.Next()
		if closeParen.ID != RightParenthesis {
			return 0, fmt.Errorf("expected closing parenthesis after %s", name)
		}
		return expressionByte(name, value)
	}

	value, err := a.parseNextExpression(asmTokens, mnemonic, PRECEDENCE_PREFIX, preprocess)
	if err != nil {
		return 0, err
	}
	return expressionByte(name, value)
}

func isExpressionByteFunction(name string) bool {
	switch name {
	case "lo", "LO", "Lo", "lO", "hi", "HI", "Hi", "hI":
		return true
	default:
		return false
	}
}

func expressionByte(name string, value int64) (int64, error) {
	switch name {
	case "lo", "LO", "Lo", "lO":
		return value & 0xff, nil
	case "hi", "HI", "Hi", "hI":
		return (value >> 8) & 0xff, nil
	default:
		return 0, fmt.Errorf("unknown byte function %s", name)
	}
}

func expressionBool(value bool) int64 {
	if value {
		return 1
	}
	return 0
}

func (a *Assembler) expressionIdentifierValue(mnemonic, identifier string, preprocess bool) (int64, error) {
	if value, ok := a.constants[identifier]; ok {
		return toInt64(value)
	}

	if address, ok := a.labels[identifier]; ok {
		if mnemonic == "" {
			return int64(address), nil
		}
		_, value, err := a.parseLabelOffset(mnemonic, address)
		if err != nil {
			return 0, err
		}
		return toInt64(value)
	}

	if preprocess {
		if mnemonic == "" {
			return 0, nil
		}
		_, value, err := a.layoutLabelSizer(mnemonic)
		if err != nil {
			return 0, err
		}
		return toInt64(value)
	}

	return 0, fmt.Errorf("undefined identifier: %s", identifier)
}
