package assembler

import (
	"fmt"

	"github.com/jrsteele09/go-lexer/lexer"
)

// parseCurrentExpression adapts assembler lexer tokens to the expression parser
// shared with source constants and source conditionals.
func (a *Assembler) parseCurrentExpression(asmTokens *Tokens, mnemonic string, preprocess bool) (int64, error) {
	return a.parseCurrentExpressionWithUnknowns(asmTokens, mnemonic, preprocess)
}

func (a *Assembler) evaluateConditionalExpression(asmTokens *Tokens) (int64, error) {
	// Layout conditions must be decidable at their source position. This avoids
	// selecting a branch using a placeholder and changing layout in generation.
	return a.parseCurrentExpressionWithUnknowns(asmTokens, "", false)
}

func (a *Assembler) parseCurrentExpressionWithUnknowns(asmTokens *Tokens, mnemonic string, allowUnknown bool) (int64, error) {
	tokens, err := collectAssemblerExpressionTokens(asmTokens)
	if err != nil {
		return 0, err
	}

	var lookupErr error
	lookup := func(name string) (int64, bool) {
		if name == "*" {
			return int64(a.programCounter), true
		}
		if value, ok := a.constants[name]; ok {
			converted, err := toInt64(value)
			if err != nil {
				lookupErr = err
				return 0, false
			}
			return converted, true
		}
		if address, ok := a.labels[name]; ok {
			if mnemonic == "" {
				return int64(address), true
			}
			_, value, err := a.parseLabelOffset(mnemonic, address)
			if err != nil {
				lookupErr = err
				return 0, false
			}
			converted, err := toInt64(value)
			if err != nil {
				lookupErr = err
				return 0, false
			}
			return converted, true
		}
		if !allowUnknown {
			return 0, false
		}
		if mnemonic == "" {
			return 0, true
		}
		_, value, err := a.layoutLabelSizer(mnemonic)
		if err != nil {
			lookupErr = err
			return 0, false
		}
		converted, err := toInt64(value)
		if err != nil {
			lookupErr = err
			return 0, false
		}
		return converted, true
	}

	parser := &sourceExpressionParser{tokens: tokens, lookup: lookup}
	value, err := parser.parseExpression(sourcePrecedenceLowest)
	if lookupErr != nil {
		return 0, lookupErr
	}
	if err != nil {
		return 0, err
	}
	if parser.peek().kind != sourceExpressionTokenEOF {
		return 0, fmt.Errorf("unexpected token %q", parser.peek().literal)
	}
	return value, nil
}

func collectAssemblerExpressionTokens(asmTokens *Tokens) ([]sourceExpressionToken, error) {
	first, err := assemblerExpressionToken(asmTokens.Current())
	if err != nil {
		return nil, err
	}
	tokens := []sourceExpressionToken{first}
	parenthesisDepth := expressionParenthesisDelta(first)
	expectOperand := expressionExpectsOperandAfter(first, true)

	for {
		next := asmTokens.Peek()
		if isTerminatorToken(next.ID) || next.ID == CommaToken {
			break
		}
		if next.ID == RightParenthesis && parenthesisDepth == 0 {
			break
		}

		token, err := assemblerExpressionToken(next)
		if err != nil {
			return nil, err
		}
		previous := tokens[len(tokens)-1]
		if !expressionTokenContinues(previous, token, expectOperand, parenthesisDepth) {
			break
		}
		asmTokens.Next()
		tokens = append(tokens, token)
		parenthesisDepth += expressionParenthesisDelta(token)
		expectOperand = expressionExpectsOperandAfter(token, expectOperand)
	}

	tokens = append(tokens, sourceExpressionToken{kind: sourceExpressionTokenEOF})
	return tokens, nil
}

func expressionTokenContinues(previous, next sourceExpressionToken, expectOperand bool, parenthesisDepth int) bool {
	if expectOperand {
		return next.kind == sourceExpressionTokenInteger ||
			next.kind == sourceExpressionTokenIdentifier ||
			next.kind == sourceExpressionTokenLeftParen ||
			next.kind == sourceExpressionTokenOperator
	}
	if next.kind == sourceExpressionTokenOperator {
		return next.literal != "~"
	}
	if next.kind == sourceExpressionTokenRightParen {
		return parenthesisDepth > 0
	}
	if previous.kind != sourceExpressionTokenIdentifier || !isSourceExpressionFunction(previous.literal) {
		return false
	}
	return next.kind == sourceExpressionTokenLeftParen ||
		next.kind == sourceExpressionTokenInteger ||
		next.kind == sourceExpressionTokenIdentifier ||
		next.kind == sourceExpressionTokenOperator
}

func expressionExpectsOperandAfter(token sourceExpressionToken, wasExpectingOperand bool) bool {
	switch token.kind {
	case sourceExpressionTokenLeftParen:
		return true
	case sourceExpressionTokenOperator:
		if wasExpectingOperand && token.literal == "*" {
			return false
		}
		return true
	default:
		return false
	}
}

func expressionParenthesisDelta(token sourceExpressionToken) int {
	switch token.kind {
	case sourceExpressionTokenLeftParen:
		return 1
	case sourceExpressionTokenRightParen:
		return -1
	default:
		return 0
	}
}

func assemblerExpressionToken(token lexer.Token) (sourceExpressionToken, error) {
	switch token.ID {
	case lexer.HexLiteral, lexer.IntegerLiteral:
		value, err := toInt64(token.Value)
		if err != nil {
			return sourceExpressionToken{}, fmt.Errorf("invalid literal: %w", err)
		}
		return sourceExpressionToken{kind: sourceExpressionTokenInteger, literal: token.Literal, value: value}, nil
	case lexer.StringLiteral:
		value, ok := token.Value.(string)
		if !ok || len([]rune(value)) != 1 {
			return sourceExpressionToken{}, fmt.Errorf("expected character literal, got %q", token.Literal)
		}
		return sourceExpressionToken{kind: sourceExpressionTokenInteger, literal: token.Literal, value: int64([]rune(value)[0])}, nil
	case IdentifierToken:
		return sourceExpressionToken{kind: sourceExpressionTokenIdentifier, literal: token.Literal}, nil
	case LeftParenthesis:
		return sourceExpressionToken{kind: sourceExpressionTokenLeftParen, literal: token.Literal}, nil
	case RightParenthesis:
		return sourceExpressionToken{kind: sourceExpressionTokenRightParen, literal: token.Literal}, nil
	case PlusToken, MinusToken, AsterixSymbolToken, DivideSymbolToken,
		AmpersandToken, PipeToken, CaretToken, TildeToken,
		LessThanToken, GreaterThanToken, EqualsSymbolToken, EqualEqualToken,
		NotEqualToken, LessEqualToken, GreaterEqualToken, ShiftLeftToken, ShiftRightToken:
		return sourceExpressionToken{kind: sourceExpressionTokenOperator, literal: token.Literal}, nil
	default:
		return sourceExpressionToken{}, fmt.Errorf("unexpected token %q in expression", token.Literal)
	}
}
