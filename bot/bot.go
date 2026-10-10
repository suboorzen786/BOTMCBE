package bot

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// State is the connection state shown in the UI.
type State int

const (
	Disconnected State = iota
	Connecting
	Connected
	Error
)

func (s State) String() string {
	switch s {
	case Connecting:
		return "Connecting..."
	case Connected:
		return "Connected"
	case Error:
		return "Error"
	}
	return "Disconnected"
}

// Event is sent from the bot to the UI.
type Event struct {
	State State
	Log   string // non-empty => append to console
}

// Config holds the user settings.
type Config struct {
	Host          string
	Port          string
	Username      string
	AutoReconnect bool
	AntiAFK       bool
}

// Bot manages one Bedrock session.
type Bot struct {
	mu     sync.Mutex
	cancel context.CancelFunc
	conn   *minecraft.Conn
	events chan Event
}

func New() *Bot { return &Bot{events: make(chan Event, 256)} }

// Events returns the channel the UI should read.
func (b *Bot) Events() <-chan Event { return b.events }

func (b *Bot) emit(s State, log string) {
	select {
	case b.events <- Event{State: s, Log: log}:
	default:
	}
}

// Start begins a session in the background.
func (b *Bot) Start(cfg Config) {
	b.Stop()
	ctx, cancel := context.WithCancel(context.Background())
	b.mu.Lock()
	b.cancel = cancel
	b.mu.Unlock()
	go b.run(ctx, cfg)
}

// Stop disconnects and disables reconnecting.
func (b *Bot) Stop() {
	b.mu.Lock()
	cancel, conn := b.cancel, b.conn
	b.cancel, b.conn = nil, nil
	b.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	if conn != nil {
		_ = conn.Close()
	}
}

// SendChat sends a chat message / command as the bot.
func (b *Bot) SendChat(msg string) error {
	b.mu.Lock()
	conn := b.conn
	b.mu.Unlock()
	if conn == nil {
		return fmt.Errorf("not connected")
	}
	return conn.WritePacket(&packet.Text{
		TextType:   packet.TextTypeChat,
		SourceName: conn.IdentityData().DisplayName,
		Message:    msg,
		XUID:       conn.IdentityData().XUID,
	})
}

func (b *Bot) run(ctx context.Context, cfg Config) {
	for {
		err := b.session(ctx, cfg)
		if ctx.Err() != nil {
			b.emit(Disconnected, "Disconnected.")
			return
		}
		if err != nil {
			b.emit(Error, "Error: "+err.Error())
		}
		if !cfg.AutoReconnect {
			b.emit(Disconnected, "Session ended.")
			return
		}
		b.emit(Connecting, "Reconnecting in 5s...")
		select {
		case <-ctx.Done():
			b.emit(Disconnected, "Disconnected.")
			return
		case <-time.After(5 * time.Second):
		}
	}
}

func (b *Bot) session(ctx context.Context, cfg Config) error {
	addr := strings.TrimSpace(cfg.Host) + ":" + strings.TrimSpace(cfg.Port)
	b.emit(Connecting, "Connecting to "+addr+" ...")

	dctx, dcancel := context.WithTimeout(ctx, 20*time.Second)
	defer dcancel()

	// Offline-style login (works on servers with xbox-auth disabled).
	// Online-mode servers need Microsoft auth, which is not added yet.
	d := minecraft.Dialer{
		IdentityData: login.IdentityData{DisplayName: cfg.Username},
	}
	conn, err := d.DialContext(dctx, "raknet", addr)
	if err != nil {
		return err
	}
	b.mu.Lock()
	b.conn = conn
	b.mu.Unlock()
	defer func() {
		_ = conn.Close()
		b.mu.Lock()
		if b.conn == conn {
			b.conn = nil
		}
		b.mu.Unlock()
	}()

	if err := conn.DoSpawnContext(dctx); err != nil {
		return err
	}
	b.emit(Connected, "Spawned in the world.")

	if cfg.AntiAFK {
		go b.antiAFK(ctx, conn)
	}

	for {
		pk, err := conn.ReadPacket()
		if err != nil {
			return err
		}
		switch p := pk.(type) {
		case *packet.Text:
			b.emit(Connected, fmt.Sprintf("<%s> %s", p.SourceName, p.Message))
		case *packet.Disconnect:
			return fmt.Errorf("kicked: %s", p.Message)
		}
	}
}

// antiAFK sends a jump action periodically so the server sees activity.
func (b *Bot) antiAFK(ctx context.Context, conn *minecraft.Conn) {
	t := time.NewTicker(25 * time.Second)
	defer t.Stop()
	id := conn.GameData().EntityRuntimeID
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if err := conn.WritePacket(&packet.PlayerAction{
				EntityRuntimeID: id,
				ActionType:      protocol.PlayerActionJump,
			}); err != nil {
				return
			}
		}
	}
}
