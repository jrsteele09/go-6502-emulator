# 6502 Assembler and Debugger

[![Go Reference](https://pkg.go.dev/badge/github.com/jrsteele09/go-6502-emulator)](https://pkg.go.dev/github.com/jrsteele09/go-6502-emulator)

A simple 6502 assembler and debugger written in Go. This project aims to
provide everything needed to write, assemble, and debug 6502 assembly
programs.

The assembler will continue to evolve over time, but the vision is to be compatible with as many 6502 assemblers as possible.

## Features

**Assembler:**

- Full 6502 instruction set support
- PRG, D64, and T64 file output formats
- Comprehensive error reporting with line numbers
- Support for labels, constants, expressions, and include files

**Debugger:**

- Interactive REPL-style debugger
- Memory examination and modification
- Disassembly with PC and breakpoint markers
- Single-step execution with register display
- Breakpoint management
- PRG loading and execution control

**Testing:**

- Functional testing with the [Klaus 6502 functional tests and decimal mode test](https://github.com/Klaus2m5/6502_65C02_functional_tests)

## Documentation

- [Assembler user manual](docs/assembler.md) — command-line usage, source
  syntax, labels, expressions, macros, directives, includes, and output formats.
- [Debugger user manual](docs/debugger.md) — loading programs, commands,
  disassembly, stepping, breakpoints, memory tools, and worked examples.
- [6502 instruction set manual](docs/instruction-set.md) — instruction
  behaviour, addressing modes, opcodes, byte counts, and cycle timings.

## Installation

### Prerequisites

- Go 1.26 or later

### Building from Source

```bash
git clone https://github.com/jrsteele09/go-6502-emulator.git
cd go-6502-emulator
go build -o asm6502 ./cmd/assembler
go build -o debug6502 ./cmd/debugger
```

Place them in an appropriate location for your system and add that location to
your command path.

## License

MIT License. See [LICENSE](LICENSE) for more information.
