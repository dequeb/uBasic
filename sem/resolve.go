package sem

import (
	"strconv"
	"uBasic/ast"
	"uBasic/ast/astutil"
	"uBasic/errors"
	"uBasic/token"
	"uBasic/types"
)

// universePos specifies a pseudo-position used for identifiers declared in the
// universe scope.
var universePos = token.Position{Line: 0, Column: 0, Absolute: -1}
var universeToken = token.Token{Position: universePos, Literal: ""}

// resolve performs identifier resolution, mapping identifiers to corresponding
// declarations.
func resolve(file *ast.File, scopes map[ast.Node]*Scope) error {
	// Pre-pass, add keyword types to universe scope.
	universe := NewScope(nil)

	longIdent := &ast.Identifier{Tok: &universeToken, Name: "Long"}
	longDecl := &ast.TypeDef{DeclType: longIdent, TypeName: longIdent, Val: &types.Basic{Kind: types.Long}}
	longIdent.Decl = longDecl

	integerIdent := &ast.Identifier{Tok: &universeToken, Name: "Integer"}
	integerDecl := &ast.TypeDef{DeclType: integerIdent, TypeName: integerIdent, Val: &types.Basic{Kind: types.Integer}}
	integerIdent.Decl = integerDecl

	currency := &ast.Identifier{Tok: &universeToken, Name: "Currency"}
	currencyDecl := &ast.TypeDef{DeclType: currency, TypeName: currency, Val: &types.Basic{Kind: types.Integer}}
	currency.Decl = currencyDecl

	doubleIdent := &ast.Identifier{Tok: &universeToken, Name: "Double"}
	doubleDecl := &ast.TypeDef{DeclType: doubleIdent, TypeName: doubleIdent, Val: &types.Basic{Kind: types.Double}}
	doubleIdent.Decl = doubleDecl

	singleIdent := &ast.Identifier{Tok: &universeToken, Name: "Single"}
	singleDecl := &ast.TypeDef{DeclType: singleIdent, TypeName: singleIdent, Val: &types.Basic{Kind: types.Double}}
	singleIdent.Decl = singleDecl

	stringIdent := &ast.Identifier{Tok: &universeToken, Name: "String"}
	stringDecl := &ast.TypeDef{DeclType: stringIdent, TypeName: stringIdent, Val: &types.Basic{Kind: types.String}}
	stringIdent.Decl = stringDecl

	dateTimeIdent := &ast.Identifier{Tok: &universeToken, Name: "DateTime"}
	dateTimeDecl := &ast.TypeDef{DeclType: dateTimeIdent, TypeName: dateTimeIdent, Val: &types.Basic{Kind: types.DateTime}}
	dateTimeIdent.Decl = dateTimeDecl

	variantIdent := &ast.Identifier{Tok: &universeToken, Name: "Variant"}
	variantDecl := &ast.TypeDef{DeclType: variantIdent, TypeName: variantIdent, Val: &types.Basic{Kind: types.Variant}}
	variantIdent.Decl = variantDecl

	booleanIdent := &ast.Identifier{Tok: &universeToken, Name: "Boolean"}
	booleanDecl := &ast.TypeDef{DeclType: booleanIdent, TypeName: booleanIdent, Val: &types.Basic{Kind: types.Boolean}}
	booleanIdent.Decl = booleanDecl

	// Pre-pass, add run-time library functions to universe scope.
	// --------------------------------
	// ------- string functions -------
	// --------------------------------
	chrIdent := &ast.Identifier{Tok: &universeToken, Name: "Chr"}
	chrParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "character"}, VarType: longIdent}
	chrFuncType := &ast.FuncType{Params: []ast.ParamItem{*chrParam1Item}, Result: stringIdent}
	chrDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: chrIdent, FuncType: chrFuncType, Body: nil}
	chrIdent.Decl = chrDecl

	inStrIdent := &ast.Identifier{Tok: &universeToken, Name: "InStr"}
	inStrParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "start"}, VarType: longIdent}
	inStrParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "string1"}, VarType: stringIdent}
	inStrParam3Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "string2"}, VarType: stringIdent}
	inStrFuncType := &ast.FuncType{Params: []ast.ParamItem{*inStrParam1Item, *inStrParam2Item, *inStrParam3Item}, Result: longIdent}
	inStrDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: inStrIdent, FuncType: inStrFuncType, Body: nil}
	inStrIdent.Decl = inStrDecl

	InStrRevIdent := &ast.Identifier{Tok: &universeToken, Name: "InStrRev"}
	InStrRevParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "stringCheck"}, VarType: stringIdent}
	InStrRevParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "stringMatch"}, VarType: stringIdent}
	InStrRevParam3Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "start"}, VarType: longIdent}
	InStrRevFuncType := &ast.FuncType{Params: []ast.ParamItem{*InStrRevParam1Item, *InStrRevParam2Item, *InStrRevParam3Item}, Result: longIdent}
	InStrRevDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: InStrRevIdent, FuncType: InStrRevFuncType, Body: nil}
	InStrRevIdent.Decl = InStrRevDecl

	LCaseIdent := &ast.Identifier{Tok: &universeToken, Name: "LCase"}
	LCaseParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "string"}, VarType: stringIdent}
	LCaseFuncType := &ast.FuncType{Params: []ast.ParamItem{*LCaseParam1Item}, Result: stringIdent}
	LCaseDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: LCaseIdent, FuncType: LCaseFuncType, Body: nil}
	LCaseIdent.Decl = LCaseDecl

	LeftIdent := &ast.Identifier{Tok: &universeToken, Name: "Left"}
	LeftParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "string"}, VarType: stringIdent}
	LeftParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "length"}, VarType: longIdent}
	LeftFuncType := &ast.FuncType{Params: []ast.ParamItem{*LeftParam1Item, *LeftParam2Item}, Result: stringIdent}
	LeftDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: LeftIdent, FuncType: LeftFuncType, Body: nil}
	LeftIdent.Decl = LeftDecl

	LenIdent := &ast.Identifier{Tok: &universeToken, Name: "Len"}
	LenParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "string"}, VarType: stringIdent}
	LenFuncType := &ast.FuncType{Params: []ast.ParamItem{*LenParam1Item}, Result: longIdent}
	LenDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: LenIdent, FuncType: LenFuncType, Body: nil}
	LenIdent.Decl = LenDecl

	LTrimIdent := &ast.Identifier{Tok: &universeToken, Name: "LTrim"}
	LTrimParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "string"}, VarType: stringIdent}
	LTrimFuncType := &ast.FuncType{Params: []ast.ParamItem{*LTrimParam1Item}, Result: stringIdent}
	LTrimDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: LTrimIdent, FuncType: LTrimFuncType, Body: nil}
	LTrimIdent.Decl = LTrimDecl

	MidIdent := &ast.Identifier{Tok: &universeToken, Name: "Mid"}
	MidParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "string"}, VarType: stringIdent}
	MidParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "start"}, VarType: longIdent}
	MidParam3Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "length"}, VarType: longIdent}
	MidFuncType := &ast.FuncType{Params: []ast.ParamItem{*MidParam1Item, *MidParam2Item, *MidParam3Item}, Result: stringIdent}
	MidDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: MidIdent, FuncType: MidFuncType, Body: nil}
	MidIdent.Decl = MidDecl

	RightIdent := &ast.Identifier{Tok: &universeToken, Name: "Right"}
	RightParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "string"}, VarType: stringIdent}
	RightParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "length"}, VarType: longIdent}
	RightFuncType := &ast.FuncType{Params: []ast.ParamItem{*RightParam1Item, *RightParam2Item}, Result: stringIdent}
	RightDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: RightIdent, FuncType: RightFuncType, Body: nil}
	RightIdent.Decl = RightDecl

	RTrimIdent := &ast.Identifier{Tok: &universeToken, Name: "RTrim"}
	RTrimParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "string"}, VarType: stringIdent}
	RTrimFuncType := &ast.FuncType{Params: []ast.ParamItem{*RTrimParam1Item}, Result: stringIdent}
	RTrimDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: RTrimIdent, FuncType: RTrimFuncType, Body: nil}
	RTrimIdent.Decl = RTrimDecl

	SpaceIdent := &ast.Identifier{Tok: &universeToken, Name: "Space"}
	SpaceParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "length"}, VarType: longIdent}
	SpaceFuncType := &ast.FuncType{Params: []ast.ParamItem{*SpaceParam1Item}, Result: stringIdent}
	SpaceDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: SpaceIdent, FuncType: SpaceFuncType, Body: nil}
	SpaceIdent.Decl = SpaceDecl

	StrCompIdent := &ast.Identifier{Tok: &universeToken, Name: "StrComp"}
	StrCompParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "string1"}, VarType: stringIdent}
	StrCompParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "string2"}, VarType: stringIdent}
	StrCompFuncType := &ast.FuncType{Params: []ast.ParamItem{*StrCompParam1Item, *StrCompParam2Item}, Result: longIdent}
	StrCompDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: StrCompIdent, FuncType: StrCompFuncType, Body: nil}
	StrCompIdent.Decl = StrCompDecl

	StringIdent := &ast.Identifier{Tok: &universeToken, Name: "String"}
	StringParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "length"}, VarType: longIdent}
	StringParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "character"}, VarType: longIdent}
	StringFuncType := &ast.FuncType{Params: []ast.ParamItem{*StringParam1Item, *StringParam2Item}, Result: stringIdent}
	StringDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: StringIdent, FuncType: StringFuncType, Body: nil}
	StringIdent.Decl = StringDecl

	StrReverseIdent := &ast.Identifier{Tok: &universeToken, Name: "StrReverse"}
	StrReverseParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "string"}, VarType: stringIdent}
	StrReverseFuncType := &ast.FuncType{Params: []ast.ParamItem{*StrReverseParam1Item}, Result: stringIdent}
	StrReverseDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: StrReverseIdent, FuncType: StrReverseFuncType, Body: nil}
	StrReverseIdent.Decl = StrReverseDecl

	TrimIdent := &ast.Identifier{Tok: &universeToken, Name: "Trim"}
	TrimParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "string"}, VarType: stringIdent}
	TrimFuncType := &ast.FuncType{Params: []ast.ParamItem{*TrimParam1Item}, Result: stringIdent}
	TrimDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: TrimIdent, FuncType: TrimFuncType, Body: nil}
	TrimIdent.Decl = TrimDecl

	UCaseIdent := &ast.Identifier{Tok: &universeToken, Name: "UCase"}
	UCaseParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "string"}, VarType: stringIdent}
	UCaseFuncType := &ast.FuncType{Params: []ast.ParamItem{*UCaseParam1Item}, Result: stringIdent}
	UCaseDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: UCaseIdent, FuncType: UCaseFuncType, Body: nil}
	UCaseIdent.Decl = UCaseDecl

	// --------------------------------
	// ------- date/time functions -------
	// --------------------------------
	DateIdent := &ast.Identifier{Tok: &universeToken, Name: "Date"}
	DateFuncType := &ast.FuncType{Params: nil, Result: dateTimeIdent}
	DateDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: DateIdent, FuncType: DateFuncType, Body: nil}
	DateIdent.Decl = DateDecl

	DateAddIdent := &ast.Identifier{Tok: &universeToken, Name: "DateAdd"}
	DateAddParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "interval"}, VarType: stringIdent}
	DateAddParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "number"}, VarType: longIdent}
	DateAddParam3Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "date"}, VarType: dateTimeIdent}
	DateAddFuncType := &ast.FuncType{Params: []ast.ParamItem{*DateAddParam1Item, *DateAddParam2Item, *DateAddParam3Item}, Result: dateTimeIdent}
	DateAddDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: DateAddIdent, FuncType: DateAddFuncType, Body: nil}
	DateAddIdent.Decl = DateAddDecl

	DateDiffIdent := &ast.Identifier{Tok: &universeToken, Name: "DateDiff"}
	DateDiffParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "interval"}, VarType: stringIdent}
	DateDiffParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "date1"}, VarType: dateTimeIdent}
	DateDiffParam3Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "date2"}, VarType: dateTimeIdent}
	DateDiffFuncType := &ast.FuncType{Params: []ast.ParamItem{*DateDiffParam1Item, *DateDiffParam2Item, *DateDiffParam3Item}, Result: longIdent}
	DateDiffDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: DateDiffIdent, FuncType: DateDiffFuncType, Body: nil}
	DateDiffIdent.Decl = DateDiffDecl

	DatePartIdent := &ast.Identifier{Tok: &universeToken, Name: "DatePart"}
	DatePartParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "interval"}, VarType: stringIdent}
	DatePartParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "date"}, VarType: dateTimeIdent}
	DatePartFuncType := &ast.FuncType{Params: []ast.ParamItem{*DatePartParam1Item, *DatePartParam2Item}, Result: longIdent}
	DatePartDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: DatePartIdent, FuncType: DatePartFuncType, Body: nil}
	DatePartIdent.Decl = DatePartDecl

	DateSerialIdent := &ast.Identifier{Tok: &universeToken, Name: "DateSerial"}
	DateSerialParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "year"}, VarType: longIdent}
	DateSerialParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "month"}, VarType: longIdent}
	DateSerialParam3Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "day"}, VarType: longIdent}
	DateSerialFuncType := &ast.FuncType{Params: []ast.ParamItem{*DateSerialParam1Item, *DateSerialParam2Item, *DateSerialParam3Item}, Result: dateTimeIdent}
	DateSerialDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: DateSerialIdent, FuncType: DateSerialFuncType, Body: nil}
	DateSerialIdent.Decl = DateSerialDecl

	DateValueIdent := &ast.Identifier{Tok: &universeToken, Name: "DateValue"}
	DateValueParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "date"}, VarType: stringIdent}
	DateValueFuncType := &ast.FuncType{Params: []ast.ParamItem{*DateValueParam1Item}, Result: dateTimeIdent}
	DateValueDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: DateValueIdent, FuncType: DateValueFuncType, Body: nil}
	DateValueIdent.Decl = DateValueDecl

	DayIdent := &ast.Identifier{Tok: &universeToken, Name: "Day"}
	DayParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "date"}, VarType: dateTimeIdent}
	DayFuncType := &ast.FuncType{Params: []ast.ParamItem{*DayParam1Item}, Result: longIdent}
	DayDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: DayIdent, FuncType: DayFuncType, Body: nil}
	DayIdent.Decl = DayDecl

	HourIdent := &ast.Identifier{Tok: &universeToken, Name: "Hour"}
	HourParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "time"}, VarType: dateTimeIdent}
	HourFuncType := &ast.FuncType{Params: []ast.ParamItem{*HourParam1Item}, Result: longIdent}
	HourDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: HourIdent, FuncType: HourFuncType, Body: nil}
	HourIdent.Decl = HourDecl

	MinuteIdent := &ast.Identifier{Tok: &universeToken, Name: "Minute"}
	MinuteParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "time"}, VarType: dateTimeIdent}
	MinuteFuncType := &ast.FuncType{Params: []ast.ParamItem{*MinuteParam1Item}, Result: longIdent}
	MinuteDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: MinuteIdent, FuncType: MinuteFuncType, Body: nil}
	MinuteIdent.Decl = MinuteDecl

	MonthIdent := &ast.Identifier{Tok: &universeToken, Name: "Month"}
	MonthParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "date"}, VarType: dateTimeIdent}
	MonthFuncType := &ast.FuncType{Params: []ast.ParamItem{*MonthParam1Item}, Result: longIdent}
	MonthDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: MonthIdent, FuncType: MonthFuncType, Body: nil}
	MonthIdent.Decl = MonthDecl

	NowIdent := &ast.Identifier{Tok: &universeToken, Name: "Now"}
	NowFuncType := &ast.FuncType{Params: nil, Result: dateTimeIdent}
	NowDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: NowIdent, FuncType: NowFuncType, Body: nil}
	NowIdent.Decl = NowDecl

	SecondIdent := &ast.Identifier{Tok: &universeToken, Name: "Second"}
	SecondParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "time"}, VarType: dateTimeIdent}
	SecondFuncType := &ast.FuncType{Params: []ast.ParamItem{*SecondParam1Item}, Result: longIdent}
	SecondDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: SecondIdent, FuncType: SecondFuncType, Body: nil}
	SecondIdent.Decl = SecondDecl

	TimeIdent := &ast.Identifier{Tok: &universeToken, Name: "Time"}
	TimeFuncType := &ast.FuncType{Params: nil, Result: dateTimeIdent}
	TimeDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: TimeIdent, FuncType: TimeFuncType, Body: nil}
	TimeIdent.Decl = TimeDecl

	TimerIdent := &ast.Identifier{Tok: &universeToken, Name: "Timer"}
	TimerFuncType := &ast.FuncType{Params: nil, Result: doubleIdent}
	TimerDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: TimerIdent, FuncType: TimerFuncType, Body: nil}
	TimerIdent.Decl = TimerDecl

	TimeSerialIdent := &ast.Identifier{Tok: &universeToken, Name: "TimeSerial"}
	TimeSerialParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "hour"}, VarType: longIdent}
	TimeSerialParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "minute"}, VarType: longIdent}
	TimeSerialParam3Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "second"}, VarType: longIdent}
	TimeSerialFuncType := &ast.FuncType{Params: []ast.ParamItem{*TimeSerialParam1Item, *TimeSerialParam2Item, *TimeSerialParam3Item}, Result: dateTimeIdent}
	TimeSerialDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: TimeSerialIdent, FuncType: TimeSerialFuncType, Body: nil}
	TimeSerialIdent.Decl = TimeSerialDecl

	TimeValueIdent := &ast.Identifier{Tok: &universeToken, Name: "TimeValue"}
	TimeValueParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "time"}, VarType: stringIdent}
	TimeValueFuncType := &ast.FuncType{Params: []ast.ParamItem{*TimeValueParam1Item}, Result: dateTimeIdent}
	TimeValueDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: TimeValueIdent, FuncType: TimeValueFuncType, Body: nil}
	TimeValueIdent.Decl = TimeValueDecl

	WeekdayIdent := &ast.Identifier{Tok: &universeToken, Name: "Weekday"}
	WeekdayParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "date"}, VarType: dateTimeIdent}
	WeekdayFuncType := &ast.FuncType{Params: []ast.ParamItem{*WeekdayParam1Item}, Result: longIdent}
	WeekdayDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: WeekdayIdent, FuncType: WeekdayFuncType, Body: nil}
	WeekdayIdent.Decl = WeekdayDecl

	YearIdent := &ast.Identifier{Tok: &universeToken, Name: "Year"}
	YearParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "date"}, VarType: dateTimeIdent}
	YearFuncType := &ast.FuncType{Params: []ast.ParamItem{*YearParam1Item}, Result: longIdent}
	YearDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: YearIdent, FuncType: YearFuncType, Body: nil}
	YearIdent.Decl = YearDecl

	// --------------------------------
	// ------- conversion functions -------
	// --------------------------------

	CBoolIdent := &ast.Identifier{Tok: &universeToken, Name: "CBool"}
	CBoolParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "expression"}, VarType: variantIdent}
	CBoolFuncType := &ast.FuncType{Params: []ast.ParamItem{*CBoolParam1Item}, Result: booleanIdent}
	CBoolDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: CBoolIdent, FuncType: CBoolFuncType, Body: nil}
	CBoolIdent.Decl = CBoolDecl

	CDateIdent := &ast.Identifier{Tok: &universeToken, Name: "CDate"}
	CDateParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "expression"}, VarType: variantIdent}
	CDateFuncType := &ast.FuncType{Params: []ast.ParamItem{*CDateParam1Item}, Result: dateTimeIdent}
	CDateDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: CDateIdent, FuncType: CDateFuncType, Body: nil}
	CDateIdent.Decl = CDateDecl

	CDblIdent := &ast.Identifier{Tok: &universeToken, Name: "CDbl"}
	CDblParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "expression"}, VarType: variantIdent}
	CDblFuncType := &ast.FuncType{Params: []ast.ParamItem{*CDblParam1Item}, Result: doubleIdent}
	CDblDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: CDblIdent, FuncType: CDblFuncType, Body: nil}
	CDblIdent.Decl = CDblDecl

	CLngIdent := &ast.Identifier{Tok: &universeToken, Name: "CLng"}
	CLngParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "expression"}, VarType: variantIdent}
	CLngFuncType := &ast.FuncType{Params: []ast.ParamItem{*CLngParam1Item}, Result: longIdent}
	CLngDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: CLngIdent, FuncType: CLngFuncType, Body: nil}
	CLngIdent.Decl = CLngDecl

	CStrIdent := &ast.Identifier{Tok: &universeToken, Name: "CStr"}
	CStrParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "expression"}, VarType: variantIdent}
	CStrFuncType := &ast.FuncType{Params: []ast.ParamItem{*CStrParam1Item}, Result: stringIdent}
	CStrDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: CStrIdent, FuncType: CStrFuncType, Body: nil}
	CStrIdent.Decl = CStrDecl

	CvarIdent := &ast.Identifier{Tok: &universeToken, Name: "Cvar"}
	CvarParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "expression"}, VarType: variantIdent}
	CvarFuncType := &ast.FuncType{Params: []ast.ParamItem{*CvarParam1Item}, Result: variantIdent}
	CvarDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: CvarIdent, FuncType: CvarFuncType, Body: nil}
	CvarIdent.Decl = CvarDecl

	AscIdent := &ast.Identifier{Tok: &universeToken, Name: "Asc"}
	AscParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "character"}, VarType: stringIdent}
	AscFuncType := &ast.FuncType{Params: []ast.ParamItem{*AscParam1Item}, Result: longIdent}
	AscDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: AscIdent, FuncType: AscFuncType, Body: nil}
	AscIdent.Decl = AscDecl

	FormatIdent := &ast.Identifier{Tok: &universeToken, Name: "Format"}
	FormatParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "expression"}, VarType: variantIdent}
	FormatParam2Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "format"}, VarType: stringIdent}
	FormatFuncType := &ast.FuncType{Params: []ast.ParamItem{*FormatParam1Item, *FormatParam2Item}, Result: stringIdent}
	FormatDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: FormatIdent, FuncType: FormatFuncType, Body: nil}
	FormatIdent.Decl = FormatDecl

	HexIdent := &ast.Identifier{Tok: &universeToken, Name: "Hex"}
	HexParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "number"}, VarType: longIdent}
	HexFuncType := &ast.FuncType{Params: []ast.ParamItem{*HexParam1Item}, Result: stringIdent}
	HexDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: HexIdent, FuncType: HexFuncType, Body: nil}
	HexIdent.Decl = HexDecl

	OctIdent := &ast.Identifier{Tok: &universeToken, Name: "Oct"}
	OctParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "number"}, VarType: longIdent}
	OctFuncType := &ast.FuncType{Params: []ast.ParamItem{*OctParam1Item}, Result: stringIdent}
	OctDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: OctIdent, FuncType: OctFuncType, Body: nil}
	OctIdent.Decl = OctDecl

	// --------------------------------
	// ------- mathematical functions -------
	// --------------------------------

	AbsIdent := &ast.Identifier{Tok: &universeToken, Name: "Abs"}
	AbsParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "number"}, VarType: doubleIdent}
	AbsFuncType := &ast.FuncType{Params: []ast.ParamItem{*AbsParam1Item}, Result: doubleIdent}
	AbsDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: AbsIdent, FuncType: AbsFuncType, Body: nil}
	AbsIdent.Decl = AbsDecl

	AtnIdent := &ast.Identifier{Tok: &universeToken, Name: "Atn"}
	AtnParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "number"}, VarType: doubleIdent}
	AtnFuncType := &ast.FuncType{Params: []ast.ParamItem{*AtnParam1Item}, Result: doubleIdent}
	AtnDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: AtnIdent, FuncType: AtnFuncType, Body: nil}
	AtnIdent.Decl = AtnDecl

	CosIdent := &ast.Identifier{Tok: &universeToken, Name: "Cos"}
	CosParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "number"}, VarType: doubleIdent}
	CosFuncType := &ast.FuncType{Params: []ast.ParamItem{*CosParam1Item}, Result: doubleIdent}
	CosDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: CosIdent, FuncType: CosFuncType, Body: nil}
	CosIdent.Decl = CosDecl

	ExpIdent := &ast.Identifier{Tok: &universeToken, Name: "Exp"}
	ExpParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "number"}, VarType: doubleIdent}
	ExpFuncType := &ast.FuncType{Params: []ast.ParamItem{*ExpParam1Item}, Result: doubleIdent}
	ExpDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: ExpIdent, FuncType: ExpFuncType, Body: nil}
	ExpIdent.Decl = ExpDecl

	FixIdent := &ast.Identifier{Tok: &universeToken, Name: "Fix"}
	FixParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "number"}, VarType: doubleIdent}
	FixFuncType := &ast.FuncType{Params: []ast.ParamItem{*FixParam1Item}, Result: longIdent}
	FixDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: FixIdent, FuncType: FixFuncType, Body: nil}
	FixIdent.Decl = FixDecl

	intIdent := &ast.Identifier{Tok: &universeToken, Name: "Int"}
	IntParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "number"}, VarType: doubleIdent}
	IntFuncType := &ast.FuncType{Params: []ast.ParamItem{*IntParam1Item}, Result: longIdent}
	intDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: intIdent, FuncType: IntFuncType, Body: nil}
	integerIdent.Decl = intDecl

	LogIdent := &ast.Identifier{Tok: &universeToken, Name: "Log"}
	LogParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "number"}, VarType: doubleIdent}
	LogFuncType := &ast.FuncType{Params: []ast.ParamItem{*LogParam1Item}, Result: doubleIdent}
	LogDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: LogIdent, FuncType: LogFuncType, Body: nil}
	LogIdent.Decl = LogDecl

	RndIdent := &ast.Identifier{Tok: &universeToken, Name: "Rnd"}
	RndFuncType := &ast.FuncType{Params: nil, Result: doubleIdent}
	RndDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: RndIdent, FuncType: RndFuncType, Body: nil}
	RndIdent.Decl = RndDecl

	SgnIdent := &ast.Identifier{Tok: &universeToken, Name: "Sgn"}
	SgnParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "number"}, VarType: doubleIdent}
	SgnFuncType := &ast.FuncType{Params: []ast.ParamItem{*SgnParam1Item}, Result: longIdent}
	SgnDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: SgnIdent, FuncType: SgnFuncType, Body: nil}
	SgnIdent.Decl = SgnDecl

	SinIdent := &ast.Identifier{Tok: &universeToken, Name: "Sin"}
	SinParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "number"}, VarType: doubleIdent}
	SinFuncType := &ast.FuncType{Params: []ast.ParamItem{*SinParam1Item}, Result: doubleIdent}
	SinDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: SinIdent, FuncType: SinFuncType, Body: nil}
	SinIdent.Decl = SinDecl

	SqrIdent := &ast.Identifier{Tok: &universeToken, Name: "Sqr"}
	SqrParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "number"}, VarType: doubleIdent}
	SqrFuncType := &ast.FuncType{Params: []ast.ParamItem{*SqrParam1Item}, Result: doubleIdent}
	SqrDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: SqrIdent, FuncType: SqrFuncType, Body: nil}
	SqrIdent.Decl = SqrDecl

	TanIdent := &ast.Identifier{Tok: &universeToken, Name: "Tan"}
	TanParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "number"}, VarType: doubleIdent}
	TanFuncType := &ast.FuncType{Params: []ast.ParamItem{*TanParam1Item}, Result: doubleIdent}
	TanDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: TanIdent, FuncType: TanFuncType, Body: nil}
	TanIdent.Decl = TanDecl

	// --------------------------------
	// ------- array functions -------
	// --------------------------------

	LBoundIdent := &ast.Identifier{Tok: &universeToken, Name: "LBound"}
	LBoundParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "array"}, VarType: variantIdent}
	LBoundParam2Item := &ast.ParamItem{Optional: true, VarName: &ast.Identifier{Tok: &universeToken, Name: "dimension"}, VarType: longIdent}
	LBoundFuncType := &ast.FuncType{Params: []ast.ParamItem{*LBoundParam1Item, *LBoundParam2Item}, Result: longIdent}
	LBoundDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: LBoundIdent, FuncType: LBoundFuncType, Body: nil}
	LBoundIdent.Decl = LBoundDecl

	UBoundIdent := &ast.Identifier{Tok: &universeToken, Name: "UBound"}
	UBoundParam1Item := &ast.ParamItem{VarName: &ast.Identifier{Tok: &universeToken, Name: "array"}, VarType: variantIdent}
	UBoundParam2Item := &ast.ParamItem{Optional: true, VarName: &ast.Identifier{Tok: &universeToken, Name: "dimension"}, VarType: longIdent}
	UBoundFuncType := &ast.FuncType{Params: []ast.ParamItem{*UBoundParam1Item, *UBoundParam2Item}, Result: longIdent}
	UBoundDecl := &ast.FuncDecl{FunctionKw: &universeToken, FuncName: UBoundIdent, FuncType: UBoundFuncType, Body: nil}
	UBoundIdent.Decl = UBoundDecl

	// --------------------------------
	// ------- boolean constant -------
	// --------------------------------

	TrueIdent := &ast.Identifier{Tok: &universeToken, Name: "True"}
	TrueDecl := &ast.ConstDeclItem{ConstName: TrueIdent, ConstType: booleanIdent, ConstValue: &ast.BasicLit{Kind: token.BooleanLit, Value: "True"}}
	TrueIdent.Decl = TrueDecl

	FalseIdent := &ast.Identifier{Tok: &universeToken, Name: "False"}
	FalseDecl := &ast.ConstDeclItem{ConstName: FalseIdent, ConstType: booleanIdent, ConstValue: &ast.BasicLit{Kind: token.BooleanLit, Value: "False"}}
	FalseIdent.Decl = FalseDecl

	universeDecls := []*ast.TypeDef{
		longDecl,
		integerDecl,
		currencyDecl,
		doubleDecl,
		singleDecl,
		stringDecl,
		dateTimeDecl,
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
		StringDecl,
		StrReverseDecl,
		TrimDecl,
		UCaseDecl,
		DateDecl,
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
		ExpDecl,
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
	// for _, node := range file.Nodes {
	// 	if decl, ok := node.(ast.Decl); ok {
	// 		if err := fileScope.Insert(decl); err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	// scope specifies the current lexical scope.
	scope := fileScope

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
			case ast.VarDecl, *ast.ConstDeclItem:
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
						ConstName:  &ast.Identifier{Name: value.Name, Tok: value.Token()},
						ConstType:  constType,
						ConstValue: &ast.BasicLit{Kind: token.LongLit, Value: strconv.Itoa(index)},
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
						if decl != nil && constDecl != nil && decl.Token().Position != constDecl.Token().Position {
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
				for _, param := range fn.GetParams() {
					// add parameters to the function scope
					if err := scope.Insert(&param); err != nil {
						return err
					}
				}
				// add the function to the scope, for return resolution
				if err := scope.Insert(n); err != nil {
					return err
				}

			}

		case *ast.Identifier:
			// if node is an expression containing an identifier
			// then we need to resolve the identifier

			decl, ok := scope.Lookup(n.Name, false)
			if !ok {
				return errors.Newf(n.Token().Position, "undeclared identifier %q", n)
			}
			// cannot use function name as a declaration
			if _, ok := decl.(*ast.FuncDecl); ok {
				if parent, ok := n.GetParent().(*ast.UserDefinedType); ok {
					return errors.Newf(n.Token().Position, "cannot use function name as a declaration in %q", parent.Name)
				}
			}
			n.Decl = decl
		}
		return nil
	}

	// after reverts to the outer scope after traversing block statements.
	after := func(n ast.Node) error {

		if _, ok := n.(ast.FuncOrSub); ok {
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

	return nil
}
