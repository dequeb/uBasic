package sem

import (
	"uBasic/ast"
	"uBasic/ast/astutil"
	"uBasic/errors"
	"uBasic/types"

	"github.com/iancoleman/strcase"
)

// A Scope maintains the set of named language entities declared in the lexical
// scope and a link to the immediately surrounding outer scope.
type Scope struct {
	// Immediately surrounding outer scope; or nil if universe scope.
	Outer *Scope
	// Identifiers declared within the current scope.
	Decls map[string]ast.Decl
	// IsDef reports whether the given declaration is a definition.
	IsDef func(ast.Node) bool
}

// NewScope returns a new lexical scope immediately surrouded by the given outer
// scope.
func NewScope(outer *Scope) *Scope {
	return &Scope{
		Outer: outer,
		Decls: make(map[string]ast.Decl),
		IsDef: astutil.IsDef,
	}
}

// Insert inserts the given declaration into the current scope.
func (s *Scope) Insert(decl ast.Decl) error {
	// Early return for first-time declarations.
	ident := decl.Name()
	if ident == nil {
		// Anonymous function parameter declaration.
		return nil
	}
	name := ident.String()
	prev, ok := s.Decls[name]
	if !ok {
		s.Decls[name] = decl
		return nil
	}
	prevIdent := prev.Name()

	if prevIdent.Name == ident.Name {
		// Identifier already added to scope.
		return nil
	}

	// Previously declared.
	prevType, err := prev.Type()
	if err != nil {
		return err
	}
	declType, err := decl.Type()
	if err != nil {
		return err
	}

	if !types.Equal(prevType, declType) {
		return errors.Newf(ident.Token().Position, "redefinition of %q with type %q instead of %q", name, declType, prevType)
	}

	// The last tentative definition becomes the definition, unless defined
	// explicitly (e.g. having an initializer or function body).
	if !s.IsDef(prev) {
		s.Decls[name] = decl
		return nil
	}

	// Definition already present in scope.
	if s.IsDef(decl) {
		// TODO: Consider adding support for multiple errors (and potentially
		// warnings and notifications).
		//
		// If support for notifications are added, add a note of the previous declaration.
		//    errors.Notef(prevIdent.Start(), "previous definition of %q", name)
		return errors.Newf(ident.Token().Position, "redefinition of %q", name)
	}

	// Declaration of previously declared identifier.
	return nil
}

// Lookup returns the declaration of name in the innermost scope of s. The
// returned boolean variable reports whether a declaration of name was located.
func (s *Scope) Lookup(name string, currentOnly bool) (ast.Decl, bool) {
	if decl, ok := s.Decls[name]; ok {
		return decl, true
	}
	// check with camel case for native type
	if decl, ok := s.Decls[strcase.ToCamel(name)]; ok {
		return decl, true
	}

	if s.Outer == nil {
		return nil, false
	}
	if !currentOnly {
		return s.Outer.Lookup(name, false)
	}
	return nil, false
}
