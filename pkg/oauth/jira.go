package oauth

import (
	"fmt"
	"net/http"

	"github.com/ebadfd/jira_sucks/lib"
	"github.com/ebadfd/jira_sucks/pkg/handlers"
	"github.com/ebadfd/jira_sucks/views"
	"golang.org/x/oauth2"
)

type JiraOAuthServiceImpl struct {
	log   lib.Logger
	conf  *lib.Configuration
	oauth oauth2.Config
}

func NewJiraOAuthServiceImpl(log lib.Logger, conf *lib.Configuration) *JiraOAuthServiceImpl {
	var jiraOAuthConfig = &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%s/auth/callback", conf.Host),
		ClientID:     conf.JiraClientId,
		ClientSecret: conf.JiraClientSecret,
		Scopes:       JiraPermissions,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://auth.atlassian.com/authorize?audience=api.atlassian.com",
			TokenURL: "https://auth.atlassian.io/oauth2/token",
		},
	}

	return &JiraOAuthServiceImpl{
		log:   log,
		conf:  conf,
		oauth: *jiraOAuthConfig,
	}
}

func (h *JiraOAuthServiceImpl) OAuthJiraLogin(w http.ResponseWriter, r *http.Request) {
	oauthState, err := generateStateOauthCookie(w, r)

	if err != nil {
		lib.Render(w, http.StatusBadRequest, views.ErrorPage(err))
		return
	}

	/*
		AuthCodeURL receive state that is a token to protect the user from CSRF attacks. You must always provide a non-empty string and
		validate that it matches the the state query parameter on your redirect callback.
	*/
	u := h.oauth.AuthCodeURL(*oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (h *JiraOAuthServiceImpl) OAuthJiraCallback(w http.ResponseWriter, r *http.Request) {
	session, _ := handlers.Store.Get(r, lib.OAuthSessionName)
	profileSession, _ := handlers.Store.Get(r, lib.ProfileSessionName)
	oauthState := session.Values[lib.OAuthStateKey]

	if r.FormValue("state") != oauthState {
		h.log.Warnf("invalid oauth cookie state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res, err := h.ExchangeToken(r.FormValue("code"))

	if err != nil {
		lib.Render(w, http.StatusBadRequest, views.ErrorPage(err))
		return
	}

	accessibleResource, err := lib.AccessibleResources(res.AccessToken)

	if err != nil {
		lib.Render(w, http.StatusBadRequest, views.ErrorPage(err))
		return
	}

	client, err := lib.JiraClient(accessibleResource.ID, res.AccessToken)

	if err != nil {
		lib.Render(w, http.StatusBadRequest, views.ErrorPage(err))
		return
	}

	profile, _, err := client.User.GetSelf()

	if err != nil {
		lib.Render(w, http.StatusBadRequest, views.ErrorPage(err))
		return
	}

	session.Values[lib.OAuthStateToken] = res.AccessToken
	session.Values[lib.OAuthCloudId] = accessibleResource.ID
	profileSession.Values[lib.ProfileUserDisplayName] = profile.DisplayName
	profileSession.Values[lib.ProfileUserImage] = profile.AvatarUrls.Four8X48

	err = session.Save(r, w)

	if err != nil {
		lib.Render(w, http.StatusBadRequest, views.ErrorPage(err))
		return
	}

	err = profileSession.Save(r, w)

	if err != nil {
		lib.Render(w, http.StatusBadRequest, views.ErrorPage(err))
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
