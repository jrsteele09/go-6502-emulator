package cpu

import (
	"fmt"
	"strings"
)

// OpCodeDef represents the definition of an opcode, including its mnemonic, addressing mode type, byte size, cycle count, and execution function.
type OpCodeDef struct {
	Mnemonic           string
	AddressingModeType AddressingModeType
	Bytes              int
	Cycles             int
	GetInstructionFunc InstructionFunctionGetter
	AddressingMode     AddressingMode
}

// InstructionFunc defines a function type for executing an instruction and returning whether it is completed and any error encountered.
type InstructionFunc = func() (Completed, error)

// InstructionFunctionGetter defines a function type for getting the instruction execution function based on the opcode definition.
type InstructionFunctionGetter func(OpCodeDef) InstructionFunc

// Mnemonic formats the mnemonic string with the addressing mode type.
func Mnemonic(mnemonic string, am AddressingModeType) string {
	return strings.TrimSpace(fmt.Sprintf("%s %s", mnemonic, string(am)))
}

// InstructionDefinition holds the function to get the addressing mode based on the addressing mode type.
type InstructionDefinition struct {
	addressingModeGetter func(am AddressingModeType) AddressingMode
}

// NewInstruction creates a new InstructionDefinition instance.
func NewInstruction(am func(am AddressingModeType) AddressingMode) InstructionDefinition {
	return InstructionDefinition{
		addressingModeGetter: am,
	}
}

// Instruction creates an OpCodeDef instance for the given mnemonic, cycles, and execution function getter.
func (id InstructionDefinition) Instruction(Mnemonic string, cycles int, execGet InstructionFunctionGetter) *OpCodeDef {
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
	oc.GetInstructionFunc = execGet

	return oc
}
