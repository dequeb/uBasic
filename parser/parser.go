package parser

import (
	"strings"
	"uBasic/ast"
	"uBasic/errors"
	"uBasic/lexer"
	"uBasic/token"
)

// first in first out
type Fifo struct {
	top  int
	data []token.Token
}

func (f *Fifo) push(t *token.Token) {
	f.data = append(f.data, *t)
}

func (f *Fifo) pop() *token.Token {
	const MAX_FIFO_SIZE = 30

	t := f.data[f.top]
	f.top++
	if f.top >= MAX_FIFO_SIZE {
		// remove all elements except last
		f.data = f.data[f.top:]
		f.top = 0
	}
	return &t
}

func (f *Fifo) peek(n int) *token.Token {
	if f.top+n >= len(f.data) {
		return &token.Token{Kind: token.EOF}
	}
	return &f.data[f.top+n]
}

func (f *Fifo) size() int {
	return len(f.data) - f.top
}

// Parser is a parser for the µBasic programming language.
type Parser struct {
	l      *lexer.Lexer
	fifo   *Fifo
	errors []error
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.fifo = &Fifo{top: 0, data: []token.Token{}}

	// Read four tokens in fifo structure
	p.loadFifo(2) // same behavior as before (load 2 tokens in fifo structure)
	return p
}

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) AddError(tok *token.Token, message string) {
	p.errors = append(p.errors, errors.New(tok.Position, message))
}

func (p *Parser) nextToken() {
	p.fifo.pop()
	p.fifo.push(p.l.NextToken())
	// skip virtual EOL token
	if p.curTokenIs(token.Underscore) && p.peekTokenIs(1, token.EOL) {
		p.fifo.pop() // skip the underscore
		p.fifo.pop() // skip the EOL
	}
}

func (p *Parser) loadFifo(n int) {
	for i := 0; i < n; i++ {
		t := p.l.NextToken()
		if t.Kind != token.EOF {
			p.fifo.push(t)
		}
	}
}

func (p *Parser) peekToken(n int) *token.Token {
	if p.fifo.size() <= n {
		p.loadFifo(n - p.fifo.size() + 1)
	}
	return p.fifo.peek(n)
}

// func (p *Parser) expectPeek(n int, k token.Kind) bool {
// 	if p.peekToken(n).Kind == k {
// 		p.nextToken()
// 		return true
// 	} else {
// 		return false
// 	}
// }

func (p *Parser) currentToken() *token.Token {
	return p.fifo.peek(0)
}

func (p *Parser) curTokenIs(k token.Kind) bool {
	return p.currentToken().Kind == k
}

func (p *Parser) peekTokenIs(n int, k token.Kind) bool {
	if n >= p.fifo.size() {
		p.loadFifo(n - p.fifo.size() + 1)
	}
	return p.peekToken(n).Kind == k
}

// ----------------------------------------------------------------------------

func (p *Parser) ParseFile() *ast.File {
	file := &ast.File{}
	file.StatementLists = []ast.StatementList{}

	for p.currentToken().Kind != token.EOF {
		stmtList := p.ParseStatementList(false, false)
		if stmtList != nil {
			file.StatementLists = append(file.StatementLists, *stmtList)
		} else {
			return nil
		}
		if !(p.curTokenIs(token.EOL) || p.curTokenIs(token.EOF)) {
			p.AddError(p.currentToken(), "expecting EOL")
			return nil
		}
		p.nextToken() // skip EOL
	}
	return file
}

func (p *Parser) ParseStatementList(inSub bool, inFunction bool) *ast.StatementList {
	statementList := &ast.StatementList{}
	statementList.Statements = []ast.Statement{}
	for !p.isEndOfStatementlist() {
		stmt := p.ParseStatement(inSub, inFunction)
		// unable to test nil conversion error: see https://go101.org/article/nil.html
		isNil := ast.IsNil(stmt)
		if !isNil {
			statementList.Statements = append(statementList.Statements, stmt)
			if p.curTokenIs(token.Colon) {
				p.nextToken()
			}
		} else {
			p.AddError(p.currentToken(), "unexpected statement: "+p.currentToken().Literal)
			return nil
		}
	}
	return statementList
}

func (p *Parser) ParseStatement(inSub bool, inFunction bool) ast.Statement {
	switch p.currentToken().Kind {
	case token.KwDim:
		return p.ParseDimStatement()
	case token.KwConst:
		return p.ParseConstStatement()
	case token.KwEnum:
		return p.ParseEnumStatement()
	case token.KwFunction:
		if !(inSub || inFunction) {
			return p.ParseFunctionStatement()
		}
		p.AddError(p.currentToken(), "does not expect function")
		return nil
	case token.KwSub:
		if !(inSub || inFunction) {
			return p.ParseSubStatement()
		}
		p.AddError(p.currentToken(), "does not expect sub")
		return nil
	case token.EOL:
		return p.ParseEmptyStatement()
	case token.KwIf:
		return p.ParseIfStatement(inSub, inFunction)
	case token.KwFor:
		return p.ParseForStatement(inSub, inFunction)
	case token.KwStop, token.KwError, token.KwResume, token.KwGoto:
		return p.ParseSpecialStatement()
	case token.KwSelect:
		return p.ParseSelectStatement(inSub, inFunction)
	case token.KwDo:
		return p.ParseDoStatement(inSub, inFunction)
	case token.KwExit:
		return p.ParseExitStatement(inSub, inFunction)
	case token.KwLet:
		return p.ParseExpressionStatement()
	case token.KwCall:
		return p.ParseSubroutineCall()
	default:
		return p.ParseSpecialStatement()
	}
}

// ----------------------------------------------------------------------------
// Variable declaration
// ----------------------------------------------------------------------------

func (p *Parser) ParseDimStatement() *ast.DimDecl {
	stmt := &ast.DimDecl{DimKw: p.currentToken()}
	stmt.Vars = make([]ast.VarDecl, 0)
	stmt.DimKw = p.currentToken()
	if !p.curTokenIs(token.KwDim) {
		p.AddError(p.currentToken(), "expected Dim")
		return nil
	}
	p.nextToken()

	for {
		if !p.curTokenIs(token.Identifier) {
			p.AddError(p.currentToken(), "expected identifier")
			return nil
		}

		var varDecl ast.VarDecl
		if p.peekTokenIs(1, token.Lparen) { // array declaration -- look for parenthesis after the identifier
			varDecl = p.ParseArrayDecl()
			// nil conversion error: see https://go101.org/article/nil.html
			if varDecl.(*ast.ArrayDecl) == nil {
				return nil
			}
		} else {
			// scalar declaration
			varDecl = p.ParseScalarDecl()
			// nil conversion error: see https://go101.org/article/nil.html
			if varDecl.(*ast.ScalarDecl) == nil {
				return nil
			}
		}
		stmt.Vars = append(stmt.Vars, varDecl)

		// look for a comma to see if there are more variables
		if !p.curTokenIs(token.Comma) {
			break
		}
		p.nextToken() // skip the comma
	}
	if !p.isEndOfStatement() {
		p.AddError(p.currentToken(), "expected EOL")
		return nil
	}
	return stmt
}

func (p *Parser) isEndOfStatement() bool {
	return p.isEndOfStatementlist() || p.curTokenIs(token.Colon)
}

func (p *Parser) isEndOfStatementlist() bool {
	return p.curTokenIs(token.EOL) || p.curTokenIs(token.EOF)
}

