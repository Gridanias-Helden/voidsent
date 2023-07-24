package main

import (
	// "context"
	"log"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	// "github.com/mediocregopher/radix/v4"
	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"

	"github.com/gridanias-helden/voidsent/pkg/config"
	"github.com/gridanias-helden/voidsent/pkg/middleware"
	"github.com/gridanias-helden/voidsent/pkg/services"
	"github.com/gridanias-helden/voidsent/pkg/services/chat"
	"github.com/gridanias-helden/voidsent/pkg/services/session"
	ws "github.com/gridanias-helden/voidsent/pkg/services/websocket"
	"github.com/gridanias-helden/voidsent/pkg/storage/memory"
	// "github.com/gridanias-helden/voidsent/pkg/storage/redis"
)

func main() {
	appConfig, err := config.New()
	if err != nil {
		log.Fatalln(err)
	}

	// redisClient, err := (radix.PoolConfig{}).New(context.Background(), "tcp", appConfig.RedisHost)
	// if err != nil {
	//	log.Fatalln("redis error", err)
	// }

	// memoryManager := services.NewMemory()
	// redisManager := redis.NewPlayers(redisClient)
	sessionService := memory.NewSessions(24 * time.Hour)
	broker := services.NewBroker()
	broker.AddService("chat", chat.New(broker))

	discordHandler := &session.Discord{
		OAuth: &oauth2.Config{
			RedirectURL:  appConfig.RedirectURL,
			Scopes:       []string{discord.ScopeIdentify},
			Endpoint:     discord.Endpoint,
			ClientID:     appConfig.DiscordClientID,
			ClientSecret: appConfig.DiscordClientSecret,
		},
		KV: make(map[string]time.Time),
		// Players:  redisManager,
		Sessions: sessionService,
	}
	wsHandler := &ws.WebSocket{
		Sessons: sessionService,
		Broker:  broker,
	}
	guestHandler := &session.GuestLogin{
		Sessions: sessionService,
	}

	router := &http.ServeMux{}
	wrappedRouter := middleware.Chain(
		router,
		middleware.WithLogging,
		middleware.WithSession(sessionService),
	)

	router.Handle("/", http.FileServer(http.Dir(appConfig.Static)))
	router.HandleFunc("/auth/login/discord", discordHandler.Auth)
	router.HandleFunc("/auth/callback/discord", discordHandler.Callback)
	router.HandleFunc("/auth/login/guest", guestHandler.Register)
	router.HandleFunc("/auth/logout", discordHandler.Logout)
	router.HandleFunc("/ws", wsHandler.HTTPRequest)

	log.Printf("Listening on %s", appConfig.Bind)
	log.Fatal(http.ListenAndServe(appConfig.Bind, wrappedRouter))
}
