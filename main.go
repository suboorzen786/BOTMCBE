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
Version string
Ping    string
Click   widget.Clickable
}

func main() {
go func() {
w := new(app.Window)
w.Option(app.Title("BOTMC - ChatCraft PRO Ultimate"))
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

servers := []ServerProfile{
{Name: "Hypixel Network", Address: "mc.hypixel.net", Port: "19132", Version: "26.52", Ping: "35 ms"},
{Name: "Rajasthan SMP", Address: "pe.rajasthansmp.fun", Port: "19132", Version: "26.52", Ping: "42 ms"},
}

var chatEditor widget.Editor
chatEditor.SingleLine = true
var sendChatBtn, addServerBtn, logoutBtn, confirmAddServerBtn, backToHomeBtn widget.Clickable

var nameInput, ipInput, portInput widget.Editor
nameInput.SingleLine = true
ipInput.SingleLine = true
portInput.SingleLine = true

var autoLoginToggle, cmdAfterLoginToggle, autoMoveToggle widget.Bool
var cmdAfterLoginInput widget.Editor
cmdAfterLoginInput.SingleLine = true

var trig1Text, trig1Time, trig2Text, trig2Time, trig3Text, trig3Time widget.Editor
trig1Text.SingleLine = true
trig1Time.SingleLine = true
trig2Text.SingleLine = true
trig2Time.SingleLine = true
trig3Text.SingleLine = true
trig3Time.SingleLine = true

currentView := 0 
var settingsBtn widget.Clickable
isAddingServer := false

var slideRightBtn, slideLeftBtn widget.Clickable

logMessages := []string{
"[System] Connected to server successfully.",
"Esegui il Login con /login <password>",
"Il login è avvenuto con successo!",
"Carra: Hello there!",
}

onlinePlayers := []string{
"ramusar", "SharliLP", "Roquexd", "Steeph", "Darcyy", "Rip__Hunter",
}

var list, playerList widget.List
list.Axis = layout.Vertical
playerList.Axis = layout.Vertical

bgColor := color.NRGBA{R: 20, G: 22, B: 28, A: 255}
headerColor := color.NRGBA{R: 190, G: 40, B: 40, A: 255}
cardColor := color.NRGBA{R: 30, G: 35, B: 45, A: 255}
inputBgColor := color.NRGBA{R: 15, G: 18, B: 24, A: 255}
textColor := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
subTextColor := color.NRGBA{R: 160, G: 170, B: 190, A: 255}
cyanAccent := color.NRGBA{R: 0, G: 240, B: 255, A: 255}

for {
switch e := w.Event().(type) {
case app.DestroyEvent:
return e.Err
case app.FrameEvent:
gtx := app.NewContext(&ops, e)

if settingsBtn.Clicked(gtx) { currentView = 4 }
if backToHomeBtn.Clicked(gtx) { currentView = 0; isAddingServer = false }

if addServerBtn.Clicked(gtx) {
isAddingServer = !isAddingServer
}

if confirmAddServerBtn.Clicked(gtx) {
if nameInput.Text() != "" && ipInput.Text() != "" {
portStr := portInput.Text()
if portStr == "" { portStr = "19132" }
servers = append(servers, ServerProfile{
Name:    nameInput.Text(),
Address: ipInput.Text(),
Port:    portStr,
Version: "26.52",
Ping:    "28 ms",
})
nameInput.SetText("")
ipInput.SetText("")
portInput.SetText("")
isAddingServer = false
}
}

for i := range servers {
if servers[i].Click.Clicked(gtx) {
currentView = 1 
logMessages = append([]string{"[Joined] Connected to " + servers[i].Name}, logMessages...)
}
}

if slideRightBtn.Clicked(gtx) {
if currentView >= 1 && currentView < 3 {
currentView++
}
}
if slideLeftBtn.Clicked(gtx) {
if currentView > 1 && currentView <= 3 {
currentView--
}
}

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
titleText := "BOTMC PRO - Home"
if currentView == 1 { titleText = "Chat Console (Slide ➡)" }
if currentView == 2 { titleText = "Minimap & Radar" }
if currentView == 3 { titleText = "Online Players" }
if currentView == 4 { titleText = "Advanced Settings" }

lbl := material.H6(th, titleText)
lbl.Color = textColor
return lbl.Layout(gtx)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
func() layout.Widget {
if currentView != 0 && currentView != 4 {
return func(gtx layout.Context) layout.Dimensions {
b := material.Button(th, &backToHomeBtn, "🏠")
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
case 0:
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 12, Bottom: 12, Left: 14, Right: 14}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(material.Body1(th, "👤 PlayerBot#786").Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Caption(th, "Microsoft Account Connected")
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
func() layout.Widget {
if !isAddingServer {
return func(gtx layout.Context) layout.Dimensions { return layout.Dimensions{} }
}
return func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 10, Bottom: 10, Left: 12, Right: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(material.Body2(th, "Add Server (Name, IP, Port)")),
layout.Rigid(layout.Spacer{Height: 4}.Layout),
layout.Rigid(material.Editor(th, &nameInput, "Server Name").Layout),
layout.Rigid(layout.Spacer{Height: 4}.Layout),
layout.Rigid(material.Editor(th, &ipInput, "IP Address").Layout),
layout.Rigid(layout.Spacer{Height: 4}.Layout),
layout.Rigid(material.Editor(th, &portInput, "Port (19132)").Layout),
layout.Rigid(layout.Spacer{Height: 6}.Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
b := material.Button(th, &confirmAddServerBtn, "✔ Save & Tick")
b.Background = color.NRGBA{R: 50, G: 200, B: 100, A: 255}
return b.Layout(gtx)
}),
layout.Rigid(layout.Spacer{Height: 6}.Layout),
)
})
})
}
}(),
layout.Rigid(layout.Spacer{Height: 4}.Layout),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return material.List(th, &list).Layout(gtx, len(servers), func(gtx layout.Context, index int) layout.Dimensions {
srv := &servers[index]
return layout.Inset{Bottom: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 10, Bottom: 10, Left: 12, Right: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
layout.Rigid(material.H6(th, srv.Name).Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body2(th, srv.Ping)
lbl.Color = cyanAccent
return lbl.Layout(gtx)
}),
)
}),
layout.Rigid(material.Caption(th, "IP: "+srv.Address+":"+srv.Port)),
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
return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
layout.Rigid(material.Caption(th, "● Online (Chat Console)")),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
layout.Rigid(material.Button(th, &slideRightBtn, "Slide Right ➡").Layout),
)
}),
)
})
})
}),
layout.Rigid(layout.Spacer{Height: 6}.Layout),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 8, Bottom: 8, Left: 8, Right: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return material.List(th, &list).Layout(gtx, len(logMessages), func(gtx layout.Context, index int) layout.Dimensions {
lbl := material.Body2(th, logMessages[index])
lbl.Color = textColor
return layout.Inset{Top: 3, Bottom: 3}.Layout(gtx, lbl.Layout)
})
})
})
}),
layout.Rigid(layout.Spacer{Height: 6}.Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, inputBgColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 8, Bottom: 8, Left: 8, Right: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
ed := material.Editor(th, &chatEditor, "Type chat message or /command...")
ed.Color = textColor
return ed.Layout(gtx)
})
})
}),
layout.Rigid(layout.Spacer{Width: 6}.Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
b := material.Button(th, &sendChatBtn, "SEND")
b.Background = cyanAccent
b.Color = color.NRGBA{R: 20, G: 22, B: 28, A: 255}
return b.Layout(gtx)
}),
)
}),
)

