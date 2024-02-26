package ast

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"strings"
	"uBasic/token"
)

// ReadCompact(d Dict) error

const (
	errorValue = iota
	BasicLitValue
	BinaryExprValue
	ForNextExprValue
	ForEachExprValue
	ExitStmtValue
	SpecialStmtValue
	CallOrIndexExprValue
	CallSelectorExprValue
	EmptyStmtValue
	ExprStmtValue
	FileValue
	StatementListValue
	SelectStmtValue
	CaseExprValue
	CallSubStmtValue
	SubReturnTypeValue
	UserDefinedTypeValue
	FuncDeclValue
	SubDeclValue
	FuncTypeValue
	SubTypeValue
	IdentifierValue
	IfStmtValue
	ElseIfStmtValue
	ParenExprValue
	TypeDefValue
	UnaryExprValue
	EnumDeclValue
	ConstDeclValue
	ArrayDeclValue
	ArrayTypeValue
	ScalarDeclValue
	ConstDeclItemValue
	DimDeclValue
	ParamItemValue
	WhileStmtValue
	UntilStmtValue
	DoWhileStmtValue
	DoUntilStmtValue
	ForStmtValue
	LabelDeclValue

	TokenValue
)

// ReadCompact reads a string and returns a File
func (n *File) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var stmtList Node

	// skip opening parenthesis
	if !s.PopVerify('(') {
		return fmt.Errorf("error reading File")
	}

	// get statement lists
	for {
		if s.Peek() == ')' {
			break
		}
		stmtList, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		n.StatementLists = append(n.StatementLists, *stmtList.(*StatementList))
	}
	// remove last brackets from stack
	if !s.PopVerifyString(")]") {
		return fmt.Errorf("error reading File")
	}
	return nil
}

func (n *File) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:(", FileValue))
	for _, stmtList := range n.StatementLists {
		buf.WriteString(stmtList.WriteCompact(d))
	}
	buf.WriteString(")]")
	return WriteCompactDict(d) + buf.String()
}

// ReadCompact reads a string and returns a statement list
func (n *StatementList) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var stmt Node

	// skip opening parenthesis
	if !s.PopVerify('(') {
		return fmt.Errorf("error reading File")
	}

	// get statements
	for {
		if s.Peek() == ')' {
			break
		}
		stmt, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		n.Statements = append(n.Statements, stmt.(Statement))
	}
	// remove last brackets from stack
	if !s.PopVerifyString(")]") {
		return fmt.Errorf("error reading File")
	}
	return nil
}

func (n *StatementList) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:(", StatementListValue))
	for _, stmt := range n.Statements {
		buf.WriteString(stmt.WriteCompact(d))
	}
	buf.WriteString(")]")
	return buf.String()
}

//		*EmptyStmt
//		*ExprStmt
//		*IfStmt
//		*ElseIFStmt
//		*SelectStmt
//		*CaseExpr
//		*ExitStmt
//		*SpecialStmt
//	 	*DoWhileStmt
//		*DoUntilStmt
//		*WhileStmt
//	    *UntilStmt
//		*ForStmt
//	    *DimDecl
//		*ConstDecl
//		*EnumDecl
//	    *LabelDecl
//	    *CallSubStmt

// ReadCompact reads a string and returns a EmptyStmt
func (n *EmptyStmt) ReadCompact(s *CharStack, d *IdList) error {
	// read token
	var err error
	n.EOL, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}

	// remove last bracket from stack
	if !s.PopVerify(']') {
		return fmt.Errorf("error reading EmptyStmt")
	}
	return nil
}

func (n *EmptyStmt) WriteCompact(d *Dict) string {
	return fmt.Sprintf("[%X:%s]", EmptyStmtValue, WriteCompactToken(n.EOL, d))
}

// ReadCompact reads a string and returns a ExprStmt
func (n *ExprStmt) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	var ok bool
	// get X expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.Expression, ok = node.(Expression)
	if !ok {
		return fmt.Errorf("error reading ExprStmt.X")
	}
	// remove last bracket from stack
	if !s.PopVerify(']') {
		return fmt.Errorf("error reading ExprStmt")
	}
	return nil
}

func (n *ExprStmt) WriteCompact(d *Dict) string {
	return fmt.Sprintf("[%X:%v]", ExprStmtValue, n.Expression.WriteCompact(d))
}

// ReadCompact reads a string and returns an IfStmt
func (n *IfStmt) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	var ok bool
	// read token
	n.IfKw, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}

	// get Cond expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.Condition, ok = node.(Expression)
	if !ok {
		return fmt.Errorf("error reading IfStmt.Cond")
	}
	// skip opening parenthesis
	if !s.PopVerify('(') {
		return fmt.Errorf("error reading IfStmt")
	}

	// get Body expressions
	for {
		if s.Peek() == ')' {
			break
		}
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		stmtList := node.(*StatementList)
		n.Body = append(n.Body, *stmtList)
	}
	// remove last parenthesis from stack
	if !s.PopVerify(')') {
		return fmt.Errorf("error reading IfStmt")
	}
	// skip opening parenthesis
	if !s.PopVerify('(') {
		return fmt.Errorf("error reading IfStmt")
	}
	for {
		if s.Peek() == ')' {
			break
		}
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		stmtList := node.(*ElseIfStmt)
		n.ElseIf = append(n.ElseIf, *stmtList)
	}
	// remove last bracket from stack
	if !s.PopVerify(')') {
		return fmt.Errorf("error reading IfStmt")
	}

	// get Else expressions
	if s.Peek() != ']' {
		// skip opening parenthesis
		if !s.PopVerify('(') {
			return fmt.Errorf("error reading IfStmt")
		}

		// get Else expressions
		for {
			if s.Peek() == ')' {
				break
			}
			node, err = ReadCompact(s, d)
			if err != nil {
				return err
			}
			stmtList := node.(*StatementList)
			n.Else = append(n.Else, *stmtList)
		}
		// remove last bracket from stack
		if !s.PopVerify(')') {
			return fmt.Errorf("error reading IfStmt")
		}
	}
	// remove last bracket from stack
	if !s.PopVerify(']') {
		return fmt.Errorf("error reading IfStmt")
	}

	return nil
}

func (n *IfStmt) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:%s", IfStmtValue, WriteCompactToken(n.Token(), d)))
	buf.WriteString(n.Condition.WriteCompact(d))
	buf.WriteString("(")
	for _, stmtList := range n.Body {
		buf.WriteString(stmtList.WriteCompact(d))
	}
	buf.WriteString(")(")
	if n.ElseIf != nil {
		for _, stmt := range n.ElseIf {
			buf.WriteString(stmt.WriteCompact(d))
		}
	}
	buf.WriteString(")(")
	if n.Else != nil {
		for _, stmtList := range n.Else {
			buf.WriteString(stmtList.WriteCompact(d))
		}
	}

	buf.WriteString(")]")
	return buf.String()
}

// ReadCompact reads a string and returns a ElseIfStmt
func (n *ElseIfStmt) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	var ok bool
	// read token
	n.ElseIfKw, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}

	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.Condition, ok = node.(Expression)
	if !ok {
		return fmt.Errorf("error reading ElseIfStmt.Condition")
	}
	// skip two parenthesis
	if !s.PopVerifyString("(") {
		return fmt.Errorf("error reading ElseIfStmt")
	}

	// get Body expressions
	for {
		if s.Peek() == ')' {
			break
		}
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		stmtList := node.(*StatementList)
		n.Body = append(n.Body, *stmtList)
	}
	// remove last bracket from stack
	if !s.PopVerifyString(")]") {
		return fmt.Errorf("error reading ElseIfStmt")
	}
	return nil
}

