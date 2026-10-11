package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type SavedServer struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    string `json:"port"`
}

type TriggerConfig struct {
	Command string `json:"command"`
	Value   int    `json:"value"`   // 1 to 60
	IsMin   bool   `json:"is_min"`  // true = min, false = sec
}

type AppConfig struct {
	IsLoggedIn      bool          `json:"is_logged_in"`
	MicrosoftUser   string        `json:"microsoft_user"`
	RefreshToken    string        `json:"refresh_token"`
	AutoLogin       bool          `json:"auto_login"`
	AutoLoginPass   string        `json:"auto_login_pass"`
	CmdAfterLogin   bool          `json:"cmd_after_login"`
	CmdAfterText    string        `json:"cmd_after_text"`
	AutoMovement    bool          `json:"auto_movement"`
	Triggers        [3]TriggerConfig `json:"triggers"`
	Servers         []SavedServer `json:"servers"`
}

func getConfigPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = "."
	}
	appDir := filepath.Join(dir, "botmc_pro")
	_ = os.MkdirAll(appDir, 0755)
	return filepath.Join(appDir, "config.json")
}

func LoadConfig() AppConfig {
	path := getConfigPath()
	file, err := os.ReadFile(path)
	if err != nil {
		// Default config if file doesn't exist
		return AppConfig{
			IsLoggedIn: false,
			Servers: []SavedServer{
				{Name: "Hypixel Network", Address: "mc.hypixel.net", Port: "19132"},
			},
		}
	}
	var cfg AppConfig
	_ = json.Unmarshal(file, &cfg)
	return cfg
}

func SaveConfig(cfg AppConfig) error {
	path := getConfigPath()
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

