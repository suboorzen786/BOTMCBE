package bot

import (
"fmt"
"time"
)

type BotSession struct {
ServerIP    string
ServerPort  string
Username    string
Version     string
IsConnected bool
Coords      string
logChan     chan<- string
stopChan    chan struct{}
}

func NewBotSession(ip, port, username, version string, logChan chan<- string) *BotSession {
return &BotSession{
ServerIP:   ip,
ServerPort: port,
Username:   username,
Version:    version,
logChan:    logChan,
stopChan:   make(chan struct{}),
Coords:     "X: 0, Y: 0, Z: 0",
}
}

func (b *BotSession) Connect() {
b.logChan <- fmt.Sprintf("[*] Connecting to %s:%s as %s...", b.ServerIP, b.ServerPort, b.Username)

// Simulating connection handshake for Bedrock protocol
time.Sleep(1 * time.Second)

b.IsConnected = true
b.logChan <- "[+] Successfully connected to the MCPE server!"

// Simulated game loop / position tracker
ticker := time.NewTicker(2 * time.Second)
defer ticker.Stop()

// Dummy coordinates increment to test live UI updates
x, y, z := 100.0, 64.0, 200.0

for {
select {
case <-b.stopChan:
b.logChan <- "[-] Bot disconnected."
return
case <-ticker.C:
if !b.IsConnected {
return
}
x += 0.5
z += 0.5
b.Coords = fmt.Sprintf("X: %.1f, Y: %.1f, Z: %.1f", x, y, z)
}
}
}

func (b *BotSession) SendChat(message string) {
if !b.IsConnected {
b.logChan <- "[!] Cannot send chat: Bot is offline."
return
}
b.logChan <- fmt.Sprintf("<%s> %s", b.Username, message)
// Yahan par Bedrock chat packet dispatch hoga
}

func (b *BotSession) Disconnect() {
if b.IsConnected {
b.IsConnected = false
close(b.stopChan)
}
}