func (n *ElseIfStmt) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:%s", ElseIfStmtValue, WriteCompactToken(n.Token(), d)))
	buf.WriteString(n.Condition.WriteCompact(d))
	buf.WriteString("(")
	for _, stmtList := range n.Body {
		buf.WriteString(stmtList.WriteCompact(d))
	}
	buf.WriteString(")]")
	return buf.String()
}

// ReadCompact reads a string and returns a SelectStmt
func (n *SelectStmt) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	var ok bool
	// read token
	n.SelectKw, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}

	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.Condition, ok = node.(Expression)
	if !ok {
		return fmt.Errorf("error reading Select.Condition")
	}
	// skip two parenthesis
	if !s.PopVerifyString("(") {
		return fmt.Errorf("error reading SelectStmt")
	}

	// get Body expressions
	for {
		if s.Peek() == ')' {
			break
		}
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		stmtList := node.(*CaseStmt)
		n.Body = append(n.Body, *stmtList)
	}
	// remove last bracket from stack
	if !s.PopVerifyString(")]") {
		return fmt.Errorf("error reading SelectStmt")
	}
	return nil
}

func (n *SelectStmt) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:%s", SelectStmtValue, WriteCompactToken(n.Token(), d)))
	buf.WriteString(n.Condition.WriteCompact(d))
	buf.WriteString("(")
	for _, stmtList := range n.Body {
		buf.WriteString(stmtList.WriteCompact(d))
	}
	buf.WriteString(")]")
	return buf.String()
}

// ReadCompact reads a string and returns a ElseIfStmt
func (n *CaseStmt) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	var ok bool
	// read token
	n.CaseKw, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}

	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.Condition, ok = node.(Expression)
	if !ok {
		return fmt.Errorf("error reading ElseIfStmt.Condition")
	}
	// skip two parenthesis
	if !s.PopVerifyString("(") {
		return fmt.Errorf("error reading ElseIfStmt")
	}

	// get Body expressions
	for {
		if s.Peek() == ')' {
			break
		}
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		stmtList := node.(*StatementList)
		n.Body = append(n.Body, *stmtList)
	}
	// remove last bracket from stack
	if !s.PopVerifyString(")]") {
		return fmt.Errorf("error reading ExprStmt")
	}
	return nil
}

func (n *CaseStmt) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:%s", CaseExprValue, WriteCompactToken(n.Token(), d)))
	buf.WriteString(n.Condition.WriteCompact(d))
	buf.WriteString("(")
	for _, stmtList := range n.Body {
		buf.WriteString(stmtList.WriteCompact(d))
	}
	buf.WriteString(")]")
	return buf.String()
}

// ReadCompact reads a string and returns a ExitStmt
func (n *ExitStmt) ReadCompact(s *CharStack, d *IdList) error {
	// read token
	var err error

	// get token expression
	n.ExitKw, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}

	// remove comma
	if !s.PopVerify(',') {
		return fmt.Errorf("error reading ExitStmt")
	}

	// get ExitType expression
	// get token expression
	n.ExitType, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}

	// remove last bracket from stack
	if !s.PopVerify(']') {
		return fmt.Errorf("error reading ExitStmt")
	}

	return nil
}
func (n *ExitStmt) WriteCompact(d *Dict) string {
	return fmt.Sprintf("[%X:%s,%s]", ExitStmtValue, WriteCompactToken(n.ExitKw, d), WriteCompactToken(n.ExitType, d))
}

// ReadCompact reads a string and returns a SpecialStmt
func (n *SpecialStmt) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	// read token
	n.Keyword1, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}

	// get Keyword2 expression
	var k2 int
	fmt.Sscanf(s.String(), ",%X(", &k2)
	if !s.PopVerifyString(fmt.Sprintf(",%X(", k2)) {
		return fmt.Errorf("error reading SpecialStmt")
	}
	n.Keyword2 = d.get(k2)

	// get Args expressions
	for {
		if s.Peek() == ')' {
			break
		}
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		expr := node.(Expression)
		n.Args = append(n.Args, expr)
	}
	// remove last brackets from stack
	if !s.PopVerify(')') {
		return fmt.Errorf("error reading SpecialStmt")
	}
	if s.Peek() != ']' {
		// read semicolon
		n.Semicolon, err = ReadCompactToken(s, d)
		if err != nil {
			return err
		}
	}

	if !s.PopVerify(']') {
		return fmt.Errorf("error reading SpecialStmt")
	}
	return nil
}

func (n *SpecialStmt) WriteCompact(d *Dict) string {
	k2 := d.get(n.Keyword2)
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:%s", SpecialStmtValue, WriteCompactToken(n.Keyword1, d)))
	buf.WriteString(fmt.Sprintf(",%X(", k2))

	for _, arg := range n.Args {
		buf.WriteString(arg.WriteCompact(d))
	}
	buf.WriteString(")")
	if n.Semicolon != nil {
		buf.WriteString(WriteCompactToken(n.Semicolon, d))
	}
	buf.WriteString("]")
	return buf.String()
}

// ReadCompact reads a string and returns a WhileStmt
func (n *WhileStmt) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	// read token
	n.DoKw, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}

	// get Cond expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.Condition, _ = node.(Expression)
	// skip opening parenthesis
	if !s.PopVerify('(') {
		return fmt.Errorf("error reading WhileStmt")
	}

	// get Body expressions
	for {
		if s.Peek() == ')' {
			break
		}
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		stmtList := node.(*StatementList)
		n.Body = append(n.Body, *stmtList)
	}
	// remove last brackets from stack
	if !s.PopVerifyString(")]") {
		return fmt.Errorf("error reading WhileStmt")
	}
	return nil
}

func (n *WhileStmt) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:%s", WhileStmtValue, WriteCompactToken(n.Token(), d)))
	buf.WriteString(n.Condition.WriteCompact(d))
	buf.WriteString("(")
	for _, stmt := range n.Body {
		buf.WriteString(stmt.WriteCompact(d))
	}
	buf.WriteString(")]")
	return buf.String()
}

// ReadCompact reads a string and returns a UntilStmt
func (n *UntilStmt) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	var ok bool
	// read token
	n.DoKw, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}

	// get Cond expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.Condition, ok = node.(Expression)
	if !ok {
		return fmt.Errorf("error reading UntilStmt.Cond")
	}
	// skip opening parenthesis
	if !s.PopVerify('(') {
		return fmt.Errorf("error reading UntilStmt")
	}

	// get Body expressions
	for {
		if s.Peek() == ')' {
			break
		}
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		stmtList := node.(*StatementList)
		n.Body = append(n.Body, *stmtList)
	}
	// remove last brackets from stack
	if !s.PopVerifyString(")]") {
		return fmt.Errorf("error reading UntilStmt")
	}

	return nil
}

func (n *UntilStmt) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:%s", UntilStmtValue, WriteCompactToken(n.Token(), d)))
	buf.WriteString(n.Condition.WriteCompact(d))
	buf.WriteString("(")
	for _, stmtList := range n.Body {
		buf.WriteString(stmtList.WriteCompact(d))
	}
	buf.WriteString(")]")
	return buf.String()
}

