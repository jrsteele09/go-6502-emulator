package assembler_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/jrsteele09/go-6502-emulator/assembler"
	"github.com/jrsteele09/go-6502-emulator/utils"
	"github.com/stretchr/testify/require"
)

func TestAssemble_SourcePreprocessorAndBareDirectives(t *testing.T) {
	_, cpu := createHardware()
	asm := assembler.New(cpu.OpCodes())
	resolver := utils.NewOSFileResolver("./test_assembly_files/TestSourcePreprocessor")

	segments, err := asm.AssembleFile("main.asm", resolver)

	require.NoError(t, err, "AssembleFile failed")
	require.Len(t, segments, 1)
	require.Equal(t, uint16(0x1000), segments[0].StartAddress)
	require.Equal(t, []byte{0xA9, 0x42, 0x8D, 0x00, 0xC0, 0xEA, 0x7E, 0x00, 0x01}, segments[0].Data.Bytes())
}

func TestAssemble_SourceExpressionsMacrosAndConditionalsWorkTogether(t *testing.T) {
	_, cpu := createHardware()
	asm := assembler.New(cpu.OpCodes())
	resolver := utils.NewOSFileResolver("./test_assembly_files/TestSourceExpressions")

	segments, err := asm.AssembleFile("main.asm", resolver)

	require.NoError(t, err)
	require.Len(t, segments, 1)
	require.Equal(t, uint16(0x1200), segments[0].StartAddress)
	require.Equal(t, []byte{0xA9, 0x42, 0x8D, 0x10, 0xC0, 0xA5, 0x08, 0x4F, 0xA9, 0xFB}, segments[0].Data.Bytes())
}

func TestAssemble_SourcePreprocessorIgnoresCommentsWhenClassifyingLines(t *testing.T) {
	_, cpu := createHardware()
	asm := assembler.New(cpu.OpCodes())
	resolver := utils.NewOSFileResolver("./test_assembly_files/TestSourceComments")

	segments, err := asm.AssembleFile("main.asm", resolver)

	require.NoError(t, err)
	require.Len(t, segments, 1)
	require.Equal(t, uint16(0x1300), segments[0].StartAddress)
	require.Equal(t, []byte{0xA9, 0x7F, 0x2F, 0x3B, ';', '/', '/', '/', '*'}, segments[0].Data.Bytes())
}

func TestAssemble_LabelledStatements(t *testing.T) {
	_, cpu := createHardware()
	asm := assembler.New(cpu.OpCodes())
	resolver := utils.NewOSFileResolver("./test_assembly_files/TestLabelledStatements")

	segments, err := asm.AssembleFile("main.asm", resolver)

	require.NoError(t, err)
	require.Len(t, segments, 1)
	require.Equal(t, uint16(0x1400), segments[0].StartAddress)
	require.Equal(t, []byte{
		0xA9, 0x01,
		0xEA,
		0x02, 0x03,
		0x00, 0x14,
		0xAD, 0x03, 0x14,
		0x00, 0x00,
		0xA9, 0xFF,
		0x20, 0x00, 0x14,
	}, segments[0].Data.Bytes())
}

func TestAssemble_AssemblerExpressions(t *testing.T) {
	_, cpu := createHardware()
	asm := assembler.New(cpu.OpCodes())
	resolver := utils.NewOSFileResolver("./test_assembly_files/TestAssemblerExpressions")

	segments, err := asm.AssembleFile("main.asm", resolver)

	require.NoError(t, err)
	require.Len(t, segments, 1)
	require.Equal(t, uint16(0x1500), segments[0].StartAddress)
	require.Equal(t, []byte{
		0xA9, 0xBE,
		0x29, 0x0B,
		0x09, 0x09,
		0xC9, 0x4B,
		0xC9, 0xEC,
		0xA9, 0x15,
		0xA2, 0x00,
		0xFB, 0x42, 0x03,
		0x00, 0x15,
		0x01, 0x15,
		0x00, 0xEE, 0xFF,
	}, segments[0].Data.Bytes())
}

func TestAssemble_AssemblyConditionalsCanUseLabels(t *testing.T) {
	_, cpu := createHardware()
	asm := assembler.New(cpu.OpCodes())
	resolver := utils.NewOSFileResolver("./test_assembly_files/TestAssemblyConditionals")

	segments, err := asm.AssembleFile("main.asm", resolver)

	require.NoError(t, err)
	require.Len(t, segments, 1)
	require.Equal(t, uint16(0x1600), segments[0].StartAddress)
	require.Equal(t, []byte{
		0x01, 0x02, 0x03,
		0x00, 0x00, 0x00,
		0xA9, 0x03,
		0xA2, 0x16,
	}, segments[0].Data.Bytes())
}

