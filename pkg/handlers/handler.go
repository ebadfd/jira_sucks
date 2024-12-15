package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/context"

	"github.com/ebadfd/jira_sucks/lib"
	"github.com/ebadfd/jira_sucks/views"
	"github.com/gorilla/sessions"
)

// TODO: update this
var Store = sessions.NewCookieStore([]byte("secret_key"))

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := Store.Get(r, lib.OAuthSessionName)

		if err != nil {
			lib.Render(w, http.StatusBadRequest, views.ErrorPage(err))
			return
		}

		profileSession, err := Store.Get(r, lib.ProfileSessionName)

		if err != nil {
			lib.Render(w, http.StatusBadRequest, views.ErrorPage(err))
			return
		}

		token := session.Values[lib.OAuthStateToken]
		displayName := profileSession.Values[lib.ProfileUserDisplayName]
		image := profileSession.Values[lib.ProfileUserImage]
		cloudId := session.Values[lib.OAuthCloudId]

		if token == nil {
			http.Redirect(w, r, "/auth/jira", http.StatusTemporaryRedirect)
			return
		}

		if cloudId == nil {
			http.Redirect(w, r, "/auth/jira", http.StatusTemporaryRedirect)
			return
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

func Error(w http.ResponseWriter, r *http.Request) {
	lib.Render(w, http.StatusOK, views.ErrorPage(errors.New("something werong wrong")))
}
