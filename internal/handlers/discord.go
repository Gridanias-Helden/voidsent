package handlers

import (
	"context"
	"crypto/rand"
	"io"
	"net/http"
	"time"

	"github.com/oklog/ulid"
	"golang.org/x/oauth2"

	"github.com/gridanias-helden/voidsent/internal/services"
)

type Discord struct {
	OAuth   *oauth2.Config
	Service services.Service
	KV      map[string]time.Time
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
	if timestamp, ok := d.KV[state]; !ok || timestamp.Add(5*time.Minute).After(time.Now()) {
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

	// player := models.Player{
	// 	ID:     "",
	// 	Name:   "",
	// 	Avatar: "",
	// }

	// cookie := http.Cookie{
	// 	Name:  "voidsent_session",
	// 	Value: sessionID,
	// }

	w.Write(body)
}
