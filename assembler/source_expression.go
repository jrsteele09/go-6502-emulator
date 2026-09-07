package assembler

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type SourceSymbolLookup func(name string) (int64, bool)

type UndefinedSymbolError struct {
	Name string
}

func (e *UndefinedSymbolError) Error() string {
	return fmt.Sprintf("undefined symbol %q", e.Name)
}

type sourceExpressionTokenKind int

const (
	sourceExpressionTokenEOF sourceExpressionTokenKind = iota
	sourceExpressionTokenInteger
	sourceExpressionTokenIdentifier
	sourceExpressionTokenLeftParen
	sourceExpressionTokenRightParen
	sourceExpressionTokenOperator
)

type sourceExpressionToken struct {
	kind    sourceExpressionTokenKind
	literal string
	value   int64
}

const (
	sourcePrecedenceLowest = iota
	sourcePrecedenceCompare
	sourcePrecedenceBitOr
	sourcePrecedenceBitXor
	sourcePrecedenceBitAnd
	sourcePrecedenceShift
	sourcePrecedenceSum
	sourcePrecedenceProduct
	sourcePrecedencePrefix
	sourcePrecedenceCall
)

func EvaluateSourceExpression(expression string, lookup SourceSymbolLookup) (int64, error) {
	tokens, err := tokenizeSourceExpression(expression)
	if err != nil {
		return 0, err
	}
	parser := &sourceExpressionParser{
		tokens: tokens,
		lookup: lookup,
	}
	value, err := parser.parseExpression(sourcePrecedenceLowest)
	if err != nil {
		return 0, err
	}
	if parser.peek().kind != sourceExpressionTokenEOF {
		return 0, fmt.Errorf("unexpected token %q", parser.peek().literal)
	}
	return value, nil
}

func EvaluateSourceCondition(expression string, lookup SourceSymbolLookup) (bool, error) {
	value, err := EvaluateSourceExpression(expression, lookup)
	if err != nil {
		return false, err
	}
	return value != 0, nil
}

func tokenizeSourceExpression(expression string) ([]sourceExpressionToken, error) {
	var tokens []sourceExpressionToken
	for i := 0; i < len(expression); {
		r := rune(expression[i])
		if unicode.IsSpace(r) {
			i++
			continue
		}
		if expression[i] == ';' {
			break
		}
		if isSourceExpressionIdentifierStart(r) {
			start := i
			i++
			for i < len(expression) && isSourceExpressionIdentifierPart(rune(expression[i])) {
				i++
			}
			tokens = append(tokens, sourceExpressionToken{
				kind:    sourceExpressionTokenIdentifier,
				literal: expression[start:i],
			})
			continue
		}
		if unicode.IsDigit(r) || expression[i] == '$' || expression[i] == '%' {
			token, next, err := scanSourceExpressionInteger(expression, i)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)
			i = next
			continue
		}
		if expression[i] == '\'' {
			token, next, err := scanSourceExpressionCharacter(expression, i)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)
			i = next
			continue
		}
		if i+1 < len(expression) {
			two := expression[i : i+2]
			switch two {
			case "==", "!=", "<=", ">=", "<<", ">>":
				tokens = append(tokens, sourceExpressionToken{kind: sourceExpressionTokenOperator, literal: two})
				i += 2
				continue
			}
		}
		switch expression[i] {
		case '(':
			tokens = append(tokens, sourceExpressionToken{kind: sourceExpressionTokenLeftParen, literal: "("})
		case ')':
			tokens = append(tokens, sourceExpressionToken{kind: sourceExpressionTokenRightParen, literal: ")"})
		case '+', '-', '*', '/', '&', '|', '^', '~', '<', '>', '=':
			tokens = append(tokens, sourceExpressionToken{kind: sourceExpressionTokenOperator, literal: expression[i : i+1]})
		default:
			return nil, fmt.Errorf("unexpected character %q", expression[i])
		}
		i++
	}
	tokens = append(tokens, sourceExpressionToken{kind: sourceExpressionTokenEOF})
	return tokens, nil
}