// ReadCompact reads a string and returns a DoWhileStmt
func (n *DoWhileStmt) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	var ok bool
	// read token
	n.DoKw, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}

	// skip opening parenthesis
	if !s.PopVerify('(') {
		return fmt.Errorf("error reading DoWhileStmt")
	}

	// get Body expressions
	for {
		if s.Peek() == ')' {
			break
		}
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		stmtList := node.(*StatementList)
		n.Body = append(n.Body, *stmtList)
	}
	// verify last parenthesis
	if !s.PopVerify(')') {
		return fmt.Errorf("error reading DoWhileStmt")
	}

	// get Cond expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.Condition, ok = node.(Expression)
	if !ok {
		return fmt.Errorf("error reading UntilStmt.Cond")
	}

	// remove last brackets from stack
	if !s.PopVerify(']') {
		return fmt.Errorf("error reading DoWhileStmt")
	}
	return nil
}

func (n *DoWhileStmt) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:%s", DoWhileStmtValue, WriteCompactToken(n.Token(), d)))
	buf.WriteString("(")
	for _, stmt := range n.Body {
		buf.WriteString(stmt.WriteCompact(d))
	}
	buf.WriteString(")")
	buf.WriteString(n.Condition.WriteCompact(d))
	buf.WriteString("]")
	return buf.String()
}

// ReadCompact reads a string and returns a DoUntilStmt
func (n *DoUntilStmt) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	var ok bool
	// read token
	n.DoKw, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}

	// remove opening parenthesis
	if !s.PopVerify('(') {
		return fmt.Errorf("error reading DoUntilStmt")
	}

	// get Body expressions
	for {
		if s.Peek() == ')' {
			break
		}
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		stmtList := node.(*StatementList)
		n.Body = append(n.Body, *stmtList)
	}
	// remove last parenthesis from stack
	if !s.PopVerify(')') {
		return fmt.Errorf("error reading DoUntilStmt")
	}

	// get Cond expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.Condition, ok = node.(Expression)
	if !ok {
		return fmt.Errorf("error reading UntilStmt.Cond")
	}
	// remove last bracket from stack
	if !s.PopVerify(']') {
		return fmt.Errorf("error reading DoUntilStmt")
	}
	return nil
}

func (n *DoUntilStmt) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:%s", DoUntilStmtValue, WriteCompactToken(n.Token(), d)))
	buf.WriteString("(")
	for _, stmt := range n.Body {
		buf.WriteString(stmt.WriteCompact(d))
	}
	buf.WriteString(")")
	buf.WriteString(n.Condition.WriteCompact(d))
	buf.WriteString("]")
	return buf.String()
}

// ReadCompact reads a string and returns a ForStmt
func (n *ForStmt) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	var ok bool
	// read token
	n.ForKw, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}

	// get ForExpr expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.ForExpression, ok = node.(ForExpr)
	if !ok {
		return fmt.Errorf("error reading ForStmt.ForExpr")
	}
	// skip opening parenthesis
	if !s.PopVerify('(') {
		return fmt.Errorf("error reading ForStmt")
	}

	// get Body expressions
	for {
		if s.Peek() == ')' {
			break
		}
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		stmtList := node.(*StatementList)
		n.Body = append(n.Body, *stmtList)
	}
	// remove last parenthesis from stack
	if !s.PopVerify(')') {
		return fmt.Errorf("error reading ForStmt")
	}

	// get Next expression
	if s.Peek() != ']' {
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		n.Next, ok = node.(*Identifier)
		if !ok {
			return fmt.Errorf("error reading ForStmt.Next")
		}
	}
	// remove last bracket from stack
	if !s.PopVerify(']') {
		return fmt.Errorf("error reading ForStmt")
	}

	return nil
}

func (n *ForStmt) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:%s", ForStmtValue, WriteCompactToken(n.Token(), d)))
	buf.WriteString(n.ForExpression.WriteCompact(d))
	buf.WriteString("(")
	for _, stmt := range n.Body {
		buf.WriteString(stmt.WriteCompact(d))
	}
	buf.WriteString(")")
	if n.Next != nil {
		buf.WriteString(n.Next.WriteCompact(d))
	}
	buf.WriteString("]")
	return buf.String()
}

// ReadCompact reads a string and returns a DimDecl
func (n *DimDecl) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	// read token
	n.DimKw, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}

	// get Var expression
	// remove parenthesis from stack
	if !s.PopVerify('(') {
		return fmt.Errorf("error reading DimDecl")
	}

	for {
		// check if there is a closing parenthesis
		if s.Peek() == ')' {
			break
		}
		// get Var declaration
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		varDecl, ok := node.(VarDecl)
		if !ok {
			return fmt.Errorf("error reading DimDecl.Var")
		}
		n.Vars = append(n.Vars, varDecl)
	}
	// remove last bracket and last parenthese from stack
	if !s.PopVerifyString(")]") {
		return fmt.Errorf("error reading DimDecl")
	}
	return nil
}

func (n *DimDecl) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:%s(", DimDeclValue, WriteCompactToken(n.Token(), d)))
	for _, varDecl := range n.Vars {
		buf.WriteString(varDecl.WriteCompact(d))
	}
	buf.WriteString(")]")
	return buf.String()
}

// ReadCompact reads a string and returns a ConstDecl
func (n *ConstDecl) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	// read token
	n.ConstKw, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}
	// skip opening parenthesis
	if !s.PopVerify('(') {
		return fmt.Errorf("error reading ConstDecl")
	}

	// get Consts expressions
	for {
		if s.Peek() == ')' {
			break
		}
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		constDecl := node.(*ConstDeclItem)
		n.Consts = append(n.Consts, *constDecl)
	}
	// remove last brackets from stack
	if !s.PopVerifyString(")]") {
		return fmt.Errorf("error reading ConstDecl")
	}

	return nil
}

func (n *ConstDecl) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:%s(", ConstDeclValue, WriteCompactToken(n.Token(), d)))
	for _, constDecl := range n.Consts {
		buf.WriteString(constDecl.WriteCompact(d))
	}
	buf.WriteString(")]")
	return buf.String()
}

// ReadCompact reads a string and returns a ConstDeclItem
func (n *ConstDeclItem) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	var ok bool
	// read token
	_, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}

	// get ConstName expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.ConstName, ok = node.(*Identifier)
	if !ok {
		return fmt.Errorf("error reading ConstDeclItem.ConstName")
	}
	// get ConstType expression
	if s.Peek() != ',' {
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		n.ConstType, ok = node.(Type)
		if !ok {
			return fmt.Errorf("error reading ConstDeclItem.ConstType")
		}
	}
	// remove comma from stack
	if !s.PopVerify(',') {
		return fmt.Errorf("error reading ConstDeclItem")
	}

	// get ConstValue expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.ConstValue, ok = node.(Expression)
	if !ok {
		return fmt.Errorf("error reading ConstDeclItem.ConstValue")
	}
	// remove last bracket from stack
	if !s.PopVerify(']') {
		return fmt.Errorf("error reading ConstDeclItem")
	}

	return nil
}

