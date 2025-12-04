package config

import (
	"os"
	"sync"
)

type FonnteConfig struct {
	BaseURL string
	Token   string
}

var (
	fonnteConfig *FonnteConfig
	fonnteOnce   sync.Once
)

func LoadFonnteConfig() *FonnteConfig {
	fonnteOnce.Do(func() {
		fonnteConfig = &FonnteConfig{
			BaseURL: "https://api.fonnte.com",
			Token:   os.Getenv("FONNTE_TOKEN"),
		}
	})
	return fonnteConfig
}
