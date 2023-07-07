package handlers

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/oklog/ulid"
	"golang.org/x/oauth2"

	"github.com/gridanias-helden/voidsent/internal/models"
	"github.com/gridanias-helden/voidsent/internal/services"
)

type discordUser struct {
	ID            string  `json:"id"`
	Username      string  `json:"username"`
	Discriminator string  `json:"discriminator"`
	GlobalName    *string `json:"global_name"`
	Avatar        *string `json:"avatar"`
	Bot           *bool   `json:"bot"`
	System        *bool   `json:"system"`
	MFAEnabled    *bool   `json:"mfa_enabled"`
	Banner        *string `json:"banner"`
	AccentColor   *int    `json:"accent_color"`
	Locale        *string `json:"locale"`
	Flags         *int    `json:"flags"`
	PremiumType   int     `json:"premium_type"`
	PublicFlags   int     `json:"public_flags"`
}

type Discord struct {
	OAuth   *oauth2.Config
	Service services.Service
	KV      map[string]time.Time
}

func (d *Discord) Auth(w http.ResponseWriter, r *http.Request) {
	// Is user already logged in?
	c, err := r.Cookie("voidsent_session")
	if err == nil {
		_, err := d.Service.LoadSessionByID(r.Context(), c.Value)
		if err == nil {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
	}

	state, err := ulid.New(uint64(time.Now().UnixMilli()), rand.Reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	d.KV[state.String()] = time.Now()

	http.Redirect(w, r, d.OAuth.AuthCodeURL(state.String()), http.StatusTemporaryRedirect)
}

func (d *Discord) Callback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if timestamp, ok := d.KV[state]; !ok || time.Now().After(timestamp.Add(5*time.Minute)) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("State does not match or is expired."))
		return
	}

	// We exchange the code we got for an access token
	// Then we can use the access token to do actions, limited to scopes we requested
	token, err := d.OAuth.Exchange(context.Background(), r.FormValue("code"))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Use the access token, here we use it to get the logged in user's info.
	res, err := d.OAuth.Client(context.Background(), token).Get("https://discord.com/api/users/@me")

	if err != nil || res.StatusCode != 200 {
		w.WriteHeader(http.StatusInternalServerError)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte(res.Status))
		}
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var user discordUser
	if err := json.Unmarshal(body, &user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	player, err := d.Service.LoadPlayerByID(r.Context(), user.ID)
	if err != nil {
		var avatar string
		if user.Avatar != nil {
			avatar = *user.Avatar
		}
		player = &models.Player{
			ID:     user.ID,
			Name:   user.Username,
			Avatar: avatar,
		}
		_, err := d.Service.SavePlayer(r.Context(), player)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/errors/"+err.Error(), http.StatusTemporaryRedirect)
			return
		}
	}

	session := &models.Session{
		ID:      ulid.MustNew(uint64(time.Now().UnixMilli()), rand.Reader).String(),
		Player:  player,
		Started: time.Now(),
	}

	if _, err := d.Service.SaveSession(r.Context(), session); err != nil {
		log.Println(err)
		http.Redirect(w, r, "/errors/"+err.Error(), http.StatusTemporaryRedirect)
		return
	}

	sessionCookie := http.Cookie{
		Name:     "voidsent_session",
		Value:    session.ID,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		MaxAge:   24 * 60 * 60,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
	}

	http.SetCookie(w, &sessionCookie)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
