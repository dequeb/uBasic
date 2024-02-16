package ast

// EqualsNode returns true if the two nodes are equal.
func (n *ArrayType) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*ArrayType)
	if !ok {
		return n
	}

	if len(n.Dimensions) != len(n2.Dimensions) {
		return n
	}
	for i, v := range n.Dimensions {
		dim := n2.Dimensions[i].(Node) // Change the type of dim from Expr to *Node
		if res := v.EqualsNode(dim); res != nil {
			return res
		}
	}
	var res Node
	if n.Type != nil {
		res = n.Type.EqualsNode(n2.Type)
		if res != nil {
			return res
		}
	} else {
		if n2.Type != nil {
			return n
		}
	}
	if n.Lparen == nil || n.Rparen == nil {
		return n
	}
	if !n.Lparen.Equals(n2.Lparen) {
		return n
	}
	if !n.Rparen.Equals(n2.Rparen) {
		return n
	}
	return nil
}

// EqualsNode returns true if the two nodes are equal.
func (n *BasicLit) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*BasicLit)
	if !ok {
		return n
	}
	if n.Kind == n2.Kind && n.Value == n2.Value && n.ValPos.Equals(n2.ValPos) {
		return nil
	}
	return n
}

// EqualsNode returns true if the two nodes are equal.
func (n *BinaryExpr) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*BinaryExpr)
	if !ok {
		return n
	}
	if n.Left != nil {
		res := n.Left.EqualsNode(n2.Left)
		if res != nil {
			return res
		}
	} else {
		if n2.Left != nil {
			return n
		}
	}

	if n.Right != nil {
		res := n.Right.EqualsNode(n2.Right)
		if res != nil {
			return res
		}
	} else {
		if n2.Right != nil {
			return n
		}
	}
	if n.OpKind != n2.OpKind {
		return n
	}
	return nil
}

// EqualsNode returns true if the two nodes are equal.
func (n *ExitStmt) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*ExitStmt)
	if !ok {
		return n
	}

	if n.ExitKw == nil {
		return n
	}
	if n.ExitType == n2.ExitType && n.ExitKw.Equals(n2.ExitKw) {
		return nil
	}
	return n
}

// EqualsNode returns true if the two nodes are equal.
func (n *SpecialStmt) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*SpecialStmt)
	if !ok {
		return n
	}

	if len(n.Args) != len(n2.Args) {
		return n
	}
	for i, v := range n.Args {
		res := v.EqualsNode(n2.Args[i])
		if res != nil {
			return res
		}
	}

	if n.Keyword1 == nil || n2.Keyword1 == nil {
		return n
	}
	if n.Semicolon != nil {
		if !n.Semicolon.Equals(n2.Semicolon) {
			return n
		}
	} else {
		if n2.Semicolon != nil {
			return n
		}
	}

	if n.Keyword1.Equals(n2.Keyword1) && n.Keyword2 == n2.Keyword2 {
		return nil
	}

	return n
}

// EqualsNode returns true if the two nodes are equal.
func (n *CallSubStmt) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}

	n2, ok := node.(*CallSubStmt)
	if !ok {
		return n
	}

	if len(n.Args) != len(n2.Args) {
		return n
	}

	for i, v := range n.Args {
		res := v.EqualsNode(n2.Args[i])
		if res != nil {
			return res
		}
	}
	if n.CallKw == nil || n.Name == nil {
		return n
	}

	res := n.Name.EqualsNode(n2.Name)
	if res != nil {
		return res
	}

	if n.CallKw.Equals(n2.CallKw) {
		return nil
	}
	return n
}

// EqualsNode returns true if the two nodes are equal.
func (n *CallOrIndexExpr) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*CallOrIndexExpr)
	if !ok {
		return n
	}

	if len(n.Args) != len(n2.Args) {
		return n
	}
	for i, v := range n.Args {
		res := v.EqualsNode(n2.Args[i])
		if res != nil {
			return res
		}
	}
	if n.Name == nil {
		return n
	}

	return n.Name.EqualsNode(n2.Name)
}