case 2:
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 6, Bottom: 6, Left: 10, Right: 10}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
layout.Rigid(material.Button(th, &slideLeftBtn, "⬅ Slide Left").Layout),
layout.Rigid(material.Caption(th, "Slide Right -> Players")),
layout.Rigid(material.Button(th, &slideRightBtn, "Slide Right ➡").Layout),
)
})
})
}),
layout.Rigid(layout.Spacer{Height: 12}.Layout),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 20, Bottom: 20, Left: 16, Right: 16}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.H6(th, "[ Minimap & Radar Active ]")
lbl.Color = cyanAccent
return layout.Inset{Bottom: 8}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body1(th, "Coordinates: X=3745.94 Y=108.00 Z=-2366.58")
lbl.Color = textColor
return layout.Inset{Bottom: 6}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body2(th, "Server Ping: 35 ms")
lbl.Color = subTextColor
return lbl.Layout(gtx)
}),
)
})
})
}),
)

case 3:
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 6, Bottom: 6, Left: 10, Right: 10}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
layout.Rigid(material.Button(th, &slideLeftBtn, "⬅ Slide Left").Layout),
layout.Rigid(material.Caption(th, "Online Players List")),
layout.Rigid(func(gtx layout.Context) layout.Dimensions { return layout.Dimensions{} }),
)
})
})
}),
layout.Rigid(layout.Spacer{Height: 8}.Layout),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 8, Bottom: 8, Left: 8, Right: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return material.List(th, &playerList).Layout(gtx, len(onlinePlayers), func(gtx layout.Context, index int) layout.Dimensions {
lbl := material.Body1(th, "👤  "+onlinePlayers[index])
lbl.Color = textColor
return layout.Inset{Top: 6, Bottom: 6}.Layout(gtx, lbl.Layout)
})
})
})
}),
)

