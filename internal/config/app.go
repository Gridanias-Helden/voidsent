package config

import (
	"fmt"
	"os"
	"strings"
)

type AppConfig struct {
	RedirectURL         string
	DiscordClientID     string
	DiscordClientSecret string
	Host                string
	Bind                string
	Static              string
}

func New() (*AppConfig, error) {
	conf := AppConfig{
		DiscordClientID:     strings.TrimSpace(os.Getenv("VOIDSENT_DISCORD_CLIENT_ID")),
		DiscordClientSecret: strings.TrimSpace(os.Getenv("VOIDSENT_DISCORD_CLIENT_SECRET")),
		Host:                strings.TrimSpace(os.Getenv("VOIDSENT_HOST")),
		Bind:                strings.TrimSpace(os.Getenv("VOIDSENT_BIND")),
		Static:              strings.TrimSpace(os.Getenv("VOIDSENT_DATA")),
	}

	if conf.Bind == "" {
		conf.Bind = ":3080"
	}

	if conf.Host == "" {
		conf.Host = "localhost"
	}

	if conf.Static == "" {
		conf.Static = "./static"
	}

	if conf.DiscordClientID == "" {
		return nil, fmt.Errorf("missing discord client id (VOIDSENT_DISCORD_CLIENT_ID)")
	}

	if conf.DiscordClientSecret == "" {
		return nil, fmt.Errorf("missing discord client secret (VOIDSENT_DISCORD_CLIENT_SECRET)")
	}

	return &conf, nil
}
