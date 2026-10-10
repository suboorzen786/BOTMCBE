package main

import (
"encoding/json"
"os"
)

type AppConfig struct {
ServerIP          string `json:"server_ip"`
ServerPort        string `json:"server_port"`
BotUsername        string `json:"bot_username"`
MCPEVersion       string `json:"mcpe_version"`
AuthToken         string `json:"auth_token"`
AutoLoginEnabled  bool   `json:"auto_login_enabled"`
AutoLoginPassword string `json:"auto_login_password"`
AutoLoginDelaySec int    `json:"auto_login_delay_sec"`
RepeatCmdEnabled  bool   `json:"repeat_cmd_enabled"`
RepeatCommand     string `json:"repeat_command"`
RepeatIntervalSec int    `json:"repeat_interval_sec"`
ShowMinimap       bool   `json:"show_minimap"`
ShowCoords        bool   `json:"show_coords"`
}

func LoadConfig() AppConfig {
cfg := AppConfig{
ServerIP:          "play.example.net",
ServerPort:        "19132",
BotUsername:        "BOTMC_Bot",
MCPEVersion:       "1.21.70",
AutoLoginDelaySec: 5,
RepeatIntervalSec: 168, // Example: 2m 48s
ShowMinimap:       true,
ShowCoords:        true,
}
data, err := os.ReadFile("botmc_config.json")
if err == nil {
json.Unmarshal(data, &cfg)
}
return cfg
}

func (c *AppConfig) Save() {
data, _ := json.MarshalIndent(c, "", "  ")
os.WriteFile("botmc_config.json", data, 0644)
}