// EqualsNode returns true if the two nodes are equal.
func (n *CallSelectorExpr) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*CallSelectorExpr)
	if !ok {
		return n
	}

	if n.Selector == nil {
		return n
	}
	return n.Selector.EqualsNode(n2.Selector)
}

// EqualsNode returns true if the two nodes are equal.
func (n *EmptyStmt) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*EmptyStmt)
	if !ok {
		return n
	}
	if n2.EOL == nil {
		return n
	}
	if n.EOL.Equals(n2.EOL) {
		return nil
	}
	return n
}

// EqualsNode returns true if the two nodes are equal.
func (n *ExprStmt) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*ExprStmt)
	if !ok {
		return n
	}
	if n.Expression == nil {
		return n
	}
	return n.Expression.EqualsNode(n2.Expression)
}

// EqualsNode returns true if the two nodes are equal.
func (n *File) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*File)
	if !ok {
		return n
	}

	if len(n.StatementLists) != len(n2.StatementLists) {
		return n
	}
	for i, v := range n.StatementLists {
		res := v.EqualsNode(&n2.StatementLists[i])
		if res != nil {
			return res
		}
	}
	// file name may be different
	return nil
}

// EqualsNode returns true if the two nodes are equal.
func (n *StatementList) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*StatementList)
	if !ok {
		return n
	}

	if len(n.Statements) != len(n2.Statements) {
		return n
	}
	for i, v := range n.Statements {
		res := v.EqualsNode(n2.Statements[i])
		if res != nil {
			return res
		}
	}
	return nil
}

// EqualsNode returns true if the two nodes are equal.
func (n *FuncDecl) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*FuncDecl)
	if !ok {
		return n
	}

	if len(n.Body) != len(n2.Body) {
		return n
	}
	for i, v := range n.Body {
		res := v.EqualsNode(&n2.Body[i])
		if res != nil {
			return res
		}
	}
	if n.FuncName != nil {
		res := n.FuncName.EqualsNode(n2.FuncName)
		if res != nil {
			return res
		}
	} else {
		if n2.FuncName != nil {
			return n
		}

	}
	if n.FuncType != nil {
		res := n.FuncType.EqualsNode(n2.FuncType)
		if res != nil {
			return res
		}
	} else {
		if n2.FuncType != nil {
			return n
		}
	}
	if n.FunctionKw == nil {
		return n
	}
	if n.FunctionKw.Equals(n2.FunctionKw) {
		return nil
	}
	return n
}

// EqualsNode returns true if the two nodes are equal.
func (n *SubDecl) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*SubDecl)
	if !ok {
		return n
	}

	if len(n.Body) != len(n2.Body) {
		return n
	}
	for i, v := range n.Body {
		res := v.EqualsNode(&n2.Body[i])
		if res != nil {
			return res
		}
	}
	if n.SubName != nil {
		res := n.SubName.EqualsNode(n2.SubName)
		if res != nil {
			return res
		}
	} else {
		if n2.SubName != nil {
			return n
		}
	}
	if n.SubType != nil {
		res := n.SubType.EqualsNode(n2.SubType)
		if res != nil {
			return res
		}
	} else {
		if n2.SubType != nil {
			return n
		}
	}

	if n.SubKw.Equals(n2.SubKw) {
		return nil
	}
	return n
}

// EqualsNode returns true if the two nodes are equal.
func (n *FuncType) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*FuncType)
	if !ok {
		return n
	}

	if len(n.Params) != len(n2.Params) {
		return n
	}
	for i, v := range n.Params {
		res := v.EqualsNode(&n2.Params[i])
		if res != nil {
			return res
		}
	}
	if !(n.Lparen.Equals(n2.Lparen) && n.Rparen.Equals(&n2.Rparen)) {
		return n
	}
	return n.Result.EqualsNode(n2.Result)
}

// EqualsNode returns true if the two nodes are equal.
func (n *SubType) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*SubType)
	if !ok {
		return n
	}

	if len(n.Params) != len(n2.Params) {
		return n
	}
	for i, v := range n.Params {
		res := v.EqualsNode(&n2.Params[i])
		if res != nil {
			return res
		}
	}
	if n.Lparen.Equals(n2.Lparen) && n.Rparen.Equals(&n2.Rparen) {
		return nil
	}
	return n
}

