package debugger

import (
	"testing"
	"time"

	"github.com/jrsteele09/go-6502-emulator/cpu"
	"github.com/stretchr/testify/require"
)

func runDebuggerWithTimeout(t *testing.T, dbg *Debugger) string {
	t.Helper()

	result := make(chan string, 1)
	go func() {
		result <- dbg.Go(nil)
	}()

	select {
	case output := <-result:
		return output
	case <-time.After(time.Second):
		dbg.Stop()
		t.Fatal("debugger did not stop")
		return ""
	}
}

func TestSetProgramCounter(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected uint16
	}{
		{name: "hexadecimal", value: "$C000", expected: 0xC000},
		{name: "decimal", value: "49152", expected: 0xC000},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dbg := NewDebugger()

			output := dbg.SetProgramCounter(test.value)

			require.Equal(t, test.expected, dbg.GetCPU().Registers().PC)
			require.Equal(t, test.expected, dbg.GetLastDisasmAddr())
			require.Equal(t, "PC set to $C000\n", output)
		})
	}
}

func TestSetProgramCounterRejectsInvalidAddress(t *testing.T) {
	dbg := NewDebugger()
	dbg.GetCPU().Registers().PC = 0x1234
	dbg.SetLastDisasmAddr(0x5678)

	output := dbg.SetProgramCounter("$NOPE")

	require.Contains(t, output, "Error setting PC:")
	require.Equal(t, uint16(0x1234), dbg.GetCPU().Registers().PC)
	require.Equal(t, uint16(0x5678), dbg.GetLastDisasmAddr())
}

func TestShowRegistersVerbose(t *testing.T) {
	dbg := NewDebugger()
	regs := dbg.GetCPU().Registers()
	regs.A = 0x80
	regs.X = 0x10
	regs.Y = 0x20
	regs.S = 0xFD
	regs.PC = 0xC000
	regs.Status = byte(cpu.NegativeFlag | cpu.UnusedFlag | cpu.DecimalFlag | cpu.ZeroFlag)

	output := dbg.ShowRegistersVerbose()

	require.Contains(t, output, "Accumulator          A    $80  (128)")
	require.Contains(t, output, "X index register     X    $10  (16)")
	require.Contains(t, output, "Y index register     Y    $20  (32)")
	require.Contains(t, output, "Stack pointer        S    $FD  (253)")
	require.Contains(t, output, "Program counter      PC   $C000  (49152)")
	require.Contains(t, output, "Processor status     P    $AA  (%10101010)")

	expectedFlags := []string{
		"Negative             N    set",
		"Overflow             V    clear",
		"Unused               U    set",
		"Break                B    clear",
		"Decimal mode         D    set",
		"Interrupt disable    I    clear",
		"Zero                 Z    set",
		"Carry                C    clear",
	}
	for _, expected := range expectedFlags {
		require.Contains(t, output, expected)
	}
}

func TestGoStopsBeforeBRK(t *testing.T) {
	dbg := NewDebugger()
	dbg.GetMemory().Write(0x1000, brkOpcode)
	dbg.GetCPU().Registers().PC = 0x1000

	output := runDebuggerWithTimeout(t, dbg)

	require.Contains(t, output, "BRK encountered at $1000")
	require.Equal(t, uint16(0x1000), dbg.GetCPU().Registers().PC)
	require.Equal(t, byte(0xFF), dbg.GetCPU().Registers().S)
	require.False(t, dbg.GetCPU().Registers().IsSet(cpu.BreakFlag))
	require.False(t, dbg.IsRunning())
}

func TestStepStopsBeforeBRK(t *testing.T) {
	dbg := NewDebugger()
	dbg.GetMemory().Write(0x1000, brkOpcode)
	dbg.GetCPU().Registers().PC = 0x1000

	output := dbg.Step(nil)

	require.Contains(t, output, "BRK encountered at $1000")
	require.Equal(t, uint16(0x1000), dbg.GetCPU().Registers().PC)
	require.Equal(t, byte(0xFF), dbg.GetCPU().Registers().S)
	require.False(t, dbg.GetCPU().Registers().IsSet(cpu.BreakFlag))
}

func TestGoStopsBeforeTopLevelRTS(t *testing.T) {
	dbg := NewDebugger()
	dbg.GetMemory().Write(0x1000, rtsOpcode)
	dbg.GetCPU().Registers().PC = 0x1000

	output := runDebuggerWithTimeout(t, dbg)

	require.Contains(t, output, "Top-level RTS encountered at $1000")
	require.Equal(t, uint16(0x1000), dbg.GetCPU().Registers().PC)
	require.Equal(t, byte(0xFF), dbg.GetCPU().Registers().S)
	require.False(t, dbg.IsRunning())
}

func TestStepStopsBeforeTopLevelRTS(t *testing.T) {
	dbg := NewDebugger()
	dbg.GetMemory().Write(0x1000, rtsOpcode)
	dbg.GetCPU().Registers().PC = 0x1000

	output := dbg.Step(nil)

	require.Contains(t, output, "Top-level RTS encountered at $1000")
	require.Equal(t, uint16(0x1000), dbg.GetCPU().Registers().PC)
	require.Equal(t, byte(0xFF), dbg.GetCPU().Registers().S)
}

func TestStepExecutesRTSWhenReturnAddressIsOnStack(t *testing.T) {
	dbg := NewDebugger()
	// JSR $1004; BRK; RTS. Two steps execute JSR and its matching RTS,
	// leaving PC on the instruction following JSR with an empty stack.
	dbg.GetMemory().Write(0x1000, 0x20, 0x04, 0x10, brkOpcode, rtsOpcode)
	dbg.GetCPU().Registers().PC = 0x1000

	output := dbg.Step([]string{"2"})

	require.NotContains(t, output, "Top-level RTS")
	require.Equal(t, uint16(0x1003), dbg.GetCPU().Registers().PC)
	require.Equal(t, byte(0xFF), dbg.GetCPU().Registers().S)
}

func TestGoExecutesRTSWhenReturnAddressIsOnStack(t *testing.T) {
	dbg := NewDebugger()
	// JSR $1004; BRK; RTS. The nested RTS must return to the BRK rather than
	// being mistaken for the top-level return.
	dbg.GetMemory().Write(0x1000, 0x20, 0x04, 0x10, brkOpcode, rtsOpcode)
	dbg.GetCPU().Registers().PC = 0x1000

	output := runDebuggerWithTimeout(t, dbg)

	require.Contains(t, output, "BRK encountered at $1003")
	require.NotContains(t, output, "Top-level RTS")
	require.Equal(t, uint16(0x1003), dbg.GetCPU().Registers().PC)
	require.Equal(t, byte(0xFF), dbg.GetCPU().Registers().S)
	require.False(t, dbg.IsRunning())
}

func TestGoCanBeStoppedAsynchronously(t *testing.T) {
	dbg := NewDebugger()
	dbg.GetMemory().Write(0x1000, 0x4c) // JMP $1000
	dbg.GetMemory().Write(0x1001, 0x00)
	dbg.GetMemory().Write(0x1002, 0x10)
	dbg.GetCPU().Registers().PC = 0x1000

	result := make(chan string, 1)
	go func() {
		result <- dbg.Go(nil)
	}()

	require.Eventually(t, dbg.IsRunning, time.Second, time.Millisecond)
	dbg.Stop()

	select {
	case output := <-result:
		require.Contains(t, output, "Execution stopped at $1000")
	case <-time.After(time.Second):
		t.Fatal("debugger did not stop")
	}
	require.False(t, dbg.IsRunning())
}
