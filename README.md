# go-6502-emulator

[![Go Report Card](https://goreportcard.com/badge/github.com/jrsteele09/go-6502-emulator)](https://goreportcard.com/report/github.com/jrsteele09/go-6502-emulator)
[![GoDoc](https://pkg.go.dev/badge/github.com/jrsteele09/go-6502-emulator)](https://pkg.go.dev/github.com/jrsteele09/go-6502-emulator)

`go-6502-emulator` is a Go library that emulates the 6502 processor, providing memory management, Cycle Exact CPU operations, and a disassembler for interpreting machine code.


## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
  - [Example](#example)
  - [Disassembler](#disassembler)
- [Contributing](#contributing)
- [License](#license)

## Installation

To install the package, run:

```bash
go get -u github.com/jrsteele09/go-6502-emulator
```
## Usage

### Example
```go
package main

import (
    "fmt"
    "go-6502-emulator/cpu"
    "go-6502-emulator/memory"
)

func main() {
    // Initialize memory with a size of 64KB
    mem := memory.NewMemory[uint16](64 * 1024)
    
    // Initialize the CPU with the memory
    cpu := cpu.NewCpu(mem)
    
    // Write some instructions to memory (example instructions)
    startAddress := uint16(0x0200)
    mem.Write(startAddress, 0xA9, 0x05, 0x00) // LDA #$05; BRK

    // Set the program counter to the start address
    cpu.Registers().PC = startAddress
    
    // Execute instructions until completion
    var complete cpu.Completed
    for !complete {
        c, err := cpu.Execute()
        if err != nil {
            fmt.Println("Error executing instruction:", err)
            break
        }
        complete = c
    }

}
```
### Disassembler

The library includes a disassembler to interpret 6502 machine code into human-readable assembly instructions.

```go

	m := memory.NewMemory[uint16](64 * 1024)
	p := cpu.NewCpu(m)
	m.Write(0xC000,
		0x02,
		0xA9, 0x01,
		0xA9, 0x80,
		0xA9, 0x00,
		0xA5, 0x80,
		0xB5, 0x80,
		0xAD, 0x80, 0x00,
		0xBD, 0x80, 0x00,
		0xBD, 0x01, 0x00,
		0xB9, 0x80, 0x00,
		0xB9, 0x01, 0x00,
		0xA1, 0x05,
		0xB1, 0x05,
	) //0xc01d

	dissassembler := NewDisassembler(m, cpu.OpCodes(p))
	dissassembledCode := ""

	address := uint16(0xC000)
	for address < uint16(0xC01D) {
		line, bytes := dissassembler.Disassemble(address)
		dissassembledCode += fmt.Sprintf("%s\n",line)
		address += uint16(bytes)
	}

```