#!/bin/bash

go build -o ./bin/asm6502 ./cmd/assembler
go build -o ./bin/debug6502 ./cmd/debugger

BIN_PATH="$(pwd)/bin"
case ":$PATH:" in
    *":$BIN_PATH:"*)
        # Already in PATH
        ;;
    *)
        export PATH="$PATH:$BIN_PATH"
        ;;
esac