func scanSourceExpressionInteger(expression string, start int) (sourceExpressionToken, int, error) {
	base := 10
	i := start
	switch expression[i] {
	case '$':
		base = 16
		i++
	case '%':
		base = 2
		i++
	}
	digitsStart := i
	for i < len(expression) && isSourceExpressionDigit(rune(expression[i]), base) {
		i++
	}
	if digitsStart == i {
		return sourceExpressionToken{}, start, fmt.Errorf("expected digits after %q", expression[start])
	}
	value, err := strconv.ParseInt(expression[digitsStart:i], base, 64)
	if err != nil {
		return sourceExpressionToken{}, start, err
	}
	return sourceExpressionToken{
		kind:    sourceExpressionTokenInteger,
		literal: expression[start:i],
		value:   value,
	}, i, nil
}

func scanSourceExpressionCharacter(expression string, start int) (sourceExpressionToken, int, error) {
	i := start + 1
	if i >= len(expression) {
		return sourceExpressionToken{}, start, fmt.Errorf("unterminated character literal")
	}
	var value rune
	if expression[i] == '\\' {
		i++
		if i >= len(expression) {
			return sourceExpressionToken{}, start, fmt.Errorf("unterminated character literal")
		}
		switch expression[i] {
		case 'n':
			value = '\n'
		case 'r':
			value = '\r'
		case 't':
			value = '\t'
		case '\\':
			value = '\\'
		case '\'':
			value = '\''
		default:
			value = rune(expression[i])
		}
	} else {
		value = rune(expression[i])
	}
	i++
	if i >= len(expression) || expression[i] != '\'' {
		return sourceExpressionToken{}, start, fmt.Errorf("unterminated character literal")
	}
	i++
	return sourceExpressionToken{
		kind:    sourceExpressionTokenInteger,
		literal: expression[start:i],
		value:   int64(value),
	}, i, nil
}

func isSourceExpressionDigit(r rune, base int) bool {
	switch {
	case r >= '0' && r <= '1':
		return true
	case r >= '2' && r <= '9':
		return base >= 10
	case r >= 'a' && r <= 'f':
		return base == 16
	case r >= 'A' && r <= 'F':
		return base == 16
	default:
		return false
	}
}