func (n *ConstDeclItem) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:%s", ConstDeclItemValue, WriteCompactToken(n.Token(), d)))
	buf.WriteString(n.ConstName.WriteCompact(d))
	if n.ConstType != nil {
		buf.WriteString(n.ConstType.WriteCompact(d))
	}
	buf.WriteString(",")
	buf.WriteString(n.ConstValue.WriteCompact(d))
	buf.WriteString("]")
	return buf.String()
}

// ReadCompact reads a string and returns an EnumDecl
func (n *EnumDecl) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	var ok bool
	// read token
	n.EnumKw, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}

	// get EnumName expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.Identifier, ok = node.(*Identifier)
	if !ok {
		return fmt.Errorf("error reading EnumDecl.EnumName")
	}
	// skip opening parenthesis
	if !s.PopVerify('(') {
		return fmt.Errorf("error reading EnumDecl")
	}

	// get Values expressions
	for {
		if s.Peek() == ')' {
			break
		}
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		value := node.(*Identifier)
		n.Values = append(n.Values, *value)
	}
	// remove last brackets from stack
	if !s.PopVerifyString(")]") {
		return fmt.Errorf("error reading EnumDecl")
	}
	return nil
}

func (n *EnumDecl) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:%s", EnumDeclValue, WriteCompactToken(n.Token(), d)))
	buf.WriteString(n.Identifier.WriteCompact(d))
	buf.WriteString("(")
	for _, value := range n.Values {
		buf.WriteString(value.WriteCompact(d))
	}
	buf.WriteString(")]")
	return buf.String()
}

// ReadCompact reads a string and returns a EmptyStmt
func (n *LabelDecl) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var ok bool
	var node Node
	// read identifier
	// get EnumName expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.LabelName, ok = node.(*Identifier)
	if !ok {
		return fmt.Errorf("error reading LabelDecl.LabelName")
	}

	// remove last bracket from stack
	if !s.PopVerify(']') {
		return fmt.Errorf("error reading EmptyStmt")
	}
	return nil
}

func (n *LabelDecl) WriteCompact(d *Dict) string {
	return fmt.Sprintf("[%X:%s]", LabelDeclValue, n.LabelName.WriteCompact(d))
}

// ReadCompact reads a string and returns a ElseIfStmt
func (n *CallSubStmt) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	var ok bool
	// read token
	n.CallKw, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}

	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.Definition, ok = node.(*Identifier)
	if !ok {
		return fmt.Errorf("error reading CallSubStmt.Definition")
	}
	// remove last bracket from stack
	if !s.PopVerify(']') {
		return fmt.Errorf("error reading CallSubStmt")
	}
	return nil
}

func (n *CallSubStmt) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:%s", CallSubStmtValue, WriteCompactToken(n.Token(), d)))
	buf.WriteString(n.Definition.WriteCompact(d))
	buf.WriteString("]")
	return buf.String()
}

//	*BasicLit
//	*BinaryExpr
//	*Identifier
//	*CallOrIndexExpr
//	*ParenExpr
//	*UnaryExpr

// ReadCompact reads a string and returns a BasicLit
func (n *BasicLit) ReadCompact(s *CharStack, d *IdList) error {
	var nVal int
	// read token
	var err error
	n.ValPos, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}
	// remove comma
	if !s.PopVerify(',') {
		return fmt.Errorf("error reading BasicLit")
	}

	fmt.Sscanf(s.String(), "%X", &nVal)
	n.Value = d.get(nVal)

	// calculate len of read string
	s.PopVerifyString(fmt.Sprintf("%X", nVal))

	// remove comma
	if !s.PopVerify(',') {
		return fmt.Errorf("error reading BasicLit")
	}
	var nKind int
	fmt.Sscanf(s.String(), "%X]", &nKind)
	n.Kind = token.Kind(nKind)

	// remove string from stack
	if !s.PopVerifyString(fmt.Sprintf("%X", nKind)) {
		return fmt.Errorf("error reading BasicLit")
	}

	// remove last bracket from stack
	if !s.PopVerify(']') {
		return fmt.Errorf("error reading BasicLit")
	}

	return nil
}

func (n *BasicLit) WriteCompact(d *Dict) string {
	v := d.get(n.String())
	return fmt.Sprintf("[%X:%s,%X,%X]", BasicLitValue, WriteCompactToken(n.ValPos, d), v, int(n.Kind))
}

// ReadCompact reads a string and returns a BinaryExpr
func (n *BinaryExpr) ReadCompact(s *CharStack, d *IdList) error {
	var ok bool
	var err error
	// read token
	n.OpToken, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}

	// skip comma
	if !s.PopVerify(',') {
		return fmt.Errorf("error reading BinaryExpr")
	}

	// read n.Left expression
	nLeft, err := ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.Left, ok = nLeft.(Expression)
	if !ok {
		return fmt.Errorf("error reading BinaryExpr.X")
	}
	// skip comma
	if !s.PopVerify(',') {
		return fmt.Errorf("error reading BinaryExpr")
	}

	// read kind
	var nKind uint16
	fmt.Sscanf(s.String(), "%X", &nKind)
	n.OpKind = token.Kind(nKind)
	s.PopVerifyString(fmt.Sprintf("%X", nKind))

	// skip comma
	if !s.PopVerify(',') {
		return fmt.Errorf("error reading BinaryExpr")
	}

	// read n.Right expression
	nRight, err := ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.Right, ok = nRight.(Expression)
	if !ok {
		return fmt.Errorf("error reading BinaryExpr.Y")
	}

	// remove last bracket from stack
	if !s.PopVerify(']') {
		return fmt.Errorf("error reading BinaryExpr")
	}

	return nil
}

func (n *BinaryExpr) WriteCompact(d *Dict) string {
	return fmt.Sprintf("[%X:%s,%s,%X,%v]", BinaryExprValue, WriteCompactToken(n.OpToken, d), n.Left.WriteCompact(d), uint16(n.OpKind), n.Right.WriteCompact(d))
}

// ReadCompact reads a string and returns an Ident
func (n *Identifier) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	// read token
	n.Tok, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}

	// check if there is a declaration
	if s.Peek() != ',' {
		decl, err := ReadCompact(s, d)
		if err != nil {
			return err
		}
		n.Decl = decl.(Decl)
	}
	// skip comma
	if !s.PopVerify(',') {
		return fmt.Errorf("error reading Ident")
	}

	// read n.Name
	var nName int
	fmt.Sscanf(s.String(), "%X]", &nName)
	n.Name = d.get(nName)
	// calculate len of read string
	length := len(fmt.Sprintf("%X", nName))
	// remove read string from stack
	_ = s.PopString(length)
	if !s.PopVerify(']') {
		return fmt.Errorf("error reading Ident")
	}
	return nil
}

func (n *Identifier) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:%s", IdentifierValue, WriteCompactToken(n.Token(), d)))

	if n.Decl != nil {
		// type declaration structure points to itself
		// so we need to check if the declaration is not a type
		// and if the position of the declaration is different
		// from the position of the identifier
		scalarDecl, ok := n.Decl.(*ScalarDecl)
		if ok {
			pos := scalarDecl.VarName.Tok.Position
			if pos != n.Tok.Position {
				n.Decl.WriteCompact(d)
			}
			// } else {
			// 	enumDecl, ok := n.Decl.(*EnumDecl)
			// 	if ok {
			// 		pos := enumDecl.Enum
			// 		if pos != n.NamePos {
			// 			n.Decl.WriteCompact(d)
			// 		}
			// 	} else {
			// 		_, ok := n.Decl.(*TypeDef)
			// 		if !ok {
			// 			n.Decl.WriteCompact(d)
			// 		}
			// 	}
		}
	}
	v := d.get(n.Name)
	buf.WriteString(fmt.Sprintf(",%X]", v))
	return buf.String()
}

