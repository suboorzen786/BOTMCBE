package main

import (
"fmt"
"log"
"os"

"gioui.org/app"
"gioui.org/font/gofont"
"gioui.org/layout"
"gioui.org/op"
"gioui.org/text"
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
th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

cfg := LoadConfig()

var ops op.Ops

var ipEditor, portEditor, userEditor widget.Editor
ipEditor.SetText(cfg.ServerIP)
portEditor.SetText(cfg.ServerPort)
userEditor.SetText(cfg.BotUsername)

var chatEditor widget.Editor
var joinBtn, stopBtn, rejoinBtn widget.Clickable
var sendChatBtn widget.Clickable

logChan := make(chan string, 100)
var logMessages []string
var activeBot *bot.BotSession

var list widget.List
list.Axis = layout.Vertical

for {
switch e := w.Event().(type) {
case app.DestroyEvent:
return e.Err
case app.FrameEvent:
gtx := app.NewContext(&ops, e)

// Non-blocking read from log channel
select {
case msg := <-logChan:
logMessages = append(logMessages, msg)
if len(logMessages) > 100 {
logMessages = logMessages[len(logMessages)-100:]
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

layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 16, Bottom: 12, Left: 16}.Layout(gtx,
material.H5(th, "BOTMC Dashboard").Layout,
)
}),

layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Left: 16, Right: 16, Bottom: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(material.Editor(th, &ipEditor, "Server IP").Layout),
layout.Rigid(material.Editor(th, &portEditor, "Port").Layout),
layout.Rigid(material.Editor(th, &userEditor, "Username").Layout),
)
})
}),

layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 6, Bottom: 6, Left: 16, Right: 16}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
layout.Flexed(1, material.Button(th, &joinBtn, "JOIN").Layout),
layout.Flexed(1, material.Button(th, &stopBtn, "STOP").Layout),
layout.Flexed(1, material.Button(th, &rejoinBtn, "REJOIN").Layout),
)
})
}),

layout.Rigid(func(gtx layout.Context) layout.Dimensions {
coordsText := "Status: Offline"
if activeBot != nil && activeBot.IsConnected {
coordsText = fmt.Sprintf("Live %s", activeBot.Coords)
}
return layout.Inset{Left: 16, Top: 8, Bottom: 8}.Layout(gtx,
material.Body2(th, coordsText).Layout,
)
}),

// Live Logs Console Area
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Left: 16, Right: 16, Bottom: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return material.List(th, &list).Layout(gtx, len(logMessages), func(gtx layout.Context, index int) layout.Dimensions {
return material.Body2(th, logMessages[index]).Layout(gtx)
})
})
}),

layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Left: 16, Right: 16, Bottom: 16}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
layout.Flexed(1, material.Editor(th, &chatEditor, "Chat / Command...").Layout),
layout.Rigid(material.Button(th, &sendChatBtn, "Send").Layout),
)
})
}),
)

e.Frame(gtx.Ops)
}
}
}