// EqualsNode returns true if the two nodes are equal.
func (n *Identifier) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*Identifier)
	if !ok {
		return n
	}
	if n.Name == n2.Name && n.Tok.Equals(n2.Tok) {
		return nil
	}
	return n
}

// EqualsNode returns true if the two nodes are equal.
func (n *IfStmt) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*IfStmt)
	if !ok {
		return n
	}

	// compare body
	if len(n.Body) != len(n2.Body) {
		return n
	}
	for i, v := range n.Body {
		res := v.EqualsNode(&n2.Body[i])
		if res != nil {
			return res
		}
	}
	// compare else
	if len(n.Else) != len(n2.Else) {
		return n
	}
	for i, v := range n.Else {
		res := v.EqualsNode(&n2.Else[i])
		if res != nil {
			return res
		}
	}
	if !n.IfKw.Equals(n2.IfKw) {
		return n
	}
	return n.Condition.EqualsNode(n2.Condition)
}

// EqualsNode returns true if the two nodes are equal.
func (n *SelectStmt) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}

	n2, ok := node.(*SelectStmt)
	if !ok {
		return n
	}

	if len(n.Body) != len(n2.Body) {
		return n
	}
	for i, v := range n.Body {
		res := v.EqualsNode(&n2.Body[i])
		if res != nil {
			return res
		}
	}
	if !n.SelectKw.Equals(n2.SelectKw) {
		return n
	}
	return n.Condition.EqualsNode(n2.Condition)
}

// EqualsNode returns true if the two nodes are equal.
func (n *CaseExpr) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}

	n2, ok := node.(*CaseExpr)
	if !ok {
		return n
	}

	if len(n.Body) != len(n2.Body) {
		return n
	}
	for i, v := range n.Body {
		res := v.EqualsNode(&n2.Body[i])
		if res != nil {
			return res
		}
	}

	if n.Condition != nil {
		res := n.Condition.EqualsNode(n2.Condition)
		if res != nil {
			return res
		}
	} else {
		if n2.Condition != nil {
			return n
		}
	}
	if n.CaseKw.Equals(n2.CaseKw) {
		return nil
	}
	return n
}

// EqualsNode returns true if the two nodes are equal.
func (n *ElseIfStmt) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*ElseIfStmt)
	if !ok {
		return n
	}

	// compare body
	if len(n.Body) != len(n2.Body) {
		return n
	}
	for i, v := range n.Body {
		res := v.EqualsNode(&n2.Body[i])
		if res != nil {
			return res
		}
	}
	if !n.ElseIfKw.Equals(n2.ElseIfKw) {
		return n
	}
	return n.Condition.EqualsNode(n2.Condition)
}

// EqualsNode returns true if the two nodes are equal.
func (n *ParenExpr) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*ParenExpr)
	if !ok {
		return n
	}

	if !(n.Lparen.Equals(n2.Lparen) && n.Rparen.Equals(&n2.Rparen)) {
		return n
	}
	return n.Expr.EqualsNode(n2.Expr)
}

// EqualsNode returns true if the two nodes are equal.
func (n *TypeDef) EqualsNode(node Node) Node {
	panic("Not implemented")

}

// EqualsNode returns true if the two nodes are equal.
func (n *UnaryExpr) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*UnaryExpr)
	if !ok {
		return n
	}

	if !(n.OpKind == n2.OpKind && n.OpToken.Equals(n2.OpToken)) {
		return n
	}
	return n.Right.EqualsNode(n2.Right)
}

// EqualsNode returns true if the two nodes are equal.
func (n *DimDecl) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*DimDecl)
	if !ok {
		return n
	}

	if len(n.Vars) != len(n2.Vars) {
		return n
	}
	for i, v := range n.Vars {
		res := v.EqualsNode(n2.Vars[i])
		if res != nil {
			return res
		}
	}

	if n.DimKw.Equals(n2.DimKw) {
		return nil
	}
	return n
}

