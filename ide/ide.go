package ide

import (
	"image/color"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type D layout.Dimensions
type C layout.Context

type Debugger struct {
	// sourceCode  []string
	// breakpoints []bool
	// variables   map[string]string
	// callStack   []string
}

type FourPaneArea struct {
	mainSplit   *Split
	leftSplit   *Split
	rightSplit  *Split
	sourceCode  widget.Editor
	variables   widget.Editor
	callStack   widget.Editor
	terminal    widget.Editor
	breakpoints widget.Editor
}

func NewFourPaneArea() *FourPaneArea {
	a := &FourPaneArea{
		mainSplit: &Split{
			Ratio: -0.5, direction: Vertical, Bar: 7},
		leftSplit: &Split{
			Ratio: 0, direction: Horizontal, Bar: 10},
		rightSplit: &Split{
			Ratio: 0.5, direction: Horizontal, Bar: 10},
	}
	return a
}

func (area *FourPaneArea) MainSplit(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return area.mainSplit.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return area.LeftSplit(gtx, th)
	}, func(gtx layout.Context) layout.Dimensions {
		return area.RightPanes(gtx, th)
	})
}

func (area *FourPaneArea) LeftSplit(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return area.leftSplit.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return area.Variables(gtx, th)
	}, func(gtx layout.Context) layout.Dimensions {
		return area.CallStack(gtx, th)
	})
}

func (area *FourPaneArea) RightPanes(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return area.rightSplit.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

		return area.SouceEditor(gtx, th)
	}, func(gtx layout.Context) layout.Dimensions {
		return area.Terminal(gtx, th)
	})
}
func (area *FourPaneArea) SouceEditor(gtx layout.Context, th *material.Theme) layout.Dimensions {
	ed := material.Editor(th, &area.sourceCode, "source code")
	// Define insets ...
	margins := layout.Inset{
		Top:    unit.Dp(3),
		Right:  unit.Dp(3),
		Bottom: unit.Dp(0),
		Left:   unit.Dp(0),
	}

	// ... and borders ...
	border := widget.Border{
		Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
		CornerRadius: unit.Dp(3),
		Width:        unit.Dp(2),
	}
	// ... before laying it out, one inside the other
	return margins.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			return border.Layout(gtx, ed.Layout)
		},
	)
}

func (area *FourPaneArea) Variables(gtx layout.Context, th *material.Theme) layout.Dimensions {
	area.variables.ReadOnly = true // make it read-only
	ed := material.Editor(th, &area.variables, "variables")
	// Define insets ...
	margins := layout.Inset{
		Top:    unit.Dp(3),
		Right:  unit.Dp(0),
		Bottom: unit.Dp(0),
		Left:   unit.Dp(3),
	}

	// ... and borders ...
	border := widget.Border{
		Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
		CornerRadius: unit.Dp(3),
		Width:        unit.Dp(2),
	}
	// ... before laying it out, one inside the other
	return margins.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			return border.Layout(gtx, ed.Layout)
		},
	)
}

func (area *FourPaneArea) CallStack(gtx layout.Context, th *material.Theme) layout.Dimensions {
	area.variables.ReadOnly = true // make it read-only
	ed1 := material.Editor(th, &area.callStack, "call stack")
	ed1.Editor.SetText("call stack\ncall stack\ncall stack\ncall stack\ncall stack\ncall stack\ncall stack\ncall stack\ncall stack")

	area.breakpoints.ReadOnly = true // make it read-only
	ed2 := material.Editor(th, &area.breakpoints, "breakpoints")
	ed2.Editor.SetText("breakpoints\nbreakpoints\nbreakpoints\nbreakpoints\nbreakpoints\nbreakpoints\nbreakpoints\nbreakpoints\nbreakpoints")
	// Define insets ...
	margins := layout.Inset{
		Top:    unit.Dp(0),
		Right:  unit.Dp(0),
		Bottom: unit.Dp(3),
		Left:   unit.Dp(3),
	}

	// ... and borders ...
	border := widget.Border{
		Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
		CornerRadius: unit.Dp(3),
		Width:        unit.Dp(2),
	}

	flexbox := layout.Flex{Axis: layout.Vertical}

	// ... before laying it out, one inside the other
	return margins.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			return flexbox.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return border.Layout(gtx, ed1.Layout)
				},
				),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return border.Layout(gtx, ed2.Layout)
				},
				),
			)
		},
	)
}

func (area *FourPaneArea) Terminal(gtx layout.Context, th *material.Theme) layout.Dimensions {
	area.variables.ReadOnly = true // make it read-only
	ed := material.Editor(th, &area.terminal, "terminal")
	// Define insets ...
	margins := layout.Inset{
		Top:    unit.Dp(0),
		Right:  unit.Dp(3),
		Bottom: unit.Dp(3),
		Left:   unit.Dp(0),
	}

	// ... and borders ...
	border := widget.Border{
		Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
		CornerRadius: unit.Dp(3),
		Width:        unit.Dp(2),
	}
	// ... before laying it out, one inside the other
	return margins.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			return border.Layout(gtx, ed.Layout)
		},
	)
}

func (d *Debugger) Run(w *app.Window) error {
	th := material.NewTheme()
	area := NewFourPaneArea()
	var ops op.Ops
	for {
		switch e := w.NextEvent().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)

			area.MainSplit(gtx, th)

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}
