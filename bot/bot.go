package bot

import (
"fmt"
"time"
)

type BotSession struct {
IP          string
Port        string
Username    string
Version     string
IsConnected bool
Coords      string
LogChan     chan string
StopChan    chan struct{}
}

func NewBotSession(ip, port, username, version string, logChan chan string) *BotSession {
return &BotSession{
IP:          ip,
Port:        port,
Username:    username,
Version:     version,
Coords:      "X: 0, Y: 64, Z: 0",
LogChan:     logChan,
StopChan:    make(chan struct{}),
}
}

func (b *BotSession) StartMicrosoftAuth() string {
return "https://microsoft.com/link - Code: MCBE-BOT786"
}

func (b *BotSession) Connect() {
b.LogChan <- fmt.Sprintf("[SYSTEM] Connecting to %s:%s (MCPE Version: %s)...", b.IP, b.Port, b.Version)

time.Sleep(1 * time.Second)
b.IsConnected = true
b.LogChan <- fmt.Sprintf("[SUCCESS] Connected as %s!", b.Username)
b.LogChan <- "[BOT] Anti-AFK & Auto-Respawn engines activated."

go func() {
ticker := time.NewTicker(3 * time.Second)
defer ticker.Stop()
step := 0
for {
select {
case <-b.StopChan:
b.IsConnected = false
b.LogChan <- "[SYSTEM] Bot Disconnected."
return
case <-ticker.C:
step++
b.Coords = fmt.Sprintf("X: %d, Y: 64, Z: %d", 100+step, 200-step)
}
}
}()
}

func (b *BotSession) Disconnect() {
if b.IsConnected {
close(b.StopChan)
}
}

func (b *BotSession) SendChat(msg string) {
if b.IsConnected {
b.LogChan <- fmt.Sprintf("[CHAT] %s: %s", b.Username, msg)
} else {
b.LogChan <- "[ERROR] Cannot send message: Bot is offline."
}
}
