package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/auth"
	"golang.org/x/oauth2"
)

type BotRequest struct {
	SessionID string `json:"sessionId"`
	Server    struct {
		IP   string `json:"ip"`
		Port int    `json:"port"`
	} `json:"server"`
	UseMicrosoft bool `json:"useMicrosoft"`
}

var (
	activeBots = make(map[string]*minecraft.Conn)
	mutex      sync.Mutex
)

func main() {
	fmt.Println("BOTMC-GitHubAPK Engine Started")

	http.HandleFunc("/api/bot/start", startBotHandler)
	http.HandleFunc("/api/bot/stop", stopBotHandler)

	fmt.Println("Server listening on port :8080...")
	http.ListenAndServe(":8080", nil)
}

func startBotHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var req BotRequest
	json.NewDecoder(r.Body).Decode(&req)

	go func() {
		var token *oauth2.Token
		var err error

		tokenFile := fmt.Sprintf("token_%s.json", req.SessionID)

		if _, statErr := os.Stat(tokenFile); statErr == nil {
			fileData, readErr := os.ReadFile(tokenFile)
			if readErr == nil {
				var t oauth2.Token
				if json.Unmarshal(fileData, &t) == nil {
					token = &t
				}
			}
		}

		if token == nil && req.UseMicrosoft {
			token, err = auth.RequestToken(func(url string) {
				fmt.Printf("\n[MICROSOFT LOGIN REQUIRED FOR SESSION: %s]\nOpen this link: %s\n\n", req.SessionID, url)
			})
			if err != nil {
				return
			}

			if data, marshalErr := json.Marshal(token); marshalErr == nil {
				_ = os.WriteFile(tokenFile, data, 0600)
			}
		}

		runUserBot(req, token)
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":    "Bot process initiated",
		"sessionId": req.SessionID,
	})
}

func stopBotHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("sessionId")
	mutex.Lock()
	if conn, exists := activeBots[sessionID]; exists {
		conn.Close()
		delete(activeBots, sessionID)
	}
	mutex.Unlock()
	w.Write([]byte(`{"status": "Bot stopped"}`))
}

func runUserBot(req BotRequest, token *oauth2.Token) {
	address := fmt.Sprintf("%s:%d", req.Server.IP, req.Server.Port)

	dialer := minecraft.Dialer{
		Token: token,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	conn, err := dialer.DialContext(ctx, "raknet", address)
	if err != nil {
		return
	}

	mutex.Lock()
	activeBots[req.SessionID] = conn
	mutex.Unlock()

	for {
		_, err := conn.ReadPacket()
		if err != nil {
			break
		}
	}
}

