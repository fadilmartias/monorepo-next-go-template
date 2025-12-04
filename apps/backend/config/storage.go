package config

import (
	"os"
	"sync"
)

type CloudflareR2Config struct {
	Token      string
	BucketName string
	AccountID  string
	AccessKey  string
	SecretKey  string
}

var (
	cloudflareR2Config *CloudflareR2Config
	cloudflareR2Once   sync.Once
)

func LoadCloudflareR2Config() *CloudflareR2Config {
	cloudflareR2Once.Do(func() {
		cloudflareR2Config = &CloudflareR2Config{
			Token:      os.Getenv("CLOUDFLARE_R2_TOKEN"),
			BucketName: os.Getenv("CLOUDFLARE_R2_BUCKET_NAME"),
			AccountID:  os.Getenv("CLOUDFLARE_R2_ACCOUNT_ID"),
			AccessKey:  os.Getenv("CLOUDFLARE_R2_ACCESS_KEY"),
			SecretKey:  os.Getenv("CLOUDFLARE_R2_SECRET_KEY"),
		}
	})
	return cloudflareR2Config
}
