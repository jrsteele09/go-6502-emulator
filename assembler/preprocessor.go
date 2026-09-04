package assembler

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

type sourceMacro struct {
	parameters []string
	body       []string
}

type sourceConstant struct {
	name       string
	expression string
}

type conditionalState struct {
	parentActive bool
	condition    bool
	elseSeen     bool
	deferred     bool
}

// preprocessSource expands the small, source-level features shared by several
// 6502 assemblers before the token stream is consumed by the two assembler passes.
func preprocessSource(source string, filename string) (string, error) {
	macros := make(map[string]sourceMacro)
	constants := make(map[string]int64)
	return processSourceLines(strings.Split(source, "\n"), filename, macros, constants, 0)
}

func processSourceLines(lines []string, filename string, macros map[string]sourceMacro, constants map[string]int64, depth int) (string, error) {
	if depth > 32 {
		return "", fmt.Errorf("source macro expansion exceeded maximum depth in %s", filename)
	}

	processor := &sourceProcessor{
		lines:     lines,
		filename:  filename,
		macros:    macros,
		constants: constants,
		depth:     depth,
		active:    true,
	}
	return processor.process()
}

type sourceProcessor struct {
	lines        []string
	filename     string
	macros       map[string]sourceMacro
	constants    map[string]int64
	conditionals []conditionalState
	output       []string
	active       bool
	depth        int
	comments     sourceCommentState
}

type sourceLineKind int

const (
	sourceLineOrdinary sourceLineKind = iota
	sourceLineIf
	sourceLineElse
	sourceLineEndif
	sourceLineEndMacro
	sourceLineMacroDefinition
	sourceLineInactive
	sourceLineMacroInvocation
	sourceLineConstant
)

func (p *sourceProcessor) process() (string, error) {
	for lineNumber := 0; lineNumber < len(p.lines); lineNumber++ {
		nextLine, err := p.processLine(lineNumber)
		if err != nil {
			return "", err
		}
		lineNumber = nextLine
	}
	if len(p.conditionals) != 0 {
		return "", fmt.Errorf("%s: unterminated conditional block", p.filename)
	}
	return strings.Join(p.output, "\n"), nil
}

func (p *sourceProcessor) processLine(lineNumber int) (int, error) {
	line := p.lines[lineNumber]
	code := p.comments.codePart(line)
	trimmed := strings.TrimSpace(code)
	words := strings.Fields(trimmed)
	if len(words) == 0 {
		p.output = append(p.output, line)
		return lineNumber, nil
	}

	kind, macroName, macroParameters := p.classifyLine(trimmed, words)
	switch kind {
	case sourceLineIf:
		return lineNumber, p.handleIf(trimmed, line, lineNumber)
	case sourceLineElse:
		return lineNumber, p.handleElse(line, lineNumber)
	case sourceLineEndif:
		return lineNumber, p.handleEndif(line, lineNumber)
	case sourceLineEndMacro, sourceLineInactive:
		return lineNumber, nil
	case sourceLineMacroDefinition:
		return p.handleMacroDefinition(lineNumber, macroName, macroParameters)
	case sourceLineMacroInvocation:
		_, err := p.handleMacroInvocation(trimmed, words, lineNumber)
		return lineNumber, err
	case sourceLineConstant:
		p.output = append(p.output, p.recordConstant(trimmed, line))
		return lineNumber, nil
	case sourceLineOrdinary:
		p.output = append(p.output, line)
		return lineNumber, nil
	}
	return lineNumber, fmt.Errorf("%s:%d: unhandled source line", p.filename, lineNumber+1)
}

func (p *sourceProcessor) classifyLine(trimmed string, words []string) (sourceLineKind, string, []string) {
	switch strings.ToLower(words[0]) {
	case "if":
		return sourceLineIf, "", nil
	case "else":
		return sourceLineElse, "", nil
	case "endif":
		return sourceLineEndif, "", nil
	case "endm":
		return sourceLineEndMacro, "", nil
	}

	if name, parameters, ok := parseMacroDefinition(trimmed); ok {
		return sourceLineMacroDefinition, name, parameters
	}
	if !p.active {
		return sourceLineInactive, "", nil
	}
	if p.isMacroInvocation(words) {
		return sourceLineMacroInvocation, "", nil
	}
	if p.isConstantAssignment(trimmed) {
		return sourceLineConstant, "", nil
	}
	return sourceLineOrdinary, "", nil
}

func (p *sourceProcessor) isMacroInvocation(words []string) bool {
	name := strings.TrimSuffix(words[0], ":")
	_, ok := p.macros[strings.ToUpper(name)]
	return ok && !strings.Contains(name, ":")
}

func (p *sourceProcessor) isConstantAssignment(line string) bool {
	_, ok := parseSourceConstant(line)
	return ok
}

