package main

import (
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("BOTMC - Created by Suboorzen"))
		w.Option(app.Size(unit.Dp(360), unit.Dp(640)))

		th := material.NewTheme()
		th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

		var ops op.Ops

		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				os.Exit(0)

			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				paint.Fill(gtx.Ops, color.NRGBA{R: 18, G: 20, B: 30, A: 255})

				layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:      layout.Vertical,
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							title := material.H3(th, "BOTMC")
							title.Color = color.NRGBA{R: 100, G: 180, B: 255, A: 255}
							return title.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Body1(th, "Created by Suboorzen")
							label.Color = color.NRGBA{R: 220, G: 220, B: 230, A: 255}
							return label.Layout(gtx)
						}),
					)
				})

				e.Frame(gtx.Ops)
			}
		}
	}()

	app.Main()
}
