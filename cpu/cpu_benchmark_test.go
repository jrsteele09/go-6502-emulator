package cpu

import (
	"testing"

	"github.com/jrsteele09/go-6502-emulator/memory"
)

func BenchmarkExecuteNOP(b *testing.B) {
	mem := memory.NewMemory[uint16](64 * 1024)
	for address := uint32(0); address < 64*1024; address++ {
		mem.Write(uint16(address), 0xEA)
	}
	processor := NewCPU(mem, false)
	processor.Reg.PC = 0

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		completed := Completed(false)
		for !completed {
			var err error
			completed, err = processor.Execute()
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}
