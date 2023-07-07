package main

import (
	"log"
	"net/http"
	"time"

	"github.com/olahol/melody"

	"github.com/gridanias-helden/voidsent/internal/middleware"
	"github.com/gridanias-helden/voidsent/internal/models"

	_ "github.com/joho/godotenv/autoload"
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

	mel := melody.New()

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

	router := &http.ServeMux{}
	wrappedRouter := middleware.Chain(
		router,
		middleware.WithLogging,
		middleware.WithSession(memoryManager),
	)

	mel.HandleConnect(func(s *melody.Session) {
		session, ok := s.Request.Context().Value(middleware.SessionKey).(*models.Session)
		if !ok {
			s.Close()
			return
		}

		s.Set("session", session)
		log.Printf("connected: %s", s.Request.RemoteAddr)
	})

	mel.HandleMessage(func(s *melody.Session, msg []byte) {
		log.Printf("received: %s", string(msg))
	})

	router.Handle("/", http.FileServer(http.Dir(appConfig.Static)))
	router.HandleFunc("/auth/login/discord", discordOauth.Auth)
	router.HandleFunc("/auth/callback/discord", discordOauth.Callback)
	router.HandleFunc("/auth/logout", discordOauth.Logout)
	router.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		mel.HandleRequest(writer, request)
	})

	log.Printf("Listening on %s", appConfig.Bind)
	log.Fatal(http.ListenAndServe(appConfig.Bind, wrappedRouter))
}
