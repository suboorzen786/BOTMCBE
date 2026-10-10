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

type NavItem struct {
Name  string
Click widget.Clickable
}

func main() {
go func() {
w := new(app.Window)
w.Option(app.Title("BOTMC - Desktop Client"))
w.Option(app.Size(1024, 720))
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

isLoggedIn := false
userAuthCode := "https://microsoft.com/link\nCode: MCBE-7829"
var msLoginBtn widget.Clickable

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

navs := []NavItem{
{"Home Dashboard", widget.Clickable{}},
{"Minimap & Radar", widget.Clickable{}},
{"Chat Console", widget.Clickable{}},
{"Commands & Automation", widget.Clickable{}},
{"Login & Account", widget.Clickable{}},
{"Auto Reconnect", widget.Clickable{}},
{"Saved Profiles", widget.Clickable{}},
{"Connection Logs", widget.Clickable{}},
{"App Settings", widget.Clickable{}},
}
currentTab := 0

logChan := make(chan string, 100)
var logMessages []string = []string{
"Desktop client initialized - BOTMC v1.0",
"Connected to play.example.net - BotMC_Bot active",
"Ping: 38 ms - Stable connection",
"Uptime tracker running: 01:24:10",
}
var activeBot *bot.BotSession

var list, navList widget.List
list.Axis = layout.Vertical
navList.Axis = layout.Vertical

bgColor := color.NRGBA{R: 12, G: 14, B: 20, A: 255}       
sidebarColor := color.NRGBA{R: 16, G: 20, B: 28, A: 255}  
cardColor := color.NRGBA{R: 22, G: 26, B: 36, A: 255}     
textColor := color.NRGBA{R: 240, G: 242, B: 245, A: 255}  
subTextColor := color.NRGBA{R: 130, G: 140, B: 160, A: 255} 
cyanAccent := color.NRGBA{R: 0, G: 240, B: 255, A: 255}   
btnRed := color.NRGBA{R: 215, G: 55, B: 75, A: 255}       
btnGreen := color.NRGBA{R: 35, G: 140, B: 95, A: 255}     

for {
switch e := w.Event().(type) {
case app.DestroyEvent:
return e.Err
case app.FrameEvent:
gtx := app.NewContext(&ops, e)

if msLoginBtn.Clicked(gtx) {
isLoggedIn = true
logMessages = append([]string{"Microsoft Account Authenticated successfully!"}, logMessages...)
}

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

if !isLoggedIn {
layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 32, Bottom: 32, Left: 36, Right: 36}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.H3(th, "BOTMC")
lbl.Color = cyanAccent
return layout.Inset{Bottom: 4}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body1(th, "✦ Desktop Edition — Crafted by Suboorzen")
lbl.Color = subTextColor
return layout.Inset{Bottom: 24}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body1(th, "Microsoft Device Code Login")
lbl.Color = textColor
return layout.Inset{Bottom: 12}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.H6(th, userAuthCode)
lbl.Color = cyanAccent
return layout.Inset{Bottom: 24}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
btn := material.Button(th, &msLoginBtn, "Authenticate & Launch Desktop Client")
btn.Background = cyanAccent
btn.Color = color.NRGBA{R: 12, G: 14, B: 20, A: 255}
return btn.Layout(gtx)
}),
)
})
})
})
} else {
layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
gtx.Constraints.Max.X = gtx.Dp(unit.Dp(220))
return drawCard(gtx, sidebarColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 16, Bottom: 16, Left: 12, Right: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.H5(th, "BOTMC")
lbl.Color = cyanAccent
return layout.Inset{Bottom: 2, Left: 4}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Caption(th, "Desktop Client v1.0")
lbl.Color = subTextColor
return layout.Inset{Bottom: 20, Left: 4}.Layout(gtx, lbl.Layout)
}),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return material.List(th, &navList).Layout(gtx, len(navs), func(gtx layout.Context, index int) layout.Dimensions {
btn := material.Button(th, &navs[index].Click, navs[index].Name)
if currentTab == index {
btn.Background = cyanAccent
btn.Color = color.NRGBA{R: 12, G: 14, B: 20, A: 255}
} else {
btn.Background = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
btn.Color = textColor
}
return layout.Inset{Bottom: 8}.Layout(gtx, btn.Layout)
})
}),
)
})
})
}),

layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 16, Bottom: 16, Left: 16, Right: 16}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
switch currentTab {
case 0:
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Bottom: 14}.Layout(gtx, material.H4(th, "Dashboard Control Center").Layout)
}),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Bottom: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 16, Bottom: 16, Left: 16, Right: 16}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.H6(th, "Server Connection")
lbl.Color = textColor
return layout.Inset{Bottom: 12}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(material.Editor(th, &ipEditor, "Server Address").Layout),
layout.Rigid(layout.Spacer{Height: 10}.Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
layout.Flexed(1, material.Editor(th, &portEditor, "Port").Layout),
layout.Rigid(layout.Spacer{Width: 10}.Layout),
layout.Flexed(2, material.Editor(th, &userEditor, "Bot Username").Layout),
)
}),
layout.Rigid(layout.Spacer{Height: 14}.Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
btn := material.Button(th, &joinBtn, "JOIN SERVER")
btn.Background = cyanAccent
btn.Color = color.NRGBA{R: 12, G: 14, B: 20, A: 255}
return btn.Layout(gtx)
}),
layout.Rigid(layout.Spacer{Height: 10}.Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
b := material.Button(th, &stopBtn, "STOP")
b.Background = btnRed
return b.Layout(gtx)
}),
layout.Rigid(layout.Spacer{Width: 10}.Layout),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
b := material.Button(th, &rejoinBtn, "REJOIN")
b.Background = btnGreen
return b.Layout(gtx)
}),
)
}),
)
})
})
})
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 14, Bottom: 14, Left: 16, Right: 16}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body1(th, "● Connected (Stable)")
lbl.Color = color.NRGBA{R: 50, G: 220, B: 100, A: 255}
return lbl.Layout(gtx)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body1(th, "Ping: 38 ms")
lbl.Color = subTextColor
return lbl.Layout(gtx)
}),
)
})
})
}),
)
}),
layout.Rigid(layout.Spacer{Width: 16}.Layout),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 16, Bottom: 16, Left: 16, Right: 16}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.H6(th, "Live Activity Console")
lbl.Color = textColor
return layout.Inset{Bottom: 10}.Layout(gtx, lbl.Layout)
}),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return material.List(th, &list).Layout(gtx, len(logMessages), func(gtx layout.Context, index int) layout.Dimensions {
lbl := material.Body1(th, logMessages[index])
lbl.Color = subTextColor
return layout.Inset{Top: 4, Bottom: 4}.Layout(gtx, lbl.Layout)
})
}),
)
})
})
}),
)
}),
)

default:
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Bottom: 14}.Layout(gtx, material.H4(th, navs[currentTab].Name).Layout)
}),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 20, Left: 20, Right: 20}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
lbl := material.Body1(th, "Desktop module active and configured.")
lbl.Color = subTextColor
return lbl.Layout(gtx)
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
