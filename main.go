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
w.Option(app.Title("BOTMC - ChatCraft PRO Edition"))
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
{Name: "Hypixel Network [Example]", Address: "mc.hypixel.net", Port: "19132", Version: "26.52", Ping: "35 ms"},
{Name: "Rajasthan SMP [Example]", Address: "pe.rajasthansmp.fun", Port: "19132", Version: "26.52", Ping: "42 ms"},
}

var chatEditor widget.Editor
chatEditor.SingleLine = true
var sendChatBtn, tabBtn, upBtn, downBtn widget.Clickable

currentTab := 0 // 0: Servers/Home, 1: Chat, 2: Players, 3: Movement, 4: Settings
var navServers, navChat, navPlayers, navMovement, navSettings widget.Clickable

logMessages := []string{
"[UUIDSpoofFix] You passed the UUID check with success!",
"Your UUID: b1158376-1eca-344b-b741-a7e52097f7f4",
"Connected to Hypixel Network - Bot active",
"Esegui il Login con /login <password>",
"Il login è avvenuto con successo!",
"Carra: Hello there!",
}

onlinePlayers := []string{
"ramusar", "SharliLP", "Roquexd", 
"Steeph", "Darcyy", "Rip__Hunter",
"Reyrey04", "Diaxhunter", "laugesfar",
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

if navServers.Clicked(gtx) { currentTab = 0 }
if navChat.Clicked(gtx) { currentTab = 1 }
if navPlayers.Clicked(gtx) { currentTab = 2 }
if navMovement.Clicked(gtx) { currentTab = 3 }
if navSettings.Clicked(gtx) { currentTab = 4 }

for i := range servers {
if servers[i].Click.Clicked(gtx) {
currentTab = 1
logMessages = append([]string{"Connecting to " + servers[i].Name + "..."}, logMessages...)
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
return layout.Inset{Top: 10, Bottom: 10, Left: 12, Right: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.H6(th, "BOTMC PRO")
lbl.Color = textColor
return lbl.Layout(gtx)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
b := material.Button(th, &navChat, "💬")
b.Background = color.NRGBA{}
return b.Layout(gtx)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
b := material.Button(th, &navPlayers, "👥")
b.Background = color.NRGBA{}
return b.Layout(gtx)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
b := material.Button(th, &navMovement, "📍")
b.Background = color.NRGBA{}
return b.Layout(gtx)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
b := material.Button(th, &navSettings, "⚙")
b.Background = color.NRGBA{}
return b.Layout(gtx)
}),
)
}),
)
})
})
}),

layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 8, Bottom: 8, Left: 8, Right: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
switch currentTab {
case 0:
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body1(th, "Accounts - Persistent Login")
lbl.Color = subTextColor
return layout.Inset{Bottom: 6}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 16, Bottom: 16, Left: 16, Right: 16}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.H4(th, "+")
lbl.Color = cyanAccent
return lbl.Layout(gtx)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Caption(th, "Microsoft Account Logged In")
lbl.Color = textColor
return lbl.Layout(gtx)
}),
)
})
})
}),
layout.Rigid(layout.Spacer{Height: 12}.Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body1(th, "Servers - Click server to join")
lbl.Color = subTextColor
return layout.Inset{Bottom: 6}.Layout(gtx, lbl.Layout)
}),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return material.List(th, &list).Layout(gtx, len(servers), func(gtx layout.Context, index int) layout.Dimensions {
srv := &servers[index]
return layout.Inset{Bottom: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 12, Bottom: 12, Left: 12, Right: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
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
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Caption(th, "IP: "+srv.Address+" | Version: "+srv.Version)
lbl.Color = subTextColor
return layout.Inset{Bottom: 8}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
btn := material.Button(th, &srv.Click, "CONNECT TO TEST")
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
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Caption(th, "● Online (Hypixel)")
lbl.Color = color.NRGBA{R: 50, G: 220, B: 100, A: 255}
return lbl.Layout(gtx)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Caption(th, "Ping: 35 ms | Use top/bottom nav")
lbl.Color = cyanAccent
return lbl.Layout(gtx)
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
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
b := material.Button(th, &tabBtn, "TAB")
b.Background = color.NRGBA{R: 50, G: 120, B: 200, A: 255}
return b.Layout(gtx)
}),
layout.Rigid(layout.Spacer{Width: 6}.Layout),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, inputBgColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 8, Bottom: 8, Left: 8, Right: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
ed := material.Editor(th, &chatEditor, "Send message...")
ed.Color = textColor
return ed.Layout(gtx)
})
})
}),
layout.Rigid(layout.Spacer{Width: 6}.Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
b := material.Button(th, &sendChatBtn, " > ")
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
lbl := material.H6(th, "Online players")
lbl.Color = textColor
return layout.Inset{Bottom: 8}.Layout(gtx, lbl.Layout)
}),
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

case 3:
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 20, Bottom: 20, Left: 16, Right: 16}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.H6(th, "[ Minimap & Radar Active ]")
lbl.Color = cyanAccent
return layout.Inset{Bottom: 8}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body1(th, "X=3745.94 Y=108.00 Z=-2366.58")
lbl.Color = textColor
return lbl.Layout(gtx)
}),
)
})
})
}),
layout.Rigid(layout.Spacer{Height: 12}.Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
b := material.Button(th, &downBtn, "DOWN")
b.Background = color.NRGBA{R: 70, G: 75, B: 85, A: 255}
return b.Layout(gtx)
}),
layout.Rigid(layout.Spacer{Width: 12}.Layout),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
b := material.Button(th, &upBtn, "UP")
b.Background = color.NRGBA{R: 70, G: 75, B: 85, A: 255}
return b.Layout(gtx)
}),
)
}),
)

default:
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.H6(th, "ChatCraft PRO Settings")
lbl.Color = textColor
return layout.Inset{Bottom: 10}.Layout(gtx, lbl.Layout)
}),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 12, Bottom: 12, Left: 12, Right: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body1(th, "✓ Automatically reconnect (Enabled)")
lbl.Color = textColor
return layout.Inset{Bottom: 10}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body1(th, "✓ Auto login (/login password)")
lbl.Color = textColor
return layout.Inset{Bottom: 10}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body1(th, "• Fixed Version: 26.52")
lbl.Color = cyanAccent
return layout.Inset{Bottom: 10}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body1(th, "• Recurrent Timer: Minutes & Seconds Dial")
lbl.Color = textColor
return lbl.Layout(gtx)
}),
)
})
})
}),
)
}
})
}),

layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 6, Bottom: 6, Left: 6, Right: 6}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceEvenly}.Layout(gtx,
layout.Rigid(material.Button(th, &navServers, "Servers").Layout),
layout.Rigid(material.Button(th, &navChat, "Chat").Layout),
layout.Rigid(material.Button(th, &navPlayers, "Players").Layout),
layout.Rigid(material.Button(th, &navMovement, "Movement").Layout),
layout.Rigid(material.Button(th, &navSettings, "Settings").Layout),
)
})
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
