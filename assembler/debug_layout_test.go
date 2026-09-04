package assembler

import (
	"io"
	"strings"
	"testing"

	"github.com/jrsteele09/go-6502-emulator/cpu"
	"github.com/jrsteele09/go-6502-emulator/memory"
	"github.com/jrsteele09/go-6502-emulator/utils"
)

func TestDebugKlausLabels(t *testing.T) {
	mem := memory.NewMemory[uint16](64 * 1024)
	c := cpu.NewCPU(mem, true)
	a := New(c.OpCodes())
	resolver := utils.NewOSFileResolver("/tmp")
	reader, err := resolver.Resolve("6502_functional_test.a65")
	if err != nil {
		t.Fatal(err)
	}
	lex := NewAssemblerLexer(resolver)
	tokens, err := lex.Tokens(a.lexerConfig, reader, "6502_functional_test.a65")
	if err != nil {
		t.Fatal(err)
	}
	_, err = a.resolveLayout(tokens)
	t.Logf("err=%v tdad=%04x tdad6=%04x tdad7=%04x pc=%04x labels=%d", err, a.labels["tdad"], a.labels["tdad6"], a.labels["tdad7"], a.programCounter, len(a.labels))

	reader, _ = resolver.Resolve("6502_functional_test.a65")
	data, _ := io.ReadAll(reader)
	processed, _ := preprocessSource(string(data), "6502_functional_test.a65")
	lines := strings.Split(processed, "\n")
	for i := 7850; i <= 7935 && i <= len(lines); i++ {
		t.Logf("%d: %s", i, lines[i-1])
	}
}
