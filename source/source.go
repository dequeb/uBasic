package source

import (
	"fmt"
	"strings"
	"uBasic/token"
)

// A Source represents an input source.
type Source struct {
	// Input source text.
	Input string
	// file name
	Name string
}

// return the source line where offset is located
func (s *Source) Line(position token.Position) string {
	// universePosition
	if position.Absolute < 0 {
		return ""
	}
	// find begin of line
	begin := 0
	for i := position.Absolute - 1; i > 0; i-- {
		if s.Input[i] == '\n' {
			begin = i
			break
		}
	}
	// find end of line
	end := len(s.Input) - 1
	for i := position.Absolute; i < len(s.Input); i++ {
		if s.Input[i] == '\n' {
			end = i
			break
		}
	}

	return strings.Trim(s.Input[begin:end+1], "\n\r")
}

// return source with line numbers
func (s *Source) WithLineNumbers() string {
	lines := strings.Split(s.Input, "\n")
	for i, line := range lines {
		lines[i] = fmt.Sprintf("% 6d %s\n", i+1, strings.TrimRight(line, "\n\r"))
	}
	return strings.Join(lines, "")
}

// LineCount returns the number of lines in the source.
func (s *Source) LineCount() int {
	return strings.Count(s.Input, "\n") + 1
}