// ReadCompact reads a string and returns a CallOrIndexExpr
func (n *CallOrIndexExpr) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	var ok bool
	// get Name expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.Identifier, ok = node.(*Identifier)
	if !ok {
		return fmt.Errorf("error reading CallOrIndexExpr.Name")
	}
	// skip opening parenthesis
	if !s.PopVerify('(') {
		return fmt.Errorf("error reading CallOrIndexExpr")
	}

	// get Args expressions
	for {
		if s.Peek() == ')' {
			break
		}
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		expr := node.(Expression)
		n.Args = append(n.Args, expr)
	}
	// remove last brackets from stack
	if !s.PopVerifyString(")]") {
		return fmt.Errorf("error reading CallOrIndexExpr")
	}

	return nil
}

func (n *CallOrIndexExpr) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:", CallOrIndexExprValue))
	buf.WriteString(n.Identifier.WriteCompact(d))
	buf.WriteString("(")
	for _, arg := range n.Args {
		buf.WriteString(arg.WriteCompact(d))
	}
	buf.WriteString(")]")
	return buf.String()
}

// ReadCompact reads a string and returns a ParenExpr
func (n *ParenExpr) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	// read token
	n.Lparen, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}
	// skip comma
	if !s.PopVerify(',') {
		return fmt.Errorf("error reading ParenExpr")
	}

	// get X expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.Expr = node.(Expression)

	// skip comma
	if !s.PopVerify(',') {
		return fmt.Errorf("error reading ParenExpr")
	}

	// read right parenthesis
	var pos *token.Position
	pos, err = ReadCompactPosition(s, d)
	if err != nil {
		return err
	}
	n.Rparen = *pos

	// remove last bracket from stack
	if !s.PopVerify(']') {
		return fmt.Errorf("error reading ParenExpr")
	}

	return nil
}

func (n *ParenExpr) WriteCompact(d *Dict) string {
	return fmt.Sprintf("[%X:%s,%v,%v]", ParenExprValue, WriteCompactToken(n.Lparen, d), n.Expr.WriteCompact(d), WriteCompactPosition(&n.Rparen, d))
}

// ReadCompact reads a string and returns a UnaryExpr
func (n *UnaryExpr) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	// read token
	n.OpToken, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}

	// skip comma
	if !s.PopVerify(',') {
		return fmt.Errorf("error reading UnaryExpr")
	}

	// get X expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.Right = node.(Expression)

	// skip comma
	if !s.PopVerify(',') {
		return fmt.Errorf("error reading UnaryExpr")
	}

	// read kind
	var nKind uint16
	fmt.Sscanf(s.String(), "%X]", &nKind)
	n.OpKind = token.Kind(nKind)
	s.PopVerifyString(fmt.Sprintf("%X", nKind))

	// remove last bracket from stack
	if !s.PopVerify(']') {
		return fmt.Errorf("error reading UnaryExpr")
	}

	return nil
}

func (n *UnaryExpr) WriteCompact(d *Dict) string {
	return fmt.Sprintf("[%X:%s,%v,%X]", UnaryExprValue, WriteCompactToken(n.OpToken, d), n.Right.WriteCompact(d), uint16(n.OpKind))
}

// ReadCompact reads a string and returns a ForNextExpr
func (n *ForNextExpr) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	var ok bool
	// get Var expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.Variable, ok = node.(*Identifier)
	if !ok {
		return fmt.Errorf("error reading ForNextExpr.Var")
	}

	// get From expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.From, ok = node.(Expression)
	if !ok {
		return fmt.Errorf("error reading ForNextExpr.From")
	}
	// get To expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.To, ok = node.(Expression)
	if !ok {
		return fmt.Errorf("error reading ForNextExpr.To")
	}
	// get StepVal expression
	if s.Peek() != ']' {
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		n.Step, ok = node.(Expression)
		if !ok {
			return fmt.Errorf("error reading ForNextExpr.StepVal")
		}
	}
	// remove last bracket from stack
	if !s.PopVerify(']') {
		return fmt.Errorf("error reading ForNextExpr")
	}

	return nil
}

func (n *ForNextExpr) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:", ForNextExprValue))
	buf.WriteString(n.Variable.WriteCompact(d))
	buf.WriteString(n.From.WriteCompact(d))
	buf.WriteString(n.To.WriteCompact(d))
	if n.Step != nil {
		buf.WriteString(n.Step.WriteCompact(d))
	}
	buf.WriteString("]")
	return buf.String()
}

// ReadCompact reads a string and returns a ForEachExpr
func (n *ForEachExpr) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	var ok bool
	// get Var expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.Variable, ok = node.(*Identifier)
	if !ok {
		return fmt.Errorf("error reading ForEachExpr.Var")
	}
	// get Collection expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.Collection, ok = node.(Expression)
	if !ok {
		return fmt.Errorf("error reading ForEachExpr.Collection")
	}
	// remove last bracket from stack
	if !s.PopVerify(']') {
		return fmt.Errorf("error reading ForEachExpr")
	}
	return nil
}

func (n *ForEachExpr) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:", ForEachExprValue))
	buf.WriteString(n.Variable.WriteCompact(d))
	buf.WriteString(n.Collection.WriteCompact(d))
	buf.WriteString("]")
	return buf.String()
}

// ReadCompact reads a string and returns a CallSelectorExpr
func (n *CallSelectorExpr) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	var ok bool
	// get Root expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.Root, ok = node.(Expression)
	if !ok {
		return fmt.Errorf("error reading CallSelectorExpr.Root")
	}
	// get Selector expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.Selector, ok = node.(*Identifier)
	if !ok {
		return fmt.Errorf("error reading CallSelectorExpr.Selector")
	}
	// remove last bracket from stack
	if !s.PopVerify(']') {
		return fmt.Errorf("error reading CallSelectorExpr")
	}

	return nil
}

func (n *CallSelectorExpr) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:", CallSelectorExprValue))
	buf.WriteString(n.Root.WriteCompact(d))
	buf.WriteString(n.Selector.WriteCompact(d))
	buf.WriteString("]")
	return buf.String()
}

// ReadCompact reads a string and returns a UserDefinedType
func (n *UserDefinedType) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	var ok bool
	// get X expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.Identifier, ok = node.(*Identifier)
	if !ok {
		return fmt.Errorf("error reading UserDefinedType.Name")
	}
	// remove last bracket from stack
	if !s.PopVerify(']') {
		return fmt.Errorf("error reading UserDefinedType")
	}
	return nil
}

func (n *UserDefinedType) WriteCompact(d *Dict) string {
	return fmt.Sprintf("[%X:%v]", UserDefinedTypeValue, n.Identifier.WriteCompact(d))
}

// *FuncDecl
// *SubDecl
// *VarDecl
// *EnumDecl
// *ConstDeclItem
// *ArrayDecl
// *DimDecl
// *ScalarDecl
// *LabelDecl

