package assembler

import (
	"fmt"
	"strings"

	"github.com/jrsteele09/go-lexer/lexer"
)

type assemblyConditionalState struct {
	parentActive bool
	condition    bool
	elseSeen     bool
}

type assemblyConditionalStack struct {
	states []assemblyConditionalState
}

func (s *assemblyConditionalStack) Active() bool {
	active := true
	for _, state := range s.states {
		branchActive := state.condition
		if state.elseSeen {
			branchActive = !state.condition
		}
		active = active && state.parentActive && branchActive
	}
	return active
}

func (s *assemblyConditionalStack) Handle(a *Assembler, token lexer.Token, tokenPosition int, asmTokens *Tokens, preprocess bool) (bool, error) {
	if token.ID != IdentifierToken || tokenPosition != 1 {
		return false, nil
	}

	switch strings.ToLower(token.Literal) {
	case "if":
		return true, s.handleIf(a, asmTokens, preprocess)
	case "else":
		return true, s.handleElse(asmTokens, token)
	case "endif":
		return true, s.handleEndif(asmTokens, token)
	default:
		return false, nil
	}
}

func (s *assemblyConditionalStack) handleIf(a *Assembler, asmTokens *Tokens, preprocess bool) error {
	parentActive := s.Active()
	condition := false

	if parentActive {
		t := asmTokens.Next()
		if isTerminatorToken(t.ID) {
			return fmt.Errorf("[assembly conditional] expected expression after if")
		}
		value, err := a.EvaluateExpression(asmTokens, "", preprocess)
		if err != nil {
			return fmt.Errorf("[assembly conditional] if expression: %w", err)
		}
		condition = value != 0
	}

	s.states = append(s.states, assemblyConditionalState{
		parentActive: parentActive,
		condition:    condition,
	})
	s.skipLineRemainder(asmTokens)
	return nil
}

func (s *assemblyConditionalStack) handleElse(asmTokens *Tokens, token lexer.Token) error {
	if len(s.states) == 0 {
		return fmt.Errorf("[assembly conditional] unexpected else at %d:%d", token.SourceLine, token.SourceColumn)
	}
	state := &s.states[len(s.states)-1]
	if state.elseSeen {
		return fmt.Errorf("[assembly conditional] duplicate else at %d:%d", token.SourceLine, token.SourceColumn)
	}
	state.elseSeen = true
	s.skipLineRemainder(asmTokens)
	return nil
}

func (s *assemblyConditionalStack) handleEndif(asmTokens *Tokens, token lexer.Token) error {
	if len(s.states) == 0 {
		return fmt.Errorf("[assembly conditional] unexpected endif at %d:%d", token.SourceLine, token.SourceColumn)
	}
	s.states = s.states[:len(s.states)-1]
	s.skipLineRemainder(asmTokens)
	return nil
}

func (s *assemblyConditionalStack) Complete() error {
	if len(s.states) != 0 {
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
