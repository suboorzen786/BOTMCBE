package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type ServerProfile struct {
	Name    string
	Address string
	Port    string
	Ping    string
	Click   widget.Clickable
}

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("BOTMC PRO"))
		w.Option(app.Size(420, 850))
		if err := run(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	var ops op.Ops

	cfg := LoadConfig()

	servers := make([]ServerProfile, len(cfg.Servers))
	for i, s := range cfg.Servers {
		servers[i] = ServerProfile{
			Name:    s.Name,
			Address: s.Address,
			Port:    s.Port,
			Ping:    "35 ms",
		}
	}

	var msLoginBtn, logoutBtn, addServerBtn, confirmAddBtn, backHomeBtn, settingsBtn widget.Clickable
	var slideRightBtn, slideLeftBtn, sendChatBtn widget.Clickable
	var isAddingServer bool

	var nameInput, ipInput, portInput widget.Editor
	nameInput.SingleLine = true
	ipInput.SingleLine = true
	portInput.SingleLine = true

	var chatEditor widget.Editor
	chatEditor.SingleLine = true

	var autoLoginToggle, cmdAfterToggle, autoMoveToggle widget.Bool
	autoLoginToggle.Value = cfg.AutoLogin
	cmdAfterToggle.Value = cfg.CmdAfterLogin
	autoMoveToggle.Value = cfg.AutoMovement

	var cmdAfterInput widget.Editor
	cmdAfterInput.SetText(cfg.CmdAfterText)
	cmdAfterInput.SingleLine = true

	var trig1Cmd, trig1Val, trig2Cmd, trig2Val, trig3Cmd, trig3Val widget.Editor
	trig1Cmd.SingleLine = true
	trig1Val.SingleLine = true
	trig2Cmd.SingleLine = true
	trig2Val.SingleLine = true
	trig3Cmd.SingleLine = true
	trig3Val.SingleLine = true

	currentView := 0
	if !cfg.IsLoggedIn {
		currentView = -1
	}

	logMessages := []string{"[System] Ready to connect."}
	onlinePlayers := []string{"ramusar", "SharliLP", "Roquexd", "Steeph", "Darcyy"}

	var list, playerList widget.List
	list.Axis = layout.Vertical
	playerList.Axis = layout.Vertical

	bgColor := color.NRGBA{R: 20, G: 22, B: 28, A: 255}
	headerColor := color.NRGBA{R: 190, G: 40, B: 40, A: 255}
	cardColor := color.NRGBA{R: 30, G: 35, B: 45, A: 255}
	textColor := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	cyanAccent := color.NRGBA{R: 0, G: 240, B: 255, A: 255}

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if msLoginBtn.Clicked(gtx) {
				cfg.IsLoggedIn = true
				cfg.MicrosoftUser = "PlayerBot#786"
				_ = SaveConfig(cfg)
				currentView = 0
			}

			if logoutBtn.Clicked(gtx) {
				cfg.IsLoggedIn = false
				cfg.MicrosoftUser = ""
				_ = SaveConfig(cfg)
				currentView = -1
			}

			if settingsBtn.Clicked(gtx) { currentView = 4 }
			if backHomeBtn.Clicked(gtx) { currentView = 0; isAddingServer = false }

			if addServerBtn.Clicked(gtx) { isAddingServer = !isAddingServer }

			if confirmAddBtn.Clicked(gtx) && nameInput.Text() != "" && ipInput.Text() != "" {
				p := portInput.Text()
				if p == "" { p = "19132" }
				cfg.Servers = append(cfg.Servers, SavedServer{Name: nameInput.Text(), Address: ipInput.Text(), Port: p})
				_ = SaveConfig(cfg)
				servers = append(servers, ServerProfile{Name: nameInput.Text(), Address: ipInput.Text(), Port: p, Ping: "28 ms"})
				nameInput.SetText(""); ipInput.SetText(""); portInput.SetText("")
				isAddingServer = false
			}

			for i := range servers {
				if servers[i].Click.Clicked(gtx) {
					currentView = 1
					logMessages = append([]string{"[Joined] Connected to " + servers[i].Name}, logMessages...)
				}
			}

			if slideRightBtn.Clicked(gtx) && currentView >= 1 && currentView < 3 { currentView++ }
			if slideLeftBtn.Clicked(gtx) && currentView > 1 && currentView <= 3 { currentView-- }

			if sendChatBtn.Clicked(gtx) && chatEditor.Text() != "" {
				logMessages = append([]string{"You: " + chatEditor.Text()}, logMessages...)
				chatEditor.SetText("")
			}

			paint.Fill(gtx.Ops, bgColor)

			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return drawCard(gtx, headerColor, func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: 12, Bottom: 12, Left: 16, Right: 16}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									t := "BOTMC PRO"
									if currentView == -1 { t = "Microsoft Authentication" }
									if currentView == 1 { t = "Chat Console (Slide ➡)" }
									if currentView == 2 { t = "Minimap & Radar" }
									if currentView == 3 { t = "Online Players" }
									if currentView == 4 { t = "Advanced Settings" }
									lbl := material.H6(th, t)
									lbl.Color = textColor
									return lbl.Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if currentView == -1 { return layout.Dimensions{} }
									return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
										func() layout.Widget {
											if currentView != 0 && currentView != 4 {
												return func(gtx layout.Context) layout.Dimensions {
													b := material.Button(th, &backHomeBtn, "🏠")
													b.Background = color.NRGBA{}
													b.Color = textColor
													return b.Layout(gtx)
												}
											}
											return func(gtx layout.Context) layout.Dimensions { return layout.Dimensions{} }
										}(),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											b := material.Button(th, &settingsBtn, "⚙")
											b.Background = color.NRGBA{}
											b.Color = textColor
											return b.Layout(gtx)
										}),
									)
								}),
							)
						})
					})
				}),

				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Top: 4, Bottom: 4, Left: 4, Right: 4}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						switch currentView {
						case -1:
							return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(layout.Spacer{Height: 40}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
										return layout.Inset{Top: 20, Bottom: 20, Left: 20, Right: 20}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
												layout.Rigid(material.H6(th, "Microsoft Account Required").Layout),
												layout.Rigid(layout.Spacer{Height: 8}.Layout),
												layout.Rigid(material.Body2(th, "Connect your account to access MCPE Servers.").Layout),
												layout.Rigid(layout.Spacer{Height: 16}.Layout),
												layout.Rigid(func(gtx layout.Context) layout.Dimensions {
													b := material.Button(th, &msLoginBtn, "LOGIN WITH MICROSOFT")
													b.Background = cyanAccent
													b.Color = color.NRGBA{R: 20, G: 22, B: 28, A: 255}
													return b.Layout(gtx)
												}),
											)
										})
									})
								}),
							)

						case 0:
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
										return layout.Inset{Top: 12, Bottom: 12, Left: 14, Right: 14}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
												layout.Rigid(func(gtx layout.Context) layout.Dimensions {
													return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
														layout.Rigid(material.Body1(th, "👤 "+cfg.MicrosoftUser).Layout),
														layout.Rigid(func(gtx layout.Context) layout.Dimensions {
															lbl := material.Caption(th, "Microsoft Connected & Saved")
															lbl.Color = cyanAccent
															return lbl.Layout(gtx)
														}),
													)
												}),
												layout.Rigid(func(gtx layout.Context) layout.Dimensions {
													b := material.Button(th, &logoutBtn, "Logout")
													b.Background = color.NRGBA{R: 200, G: 50, B: 50, A: 255}
													return b.Layout(gtx)
												}),
											)
										})
									})
								}),
								layout.Rigid(layout.Spacer{Height: 8}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(material.Body1(th, "Saved Servers").Layout),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											b := material.Button(th, &addServerBtn, " + ")
											b.Background = cyanAccent
											b.Color = color.NRGBA{R: 20, G: 22, B: 28, A: 255}
											return b.Layout(gtx)
										}),
									)
								}),
								layout.Rigid(layout.Spacer{Height: 6}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if !isAddingServer { return layout.Dimensions{} }
									return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
										return layout.Inset{Top: 10, Bottom: 10, Left: 12, Right: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
												layout.Rigid(material.Editor(th, &nameInput, "Server Name").Layout),
												layout.Rigid(layout.Spacer{Height: 4}.Layout),
												layout.Rigid(material.Editor(th, &ipInput, "IP Address").Layout),
												layout.Rigid(layout.Spacer{Height: 4}.Layout),
												layout.Rigid(material.Editor(th, &portInput, "Port (19132)").Layout),
												layout.Rigid(layout.Spacer{Height: 6}.Layout),
												layout.Rigid(func(gtx layout.Context) layout.Dimensions {
													b := material.Button(th, &confirmAddBtn, "✔ Save Server")
													b.Background = color.NRGBA{R: 50, G: 200, B: 100, A: 255}
													return b.Layout(gtx)
												}),
											)
										})
									})
								}),
								layout.Rigid(layout.Spacer{Height: 4}.Layout),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									if len(servers) == 0 {
										return layout.Center.Layout(gtx, material.Body2(th, "No saved servers. Click '+' to add one.").Layout)
									}
									return material.List(th, &list).Layout(gtx, len(servers), func(gtx layout.Context, index int) layout.Dimensions {
										srv := &servers[index]
										return layout.Inset{Bottom: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
												return layout.Inset{Top: 10, Bottom: 10, Left: 12, Right: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
													return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
														layout.Rigid(material.H6(th, srv.Name).Layout),
														layout.Rigid(material.Caption(th, "IP: "+srv.Address+":"+srv.Port).Layout),
														layout.Rigid(layout.Spacer{Height: 6}.Layout),
														layout.Rigid(func(gtx layout.Context) layout.Dimensions {
															btn := material.Button(th, &srv.Click, "JOIN SERVER")
															btn.Background = cyanAccent
                               btn.Color = color.NRGBA{R: 20, G: 22, B: 28, A: 255}
															return btn.Layout(gtx)
														}),
													)
												})
											})
										})
									})
								}),
							)

						case 1:
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
										return layout.Inset{Top: 6, Bottom: 6, Left: 10, Right: 10}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
												layout.Rigid(material.Caption(th, "● Online (Chat)").Layout),
												layout.Rigid(material.Button(th, &slideRightBtn, "Slide Right ➡").Layout),
											)
										})
									})
								}),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
										return layout.Inset{Top: 8, Bottom: 8, Left: 8, Right: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											return material.List(th, &list).Layout(gtx, len(logMessages), func(gtx layout.Context, idx int) layout.Dimensions {
												return layout.Inset{Top: 3, Bottom: 3}.Layout(gtx, material.Body2(th, logMessages[idx]).Layout)
											})
										})
									})
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
										layout.Flexed(1, material.Editor(th, &chatEditor, "Type chat/command...").Layout),
										layout.Rigid(layout.Spacer{Width: 6}.Layout),
										layout.Rigid(material.Button(th, &sendChatBtn, "SEND").Layout),
									)
								}),
							)

						case 2:
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(material.Button(th, &slideRightBtn, "Slide Right -> Players").Layout),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
										return layout.Inset{Top: 20, Bottom: 20, Left: 16, Right: 16}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
												layout.Rigid(material.H6(th, "Minimap / Radar Active").Layout),
												layout.Rigid(material.Body1(th, "Coords: X=3745.94 Y=108.00 Z=-2366.58").Layout),
												layout.Rigid(material.Body2(th, "Ping: 35 ms").Layout),
											)
										})
									})
								}),
							)

						case 3:
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(material.Button(th, &slideLeftBtn, "⬅ Slide Left").Layout),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									return material.List(th, &playerList).Layout(gtx, len(onlinePlayers), func(gtx layout.Context, idx int) layout.Dimensions {
										return layout.Inset{Top: 6, Bottom: 6}.Layout(gtx, material.Body1(th, "👤 "+onlinePlayers[idx]).Layout)
									})
								}),
							)

						default:
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(material.CheckBox(th, &autoLoginToggle, "Auto Login (3s after join)").Layout),
								layout.Rigid(layout.Spacer{Height: 4}.Layout),
								layout.Rigid(material.CheckBox(th, &cmdAfterToggle, "Command After Login").Layout),
								layout.Rigid(layout.Spacer{Height: 4}.Layout),
								layout.Rigid(material.Editor(th, &cmdAfterInput, "Command (e.g. /lobby)").Layout),
								layout.Rigid(layout.Spacer{Height: 4}.Layout),
								layout.Rigid(material.CheckBox(th, &autoMoveToggle, "Auto Movement Every Min").Layout),
								layout.Rigid(layout.Spacer{Height: 10}.Layout),
								layout.Rigid(material.Body2(th, "3 Repeat Text Triggers (1-60 Min/Sec):").Layout),
								layout.Rigid(layout.Spacer{Height: 4}.Layout),
								layout.Rigid(material.Editor(th, &trig1Cmd, "Trig 1 Command").Layout),
								layout.Rigid(layout.Spacer{Height: 2}.Layout),
								layout.Rigid(material.Editor(th, &trig1Val, "Trig 1 Time (1-60)").Layout),
								layout.Rigid(layout.Spacer{Height: 4}.Layout),
								layout.Rigid(material.Editor(th, &trig2Cmd, "Trig 2 Command").Layout),
								layout.Rigid(layout.Spacer{Height: 2}.Layout),
								layout.Rigid(material.Editor(th, &trig2Val, "Trig 2 Time (1-60)").Layout),
								layout.Rigid(layout.Spacer{Height: 4}.Layout),
								layout.Rigid(material.Editor(th, &trig3Cmd, "Trig 3 Command").Layout),
								layout.Rigid(layout.Spacer{Height: 2}.Layout),
								layout.Rigid(material.Editor(th, &trig3Val, "Trig 3 Time (1-60)").Layout),
							)
						}
					})
				}),
			)

			e.Frame(gtx.Ops)
		}
	}
}

func drawCard(gtx layout.Context, bg color.NRGBA, w layout.Widget) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	dims := w(gtx)
	call := macro.Stop()
	rr := float32(gtx.Dp(unit.Dp(8)))
	paint.FillShape(gtx.Ops, bg, clip.UniformRRect(image.Rectangle{Max: dims.Size}, int(rr)).Op(gtx.Ops))
	call.Add(gtx.Ops)
	return dims
}

