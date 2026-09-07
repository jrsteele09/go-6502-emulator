package assembler

import (
	"fmt"
	"strings"

	"github.com/jrsteele09/go-lexer/lexer"
)

type assemblyConditionalStack struct {
	conditionalStack
}

func (s *assemblyConditionalStack) Handle(a *Assembler, token lexer.Token, tokenPosition int, asmTokens *Tokens) (bool, error) {
	if token.ID != IdentifierToken || tokenPosition != 1 {
		return false, nil
	}

	switch strings.ToLower(token.Literal) {
	case "if":
		return true, s.handleIf(a, asmTokens)
	case "else":
		return true, s.handleElse(asmTokens, token)
	case "endif":
		return true, s.handleEndif(asmTokens, token)
	default:
		return false, nil
	}
}

func (s *assemblyConditionalStack) handleIf(a *Assembler, asmTokens *Tokens) error {
	parentActive := s.Active()
	condition := false

	if parentActive {
		t := asmTokens.Next()
		if isTerminatorToken(t.ID) {
			return fmt.Errorf("[assembly conditional] expected expression after if")
		}
		value, err := a.evaluateConditionalExpression(asmTokens)
		if err != nil {
			return fmt.Errorf("[assembly conditional] if expression: %w", err)
		}
		condition = value != 0
	}

	s.Push(parentActive, condition, false)
	s.skipLineRemainder(asmTokens)
	return nil
}

func (s *assemblyConditionalStack) handleElse(asmTokens *Tokens, token lexer.Token) error {
	state := s.Current()
	if state == nil {
		return fmt.Errorf("[assembly conditional] unexpected else at %d:%d", token.SourceLine, token.SourceColumn)
	}
	if state.elseSeen {
		return fmt.Errorf("[assembly conditional] duplicate else at %d:%d", token.SourceLine, token.SourceColumn)
	}
	state.elseSeen = true
	s.skipLineRemainder(asmTokens)
	return nil
}

func (s *assemblyConditionalStack) handleEndif(asmTokens *Tokens, token lexer.Token) error {
	if !s.Pop() {
		return fmt.Errorf("[assembly conditional] unexpected endif at %d:%d", token.SourceLine, token.SourceColumn)
	}
	s.skipLineRemainder(asmTokens)
	return nil
}

func (s *assemblyConditionalStack) Complete() error {
	if !s.conditionalStack.Complete() {
		return fmt.Errorf("[assembly conditional] unterminated conditional block")
	}
	return nil
}

func (s *assemblyConditionalStack) skipLineRemainder(asmTokens *Tokens) {
	for {
		t := asmTokens.Peek()
		if isTerminatorToken(t.ID) {
			asmTokens.Next()
			return
		}
		asmTokens.Next()
	}
}
