package session

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/oklog/ulid"
	"golang.org/x/oauth2"

	"github.com/gridanias-helden/voidsent/pkg/middleware"
	"github.com/gridanias-helden/voidsent/pkg/models"
	"github.com/gridanias-helden/voidsent/pkg/storage"
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
	OAuth    *oauth2.Config
	Sessions storage.Sessions
	// Players  storage.Players
	KV map[string]time.Time
}

func (d *Discord) Auth(w http.ResponseWriter, r *http.Request) {
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
		log.Println("State does not match or is expired.")
		return
	}

	// We exchange the code we got for an access token
	// Then we can use the access token to do actions, limited to scopes we requested
	token, err := d.OAuth.Exchange(context.Background(), r.FormValue("code"))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Use the access token, here we use it to get the logged-in user's info.
	res, err := d.OAuth.Client(context.Background(), token).Get("https://discord.com/api/users/@me")

	if err != nil || res.StatusCode != 200 {
		w.WriteHeader(http.StatusInternalServerError)
		if err != nil {
			log.Println(err)
		} else {
			log.Println(err)
		}
		return
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var user discordUser
	if err := json.Unmarshal(body, &user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	userName := user.Username

	var avatar string
	if user.Avatar != nil {
		avatar = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png?size=64", user.ID, *user.Avatar)
	} else {
		avatar = GuestAvatar()
	}

	session := models.Session{
		ID:       ulid.MustNew(uint64(time.Now().UnixMilli()), rand.Reader).String(),
		PlayerID: "discord:" + user.ID,
		Avatar:   avatar,
		Username: userName,
	}

	if _, err := d.Sessions.SaveSession(r.Context(), session); err != nil {
		log.Println(err)
		http.Redirect(w, r, "/errors/"+err.Error(), http.StatusTemporaryRedirect)
		return
	}

	sessionCookie := http.Cookie{
		Name:     "voidsent_session",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
	}

	http.SetCookie(w, &sessionCookie)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (d *Discord) Logout(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(middleware.SessionKey).(models.Session)
	if !ok {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if session.WS != nil {
		_ = session.WS.Close()
	}

	if err := d.Sessions.DeleteSession(r.Context(), session); err != nil {
		log.Println(err)
		http.Redirect(w, r, "/errors/"+err.Error(), http.StatusTemporaryRedirect)
		return
	}

	sessionCookie := http.Cookie{
		Name:     "voidsent_session",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-24 * time.Hour),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
	}

	http.SetCookie(w, &sessionCookie)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
