package main

import (
"fmt"
"image/color"
"log"
"os"

"gioui.org/app"
"gioui.org/font/gofont"
"gioui.org/layout"
"gioui.org/op"
"gioui.org/op/paint"
"gioui.org/text"
"gioui.org/widget"
"gioui.org/widget/material"
"gioui.org/f32"
"gioui.org/unit"

"botmc/bot"
)

func main() {
go func() {
w := new(app.Window)
w.Option(app.Title("BOTMC Pro"))
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
ipEditor.SingleLine = true
portEditor.SingleLine = true
userEditor.SingleLine = true

var chatEditor widget.Editor
chatEditor.SingleLine = true

var joinBtn, stopBtn, rejoinBtn widget.Clickable
var sendChatBtn widget.Clickable

logChan := make(chan string, 100)
var logMessages []string
var activeBot *bot.BotSession

var list widget.List
list.Axis = layout.Vertical

// Colors
bgColor := color.NRGBA{R: 18, G: 18, B: 24, A: 255}       // Dark background
cardColor := color.NRGBA{R: 30, G: 30, B: 42, A: 255}     // Card background
textColor := color.NRGBA{R: 240, G: 240, B: 245, A: 255}  // Bright text
subTextColor := color.NRGBA{R: 150, G: 150, B: 170, A: 255} // Dim text
accentColor := color.NRGBA{R: 90, G: 105, B: 255, A: 255} // Neon Blue/Indigo

for {
switch e := w.Event().(type) {
case app.DestroyEvent:
return e.Err
case app.FrameEvent:
gtx := app.NewContext(&ops, e)

// Handle logs
select {
case msg := <-logChan:
logMessages = append(logMessages, msg)
if len(logMessages) > 150 {
logMessages = logMessages[len(logMessages)-150:]
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

// Fill background
paint.Fill(gtx.Ops, bgColor)

layout.Flex{Axis: layout.Vertical}.Layout(gtx,
// Header
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 20, Bottom: 16, Left: 16, Right: 16}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
title := material.H5(th, "BOTMC Pro")
title.Color = textColor
return title.Layout(gtx)
})
}),

// Connection Card
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Left: 16, Right: 16, Bottom: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 12, Bottom: 12, Left: 12, Right: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Bottom: 8}.Layout(gtx, material.Body1(th, "Server IP").Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
ed := material.Editor(th, &ipEditor, "127.0.0.1")
ed.Color = textColor
return ed.Layout(gtx)
}),
layout.Rigid(layout.Spacer{Height: 8}.Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(material.Body1(th, "Port").Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
ed := material.Editor(th, &portEditor, "19132")
ed.Color = textColor
return ed.Layout(gtx)
}),
)
}),
layout.Rigid(layout.Spacer{Width: 12}.Layout),
layout.Flexed(2, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(material.Body1(th, "Bot Username").Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
ed := material.Editor(th, &userEditor, "Username")
ed.Color = textColor
return ed.Layout(gtx)
}),
)
}),
)
}),
)
})
})
})
}),

// Action Buttons
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Left: 16, Right: 16, Bottom: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
btn := material.Button(th, &joinBtn, "JOIN")
btn.Background = accentColor
return btn.Layout(gtx)
}),
layout.Rigid(layout.Spacer{Width: 8}.Layout),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
btn := material.Button(th, &stopBtn, "STOP")
btn.Background = color.NRGBA{R: 200, G: 60, B: 60, A: 255}
return btn.Layout(gtx)
}),
layout.Rigid(layout.Spacer{Width: 8}.Layout),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
btn := material.Button(th, &rejoinBtn, "REJOIN")
btn.Background = color.NRGBA{R: 60, G: 160, B: 80, A: 255}
return btn.Layout(gtx)
}),
)
})
}),

// Status & Coordinates Bar
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
coordsText := "Status: Offline"
if activeBot != nil && activeBot.IsConnected {
coordsText = fmt.Sprintf("Online | %s", activeBot.Coords)
}
return layout.Inset{Left: 20, Right: 20, Bottom: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
lbl := material.Body2(th, coordsText)
lbl.Color = subTextColor
return lbl.Layout(gtx)
})
}),

// Live Logs Console Card
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Left: 16, Right: 16, Bottom: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 8, Bottom: 8, Left: 10, Right: 10}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return material.List(th, &list).Layout(gtx, len(logMessages), func(gtx layout.Context, index int) layout.Dimensions {
lbl := material.Body2(th, logMessages[index])
lbl.Color = textColor
return layout.Inset{Top: 2, Bottom: 2}.Layout(gtx, lbl.Layout)
})
})
})
})
}),

// Chat / Command Input Bar
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Left: 16, Right: 16, Bottom: 20}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 10, Bottom: 10, Left: 12, Right: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
ed := material.Editor(th, &chatEditor, "Type command or chat...")
ed.Color = textColor
return ed.Layout(gtx)
})
})
}),
layout.Rigid(layout.Spacer{Width: 8}.Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
btn := material.Button(th, &sendChatBtn, "Send")
btn.Background = accentColor
return btn.Layout(gtx)
}),
)
})
}),
)

e.Frame(gtx.Ops)
}
}
}

// Helper to draw rounded cards
func drawCard(gtx layout.Context, bg color.NRGBA, w layout.Widget) layout.Dimensions {
return layout.Stack{}.Layout(gtx,
layout.Stacked(func(gtx layout.Context) layout.Dimensions {
rr := float32(gtx.Dp(unit.Dp(8)))
paint.FillShape(gtx.Ops, bg,
clip.RRect{
Rect: f32.Rectangle{Max: f32.Point{
X: float32(gtx.Constraints.Max.X),
Y: float32(gtx.Constraints.Max.Y),
}},
Radius: rr,
}.Op(gtx.Ops),
)
return layout.Dimensions{Size: gtx.Constraints.Max}
}),
layout.Stacked(w),
)
}
