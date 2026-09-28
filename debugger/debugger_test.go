package debugger

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

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