func (p *Parser) ParseScalarDecl() *ast.ScalarDecl {
	// read identifier
	ident := p.ParseIdentifier()
	if ident == nil {
		// p.AddError(p.currentToken(), "expected identifier")
		return nil
	}
	// skip as keyword
	if !p.curTokenIs(token.KwAs) {
		p.AddError(p.currentToken(), "expected 'as'")
		return nil
	}
	p.nextToken()

	// read type
	typ := p.ParseType()
	if typ == nil {
		// p.AddError(p.currentToken(), "expected type")
		return nil
	}
	return &ast.ScalarDecl{VarName: ident, VarType: typ}
}

func (p *Parser) ParseArrayDecl() *ast.ArrayDecl {
	// read identifier
	ident := p.ParseIdentifier()
	if ident == nil {
		// p.AddError(p.currentToken(), "expected identifier")
		return nil
	}
	// skip left parenthesis
	if !p.curTokenIs(token.Lparen) {
		p.AddError(p.currentToken(), "expected left parenthesis")
		return nil
	}

	// read array size
	typ := p.ParseArrayType()
	if typ == nil {
		// p.AddError(p.currentToken(), "expected array type")
		return nil
	}

	// skip as keyword
	if !p.curTokenIs(token.KwAs) {
		p.AddError(p.currentToken(), "expected 'as'")
		return nil
	}
	p.nextToken()

	// read type
	typ.Type = p.ParseType()
	return &ast.ArrayDecl{VarName: ident, VarType: typ}
}

func (p *Parser) ParseArrayType() *ast.ArrayType {
	arrayType := &ast.ArrayType{}

	// get left parenthesis
	if !p.curTokenIs(token.Lparen) {
		p.AddError(p.currentToken(), "expected left parenthesis")
		return nil
	}
	lparen := p.currentToken()
	p.nextToken()

	arrayType.Dimensions = make([]ast.Expression, 0)
	arrayType.Lparen = lparen
	// read array size
	for {
		var size ast.Expression
		if !p.curTokenIs(token.Rparen) {
			// read size
			size = p.ParseExpression()
			if size == nil {
				// p.AddError(p.currentToken(), "expected expression")
				return nil
			}
		} else {
			break
		}

		// add size to the list
		arrayType.Dimensions = append(arrayType.Dimensions, size)

		// look for a comma to see if there are more sizes
		if !p.curTokenIs(token.Comma) {
			break
		}
		p.nextToken() // skip the comma
	}

	// get right parenthesis
	if !p.curTokenIs(token.Rparen) {
		p.AddError(p.currentToken(), "expected right parenthesis")
		return nil
	}
	arrayType.Rparen = &p.currentToken().Position
	p.nextToken() // skip the right parenthesis
	return arrayType
}

func (p *Parser) ParseType() ast.Type {
	tok := p.currentToken()
	var typ ast.Type
	switch tok.Kind {
	case token.KwLong, token.KwInteger, token.KwSingle, token.KwDouble, token.KwString, token.KwBoolean, token.KwDate, token.KwCurrency, token.KwVariant:
		typ = &ast.Identifier{Tok: tok, Name: tok.Literal}
	default:
		// user defined type
		identifier := p.ParseIdentifier()
		if identifier == nil {
			// p.AddError(p.currentToken(), "expected type")
			return nil
		}
		return identifier
	}
	p.nextToken()
	return typ
}

func (p *Parser) ParseIdentifier() *ast.Identifier {
	if !p.curTokenIs(token.Identifier) {
		p.AddError(p.currentToken(), "expected identifier")
		return nil
	}

	t := &ast.Identifier{Tok: p.currentToken(), Name: p.currentToken().Literal}
	p.nextToken()
	return t
}

// ----------------------------------------------------------------------------
// Constants declaration
// ----------------------------------------------------------------------------
func (p *Parser) ParseConstStatement() *ast.ConstDecl {
	stmt := &ast.ConstDecl{ConstKw: p.currentToken()}
	stmt.Consts = make([]ast.ConstDeclItem, 0)
	stmt.ConstKw = p.currentToken()
	if !p.curTokenIs(token.KwConst) {
		p.AddError(p.currentToken(), "expected Const")
		return nil
	}
	p.nextToken()

	for {
		if !p.curTokenIs(token.Identifier) {
			p.AddError(p.currentToken(), "expected identifier")
			return nil
		}

		constDecl := p.ParseConstDeclItem()
		if constDecl == nil {
			// p.AddError(p.currentToken(), "expected constant declaration")
			return nil
		}
		stmt.Consts = append(stmt.Consts, *constDecl)

		// look for a comma to see if there are more variables
		if !p.curTokenIs(token.Comma) {
			break
		}
		p.nextToken() // skip the comma
	}
	if !p.isEndOfStatement() {
		p.AddError(p.currentToken(), "expected EOL")
		return nil
	}
	return stmt
}

func (p *Parser) ParseConstDeclItem() *ast.ConstDeclItem {
	// read identifier
	ident := p.ParseIdentifier()
	if ident == nil {
		// p.AddError(p.currentToken(), "expected identifier")
		return nil
	}
	// skip as keyword
	if !p.curTokenIs(token.KwAs) {
		p.AddError(p.currentToken(), "expected 'as'")
		return nil
	}
	p.nextToken()

	// read type
	typ := p.ParseType()
	if typ == nil {
		// p.AddError(p.currentToken(), "expected type")
		return nil
	}

	// skip equal sign
	if !p.curTokenIs(token.Assign) {
		p.AddError(p.currentToken(), "expected '='")
		return nil
	}
	p.nextToken()

	// read value
	value := p.ParseExpression()
	if value == nil {
		// p.AddError(p.currentToken(), "expected expression")
		return nil
	}
	return &ast.ConstDeclItem{ConstName: ident, ConstType: typ, ConstValue: value}
}

// ----------------------------------------------------------------------------
// Enumeration declaration
// ----------------------------------------------------------------------------
func (p *Parser) ParseEnumStatement() *ast.EnumDecl {
	stmt := &ast.EnumDecl{EnumKw: p.currentToken()}
	stmt.Values = make([]ast.Identifier, 0)
	stmt.EnumKw = p.currentToken()
	if !p.curTokenIs(token.KwEnum) {
		p.AddError(p.currentToken(), "expected Enum")
		return nil
	}
	p.nextToken()

	// read identifier
	if !p.curTokenIs(token.Identifier) {
		// p.AddError(p.currentToken(), "expected identifier")
		return nil
	}
	ident := &ast.Identifier{Tok: p.currentToken(), Name: p.currentToken().Literal}
	stmt.Identifier = ident
	p.nextToken()
	for {

		// look for end of line
		if !p.curTokenIs(token.EOL) {
			p.AddError(p.currentToken(), "expected EOL")
			return nil
		}
		p.nextToken()

		// look for end of enumeration
		if p.curTokenIs(token.KwEnd) {
			break
		}

		// read identifer
		ident := p.ParseIdentifier()
		if ident == nil {
			// p.AddError(p.currentToken(), "expected identifier")
			return nil
		}
		stmt.Values = append(stmt.Values, *ident)
	}

	// skip end keyword
	if !p.curTokenIs(token.KwEnd) {
		p.AddError(p.currentToken(), "expected 'end'")
		return nil
	}
	p.nextToken()

	// skip enum keyword
	if !p.curTokenIs(token.KwEnum) {
		p.AddError(p.currentToken(), "expected 'enum'")
		return nil
	}
	p.nextToken()

	return stmt
}

