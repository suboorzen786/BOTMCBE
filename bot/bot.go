package bot

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sandertv/gophertunnel/minecraft"
)

type BotSession struct {
	ServerIP   string
	ServerPort string
	Username   string
	Version    string

	IsConnected bool
	Coords      string

	logChan chan<- string

	mu       sync.Mutex
	conn     *minecraft.Conn
	cancel   context.CancelFunc
	stopOnce sync.Once
}

func NewBotSession(ip, port, username, version string, logChan chan<- string) *BotSession {
	return &BotSession{
		ServerIP:   ip,
		ServerPort: port,
		Username:   username,
		Version:    version,
		Coords:     "Position unavailable",
		logChan:    logChan,
	}
}

func (b *BotSession) log(message string) {
	if b.logChan == nil {
		return
	}
	select {
	case b.logChan <- message:
	default:
	}
}

func (b *BotSession) Connect() {
	b.mu.Lock()
	if b.conn != nil {
		b.mu.Unlock()
		b.log("[!] A connection is already active.")
		return
	}

	port, err := strconv.Atoi(b.ServerPort)
	if err != nil || port < 1 || port > 65535 {
		b.mu.Unlock()
		b.log("[!] Invalid server port.")
		return
	}

	host := strings.TrimSpace(b.ServerIP)
	if host == "" {
		b.mu.Unlock()
		b.log("[!] Enter a server address.")
		return
	}

	address := net.JoinHostPort(host, b.ServerPort)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	b.cancel = cancel
	b.mu.Unlock()

	b.log("[*] Starting real Bedrock protocol connection to " + address)
	b.log("[!] Microsoft authentication is not configured yet. This connection does not claim to be Microsoft-authenticated.")

	dialer := minecraft.Dialer{}
	conn, err := dialer.DialContext(ctx, "raknet", address)
	cancel()

	if err != nil {
		b.mu.Lock()
		b.cancel = nil
		b.mu.Unlock()
		b.log(fmt.Sprintf("[-] Bedrock connection failed: %v", err))
		return
	}

	if err := conn.DoSpawn(); err != nil {
		_ = conn.Close()
		b.log(fmt.Sprintf("[-] Could not enter the world: %v", err))
		return
	}

	b.mu.Lock()
	b.conn = conn
	b.IsConnected = true
	b.cancel = nil
	b.mu.Unlock()

	b.log("[+] Bedrock session established. Authentication acceptance depends on the server.")

	for {
		_, err := conn.ReadPacket()
		if err != nil {
			b.mu.Lock()
			if b.conn == conn {
				b.conn = nil
				b.IsConnected = false
			}
			b.mu.Unlock()

			_ = conn.Close()
			b.log(fmt.Sprintf("[-] Bedrock session ended: %v", err))
			return
		}
	}
}

func (b *BotSession) SendChat(message string) {
	b.mu.Lock()
	connected := b.conn != nil && b.IsConnected
	b.mu.Unlock()

	if !connected {
		b.log("[!] Chat unavailable: bot is not connected.")
		return
	}

	b.log("[!] Chat packet sending is not implemented yet.")
}

func (b *BotSession) Disconnect() {
	b.stopOnce.Do(func() {
		b.mu.Lock()
		cancel := b.cancel
		conn := b.conn
		b.conn = nil
		b.IsConnected = false
		b.mu.Unlock()

		if cancel != nil {
			cancel()
		}
		if conn != nil {
			_ = conn.Close()
			b.log("[-] Bedrock connection closed.")
		}
	})
}