func (p *sourceProcessor) recordConstant(codeLine string, originalLine string) string {
	constant, ok := parseSourceConstant(codeLine)
	if !ok {
		return originalLine
	}
	parsed, err := parseSourceInteger(constant.expression, p.constants)
	if err != nil {
		return originalLine
	}
	p.constants[constant.name] = parsed
	return fmt.Sprintf("%s = %d", constant.name, parsed)
}

func (p *sourceProcessor) handleIf(trimmed string, originalLine string, lineNumber int) error {
	condition, err := evaluateCondition(strings.TrimSpace(strings.TrimPrefix(trimmed, strings.Fields(trimmed)[0])), p.constants)
	if err != nil {
		if p.active {
			p.output = append(p.output, originalLine)
		}
		p.conditionals = append(p.conditionals, conditionalState{parentActive: p.active, condition: true, deferred: true})
		return nil
	}
	p.conditionals = append(p.conditionals, conditionalState{parentActive: p.active, condition: condition})
	p.active = p.active && condition
	return nil
}

func (p *sourceProcessor) handleElse(originalLine string, lineNumber int) error {
	if len(p.conditionals) == 0 {
		return fmt.Errorf("%s:%d: unexpected else", p.filename, lineNumber+1)
	}
	state := &p.conditionals[len(p.conditionals)-1]
	if state.elseSeen {
		return fmt.Errorf("%s:%d: duplicate else", p.filename, lineNumber+1)
	}
	state.elseSeen = true
	if state.deferred {
		if state.parentActive {
			p.output = append(p.output, originalLine)
		}
		p.active = state.parentActive
		return nil
	}
	p.active = state.parentActive && !state.condition
	return nil
}

func (p *sourceProcessor) handleEndif(originalLine string, lineNumber int) error {
	if len(p.conditionals) == 0 {
		return fmt.Errorf("%s:%d: unexpected endif", p.filename, lineNumber+1)
	}
	state := p.conditionals[len(p.conditionals)-1]
	if state.deferred && state.parentActive {
		p.output = append(p.output, originalLine)
	}
	p.conditionals = p.conditionals[:len(p.conditionals)-1]
	p.active = true
	for _, state := range p.conditionals {
		if state.deferred {
			p.active = p.active && state.parentActive
			continue
		}
		p.active = p.active && state.parentActive && ((state.condition && !state.elseSeen) || (!state.condition && state.elseSeen))
	}
	return nil
}

func (p *sourceProcessor) handleMacroDefinition(lineNumber int, name string, parameters []string) (int, error) {
	body := make([]string, 0)
	endLine := lineNumber + 1
	for ; endLine < len(p.lines); endLine++ {
		bodyLine := p.lines[endLine]
		bodyCode := p.comments.codePart(bodyLine)
		if strings.EqualFold(strings.TrimSpace(bodyCode), "endm") {
			break
		}
		body = append(body, bodyLine)
	}
	if endLine >= len(p.lines) {
		return lineNumber, fmt.Errorf("%s: unterminated macro %s", p.filename, name)
	}
	if p.active {
		p.macros[strings.ToUpper(name)] = sourceMacro{parameters: parameters, body: body}
	}
	return endLine, nil
}

func (p *sourceProcessor) handleMacroInvocation(trimmed string, words []string, lineNumber int) (bool, error) {
	name := strings.TrimSuffix(words[0], ":")
	macro, ok := p.macros[strings.ToUpper(name)]
	if !ok || strings.Contains(name, ":") {
		return false, nil
	}
	args := parseMacroArguments(strings.TrimSpace(strings.TrimPrefix(trimmed, words[0])))
	expanded, err := expandMacro(macro, args)
	if err != nil {
		return true, fmt.Errorf("%s:%d: %w", p.filename, lineNumber+1, err)
	}
	expandedSource, err := processSourceLines(expanded, p.filename, p.macros, p.constants, p.depth+1)
	if err != nil {
		return true, err
	}
	if expandedSource != "" {
		p.output = append(p.output, strings.Split(strings.TrimSuffix(expandedSource, "\n"), "\n")...)
	}
	return true, nil
}

func parseMacroDefinition(line string) (string, []string, bool) {
	words := strings.Fields(line)
	if len(words) == 0 {
		return "", nil, false
	}
	if strings.EqualFold(words[0], "macro") && len(words) >= 2 {
		macroStart := strings.Index(strings.ToLower(line), "macro") + len("macro")
		remainder := strings.TrimSpace(line[macroStart:])
		name := strings.Fields(remainder)[0]
		parameters := strings.TrimSpace(strings.TrimPrefix(remainder, name))
		return name, parseMacroParameters(parameters), true
	}
	if len(words) >= 2 && strings.EqualFold(words[1], "macro") {
		macroStart := strings.Index(strings.ToLower(line), "macro") + len("macro")
		return words[0], parseMacroParameters(strings.TrimSpace(line[macroStart:])), true
	}
	return "", nil, false
}

