package main

import (
"image/color"
"os"
"strings"

"gioui.org/app"
"gioui.org/layout"
"gioui.org/op"
"gioui.org/op/clip"
"gioui.org/op/paint"
"gioui.org/unit"
"gioui.org/widget"
"gioui.org/widget/material"

"botmc/bot"
)

func main() {
go func() {
w := new(app.Window)
w.Option(app.Title("BOTMC"))
if err := loop(w); err != nil {
os.Exit(1)
}
os.Exit(0)
}()
app.Main()
}

func loop(w *app.Window) error {
th := material.NewTheme()
th.Palette.Bg = color.NRGBA{R: 18, G: 18, B: 24, A: 255}
th.Palette.Fg = color.NRGBA{R: 235, G: 235, B: 240, A: 255}
th.Palette.ContrastBg = color.NRGBA{R: 76, G: 175, B: 80, A: 255}

var (
ip     widget.Editor
port   widget.Editor
user   widget.Editor
chat   widget.Editor
start  widget.Clickable
stop   widget.Clickable
send   widget.Clickable
reconn widget.Bool
afk    widget.Bool
list   widget.List
ops    op.Ops
logs   []string
state  = bot.Disconnected
b      = bot.New()
)
ip.SingleLine, port.SingleLine, user.SingleLine, chat.SingleLine = true, true, true, true
port.SetText("19132")
user.SetText("BOTMC")
list.Axis = layout.Vertical
reconn.Value, afk.Value = true, true

// Bot events -> UI.
go func() {
for ev := range b.Events() {
state = ev.State
if ev.Log != "" {
logs = append(logs, ev.Log)
if len(logs) > 500 {
logs = logs[len(logs)-500:]
}
}
w.Invalidate()
}
}()

for {
switch e := w.Event().(type) {
case app.DestroyEvent:
b.Stop()
return e.Err
case app.FrameEvent:
gtx := app.NewContext(&ops, e)

if start.Clicked(gtx) {
b.Start(bot.Config{
Host:          strings.TrimSpace(ip.Text()),
Port:          strings.TrimSpace(port.Text()),
Username:      strings.TrimSpace(user.Text()),
AutoReconnect: reconn.Value,
AntiAFK:       afk.Value,
})
}
if stop.Clicked(gtx) {
b.Stop()
}
if send.Clicked(gtx) && strings.TrimSpace(chat.Text()) != "" {
if err := b.SendChat(chat.Text()); err != nil {
logs = append(logs, "Error: "+err.Error())
} else {
logs = append(logs, "> "+chat.Text())
chat.SetText("")
}
}

paintRect(gtx, th.Palette.Bg)
layout.UniformInset(unit.Dp(14)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
row(material.H5(th, "BOTMC").Layout),
row(material.Caption(th, "Created by Suboorzen").Layout),
row(field(th, &ip, "Server IP / Hostname")),
row(field(th, &port, "Port")),
row(field(th, &user, "Bot username")),
row(material.CheckBox(th, &reconn, "Auto-reconnect").Layout),
row(material.CheckBox(th, &afk, "Anti-AFK").Layout),
row(func(gtx layout.Context) layout.Dimensions {
return layout.Flex{Spacing: layout.SpaceEvenly}.Layout(gtx,
layout.Flexed(1, material.Button(th, &start, "Connect").Layout),
layout.Rigid(layout.Spacer{Width: 10}.Layout),
layout.Flexed(1, material.Button(th, &stop, "Disconnect").Layout),
)
}),
row(material.Body1(th, "Status: "+state.String()).Layout),
layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
return material.List(th, &list).Layout(gtx, len(logs), func(gtx layout.Context, i int) layout.Dimensions {
return material.Body2(th, logs[i]).Layout(gtx)
})
}),
row(func(gtx layout.Context) layout.Dimensions {
return layout.Flex{}.Layout(gtx,
layout.Flexed(1, field(th, &chat, "Chat / command")),
layout.Rigid(layout.Spacer{Width: 10}.Layout),
layout.Rigid(material.Button(th, &send, "Send").Layout),
)
}),
)
})
e.Frame(gtx.Ops)
}
}
}

func row(w layout.Widget) layout.FlexChild {
return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, w)
})
}

func field(th *material.Theme, ed *widget.Editor, hint string) layout.Widget {
return func(gtx layout.Context) layout.Dimensions {
return material.Editor(th, ed, hint).Layout(gtx)
}
}

func paintRect(gtx layout.Context, c color.NRGBA) {
defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
paint.Fill(gtx.Ops, c)
}
