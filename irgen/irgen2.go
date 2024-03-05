package irgen

import (
	"fmt"
	"strings"
	"time"

	"strconv"
	"uBasic/ast"
	"uBasic/object"
	"uBasic/token"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

type Context struct {
	Block  *ir.Block
	Parent *Context
	Vars   map[string]object.Object // to hold the variable values
}

func NewContext(b *ir.Block) *Context {
	return &Context{
		Block:  b,
		Parent: nil,
		Vars:   make(map[string]object.Object),
	}
}

func (c *Context) NewContext(b *ir.Block) *Context {
	ctx := NewContext(b)
	ctx.Parent = c
	return ctx
}

func (c Context) lookupVariable(name string) (object.Object, error) {
	if v, ok := c.Vars[name]; ok {
		return v, nil
	} else if c.Parent != nil {
		return c.Parent.lookupVariable(name)
	} else {
		return nil, fmt.Errorf("Variable not found: " + name)
	}
}

func (ctx *Context) Compile(node ast.Node) {
	if ctx.Parent != nil {
		return
	}
	// parentCtx := ctx.Parent
	switch node := node.(type) {
	case *ast.CaseStmt:
	case *ast.SelectStmt:
	case *ast.CallSubStmt:
	case *ast.CallSelectorExpr:
	case *ast.WhileStmt:
	case *ast.ForStmt:
	case *ast.ForNextExpr:
	case *ast.ForEachExpr:
	case *ast.UntilStmt:
	case *ast.DoWhileStmt:
	case *ast.DoUntilStmt:
	case *ast.ExitStmt:
	case *ast.ParenExpr:
	case *ast.UnaryExpr:
	case *ast.BinaryExpr:
	case *ast.ExprStmt:
	case *ast.ConstDecl:

	case *ast.ConstDeclItem:
	case *ast.ElseIfStmt:
	case *ast.IfStmt:
	case *ast.EnumDecl:
	case *ast.CallOrIndexExpr:
	case *ast.Comment:
	case *ast.BasicLit:
		compileConstant(node)
	case *ast.DimDecl:
	case *ast.ScalarDecl:
	case *ast.ArrayDecl:
	case *ast.ArrayType:
	case *ast.Identifier:
	case *ast.FuncDecl:
	case *ast.SubDecl:
	case *ast.FuncType:
	case *ast.SubType:
	case *ast.SpecialStmt:
	case *ast.File:
	case *ast.StatementList:
	case *ast.EmptyStmt:
	case *ast.JumpLabelDecl:
	case *ast.UserDefinedType:
	case *ast.ParamItem:
	case *ast.ClassDecl:
	case *ast.JumpStmt:
	}
}

func compileConstant(node *ast.BasicLit) (constant.Constant, error) {
	switch node.Kind {
	case token.BooleanLit:
		// we have no boolean in LLVM IR
		if node.Value.(bool) {
			return constant.NewInt(types.I1, 1), nil
		} else {
			return constant.NewInt(types.I1, 0), nil
		}
	case token.LongLit:
		i, err := strconv.ParseInt(node.Value.(string), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Error parsing long literal: %s", err)
		}
		return constant.NewInt(types.I64, i), nil
	case token.DoubleLit, token.CurrencyLit:
		i, err := strconv.ParseFloat(node.Value.(string), 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing double literal: %s", err)
		}
		return constant.NewFloat(types.Float, i), nil
	case token.StringLit:
		return constant.NewCharArrayFromString(node.Value.(string)), nil
	case token.DateLit:
		// structure:
		//	INT32 number of days since 0001-01-01
		//	INT32 number of seconds since midnight
		dateTime := StringToTime(node.Value.(string))
		if dateTime == nil {
			return nil, fmt.Errorf("error parsing date literal: %s", node.Value)
		}
		// number of days since 0001-01-01
		days, second := convertTimeIntoDaysAndSeconds(*dateTime)
		daysConstant := constant.NewInt(types.I32, int64(days))
		secondsConstant := constant.NewInt(types.I32, int64(second))
		DateStruct := types.NewStruct(types.I32, types.I32)
		return constant.NewStruct(DateStruct, daysConstant, secondsConstant), nil
	case token.KwNothing:
		return nil, nil
	}
	panic("unknown expression")
}

func StringToTime(s string) *time.Time {
	s = strings.Trim(s, "#")
	// List of date formats
	dateFormats := []string{
		"2006-01-02",
		"2006/01/02",
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05",
		"15:04:05",
	}

	for _, format := range dateFormats {
		dateTimeValue, err := time.Parse(format, s)
		if err == nil {
			return &dateTimeValue
		}
	}
	return nil
}

func convertTimeIntoDaysAndSeconds(dateTime time.Time) (int, int) {
	format := "2006-01-02 15:04:05"

	// calculate the number of days since 0001-01-01
	then, _ := time.Parse(format, "0001-01-01 00:00:00")
	diff := dateTime.Sub(then)
	days := int(diff.Hours() / 24)
	// calculate the number of seconds since midnight today
	return days, dateTime.Hour()*3600 + dateTime.Minute()*60 + dateTime.Second()
}