func parseMacroParameters(value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	parts := strings.FieldsFunc(value, func(r rune) bool { return r == ',' || r == ' ' || r == '\t' })
	return parts
}

func parseMacroArguments(value string) []string {
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func expandMacro(macro sourceMacro, arguments []string) ([]string, error) {
	if len(macro.parameters) > 0 && len(arguments) != len(macro.parameters) {
		return nil, fmt.Errorf("macro expects %d arguments, got %d", len(macro.parameters), len(arguments))
	}
	if len(macro.parameters) == 0 {
		requiredArguments := macro.requiredPositionalArguments()
		if len(arguments) < requiredArguments {
			return nil, fmt.Errorf("macro expects at least %d arguments, got %d", requiredArguments, len(arguments))
		}
	}

	result := make([]string, len(macro.body))
	for i, line := range macro.body {
		value := line
		for parameterIndex, parameter := range macro.parameters {
			argument := arguments[parameterIndex]
			pattern := regexp.MustCompile(`\\` + regexp.QuoteMeta(parameter) + `|\b` + regexp.QuoteMeta(parameter) + `\b`)
			value = pattern.ReplaceAllStringFunc(value, func(string) string { return argument })
		}
		for argumentIndex, argument := range arguments {
			value = strings.ReplaceAll(value, fmt.Sprintf(`\%d`, argumentIndex+1), argument)
		}
		result[i] = value
	}
	return result, nil
}

func (m sourceMacro) requiredPositionalArguments() int {
	maxArgument := 0
	pattern := regexp.MustCompile(`\\([1-9][0-9]*)`)
	for _, line := range m.body {
		for _, match := range pattern.FindAllStringSubmatch(line, -1) {
			argument, err := parseSourceInteger(match[1], nil)
			if err == nil && int(argument) > maxArgument {
				maxArgument = int(argument)
			}
		}
	}
	return maxArgument
}

func evaluateCondition(expression string, constants map[string]int64) (bool, error) {
	return EvaluateSourceCondition(expression, func(name string) (int64, bool) {
		value, ok := constants[name]
		return value, ok
	})
}

func parseSourceConstant(line string) (sourceConstant, bool) {
	if constant, ok := parseEqualsConstant(line); ok {
		return constant, true
	}
	return parseEquConstant(line)
}

func parseEqualsConstant(line string) (sourceConstant, bool) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return sourceConstant{}, false
	}
	name := strings.TrimSpace(parts[0])
	if !regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`).MatchString(name) {
		return sourceConstant{}, false
	}
	return sourceConstant{name: name, expression: strings.TrimSpace(parts[1])}, true
}

func parseEquConstant(line string) (sourceConstant, bool) {
	words := strings.Fields(line)
	if len(words) < 3 || !strings.EqualFold(words[1], "equ") {
		return sourceConstant{}, false
	}
	name := words[0]
	if !regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`).MatchString(name) {
		return sourceConstant{}, false
	}
	equOffset := strings.Index(strings.ToLower(line), "equ")
	if equOffset == -1 {
		return sourceConstant{}, false
	}
	return sourceConstant{name: name, expression: strings.TrimSpace(line[equOffset+len("equ"):])}, true
}

func parseSourceInteger(value string, constants map[string]int64) (int64, error) {
	return EvaluateSourceExpression(value, func(name string) (int64, bool) {
		constant, ok := constants[name]
		return constant, ok
	})
}

type sourceCommentState struct {
	inBlock bool
}

func (s *sourceCommentState) codePart(line string) string {
	var code strings.Builder
	for i := 0; i < len(line); {
		if s.inBlock {
			end := strings.Index(line[i:], "*/")
			if end == -1 {
				return code.String()
			}
			i += end + len("*/")
			s.inBlock = false
			code.WriteByte(' ')
			continue
		}

		if strings.HasPrefix(line[i:], "/*") {
			s.inBlock = true
			i += len("/*")
			code.WriteByte(' ')
			continue
		}
		if strings.HasPrefix(line[i:], "//") || line[i] == ';' {
			break
		}
		if line[i] == '"' || line[i] == '\'' {
			next := copyQuotedSource(&code, line, i)
			i = next
			continue
		}

		code.WriteByte(line[i])
		i++
	}
	return code.String()
}

func copyQuotedSource(dst *strings.Builder, line string, start int) int {
	quote := line[start]
	dst.WriteByte(line[start])
	i := start + 1
	for i < len(line) {
		dst.WriteByte(line[i])
		if line[i] == '\\' && i+1 < len(line) {
			i++
			dst.WriteByte(line[i])
		} else if line[i] == quote {
			i++
			break
		}
		i++
	}
	return i
}

func preprocessor(reader io.Reader, filename string) (io.Reader, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	processed, err := preprocessSource(string(data), filename)
	if err != nil {
		return nil, err
	}
	return strings.NewReader(processed), nil
}
