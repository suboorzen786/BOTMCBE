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

"botmc/bot"
)

func main() {
go func() {
w := new(app.Window)
w.Option(app.Title("MCBE BOT"))
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
var logMessages []string = []string{
"Connected to play.example.net - BotMC_Bot joined",
"Ping: 42 ms - Connection is stable",
"Reconnected - Auto reconnect completed",
"Version set to 1.21.70 - Bedrock Edition",
}
var activeBot *bot.BotSession

var list widget.List
list.Axis = layout.Vertical

// Cyberpunk Minecraft Theme Colors
bgColor := color.NRGBA{R: 16, G: 18, B: 24, A: 255}       
cardColor := color.NRGBA{R: 24, G: 28, B: 38, A: 255}     
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

layout.Flex{Axis: layout.Vertical}.Layout(gtx,
// Header Title matching UI
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 16, Bottom: 12, Left: 16, Right: 16}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
title := material.H5(th, "MCBE BOT")
title.Color = textColor
return title.Layout(gtx)
})
}),

// Connection Card
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Left: 16, Right: 16, Bottom: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 14, Bottom: 14, Left: 14, Right: 14}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body1(th, "Server Connection")
lbl.Color = textColor
return layout.Inset{Bottom: 8}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
ed := material.Editor(th, &ipEditor, "Server Address")
ed.Color = textColor
return ed.Layout(gtx)
}),
layout.Rigid(layout.Spacer{Height: 8}.Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
ed := material.Editor(th, &portEditor, "Port")
ed.Color = textColor
return ed.Layout(gtx)
}),
layout.Rigid(layout.Spacer{Width: 10}.Layout),
layout.Flexed(2, func(gtx layout.Context) layout.Dimensions {
ed := material.Editor(th, &userEditor, "Bot Username")
ed.Color = textColor
return ed.Layout(gtx)
}),
)
}),
layout.Rigid(layout.Spacer{Height: 12}.Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
btn := material.Button(th, &joinBtn, "JOIN")
btn.Background = cyanAccent
btn.Color = color.NRGBA{R: 10, G: 12, B: 18, A: 255}
return btn.Layout(gtx)
}),
layout.Rigid(layout.Spacer{Height: 8}.Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
btn := material.Button(th, &stopBtn, "STOP")
btn.Background = btnRed
return btn.Layout(gtx)
}),
layout.Rigid(layout.Spacer{Width: 8}.Layout),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
btn := material.Button(th, &rejoinBtn, "REJOIN")
btn.Background = color.NRGBA{R: 40, G: 120, B: 90, A: 255}
return btn.Layout(gtx)
}),
)
}),
)
})
})
})
}),

// Status Card
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Left: 16, Right: 16, Bottom: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 12, Bottom: 12, Left: 14, Right: 14}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body2(th, "● Connected")
lbl.Color = color.NRGBA{R: 50, G: 220, B: 100, A: 255}
return lbl.Layout(gtx)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body2(th, "Ping: 42 ms")
lbl.Color = subTextColor
return lbl.Layout(gtx)
}),
)
})
})
})
}),

// Logs Card
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Left: 16, Right: 16, Bottom: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 10, Bottom: 10, Left: 12, Right: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body1(th, "Recent Activity")
lbl.Color = textColor
return layout.Inset{Bottom: 6}.Layout(gtx, lbl.Layout)
}),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return material.List(th, &list).Layout(gtx, len(logMessages), func(gtx layout.Context, index int) layout.Dimensions {
lbl := material.Body2(th, logMessages[index])
lbl.Color = subTextColor
return layout.Inset{Top: 3, Bottom: 3}.Layout(gtx, lbl.Layout)
})
}),
)
})
})
})
}),

// Chat Bar
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Left: 16, Right: 16, Bottom: 16}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
ed := material.Editor(th, &chatEditor, "Type command or chat...")
ed.Color = textColor
return ed.Layout(gtx)
}),
layout.Rigid(layout.Spacer{Width: 8}.Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
btn := material.Button(th, &sendChatBtn, "Send")
btn.Background = cyanAccent
btn.Color = color.NRGBA{R: 10, G: 12, B: 18, A: 255}
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

func drawCard(gtx layout.Context, bg color.NRGBA, w layout.Widget) layout.Dimensions {
return layout.Stack{}.Layout(gtx,
layout.Stacked(func(gtx layout.Context) layout.Dimensions {
rr := float32(gtx.Dp(unit.Dp(10)))
bounds := gtx.Constraints.Max
paint.FillShape(gtx.Ops, bg,
clip.UniformRRect(image.Rectangle{Max: bounds}, int(rr)).Op(gtx.Ops),
)
return layout.Dimensions{Size: bounds}
}),
layout.Stacked(w),
)
}