default:
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.H6(th, "BotMC PRO Advanced Settings")
lbl.Color = textColor
return layout.Inset{Bottom: 8}.Layout(gtx, lbl.Layout)
}),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 10, Bottom: 10, Left: 12, Right: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(material.CheckBox(th, &autoLoginToggle, "Auto Login (3s after join)")),
layout.Rigid(layout.Spacer{Height: 4}.Layout),
layout.Rigid(material.CheckBox(th, &cmdAfterLoginToggle, "Command After Login")),
layout.Rigid(layout.Spacer{Height: 2}.Layout),
layout.Rigid(material.Editor(th, &cmdAfterLoginInput, "Enter custom command (e.g. /lobby)")),
layout.Rigid(layout.Spacer{Height: 6}.Layout),
layout.Rigid(material.CheckBox(th, &autoMoveToggle, "Auto Movement Every Minute")),
layout.Rigid(layout.Spacer{Height: 8}.Layout),
layout.Rigid(material.Body2(th, "3 Repeat Text Triggers (Command + 1-60 Min/Sec):")),
layout.Rigid(layout.Spacer{Height: 4}.Layout),

layout.Rigid(material.Editor(th, &trig1Text, "Trigger 1 Command (e.g. /feed)")),
layout.Rigid(layout.Spacer{Height: 2}.Layout),
layout.Rigid(material.Editor(th, &trig1Time, "Trigger 1 Time (1-60 Min/Sec)")),
layout.Rigid(layout.Spacer{Height: 6}.Layout),

layout.Rigid(material.Editor(th, &trig2Text, "Trigger 2 Command")),
layout.Rigid(layout.Spacer{Height: 2}.Layout),
layout.Rigid(material.Editor(th, &trig2Time, "Trigger 2 Time (1-60 Min/Sec)")),
layout.Rigid(layout.Spacer{Height: 6}.Layout),

layout.Rigid(material.Editor(th, &trig3Text, "Trigger 3 Command")),
layout.Rigid(layout.Spacer{Height: 2}.Layout),
layout.Rigid(material.Editor(th, &trig3Time, "Trigger 3 Time (1-60 Min/Sec)")),
)
})
})
}),
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
paint.FillShape(gtx.Ops, bg,
clip.UniformRRect(image.Rect
