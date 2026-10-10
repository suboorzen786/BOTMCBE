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
w.Option(app.Title("BOTMC - Desktop Edition"))
w.Option(app.Size(1280, 800))
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

var joinBtn, stopBtn, rejoinBtn widget.Clickable

navs := []NavItem{
{"Home Dashboard", widget.Clickable{}},
{"Minimap & Radar", widget.Clickable{}},
{"Chat Console", widget.Clickable{}},
{"Commands & Automation", widget.Clickable{}},
{"Login & Account", widget.Clickable{}},
{"Auto Reconnect", widget.Clickable{}},
{"Version Settings", widget.Clickable{}},
{"Saved Profiles", widget.Clickable{}},
{"Connection Logs", widget.Clickable{}},
{"App Settings", widget.Clickable{}},
}
currentTab := 0

logChan := make(chan string, 100)
var logMessages []string = []string{
"Desktop Client v1.0 initialized successfully",
"Connected to play.example.net - BotMC active",
"Ping: 35 ms - Connection stable",
"Uptime tracker: 02:15:40 active",
}
var activeBot *bot.BotSession

var list, navList widget.List
list.Axis = layout.Vertical
navList.Axis = layout.Vertical

bgColor := color.NRGBA{R: 10, G: 12, B: 18, A: 255}       
sidebarColor := color.NRGBA{R: 14, G: 18, B: 26, A: 255}  
cardColor := color.NRGBA{R: 20, G: 25, B: 38, A: 255}     
inputBgColor := color.NRGBA{R: 12, G: 16, B: 24, A: 255}  
textColor := color.NRGBA{R: 255, G: 255, B: 255, A: 255}  
subTextColor := color.NRGBA{R: 150, G: 165, B: 190, A: 255} 
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

paint.Fill(gtx.Ops, bgColor)

if !isLoggedIn {
layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 36, Bottom: 36, Left: 40, Right: 40}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.H3(th, "BOTMC")
lbl.Color = cyanAccent
return layout.Inset{Bottom: 4}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body1(th, "✦ Desktop Client Edition — Crafted by Suboorzen")
lbl.Color = subTextColor
return layout.Inset{Bottom: 28}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.H6(th, "Microsoft Device Code Sign-In")
lbl.Color = textColor
return layout.Inset{Bottom: 12}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body1(th, userAuthCode)
lbl.Color = cyanAccent
return layout.Inset{Bottom: 28}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
btn := material.Button(th, &msLoginBtn, "Authenticate & Launch Desktop Dashboard")
btn.Background = cyanAccent
btn.Color = color.NRGBA{R: 10, G: 12, B: 18, A: 255}
return btn.Layout(gtx)
}),
)
})
})
})
} else {
layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
gtx.Constraints.Max.X = gtx.Dp(unit.Dp(240))
return drawCard(gtx, sidebarColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 16, Bottom: 16, Left: 14, Right: 14}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.H5(th, "BOTMC")
lbl.Color = cyanAccent
return layout.Inset{Bottom: 2, Left: 4}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Caption(th, "Desktop Client v1.0")
lbl.Color = subTextColor
return layout.Inset{Bottom: 24, Left: 4}.Layout(gtx, lbl.Layout)
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
return layout.Inset{Bottom: 16}.Layout(gtx, material.H4(th, "Desktop Control Center").Layout)
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
return layout.Inset{Bottom: 4}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body2(th, "Configure your Minecraft Bedrock server parameters")
lbl.Color = subTextColor
return layout.Inset{Bottom: 14}.Layout(gtx, lbl.Layout)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, inputBgColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 8, Bottom: 8, Left: 10, Right: 10}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
ed := material.Editor(th, &ipEditor, "Server Address")
ed.Color = textColor
return ed.Layout(gtx)
})
})
}),
layout.Rigid(layout.Spacer{Height: 10}.Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, inputBgColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 8, Bottom: 8, Left: 10, Right: 10}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
ed := material.Editor(th, &portEditor, "Port")
ed.Color = textColor
return ed.Layout(gtx)
})
})
}),
layout.Rigid(layout.Spacer{Width: 10}.Layout),
layout.Flexed(2, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, inputBgColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 8, Bottom: 8, Left: 10, Right: 10}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
ed := material.Editor(th, &userEditor, "Bot Username")
ed.Color = textColor
return ed.Layout(gtx)
})
})
}),
)
}),
layout.Rigid(layout.Spacer{Height: 14}.Layout),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
btn := material.Button(th, &joinBtn, "JOIN SERVER")
btn.Background = cyanAccent
btn.Color = color.NRGBA{R: 10, G: 12, B: 18, A: 255}
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
lbl := material.Body1(th, "Uptime: 02:15:40   Ping: 35 ms")
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
return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.H6(th, "Live Activity Console")
lbl.Color = textColor
return lbl.Layout(gtx)
}),
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
lbl := material.Body2(th, "View All >")
lbl.Color = cyanAccent
return lbl.Layout(gtx)
}),
)
}),
layout.Rigid(layout.Spacer{Height: 12}.Layout),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return material.List(th, &list).Layout(gtx, len(logMessages), func(gtx layout.Context, index int) layout.Dimensions {
lbl := material.Body1(th, logMessages[index])
lbl.Color = subTextColor
return layout.Inset{Top: 6, Bottom: 6}.Layout(gtx, lbl.Layout)
})
}),
)
})
})
}),
)
}),

layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 10}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
lbl := material.Caption(th, "✦ Desktop Edition — Crafted by Suboorzen")
lbl.Color = subTextColor
return layout.Center.Layout(gtx, lbl.Layout)
})
}),
)

default:
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Bottom: 16}.Layout(gtx, material.H4(th, navs[currentTab].Name).Layout)
}),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return drawCard(gtx, cardColor, func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Top: 24, Left: 24, Right: 24}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
lbl := material.Body1(th, "Desktop module active and operational.")
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