// ReadCompact reads a string and returns a FuncDecl
func (n *FuncDecl) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	var ok bool
	// read token
	n.FunctionKw, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}

	// get FuncName expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.FuncName, ok = node.(*Identifier)
	if !ok {
		return fmt.Errorf("error reading FuncDecl.FuncName")
	}
	// get FuncType expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.FuncType, ok = node.(*FuncType)
	if !ok {
		return fmt.Errorf("error reading FuncDecl.FuncType")
	}
	// skip opening parenthesis
	if !s.PopVerify('(') {
		return fmt.Errorf("error reading FuncDecl")
	}

	// get Body expressions
	for {
		if s.Peek() == ')' {
			break
		}
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		stmtList := node.(*StatementList)
		n.Body = append(n.Body, *stmtList)
	}
	// remove last brackets from stack
	if !s.PopVerifyString(")]") {
		return fmt.Errorf("error reading FuncDecl")
	}
	return nil
}

func (n *FuncDecl) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:%s", FuncDeclValue, WriteCompactToken(n.Token(), d)))
	buf.WriteString(n.FuncName.WriteCompact(d))
	buf.WriteString(n.FuncType.WriteCompact(d))
	buf.WriteString("(")
	for _, stmt := range n.Body {
		buf.WriteString(stmt.WriteCompact(d))
	}
	buf.WriteString(")]")
	return buf.String()
}

// ReadCompact reads a string and returns a SubDecl
func (n *SubDecl) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	var ok bool
	// read token
	n.SubKw, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}

	// get SubName expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.SubName, ok = node.(*Identifier)
	if !ok {
		return fmt.Errorf("error reading SubDecl.SubName")
	}

	// skip comma
	if !s.PopVerify(',') {
		return fmt.Errorf("error reading SubDecl")
	}

	// get SubType expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.SubType, ok = node.(*SubType)
	if !ok {
		return fmt.Errorf("error reading SubDecl.SubType")
	}
	// skip opening parenthesis
	if !s.PopVerify('(') {
		return fmt.Errorf("error reading SubDecl")
	}

	// get Body expressions
	for {
		if s.Peek() == ')' {
			break
		}
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		stmtList := node.(*StatementList)
		n.Body = append(n.Body, *stmtList)
	}
	// remove last brackets from stack
	if !s.PopVerifyString(")]") {
		return fmt.Errorf("error reading SubDecl")
	}
	return nil
}

func (n *SubDecl) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:%s", SubDeclValue, WriteCompactToken(n.Token(), d)))
	buf.WriteString(n.SubName.WriteCompact(d))
	buf.WriteString(",")
	buf.WriteString(n.SubType.WriteCompact(d))
	buf.WriteString("(")
	for _, stmt := range n.Body {
		buf.WriteString(stmt.WriteCompact(d))
	}
	buf.WriteString(")]")
	return buf.String()
}

// ReadCompact reads a string and returns a FuncType
func (n *FuncType) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	var ok bool
	// read token
	n.Lparen, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}

	// skip opening parenthesis
	if !s.PopVerify('(') {
		return fmt.Errorf("error reading FuncType")
	}

	// get Params expressions
	for {
		if s.Peek() == ')' {
			break
		}
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		param := node.(*ParamItem)
		n.Params = append(n.Params, *param)
	}
	// remove last parenthesis from stack
	if !s.PopVerify(')') {
		return fmt.Errorf("error reading FuncType")
	}

	// get Result expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.Result, ok = node.(Type)
	if !ok {
		return fmt.Errorf("error reading FuncType.Result")
	}

	// read right parenthesis
	var pos *token.Position
	pos, err = ReadCompactPosition(s, d)
	n.Rparen = *pos
	if err != nil {
		return err
	}
	// remove last bracket from stack
	if !s.PopVerify(']') {
		return fmt.Errorf("error reading FuncType")
	}
	return nil
}

func (n *FuncType) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:%s(", FuncTypeValue, WriteCompactToken(n.Token(), d)))
	for _, param := range n.Params {
		buf.WriteString(param.WriteCompact(d))
	}
	buf.WriteString(")")
	buf.WriteString(n.Result.WriteCompact(d))
	buf.WriteString(WriteCompactPosition(&n.Rparen, d))
	buf.WriteString("]")
	return buf.String()
}

// ReadCompact reads a string and returns a FuncType
func (n *SubType) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	// read token
	n.Lparen, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}

	// skip opening parenthesis
	if !s.PopVerify('(') {
		return fmt.Errorf("error reading FuncType")
	}

	// get Params expressions
	for {
		if s.Peek() == ')' {
			break
		}
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		param := node.(*ParamItem)
		n.Params = append(n.Params, *param)
	}
	// remove last parenthesis from stack
	if !s.PopVerifyString(")") {
		return fmt.Errorf("error reading FuncType")
	}
	var pos *token.Position
	pos, err = ReadCompactPosition(s, d)
	n.Rparen = *pos
	if err != nil {
		return err
	}

	// remove last bracket from stack
	if !s.PopVerifyString("]") {
		return fmt.Errorf("error reading FuncType")
	}
	return nil
}

func (n *SubType) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:%s(", SubTypeValue, WriteCompactToken(n.Token(), d)))
	for _, param := range n.Params {
		buf.WriteString(param.WriteCompact(d))
	}
	buf.WriteString(")")
	buf.WriteString(WriteCompactPosition(&n.Rparen, d))
	buf.WriteString("]")
	return buf.String()
}

// ReadCompact reads a string and returns an ArrayDecl
func (n *ArrayDecl) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	var ok bool
	// get VarName expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.VarName, ok = node.(*Identifier)
	if !ok {
		return fmt.Errorf("error reading ArrayDecl.VarName")
	}
	// get VarType expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.VarType, ok = node.(*ArrayType)
	if !ok {
		return fmt.Errorf("error reading ArrayDecl.VarType")
	}
	// remove last bracket from stack
	if !s.PopVerify(']') {
		return fmt.Errorf("error reading ArrayDecl")
	}
	return nil
}

func (n *ArrayDecl) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:", ArrayDeclValue))
	buf.WriteString(n.VarName.WriteCompact(d))
	buf.WriteString(n.VarType.WriteCompact(d))
	buf.WriteString("]")
	return buf.String()
}

// ReadCompact reads a string and returns an ArrayType
func (n *ArrayType) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	var ok bool
	// read token
	n.Lparen, err = ReadCompactToken(s, d)
	if err != nil {
		return err
	}

	// skip opening parenthesis
	if !s.PopVerify('(') {
		return fmt.Errorf("error reading ArrayType")
	}

	// get Dimensions expressions
	for {
		if s.Peek() == ')' {
			break
		}
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		dim := node.(Expression)
		n.Dimensions = append(n.Dimensions, dim)
	}
	// remove parenthese from stack
	if !s.PopVerify(')') {
		return fmt.Errorf("error reading ArrayType")
	}
	// get Type expression
	if s.Peek() != ',' {
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		n.Type, ok = node.(Type)
		if !ok {
			return fmt.Errorf("error reading ArrayType.Type")
		}
	}
	// remove comma from stack
	if !s.PopVerify(',') {
		return fmt.Errorf("error reading ArrayType")
	}

	n.Rparen, err = ReadCompactPosition(s, d)
	if err != nil {
		return err
	}
	// remove last bracket from stack
	if !s.PopVerify(']') {
		return fmt.Errorf("error reading ArrayType")
	}

	return nil
}

