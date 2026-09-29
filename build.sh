#!/usr/bin/env bash

set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUILD_VERSION="$(tr -d '[:space:]' < "$PROJECT_ROOT/VERSION")"
if [[ -z "$BUILD_VERSION" ]]; then
    echo "VERSION must not be empty" >&2
    exit 1
fi

VERSION_PACKAGE="github.com/jrsteele09/go-6502-emulator/internal/buildinfo.Version"
LINKER_FLAGS="-X ${VERSION_PACKAGE}=${BUILD_VERSION}"

mkdir -p "$PROJECT_ROOT/bin"
cd "$PROJECT_ROOT"
go build -ldflags "$LINKER_FLAGS" -o "$PROJECT_ROOT/bin/asm6502" ./cmd/assembler
go build -ldflags "$LINKER_FLAGS" -o "$PROJECT_ROOT/bin/debug6502" ./cmd/debugger

BIN_PATH="$PROJECT_ROOT/bin"
case ":$PATH:" in
    *":$BIN_PATH:"*)
        # Already in PATH
        ;;
    *)
        export PATH="$PATH:$BIN_PATH"
        ;;
esac
