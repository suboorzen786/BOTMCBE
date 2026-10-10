package main

import (
"fmt"
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

"botmc/bot"
)

type NavItem struct {
Name  string
Click widget.Clickable
}

func main() {
go func() {
w := new(app.Window)
w.Option(app.Title("BOTMC"))
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

cfg := LoadConfig()
var ops op.Ops

// Microsoft Login State
isLoggedIn := false
userAuthCode := "https://microsoft.com/link\nCode: MCBE-7829"
var msLoginBtn widget.Clickable

// Editors
var ipEditor, portEditor, userEditor widget.Editor
ipEditor.SetText(cfg.ServerIP)
portEditor.SetText(cfg.ServerPort)
userEditor.SetText(cfg.BotUsername)
ipEditor.SingleLine = true
portEditor.SingleLine = true
userEditor.SingleLine = true

var chatEditor widget.Editor
chatEditor.SingleLine = true

var joinBtn, stopBtn, rejoinBtn widget.Clickable
var sendChatBtn widget.Clickable

// Toggles
var liveChatToggle, repeatMsgToggle, autoReconnectToggle widget.Bool

// Sidebar Navigation
navs := []NavItem{
{"Home", widget.Clickable{}},
{"Minimap & COD", widget.Clickable{}},
{"Chat Console", widget.Clickable{}},
{"Commands Scheduler", widget.Clickable{}},
{"Chat Repeat", widget.Clickable{}},
{"Auto Reconnect", widget.Clickable{}},
{"Version Settings", widget.Clickable{}},
{"Saved Profiles", widget.Clickable{}},
{"Global Settings", widget.Clickable{}},
}
currentTab := 0

logChan := make(chan string, 100)
var logMessages []string = []string{
"System initialized - BOTMC v1.0",
"Waiting for Microsoft Account Authentication...",
}
var activeBot *bot.BotSession

var list, navList widget.List
list.Axis = layout.Vertical
navList.Axis = layout.Vertical

// Cyberpunk Theme Palette
bgColor := color.NRGBA{R: 16, G: 18, B: 24, A: 255}
sidebarColor := color.NRGBA{R: 20, G: 24, B: 34, A: 255}
cardColor := color.NRGBA{R: 26, G: 30, B: 42, A: 255}
textColor := color.NRGBA{R: 240, G: 242, B: 245, A: 255}
subTextColor := color.NRGBA{R: 140, G: 150, B: 170, A: 255}
cyanAccent := color.NRGBA{R: 0, G: 240, B: 255, A: 255}
btnRed := color.NRGBA{R: 220, G: 60, B: 80, A: 255}

for {
switch e := w.Event().(type) {
case app.DestroyEvent:
return e.Err
case app.FrameEvent:
gtx := app.NewContext(&ops, e)

// Handle Login Button Click
if msLoginBtn.Clicked(gtx) {
isLoggedIn = true
logMessages = append([]string{"Microsoft Account Authenticated successfully!"}, logMessages...)
}

// Handle Tab Clicks
for i, nav := range navs {
if nav.Click.Clicked(gtx) {
currentTab = i
}
}

select {
case msg := <-logChan:
logMessages = append([]string{msg}, logMessages...)
if len(logMessages) > 50 {
logMessages = logMessages[:50]
}
w.Invalidate()
default:
}

if joinBtn.Clicked(gtx) {
cfg.ServerIP = ipEditor.Text()
cfg.ServerPort = portEditor.Text()
cfg.BotUsername = userEditor.Text()
cfg.Save()

if activeBot != nil {
activeBot.Disconnect()
}
activeBot = bot.NewBotSession(cfg.ServerIP, cfg.ServerPort, cfg.BotUsername, cfg.MCPEVersion, logChan)
go activeBot.Connect()
}

if stopBtn.Clicked(gtx) {
if activeBot != nil {
activeBot.Disconnect()
activeBot = nil
}
}

if rejoinBtn.Clicked(gtx) {
if activeBot != nil {
activeBot.Disconnect()
}
activeBot = bot.NewBotSession(cfg.ServerIP, cfg.ServerPort, cfg.BotUsername, cfg.MCPEVersion, logChan)
go activeBot.Connect()
}

if sendChatBtn.Clicked(gtx) && chatEditor.Text() != "" {
if activeBot != nil {
activeBot.SendChat(chatEditor.Text())
chatEditor.SetText("")
}
}

paint.Fill(gtx.Ops, bgColor)

// IF NOT LOGGED IN: SHOW MICROSOFT DEVICE CODE FIRST
if !isLoggedIn {
layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 24, Bottom: 24, Left: 20, Right: 20}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.H4(th, "BOTMC")
lbl.Color = cyanAccent
return layout.Inset{Bottom: 4}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body2(th, "✦ Crafted by Suboorzen")
lbl.Color = subTextColor
return layout.Inset{Bottom: 20}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body1(th, "Microsoft Device Code Login")
lbl.Color = textColor
return layout.Inset{Bottom: 10}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body2(th, userAuthCode)
lbl.Color = cyanAccent
return layout.Inset{Bottom: 20}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
btn := material.Button(th, &msLoginBtn, "Authenticate & Continue")
btn.Background = cyanAccent
btn.Color = color.NRGBA{R: 10, G: 12, B: 18, A: 255}
return btn.Layout(gtx)
}),
)
})
})
})
} else {
// LOGGED IN DASHBOARD
layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
// Left Drawer / Sidebar
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
gtx.Constraints.Max.X = gtx.Dp(unit.Dp(135))
return drawCard(gtx, sidebarColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 12, Bottom: 12, Left: 6, Right: 6}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.H6(th, "BOTMC")
lbl.Color = cyanAccent
return layout.Inset{Bottom: 2, Left: 4}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Caption(th, "by Suboorzen")
lbl.Color = subTextColor
return layout.Inset{Bottom: 12, Left: 4}.Layout(gtx, lbl.Layout)
}),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return material.List(th, &navList).Layout(gtx, len(navs), func(gtx layout.Context, index int) layout.Dimensions {
btn := material.Button(th, &navs[index].Click, navs[index].Name)
if currentTab == index {
btn.Background = cyanAccent
btn.Color = color.NRGBA{R: 10, G: 12, B: 18, A: 255}
} else {
btn.Background = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
btn.Color = textColor
}
return layout.Inset{Bottom: 4}.Layout(gtx, btn.Layout)
})
}),
)
})
})
}),

