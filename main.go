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
	Click   widget.Clickable
}

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("BOTMC PRO - Test"))
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
		}
	}

	var msLoginBtn, logoutBtn, addServerBtn, confirmAddBtn, backHomeBtn widget.Clickable
	var sendChatBtn widget.Clickable
	var isAddingServer bool

	var nameInput, ipInput, portInput widget.Editor
	nameInput.SingleLine = true
	ipInput.SingleLine = true
	portInput.SingleLine = true

	var chatEditor widget.Editor
	chatEditor.SingleLine = true

	currentView := 0 // 0: Home/Login, 1: Server Active Session
	if !cfg.IsLoggedIn {
		currentView = -1 // Login Screen
	}

	logMessages := []string{"[System] Connected to server successfully."}
	var list widget.List
	list.Axis = layout.Vertical

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

			// Login Action
			if msLoginBtn.Clicked(gtx) {
				cfg.IsLoggedIn = true
				cfg.MicrosoftUser = "PlayerBot#786"
				_ = SaveConfig(cfg)
				currentView = 0
			}

			// Logout Action
			if logoutBtn.Clicked(gtx) {
				cfg.IsLoggedIn = false
				cfg.MicrosoftUser = ""
				_ = SaveConfig(cfg)
				currentView = -1
			}

			if backHomeBtn.Clicked(gtx) {
				currentView = 0
				isAddingServer = false
			}

			if addServerBtn.Clicked(gtx) {
				isAddingServer = !isAddingServer
			}

			if confirmAddBtn.Clicked(gtx) && nameInput.Text() != "" && ipInput.Text() != "" {
				p := portInput.Text()
				if p == "" {
					p = "19132"
				}
				cfg.Servers = append(cfg.Servers, SavedServer{Name: nameInput.Text(), Address: ipInput.Text(), Port: p})
				_ = SaveConfig(cfg)
				servers = append(servers, ServerProfile{Name: nameInput.Text(), Address: ipInput.Text(), Port: p})
				nameInput.SetText("")
				ipInput.SetText("")
				portInput.SetText("")
				isAddingServer = false
			}

			for i := range servers {
				if servers[i].Click.Clicked(gtx) {
					currentView = 1 // Enter Active Server Chat
					logMessages = append([]string{"[Joined] Connected to " + servers[i].Name}, logMessages...)
				}
			}

			if sendChatBtn.Clicked(gtx) && chatEditor.Text() != "" {
				logMessages = append([]string{"You: " + chatEditor.Text()}, logMessages...)
				chatEditor.SetText("")
			}

			paint.Fill(gtx.Ops, bgColor)

			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				// Header
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return drawCard(gtx, headerColor, func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: 12, Bottom: 12, Left: 16, Right: 16}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									t := "BOTMC PRO - Test"
									if currentView == -1 {
										t = "Microsoft Login"
									}
									if currentView == 1 {
										t = "Server Chat Console"
									}
									lbl := material.H6(th, t)
									lbl.Color = textColor
									return lbl.Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if currentView == 1 {
										return material.Button(th, &backHomeBtn, "🏠").Layout(gtx)
									}
									return layout.Dimensions{}
								}),
							)
						})
					})
				}),

				// Body Content
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Top: 8, Bottom: 8, Left: 8, Right: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						switch currentView {
						case -1: // LOGIN SCREEN
							return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(layout.Spacer{Height: 60}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
										return layout.Inset{Top: 24, Bottom: 24, Left: 20, Right: 20}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
												layout.Rigid(material.H6(th, "Microsoft Auth Required").Layout),
												layout.Rigid(layout.Spacer{Height: 12}.Layout),
												layout.Rigid(material.Body2(th, "Sign in to manage and join servers.").Layout),
												layout.Rigid(layout.Spacer{Height: 20}.Layout),
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

						case 0: // HOME & SERVER LIST
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								// Profile Card
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
										return layout.Inset{Top: 12, Bottom: 12, Left: 14, Right: 14}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
												layout.Rigid(material.Body1(th, "👤 "+cfg.MicrosoftUser).Layout),
												layout.Rigid(material.Button(th, &logoutBtn, "Logout").Layout),
											)
										})
									})
								}),
								layout.Rigid(layout.Spacer{Height: 10}.Layout),
								// Servers Header & Add Button
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(material.Body1(th, "Saved Servers").Layout),
										layout.Rigid(material.Button(th, &addServerBtn, " + Add Server ").Layout),
									)
								}),
								layout.Rigid(layout.Spacer{Height: 8}.Layout),
								// Add Server Form Popup
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if !isAddingServer {
										return layout.Dimensions{}
									}
									return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
										return layout.Inset{Top: 12, Bottom: 12, Left: 12, Right: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
												layout.Rigid(material.Editor(th, &nameInput, "Server Name").Layout),
												layout.Rigid(layout.Spacer{Height: 6}.Layout),
												layout.Rigid(material.Editor(th, &ipInput, "IP Address").Layout),
												layout.Rigid(layout.Spacer{Height: 6}.Layout),
												layout.Rigid(material.Editor(th, &portInput, "Port (default 19132)").Layout),
												layout.Rigid(layout.Spacer{Height: 10}.Layout),
												layout.Rigid(material.Button(th, &confirmAddBtn, "Save Server").Layout),
											)
										})
									})
								}),
								layout.Rigid(layout.Spacer{Height: 6}.Layout),
								// Server List View
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									if len(servers) == 0 {
										return layout.Center.Layout(gtx, material.Body2(th, "No servers added yet. Click '+' above.").Layout)
									}
									return material.List(th, &list).Layout(gtx, len(servers), func(gtx layout.Context, index int) layout.Dimensions {
										srv := &servers[index]
										return layout.Inset{Bottom: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
												return layout.Inset{Top: 12, Bottom: 12, Left: 14, Right: 14}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
													return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
														layout.Rigid(material.H6(th, srv.Name).Layout),
														layout.Rigid(material.Caption(th, "IP: "+srv.Address+":"+srv.Port).Layout),
														layout.Rigid(layout.Spacer{Height: 8}.Layout),
														layout.Rigid(material.Button(th, &srv.Click, "JOIN SERVER").Layout),
													)
												})
											})
										})
									})
								}),
							)

						default: // ACTIVE SERVER CHAT VIEW
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(material.Caption(th, "● Connected to Server Console").Layout),
								layout.Rigid(layout.Spacer{Height: 6}.Layout),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
										return layout.Inset{Top: 8, Bottom: 8, Left: 8, Right: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											return material.List(th, &list).Layout(gtx, len(logMessages), func(gtx layout.Context, idx int) layout.Dimensions {
												return layout.Inset{Top: 4, Bottom: 4}.Layout(gtx, material.Body2(th, logMessages[idx]).Layout)
											})
										})
									})
								}),
								layout.Rigid(layout.Spacer{Height: 6}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
										layout.Flexed(1, material.Editor(th, &chatEditor, "Type message...").Layout),
										layout.Rigid(layout.Spacer{Width: 6}.Layout),
										layout.Rigid(material.Button(th, &sendChatBtn, "Send").Layout),
									)
								}),
							)
						}
					})
				}),
			)
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
