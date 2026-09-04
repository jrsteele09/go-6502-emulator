package assembler_test

import (
	"testing"

	"github.com/jrsteele09/go-6502-emulator/assembler"
	"github.com/stretchr/testify/require"
)

func TestEvaluateSourceExpression_KlausStyleOperators(t *testing.T) {
	symbols := map[string]int64{
		"data_segment":     0x0200,
		"load_data_direct": 1,
		"ROM_vectors":      1,
		"intdis":           0x04,
		"m8":               0xff,
	}
	lookup := func(name string) (int64, bool) {
		value, ok := symbols[name]
		return value, ok
	}

	tests := map[string]int64{
		"(data_segment & $ff) != 0":                  0,
		"(load_data_direct = 1) & (ROM_vectors = 1)": 1,
		"$ff&~intdis": 0xfb,
		"1<<3":        0x08,
		"lo($1234)":   0x34,
		"hi($1234)":   0x12,
		"<$1234":      0x34,
		">$1234":      0x12,
		"'R'-3":       0x4f,
		"$ff-'B'":     0xbd,
		"m8 - intdis": 0xfb,
	}

	for expression, expected := range tests {
		t.Run(expression, func(t *testing.T) {
			actual, err := assembler.EvaluateSourceExpression(expression, lookup)
			require.NoError(t, err)
			require.Equal(t, expected, actual)
		})
	}
}

func TestEvaluateSourceCondition(t *testing.T) {
	lookup := func(name string) (int64, bool) {
		return map[string]int64{"enabled": 1}[name], name == "enabled"
	}

	actual, err := assembler.EvaluateSourceCondition("enabled = 1", lookup)
	require.NoError(t, err)
	require.True(t, actual)

	actual, err = assembler.EvaluateSourceCondition("enabled = 0", lookup)
	require.NoError(t, err)
	require.False(t, actual)
}
