package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/mediocregopher/radix/v4"

	"github.com/olahol/melody"

	_ "github.com/joho/godotenv/autoload"
	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"

	"github.com/gridanias-helden/voidsent/internal/config"
	"github.com/gridanias-helden/voidsent/internal/handlers"
	"github.com/gridanias-helden/voidsent/internal/middleware"
	"github.com/gridanias-helden/voidsent/internal/services"
)

func main() {
	appConfig, err := config.New()
	if err != nil {
		log.Fatalln(err)
	}

	mel := melody.New()
	redisClient, err := (radix.PoolConfig{}).New(context.Background(), "tcp", appConfig.RedisHost)
	if err != nil {
		log.Fatalln("redis error", err)
	}

	//memoryManager := services.NewMemory()
	redisManager := services.NewRedis(redisClient)

	discordHandler := handlers.Discord{
		OAuth: &oauth2.Config{
			RedirectURL:  appConfig.RedirectURL,
			Scopes:       []string{discord.ScopeIdentify},
			Endpoint:     discord.Endpoint,
			ClientID:     appConfig.DiscordClientID,
			ClientSecret: appConfig.DiscordClientSecret,
		},
		Service: redisManager,
		KV:      make(map[string]time.Time),
	}
	wsHandler := &handlers.WebSocket{
		Melody:  mel,
		Service: redisManager,
	}

	router := &http.ServeMux{}
	wrappedRouter := middleware.Chain(
		router,
		middleware.WithLogging,
		middleware.WithSession(redisManager),
	)

	mel.HandleConnect(wsHandler.Connect)
	mel.HandleMessage(wsHandler.Message)
	mel.HandleMessageBinary(wsHandler.MessageBinary)

	router.Handle("/", http.FileServer(http.Dir(appConfig.Static)))
	router.HandleFunc("/auth/login/discord", discordHandler.Auth)
	router.HandleFunc("/auth/callback/discord", discordHandler.Callback)
	router.HandleFunc("/auth/logout", discordHandler.Logout)
	router.HandleFunc("/ws", wsHandler.HTTPRequest)

	log.Printf("Listening on %s", appConfig.Bind)
	log.Fatal(http.ListenAndServe(appConfig.Bind, wrappedRouter))
}