// ----------------------------------------------------------------------------
// Function declaration
//
// "Function" simpleIdent "(" Params ")" "As" BasicType 	eol
// "Function" simpleIdent "(" Params ")" "As" simpleIdent	eol
// ----------------------------------------------------------------------------
func (p *Parser) ParseFunctionStatement() *ast.FuncDecl {
	stmt := &ast.FuncDecl{FunctionKw: p.currentToken()}

	if !p.curTokenIs(token.KwFunction) {
		p.AddError(p.currentToken(), "expected Function")
		return nil
	}
	p.nextToken()

	stmt.FuncName = p.ParseIdentifier()
	if stmt.FuncName == nil {
		// p.AddError(p.currentToken(), "expected identifier")
		return nil
	}

	// read function type
	stmt.FuncType = p.ParseFunctionType()
	if stmt.FuncType == nil {
		// p.AddError(p.currentToken(), "expected function type")
		return nil
	}

	// look for end of line
	if !p.curTokenIs(token.EOL) {
		p.AddError(p.currentToken(), "expected EOL")
		return nil
	}
	p.nextToken()
	stmt.Body = p.ParseFunctionBody()
	if stmt.Body == nil {
		// p.AddError(p.currentToken(), "expected function body")
		return nil
	}
	return stmt
}

func (p *Parser) ParseFunctionType() *ast.FuncType {
	funcType := &ast.FuncType{}

	if !p.curTokenIs(token.Lparen) {
		p.AddError(p.currentToken(), "expected left parenthesis")
		return nil
	}
	funcType.Lparen = p.currentToken()
	p.nextToken()
	funcType.Params = make([]ast.ParamItem, 0)
	for {
		if p.curTokenIs(token.Rparen) {
			break
		}
		param := p.ParseParamItem()
		if param == nil {
			// p.AddError(p.currentToken(), "expected parameter")
			return nil
		}
		funcType.Params = append(funcType.Params, *param)
		if p.curTokenIs(token.Comma) {
			p.nextToken()
		}
	}
	if !p.curTokenIs(token.Rparen) {
		p.AddError(p.currentToken(), "expected right parenthesis")
		return nil
	}
	funcType.Rparen = p.currentToken().Position
	p.nextToken()

	// skip as keyword
	if !p.curTokenIs(token.KwAs) {
		p.AddError(p.currentToken(), "expected 'as'")
		return nil
	}
	p.nextToken()

	// read type
	typ := p.ParseType()
	if typ == nil {
		// p.AddError(p.currentToken(), "expected type")
		return nil
	}
	funcType.Result = typ
	return funcType
}

// ParseParameterItem parses a parameter item.
func (p *Parser) ParseParamItem() *ast.ParamItem {
	param := &ast.ParamItem{}

	// read Optional keyword
	if p.curTokenIs(token.KwOptional) {
		param.Optional = true
		p.nextToken()
	}

	// read byVal keyword
	if p.curTokenIs(token.KwByVal) {
		param.ByVal = true
		p.nextToken()
	} else if p.curTokenIs(token.KwByRef) {
		param.ByVal = false
		p.nextToken()
	}
	// read ParamArray keyword
	if p.curTokenIs(token.KwParamArray) {
		param.ParamArray = true
		p.nextToken()
	}

	param.VarName = p.ParseIdentifier()
	if param.VarName == nil {
		return nil
	}

	// check if it is an array
	if p.curTokenIs(token.Lparen) && p.peekTokenIs(1, token.Rparen) {
		param.IsArray = true
		p.nextToken()
		p.nextToken()
	}

	// skip as keyword
	if !p.curTokenIs(token.KwAs) {
		p.AddError(p.currentToken(), "expected 'as'")
		return nil
	}
	p.nextToken()

	param.VarType = p.ParseType()
	if param.VarType == nil {
		p.AddError(p.currentToken(), "expected type")
		return nil
	}

	// look for default value
	if p.curTokenIs(token.Assign) {
		p.nextToken()
		param.DefaultValue = p.ParseExpression()
		if param.DefaultValue == nil {
			// p.AddError(p.currentToken(), "expected expression")
			return nil
		}
	}

	return param
}

// ParseFunctionBody parses a function body.
func (p *Parser) ParseFunctionBody() []ast.StatementList {
	block := make([]ast.StatementList, 0)

	// read function statement lists
	for !(p.curTokenIs(token.KwEnd) && p.peekTokenIs(1, token.KwFunction)) {

		stmtList := p.ParseStatementList(false, true)
		if stmtList == nil {
			return nil
		}
		block = append(block, *stmtList)
		p.nextToken()
	}
	// skip end keyword
	if !p.curTokenIs(token.KwEnd) {
		p.AddError(p.currentToken(), "expected 'end'")
		return nil
	}
	p.nextToken()

	// skip function keyword
	if !p.curTokenIs(token.KwFunction) {
		p.AddError(p.currentToken(), "expected 'function'")
		return nil
	}
	p.nextToken()
	return block
}

// ParseExitStatement parses an exit statement.
func (p *Parser) ParseExitStatement(inSub bool, inFunction bool) *ast.ExitStmt {
	stmt := &ast.ExitStmt{ExitKw: p.currentToken()}
	if !p.curTokenIs(token.KwExit) {
		p.AddError(p.currentToken(), "expected Exit")
		return nil
	}
	p.nextToken()
	// get type of exit
	switch p.currentToken().Kind {
	case token.KwDo, token.KwFor:
		stmt.ExitType = p.currentToken()
	case token.KwFunction:
		if inFunction {
			stmt.ExitType = p.currentToken()
		} else {
			p.AddError(p.currentToken(), "expected 'do', 'for' or 'function'")
			return nil
		}
	case token.KwSub:
		if inSub {
			stmt.ExitType = p.currentToken()
		} else {
			p.AddError(p.currentToken(), "expected 'do', 'for' or 'sub'")
			return nil
		}
	default:
		p.AddError(p.currentToken(), "expected 'do', 'for', 'sub' or 'function'")
		return nil
	}
	p.nextToken()

	return stmt
}

// ----------------------------------------------------------------------------
// Empty statement declaration
// ----------------------------------------------------------------------------

// ParseEmptyStatement parses an empty statement.
func (p *Parser) ParseEmptyStatement() *ast.EmptyStmt {
	// check if it is an EOL
	if !p.curTokenIs(token.EOL) || p.peekTokenIs(1, token.EOF) {
		p.AddError(p.currentToken(), "expected EOL or EOF")
		return nil
	}

	stmt := &ast.EmptyStmt{EOL: p.currentToken()}
	p.nextToken()
	return stmt
}

// ----------------------------------------------------------------------------
// Subroutine declaration
// "Sub" simpleIdent "(" Params ")" 	eol
// ----------------------------------------------------------------------------
func (p *Parser) ParseSubStatement() *ast.SubDecl {
	stmt := &ast.SubDecl{SubKw: p.currentToken()}

	if !p.curTokenIs(token.KwSub) {
		p.AddError(p.currentToken(), "expected subroutine")
		return nil
	}
	p.nextToken()

	stmt.SubName = p.ParseIdentifier()
	if stmt.SubName == nil {
		// p.AddError(p.currentToken(), "expected identifier")
		return nil
	}

	// read subroutine type
	stmt.SubType = p.ParseSubType()
	if stmt.SubType == nil {
		// p.AddError(p.currentToken(), "expected function type")
		return nil
	}

	// look for end of line
	if !p.curTokenIs(token.EOL) {
		p.AddError(p.currentToken(), "expected EOL")
		return nil
	}
	p.nextToken()
	stmt.Body = p.ParseSubBody()
	if stmt.Body == nil {
		// p.AddError(p.currentToken(), "expected function body")
		return nil
	}
	return stmt
}

