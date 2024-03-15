package sem

import (
	"fmt"
	"strconv"
	"strings"
	"uBasic/ast"
	"uBasic/ast/astutil"
	"uBasic/errors"
	"uBasic/token"
	"uBasic/types"
)

// UniversePos specifies a pseudo-position used for identifiers declared in the
// universe scope.
var UniversePos = token.Position{Line: 0, Column: 0, Absolute: -1}
var UniverseToken = token.Token{Position: UniversePos, Literal: ""}

// resolve performs identifier resolution, mapping identifiers to corresponding
// declarations.
func resolve(file *ast.File, scopes map[ast.Node]*Scope) error {
	// Pre-pass, add keyword types to universe scope.
	universe := NewScope(nil)

	// Add built-in types to universe scope.
	// add "$" suffix to avoid conflict between types and built-in functions with the same name (string, date, etc.)
	longIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Long$"}
	longDecl := &ast.TypeDef{DeclType: longIdent, TypeName: longIdent, Val: &types.Basic{Kind: types.Long}}
	longIdent.Decl = longDecl

	integerIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Integer$"}
	integerDecl := &ast.TypeDef{DeclType: integerIdent, TypeName: integerIdent, Val: &types.Basic{Kind: types.Integer}}
	integerIdent.Decl = integerDecl

	currency := &ast.Identifier{Tok: &UniverseToken, Name: "Currency$"}
	currencyDecl := &ast.TypeDef{DeclType: currency, TypeName: currency, Val: &types.Basic{Kind: types.Currency}}
	currency.Decl = currencyDecl

	doubleIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Double$"}
	doubleDecl := &ast.TypeDef{DeclType: doubleIdent, TypeName: doubleIdent, Val: &types.Basic{Kind: types.Double}}
	doubleIdent.Decl = doubleDecl

	singleIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Single$"}
	singleDecl := &ast.TypeDef{DeclType: singleIdent, TypeName: singleIdent, Val: &types.Basic{Kind: types.Single}}
	singleIdent.Decl = singleDecl

	stringIdent := &ast.Identifier{Tok: &UniverseToken, Name: "String$"}
	stringDecl := &ast.TypeDef{DeclType: stringIdent, TypeName: stringIdent, Val: &types.Basic{Kind: types.String}}
	stringIdent.Decl = stringDecl

	dateIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Date$"}
	dateDecl := &ast.TypeDef{DeclType: dateIdent, TypeName: dateIdent, Val: &types.Basic{Kind: types.Date}}
	dateIdent.Decl = dateDecl

	variantIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Variant$"}
	variantDecl := &ast.TypeDef{DeclType: variantIdent, TypeName: variantIdent, Val: &types.Basic{Kind: types.Variant}}
	variantIdent.Decl = variantDecl

	booleanIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Boolean$"}
	booleanDecl := &ast.TypeDef{DeclType: booleanIdent, TypeName: booleanIdent, Val: &types.Basic{Kind: types.Boolean}}
	booleanIdent.Decl = booleanDecl

	// Pre-pass, add run-time library functions to universe scope.
	// --------------------------------
	// ------- string functions -------
	// --------------------------------
	chrIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Chr"}
	chrParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "character"}, VarType: longIdent}
	chrFuncType := &ast.FuncType{Params: []ast.ParamItem{*chrParam1Item}, Result: stringIdent}
	chrDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: chrIdent, FuncType: chrFuncType, Body: nil}
	chrIdent.Decl = chrDecl

	inStrIdent := &ast.Identifier{Tok: &UniverseToken, Name: "InStr"}
	inStrParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "start"}, VarType: longIdent}
	inStrParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "string1"}, VarType: stringIdent}
	inStrParam3Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "string2"}, VarType: stringIdent}
	inStrParam4Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "compare"}, VarType: longIdent, Optional: true, DefaultValue: &ast.BasicLit{Kind: token.LongLit, ValTok: &UniverseToken, Value: "0"}}
	inStrFuncType := &ast.FuncType{Params: []ast.ParamItem{*inStrParam1Item, *inStrParam2Item, *inStrParam3Item, *inStrParam4Item}, Result: longIdent}
	inStrDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: inStrIdent, FuncType: inStrFuncType, Body: nil}
	inStrIdent.Decl = inStrDecl

	InStrRevIdent := &ast.Identifier{Tok: &UniverseToken, Name: "InStrRev"}
	InStrRevParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "stringCheck"}, VarType: stringIdent}
	InStrRevParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "stringMatch"}, VarType: stringIdent}
	InStrRevParam3Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "start"}, VarType: longIdent}
	InstrRevParam4Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "compare"}, VarType: longIdent, Optional: true, DefaultValue: &ast.BasicLit{Kind: token.LongLit, ValTok: &UniverseToken, Value: "0"}}
	InStrRevFuncType := &ast.FuncType{Params: []ast.ParamItem{*InStrRevParam1Item, *InStrRevParam2Item, *InStrRevParam3Item, *InstrRevParam4Item}, Result: longIdent}
	InStrRevDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: InStrRevIdent, FuncType: InStrRevFuncType, Body: nil}
	InStrRevIdent.Decl = InStrRevDecl

	LCaseIdent := &ast.Identifier{Tok: &UniverseToken, Name: "LCase"}
	LCaseParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "string"}, VarType: stringIdent}
	LCaseFuncType := &ast.FuncType{Params: []ast.ParamItem{*LCaseParam1Item}, Result: stringIdent}
	LCaseDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: LCaseIdent, FuncType: LCaseFuncType, Body: nil}
	LCaseIdent.Decl = LCaseDecl

	LeftIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Left"}
	LeftParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "string"}, VarType: stringIdent}
	LeftParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "length"}, VarType: longIdent}
	LeftFuncType := &ast.FuncType{Params: []ast.ParamItem{*LeftParam1Item, *LeftParam2Item}, Result: stringIdent}
	LeftDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: LeftIdent, FuncType: LeftFuncType, Body: nil}
	LeftIdent.Decl = LeftDecl

	LenIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Len"}
	LenParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "string"}, VarType: stringIdent}
	LenFuncType := &ast.FuncType{Params: []ast.ParamItem{*LenParam1Item}, Result: longIdent}
	LenDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: LenIdent, FuncType: LenFuncType, Body: nil}
	LenIdent.Decl = LenDecl

	LTrimIdent := &ast.Identifier{Tok: &UniverseToken, Name: "LTrim"}
	LTrimParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "string"}, VarType: stringIdent}
	LTrimFuncType := &ast.FuncType{Params: []ast.ParamItem{*LTrimParam1Item}, Result: stringIdent}
	LTrimDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: LTrimIdent, FuncType: LTrimFuncType, Body: nil}
	LTrimIdent.Decl = LTrimDecl

	MidIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Mid"}
	MidParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "string"}, VarType: stringIdent}
	MidParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "start"}, VarType: longIdent}
	MidParam3Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "length"}, VarType: longIdent, Optional: true, DefaultValue: &ast.BasicLit{Kind: token.LongLit, ValTok: &UniverseToken, Value: "-1"}}
	MidFuncType := &ast.FuncType{Params: []ast.ParamItem{*MidParam1Item, *MidParam2Item, *MidParam3Item}, Result: stringIdent}
	MidDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: MidIdent, FuncType: MidFuncType, Body: nil}
	MidIdent.Decl = MidDecl

	RightIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Right"}
	RightParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "string"}, VarType: stringIdent}
	RightParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "length"}, VarType: longIdent}
	RightFuncType := &ast.FuncType{Params: []ast.ParamItem{*RightParam1Item, *RightParam2Item}, Result: stringIdent}
	RightDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: RightIdent, FuncType: RightFuncType, Body: nil}
	RightIdent.Decl = RightDecl

	RTrimIdent := &ast.Identifier{Tok: &UniverseToken, Name: "RTrim"}
	RTrimParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "string"}, VarType: stringIdent}
	RTrimFuncType := &ast.FuncType{Params: []ast.ParamItem{*RTrimParam1Item}, Result: stringIdent}
	RTrimDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: RTrimIdent, FuncType: RTrimFuncType, Body: nil}
	RTrimIdent.Decl = RTrimDecl

	SpaceIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Space"}
	SpaceParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "length"}, VarType: longIdent}
	SpaceFuncType := &ast.FuncType{Params: []ast.ParamItem{*SpaceParam1Item}, Result: stringIdent}
	SpaceDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: SpaceIdent, FuncType: SpaceFuncType, Body: nil}
	SpaceIdent.Decl = SpaceDecl

	StrCompIdent := &ast.Identifier{Tok: &UniverseToken, Name: "StrComp"}
	StrCompParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "string1"}, VarType: stringIdent}
	StrCompParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "string2"}, VarType: stringIdent}
	strCompParam3Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "compare"}, VarType: longIdent, Optional: true, DefaultValue: &ast.BasicLit{Kind: token.LongLit, ValTok: &UniverseToken, Value: "0"}}
	StrCompFuncType := &ast.FuncType{Params: []ast.ParamItem{*StrCompParam1Item, *StrCompParam2Item, *strCompParam3Item}, Result: longIdent}
	StrCompDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: StrCompIdent, FuncType: StrCompFuncType, Body: nil}
	StrCompIdent.Decl = StrCompDecl

	StrngIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Strng"}
	StrngParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "length"}, VarType: longIdent}
	StrngParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "character"}, VarType: stringIdent}
	StrngFuncType := &ast.FuncType{Params: []ast.ParamItem{*StrngParam1Item, *StrngParam2Item}, Result: stringIdent}
	StrngDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: StrngIdent, FuncType: StrngFuncType, Body: nil}
	StrngIdent.Decl = StrngDecl

	StrReverseIdent := &ast.Identifier{Tok: &UniverseToken, Name: "StrReverse"}
	StrReverseParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "string"}, VarType: stringIdent}
	StrReverseFuncType := &ast.FuncType{Params: []ast.ParamItem{*StrReverseParam1Item}, Result: stringIdent}
	StrReverseDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: StrReverseIdent, FuncType: StrReverseFuncType, Body: nil}
	StrReverseIdent.Decl = StrReverseDecl

	TrimIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Trim"}
	TrimParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "string"}, VarType: stringIdent}
	TrimFuncType := &ast.FuncType{Params: []ast.ParamItem{*TrimParam1Item}, Result: stringIdent}
	TrimDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: TrimIdent, FuncType: TrimFuncType, Body: nil}
	TrimIdent.Decl = TrimDecl

	UCaseIdent := &ast.Identifier{Tok: &UniverseToken, Name: "UCase"}
	UCaseParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "string"}, VarType: stringIdent}
	UCaseFuncType := &ast.FuncType{Params: []ast.ParamItem{*UCaseParam1Item}, Result: stringIdent}
	UCaseDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: UCaseIdent, FuncType: UCaseFuncType, Body: nil}
	UCaseIdent.Decl = UCaseDecl

	// --------------------------------
	// ------- date/time functions -------
	// --------------------------------
	DteFunctIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Dte"}
	DteFuncType := &ast.FuncType{Params: nil, Result: dateIdent}
	DteFunctDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: DteFunctIdent, FuncType: DteFuncType, Body: nil}
	DteFunctIdent.Decl = DteFunctDecl

	DateAddIdent := &ast.Identifier{Tok: &UniverseToken, Name: "DateAdd"}
	DateAddParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "interval"}, VarType: stringIdent}
	DateAddParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "number"}, VarType: longIdent}
	DateAddParam3Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "date"}, VarType: variantIdent}
	DateAddFuncType := &ast.FuncType{Params: []ast.ParamItem{*DateAddParam1Item, *DateAddParam2Item, *DateAddParam3Item}, Result: dateIdent}
	DateAddDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: DateAddIdent, FuncType: DateAddFuncType, Body: nil}
	DateAddIdent.Decl = DateAddDecl

	DateDiffIdent := &ast.Identifier{Tok: &UniverseToken, Name: "DateDiff"}
	DateDiffParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "interval"}, VarType: stringIdent}
	DateDiffParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "date1"}, VarType: dateIdent}
	DateDiffParam3Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "date2"}, VarType: dateIdent}
	DateDiffFuncType := &ast.FuncType{Params: []ast.ParamItem{*DateDiffParam1Item, *DateDiffParam2Item, *DateDiffParam3Item}, Result: longIdent}
	DateDiffDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: DateDiffIdent, FuncType: DateDiffFuncType, Body: nil}
	DateDiffIdent.Decl = DateDiffDecl

	DatePartIdent := &ast.Identifier{Tok: &UniverseToken, Name: "DatePart"}
	DatePartParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "interval"}, VarType: stringIdent}
	DatePartParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "date"}, VarType: dateIdent}
	DatePartFuncType := &ast.FuncType{Params: []ast.ParamItem{*DatePartParam1Item, *DatePartParam2Item}, Result: longIdent}
	DatePartDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: DatePartIdent, FuncType: DatePartFuncType, Body: nil}
	DatePartIdent.Decl = DatePartDecl

	DateSerialIdent := &ast.Identifier{Tok: &UniverseToken, Name: "DateSerial"}
	DateSerialParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "year"}, VarType: longIdent}
	DateSerialParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "month"}, VarType: longIdent}
	DateSerialParam3Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "day"}, VarType: longIdent}
	DateSerialFuncType := &ast.FuncType{Params: []ast.ParamItem{*DateSerialParam1Item, *DateSerialParam2Item, *DateSerialParam3Item}, Result: dateIdent}
	DateSerialDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: DateSerialIdent, FuncType: DateSerialFuncType, Body: nil}
	DateSerialIdent.Decl = DateSerialDecl

	DateValueIdent := &ast.Identifier{Tok: &UniverseToken, Name: "DateValue"}
	DateValueParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "date"}, VarType: stringIdent}
	DateValueFuncType := &ast.FuncType{Params: []ast.ParamItem{*DateValueParam1Item}, Result: dateIdent}
	DateValueDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: DateValueIdent, FuncType: DateValueFuncType, Body: nil}
	DateValueIdent.Decl = DateValueDecl

	DayIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Day"}
	DayParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "date"}, VarType: dateIdent}
	DayFuncType := &ast.FuncType{Params: []ast.ParamItem{*DayParam1Item}, Result: longIdent}
	DayDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: DayIdent, FuncType: DayFuncType, Body: nil}
	DayIdent.Decl = DayDecl

	HourIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Hour"}
	HourParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "time"}, VarType: dateIdent}
	HourFuncType := &ast.FuncType{Params: []ast.ParamItem{*HourParam1Item}, Result: longIdent}
	HourDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: HourIdent, FuncType: HourFuncType, Body: nil}
	HourIdent.Decl = HourDecl

	MinuteIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Minute"}
	MinuteParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "time"}, VarType: dateIdent}
	MinuteFuncType := &ast.FuncType{Params: []ast.ParamItem{*MinuteParam1Item}, Result: longIdent}
	MinuteDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: MinuteIdent, FuncType: MinuteFuncType, Body: nil}
	MinuteIdent.Decl = MinuteDecl

	MonthIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Month"}
	MonthParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "date"}, VarType: dateIdent}
	MonthFuncType := &ast.FuncType{Params: []ast.ParamItem{*MonthParam1Item}, Result: longIdent}
	MonthDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: MonthIdent, FuncType: MonthFuncType, Body: nil}
	MonthIdent.Decl = MonthDecl

	NowIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Now"}
	NowFuncType := &ast.FuncType{Params: nil, Result: dateIdent}
	NowDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: NowIdent, FuncType: NowFuncType, Body: nil}
	NowIdent.Decl = NowDecl

	SecondIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Second"}
	SecondParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "time"}, VarType: dateIdent}
	SecondFuncType := &ast.FuncType{Params: []ast.ParamItem{*SecondParam1Item}, Result: longIdent}
	SecondDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: SecondIdent, FuncType: SecondFuncType, Body: nil}
	SecondIdent.Decl = SecondDecl

	TimeIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Time"}
	TimeFuncType := &ast.FuncType{Params: nil, Result: dateIdent}
	TimeDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: TimeIdent, FuncType: TimeFuncType, Body: nil}
	TimeIdent.Decl = TimeDecl

	TimerIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Timer"}
	TimerFuncType := &ast.FuncType{Params: nil, Result: doubleIdent}
	TimerDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: TimerIdent, FuncType: TimerFuncType, Body: nil}
	TimerIdent.Decl = TimerDecl

	TimeSerialIdent := &ast.Identifier{Tok: &UniverseToken, Name: "TimeSerial"}
	TimeSerialParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "hour"}, VarType: longIdent}
	TimeSerialParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "minute"}, VarType: longIdent}
	TimeSerialParam3Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "second"}, VarType: longIdent}
	TimeSerialFuncType := &ast.FuncType{Params: []ast.ParamItem{*TimeSerialParam1Item, *TimeSerialParam2Item, *TimeSerialParam3Item}, Result: dateIdent}
	TimeSerialDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: TimeSerialIdent, FuncType: TimeSerialFuncType, Body: nil}
	TimeSerialIdent.Decl = TimeSerialDecl

	TimeValueIdent := &ast.Identifier{Tok: &UniverseToken, Name: "TimeValue"}
	TimeValueParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "time"}, VarType: stringIdent}
	TimeValueFuncType := &ast.FuncType{Params: []ast.ParamItem{*TimeValueParam1Item}, Result: dateIdent}
	TimeValueDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: TimeValueIdent, FuncType: TimeValueFuncType, Body: nil}
	TimeValueIdent.Decl = TimeValueDecl

	WeekdayIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Weekday"}
	WeekdayParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "date"}, VarType: dateIdent}
	WeekdayFuncType := &ast.FuncType{Params: []ast.ParamItem{*WeekdayParam1Item}, Result: longIdent}
	WeekdayDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: WeekdayIdent, FuncType: WeekdayFuncType, Body: nil}
	WeekdayIdent.Decl = WeekdayDecl

	YearIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Year"}
	YearParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "date"}, VarType: dateIdent}
	YearFuncType := &ast.FuncType{Params: []ast.ParamItem{*YearParam1Item}, Result: longIdent}
	YearDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: YearIdent, FuncType: YearFuncType, Body: nil}
	YearIdent.Decl = YearDecl

	// --------------------------------
	// ------- conversion functions -------
	// --------------------------------

	CBoolIdent := &ast.Identifier{Tok: &UniverseToken, Name: "CBool"}
	CBoolParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "expression"}, VarType: variantIdent}
	CBoolFuncType := &ast.FuncType{Params: []ast.ParamItem{*CBoolParam1Item}, Result: booleanIdent}
	CBoolDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: CBoolIdent, FuncType: CBoolFuncType, Body: nil}
	CBoolIdent.Decl = CBoolDecl

	CDateIdent := &ast.Identifier{Tok: &UniverseToken, Name: "CDate"}
	CDateParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "expression"}, VarType: variantIdent}
	CDateFuncType := &ast.FuncType{Params: []ast.ParamItem{*CDateParam1Item}, Result: dateIdent}
	CDateDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: CDateIdent, FuncType: CDateFuncType, Body: nil}
	CDateIdent.Decl = CDateDecl

	CDblIdent := &ast.Identifier{Tok: &UniverseToken, Name: "CDbl"}
	CDblParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "expression"}, VarType: variantIdent}
	CDblFuncType := &ast.FuncType{Params: []ast.ParamItem{*CDblParam1Item}, Result: doubleIdent}
	CDblDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: CDblIdent, FuncType: CDblFuncType, Body: nil}
	CDblIdent.Decl = CDblDecl

	CLngIdent := &ast.Identifier{Tok: &UniverseToken, Name: "CLng"}
	CLngParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "expression"}, VarType: variantIdent}
	CLngFuncType := &ast.FuncType{Params: []ast.ParamItem{*CLngParam1Item}, Result: longIdent}
	CLngDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: CLngIdent, FuncType: CLngFuncType, Body: nil}
	CLngIdent.Decl = CLngDecl

	CStrIdent := &ast.Identifier{Tok: &UniverseToken, Name: "CStr"}
	CStrParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "expression"}, VarType: variantIdent}
	CStrFuncType := &ast.FuncType{Params: []ast.ParamItem{*CStrParam1Item}, Result: stringIdent}
	CStrDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: CStrIdent, FuncType: CStrFuncType, Body: nil}
	CStrIdent.Decl = CStrDecl

	CvarIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Cvar"}
	CvarParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "expression"}, VarType: variantIdent}
	CvarFuncType := &ast.FuncType{Params: []ast.ParamItem{*CvarParam1Item}, Result: variantIdent}
	CvarDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: CvarIdent, FuncType: CvarFuncType, Body: nil}
	CvarIdent.Decl = CvarDecl

	AscIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Asc"}
	AscParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "character"}, VarType: stringIdent}
	AscFuncType := &ast.FuncType{Params: []ast.ParamItem{*AscParam1Item}, Result: longIdent}
	AscDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: AscIdent, FuncType: AscFuncType, Body: nil}
	AscIdent.Decl = AscDecl

	FormatIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Format"}
	FormatParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "expression"}, VarType: variantIdent}
	FormatParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "format"}, VarType: stringIdent}
	FormatFuncType := &ast.FuncType{Params: []ast.ParamItem{*FormatParam1Item, *FormatParam2Item}, Result: stringIdent}
	FormatDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: FormatIdent, FuncType: FormatFuncType, Body: nil}
	FormatIdent.Decl = FormatDecl

	HexIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Hex"}
	HexParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "number"}, VarType: longIdent}
	HexFuncType := &ast.FuncType{Params: []ast.ParamItem{*HexParam1Item}, Result: stringIdent}
	HexDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: HexIdent, FuncType: HexFuncType, Body: nil}
	HexIdent.Decl = HexDecl

	OctIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Oct"}
	OctParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "number"}, VarType: longIdent}
	OctFuncType := &ast.FuncType{Params: []ast.ParamItem{*OctParam1Item}, Result: stringIdent}
	OctDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: OctIdent, FuncType: OctFuncType, Body: nil}
	OctIdent.Decl = OctDecl

	// --------------------------------
	// ------- mathematical functions -
	// --------------------------------

	AbsIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Abs"}
	AbsParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "number"}, VarType: doubleIdent}
	AbsFuncType := &ast.FuncType{Params: []ast.ParamItem{*AbsParam1Item}, Result: doubleIdent}
	AbsDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: AbsIdent, FuncType: AbsFuncType, Body: nil}
	AbsIdent.Decl = AbsDecl

	AtnIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Atn"}
	AtnParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "number"}, VarType: doubleIdent}
	AtnFuncType := &ast.FuncType{Params: []ast.ParamItem{*AtnParam1Item}, Result: doubleIdent}
	AtnDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: AtnIdent, FuncType: AtnFuncType, Body: nil}
	AtnIdent.Decl = AtnDecl

	CosIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Cos"}
	CosParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "number"}, VarType: doubleIdent}
	CosFuncType := &ast.FuncType{Params: []ast.ParamItem{*CosParam1Item}, Result: doubleIdent}
	CosDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: CosIdent, FuncType: CosFuncType, Body: nil}
	CosIdent.Decl = CosDecl

	ExpnIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Expn"}
	ExpnParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "number"}, VarType: doubleIdent}
	ExpnFuncType := &ast.FuncType{Params: []ast.ParamItem{*ExpnParam1Item}, Result: doubleIdent}
	ExpnDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: ExpnIdent, FuncType: ExpnFuncType, Body: nil}
	ExpnIdent.Decl = ExpnDecl

	FixIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Fix"}
	FixParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "number"}, VarType: doubleIdent}
	FixFuncType := &ast.FuncType{Params: []ast.ParamItem{*FixParam1Item}, Result: longIdent}
	FixDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: FixIdent, FuncType: FixFuncType, Body: nil}
	FixIdent.Decl = FixDecl

	intIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Int"}
	IntParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "number"}, VarType: doubleIdent}
	IntFuncType := &ast.FuncType{Params: []ast.ParamItem{*IntParam1Item}, Result: variantIdent}
	intDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: intIdent, FuncType: IntFuncType, Body: nil}
	integerIdent.Decl = intDecl

	LogIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Log"}
	LogParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "number"}, VarType: doubleIdent}
	LogFuncType := &ast.FuncType{Params: []ast.ParamItem{*LogParam1Item}, Result: doubleIdent}
	LogDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: LogIdent, FuncType: LogFuncType, Body: nil}
	LogIdent.Decl = LogDecl

	RndIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Rnd"}
	RndParam1Item := &ast.ParamItem{Optional: true, VarName: &ast.Identifier{Tok: &UniverseToken, Name: "number"}, VarType: doubleIdent, DefaultValue: &ast.BasicLit{Kind: token.DoubleLit, ValTok: &UniverseToken, Value: "0.0"}}
	RndFuncType := &ast.FuncType{Params: []ast.ParamItem{*RndParam1Item}, Result: doubleIdent}
	RndDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: RndIdent, FuncType: RndFuncType, Body: nil}
	RndIdent.Decl = RndDecl

	SgnIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Sgn"}
	SgnParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "number"}, VarType: doubleIdent}
	SgnFuncType := &ast.FuncType{Params: []ast.ParamItem{*SgnParam1Item}, Result: longIdent}
	SgnDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: SgnIdent, FuncType: SgnFuncType, Body: nil}
	SgnIdent.Decl = SgnDecl

	SinIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Sin"}
	SinParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "number"}, VarType: doubleIdent}
	SinFuncType := &ast.FuncType{Params: []ast.ParamItem{*SinParam1Item}, Result: doubleIdent}
	SinDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: SinIdent, FuncType: SinFuncType, Body: nil}
	SinIdent.Decl = SinDecl

	SqrIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Sqr"}
	SqrParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "number"}, VarType: doubleIdent}
	SqrFuncType := &ast.FuncType{Params: []ast.ParamItem{*SqrParam1Item}, Result: doubleIdent}
	SqrDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: SqrIdent, FuncType: SqrFuncType, Body: nil}
	SqrIdent.Decl = SqrDecl

	TanIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Tan"}
	TanParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "number"}, VarType: doubleIdent}
	TanFuncType := &ast.FuncType{Params: []ast.ParamItem{*TanParam1Item}, Result: doubleIdent}
	TanDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: TanIdent, FuncType: TanFuncType, Body: nil}
	TanIdent.Decl = TanDecl

	// --------------------------------
	// ------- array functions -------
	// --------------------------------

	LBoundIdent := &ast.Identifier{Tok: &UniverseToken, Name: "LBound"}
	LBoundParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "array"}, VarType: variantIdent}
	LBoundParam2Item := &ast.ParamItem{Optional: true, VarName: &ast.Identifier{Tok: &UniverseToken, Name: "dimension"}, VarType: longIdent, DefaultValue: &ast.BasicLit{Kind: token.LongLit, ValTok: &UniverseToken, Value: "1"}}
	LBoundFuncType := &ast.FuncType{Params: []ast.ParamItem{*LBoundParam1Item, *LBoundParam2Item}, Result: longIdent}
	LBoundDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: LBoundIdent, FuncType: LBoundFuncType, Body: nil}
	LBoundIdent.Decl = LBoundDecl

	UBoundIdent := &ast.Identifier{Tok: &UniverseToken, Name: "UBound"}
	UBoundParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "array"}, VarType: variantIdent}
	UBoundParam2Item := &ast.ParamItem{Optional: true, VarName: &ast.Identifier{Tok: &UniverseToken, Name: "dimension"}, VarType: longIdent, DefaultValue: &ast.BasicLit{Kind: token.LongLit, ValTok: &UniverseToken, Value: "1"}}
	UBoundFuncType := &ast.FuncType{Params: []ast.ParamItem{*UBoundParam1Item, *UBoundParam2Item}, Result: longIdent}
	UBoundDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: UBoundIdent, FuncType: UBoundFuncType, Body: nil}
	UBoundIdent.Decl = UBoundDecl

	// --------------------------------
	// ------- boolean constant -------
	// --------------------------------

	TrueIdent := &ast.Identifier{Tok: &UniverseToken, Name: "True"}
	TrueDecl := &ast.ConstDeclItem{ConstName: TrueIdent, ConstType: booleanIdent, ConstValue: &ast.BasicLit{Kind: token.BooleanLit, Value: "True"}}
	TrueIdent.Decl = TrueDecl

	FalseIdent := &ast.Identifier{Tok: &UniverseToken, Name: "False"}
	FalseDecl := &ast.ConstDeclItem{ConstName: FalseIdent, ConstType: booleanIdent, ConstValue: &ast.BasicLit{Kind: token.BooleanLit, Value: "False"}}
	FalseIdent.Decl = FalseDecl

	// --------------------------------
	// ------- input output -----------
	// --------------------------------

	InputIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Input"}
	InputParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &UniverseToken, Name: "prompt"}, VarType: stringIdent}
	InputFuncType := &ast.FuncType{Params: []ast.ParamItem{*InputParam1Item}, Result: stringIdent}
	InputDecl := &ast.FuncDecl{FunctionKw: &UniverseToken, FuncName: InputIdent, FuncType: InputFuncType, Body: nil}
	InputIdent.Decl = InputDecl

	// --------------------------------
	// ------- Internal classes -------
	// --------------------------------

	// Application class
	//    - Application class is a global class that is used to access the application object
	// Supported properties:
	//    - Name: string
	//    - Version: string
	//    - User: string
	// Supported methods:
	//    - getOS() as string

	NameIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Name"}
	NameDecl := &ast.ScalarDecl{VarName: NameIdent, VarType: stringIdent}
	VersionIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Version"}
	VersionDecl := &ast.ScalarDecl{VarName: VersionIdent, VarType: stringIdent}
	UserIdent := &ast.Identifier{Tok: &UniverseToken, Name: "User"}
	UserDecl := &ast.ScalarDecl{VarName: UserIdent, VarType: stringIdent}
	GetOSDecl := &ast.FuncDecl{
		FunctionKw: &UniverseToken,
		FuncName:   &ast.Identifier{Tok: &UniverseToken, Name: "getOS"},
		FuncType:   &ast.FuncType{Params: nil, Result: stringIdent},
		Body:       nil,
	}

	ApplicationIdent := &ast.Identifier{Tok: &UniverseToken, Name: "Application"}
	ApplicationClass := &ast.ClassDecl{
		ClassKw:   &UniverseToken,
		ClassName: ApplicationIdent,
		Members:   map[string]ast.Decl{},
	}
	// key names must be in lowercase
	ApplicationClass.Members["name"] = NameDecl
	ApplicationClass.Members["version"] = VersionDecl
	ApplicationClass.Members["user"] = UserDecl
	ApplicationClass.Members["getos"] = GetOSDecl
	ApplicationIdent.Decl = ApplicationClass

	// --------------------------------
	// ------- Constants --------------
	// --------------------------------

	vbTextCompare := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbTextCompare"}, ConstType: longIdent, ConstValue: &ast.BasicLit{Kind: token.DoubleLit, Value: "1", ValTok: &UniverseToken}}
	vbBinaryCompare := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbBinaryCompare"}, ConstType: longIdent, ConstValue: &ast.BasicLit{Kind: token.DoubleLit, Value: "0", ValTok: &UniverseToken}}
	vbUseSystem := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbUseSystem"}, ConstType: longIdent, ConstValue: &ast.BasicLit{Kind: token.DoubleLit, Value: "-2", ValTok: &UniverseToken}}
	vbSunday := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbSunday"}, ConstType: longIdent, ConstValue: &ast.BasicLit{Kind: token.DoubleLit, Value: "1", ValTok: &UniverseToken}}
	vbMonday := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbMonday"}, ConstType: longIdent, ConstValue: &ast.BasicLit{Kind: token.DoubleLit, Value: "2", ValTok: &UniverseToken}}
	vbTuesday := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbTuesday"}, ConstType: longIdent, ConstValue: &ast.BasicLit{Kind: token.DoubleLit, Value: "3", ValTok: &UniverseToken}}
	vbWednesday := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbWednesday"}, ConstType: longIdent, ConstValue: &ast.BasicLit{Kind: token.DoubleLit, Value: "4", ValTok: &UniverseToken}}
	vbThursday := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbThursday"}, ConstType: longIdent, ConstValue: &ast.BasicLit{Kind: token.DoubleLit, Value: "5", ValTok: &UniverseToken}}
	vbFriday := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbFriday"}, ConstType: longIdent, ConstValue: &ast.BasicLit{Kind: token.DoubleLit, Value: "6", ValTok: &UniverseToken}}
	vbSaturday := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbSaturday"}, ConstType: longIdent, ConstValue: &ast.BasicLit{Kind: token.DoubleLit, Value: "7", ValTok: &UniverseToken}}
	vbFirstJan1 := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbFirstJan1"}, ConstType: longIdent, ConstValue: &ast.BasicLit{Kind: token.DoubleLit, Value: "1", ValTok: &UniverseToken}}
	vbFirstFourDays := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbFirstFourDays"}, ConstType: longIdent, ConstValue: &ast.BasicLit{Kind: token.DoubleLit, Value: "2", ValTok: &UniverseToken}}
	vbFirstFullWeek := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbFirstFullWeek"}, ConstType: longIdent, ConstValue: &ast.BasicLit{Kind: token.DoubleLit, Value: "3", ValTok: &UniverseToken}}
	vbGeneralDate := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbGeneralDate"}, ConstType: longIdent, ConstValue: &ast.BasicLit{Kind: token.DoubleLit, Value: "0", ValTok: &UniverseToken}}
	vbLongDate := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbLongDate"}, ConstType: longIdent, ConstValue: &ast.BasicLit{Kind: token.DoubleLit, Value: "1", ValTok: &UniverseToken}}
	vbShortDate := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbShortDate"}, ConstType: longIdent, ConstValue: &ast.BasicLit{Kind: token.DoubleLit, Value: "2", ValTok: &UniverseToken}}
	vbLongTime := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbLongTime"}, ConstType: longIdent, ConstValue: &ast.BasicLit{Kind: token.DoubleLit, Value: "3", ValTok: &UniverseToken}}
	vbShortTime := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbShortTime"}, ConstType: longIdent, ConstValue: &ast.BasicLit{Kind: token.DoubleLit, Value: "4", ValTok: &UniverseToken}}
	vbObjectError := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbObjectError"}, ConstType: longIdent, ConstValue: &ast.BasicLit{Kind: token.DoubleLit, Value: "-2147221504", ValTok: &UniverseToken}}
	vbDataObjectError := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbDataObjectError"}, ConstType: longIdent, ConstValue: &ast.BasicLit{Kind: token.DoubleLit, Value: "-2147221500", ValTok: &UniverseToken}}

	vbTab := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbTab"}, ConstType: stringIdent, ConstValue: &ast.BasicLit{Kind: token.StringLit, Value: "\"\t\"", ValTok: &UniverseToken}}
	vbLf := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbLf"}, ConstType: stringIdent, ConstValue: &ast.BasicLit{Kind: token.StringLit, Value: "\"\n\"", ValTok: &UniverseToken}}
	vbCr := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbCr"}, ConstType: stringIdent, ConstValue: &ast.BasicLit{Kind: token.StringLit, Value: "\"\r\"", ValTok: &UniverseToken}}
	vbCrLf := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbCrLf"}, ConstType: stringIdent, ConstValue: &ast.BasicLit{Kind: token.StringLit, Value: "\"\r\n\"", ValTok: &UniverseToken}}
	vbVerticalTab := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbVerticalTab"}, ConstType: stringIdent, ConstValue: &ast.BasicLit{Kind: token.StringLit, Value: "\"\v\"", ValTok: &UniverseToken}}
	vbFormFeed := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbFormFeed"}, ConstType: stringIdent, ConstValue: &ast.BasicLit{Kind: token.StringLit, Value: "\"\f\"", ValTok: &UniverseToken}}
	vbBack := &ast.ConstDeclItem{ConstName: &ast.Identifier{Tok: &UniverseToken, Name: "vbBack"}, ConstType: stringIdent, ConstValue: &ast.BasicLit{Kind: token.StringLit, Value: "\"\b\"", ValTok: &UniverseToken}}

	universeConstants := []ast.Decl{
		vbTab,
		vbLf,
		vbCr,
		vbCrLf,
		vbVerticalTab,
		vbFormFeed,
		vbBack,

		vbTextCompare,
		vbBinaryCompare,
		vbUseSystem,
		vbSunday,
		vbMonday,
		vbTuesday,
		vbWednesday,
		vbThursday,
		vbFriday,
		vbSaturday,
		vbFirstJan1,
		vbFirstFourDays,
		vbFirstFullWeek,
		vbGeneralDate,
		vbLongDate,
		vbShortDate,
		vbLongTime,
		vbShortTime,
		vbObjectError,
		vbDataObjectError,
	}

	universeDecls := []*ast.TypeDef{
		longDecl,
		integerDecl,
		currencyDecl,
		doubleDecl,
		singleDecl,
		stringDecl,
		dateDecl,
		variantDecl,
		booleanDecl,
	}

	universeRuntime := []ast.Decl{
		chrDecl,
		inStrDecl,
		InStrRevDecl,
		LCaseDecl,
		LeftDecl,
		LenDecl,
		LTrimDecl,
		MidDecl,
		RightDecl,
		RTrimDecl,
		SpaceDecl,
		StrCompDecl,
		StrngDecl,
		StrReverseDecl,
		TrimDecl,
		UCaseDecl,
		DteFunctDecl,
		DateAddDecl,
		DateDiffDecl,
		DatePartDecl,
		DateSerialDecl,
		DateValueDecl,
		DayDecl,
		HourDecl,
		MinuteDecl,
		MonthDecl,
		NowDecl,
		SecondDecl,
		TimeDecl,
		TimerDecl,
		TimeSerialDecl,
		TimeValueDecl,
		WeekdayDecl,
		YearDecl,
		CBoolDecl,
		CDateDecl,
		CDblDecl,
		CLngDecl,
		CStrDecl,
		CvarDecl,
		AscDecl,
		FormatDecl,
		HexDecl,
		OctDecl,
		AbsDecl,
		AtnDecl,
		CosDecl,
		ExpnDecl,
		FixDecl,
		integerDecl,
		LogDecl,
		RndDecl,
		SgnDecl,
		SinDecl,
		SqrDecl,
		TanDecl,
		LBoundDecl,
		UBoundDecl,
		InputDecl,
		intDecl,
		ApplicationClass,
	}
	for _, decl := range universeDecls {
		if err := universe.Insert(decl); err != nil {
			return err
		}
	}
	for _, decl := range universeRuntime {
		if err := universe.Insert(decl); err != nil {
			return err
		}
	}

	for _, decl := range universeConstants {
		if err := universe.Insert(decl); err != nil {
			return err
		}
	}

	// First pass, add global declarations to file scope.
	fileScope := NewScope(universe)
	scopes[file] = fileScope
	fileScope.IsDef = func(decl ast.Node) bool {
		// Consider variable declarations as tentative definitions; i.e. return
		// false, unless variable definition.
		if decl, ok := decl.(*ast.ScalarDecl); ok {
			return decl.Value() != nil
		}
		return false
	}

	// scope specifies the current lexical scope.
	scope := fileScope
	labelForwardDeclarations := make(map[string]*ast.Identifier)

	secondPass := func() error {
		// Second pass, resolve forward declarations in file context.
		for name, ident := range labelForwardDeclarations {
			// find label in scope
			if decl, ok := scope.Lookup(name, true); ok {
				// verify that identifier is a label
				if _, ok := decl.(*ast.JumpLabelDecl); !ok {
					return errors.Newf(ident.Token().Position, "undeclared label %q", ident.Name)
				}
				ident.Decl = decl
				// save to delete: https://stackoverflow.com/questions/23229975/is-it-safe-to-remove-selected-keys-from-map-within-a-range-loop
				delete(labelForwardDeclarations, name)
			} else {
				return errors.Newf(ident.Token().Position, "undefined label %q", name)
			}
		}
		return nil
	}

	// resolve performs identifier resolution, mapping identifiers to the
	// corresponding declarations of the closest lexical scope.
	resolve := func(n ast.Node) error {
		switch n := n.(type) {
		case ast.Decl:
			switch n := n.(type) {
			case ast.FuncOrSub:
				// if node is a function or sub declaration
				// then we need to resolve the function or sub declaration
				// verify that the declaration is not a keyword
				name := n.Name().Name
				if token.IsKeyword(name) {
					return errors.Newf(n.Token().Position, "cannot declare %v, it is a reserved keyword", name)
				}
				// verify that the declaration is not a redeclaration
				if decl, ok := scope.Lookup(n.Name().Name, true); ok {
					return errors.Newf(n.Token().Position, "redeclaration of %v", n.Name())
				} else {
					// verify that the declaration offset is the same
					if decl != nil && n != nil && decl.Token().Position != n.Token().Position {
						return errors.Newf(n.Token().Position, "redeclaration of %v", n)
					}
				}
			case ast.VarDecl, *ast.ConstDeclItem, *ast.ParamItem:
				// if node is a variable or constant declaration item
				// then we need to resolve the type of the item
				// verify that the declaration is not a keyword
				name := n.Name().Name
				if token.IsKeyword(name) {
					return errors.Newf(n.Token().Position, "cannot declare %v, it is a reserved keyword", name)
				}

				// verify that the declaration is not a redeclaration
				if decl, ok := scope.Lookup(n.Name().Name, true); ok {
					return errors.Newf(n.Token().Position, "redeclaration of %v", n.Name())
				} else {
					// verify that the declaration offset is the same
					if decl != nil && n != nil && decl.Token().Position != n.Token().Position {
						return errors.Newf(n.Token().Position, "redeclaration of %v", n)
					}
				}
			case *ast.JumpLabelDecl:
				// if node is a jump label declaration
				// then we need to resolve the label declaration
				// verify that the declaration is not a keyword
				name := n.Label.Name
				if token.IsKeyword(name) {
					return errors.Newf(n.Token().Position, "cannot declare %v, it is a reserved keyword", name)
				}
				// verify that the declaration is not a redeclaration
				if decl, ok := scope.Lookup(n.Label.Name, true); ok {
					return errors.Newf(n.Token().Position, "redeclaration of %v", n.Label)
				} else {
					// verify that the declaration offset is the same
					if decl != nil && n != nil && decl.Token().Position != n.Token().Position {
						return errors.Newf(n.Token().Position, "redeclaration of %v", n.Label)
					}
				}
			case *ast.EnumDecl:
				// verify that the declaration is not a keyword
				name := n.Name().Name
				if token.IsKeyword(name) {
					return errors.Newf(n.Token().Position, "cannot declare %v, it is a reserved keyword", name)
				}

				// verify that the declaration is not a redeclaration
				if _, ok := scope.Lookup(n.Name().Name, true); ok {
					return errors.Newf(n.Token().Position, "redeclaration of %v", n.Name())
				}

				// if node is an enum declaration
				// then we need to resolve the enum items
				enum := n
				constType := &ast.Identifier{Name: enum.Name().Name, Tok: enum.Token()}
				for index, value := range enum.Values {

					constDecl := &ast.ConstDeclItem{
						ConstName:  &ast.Identifier{Name: name + "." + value.Name, Tok: value.Token(), Decl: n},
						ConstType:  constType,
						ConstValue: &ast.BasicLit{Kind: token.LongLit, Value: strconv.Itoa(index), ValTok: value.Token()},
					}
					// verify that the declaration is not a keyword
					constName := constDecl.ConstName.Name
					if token.IsKeyword(constName) {
						return errors.Newf(n.Token().Position, "cannot declare %v, it is a reserved keyword", constName)
					}

					// verify that the declaration is not a redeclaration
					if decl, ok := scope.Lookup(constDecl.ConstName.Name, true); ok {
						return errors.Newf(n.Token().Position, "redeclaration of %v", constName)
					} else {
						// verify that the declaration offset is the same
						if decl != nil && decl.Token().Position != constDecl.Token().Position {
							return errors.Newf(n.Token().Position, "redeclaration of %v", constName)
						}
					}

					if err := scope.Insert(constDecl); err != nil {
						return err
					}
					enum.Values[index].Decl = constDecl
				}
			}

			// Insert declaration into the scope
			// if scope != fileScope {

			if err := scope.Insert(n); err != nil {
				return err
			}
			// Create nested scope for function definitions.
			if fn, ok := n.(ast.FuncOrSub); ok {
				scope = NewScope(scope)
				scopes[fn] = scope
				// for _, param := range fn.GetParams() {
				// 	// add parameters to the function scope
				// 	if err := scope.Insert(&param); err != nil {
				// 		return err
				// 	}
				// }
				// add the function to the scope, for return resolution
				if err := scope.Insert(n); err != nil {
					return err
				}

			}

		case *ast.Identifier:
			// if node is an expression containing an identifier
			// then we need to resolve the identifier

			if n.Decl == nil {
				// if it is not a label
				labelParent := n.GetParent()
				if _, ok := labelParent.(*ast.JumpStmt); ok {
					break // will be validated in second pass
				} else {

					// check if it is a basic type first
					// the $ sign had to be added to break conflict between Date basic type and Date function; string type and string function
					// as type and function names are in the same scope object
					decl, ok := scope.Lookup(n.Name+"$", false)
					if !ok {
						decl, ok = scope.Lookup(n.Name, false)
						if !ok {
							return errors.Newf(n.Token().Position, "undeclared identifier %q", n)
						}
					}
					// cannot use function name as a declaration
					if _, ok := decl.(*ast.FuncDecl); ok {
						if parent, ok := n.GetParent().(*ast.UserDefinedType); ok {
							return errors.Newf(n.Token().Position, "cannot use function name as a declaration in %q", parent.Identifier)
						}
					}
					n.Decl = decl
				}
			}
		case *ast.CallSelectorExpr:
			switch root := n.Root.(type) {
			case *ast.CallOrIndexExpr:
				fmt.Println("call or index expression " + root.Identifier.Name + "()." + n.Selector.String())
				panic("not implemented")
			case *ast.CallSelectorExpr:
				fmt.Println("call Selector " + root.Root.String() + "." + root.Selector.String() + "()." + n.Selector.String())
				panic("not implemented")
			case *ast.Identifier:
				// find root in scope
				decl, ok := scope.Lookup(root.Name, false)
				if !ok {
					return errors.Newf(root.Token().Position, "undeclared identifier %q", root)
				}
				// find selector in decl members
				classDecl, ok := decl.(*ast.ClassDecl)
				if !ok {
					// validate if the root is an enum
					enumDecl, ok := decl.(*ast.EnumDecl)
					if ok {
						// transform the selector into identifier
						selector, ok := n.Selector.(*ast.Identifier)
						if !ok {
							return errors.Newf(n.Selector.Token().Position, "undeclared identifier %q", n.Selector)
						}
						// find selector in decl members
						for _, value := range enumDecl.Values {
							if strings.EqualFold(value.Name, selector.Name) {
								selector.Decl = value.Decl
								return nil
							}
						}
						return errors.Newf(n.Selector.Token().Position, "undeclared identifier %q", n.Selector)
					} else {
						return errors.Newf(root.Token().Position, "undeclared identifier %q", root)
					}
				}
				root.Decl = classDecl

				// find selector in decl members
				switch selector := n.Selector.(type) {
				case *ast.Identifier:
					decl, ok = classDecl.Members[strings.ToLower(selector.Name)]
					if !ok {
						return errors.Newf(selector.Token().Position, "undeclared identifier %q", selector)
					}
					selector.Decl = decl
				case *ast.CallOrIndexExpr:
					// find selector in decl members
					decl, ok = classDecl.Members[strings.ToLower(selector.Identifier.Name)].(*ast.FuncDecl)
					if !ok {
						return errors.Newf(selector.Token().Position, "undeclared identifier %q", selector)
					}
					selector.Identifier.Decl = decl
				default:
					return errors.Newf(n.Selector.Token().Position, "undeclared identifier %q", n.Selector)
				}

			}
		case *ast.JumpStmt:
			// if node is a jump statement
			// then we need to resolve the label
			if n.Label != nil {
				// put the label in the forward declarations scope
				labelForwardDeclarations[n.Label.Name] = n.Label
			}
			return nil
		}
		return nil
	}

	// after reverts to the outer scope after traversing block statements.
	after := func(n ast.Node) error {
		if _, ok := n.(ast.FuncOrSub); ok {
			// Second pass, resolve forward declarations.
			if err := secondPass(); err != nil {
				return err
			}

			scope = scope.Outer
			// } else if fn, ok := n.(*ast.FuncDecl); ok && !astutil.IsDef(fn) {
			// 	scope = scope.Outer
		}
		return nil
	}

	// Walk the AST of the given file to resolve identifiers.
	if err := astutil.WalkBeforeAfter(file, resolve, after); err != nil {
		return err
	}

	if err := secondPass(); err != nil {
		return err
	}
	return nil
}
