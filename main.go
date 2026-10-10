package main

import (
"fmt"
"image/color"
"os"

"gioui.org/app"
"gioui.org/font/gofont"
"gioui.org/io/system"
"gioui.org/layout"
"gioui.org/op"
"gioui.org/text"
"gioui.org/widget"
"gioui.org/widget/material"
)

type BotState struct {
IP        widget.Editor
Port      widget.Editor
BotName   widget.Editor
Status    string
IsRunning bool
StartBtn  widget.Clickable
}

func main() {
bot := &BotState{
Status: "Status: Disconnected",
}
bot.IP.SetText("play.nethergames.org")
bot.Port.SetText("19132")
bot.BotName.SetText("BOTMC_Player")

go func() {
w := app.NewWindow(
app.Title("BOTMC - Minecraft Bedrock Bot"),
app.Size(400, 700),
)
th := material.NewTheme(gofont.Collection())

var ops op.Ops
for {
e := <-w.Events()
switch e := e.(type) {
case system.DestroyEvent:
os.Exit(0)
case system.FrameEvent:
gtx := layout.NewContext(&ops, e)

if bot.StartBtn.Clicked() {
if !bot.IsRunning {
bot.IsRunning = true
bot.Status = fmt.Sprintf("Status: Connected to %s:%s", bot.IP.Text(), bot.Port.Text())
} else {
bot.IsRunning = false
bot.Status = "Status: Disconnected"
}
}

layout.Flex{
Axis:    layout.Vertical,
Spacing: layout.SpaceStart,
}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
h := material.H5(th, "BOTMC Controller")
h.Color = color.NRGBA{R: 76, G: 175, B: 80, A: 255}
return h.Layout(gtx)
}),
layout.Rigid(layout.Spacer{Height: 16}.Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return material.Editor(th, &bot.IP, "Server IP / Hostname").Layout(gtx)
}),
layout.Rigid(layout.Spacer{Height: 8}.Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return material.Editor(th, &bot.Port, "Port (Default 19132)").Layout(gtx)
}),
layout.Rigid(layout.Spacer{Height: 8}.Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return material.Editor(th, &bot.BotName, "Bot Username").Layout(gtx)
}),
layout.Rigid(layout.Spacer{Height: 16}.Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
btnText := "Start Bot"
if bot.IsRunning {
btnText = "Disconnect Bot"
}
return material.Button(th, &bot.StartBtn, btnText).Layout(gtx)
}),
layout.Rigid(layout.Spacer{Height: 16}.Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body1(th, bot.Status)
lbl.Alignment = text.Middle
return lbl.Layout(gtx)
}),
)

e.Frame(gtx.Ops)
}
}
}()
app.Main()
}