func (p *Parser) ParseSubType() *ast.SubType {
	SubType := &ast.SubType{}

	if !p.curTokenIs(token.Lparen) {
		p.AddError(p.currentToken(), "expected left parenthesis")
		return nil
	}
	SubType.Lparen = p.currentToken()
	p.nextToken()
	SubType.Params = make([]ast.ParamItem, 0)
	for {
		if p.curTokenIs(token.Rparen) {
			break
		}
		param := p.ParseParamItem()
		if param == nil {
			// p.AddError(p.currentToken(), "expected parameter")
			return nil
		}
		SubType.Params = append(SubType.Params, *param)
		if p.curTokenIs(token.Comma) {
			p.nextToken()
		}
	}
	if !p.curTokenIs(token.Rparen) {
		p.AddError(p.currentToken(), "expected right parenthesis")
		return nil
	}
	SubType.Rparen = p.currentToken().Position
	p.nextToken()

	return SubType
}

// ParseSubBody parses a function body.
func (p *Parser) ParseSubBody() []ast.StatementList {
	block := make([]ast.StatementList, 0)

	// read subroutine statement lists
	for !(p.curTokenIs(token.KwEnd) && p.peekTokenIs(1, token.KwSub)) {

		stmtList := p.ParseStatementList(true, false)
		if stmtList == nil {
			return nil
		}
		block = append(block, *stmtList)
		p.nextToken()

		if p.curTokenIs(token.EOF) {
			p.AddError(p.currentToken(), "expected 'end'")
			return nil
		}
	}
	// skip end keyword
	if !p.curTokenIs(token.KwEnd) {
		p.AddError(p.currentToken(), "expected 'end'")
		return nil
	}
	p.nextToken()

	// skip sub keyword
	if !p.curTokenIs(token.KwSub) {
		p.AddError(p.currentToken(), "expected 'sub'")
		return nil
	}
	p.nextToken()
	return block
}

// ParseInFunctionExitStatement parses an exit statement.
func (p *Parser) ParseInSubExitStatement() *ast.ExitStmt {
	stmt := &ast.ExitStmt{ExitKw: p.currentToken()}
	if !p.curTokenIs(token.KwExit) {
		p.AddError(p.currentToken(), "expected Exit")
		return nil
	}
	p.nextToken()
	// get type of exit
	switch p.currentToken().Kind {
	case token.KwDo, token.KwFor, token.KwSub:
		stmt.ExitType = p.currentToken()
	default:
		p.AddError(p.currentToken(), "expected 'do', 'for' or 'sub'")
		return nil
	}
	p.nextToken()

	return stmt
}

// ----------------------------------------------------------------------------
// if declaration
// ----------------------------------------------------------------------------
func (p *Parser) ParseIfStatement(inSub bool, inFunction bool) *ast.IfStmt {
	stmt := &ast.IfStmt{IfKw: p.currentToken()}
	if !p.curTokenIs(token.KwIf) {
		p.AddError(p.currentToken(), "expected If")
		return nil
	}
	p.nextToken()

	// read condition
	stmt.Condition = p.ParseExpression()
	if stmt.Condition == nil {
		// p.AddError(p.currentToken(), "expected expression")
		return nil
	}
	// look for then keyword
	if !p.curTokenIs(token.KwThen) {
		p.AddError(p.currentToken(), "expected 'then'")
		return nil
	}
	p.nextToken()

	// look for end of line
	if !p.curTokenIs(token.EOL) {
		p.AddError(p.currentToken(), "expected EOL")
		return nil
	}
	p.nextToken()

	stmt.Body = p.ParseIfBody(inSub, inFunction)
	if stmt.Body == nil {
		// p.AddError(p.currentToken(), "expected if body")
		return nil
	}

	// look for ElseIf keyword
	for p.curTokenIs(token.KwElseIf) {
		stmt.ElseIf = append(stmt.ElseIf, *p.ParseElseIfStatement(inSub, inFunction))
	}

	// look for Else keyword
	if p.curTokenIs(token.KwElse) {
		stmt.Else = p.ParseElseStatement(inSub, inFunction)
	}

	// skip end keyword
	if !p.curTokenIs(token.KwEnd) {
		p.AddError(p.currentToken(), "expected 'end'")
		return nil
	}
	p.nextToken()

	// skip if keyword
	if !p.curTokenIs(token.KwIf) {
		p.AddError(p.currentToken(), "expected 'if'")
		return nil
	}
	p.nextToken()

	return stmt
}

// ParseIfBody parses an if body.
func (p *Parser) ParseIfBody(inSub bool, inFunction bool) []ast.StatementList {
	block := make([]ast.StatementList, 0)

	// read if statement lists
	for !(p.curTokenIs(token.KwEnd) && p.peekTokenIs(1, token.KwIf)) && !p.curTokenIs(token.KwElseIf) && !p.curTokenIs(token.KwElse) {

		stmtList := p.ParseStatementList(inSub, inFunction)
		if stmtList == nil {
			return nil
		}
		block = append(block, *stmtList)
		p.nextToken()
		if p.curTokenIs(token.EOF) {
			p.AddError(p.currentToken(), "expected 'end'")
			return nil
		}

	}
	return block
}

// ParseElseIfStatement parses an else if statement.
func (p *Parser) ParseElseIfStatement(inSub bool, inFunction bool) *ast.ElseIfStmt {
	stmt := &ast.ElseIfStmt{ElseIfKw: p.currentToken()}
	if !p.curTokenIs(token.KwElseIf) {
		p.AddError(p.currentToken(), "expected ElseIf")
		return nil
	}
	p.nextToken()

	// read condition
	stmt.Condition = p.ParseExpression()
	if stmt.Condition == nil {
		// p.AddError(p.currentToken(), "expected expression")
		return nil
	}
	// look for then keyword
	if !p.curTokenIs(token.KwThen) {
		p.AddError(p.currentToken(), "expected 'then'")
		return nil
	}
	p.nextToken()

	// look for end of line
	if !p.curTokenIs(token.EOL) {
		p.AddError(p.currentToken(), "expected EOL")
		return nil
	}

	p.nextToken()
	stmt.Body = p.ParseElseIfBody(inSub, inFunction)
	if stmt.Body == nil {
		// p.AddError(p.currentToken(), "expected if body")
		return nil
	}
	return stmt
}

// ParseElseIfBody parses an elseif body.
func (p *Parser) ParseElseIfBody(inSub bool, inFunction bool) []ast.StatementList {
	block := make([]ast.StatementList, 0)

	// read if statement lists
	for !(p.curTokenIs(token.KwEnd) && p.peekTokenIs(1, token.KwIf)) && !p.curTokenIs(token.KwElseIf) && !p.curTokenIs(token.KwElse) {

		stmtList := p.ParseStatementList(inSub, inFunction)
		if stmtList == nil {
			return nil
		}
		block = append(block, *stmtList)
		p.nextToken()
	}
	return block
}

// ParseElseStatement parses an else statement.
func (p *Parser) ParseElseStatement(inSub bool, inFunction bool) []ast.StatementList {
	block := make([]ast.StatementList, 0)

	// expect else keyword
	if !p.curTokenIs(token.KwElse) {
		p.AddError(p.currentToken(), "expected Else")
		return nil
	}
	p.nextToken()

	// look for end of line
	if !p.curTokenIs(token.EOL) {
		p.AddError(p.currentToken(), "expected EOL")
		return nil
	}
	p.nextToken()

	// read else statement lists
	for !(p.curTokenIs(token.KwEnd) && p.peekTokenIs(1, token.KwIf)) && !p.curTokenIs(token.KwElseIf) && !p.curTokenIs(token.KwElse) {

		stmtList := p.ParseStatementList(inSub, inFunction)
		if stmtList == nil {
			return nil
		}
		block = append(block, *stmtList)
		p.nextToken()
		if p.curTokenIs(token.EOF) {
			p.AddError(p.currentToken(), "expected 'end'")
			return nil
		}

	}
	return block
}

// ----------------------------------------------------------------------------
// For declaration
// ----------------------------------------------------------------------------

