package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/context"

	"github.com/ebadfd/jira_sucks/lib"
	"github.com/gorilla/sessions"
)

type AuthSession struct {
	Token        string
	DisplayName  string
	ProfileImage string
}

func Test(w http.ResponseWriter, r *http.Request) {
	adm := context.Get(r, lib.AuthResults).(AuthSession)

	fmt.Println(adm)

	w.Write([]byte("This is my home page"))
}

// TODO: update this
var Store = sessions.NewCookieStore([]byte("secret_key"))

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := Store.Get(r, lib.OAuthSessionName)

		if err != nil {
			panic(err)
		}

		profileSession, err := Store.Get(r, lib.ProfileSessionName)

		if err != nil {
			panic(err)
		}

		token := session.Values[lib.OAuthStateToken]
		displayName := profileSession.Values[lib.ProfileUserDisplayName]
		image := profileSession.Values[lib.ProfileUserImage]

		if token == nil {
			panic("no auth")
		}

		authResults := AuthSession{
			Token:        token.(string),
			DisplayName:  displayName.(string),
			ProfileImage: image.(string),
		}

		context.Set(r, lib.AuthResults, authResults)

		// handle validations
		next.ServeHTTP(w, r)
	})
}

func TrailingSlashMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		next.ServeHTTP(w, r)
	})
}
