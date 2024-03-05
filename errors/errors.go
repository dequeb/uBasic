// Package errors provides pretty-printing of semantic analysis errors.
package errors

import (
	"fmt"
	"strings"
	"uBasic/source"
	"uBasic/token"

	"github.com/mewkiz/pkg/term"
)

// UseColor indicates if error messages should use colors.
var UseColor = true

// An Error represents a semantic analysis error.
type Error struct {
	// Input source position (in bytes).
	Pos token.Position
	// Error message.
	Text string
	// Input source.
	Src *source.Source
}

// New returns a new error based on the given positional information (offset in
// bytes).
func New(pos token.Position, text string) *Error {
	err := &Error{
		Pos:  pos,
		Text: text,
	}
	return err
}

// Newf returns a new formatted error based on the given positional information
// (offset in bytes).
func Newf(pos token.Position, format string, a ...interface{}) *Error {
	err := &Error{
		Pos:  pos,
		Text: fmt.Sprintf(format, a...),
	}
	return err
}

// Error returns an error string with position information.
//
// The error format is as follows.
//
//	(file:line:column): error: text
func (e *Error) Error() string {
	// sanity check
	if e.Pos.Absolute == 0 && (e.Pos.Line != 1 || e.Pos.Column != 1) {
		// invalid position
		panic("invalid position of error " + e.Text + " " + e.Pos.String())
	}
	// Use colors.
	pos := e.Pos.String()
	prefix := "error:"
	text := e.Text
	if UseColor {
		pos = term.Color(pos, term.Bold)
		prefix = term.RedBold(prefix)
		text = term.Color(text, term.Bold)
	}
	src := e.Src
	if src == nil {
		// If Src is nil, the error format is as follows.
		//
		//    (byte offset %d) error: text
		return fmt.Sprintf("%s %s %s", pos, prefix, text)
	}
	// The error format is as follows.
	//
	//    (line) error: text
	//       1 = y
	//         ^
	line := e.Pos.Line
	col := e.Pos.Column
	name := src.Name
	srcLine := src.Line(e.Pos)
	srcLine = strings.Replace(srcLine, "\t", " ", -1)
	srcLine = strings.TrimRight(srcLine, "\n\r")
	arrow := fmt.Sprintf("%*s", col, "⩓")
	pos = fmt.Sprintf("(line %d)", line)
	if UseColor {
		pos = term.Color(pos, term.Bold)
		arrow = term.Color(arrow, term.Bold)
	}
	return fmt.Sprintf("%s: %s %s %s\n%s\n%s", name, pos, prefix, text, srcLine, arrow)
}