func (p *Parser) ParseForStatement(inSub bool, inFunction bool) *ast.ForStmt {
	stmt := &ast.ForStmt{ForKw: p.currentToken()}
	if !p.curTokenIs(token.KwFor) {
		p.AddError(p.currentToken(), "expected For")
		return nil
	}
	p.nextToken()

	// if it is a variable iteration
	if p.curTokenIs(token.Identifier) && p.peekTokenIs(1, token.Assign) {
		stmt.ForExpression = p.ParseForVariableStatement(inSub, inFunction)
		if ast.IsNil(stmt.ForExpression) {
			// p.AddError(p.currentToken(), "expected for variable iteration")
			return nil
		}
	} else if p.curTokenIs(token.KwEach) {
		stmt.ForExpression = p.ParseForEachStatement(inSub, inFunction)
		if ast.IsNil(stmt.ForExpression) {
			// p.AddError(p.currentToken(), "expected for each iteration")
			return nil
		}
	} else {
		// p.AddError(p.currentToken(), "expected for variable or for each iteration")
		return nil
	}

	// look for end of line
	if !p.curTokenIs(token.EOL) {
		p.AddError(p.currentToken(), "expected EOL")
		return nil
	}
	p.nextToken()

	stmt.Body = p.ParseForBody(inSub, inFunction)
	if stmt.Body == nil {
		// p.AddError(p.currentToken(), "expected for body")
		return nil
	}

	// skip next keyword
	if !p.curTokenIs(token.KwNext) {
		p.AddError(p.currentToken(), "expected 'next'")
		return nil
	}
	p.nextToken()

	// look optional for iteration variable
	if p.curTokenIs(token.Identifier) {
		stmt.Next = p.ParseIdentifier()
		if stmt.Next == nil {
			// p.AddError(p.currentToken(), "expected identifier")
			return nil
		}
	}

	return stmt
}

func (p *Parser) ParseForVariableStatement(inSub bool, inFunction bool) *ast.ForNextExpr {
	stmt := &ast.ForNextExpr{}
	if !p.curTokenIs(token.Identifier) {
		p.AddError(p.currentToken(), "expected identifier")
		return nil
	}
	stmt.Variable = p.ParseIdentifier()

	// skip equal sign
	if !p.curTokenIs(token.Assign) {
		p.AddError(p.currentToken(), "expected '='")
		return nil
	}
	p.nextToken()

	// read start value
	stmt.From = p.ParseExpression()
	if stmt.From == nil {
		// p.AddError(p.currentToken(), "expected expression")
		return nil
	}

	// skip to keyword
	if !p.curTokenIs(token.KwTo) {
		p.AddError(p.currentToken(), "expected 'to'")
		return nil
	}
	p.nextToken()

	// read end value
	stmt.To = p.ParseExpression()
	if stmt.To == nil {
		// p.AddError(p.currentToken(), "expected expression")
		return nil
	}

	// look for step keyword
	if p.curTokenIs(token.KwStep) {
		p.nextToken()
		stmt.Step = p.ParseExpression()
		if ast.IsNil(stmt.Step) {
			// p.AddError(p.currentToken(), "expected expression")
			return nil
		}
	}

	return stmt
}

func (p *Parser) ParseForEachStatement(inSub bool, inFunction bool) *ast.ForEachExpr {
	stmt := &ast.ForEachExpr{}
	if !p.curTokenIs(token.KwEach) {
		p.AddError(p.currentToken(), "expected Each")
		return nil
	}
	p.nextToken()

	// read iteration variable
	stmt.Variable = p.ParseIdentifier()
	if stmt.Variable == nil {
		// p.AddError(p.currentToken(), "expected identifier")
		return nil
	}

	// skip in keyword
	if !p.curTokenIs(token.KwIn) {
		p.AddError(p.currentToken(), "expected 'in'")
		return nil
	}
	p.nextToken()

	// read collection
	stmt.Collection = p.ParseExpression()
	if stmt.Collection == nil {
		// p.AddError(p.currentToken(), "expected expression")
		return nil
	}

	return stmt
}

func (p *Parser) ParseForBody(inSub bool, inFunction bool) []ast.StatementList {
	block := make([]ast.StatementList, 0)

	// read for statement lists
	for !p.curTokenIs(token.KwNext) {

		stmtList := p.ParseStatementList(inSub, inFunction)
		if stmtList == nil {
			return nil
		}
		block = append(block, *stmtList)
		p.nextToken()
		if p.curTokenIs(token.EOF) {
			p.AddError(p.currentToken(), "expected 'next'")
			return nil
		}

	}
	return block
}

// ----------------------------------------------------------------------------
// Select / Case declaration
// ----------------------------------------------------------------------------
func (p *Parser) ParseSelectStatement(inSub bool, inFunction bool) *ast.SelectStmt {
	stmt := &ast.SelectStmt{SelectKw: p.currentToken()}
	if !p.curTokenIs(token.KwSelect) {
		p.AddError(p.currentToken(), "expected Select")
		return nil
	}
	p.nextToken()

	// skip case keyword
	if !p.curTokenIs(token.KwCase) {
		p.AddError(p.currentToken(), "expected 'case'")
		return nil
	}
	p.nextToken()

	// read expression
	stmt.Condition = p.ParseExpression()
	if stmt.Condition == nil {
		// p.AddError(p.currentToken(), "expected expression")
		return nil
	}

	// look for end of line
	if !p.curTokenIs(token.EOL) {
		p.AddError(p.currentToken(), "expected EOL")
		return nil
	}
	p.nextToken()

	stmt.Body = p.ParseSelectBody(inSub, inFunction)
	if stmt.Body == nil {
		// p.AddError(p.currentToken(), "expected select body")
		return nil
	}

	// skip end keyword
	if !p.curTokenIs(token.KwEnd) {
		p.AddError(p.currentToken(), "expected 'end'")
		return nil
	}
	p.nextToken()

	// skip select keyword
	if !p.curTokenIs(token.KwSelect) {
		p.AddError(p.currentToken(), "expected 'select'")
		return nil
	}
	p.nextToken()
	return stmt
}

//	caseExpr struct {
//		// Position of `case` keyword.
//		Case *token.Token
//		// Condition.
//		Condition Expression
//		// Case body.
//		Body []StatementList
//		// parent node
//		Parent Node
//	}
func (p *Parser) ParseSelectBody(inSub bool, inFunction bool) []ast.CaseStmt {
	block := make([]ast.CaseStmt, 0)
	gotElse := false
	for {
		stmt := &ast.CaseStmt{CaseKw: p.currentToken()}
		// skip case keyword
		if !p.curTokenIs(token.KwCase) {
			p.AddError(p.currentToken(), "expected 'case'")
			return nil
		}
		p.nextToken()

		// look for else keyword
		if p.curTokenIs(token.KwElse) {
			gotElse = true
			p.nextToken()
		} else {
			// read expression
			stmt.Condition = p.ParseExpression()
			if stmt.Condition == nil {
				// p.AddError(p.currentToken(), "expected expression")
				return nil
			}
		}

		// read end of line
		if !p.curTokenIs(token.EOL) {
			p.AddError(p.currentToken(), "expected EOL")
			return nil
		}
		p.nextToken()

		// read case statement lists
		for !(p.curTokenIs(token.KwEnd) && p.peekTokenIs(1, token.KwSelect)) && !p.curTokenIs(token.KwCase) {
			stmtList := p.ParseStatementList(inSub, inFunction)
			if stmtList == nil {
				return nil
			}
			stmt.Body = append(stmt.Body, *stmtList)
			p.nextToken()
		}
		block = append(block, *stmt)

		// look for end of select
		if p.curTokenIs(token.KwEnd) && p.peekTokenIs(1, token.KwSelect) {
			break
		}
		if gotElse {
			p.AddError(p.currentToken(), "expected 'end select'")
			return nil
		}
	}
	return block
}

