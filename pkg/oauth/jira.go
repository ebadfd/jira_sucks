package oauth

import (
	"net/http"

	"github.com/ebadfd/jira_sucks/lib"
	"github.com/ebadfd/jira_sucks/pkg/handlers"
	"golang.org/x/oauth2"
)

type JiraOAuthServiceImpl struct {
	log   lib.Logger
	conf  *lib.Configuration
	oauth oauth2.Config
}

func NewJiraOAuthServiceImpl(log lib.Logger, conf *lib.Configuration) *JiraOAuthServiceImpl {
	var jiraOAuthConfig = &oauth2.Config{
		RedirectURL:  "http://127.0.0.1:3000/auth/callback",
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
	oauthState := generateStateOauthCookie(w, r)

	/*
		AuthCodeURL receive state that is a token to protect the user from CSRF attacks. You must always provide a non-empty string and
		validate that it matches the the state query parameter on your redirect callback.
	*/
	u := h.oauth.AuthCodeURL(oauthState)
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
		panic(err)
	}

	accessibleResource, err := lib.AccessibleResources(res.AccessToken)

	if err != nil {
		panic(err)
	}

	client := lib.JiraClient(accessibleResource.ID, res.AccessToken)

	profile, _, err := client.User.GetSelf()

	if err != nil {
		panic(err)
	}

	session.Values[lib.OAuthStateToken] = res.AccessToken
	session.Values[lib.OAuthCloudId] = accessibleResource.ID
	profileSession.Values[lib.ProfileUserDisplayName] = profile.DisplayName
	profileSession.Values[lib.ProfileUserImage] = profile.AvatarUrls.Four8X48

	err = session.Save(r, w)

	if err != nil {
		panic(err)
	}

	err = profileSession.Save(r, w)

	if err != nil {
		panic(err)
	}

	w.Write([]byte("login success"))
}
