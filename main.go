package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"

	"github.com/gridanias-helden/voidsent/internal/config"
	"github.com/gridanias-helden/voidsent/internal/handlers"
	"github.com/gridanias-helden/voidsent/internal/services"
)

func main() {
	appConfig, err := config.New()
	if err != nil {
		log.Fatalln(err)
	}

	// yamlService, err := services.NewYAML("./store.yml")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	memoryManager := services.NewMemory()

	discordOauth := handlers.Discord{
		OAuth: &oauth2.Config{
			RedirectURL:  appConfig.RedirectURL,
			Scopes:       []string{discord.ScopeIdentify},
			Endpoint:     discord.Endpoint,
			ClientID:     appConfig.DiscordClientID,
			ClientSecret: appConfig.DiscordClientSecret,
		},
		Service: memoryManager,
		KV:      make(map[string]time.Time),
	}

	http.Handle("/", http.FileServer(http.Dir(appConfig.Static)))
	http.HandleFunc("/auth/login/discord", discordOauth.Auth)
	http.HandleFunc("/auth/callback/discord", discordOauth.Callback)

	log.Printf("Listening on %s", appConfig.Bind)
	log.Fatal(http.ListenAndServe(appConfig.Bind, nil))
}
