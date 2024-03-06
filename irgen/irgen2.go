package irgen

import (
	"fmt"
	"strings"
	"time"

	"strconv"
	"uBasic/ast"
	"uBasic/token"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Context struct {
	Block           *ir.Block
	Module          *ir.Module
	Parent          *Context
	Vars            map[string]value.Value
	SubAndFuncsDecl bool
	Functs          map[string]*ir.Func
}

func NewContext(b *ir.Block) *Context {
	return &Context{
		Block:  b,
		Parent: nil,
		Vars:   make(map[string]value.Value),
		Functs: make(map[string]*ir.Func),
	}
}

func (c *Context) NewContext(b *ir.Block) *Context {
	ctx := NewContext(b)
	ctx.Parent = c
	ctx.Module = c.Module
	ctx.SubAndFuncsDecl = c.SubAndFuncsDecl
	return ctx
}

func (c Context) lookupVariable(name string) value.Value {
	if v, ok := c.Vars[name]; ok {
		return v
	} else if c.Parent != nil {
		return c.Parent.lookupVariable(name)
	} else {
		return nil // , fmt.Errorf("Variable not found: " + name)
	}
}

func (c *Context) lookupFunction(name string) *ir.Func {
	if v, ok := c.Functs[name]; ok {
		return v
	} else if c.Parent != nil {
		return c.Parent.lookupFunction(name)
	} else {
		return nil
	}
}

func Compile(file *ast.File) *ir.Module {
	ctx := NewContext(nil)
	ctx.Compile(file)
	return ctx.Module
}

func (ctx *Context) Compile(node ast.Node) {
	switch node := node.(type) {
	// case *ast.CaseStmt:
	// case *ast.SelectStmt:
	// case *ast.CallSubStmt:
	// case *ast.CallSelectorExpr:
	// case *ast.WhileStmt:
	// case *ast.ForStmt:
	// case *ast.ForNextExpr:
	// case *ast.ForEachExpr:
	// case *ast.UntilStmt:
	// case *ast.DoWhileStmt:
	// case *ast.DoUntilStmt:
	// case *ast.ExitStmt:
	// case *ast.ParenExpr:
	// case *ast.UnaryExpr:
	// case *ast.BinaryExpr:
	// case *ast.ExprStmt:
	// case *ast.ConstDecl:

	// case *ast.ConstDeclItem:
	// case *ast.ElseIfStmt:
	// case *ast.IfStmt:
	// case *ast.EnumDecl:
	// case *ast.CallOrIndexExpr:
	case *ast.Comment:
		// nothing to do
	case *ast.BasicLit:
		compileBasicLit(ctx, node)
	// case *ast.DimDecl:
	// case *ast.ScalarDecl:
	// case *ast.ArrayDecl:
	// case *ast.ArrayType:
	// case *ast.Identifier:
	// case *ast.FuncDecl:
	// case *ast.SubDecl:
	// case *ast.FuncType:
	// case *ast.SubType:
	case *ast.SpecialStmt:
		compileSpecialStmt(ctx, node)
	case *ast.File:
		compileFile(ctx, node)
	case *ast.StatementList:
		compileStatementList(ctx, node)
	// case *ast.EmptyStmt:
	// case *ast.JumpLabelDecl:
	// case *ast.UserDefinedType:
	// case *ast.ParamItem:
	// case *ast.ClassDecl:
	// case *ast.JumpStmt:
	default:
		panic(fmt.Sprintf("unknown node type: %T", node))
	}
}

func compileBasicLit(ctx *Context, node *ast.BasicLit) {
	panic("not implemented yet")
}

func compileConstant(node *ast.BasicLit) constant.Constant {
	switch node.Kind {
	case token.BooleanLit:
		// we have no boolean in LLVM IR
		if node.Value.(bool) {
			return constant.NewInt(types.I1, 1)
		} else {
			return constant.NewInt(types.I1, 0)
		}
	case token.LongLit:
		i, err := strconv.ParseInt(node.Value.(string), 10, 64)
		if err != nil {

			return nil // , fmt.Errorf("error parsing long literal: %s", err)
		}
		return constant.NewInt(types.I64, i)
	case token.DoubleLit, token.CurrencyLit:
		i, err := strconv.ParseFloat(node.Value.(string), 64)
		if err != nil {
			return nil //, fmt.Errorf("error parsing double literal: %s", err)
		}
		return constant.NewFloat(types.Float, i)
	case token.StringLit:
		return constant.NewCharArrayFromString(node.Value.(string))
	case token.DateLit:
		// structure:
		//	INT32 number of days since 0001-01-01
		//	INT32 number of seconds since midnight
		dateTime := StringToTime(node.Value.(string))
		if dateTime == nil {
			return nil // , fmt.Errorf("error parsing date literal: %s", node.Value)
		}
		// number of days since 0001-01-01
		days, second := convertTimeIntoDaysAndSeconds(*dateTime)
		daysConstant := constant.NewInt(types.I32, int64(days))
		secondsConstant := constant.NewInt(types.I32, int64(second))
		DateStruct := types.NewStruct(types.I32, types.I32)
		return constant.NewStruct(DateStruct, daysConstant, secondsConstant)
	case token.KwNothing:
		return nil
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

func compileFile(ctx *Context, node *ast.File) {
	ctx.Module = ir.NewModule()
	ctx.SubAndFuncsDecl = true
	compileExternals(ctx)

	// process function and sub declarations
	for _, statementList := range node.StatementLists {
		ctx.Compile(&statementList)
	}

	funcMain := ctx.Module.NewFunc("main", types.Void)
	ctx = ctx.NewContext(funcMain.NewBlock(""))
	ctx.SubAndFuncsDecl = false
	// process statements a
	for _, statementList := range node.StatementLists {
		ctx.Compile(&statementList)
	}
	ctx.Block.NewRet(constant.NewInt(types.I32, 0))
}

func compileExternals(ctx *Context) {
	puts := ctx.Module.NewFunc("puts", types.I32, ir.NewParam("", types.NewPointer(types.I8)))
	ctx.Functs["puts"] = puts

}

func compileStatementList(ctx *Context, node *ast.StatementList) {
	for _, statement := range node.Statements {
		if ctx.SubAndFuncsDecl {
			switch statement := statement.(type) {
			case *ast.FuncDecl, *ast.SubDecl:
				ctx.Compile(statement)
			}
		} else {
			switch statement := statement.(type) {
			case *ast.FuncDecl:
				// do nothing
			case *ast.SubDecl:
				// do nothing
			default:
				ctx.Compile(statement)
			}
		}
	}
}

func compileSpecialStmt(ctx *Context, node *ast.SpecialStmt) {
	switch strings.ToLower(node.Keyword1.Literal) {
	case "print", "debug.print", "msgbox":
		compilePrintStmt(ctx, node)
	default:
		panic("unknown special statement")
	}
}

func compilePrintStmt(ctx *Context, node *ast.SpecialStmt) {
	zero := constant.NewInt(types.I64, 0)
	puts := ctx.lookupFunction("puts")
	for _, arg := range node.Args {
		// prepare the string to be printed
		basicText := strings.Trim(arg.String(), "\"")
		basicText = strings.Replace(basicText, "\"\"", "\"", -1)
		basicText = strings.Replace(basicText, "\x00", "", -1)

		text := constant.NewCharArrayFromString(basicText + "\x00")
		variable := ctx.Module.NewGlobalDef(basicText, text)
		gep := constant.NewGetElementPtr(text.Typ, variable, zero, zero)
		ctx.Block.NewCall(puts, gep)
	}
}
