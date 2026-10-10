package main

import (
	"fmt"
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type C = layout.Context
type D = layout.Dimensions

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("BOTMC - Bedrock Edition"))
		w.Option(app.Size(unit.Dp(380), unit.Dp(720)))

		th := material.NewTheme()
		th.Shaper = text.NewShaper(
			text.WithCollection(gofont.Collection()),
		)

		var ops op.Ops
		var list widget.List
		list.List.Axis = layout.Vertical

		var ip, port, username widget.Editor
		ip.SingleLine = true
		port.SingleLine = true
		username.SingleLine = true
		ip.SetText("play.example.com")
		port.SetText("19132")
		username.SetText("BOTMC")

		var join, stop, rejoin, versionButton widget.Clickable
		var autoReconnect, notifications, haptics widget.Bool
		autoReconnect.Value = true
		notifications.Value = true
		haptics.Value = false

		versions := []string{
			"Auto (recommended)",
			"Bedrock 1.21.x",
			"Bedrock 1.20.x",
		}
		versionIndex := 0
		status := "Ready - enter your Bedrock server"
		var slots [27]widget.Clickable

		slotNames := []string{
			"👤", "👤", "👤", "🟩", "🟩", "🟩", "🏠", "🌍", "❌",
			"👤", "👤", "👤", "🟩", "🟩", "🟩", "🏠", "🌍", "❌",
			"👤", "👤", "👤", "🟩", "🟩", "🟩", "🏠", "🌍", "❌",
		}

		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				os.Exit(0)

			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)

				if join.Clicked(gtx) {
					status = "Connection engine pending - not connected"
				}
				if stop.Clicked(gtx) {
					status = "Stop requested"
				}
				if rejoin.Clicked(gtx) {
					status = "Rejoin requested - connection engine pending"
				}
				if versionButton.Clicked(gtx) {
					versionIndex = (versionIndex + 1) % len(versions)
				}
				for i := range slots {
					if slots[i].Clicked(gtx) {
						status = fmt.Sprintf(
							"Selected GUI slot %d: %s",
							i+1, slotNames[i],
						)
					}
				}

				paint.Fill(
					gtx.Ops,
					color.NRGBA{R: 15, G: 19, B: 30, A: 255},
				)

				material.List(th, &list).Layout(gtx, 24,
					func(gtx C, index int) D {
						return layout.Inset{
							Top:    unit.Dp(5),
							Bottom: unit.Dp(5),
							Left:   unit.Dp(12),
							Right:  unit.Dp(12),
						}.Layout(gtx, func(gtx C) D {
							switch index {
							case 0:
								title := material.H4(th, "BOTMC")
								title.Color = color.NRGBA{
									R: 90, G: 175, B: 255, A: 255,
								}
								return title.Layout(gtx)

							case 1:
								label := material.Body2(
									th, "BEDROCK / MCPE  •  Created by Suboorzen",
								)
								label.Color = color.NRGBA{
									R: 190, G: 200, B: 215, A: 255,
								}
								return label.Layout(gtx)

							case 2:
								return material.Body1(th, "Server IP").Layout(gtx)

							case 3:
								return material.Editor(
									th, &ip, "Enter server IP",
								).Layout(gtx)

							case 4:
								return material.Body1(th, "Server Port").Layout(gtx)

							case 5:
								return material.Editor(
									th, &port, "19132",
								).Layout(gtx)

							case 6:
								return material.Body1(
									th, "Bedrock Version",
								).Layout(gtx)

							case 7:
								return material.Button(
									th, &versionButton, versions[versionIndex],
								).Layout(gtx)

							case 8:
								return material.Body1(th, "Bot Username").Layout(gtx)

							case 9:
								return material.Editor(
									th, &username, "Enter username",
								).Layout(gtx)

							case 10:
								return layout.Flex{
									Axis: layout.Horizontal,
								}.Layout(gtx,
									layout.Flexed(1, func(gtx C) D {
										return material.Button(
											th, &join, "JOIN",
										).Layout(gtx)
									}),
									layout.Flexed(1, func(gtx C) D {
										return material.Button(
											th, &stop, "STOP",
										).Layout(gtx)
									}),
									layout.Flexed(1, func(gtx C) D {
										return material.Button(
											th, &rejoin, "REJOIN",
										).Layout(gtx)
									}),
								)

							case 11:
								return material.Body2(
									th, "Status: "+status,
								).Layout(gtx)

							case 12:
								return material.H6(
									th, "⚙ SETTINGS",
								).Layout(gtx)

							case 13:
								return material.CheckBox(
									th, &autoReconnect, "Auto reconnect",
								).Layout(gtx)

							case 14:
								return material.CheckBox(
									th, &notifications, "Connection notifications",
								).Layout(gtx)

							case 15:
								return material.CheckBox(
									th, &haptics, "Haptic feedback",
								).Layout(gtx)

							case 16:
								return material.Body2(
									th, "Other planned settings: timeout, saved servers, theme and logs.",
								).Layout(gtx)

							case 17:
								return material.H6(
									th, "🎒 CHEST GUI • TPA / RTP",
								).Layout(gtx)

							case 18, 19, 20:
								row := index - 18
								children := make([]layout.FlexChild, 0, 9)
								for col := 0; col < 9; col++ {
									slot := row*9 + col
									children = append(children,
										layout.Flexed(1, func(gtx C) D {
											return material.Button(
												th, &slots[slot], slotNames[slot],
											).Layout(gtx)
										}),
									)
								}
								return layout.Flex{
									Axis: layout.Horizontal,
								}.Layout(gtx, children...)

							case 21:
								return material.Body2(
									th, "GUI actions need server connection support.",
								).Layout(gtx)

							case 22:
								return material.Body2(
									th, "Saved profiles and real menu actions are next.",
								).Layout(gtx)

							default:
								return D{}
							}
						})
					},
				)

				e.Frame(gtx.Ops)
			}
		}
	}()
	app.Main()
}
