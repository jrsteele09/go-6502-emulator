package disassembler

import (
	"fmt"
	"strings"

	"github.com/jrsteele09/go-65xx-emulator/cpu65xxx"
	"github.com/jrsteele09/go-65xx-emulator/memory"
)

const (
	disassemblyFormat = "%-6s %-10s %-10s"
)

type Disassembler struct {
	mem     *memory.Memory[uint16]
	opCodes []*cpu65xxx.OpCodeDef
}

func NewDissassembler(mem *memory.Memory[uint16], opCodes []*cpu65xxx.OpCodeDef) *Disassembler {
	return &Disassembler{
		mem:     mem,
		opCodes: opCodes,
	}
}

func (d *Disassembler) Disassemble(address uint16) (string, int) {
	b := d.mem.Read(address)
	opCode := d.opCodes[b]

	if opCode == nil {
		return strings.TrimSpace(fmt.Sprintf(disassemblyFormat,
			strings.ToUpper(strings.Replace(fmt.Sprintf("%x:", address), "0x", "", 1)),
			strings.ToUpper(d.operandsToByteString(b, []byte{}, 1)),
			"???")), 1
	}

	operands := make([]byte, opCode.Bytes)
	for i := 0; i < opCode.Bytes-1; i++ {
		operands[i] = d.mem.Read(uint16(address + uint16(1+i)))
	}

	return strings.TrimSpace(fmt.Sprintf(disassemblyFormat,
		strings.ToUpper(strings.Replace(fmt.Sprintf("%x:", address), "0x", "", 1)),
		strings.ToUpper(d.operandsToByteString(b, operands, opCode.Bytes)),
		strings.ToUpper(d.opCodeToString(*opCode, operands, address)))), opCode.Bytes
}

func (d *Disassembler) operandsToByteString(opcode byte, operands []byte, l int) string {
	byteArray := []byte{opcode}
	byteArray = append(byteArray, operands...)
	bytesStr := ""
	for i := 0; i < l; i++ {
		bytesStr += fmt.Sprintf("%s ", operandToHexString(i, byteArray...))
	}
	return strings.TrimSpace(bytesStr)
}

func (d *Disassembler) opCodeToString(opcode cpu65xxx.OpCodeDef, operands []byte, address uint16) string {
	addressingModeString := ""

	switch opcode.AddressingModeType {
	case cpu65xxx.RelativeModeStr:
		address = uint16(int64(address) + int64(int8(operands[0])))
		addressingModeString = strings.Replace(fmt.Sprintf("%.2x", address), "0x", "", 1)
	default:
		if strings.Contains(string(opcode.AddressingModeType), cpu65xxx.WordAddressing) {
			operandString := fmt.Sprintf("$%s%s", operandToHexString(1, operands...), operandToHexString(0, operands...))
			addressingModeString = strings.Replace(string(opcode.AddressingModeType), cpu65xxx.WordAddressing, operandString, 1)
		} else if strings.Contains(string(opcode.AddressingModeType), cpu65xxx.ByteAddressing) {
			operandString := fmt.Sprintf("$%s", operandToHexString(0, operands...))
			addressingModeString = strings.Replace(string(opcode.AddressingModeType), cpu65xxx.ByteAddressing, operandString, 1)
		}
	}

	return fmt.Sprintf("%s %s", opcode.Mnemonic, addressingModeString)
}

func operandToHexString(idx int, operands ...byte) string {
	if len(operands) == 0 || uint(idx) > uint(len(operands)) {
		return "??"
	}
	return strings.Replace(fmt.Sprintf("%.2x", operands[idx]), "0x", "", 1)
}
