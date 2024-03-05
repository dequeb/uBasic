package ide

import (
	"image/color"
	"io"
	"time"
	"uBasic/exec"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/explorer"
	"gioui.org/x/richtext"
)

type FourPaneArea struct {
	mainSplit        *Split
	leftSplit        *Split
	rightSplit       *Split
	sourceCode       widget.Editor
	variables        widget.Editor
	callStack        widget.Editor
	terminal         widget.Editor
	breakpoints      widget.Editor
	stopButton       widget.Clickable
	continueButton   widget.Clickable
	stepButton       widget.Clickable
	breakpointButton widget.Clickable
	openButton       widget.Clickable
	saveButton       widget.Clickable
	closeButton      widget.Clickable
}

func newFourPaneArea() *FourPaneArea {
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

var (
	fonts = gofont.Collection()
	th    = material.NewTheme()
)

func (area *FourPaneArea) ButtonBar(gtx layout.Context, th *material.Theme) layout.Dimensions {
	flexbox := layout.Flex{Axis: layout.Horizontal}

	// ... before laying it out, one inside the other
	return flexbox.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			for area.stopButton.Clicked(gtx) {
				exec.Debug.Stop()
			}
			return material.Button(th, &area.stopButton, "Stop").Layout(gtx)
		},
		),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			for area.continueButton.Clicked(gtx) {
				exec.Debug.Continue()
			}
			return material.Button(th, &area.continueButton, "Continue").Layout(gtx)
		},
		),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			for area.stepButton.Clicked(gtx) {
				exec.Debug.Step()
			}
			return material.Button(th, &area.stepButton, "Step").Layout(gtx)
		},
		),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			button := material.Button(th, &area.breakpointButton, "Breakpoint")
			gtx = gtx.Disabled()
			return button.Layout(gtx)
		},
		),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			exp := explorer.NewExplorer(&app.Window{})
			for area.openButton.Clicked(gtx) {
				exp.ChooseFile(".bas")
			}

			button := material.Button(th, &area.openButton, "Open")
			gtx = gtx.Disabled()
			return button.Layout(gtx)
		},
		),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			exp := explorer.NewExplorer(&app.Window{})
			for area.saveButton.Clicked(gtx) {
				var file io.WriteCloser
				var err error
				file, err = exp.CreateFile("test.base")
				if err != nil {
					panic(err)
				}
				defer file.Close()
				_, err = file.Write([]byte(area.sourceCode.Text()))
				if err != nil {
					panic(err)
				}
			}
			button := material.Button(th, &area.saveButton, "Save")
			gtx = gtx.Disabled()
			return button.Layout(gtx)
		},
		),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			for area.closeButton.Clicked(gtx) {

			}
			button := material.Button(th, &area.closeButton, "Close")
			// gtx = gtx.Disabled()
			return button.Layout(gtx)
		},
		),
	)

}

func (area *FourPaneArea) MainSplit(gtx layout.Context, th *material.Theme) layout.Dimensions {
	flexbox := layout.Flex{Axis: layout.Vertical}

	// Define insets ...
	margins := layout.Inset{
		Top:    unit.Dp(1),
		Right:  unit.Dp(1),
		Bottom: unit.Dp(1),
		Left:   unit.Dp(1),
	}

	// ... before laying it out, one inside the other
	return margins.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			return flexbox.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return area.ButtonBar(gtx, th)
				},
				),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return area.mainSplit.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return area.LeftSplit(gtx, th)
					}, func(gtx layout.Context) layout.Dimensions {
						return area.RightPanes(gtx, th)
					})
				},
				),
			)
		},
	)
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
	th.Shaper = text.NewShaper(text.WithCollection(fonts))
	// allocate persistent state for interactive text. This
	// needs to be persisted across frames.
	var state richtext.InteractiveText
	// define the colors for the interactive text

	// https://pkg.go.dev/gioui.org/x@v0.5.0/richtext
	var currentLine int
	if exec.Debug.Env.From.Token() != nil {
		currentLine = exec.Debug.Env.From.Token().Position.Line
	}
	var spans = ColorText(exec.Debug.Ast, currentLine)

	// render the rich text into the operation list
	ed := richtext.Text(&state, th.Shaper, spans...)

	// area.sourceCode.SetText(sourceCode)
	// area.sourceCode.ReadOnly = true // make it read-only
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
	if !exec.Debug.Running {
		area.variables.SetText(exec.Debug.Env.String())
	}
	area.variables.WrapPolicy = text.WrapHeuristically
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
	area.callStack.ReadOnly = true // make it read-only
	ed1 := material.Editor(th, &area.callStack, "call stack")
	if !exec.Debug.Running {
		area.callStack.SetText(exec.Debug.Env.CallStack())
	}
	area.breakpoints.ReadOnly = true // make it read-only
	ed2 := material.Editor(th, &area.breakpoints, "breakpoints")
	if !exec.Debug.Running {
		area.breakpoints.SetText(exec.Debug.Breakpoints.String())
	}
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
	area.terminal.SetText(exec.Debug.Terminal.String())

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

var running chan bool

func Run(w *app.Window) error {
	th := material.NewTheme()
	area := newFourPaneArea()
	var ops op.Ops

	// listen for events in runner channel
	go func() {
		for range running {
			time.Sleep(time.Second * 1)
			w.Invalidate()
		}
	}()

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
