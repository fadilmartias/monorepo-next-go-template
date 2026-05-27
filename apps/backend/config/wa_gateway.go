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

type DilztifyConfig struct {
	BaseURL  string
	APIKey   string
	DeviceID string
}

var (
	dilztifyConfig *DilztifyConfig
	dilztifyOnce   sync.Once
)

func LoadDilztifyConfig() *DilztifyConfig {
	dilztifyOnce.Do(func() {
		dilztifyConfig = &DilztifyConfig{
			BaseURL:  "https://api.whatsapp.dilztopup.com/v1/wa",
			APIKey:   os.Getenv("DILZTIFY_API_KEY"),
			DeviceID: os.Getenv("DILZTIFY_DEVICE_ID"),
		}
	})
	return dilztifyConfig
}
