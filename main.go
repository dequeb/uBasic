package main

import (
	"fmt"
	"os"
	"strings"
	"uBasic/ast"
	"uBasic/errors"
	"uBasic/eval"
	"uBasic/lexer"
	"uBasic/object"
	"uBasic/parser"
	"uBasic/sem"
	"uBasic/source"
)

type DebugCommand int

const (
	Continue DebugCommand = iota
	Step
)

type debuger struct {
	// line number
	breakpoints   []bool
	Env           *object.Environment
	LastCommand   DebugCommand
	ContinueCount int
	source        *source.Source
	File          *ast.File
}

var debug debuger

const prompt = "?: Display this line,  E: environnement, C: continue, D: Display, S: step, B: Breakpoint, Q: quit >>"

func callback(node ast.Node) bool {
	needToStop := false
	// check if we have a breakpoint
	if node != nil && node.Token() != nil {
		// check if we have a breakpoint
		needToStop = debug.breakpoints[node.Token().Position.Line]
	}

	// check if we have a command
	if debug.LastCommand == Step {
		needToStop = true
	}

	for needToStop {
		// display current line
		// fmt.Printf("% 4d: %s", node.Token().Position.Line, debug.source.Line(node.Token().Position))
		line := debug.source.Line(node.Token().Position)
		line = strings.Trim(line, "\n\r")
		lineNumber := node.Token().Position.Line
		fmt.Printf("% 3d: %s", lineNumber, line)

		response := ""
		fmt.Scanln(&response)
		response = strings.ToUpper(response)
		if len(response) == 0 {
			response = "S" // enter by default will step
		}

		// return false to stop, true to continue
		switch response[0] {
		case '?':
			fmt.Print(prompt)
		case 'E':
			fmt.Println(debug.Env.String())
		case 'C':
			debug.LastCommand = Continue
			return true
		case 'D':
			fmt.Println(debug.source.WithLineNumbers())
		case 'B':
			line := 0

			if len(response) > 1 {
				fmt.Sscanf(response[1:], "%d", &line)
				debug.breakpoints[line] = true
			} else {
				fmt.Println("Enter the line number")
				fmt.Scanln(&line)
				debug.breakpoints[line] = true
			}
		case 'S':
			debug.LastCommand = Step
			return true
		case 'Q':
			debug.LastCommand = Continue
			return false
		}
	}
	return true // continue
}

func (d *debuger) LoadInterpreter(fileName string) error {
	// read the file content
	filebytes, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	input := string(filebytes)

	// read the file into the interpreter
	l := lexer.New(input)
	p := parser.New(l)
	file := p.ParseFile()
	src := &source.Source{Input: input, Name: "filename"}
	if file == nil {
		for _, err := range p.Errors() {
			e := err.(*errors.Error)
			e.Src = src
			fmt.Println(err)
		}
	} else {
		info, err := sem.Check(file)
		if err != nil {
			e := err.(*errors.Error)
			e.Src = src
			fmt.Println(err)
		} else {
			file.Name = fileName
			d.Env = eval.Define(info)
			d.breakpoints = make([]bool, 100)
			d.LastCommand = Step // start with a step
			d.File = file
			d.source = &source.Source{Input: input, Name: fileName}
			return nil
		}
	}
	return fmt.Errorf("error in file %s", fileName)
}

func main() {
	// read command line to get the file name
	if len(os.Args) > 1 {

		fileName := os.Args[1]
		// fileName := "testdata/samples/fibo.bas"
		// read the file into the interpreter
		err := debug.LoadInterpreter(fileName)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(prompt)
			// run the interpreter
			eval.Run(debug.File, debug.Env, callback)
		}
	} else {
		fmt.Println("Usage: uBasic filename")
	}

}
