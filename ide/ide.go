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

type Debugger struct {
	// sourceCode  []string
	// breakpoints []bool
	// variables   map[string]string
	// callStack   []string
}

var (
	hsplit HSplit
	vsplit VSplit
	red    = color.NRGBA{R: 255, A: 255}
	blue   = color.NRGBA{B: 255, A: 255}
)

func horizontalSplit(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return hsplit.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return FillWithLabel(gtx, th, "Left", red)
	}, func(gtx layout.Context) layout.Dimensions {
		return FillWithLabel(gtx, th, "Right", blue)
	})
}

func VerticalSplit(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return vsplit.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return FillWithLabel(gtx, th, "up", red)
	}, func(gtx layout.Context) layout.Dimensions {
		return FillWithLabel(gtx, th, "Down", blue)
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
	var ops op.Ops
	for {
		switch e := w.NextEvent().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)

			VerticalSplit(gtx, th)

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
