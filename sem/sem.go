// Package sem implements a set of semantic analysis passes.
package sem

import (
	"uBasic/ast"
	"uBasic/sem/semcheck"
	"uBasic/sem/typecheck"
	"uBasic/types"
)

type infoTypes map[ast.Expression]types.Type
type scopeTypes map[ast.Node]*Scope

// Save performs a static semantic analysis check on the given file.
// and save the file in binary format.
func Check(file *ast.File) (*Info, error) {
	// Semantic analysis is done in two passes to allow for forward references.
	// Firstly, the global declarations are added to the file-scope. Secondly,
	// the global function declaration bodies are traversed to resolve
	// identifiers and deduce the types of expressions.

	// set parent of each node
	if err := SetParents(file); err != nil {
		return nil, err
	}
	// Identifier resolution.
	info := &Info{
		Types:  make(infoTypes),
		Scopes: make(scopeTypes),
	}
	if err := resolve(file, info.Scopes); err != nil {
		return nil, err
	}

	// Type-checking.
	if err := typecheck.Check(file, info.Types); err != nil {
		return nil, err
	}

	// Semantic analysis.
	if err := semcheck.Check(file); err != nil {
		return nil, err
	}
	return info, nil
}

// func Save(file *ast.File) (*Info, error) {
// 	// Semantic analysis is done in two passes to allow for forward references.
// 	// Firstly, the global declarations are added to the file-scope. Secondly,
// 	// the global function declaration bodies are traversed to resolve
// 	// identifiers and deduce the types of expressions.

// 	// set parent of each node
// 	if err := SetParents(file); err != nil {
// 		return nil, err
// 	}
// 	// Identifier resolution.
// 	info := &Info{
// 		Types:  make(infoTypes),
// 		Scopes: make(scopeTypes),
// 	}
// 	if err := resolve(file, info.Scopes); err != nil {
// 		return nil, err
// 	}

// 	// Type-checking.
// 	if err := typecheck.Check(file, info.Types); err != nil {
// 		return nil, err
// 	}

// 	// Semantic analysis.
// 	if err := semcheck.Check(file); err != nil {
// 		return nil, err
// 	}

// 	if err := SaveFile(file, info.Types); err != nil {
// 		return nil, err
// 	}

// 	return info, nil
// }

// Info holds semantic information of a type-checked program.
type Info struct {
	// Types maps expression nodes to types.
	Types map[ast.Expression]types.Type
	// Scopes maps nodes to the scope they define.
	//
	// The following nodes define scopes.
	//
	//    *ast.File
	//    *ast.FuncDecl
	//    *ast.ForStmt
	//    *ast.IfStmt
	Scopes map[ast.Node]*Scope
}