// ----------------------------------------------------------------------------
// Do / Loop declaration
// ----------------------------------------------------------------------------
func (p *Parser) ParseDoStatement(inSub bool, inFunction bool) ast.Statement {
	var stmt ast.Statement

	// check if it is a do while or do until
	if p.peekTokenIs(1, token.KwWhile) {
		stmt = p.ParseDoWhileStatement(inSub, inFunction)
	} else if p.peekTokenIs(1, token.KwUntil) {
		stmt = p.ParseDoUntilStatement(inSub, inFunction)
	} else {
		// at this point we don't know if it is a do loop while or do loop until
		stmt = p.ParseDoLoopStatement(inSub, inFunction)
	}
	return stmt
}

func (p *Parser) ParseDoWhileStatement(inSub bool, inFunction bool) *ast.WhileStmt {
	stmt := &ast.WhileStmt{DoKw: p.currentToken()}
	if !p.curTokenIs(token.KwDo) {
		p.AddError(p.currentToken(), "expected Do")
		return nil
	}
	p.nextToken()

	// read while keyword
	if !p.curTokenIs(token.KwWhile) {
		p.AddError(p.currentToken(), "expected 'while'")
		return nil
	}
	p.nextToken()

	// read condition
	stmt.Condition = p.ParseExpression()
	if stmt.Condition == nil {
		// p.AddError(p.currentToken(), "expected expression")
		return nil
	}

	// look for end of line
	if !p.curTokenIs(token.EOL) {
		p.AddError(p.currentToken(), "expected EOL")
		return nil
	}
	p.nextToken()

	stmt.Body = p.ParseDoBody(inSub, inFunction)
	if stmt.Body == nil {
		// p.AddError(p.currentToken(), "expected do body")
		return nil
	}

	// skip loop keyword
	if !p.curTokenIs(token.KwLoop) {
		p.AddError(p.currentToken(), "expected 'loop'")
		return nil
	}
	p.nextToken()
	return stmt
}

func (p *Parser) ParseDoUntilStatement(inSub bool, inFunction bool) *ast.UntilStmt {
	stmt := &ast.UntilStmt{DoKw: p.currentToken()}
	if !p.curTokenIs(token.KwDo) {
		p.AddError(p.currentToken(), "expected Do")
		return nil
	}
	p.nextToken()

	// read until keyword
	if !p.curTokenIs(token.KwUntil) {
		p.AddError(p.currentToken(), "expected 'until'")
		return nil
	}
	p.nextToken()

	// read condition
	stmt.Condition = p.ParseExpression()
	if stmt.Condition == nil {
		// p.AddError(p.currentToken(), "expected expression")
		return nil
	}

	// look for end of line
	if !p.curTokenIs(token.EOL) {
		p.AddError(p.currentToken(), "expected EOL")
		return nil
	}
	p.nextToken()

	stmt.Body = p.ParseDoBody(inSub, inFunction)
	if stmt.Body == nil {
		// p.AddError(p.currentToken(), "expected do body")
		return nil
	}

	// skip loop keyword
	if !p.curTokenIs(token.KwLoop) {
		p.AddError(p.currentToken(), "expected 'loop'")
		return nil
	}
	p.nextToken()
	return stmt
}

func (p *Parser) ParseDoLoopStatement(inSub bool, inFunction bool) ast.Statement {
	// we will find the time of node at the end of the loop
	var stmt ast.Statement
	tok := p.currentToken()
	// read do keyword
	if !p.curTokenIs(token.KwDo) {
		p.AddError(p.currentToken(), "expected Do")
		return nil
	}
	p.nextToken()

	// read end of line
	if !p.curTokenIs(token.EOL) {
		p.AddError(p.currentToken(), "expected EOL")
		return nil
	}
	p.nextToken()

	// read do statement lists
	body := p.ParseDoBody(inSub, inFunction)
	if body == nil {
		// p.AddError(p.currentToken(), "expected do body")
		return nil
	}

	// skip loop keyword
	if !p.curTokenIs(token.KwLoop) {
		p.AddError(p.currentToken(), "expected 'loop'")
		return nil
	}
	p.nextToken()

	// determine type of loop
	if p.curTokenIs(token.KwWhile) {
		stmt = &ast.DoWhileStmt{DoKw: tok, Body: body}
		p.nextToken()
		stmtDo := stmt.(*ast.DoWhileStmt)
		// read condition
		stmtDo.Condition = p.ParseExpression()
		if stmtDo.Condition == nil {
			// p.AddError(p.currentToken(), "expected expression")
			return nil
		}
	} else if p.curTokenIs(token.KwUntil) {
		p.nextToken()
		stmt = &ast.DoUntilStmt{DoKw: tok, Body: body}
		stmtDo := stmt.(*ast.DoUntilStmt)
		// read condition
		stmtDo.Condition = p.ParseExpression()
		if stmtDo.Condition == nil {
			// p.AddError(p.currentToken(), "expected expression")
			return nil
		}
	} else {
		p.AddError(p.currentToken(), "expected 'while' or 'until'")
		return nil
	}
	return stmt
}

func (p *Parser) ParseDoBody(inSub bool, inFunction bool) []ast.StatementList {
	block := make([]ast.StatementList, 0)
	for !(p.curTokenIs(token.KwLoop)) {
		stmtList := p.ParseStatementList(inSub, inFunction)
		if stmtList == nil {
			return nil
		}
		block = append(block, *stmtList)
		p.nextToken()
		if p.curTokenIs(token.EOF) {
			p.AddError(p.currentToken(), "expected 'Loop'")
			return nil
		}

	}
	return block
}

// ----------------------------------------------------------------------------
// Subroutine call
// ----------------------------------------------------------------------------

// use Call keyword to simplify parsing
func (p *Parser) ParseSubroutineCall() *ast.CallSubStmt {
	stmt := &ast.CallSubStmt{}

	// expect call keyword
	if !p.curTokenIs(token.KwCall) {
		p.AddError(p.currentToken(), "expected Call")
		return nil
	}
	stmt.CallKw = p.currentToken()
	p.nextToken()

	// read identifier
	stmt.Definition = p.ParsePrimaryExpr()
	if stmt.Definition == nil {
		// p.AddError(p.currentToken(), "expected identifier")
		return nil
	}

	return stmt
}

// ----------------------------------------------------------------------------
// Error handling
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Special statement
// ----------------------------------------------------------------------------

func (p *Parser) ParseSpecialStatement() ast.Statement {
	stmt := &ast.SpecialStmt{Keyword1: p.currentToken()}
	stmt.Args = make([]ast.Expression, 0)
	switch p.currentToken().Kind {
	case token.KwStop:
		p.nextToken()
		return stmt
	case token.KwResume:
		p.nextToken()
		if p.curTokenIs(token.KwNext) {
			stmt.Keyword2 = p.peekToken(1).Literal
			p.nextToken()
		} else {
			// if identifier -> resume label
			stmt.Args = append(stmt.Args, p.ParseIdentifier())
		}
		return stmt
	case token.KwGoto, token.KwErase:
		p.nextToken()
		stmt.Args = append(stmt.Args, p.ParseIdentifier())
		return stmt
	case token.KwRedim:
		p.nextToken()
		if p.currentToken().Kind == token.KwPreserve {
			stmt.Keyword2 = p.currentToken().Literal
			p.nextToken()
		}
		identifier := p.ParseIdentifier()
		stmt.Args = append(stmt.Args, p.ParseArrayOrSubroutineCall(identifier))

		return stmt
	case token.Identifier:
		if p.peekTokenIs(1, token.Colon) && p.peekTokenIs(2, token.EOL) {
			return p.ParseLabelStatement()
		} else if strings.ToLower(p.currentToken().Literal) == "debug" &&
			p.peekToken(1).Kind == token.Dot &&
			strings.ToLower(p.peekToken(2).Literal) == "print" {
			return p.ParsePrintStatement()
		} else if strings.ToLower(p.currentToken().Literal) == "print" {
			return p.ParsePrintStatement()
		} else if strings.ToLower(p.currentToken().Literal) == "msgbox" {
			return p.ParseMsgBoxStatement()
		} else if strings.ToLower(p.currentToken().Literal) == "input" {
			return p.ParseInputStatement()
		}
	}
	return nil // no error, but no statement
}