func isSourceExpressionIdentifierStart(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isSourceExpressionIdentifierPart(r rune) bool {
	return isSourceExpressionIdentifierStart(r) || unicode.IsDigit(r)
}

type sourceExpressionParser struct {
	tokens []sourceExpressionToken
	index  int
	lookup SourceSymbolLookup
}

func (p *sourceExpressionParser) parseExpression(precedence int) (int64, error) {
	token := p.next()
	left, err := p.parsePrefix(token)
	if err != nil {
		return 0, err
	}
	for {
		next := p.peek()
		nextPrecedence := sourceExpressionPrecedence(next)
		if next.kind == sourceExpressionTokenEOF || next.kind == sourceExpressionTokenRightParen || nextPrecedence <= precedence {
			break
		}
		operator := p.next()
		left, err = p.parseInfix(left, operator, nextPrecedence)
		if err != nil {
			return 0, err
		}
	}
	return left, nil
}

func (p *sourceExpressionParser) parsePrefix(token sourceExpressionToken) (int64, error) {
	switch token.kind {
	case sourceExpressionTokenInteger:
		return token.value, nil
	case sourceExpressionTokenIdentifier:
		if isSourceExpressionFunction(token.literal) {
			return p.parseFunctionCall(token.literal)
		}
		if p.lookup == nil {
			return 0, &UndefinedSymbolError{Name: token.literal}
		}
		value, ok := p.lookup(token.literal)
		if !ok {
			return 0, &UndefinedSymbolError{Name: token.literal}
		}
		return value, nil
	case sourceExpressionTokenLeftParen:
		value, err := p.parseExpression(sourcePrecedenceLowest)
		if err != nil {
			return 0, err
		}
		if p.next().kind != sourceExpressionTokenRightParen {
			return 0, fmt.Errorf("expected closing parenthesis")
		}
		return value, nil
	case sourceExpressionTokenOperator:
		switch token.literal {
		case "*":
			if p.lookup == nil {
				return 0, &UndefinedSymbolError{Name: token.literal}
			}
			value, ok := p.lookup(token.literal)
			if !ok {
				return 0, &UndefinedSymbolError{Name: token.literal}
			}
			return value, nil
		case "-":
			value, err := p.parseExpression(sourcePrecedencePrefix)
			if err != nil {
				return 0, err
			}
			return -value, nil
		case "~":
			value, err := p.parseExpression(sourcePrecedencePrefix)
			if err != nil {
				return 0, err
			}
			return ^value, nil
		case "<":
			value, err := p.parseExpression(sourcePrecedencePrefix)
			if err != nil {
				return 0, err
			}
			return value & 0xff, nil
		case ">":
			value, err := p.parseExpression(sourcePrecedencePrefix)
			if err != nil {
				return 0, err
			}
			return (value >> 8) & 0xff, nil
		}
	}
	return 0, fmt.Errorf("unexpected token %q", token.literal)
}

func (p *sourceExpressionParser) parseFunctionCall(name string) (int64, error) {
	precedence := sourcePrecedencePrefix
	parenthesized := p.peek().kind == sourceExpressionTokenLeftParen
	if parenthesized {
		p.next()
		precedence = sourcePrecedenceLowest
	}
	value, err := p.parseExpression(precedence)
	if err != nil {
		return 0, err
	}
	if parenthesized && p.next().kind != sourceExpressionTokenRightParen {
		return 0, fmt.Errorf("expected closing parenthesis after %s", name)
	}
	switch strings.ToLower(name) {
	case "lo":
		return value & 0xff, nil
	case "hi":
		return (value >> 8) & 0xff, nil
	default:
		return 0, fmt.Errorf("unknown function %s", name)
	}
}

func (p *sourceExpressionParser) parseInfix(left int64, operator sourceExpressionToken, precedence int) (int64, error) {
	right, err := p.parseExpression(precedence)
	if err != nil {
		return 0, err
	}
	switch operator.literal {
	case "*":
		return left * right, nil
	case "/":
		if right == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return left / right, nil
	case "+":
		return left + right, nil
	case "-":
		return left - right, nil
	case "<<":
		if right < 0 {
			return 0, fmt.Errorf("negative shift count")
		}
		return left << uint(right), nil
	case ">>":
		if right < 0 {
			return 0, fmt.Errorf("negative shift count")
		}
		return left >> uint(right), nil
	case "&":
		return left & right, nil
	case "^":
		return left ^ right, nil
	case "|":
		return left | right, nil
	case "=", "==":
		return sourceBool(left == right), nil
	case "!=":
		return sourceBool(left != right), nil
	case "<":
		return sourceBool(left < right), nil
	case "<=":
		return sourceBool(left <= right), nil
	case ">":
		return sourceBool(left > right), nil
	case ">=":
		return sourceBool(left >= right), nil
	default:
		return 0, fmt.Errorf("unknown operator %q", operator.literal)
	}
}

func (p *sourceExpressionParser) peek() sourceExpressionToken {
	if p.index >= len(p.tokens) {
		return sourceExpressionToken{kind: sourceExpressionTokenEOF}
	}
	return p.tokens[p.index]
}

func (p *sourceExpressionParser) next() sourceExpressionToken {
	token := p.peek()
	if p.index < len(p.tokens) {
		p.index++
	}
	return token
}

func sourceExpressionPrecedence(token sourceExpressionToken) int {
	if token.kind == sourceExpressionTokenLeftParen {
		return sourcePrecedenceCall
	}
	if token.kind != sourceExpressionTokenOperator {
		return sourcePrecedenceLowest
	}
	switch token.literal {
	case "=", "==", "!=", "<", "<=", ">", ">=":
		return sourcePrecedenceCompare
	case "|":
		return sourcePrecedenceBitOr
	case "^":
		return sourcePrecedenceBitXor
	case "&":
		return sourcePrecedenceBitAnd
	case "<<", ">>":
		return sourcePrecedenceShift
	case "+", "-":
		return sourcePrecedenceSum
	case "*", "/":
		return sourcePrecedenceProduct
	default:
		return sourcePrecedenceLowest
	}
}

func isSourceExpressionFunction(name string) bool {
	switch strings.ToLower(name) {
	case "lo", "hi":
		return true
	default:
		return false
	}
}

func sourceBool(value bool) int64 {
	if value {
		return 1
	}
	return 0
}