// EqualsNode returns true if the two nodes are equal.
func (n *ConstDeclItem) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*ConstDeclItem)
	if !ok {
		return n
	}

	if n.ConstName != nil {
		res := n.ConstName.EqualsNode(n2.ConstName)
		if res != nil {
			return res
		}
	} else {
		if n2.ConstName != nil {
			return n
		}
	}
	if n.ConstValue != nil {
		res := n.ConstValue.EqualsNode(n2.ConstValue)
		if res != nil {
			return res
		}
	} else {
		if n2.ConstValue != nil {
			return n
		}
	}
	if n.ConstType != nil {
		res := n.ConstType.EqualsNode(n2.ConstType)
		if res != nil {
			return res
		}
	} else {
		if n2.ConstType != nil {
			return n
		}
	}
	return nil
}

// EqualsNode returns true if the two nodes are equal.
func (n *ArrayDecl) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*ArrayDecl)
	if !ok {
		return n
	}
	if n.VarName != nil {
		res := n.VarName.EqualsNode(n2.VarName)
		if res != nil {
			return res
		}
	} else {
		if n2.VarName != nil {
			return n
		}
	}
	if n.VarType != nil {
		return n.VarType.EqualsNode(n2.VarType)
	} else {
		if n2.VarType != nil {
			return n
		}
	}
	return nil
}

// EqualsNode returns true if the two nodes are equal.
func (n *ScalarDecl) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*ScalarDecl)
	if !ok {
		return n
	}

	if n.VarName != nil {
		res := n.VarName.EqualsNode(n2.VarName)
		if res != nil {
			return res
		}
	}
	if n.VarType != nil {
		res := n.VarType.EqualsNode(n2.VarType)
		if res != nil {
			return res
		}
	}
	if n.Value() != nil {
		return n.Value().EqualsNode(n2.Value())
	}
	return nil
}

// EqualsNode returns true if the two nodes are equal.
func (n *WhileStmt) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*WhileStmt)
	if !ok {
		return n
	}

	if len(n.Body) != len(n2.Body) {
		return n
	}
	for i, v := range n.Body {
		res := v.EqualsNode(&n2.Body[i])
		if res != nil {
			return res
		}
	}
	if !n.DoKw.Equals(n2.DoKw) {
		return n
	}
	if n.Condition != nil {
		return n.Condition.EqualsNode(n2.Condition)
	} else {
		if n2.Condition != nil {
			return n
		}
	}
	return nil
}

// EqualsNode returns true if the two nodes are equal.
func (n *UntilStmt) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*UntilStmt)
	if !ok {
		return n
	}

	if len(n.Body) != len(n2.Body) {
		return n
	}
	for i, v := range n.Body {
		res := v.EqualsNode(&n2.Body[i])
		if res != nil {
			return res
		}
	}
	if !n.DoKw.Equals(n2.DoKw) {
		return n
	}

	if n.Condition != nil {
		return n.Condition.EqualsNode(n2.Condition)
	} else {
		if n2.Condition != nil {
			return n
		}
	}
	return nil
}

// EqualsNode returns true if the two nodes are equal.
func (n *DoWhileStmt) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*DoWhileStmt)
	if !ok {
		return n
	}

	if len(n.Body) != len(n2.Body) {
		return n
	}
	for i, v := range n.Body {
		res := v.EqualsNode(&n2.Body[i])
		if res != nil {
			return res
		}
	}
	if !n.DoKw.Equals(n2.DoKw) {
		return n
	}

	if n.Condition != nil {
		return n.Condition.EqualsNode(n2.Condition)
	} else {
		if n2.Condition != nil {
			return n
		}
	}
	return nil
}

// EqualsNode returns true if the two nodes are equal.
func (n *DoUntilStmt) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*DoUntilStmt)
	if !ok {
		return n
	}

	if len(n.Body) != len(n2.Body) {
		return n
	}
	for i, v := range n.Body {
		res := v.EqualsNode(&n2.Body[i])
		if res != nil {
			return res
		}
	}
	if !n.DoKw.Equals(n2.DoKw) {
		return n
	}

	if n.Condition != nil {
		return n.Condition.EqualsNode(n2.Condition)
	} else {
		if n2.Condition != nil {
			return n
		}
	}
	return nil
}

