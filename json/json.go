package json

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Title  string `json:"title"`
	Width  string `json:"width"`
	Height string `json:"height"`
	Rd     string `json:"rd"`
	Cx     string `json:"cx"`
	Cy     string `json:"cy"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func SaveConfig(filename string, cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

var Cfg Config

func IniziaJson() {
	// Lettura
	cfg, err := LoadConfig("config.json")
	if err != nil {
		fmt.Println("Errore lettura:", err)
		return
	}
	Cfg = *cfg

	fmt.Println("Config letta:\n", cfg, "\n ", Cfg)

	// // Modifica
	// cfg.Title = "Nuovo titolo"
	// // cfg.Width = 1024

	// // Scrittura
	// if err := SaveConfig("config.json", cfg); err != nil {
	// 	fmt.Println("Errore scrittura:", err)
	// 	return
	// }

	// fmt.Println("Config salvata.")
}
