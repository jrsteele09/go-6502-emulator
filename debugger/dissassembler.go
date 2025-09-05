package debugger

import (
	"fmt"
	"strings"

	"github.com/jrsteele09/go-6502-emulator/cpu"
	"github.com/jrsteele09/go-6502-emulator/memory"
)

const (
	// disassemblyFormat defines the format for disassembly output.
	disassemblyFormat = "%-6s %-10s %-10s"
)

// Disassembler is used to convert machine code into human-readable assembly instructions.
type Disassembler struct {
	mem     *memory.Memory[uint16]
	opCodes []*cpu.OpCodeDef
}

// NewDisassembler creates a new Disassembler instance.
func NewDisassembler(mem *memory.Memory[uint16], opCodes []*cpu.OpCodeDef) *Disassembler {
	return &Disassembler{
		mem:     mem,
		opCodes: opCodes,
	}
}

// Disassemble disassembles the machine code at the given address and returns the assembly instruction and its length.
func (d *Disassembler) Disassemble(address uint16) (string, int) {
	b := d.mem.Read(address)
	opCode := d.opCodes[b]

	if opCode == nil {
		return strings.TrimSpace(fmt.Sprintf(disassemblyFormat,
			fmt.Sprintf("$%04X:", address),
			strings.ToUpper(d.operandsToByteString(b, []byte{}, 1)),
			"???")), 1
	}

	operands := make([]byte, opCode.Bytes)
	for i := 0; i < opCode.Bytes-1; i++ {
		operands[i] = d.mem.Read(uint16(address + uint16(1+i)))
	}

	return strings.TrimSpace(fmt.Sprintf(disassemblyFormat,
		fmt.Sprintf("$%04X:", address),
		strings.ToUpper(d.operandsToByteString(b, operands, opCode.Bytes)),
		strings.ToUpper(d.opCodeToString(*opCode, operands, address)))), opCode.Bytes
}

// operandsToByteString converts the operands to a string of hexadecimal values.
func (d *Disassembler) operandsToByteString(opcode byte, operands []byte, l int) string {
	byteArray := []byte{opcode}
	byteArray = append(byteArray, operands...)
	bytesStr := ""
	for i := 0; i < l; i++ {
		bytesStr += fmt.Sprintf("%s ", operandToHexString(i, byteArray...))
	}
	return strings.TrimSpace(bytesStr)
}

// opCodeToString converts the opcode and operands to a human-readable string based on the addressing mode.
func (d *Disassembler) opCodeToString(opcode cpu.OpCodeDef, operands []byte, address uint16) string {
	addressingModeString := ""

	switch opcode.AddressingModeType {
	case cpu.RelativeModeStr:
		// Address is the address of the opcode, so to calculate a relative address,
		// We need to add on two bytes for the opcode and the operand.
		address = uint16(int64(address) + 2 + int64(int8(operands[0])))
		addressingModeString = strings.Replace(fmt.Sprintf("$%.4x", address), "0x", "", 1)
	default:
		if strings.Contains(string(opcode.AddressingModeType), cpu.WordAddressing) {
			operandString := fmt.Sprintf("$%s%s", operandToHexString(1, operands...), operandToHexString(0, operands...))
			addressingModeString = strings.Replace(string(opcode.AddressingModeType), cpu.WordAddressing, operandString, 1)
		} else if strings.Contains(string(opcode.AddressingModeType), cpu.ByteAddressing) {
			operandString := fmt.Sprintf("$%s", operandToHexString(0, operands...))
			addressingModeString = strings.Replace(string(opcode.AddressingModeType), cpu.ByteAddressing, operandString, 1)
		}
	}

	return fmt.Sprintf("%s %s", opcode.Mnemonic, addressingModeString)
}

// operandToHexString converts a single operand to its hexadecimal string representation.
func operandToHexString(idx int, operands ...byte) string {
	if len(operands) == 0 || uint(idx) > uint(len(operands)) {
		return "??"
	}
	return strings.Replace(fmt.Sprintf("%.2x", operands[idx]), "0x", "", 1)
}