func (p *Parser) ParsePrintStatement() *ast.SpecialStmt {
	stmt := &ast.SpecialStmt{Keyword1: p.currentToken()}
	currentToken := strings.ToLower(p.currentToken().Literal)
	ok := false
	if currentToken == "print" {
		ok = true
		stmt.Keyword1 = p.currentToken()
		p.nextToken()
	} else if currentToken == "debug" && p.peekToken(1).Kind == token.Dot &&
		strings.ToLower(p.peekToken(2).Literal) == "print" {
		ok = true
		stmt.Keyword1 = &token.Token{Literal: "Debug.Print", Kind: token.Identifier, Position: p.currentToken().Position.Copy()}
		p.nextToken()
		p.nextToken()
		p.nextToken()
	}
	if !ok {
		p.AddError(p.currentToken(), "expected Print or Debug.Print")
		return nil
	}

	// read expression if not empty
	if !p.curTokenIs(token.Semicolon) && !p.curTokenIs(token.EOL) && !p.curTokenIs(token.EOF) {
		for {
			expr := p.ParseExpression()
			if expr == nil {
				return nil
			}
			stmt.Args = append(stmt.Args, expr)

			// look for a comma to see if there are more expressions
			if !p.curTokenIs(token.Comma) {
				break
			}
			p.nextToken() // skip the comma
		}
	}
	// if last character is a semicolon, add it to the statement
	if p.curTokenIs(token.Semicolon) {
		stmt.Semicolon = p.currentToken()
		p.nextToken()
	}
	return stmt
}

func (p *Parser) ParseMsgBoxStatement() *ast.SpecialStmt {
	stmt := &ast.SpecialStmt{Keyword1: p.currentToken()}
	stmt.Keyword1.Literal = "MsgBox"
	currentToken := strings.ToLower(p.currentToken().Literal)
	if currentToken == "msgbox" {
		stmt.Keyword1 = p.currentToken()
		p.nextToken()
	} else {
		p.AddError(p.currentToken(), "expected MsgBox")
		return nil
	}

	// read expression
	for {
		expr := p.ParseExpression()
		if expr == nil {
			return nil
		}
		stmt.Args = append(stmt.Args, expr)

		// look for a comma to see if there are more expressions
		if !p.curTokenIs(token.Comma) {
			break
		}
		p.nextToken() // skip the comma
	}
	return stmt
}

func (p *Parser) ParseInputStatement() *ast.SpecialStmt {
	stmt := &ast.SpecialStmt{Keyword1: p.currentToken()}
	stmt.Keyword1.Literal = "Input"
	currentToken := strings.ToLower(p.currentToken().Literal)
	if currentToken == "input" {
		stmt.Keyword1 = p.currentToken()
		p.nextToken()
	} else {
		p.AddError(p.currentToken(), "expected Input")
		return nil
	}

	// read expression
	count := 0
	for {
		count++
		if count > 2 {
			p.AddError(p.currentToken(), "expected 1 or 2 arguments")
			return nil
		}
		expr := p.ParseExpression()
		if expr == nil {
			return nil
		}
		stmt.Args = append(stmt.Args, expr)

		// look for a comma to see if there are more expressions
		if !p.curTokenIs(token.Comma) {
			break
		}
		p.nextToken() // skip the comma
	}
	return stmt
}

func (p *Parser) ParseLabelStatement() *ast.LabelDecl {
	stmt := &ast.LabelDecl{}
	stmt.LabelName = p.ParseIdentifier()
	if stmt.LabelName == nil {
		p.AddError(p.currentToken(), "expected identifier")
		return nil
	}
	p.nextToken()
	if !p.curTokenIs(token.Colon) {
		p.AddError(p.currentToken(), "expected colon")
		return nil
	}
	p.nextToken()
	return stmt
}

// ----------------------------------------------------------------------------
// Expressions
// ----------------------------------------------------------------------------

func (p *Parser) ParseExpressionStatement() *ast.ExprStmt {
	stmt := &ast.ExprStmt{}

	// skip let keyword
	if !p.curTokenIs(token.KwLet) {
		p.AddError(p.currentToken(), "expected Let")
		return nil
	}
	p.nextToken()

	stmt.Expression = p.ParseExpression()
	if stmt.Expression == nil {
		// p.AddError(p.currentToken(), "expected expression")
		return nil
	}
	return stmt
}

// Expr
//
//	: Expr2R
func (p *Parser) ParseExpression() ast.Expression {
	return p.ParseExpression2R()
}

// Expr2R
//
//	: Expr3L
//	| PrimaryExpr "=" Expr3L       											<< astx.NewBinaryExpr($0, $1, $2) >>
//	;
func (p *Parser) ParseExpression2R() ast.Expression {
	expr := p.ParseExpression3L()
	if expr == nil {
		// p.AddError(p.currentToken(), "expected expression")
		return nil
	}

	// check if it is an assignment
	if p.curTokenIs(token.Assign) {
		assignment := ast.BinaryExpr{}
		assignment.Left = expr
		assignment.OpKind = token.Assign
		assignment.OpToken = p.currentToken()
		p.nextToken()
		assignment.Right = p.ParseExpression3L()
		if assignment.Right == nil {
			// p.AddError(p.currentToken(), "expected expression")
			return nil
		}
		return &assignment
	}
	return expr
}

// Expr3L
//
//	: Expr5L
//	| Expr3L "Or" Expr5L
func (p *Parser) ParseExpression3L() ast.Expression {
	expr := p.ParseExpression5L()
	if expr != nil {
		for p.curTokenIs(token.Or) {
			tok := p.currentToken()
			p.nextToken()
			right := p.ParseExpression5L()
			expr = &ast.BinaryExpr{Left: expr, OpKind: token.Or, OpToken: tok, Right: right}
			if right == nil {
				// p.AddError(p.currentToken(), "expected expression")
				return nil
			}
		}
	}
	return expr
}

// Expr5L
//
//	: Expr9L
//	| Expr5L "And" Expr9L
func (p *Parser) ParseExpression5L() ast.Expression {
	expr := p.ParseExpression9L()
	if expr != nil {
		for p.curTokenIs(token.And) {
			tok := p.currentToken()
			p.nextToken()
			right := p.ParseExpression9L()
			expr = &ast.BinaryExpr{Left: expr, OpKind: token.And, OpToken: tok, Right: right}
			if right == nil {
				// p.AddError(p.currentToken(), "expected expression")
				return nil
			}
		}
	}
	return expr
}

// Expr9L
//
//	: Expr10L
//	| Expr9L "=" Expr10L
//	| Expr9L "<>" Expr10L
func (p *Parser) ParseExpression9L() ast.Expression {
	expr := p.ParseExpression10L()
	if expr != nil {
		for p.curTokenIs(token.Eq) || p.curTokenIs(token.Neq) {
			tok := p.currentToken()
			p.nextToken()
			right := p.ParseExpression10L()
			expr = &ast.BinaryExpr{Left: expr, OpKind: tok.Kind, OpToken: tok, Right: right}
			if right == nil {
				// p.AddError(p.currentToken(), "expected expression")
				return nil
			}
		}
	}
	return expr
}

