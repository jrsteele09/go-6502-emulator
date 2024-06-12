package cpu65xxx

import (
	"fmt"
	"strings"
)

type OpCodeDef struct {
	Mnemonic           string
	AddressingModeType AddressingModeType
	Bytes              int
	Cycles             int
	ExecGetter         InstructionGetter
	AddressingMode     AddressingMode
}

type InstructionFunc = func() (Completed, error)
type InstructionGetter func(OpCodeDef) InstructionFunc

func Mnemonic(mnemonic string, am AddressingModeType) string {
	return strings.TrimSpace(fmt.Sprintf("%s %s", mnemonic, string(am)))
}

type InstructionDefinition struct {
	addressingModeGetter func(am AddressingModeType) AddressingMode
}

func NewInstructionDefinition(am func(am AddressingModeType) AddressingMode) InstructionDefinition {
	return InstructionDefinition{
		addressingModeGetter: am,
	}
}

func (id InstructionDefinition) Instruction(Mnemonic string, cycles int, execGet InstructionGetter) *OpCodeDef {
	oc := &OpCodeDef{Mnemonic: Mnemonic}

	components := strings.Split(strings.TrimSpace(Mnemonic), " ")

	oc.Mnemonic = components[0]
	if len(components) > 1 {
		oc.AddressingModeType = AddressingModeType(components[1])
	}

	bytes := 1
	if strings.Contains(string(oc.AddressingModeType), WordAddressing) {
		bytes += 2
	} else if strings.Contains(string(oc.AddressingModeType), ByteAddressing) {
		bytes++
	}
	oc.Bytes = bytes
	oc.Cycles = cycles

	oc.AddressingMode = id.addressingModeGetter(oc.AddressingModeType)
	oc.ExecGetter = execGet

	return oc
}
