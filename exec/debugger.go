package exec

import (
	"fmt"
	"os"
	"strings"
	"time"
	"uBasic/ast"
	"uBasic/errors"
	"uBasic/eval"
	"uBasic/lexer"
	"uBasic/object"
	"uBasic/parser"
	"uBasic/sem"
	"uBasic/source"
)

type DebugCommand uint8

const (
	Continue DebugCommand = iota
	Step
	Stop
)

type Breakpoints []bool

type Debugger struct {
	// line number
	started     bool
	Breakpoints Breakpoints
	Env         *object.Environment
	lastCommand DebugCommand
	Source      *source.Source
	File        *ast.File
	Running     bool // pause the debugger
	Terminal    strings.Builder
	Info        *sem.Info
}

func (b *Breakpoints) String() string {
	var sb strings.Builder
	for i, v := range *b {
		if v {
			sb.WriteString(fmt.Sprintf("%d ", i))
		}
	}
	return sb.String()
}

var Debug = Debugger{} // keep a reference to the debugger

func (d *Debugger) Run() error {
	// run the interpreter
	if !d.started {
		d.started = true
		result := eval.Run(d.File, d.Env, callback)
		d.started = false
		if result != nil {
			switch result := result.(type) {
			case *object.Error:
				return fmt.Errorf(result.String())
			case *object.Nothing:
				return fmt.Errorf("\nprogram ended with Nothing")
			default:
				return fmt.Errorf("\nprogram ended with value: %v", result.String())
			}
		} else {
			return fmt.Errorf("\nprogram ended with no returned value")
		}
	}
	return nil
}

func (d *Debugger) Stop() {
	d.lastCommand = Stop
	d.Running = true
}

func (d *Debugger) Continue() {
	d.lastCommand = Continue
	d.Running = true
}

func (d *Debugger) Step() {
	d.lastCommand = Step
	d.Running = true
}

func callback(node ast.Node, env *object.Environment) bool {
	Debug.Env = env
	needToStop := false
	// check if we have a breakpoint
	if node != nil && node.Token() != nil {
		// check if we have a breakpoint
		token := node.Token()
		if token != nil {
			needToStop = Debug.Breakpoints[token.Position.Line]
		}
	}

	if needToStop {
		Debug.Running = false
	}
	// wait till we get a command
	for !Debug.Running {
		// wait for 0,25 seconds
		time.Sleep(time.Second / 4)
	}

	// read last command
	switch Debug.lastCommand {
	case Continue:
		Debug.Running = true
		return true // continue
	case Step:
		Debug.Running = false
		return true // continue
	case Stop:
		// stop
		Debug.Running = true
		return false // stop
	}
	// wait for
	return true // continue
}

func LoadInterpreter(fileName string) error {
	if fileName == "" {
		return fmt.Errorf("missing file name")
	}
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
	src := &source.Source{Input: input, Name: fileName}
	if file == nil {
		for _, err := range p.Errors() {
			e := err.(*errors.Error)
			e.Source = src
			fmt.Println(err)
		}
	} else {
		info, err := sem.Check(file)
		if err != nil {
			e := err.(*errors.Error)
			e.Source = src
			fmt.Println(err)
		} else {
			file.Name = fileName
			Debug.Info = info
			Debug.Terminal = strings.Builder{}
			Debug.Env = eval.Define(info, os.Stdin, &Debug.Terminal)
			Debug.Source = &source.Source{Input: input, Name: fileName}
			Debug.Breakpoints = make([]bool, Debug.Source.LineCount()+1)
			Debug.lastCommand = Step // start with a step
			Debug.File = file
			return nil
		}
	}
	return fmt.Errorf("error in file %s", fileName)
}
