package config

import (
	"os"
	"sync"
)

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
}

type DiscordConfig struct {
	ClientID     string
	ClientSecret string
}

type FacebookConfig struct {
	ClientID     string
	ClientSecret string
}

type SteamConfig struct {
	ApiKey string
}

type TwitchConfig struct {
	ClientID     string
	ClientSecret string
}

var (
	googleConfig *GoogleConfig
	googleOnce   sync.Once
)

func LoadGoogleConfig() *GoogleConfig {
	googleOnce.Do(func() {
		googleConfig = &GoogleConfig{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		}
	})
	return googleConfig
}

var (
	discordConfig *DiscordConfig
	discordOnce   sync.Once
)

func LoadDiscordConfig() *DiscordConfig {
	discordOnce.Do(func() {
		discordConfig = &DiscordConfig{
			ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
			ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		}
	})
	return discordConfig
}

var (
	facebookConfig *FacebookConfig
	facebookOnce   sync.Once
)

func LoadFacebookConfig() *FacebookConfig {
	facebookOnce.Do(func() {
		facebookConfig = &FacebookConfig{
			ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
			ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
		}
	})
	return facebookConfig
}

var (
	steamConfig *SteamConfig
	steamOnce   sync.Once
)

func LoadSteamConfig() *SteamConfig {
	steamOnce.Do(func() {
		steamConfig = &SteamConfig{
			ApiKey: os.Getenv("STEAM_API_KEY"),
		}
	})
	return steamConfig
}

var (
	twitchConfig *TwitchConfig
	twitchOnce   sync.Once
)

func LoadTwitchConfig() *TwitchConfig {
	twitchOnce.Do(func() {
		twitchConfig = &TwitchConfig{
			ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
			ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		}
	})
	return twitchConfig
}