// EqualsNode returns true if the two nodes are equal.
func (n *ForStmt) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*ForStmt)
	if !ok {
		return n
	}

	if len(n.Body) != len(n2.Body) {
		return n
	}
	for i, v := range n.Body {
		res := v.EqualsNode(&n2.Body[i])
		if res != nil {
			return res
		}
	}
	if !n.ForKw.Equals(n2.ForKw) {
		return n
	}
	if n.Next != nil {
		res := n.Next.EqualsNode(n2.Next)
		if res != nil {
			return res
		}
	} else {
		if n2.Next != nil {
			return n
		}
	}

	return n.ForExpr.EqualsNode(n2.ForExpr)
}

// EqualsNode returns true if the two nodes are equal.
func (n *EnumDecl) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*EnumDecl)
	if !ok {
		return n
	}

	if len(n.Values) != len(n2.Values) {
		return n
	}
	for i, v := range n.Values {
		res := v.EqualsNode(&n2.Values[i])
		if res != nil {
			return res
		}
	}
	if !n.EnumKw.Equals(n2.EnumKw) {
		return n
	}
	return n.EnumName.EqualsNode(n2.EnumName)
}

// EqualsNode returns true if the two nodes are equal.
func (n *LabelDecl) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}

	n2, ok := node.(*LabelDecl)
	if !ok {
		return n
	}

	return n.LabelName.EqualsNode(n2.LabelName)
}

// EqualsNode returns true if the two nodes are equal.
func (n *UserDefinedType) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*UserDefinedType)
	if !ok {
		return n
	}
	return n.Name.EqualsNode(n2.Name)
}

// EqualsNode returns true if the two nodes are equal.
func (n *ParamItem) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*ParamItem)
	if !ok {
		return n
	}
	if n.VarName != nil {
		res := n.VarName.EqualsNode(n2.VarName)
		if res != nil {
			return res
		}
	} else {
		if n2.VarName != nil {
			return n
		}
	}
	if n.VarType != nil {
		res := n.VarType.EqualsNode(n2.VarType)
		if res != nil {
			return res
		}
	} else {
		if n2.VarType != nil {
			return n
		}
	}
	if n.Value() != nil {
		res := n.Value().EqualsNode(n2.Value())
		if res != nil {
			return res
		}
	} else {
		if n2.Value() != nil {
			return n
		}
	}
	if !(n.ByVal == n2.ByVal && n.Optional == n2.Optional && n.ParamArray == n2.ParamArray && n.IsArray == n2.IsArray) {
		return n
	}
	return nil
}

// EqualsNode returns true if the two nodes are equal.
func (n *ConstDecl) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*ConstDecl)
	if !ok {
		return n
	}

	if len(n.Consts) != len(n2.Consts) {
		return n
	}
	for i, v := range n.Consts {
		c := n2.Consts[i]
		res := v.EqualsNode(&c)
		if res != nil {
			return res
		}
	}
	if n.ConstKw.Equals(n2.ConstKw) {
		return nil
	}
	return n
}

// EqualsNode returns true if the two nodes are equal.
func (n *ForEachExpr) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*ForEachExpr)
	if !ok {
		return n
	}

	res := n.Collection.EqualsNode(n2.Collection)
	if res != nil {
		return res
	}

	return n.Var.EqualsNode(n2.Var)
}

// EqualsNode returns true if the two nodes are equal.
func (n *ForNextExpr) EqualsNode(node Node) Node {
	if n == nil && node == nil {
		return nil
	}
	if n == nil || node == nil {
		return n
	}
	n2, ok := node.(*ForNextExpr)
	if !ok {
		return n
	}

	if n.Var != nil {
		res := n.Var.EqualsNode(n2.Var)
		if res != nil {
			return res
		}
	} else {
		if n2.Var != nil {
			return n
		}
	}
	if n.Step != nil {
		res := n.Step.EqualsNode(n2.Step)
		if res != nil {
			return res
		}
	} else {
		if n2.Step != nil {
			return n
		}
	}
	if n.To != nil {
		return n.To.EqualsNode(n2.To)
	} else {
		if n2.To != nil {
			return n
		}
	}
	return nil
}
