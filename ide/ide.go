package ide

import (
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
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
	mainSplit  *Split
	leftSplit  *Split
	rightSplit *Split
	widgets    []func(gtx layout.Context, th *material.Theme, text string, backgroundColor color.NRGBA) layout.Dimensions
}

var (
	red  = color.NRGBA{R: 255, A: 255}
	blue = color.NRGBA{B: 255, A: 255}
)

func NewFourPaneArea() *FourPaneArea {
	a := &FourPaneArea{
		mainSplit: &Split{
			Ratio: -0.5, direction: Vertical, Bar: 7},
		leftSplit: &Split{
			Ratio: 0, direction: Horizontal, Bar: 10},
		rightSplit: &Split{
			Ratio: 0.5, direction: Horizontal, Bar: 10},
	}
	a.widgets = make([]func(gtx layout.Context, th *material.Theme, text string, backgroundColor color.NRGBA) layout.Dimensions, 4)
	a.widgets[0] = FillWithLabel
	a.widgets[1] = FillWithLabel
	a.widgets[2] = FillWithLabel
	a.widgets[3] = FillWithLabel
	return a
}

func (area *FourPaneArea) MainSplit(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return area.mainSplit.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return area.leftSplit.LeftSplit(gtx, th)
	}, func(gtx layout.Context) layout.Dimensions {
		return area.rightSplit.RightSplit(gtx, th)
	})
}

func (split *Split) LeftSplit(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return split.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return FillWithLabel(gtx, th, "top Left ", red)
	}, func(gtx layout.Context) layout.Dimensions {
		return FillWithLabel(gtx, th, "bottom left", blue)
	})
}

func (split *Split) RightSplit(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return split.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return FillWithLabel(gtx, th, "top right", red)
	}, func(gtx layout.Context) layout.Dimensions {
		return FillWithLabel(gtx, th, "bottom Right", blue)
	})
}
func FillWithLabel(gtx layout.Context, th *material.Theme, text string, backgroundColor color.NRGBA) layout.Dimensions {
	ColorBox(gtx, gtx.Constraints.Max, backgroundColor)
	return layout.Center.Layout(gtx, material.H3(th, text).Layout)
}

func ColorBox(gtx layout.Context, point image.Point, backgroundColor color.NRGBA) {
	d := image.Point{X: point.X, Y: point.Y}
	rect := clip.Rect(image.Rectangle{Max: d})
	paint.FillShape(gtx.Ops, backgroundColor, rect.Op())
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

			// VerticalSplit(gtx, th)

			// Define an large label with an appropriate text:
			title := material.H1(th, "Hello, Gio")

			// Change the color of the label.
			maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
			title.Color = maroon

			// Change the position of the label.
			title.Alignment = text.Middle

			// Draw the label to the graphics context.
			title.Layout(gtx)

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}
