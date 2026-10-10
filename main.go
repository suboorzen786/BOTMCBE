package main

import (
"fmt"
"log"
"os"

"gioui.org/app"
"gioui.org/font/gofont"
"gioui.org/io/system"
"gioui.org/layout"
"gioui.org/op"
"gioui.org/widget"
"gioui.org/widget/material"

"botmc/bot"
)

func main() {
go func() {
w := new(app.Window)
w.Option(app.Title("BOTMC"))
w.Option(app.Size(400, 800))
if err := run(w); err != nil {
log.Fatal(err)
}
os.Exit(0)
}()
app.Main()
}

func run(w *app.Window) error {
th := material.NewTheme()
th.Shaper = gofont.Collection()

cfg := LoadConfig()

var ops op.Ops

// Form Controls
var ipEditor, portEditor, userEditor widget.Editor
ipEditor.SetText(cfg.ServerIP)
portEditor.SetText(cfg.ServerPort)
userEditor.SetText(cfg.BotUsername)

var chatEditor widget.Editor

// Action Buttons
var joinBtn, stopBtn, rejoinBtn widget.Clickable
var sendChatBtn widget.Clickable

logChan := make(chan string, 100)
var activeBot *bot.BotSession

for {
switch e := w.Event().(type) {
case system.DestroyEvent:
return e.Err
case system.FrameEvent:
gtx := layout.NewContext(&ops, e)

// Button Handlers
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

// Render Screen Layout
layout.Flex{Axis: layout.Vertical}.Layout(gtx,
// Header Title
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 16, Bottom: 16, Left: 16}.Layout(gtx,
material.H5(th, "BOTMC Dashboard").Layout,
)
}),

// Inputs Form
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Left: 16, Right: 16}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(material.Editor(th, &ipEditor, "Server Address").Layout),
layout.Rigid(material.Editor(th, &portEditor, "Port").Layout),
layout.Rigid(material.Editor(th, &userEditor, "Bot Username").Layout),
)
})
}),

// Action Buttons (JOIN, STOP, REJOIN)
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 10, Bottom: 10, Left: 16, Right: 16}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
layout.Flexed(1, material.Button(th, &joinBtn, "JOIN").Layout),
layout.Flexed(1, material.Button(th, &stopBtn, "STOP").Layout),
layout.Flexed(1, material.Button(th, &rejoinBtn, "REJOIN").Layout),
)
})
}),

// Coordinates & Status Overlay
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
coordsText := "Bot Offline"
if activeBot != nil && activeBot.IsConnected {
coordsText = fmt.Sprintf("Position: %s", activeBot.Coords)
}
return layout.Inset{Left: 16, Top: 8, Bottom: 8}.Layout(gtx,
material.Body2(th, coordsText).Layout,
)
}),

// Chat Console Input Box
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Left: 16, Right: 16, Top: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
layout.Flexed(1, material.Editor(th, &chatEditor, "Type message / command...").Layout),
layout.Rigid(material.Button(th, &sendChatBtn, "Send").Layout),
)
})
}),

// Footer Watermark
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return layout.Align(layout.South).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Bottom: 16}.Layout(gtx,
material.Caption(th, "✦ Crafted by Suboorzen").Layout,
)
})
}),
)

e.Frame(gtx.Ops)
}
}
}