func TestAssemble_ReassignableConstants(t *testing.T) {
	_, cpu := createHardware()
	asm := assembler.New(cpu.OpCodes())
	resolver := utils.NewOSFileResolver("./test_assembly_files/TestReassignableConstants")

	segments, err := asm.AssembleFile("main.asm", resolver)

	require.NoError(t, err)
	require.Len(t, segments, 1)
	require.Equal(t, uint16(0x1700), segments[0].StartAddress)
	require.Equal(t, []byte{0xA9, 0x00, 0xA2, 0x01, 0xA0, 0x02}, segments[0].Data.Bytes())
}

func TestAssemble_SourcePreprocessorErrorsOnUnterminatedConditional(t *testing.T) {
	_, cpu := createHardware()
	asm := assembler.New(cpu.OpCodes())

	_, err := asm.Assemble(bytes.NewBufferString("if 1\nlda #$01\n"), "unterminated.asm")

	require.Error(t, err)
	require.Contains(t, err.Error(), "unterminated conditional block")
}

func TestAssemble_SourcePreprocessorDoesNotDeferInvalidExpression(t *testing.T) {
	_, cpu := createHardware()
	asm := assembler.New(cpu.OpCodes())

	_, err := asm.Assemble(bytes.NewBufferString("if 1 / 0\n db $01\nendif\n"), "invalid-if.asm")

	require.ErrorContains(t, err, "invalid conditional expression: division by zero")
}

func TestAssemble_ConditionalUsingBackwardLabelIsResolvedDuringLayout(t *testing.T) {
	_, cpu := createHardware()
	asm := assembler.New(cpu.OpCodes())
	source := "org $1000\nstart nop\nif start = $1000\n db $42\nelse\n db $99\nendif\n"

	segments, err := asm.Assemble(strings.NewReader(source), "layout-if.asm")

	require.NoError(t, err)
	require.Len(t, segments, 1)
	require.Equal(t, []byte{0xEA, 0x42}, segments[0].Data.Bytes())
}

func TestAssemble_ConditionalRejectsUnresolvedForwardLabelDuringLayout(t *testing.T) {
	_, cpu := createHardware()
	asm := assembler.New(cpu.OpCodes())
	source := "org $1000\nif target = $1001\n db $42\nendif\ntarget nop\n"

	_, err := asm.Assemble(strings.NewReader(source), "forward-if.asm")

	require.ErrorContains(t, err, `undefined symbol "target"`)
}

func TestAssemble_MacroDefinedInIncludeIsAvailableToParent(t *testing.T) {
	_, cpu := createHardware()
	asm := assembler.New(cpu.OpCodes())
	resolver := utils.NewMemoryFileResolver(map[string]string{
		"main.asm":   "org $1000\n.include \"macros.asm\"\nEMIT $42\n",
		"macros.asm": "EMIT macro value\n db value\nendm\n",
	})

	segments, err := asm.AssembleFile("main.asm", resolver)

	require.NoError(t, err)
	require.Len(t, segments, 1)
	require.Equal(t, []byte{0x42}, segments[0].Data.Bytes())
}

func TestAssemble_InactiveIncludeIsNotResolved(t *testing.T) {
	_, cpu := createHardware()
	asm := assembler.New(cpu.OpCodes())
	resolver := utils.NewMemoryFileResolver(map[string]string{
		"main.asm": "org $1000\nif 0\n.include \"missing.asm\"\nendif\ndb $42\n",
	})

	segments, err := asm.AssembleFile("main.asm", resolver)

	require.NoError(t, err)
	require.Len(t, segments, 1)
	require.Equal(t, []byte{0x42}, segments[0].Data.Bytes())
}

func TestAssemble_EquConstantNameContainingEqu(t *testing.T) {
	_, cpu := createHardware()
	asm := assembler.New(cpu.OpCodes())

	segments, err := asm.Assemble(bytes.NewBufferString("SEQUENCE EQU $42\norg $1000\ndb SEQUENCE\n"), "equ.asm")

	require.NoError(t, err)
	require.Len(t, segments, 1)
	require.Equal(t, []byte{0x42}, segments[0].Data.Bytes())
}

func TestAssemble_EquCompatibilityForms(t *testing.T) {
	_, cpu := createHardware()
	asm := assembler.New(cpu.OpCodes())
	source := `
BASE = $40
FIRST EQU BASE + 1
SECOND .EQU FIRST + 1
.EQU THIRD, SECOND + 1
.EQU FOURTH = THIRD + 1
org $1000
db FIRST, SECOND, THIRD, FOURTH
`

	segments, err := asm.Assemble(strings.NewReader(source), "equ-forms.asm")

	require.NoError(t, err)
	require.Len(t, segments, 1)
	require.Equal(t, []byte{0x41, 0x42, 0x43, 0x44}, segments[0].Data.Bytes())
}
