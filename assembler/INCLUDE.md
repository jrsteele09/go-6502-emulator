# Assembler Include Directive Support

The assembler now supports file inclusion using multiple directive formats, providing robust and flexible code organization capabilities.

## Supported Include Directives

### 1. Hash Include (C-style)
```assembly
#include "filename.asm"
#include 'filename.asm'
#include filename.asm
```

### 2. Dot Include (Traditional assembler)
```assembly
.include "filename.asm"
.include 'filename.asm'
.include filename.asm
```

### 3. Uppercase Dot Include
```assembly
.INCLUDE "filename.asm"
.INCLUDE 'filename.asm'
.INCLUDE filename.asm
```

## Usage

### Basic Usage

```go
// Create assembler
asm := assembler.New(opcodes)

// For filesystem-based includes
resolver := assembler.NewOSFileResolver("/path/to/base/directory")
segments, err := asm.AssembleWithPreprocessor(reader, resolver)

// For memory-based includes (testing)
files := map[string]string{
    "main.asm": "LDA #$10\n#include \"lib.asm\"",
    "lib.asm":  "CMP #$20",
}
resolver := assembler.NewMemoryFileResolver(files)
segments, err := asm.AssembleWithPreprocessor(reader, resolver)
```

### Backward Compatibility

The original `Assemble` method continues to work without preprocessing:

```go
segments, err := asm.Assemble(reader) // No include support
```

## File Resolvers

### OSFileResolver
Resolves files from the filesystem relative to a base directory.

```go
resolver := assembler.NewOSFileResolver("/projects/my6502/src")
```

### MemoryFileResolver
Resolves files from an in-memory map (useful for testing or embedded content).

```go
files := map[string]string{
    "main.asm": "...",
    "lib.asm":  "...",
}
resolver := assembler.NewMemoryFileResolver(files)
```

### Custom Resolvers
Implement the `FileResolver` interface for custom file resolution:

```go
type FileResolver interface {
    Resolve(path string) (io.Reader, error)
}
```

## Features

### Nested Includes
Files can include other files, creating a tree of dependencies:

```assembly
; main.asm
#include "graphics.asm"
#include "sound.asm"

; graphics.asm
#include "sprites.asm"
#include "background.asm"
```

### Circular Include Detection
The preprocessor automatically detects and prevents circular includes:

```assembly
; file1.asm
#include "file2.asm"

; file2.asm
#include "file1.asm"  ; ERROR: Circular include detected
```

### Maximum Depth Protection
Prevents infinite recursion with configurable maximum include depth:

```go
preprocessor := assembler.NewPreprocessor(resolver)
preprocessor.SetMaxDepth(20) // Default is 10
```

### Path Resolution
- Quoted paths: `#include "path/to/file.asm"`
- Single-quoted paths: `#include 'path/to/file.asm'`
- Unquoted paths: `#include path/to/file.asm`
- Relative paths resolved from base directory
- Absolute paths used as-is

## Error Handling

The preprocessor provides detailed error messages:

- **File not found**: `"include file 'missing.asm' not found"`
- **Circular includes**: `"circular include detected: 'file.asm'"`
- **Maximum depth**: `"maximum include depth (10) exceeded"`
- **Line numbers**: `"line 42: failed to open include file"`

## Example

```assembly
; main.asm
    .ORG $8000
    
start:
    LDA #$00
    #include "init.asm"
    
main_loop:
    .include "input.asm"
    .INCLUDE "graphics.asm"
    JMP main_loop

; init.asm
    STA $0200    ; Clear sprite 0
    STA $0201
    
; input.asm  
    LDA $4016    ; Read controller
    AND #$01
    
; graphics.asm
    LDA #$3F     ; Load palette
    STA $2006
```

This assembles to a single program with all includes expanded inline.

## Architecture

The include system uses a two-phase approach:

1. **Preprocessing Phase**: Expands all includes before tokenization
2. **Assembly Phase**: Normal two-pass assembly on expanded source

This design ensures that includes are processed at the source level, maintaining clean separation between file handling and assembly logic.

## Testing

Comprehensive test coverage includes:
- Basic include functionality
- Nested includes
- Circular include detection
- Error conditions
- All directive formats
- Path resolution scenarios
- Integration with assembler

Run tests with:
```bash
go test -v -run TestPreprocessor
go test -v -run TestAssembler_
```
