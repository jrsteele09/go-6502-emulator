package assembler

import (
	"fmt"

	"github.com/jrsteele09/go-lexer/lexer"
)

// Operator precedence levels for Pratt parser
const (
	PRECEDENCE_LOWEST  = iota
	PRECEDENCE_SUM     // +, -
	PRECEDENCE_PRODUCT // *, /
	PRECEDENCE_PREFIX  // -x (unary minus)
)

// getPrecedence returns the precedence of an operator token
func (a *Assembler) getPrecedence(tokenID lexer.TokenIdentifier) int {
	switch tokenID {
	case PlusToken, MinusToken:
		return PRECEDENCE_SUM
	case AsterixSymbolToken, DivideSymbolToken:
		return PRECEDENCE_PRODUCT
	default:
		return PRECEDENCE_LOWEST
	}
}

func (a *Assembler) parseNextExpression(asmTokens *AssemblerTokens, mnemonic string, precedence int, preprocess bool) (int64, error) {
	asmTokens.Next() // Advance to next token
	return a.parseCurrentExpression(asmTokens, mnemonic, precedence, preprocess)
}

// parseExpression implements a Pratt parser for mathematical expressions
func (a *Assembler) parseCurrentExpression(asmTokens *AssemblerTokens, mnemonic string, precedence int, preprocess bool) (int64, error) {
	// Parse prefix expression (primary)
	left, err := a.parsePrimary(asmTokens, mnemonic, preprocess)
	if err != nil {
		return 0, err
	}

	// Parse infix expressions based on precedence
	for {
		nextToken := asmTokens.Peek()
		if nextToken == nil || nextToken.ID == lexer.EndOfLineType || nextToken.ID == lexer.EOFType {
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
		default:
			return 0, fmt.Errorf("[parseExpression] unknown operator: %s", operatorToken.Literal)
		}
	}

	return left, nil
}

// parsePrimary parses primary expressions (literals, identifiers, parentheses, unary minus)
func (a *Assembler) parsePrimary(asmTokens *AssemblerTokens, mnemonic string, preprocess bool) (int64, error) {
	token := asmTokens.Current()
	if token == nil {
		return 0, fmt.Errorf("[parsePrimary] unexpected end of expression")
	}

	switch token.ID {
	case lexer.HexLiteral, lexer.IntegerLiteral:
		// Convert literal to int64
		value, err := toInt64(token.Value)
		if err != nil {
			return 0, fmt.Errorf("[parsePrimary] invalid literal: %w", err)
		}
		return value, nil

	case IdentifierToken:
		_, value, err := a.LabelOrConstantIdentifier(mnemonic, token.Literal, preprocess)
		if err != nil || value == nil {
			return 0, fmt.Errorf("[parsePrimary] identifier lookup failed: %w", err)
		}

		return toInt64(value)

	case MinusToken:
		// Unary minus
		right, err := a.parseNextExpression(asmTokens, mnemonic, PRECEDENCE_PREFIX, preprocess)
		if err != nil {
			return 0, err
		}
		return -right, nil

	case LeftParenthesis:
		// Parenthesized expression
		result, err := a.parseNextExpression(asmTokens, mnemonic, PRECEDENCE_LOWEST, preprocess)
		if err != nil {
			return 0, err
		}

		// Expect closing parenthesis
		closeParen := asmTokens.Next()
		if closeParen == nil || closeParen.ID != RightParenthesis {
			return 0, fmt.Errorf("[parsePrimary] expected closing parenthesis")
		}

		return result, nil

	default:
		return 0, fmt.Errorf("[parsePrimary] unexpected token: %s", token.Literal)
	}
}
