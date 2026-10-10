package main

import (
"encoding/json"
"os"
)

type Config struct {
ServerIP    string `json:"server_ip"`
ServerPort  string `json:"server_port"`
BotUsername string `json:"bot_username"`
MCPEVersion string `json:"mcpe_version"`
}

func LoadConfig() Config {
file, err := os.ReadFile("config.json")
if err != nil {
return Config{
ServerIP:    "127.0.0.1",
ServerPort:  "19132",
BotUsername: "BOTMC_User",
MCPEVersion: "1.20.0",
}
}
var cfg Config
json.Unmarshal(file, &cfg)
return cfg
}

func (c *Config) Save() {
data, _ := json.MarshalIndent(c, "", "  ")
os.WriteFile("config.json", data, 0644)
}