func (n *ArrayType) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:%s", ArrayTypeValue, WriteCompactToken(n.Token(), d)))
	buf.WriteString("(")
	for _, dim := range n.Dimensions {
		buf.WriteString(dim.WriteCompact(d))
	}
	buf.WriteString(")")
	if n.Type != nil {
		buf.WriteString(n.Type.WriteCompact(d))
	}
	buf.WriteString(fmt.Sprintf(",%s]", WriteCompactPosition(n.Rparen, d)))
	return buf.String()
}

// ReadCompact reads a string and returns a ScalarDecl
func (n *ScalarDecl) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	var ok bool
	// get VarName expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.VarName, ok = node.(*Identifier)
	if !ok {
		return fmt.Errorf("error reading ScalarDecl.VarName")
	}

	// remove comma from stack
	if !s.PopVerify(',') {
		return fmt.Errorf("error reading ScalarDecl")
	}

	// get VarType expression
	if s.Peek() != ',' {
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		n.VarType, ok = node.(Type)
		if !ok {
			return fmt.Errorf("error reading Type.VarType")
		}
	}
	// remove comma from stack
	if !s.PopVerify(',') {
		return fmt.Errorf("error reading ScalarDecl")
	}

	// get VarValue expression
	if s.Peek() != ']' {
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		n.VarValue, ok = node.(Expression)
		if !ok {
			return fmt.Errorf("error reading ScalarDecl.VarValue")
		}
	}
	// remove last bracket from stack
	if !s.PopVerify(']') {
		return fmt.Errorf("error reading ScalarDecl")
	}

	return nil
}

func (n *ScalarDecl) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:", ScalarDeclValue))
	buf.WriteString(n.VarName.WriteCompact(d))
	buf.WriteString(",")
	if n.VarType != nil {
		buf.WriteString(n.VarType.WriteCompact(d))
	}
	buf.WriteString(",")
	if n.VarValue != nil {
		buf.WriteString(n.VarValue.WriteCompact(d))
	}
	buf.WriteString("]")
	return buf.String()
}

// ReadCompact reads a string and returns a ParamItem
func (n *ParamItem) ReadCompact(s *CharStack, d *IdList) error {
	var err error
	var node Node
	var ok bool
	// get Optional expression
	nOptional := s.Pop()
	n.Optional = nOptional == '1'
	// get ByRef expression
	nByVal := s.Pop()
	n.ByVal = nByVal == '1'
	// get ParamArray expression
	nParamArray := s.Pop()
	n.ParamArray = nParamArray == '1'
	// get IsArray expression
	nIsArray := s.Pop()
	n.IsArray = nIsArray == '1'

	// get VarName expression
	node, err = ReadCompact(s, d)
	if err != nil {
		return err
	}
	n.VarName, ok = node.(*Identifier)
	if !ok {
		return fmt.Errorf("error reading ParamItem.VarName")
	}
	// remove comma from stack
	if !s.PopVerify(',') {
		return fmt.Errorf("error reading ParamItem")
	}

	// get VarType expression
	if s.Peek() != ',' {
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		n.VarType, ok = node.(Type)
		if !ok {
			return fmt.Errorf("error reading ParamItem.VarType")
		}
	}
	// remove comma from stack
	if !s.PopVerify(',') {
		return fmt.Errorf("error reading ParamItem")
	}

	// get VarValue expression
	if s.Peek() != ']' {
		node, err = ReadCompact(s, d)
		if err != nil {
			return err
		}
		n.DefaultValue, ok = node.(Expression)
		if !ok {
			return fmt.Errorf("error reading ParamItem.VarValue")
		}
	}
	// remove last bracket from stack
	if !s.PopVerify(']') {
		return fmt.Errorf("error reading ParamItem")
	}
	return nil
}

func (n *ParamItem) WriteCompact(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:", ParamItemValue))
	if n.Optional {
		buf.WriteString("1")
	} else {
		buf.WriteString("0")
	}
	if n.ByVal {
		buf.WriteString("1")
	} else {
		buf.WriteString("0")
	}
	if n.ParamArray {
		buf.WriteString("1")
	} else {
		buf.WriteString("0")
	}
	if n.IsArray {
		buf.WriteString("1")
	} else {
		buf.WriteString("0")
	}
	buf.WriteString(n.VarName.WriteCompact(d))
	buf.WriteString(",")
	if n.VarType != nil {
		buf.WriteString(n.VarType.WriteCompact(d))
	}
	buf.WriteString(",")
	if n.DefaultValue != nil {
		buf.WriteString(n.DefaultValue.WriteCompact(d))
	}
	buf.WriteString("]")
	return buf.String()
}

// WriteCompactDict writes the dictionary to a string
func WriteCompactDict(d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString("[")
	for k, v := range *d {
		buf.WriteString(fmt.Sprintf("(%X:%q)", v, k))
	}
	buf.WriteString("]")
	return buf.String()
}

// ReacCompactDict reads the ident list from a string
func ReadCompactDict(cs *CharStack) IdList {
	d := make(IdList)
	// remove opening bracket
	if !cs.PopVerify('[') {
		return nil
	}

	// get strings
	for {
		if cs.Peek() == ']' {
			break
		}
		var k string
		var v int
		fmt.Sscanf(cs.String(), "(%X:%q)", &v, &k)
		d[v] = k
		// calculate len of read string
		length := len(fmt.Sprintf("(%X:%q)", v, k))
		// remove read string from stack
		_ = cs.PopString(length)
	}
	// remove last bracket from stack
	if !cs.PopVerify(']') {
		return nil
	}

	return d
}

// get info from dictionary. If not found, add it
func (d *Dict) get(value string) int {
	if _, ok := (*d)[value]; !ok {
		(*d)[value] = len(*d)
	}
	return (*d)[value]
}

// CompresssToFile compresses a string to a file
func CompresssToFile(s string, filename string) error {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)

	// Add the file to the archive.
	f, err := w.Create(filename)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(s))
	if err != nil {
		return err
	}

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		return err
	}

	file, err := os.Create(filename + ".zip")
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = buf.WriteTo(file)
	return err
}

// DecompressFile decompresses a file
func DecompressFile(filename string) (string, error) {
	r, err := zip.OpenReader(filename)
	if err != nil {
		return "", err
	}
	defer r.Close()
	inFile := filename[:len(filename)-4] // remove .zip
	for _, f := range r.File {
		if f.FileHeader.Name == inFile {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()
			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(rc)
			if err != nil {
				return "", err
			}
			return buf.String(), nil
		}
	}
	return "", fmt.Errorf("file not found in zip")
}

func identifyNodeType(s *CharStack) int {
	var n int
	fmt.Sscanf(s.String(), "[%X:", &n)
	if !s.PopVerify('[') {
		return -n
	}

	// calculate len of read string
	length := len(fmt.Sprintf("%X", n))
	// remove read string from stack
	_ = s.PopString(length)
	if !s.PopVerify(':') {
		return -n
	}
	return n
}

// InitReadCompact reads a string and returns an initialized Node
func InitReadCompact(s *CharStack, d *IdList, filename string) (*File, error) {
	n, err := InitializeNode(s)
	if err != nil {
		return nil, err
	}
	// set filename
	f := n.(*File)
	f.Name = filename

	err = f.ReadCompact(s, d)
	return f, err
}

