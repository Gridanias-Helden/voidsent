package session

import (
	"crypto/rand"
	"fmt"
	rnd "math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/oklog/ulid"

	"github.com/gridanias-helden/voidsent/pkg/models"
	"github.com/gridanias-helden/voidsent/pkg/storage"
)

var nameGen = map[string]func() string{
	"de": GenNameDE,
	"en": GenNameEN,
}

type GuestLogin struct {
	Sessions storage.Sessions
}

func GuestAvatar() string {
	return fmt.Sprintf("/assets/images/avatars/con%d.png", rnd.Intn(42)+1)
}

func (gl *GuestLogin) Register(w http.ResponseWriter, r *http.Request) {
	languages := strings.Split(r.Header.Get("Accept-Language"), ",")
	lang := "en"
	for _, l := range languages {
		if _, ok := nameGen[l]; ok {
			lang = l
			break
		} else {
		}
	}
	playerName := nameGen[lang]()

	sess := models.Session{
		ID:       ulid.MustNew(ulid.Now(), rand.Reader).String(),
		PlayerID: "guest:" + playerName,
		Username: playerName,
		Avatar:   GuestAvatar(),
		Updated:  time.Now().UTC(),
	}

	_, _ = gl.Sessions.SaveSession(r.Context(), sess)

	sessionCookie := http.Cookie{
		Name:     "voidsent_session",
		Value:    sess.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
	}

	http.SetCookie(w, &sessionCookie)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
