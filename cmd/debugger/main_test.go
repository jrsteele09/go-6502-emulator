package main

import (
	"strings"
	"testing"

	"github.com/jrsteele09/go-6502-emulator/debugger"
	"github.com/stretchr/testify/require"
)

func TestParseProgramCounterCommand(t *testing.T) {
	tests := []struct {
		input   string
		value   string
		matched bool
	}{
		{input: "PC=$C000", value: "$C000", matched: true},
		{input: "pc = $c000", value: "$c000", matched: true},
		{input: "PC $C000", value: "$C000", matched: true},
		{input: "PC=49152", value: "49152", matched: true},
		{input: "PC", value: "", matched: true},
		{input: "PCM=$C000", value: "", matched: false},
		{input: "G $C000", value: "", matched: false},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			value, matched := parseProgramCounterCommand(test.input)
			require.Equal(t, test.matched, matched)
			require.Equal(t, test.value, value)
		})
	}
}

func TestVerboseRegisterArguments(t *testing.T) {
	for _, argument := range []string{"V", "verbose", "D", "descriptive", "DETAILED"} {
		require.True(t, isVerboseRegisterArgument(argument), argument)
	}
	require.False(t, isVerboseRegisterArgument("compact"))
}

func TestPromptShowsCurrentInstruction(t *testing.T) {
	dbg := debugger.NewDebugger()
	dbg.GetMemory().Write(0xC000, 0xA9, 0x42)
	dbg.GetCPU().Registers().PC = 0xC000
	repl := &DebuggerRepl{debugger: dbg}

	prompt := repl.prompt()

	require.Contains(t, prompt, "$C000: A9 42      LDA #$42")
	require.True(t, strings.HasSuffix(prompt, Reset+" > "))
}