// ReadCompact reads a string and returns an initialized Node
func ReadCompact(s *CharStack, d *IdList) (Node, error) {
	n, err := InitializeNode(s)
	if err != nil {
		return nil, err
	}
	err = n.ReadCompact(s, d)
	return n, err
}

// for type compatibility only
func (n *TypeDef) ReadCompact(s *CharStack, d *IdList) error {
	panic("not implemented")
}

func (n *TypeDef) WriteCompact(d *Dict) string {
	panic("not implemented")
}

// for type compatibility only
func (n *ClassDecl) ReadCompact(s *CharStack, d *IdList) error {
	panic("not implemented")
}

func (n *ClassDecl) WriteCompact(d *Dict) string {
	panic("not implemented")
}

// Initialize Node, instanciate a node from a string
func InitializeNode(s *CharStack) (Node, error) {
	n := identifyNodeType(s)
	switch n {
	case FileValue:
		return &File{}, nil
	case StatementListValue:
		return &StatementList{}, nil
	// statements
	case EmptyStmtValue:
		return &EmptyStmt{}, nil
	case ExprStmtValue:
		return &ExprStmt{}, nil
	case IfStmtValue:
		return &IfStmt{}, nil
	case ElseIfStmtValue:
		return &ElseIfStmt{}, nil
	case SelectStmtValue:
		return &SelectStmt{}, nil
	case CaseExprValue:
		return &CaseStmt{}, nil
	case ExitStmtValue:
		return &ExitStmt{}, nil
	case SpecialStmtValue:
		return &SpecialStmt{}, nil
	case WhileStmtValue:
		return &WhileStmt{}, nil
	case UntilStmtValue:
		return &UntilStmt{}, nil
	case DoWhileStmtValue:
		return &DoWhileStmt{}, nil
	case DoUntilStmtValue:
		return &DoUntilStmt{}, nil
	case ForStmtValue:
		return &ForStmt{}, nil
	case ForNextExprValue:
		return &ForNextExpr{}, nil
	case ForEachExprValue:
		return &ForEachExpr{}, nil
	case DimDeclValue:
		return &DimDecl{}, nil
	case ConstDeclValue:
		return &ConstDecl{}, nil
	case ConstDeclItemValue:
		return &ConstDeclItem{}, nil
	case EnumDeclValue:
		return &EnumDecl{}, nil
	case LabelDeclValue:
		return &LabelDecl{}, nil
	case CallSubStmtValue:
		return &CallSubStmt{}, nil
	case BasicLitValue:
		return &BasicLit{}, nil
	// expressions
	case BinaryExprValue:
		return &BinaryExpr{}, nil
	case IdentifierValue:
		return &Identifier{}, nil
	case CallOrIndexExprValue:
		return &CallOrIndexExpr{}, nil
	case ParenExprValue:
		return &ParenExpr{}, nil
	case UnaryExprValue:
		return &UnaryExpr{}, nil
	case CallSelectorExprValue:
		return &CallSelectorExpr{}, nil
	case UserDefinedTypeValue:
		return &UserDefinedType{}, nil
	case FuncDeclValue:
		return &FuncDecl{}, nil
	case SubDeclValue:
		return &SubDecl{}, nil
	case FuncTypeValue:
		return &FuncType{}, nil
	case SubTypeValue:
		return &SubType{}, nil
	case ArrayDeclValue:
		return &ArrayDecl{}, nil
	case ArrayTypeValue:
		return &ArrayType{}, nil
	case ScalarDeclValue:
		return &ScalarDecl{}, nil

	case ParamItemValue:
		return &ParamItem{}, nil
		// case TypeDefValue:
		// 	return &TypeDef{}, nil
	}
	return nil, fmt.Errorf("unknown node type %v", n)
}

// get info from IdList
func (d *IdList) get(value int) string {
	return (*d)[value]
}

// LoadFile reads a file and return a File node
func LoadFile(filename string) (*File, error) {
	var err error

	// read compressed file
	// TODO: implement GZIP compression
	//s, err = DecompressFile(filename)
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// read file content into stack
	cs := CharStack{}
	cs.PushString(string(buf))

	// read dictionary
	idList := ReadCompactDict(&cs)
	// read file

	file, err := InitReadCompact(&cs, &idList, filename)
	if err != nil {
		return nil, err
	}

	// ensure that the stack is empty
	if cs.Len() != 0 {
		return nil, fmt.Errorf("stack not empty")
	}
	return file, nil
}

//  -----------------------------------------------------------------
//  token.Token
//  -----------------------------------------------------------------

// ReadCompactToken reads a string and returns a token.Token
func ReadCompactToken(s *CharStack, d *IdList) (*token.Token, error) {
	t := token.Token{}
	// remove header
	if !s.PopVerifyString(fmt.Sprintf("[%X:", TokenValue)) {
		return nil, fmt.Errorf("error reading token.Token")
	}

	// read literal
	var tLiteral int
	fmt.Sscanf(s.String(), "%X,", &tLiteral)
	t.Literal = d.get(tLiteral)
	// remove read string from stack
	if !s.PopVerifyString(fmt.Sprintf("%X,", tLiteral)) {
		return nil, fmt.Errorf("error reading token.Token")
	}

	// read kind
	var tKind uint16
	fmt.Sscanf(s.String(), "%X", &tKind)
	t.Kind = token.Kind(tKind)
	if !s.PopVerifyString(fmt.Sprintf("%X", tKind)) {
		return nil, fmt.Errorf("error reading token.Token")
	}

	// read position
	var err error
	var tPosition *token.Position
	tPosition, err = ReadCompactPosition(s, d)
	t.Position = *tPosition
	if err != nil {
		return nil, err
	}
	s.PopVerify(']')

	return &t, nil
}

func WriteCompactToken(t *token.Token, d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%X:", TokenValue))
	v := d.get(t.Literal)
	buf.WriteString(fmt.Sprintf("%X,", v)) // literal
	buf.WriteString(fmt.Sprintf("%X", uint16(t.Kind)))
	buf.WriteString(WriteCompactPosition(&t.Position, d))
	buf.WriteString("]")
	return buf.String()
}

func ReadCompactPosition(s *CharStack, d *IdList) (*token.Position, error) {
	var tPosition token.Position
	fmt.Sscanf(s.String(), "(%X,", &tPosition.Line)
	if !s.PopVerifyString(fmt.Sprintf("(%X,", tPosition.Line)) {
		return nil, fmt.Errorf("error reading token.Position")
	}

	fmt.Sscanf(s.String(), "%X,", &tPosition.Column)
	if !s.PopVerifyString(fmt.Sprintf("%X,", tPosition.Column)) {
		return nil, fmt.Errorf("error reading token.Position")
	}

	fmt.Sscanf(s.String(), "%X)", &tPosition.Absolute)
	if !s.PopVerifyString(fmt.Sprintf("%X)", tPosition.Absolute)) {
		return nil, fmt.Errorf("error reading token.Position")
	}

	return &tPosition, nil
}

func WriteCompactPosition(p *token.Position, d *Dict) string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("(%X,", p.Line))
	buf.WriteString(fmt.Sprintf("%X,", p.Column))
	buf.WriteString(fmt.Sprintf("%X)", p.Absolute))
	return buf.String()
}
