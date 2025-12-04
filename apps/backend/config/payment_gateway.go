package config

import (
	"os"
	"strings"
	"sync"
)

type MidtransConfig struct {
	Env          string
	SecretKey    string
	PublicKey    string
	BaseURL      string
	SnapAPIURL   string
	WebhookToken string
	MerchantID   string
}

type IPaymuConfig struct {
	Env     string
	VA      string
	APIKey  string
	BaseURL string
}

var (
	midtransConfig *MidtransConfig
	ipaymuConfig   *IPaymuConfig
	midtransOnce   sync.Once
	ipaymuOnce     sync.Once
)

func LoadMidtransConfig() *MidtransConfig {
	midtransOnce.Do(func() {
		midtransEnv := os.Getenv("MIDTRANS_ENV")
		secretKey := os.Getenv("MIDTRANS_SECRET_KEY_" + strings.ToUpper(midtransEnv))
		publicKey := os.Getenv("MIDTRANS_PUBLIC_KEY_" + strings.ToUpper(midtransEnv))
		baseURL := os.Getenv("MIDTRANS_BASE_URL_" + strings.ToUpper(midtransEnv))
		snapAPIURL := os.Getenv("MIDTRANS_SNAP_API_URL_" + strings.ToUpper(midtransEnv))
		midtransConfig = &MidtransConfig{
			Env:        midtransEnv,
			SecretKey:  secretKey,
			PublicKey:  publicKey,
			BaseURL:    baseURL,
			SnapAPIURL: snapAPIURL,
			MerchantID: os.Getenv("MIDTRANS_MERCHANT_ID"),
		}
	})
	return midtransConfig
}

func LoadIPaymuConfig() *IPaymuConfig {
	ipaymuOnce.Do(func() {
		ipaymuEnv := os.Getenv("IPAYMU_ENV")
		VA := os.Getenv("IPAYMU_VA_" + strings.ToUpper(ipaymuEnv))
		APIKey := os.Getenv("IPAYMU_API_KEY_" + strings.ToUpper(ipaymuEnv))
		baseURL := os.Getenv("IPAYMU_BASE_URL_" + strings.ToUpper(ipaymuEnv))
		ipaymuConfig = &IPaymuConfig{
			Env:     ipaymuEnv,
			VA:      VA,
			APIKey:  APIKey,
			BaseURL: baseURL,
		}
	})
	return ipaymuConfig
}
