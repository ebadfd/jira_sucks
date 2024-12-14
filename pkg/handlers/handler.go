package handlers

import (
	"net/http"
	"strings"

	"github.com/gorilla/context"

	"github.com/ebadfd/jira_sucks/lib"
	"github.com/gorilla/sessions"
)

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
		cloudId := session.Values[lib.OAuthCloudId]

		if token == nil {
			panic("no auth")
		}

		if cloudId == nil {
			panic("no auth")
		}

		authResults := lib.AuthSession{
			CloudId:      cloudId.(string),
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