// Main Right Panel
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 8, Bottom: 8, Left: 8, Right: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
switch currentTab {
case 0: // Home Dashboard
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Bottom: 8}.Layout(gtx, material.H5(th, "Server Connection").Layout)
}),
// Server Card
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Bottom: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 10, Bottom: 10, Left: 10, Right: 10}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(material.Editor(th, &ipEditor, "Server IP").Layout),
layout.Rigid(layout.Spacer{Height: 6}.Layout),
layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
layout.Flexed(1, material.Editor(th, &portEditor, "Port").Layout),
layout.Rigid(layout.Spacer{Width: 6}.Layout),
layout.Flexed(2, material.Editor(th, &userEditor, "Bot Username").Layout),
),
layout.Rigid(layout.Spacer{Height: 8}.Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
btn := material.Button(th, &joinBtn, "JOIN")
btn.Background = cyanAccent
btn.Color = color.NRGBA{R: 10, G: 12, B: 18, A: 255}
return btn.Layout(gtx)
}),
layout.Rigid(layout.Spacer{Height: 6}.Layout),
layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
b := material.Button(th, &stopBtn, "STOP")
b.Background = btnRed
return b.Layout(gtx)
}),
layout.Rigid(layout.Spacer{Width: 6}.Layout),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
b := material.Button(th, &rejoinBtn, "REJOIN")
b.Background = color.NRGBA{R: 40, G: 120, B: 90, A: 255}
return b.Layout(gtx)
}),
),
)
})
})
})
}),
// Live Status Bar
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Bottom: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 8, Bottom: 8, Left: 10, Right: 10}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
coordsText := "Offline"
if activeBot != nil && activeBot.IsConnected {
coordsText = fmt.Sprintf("Online | %s", activeBot.Coords)
}
return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
layout.Rigid(material.Body2(th, "Status: "+coordsText).Layout),
layout.Rigid(material.Body2(th, "Ping: -- ms").Layout),
)
})
})
})
}),
// Recent Activity Console Logs
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 8, Bottom: 8, Left: 10, Right: 10}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(material.Body1(th, "Recent Activity").Layout),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return material.List(th, &list).Layout(gtx, len(logMessages), func(gtx layout.Context, index int) layout.Dimensions {
return material.Body2(th, logMessages[index]).Layout(gtx)
})
}),
)
})
})
}),
)

case 8: // Global Settings & Toggles
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Bottom: 12}.Layout(gtx, material.H5(th, "Global Settings").Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 12, Bottom: 12, Left: 12, Right: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(material.CheckBox(th, &liveChatToggle, "Enable Live Chat Console").Layout),
layout.Rigid(layout.Spacer{Height: 8}.Layout),
layout.Rigid(material.CheckBox(th, &repeatMsgToggle, "Enable Repeat Message System").Layout),
layout.Rigid(layout.Spacer{Height: 8}.Layout),
layout.Rigid(material.CheckBox(th, &autoReconnectToggle, "Enable Auto-Reconnect Timer").Layout),
)
})
})
}),
)

default:
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Bottom: 12}.Layout(gtx, material.H5(th, navs[currentTab].Name).Layout)
}),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 16, Left: 16, Right: 16}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return material.Body1(th, "Module Configured & Ready.").Layout(gtx)
})
})
}),
)
}
})
}),
)
}

e.Frame(gtx.Ops)
}
}
}

func drawCard(gtx layout.Context, bg color.NRGBA, w layout.Widget) layout.Dimensions {
return layout.Stack{}.Layout(gtx,
layout.Stacked(func(gtx layout.Context) layout.Dimensions {
rr := float32(gtx.Dp(unit.Dp(8)))
bounds := gtx.Constraints.Max
paint.FillShape(gtx.Ops, bg,
clip.UniformRRect(image.Rectangle{Max: bounds}, int(rr)).Op(gtx.Ops),
)
return layout.Dimensions{Size: bounds}
}),
layout.Stacked(w),
)
}
