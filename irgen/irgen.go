// Package irgen implements a µC to LLVM IR generator.
package irgen

import (
	"fmt"
	"log"
	"os"

	"uBasic/ast"
	"uBasic/object"
	"uBasic/sem"

	"github.com/llir/llvm/ir"
	irtypes "github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mewkiz/pkg/term"
)

// TODO: Remove debug output.
// dbg is a logger which prefixes debug messages with "irgen:".
var dbg = log.New(os.Stdout, term.WhiteBold("irgen:"), log.Lshortfile)

// A Module represents an LLVM IR module generator.
type Module struct {
	*ir.Module
	// info holds semantic information about the program from the type-checker.
	info *sem.Info
	// Maps from identifier source code position to the associated value.
	idents map[int]value.Value
	// SkipLocalVariables is a indicator to skip local variable declaration.
	SkipLocalVariables bool
	// constant values
	env *object.Environment
}

// NewModule returns a new module generator.
func NewModule(info *sem.Info) *Module {
	return &Module{
		Module: ir.NewModule(),
		info:   info,
		idents: make(map[int]value.Value),
		env:    object.NewEnvironment(),
	}
}

// emitFunc emits to m the given function.
func (m *Module) emitFunc(f *Function) {
	m.Funcs = append(m.Funcs, f.Func)
}

// emitGlobal emits to m the given global variable declaration.
func (m *Module) emitGlobal(global *ir.Global) {
	m.Globals = append(m.Globals, global)
}

// A Function represents an LLVM IR function generator.
type Function struct {
	// Function being generated.
	*ir.Func
	// Current basic block being generated.
	currentBlock *Block
	// Maps from identifier source code position to the associated value.
	idents map[int]value.Value
	// Map of existing local variable names.
	exists map[string]bool
}

// NewFunc returns a new function generator based on the given function name and
// signature.
//
// The caller is responsible for initializing basic blocks.
func NewFunc(name string, retType irtypes.Type, params ...*ir.Param) *Function {
	f := ir.NewFunc(name, retType, params...)
	return &Function{Func: f, idents: make(map[int]value.Value), exists: make(map[string]bool)}
}

// startBody initializes the generation of the function body.
func (f *Function) startBody() {
	entry := f.NewBlock("") // "entry"
	f.currentBlock = entry
}

// endBody finalizes the generation of the function body.
func (m *Module) endBody(f *Function, resultIdentifier *ast.Identifier) error {
	if block := f.currentBlock; block != nil && block.Term == nil {
		switch {
		case f.Func.Name() == "main":
			// From C11 spec $5.1.2.2.3.
			//
			// "If the return type of the main function is a type compatible with
			// int, a return from the initial call to the main function is
			// equivalent to calling the exit function with the value returned by
			// the main function as its argument; reaching the } that terminates
			// the main function returns a value of 0."

			m.GCstop(f)
			result := f.Sig.RetType
			zero := constZero(result)
			termRet := ir.NewRet(zero)
			block.SetTerm(termRet)
		default:
			// Add void return terminator to the current basic block, if a
			// terminator is missing.
			switch result := f.Sig.RetType; {
			case result.Equal(irtypes.Void):
				termRet := ir.NewRet(nil)
				block.SetTerm(termRet)
			default:
				declPos := resultIdentifier.Decl.Name().Token().Position.Absolute

				// is it local or global
				if value1, ok := f.idents[declPos]; ok {
					variableValue := f.currentBlock.NewLoad(result, value1)
					termRet := ir.NewRet(variableValue)
					block.SetTerm(termRet)
				} else {
					panic("unknown identifier")
				}
			}
		}
	}
	f.currentBlock = nil
	return nil
}

// emitLocal emits to f the given named value instruction.
func (f *Function) emitLocal(ident *ast.Identifier, inst valueInst) value.Value {
	return f.currentBlock.emitLocal(ident, inst)
}

// A Block represents an LLVM IR basic block generator.
type Block struct {
	// Basic block being generated.
	*ir.Block
	// Parent function of the basic block.
	parent *Function
}

// NewBlock returns a new basic block generator based on the given name and
// parent function.
func (f *Function) NewBlock(name string) *Block {
	block := ir.NewBlock(name)
	return &Block{Block: block, parent: f}
}

// valueInst represents an instruction producing a value.
type valueInst interface {
	ir.Instruction
	value.Named
}

// emitLocal emits to b the given named value instruction.
func (b *Block) emitLocal(ident *ast.Identifier, inst valueInst) value.Value {
	name := b.parent.genUnique(ident)
	inst.SetName(name)
	b.parent.setIdentValue(ident, inst)
	return inst
}

// SetTerm sets the terminator of the basic block.
func (b *Block) SetTerm(term ir.Terminator) {
	if b.Term != nil {
		panic(fmt.Sprintf("terminator instruction already set for basic block; old term (%v), new term (%v), basic block (%v)", term, b.Term, b))
	}
	b.Block.Term = term
	b.parent.Blocks = append(b.parent.Blocks, b.Block)
}

// changeBlock changes the current basic block of the function to the basic
// block of the given statement.
func (f *Function) changeBlock(block *Block) {
	f.Blocks = append(f.Blocks, f.currentBlock.Block)
	f.currentBlock = block
}

// LookupFunction returns the LLVM IR function associated with the given name.
func (m *Module) LookupFunction(name string) *ir.Func {
	for _, f := range m.Funcs {
		if f.Name() == name {
			return f
		}
	}
	return nil
}

// LookupGlobal returns the LLVM IR global variable associated with the given
// name.
func (m *Module) LookupGlobal(name string) *ir.Global {
	for _, global := range m.Globals {
		if global.Name() == name {
			return global
		}
	}
	return nil
}