// Expr10L
//
//	: Expr12L
//	| Expr10L "<" Expr12L
//	| Expr10L ">" Expr12L
//	| Expr10L "<=" Expr12L
//	| Expr10L ">=" Expr12L
func (p *Parser) ParseExpression10L() ast.Expression {
	expr := p.ParseExpression12L()
	if expr != nil {
		for p.curTokenIs(token.Lt) || p.curTokenIs(token.Gt) || p.curTokenIs(token.Le) || p.curTokenIs(token.Ge) {
			tok := p.currentToken()
			p.nextToken()
			right := p.ParseExpression12L()
			expr = &ast.BinaryExpr{Left: expr, OpKind: tok.Kind, OpToken: tok, Right: right}
			if right == nil {
				// p.AddError(p.currentToken(), "expected expression")
				return nil
			}
		}
	}
	return expr
}

// Expr12L
//
//	: Expr13L
//	| Expr12L "+" Expr13L
//	| Expr12L "&" Expr13L
//	| Expr12L "-" Expr13L
func (p *Parser) ParseExpression12L() ast.Expression {
	expr := p.ParseExpression13L()
	if expr != nil {
		for p.curTokenIs(token.Add) || p.curTokenIs(token.Concat) || p.curTokenIs(token.Minus) {
			tok := p.currentToken()
			p.nextToken()
			right := p.ParseExpression13L()
			expr = &ast.BinaryExpr{Left: expr, OpKind: tok.Kind, OpToken: tok, Right: right}
			if right == nil {
				// p.AddError(p.currentToken(), "expected expression")
				return nil
			}
		}
	}
	return expr
}

// Expr13L
//
//	: Expr14L
//	| Expr13L "*" Expr14L
//	| Expr13L "Mod" Expr14L
//	| Expr13L "Div" Expr14L
//	| Expr13L "/" Expr14L
func (p *Parser) ParseExpression13L() ast.Expression {
	expr := p.ParseExpression14L()
	if expr != nil {
		for p.curTokenIs(token.Mul) || p.curTokenIs(token.Mod) || p.curTokenIs(token.Div) || p.curTokenIs(token.IntDiv) {
			tok := p.currentToken()
			p.nextToken()
			right := p.ParseExpression14L()
			expr = &ast.BinaryExpr{Left: expr, OpKind: tok.Kind, OpToken: tok, Right: right}
			if right == nil {
				// p.AddError(p.currentToken(), "expected expression")
				return nil
			}
		}
	}
	return expr
}

// Expr14L
//
//	: Expr15
//	| Expr14L "Exp" Expr15
func (p *Parser) ParseExpression14L() ast.Expression {
	expr := p.ParseExpression15()
	if expr != nil {
		for p.curTokenIs(token.Exponent) {
			tok := p.currentToken()
			p.nextToken()
			right := p.ParseExpression15()
			expr = &ast.BinaryExpr{Left: expr, OpKind: tok.Kind, OpToken: tok, Right: right}
			if right == nil {
				// p.AddError(p.currentToken(), "expected expression")
				return nil
			}
		}
	}
	return expr
}

// Expr15
//
//	: Expr16
//	| "-" Expr16
//	| "Not" Expr16
func (p *Parser) ParseExpression15() ast.Expression {
	if p.curTokenIs(token.Minus) || p.curTokenIs(token.Not) {
		expr := &ast.UnaryExpr{OpKind: p.currentToken().Kind, OpToken: p.currentToken()}
		p.nextToken()
		right := p.ParseExpression16()
		if right == nil {
			// p.AddError(p.currentToken(), "expected expression")
			return nil
		}
		expr.Right = right
		return expr
	}
	return p.ParseExpression16()
}

// Expr16
//
//	: PrimaryExpr
//	| long_lit
//	| double_lit
//	| string_lit
//	| dateTime_lit
//	| "True"
//	| "False"
//	| "Nothing"
//	| ParenExpr
func (p *Parser) ParseExpression16() ast.Expression {
	switch p.currentToken().Kind {
	case token.DoubleLit, token.StringLit, token.BooleanLit, token.DateLit, token.LongLit, token.KwTrue, token.KwFalse, token.KwNothing, token.CurrencyLit:
		expr := &ast.BasicLit{Kind: p.currentToken().Kind, ValPos: p.currentToken(), Value: p.currentToken().Literal}
		p.nextToken()
		return expr
	case token.Identifier:
		return p.ParsePrimaryExpr()
	case token.Lparen:
		return p.ParseParenExpr()
	default:
		p.AddError(p.currentToken(), "expected expression")
		return nil
	}
}

// ParenExpr
func (p *Parser) ParseParenExpr() *ast.ParenExpr {
	expr := &ast.ParenExpr{Lparen: p.currentToken()}
	if !p.curTokenIs(token.Lparen) {
		p.AddError(p.currentToken(), "expected left parenthesis")
		return nil
	}
	p.nextToken()
	expr.Expr = p.ParseExpression()
	if expr.Expr == nil {
		// p.AddError(p.currentToken(), "expected expression")
		return nil
	}
	if !p.curTokenIs(token.Rparen) {
		p.AddError(p.currentToken(), "expected right parenthesis")
		return nil
	}
	expr.Rparen = p.currentToken().Position
	p.nextToken()
	return expr
}

func (p *Parser) ParseArrayOrSubroutineCall(identifier *ast.Identifier) *ast.CallOrIndexExpr {
	expr := ast.CallOrIndexExpr{}

	// read identifier
	expr.Identifier = identifier
	if expr.Identifier == nil {
		// p.AddError(p.currentToken(), "expected identifier")
		return nil
	}
	// look for parenthesis
	if p.curTokenIs(token.Lparen) {
		expr.Lparen = p.currentToken()
		p.nextToken()
	} else {
		return nil
	}
	// read arguments
	for !p.curTokenIs(token.Rparen) {

		arg := p.ParseExpression()
		if arg == nil {
			// p.AddError(p.currentToken(), "expected expression")
			return nil
		}
		expr.Args = append(expr.Args, arg)
		// look for a comma to see if there are more arguments
		if p.curTokenIs(token.Comma) {
			p.nextToken() // skip the comma
		}
		if p.curTokenIs(token.EOF) {
			p.AddError(p.currentToken(), "expected right parenthesis")
			return nil
		}
	}
	// get right parenthesis
	if !p.curTokenIs(token.Rparen) {
		p.AddError(p.currentToken(), "expected right parenthesis")
		return nil
	}
	expr.Rparen = p.currentToken().Position
	p.nextToken()
	return &expr
}

// PrimaryExpr
//
//		: simpleIdent
//	 | ArrayOrCallExpr
//		| PrimaryExpr "." simpleIdent
//		| PrimaryExpr "." ArrayOrCallExpr
func (p *Parser) ParsePrimaryExpr() ast.Expression {
	var expr ast.Expression
	expr1 := p.ParseIdentifier()
	if expr1 == nil {
		return nil
	}

	// check if it is an array or a function call
	expr2 := p.ParseArrayOrSubroutineCall(expr1)
	if expr2 != nil {
		expr = expr2
	} else {
		expr = expr1
	}

	// check if we have a dot
	if p.curTokenIs(token.Dot) {
		memberAccess := ast.CallSelectorExpr{}
		memberAccess.Root = expr
		memberAccess.Dot = p.currentToken().Position
		p.nextToken()
		memberAccess.Selector = p.ParsePrimaryExpr()
		return &memberAccess
	}
	return expr
}
