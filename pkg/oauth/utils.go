package oauth

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/ebadfd/jira_sucks/lib"
	"github.com/ebadfd/jira_sucks/pkg/handlers"
)

func generateStateOauthCookie(w http.ResponseWriter, r *http.Request) string {
	session, err := handlers.Store.Get(r, lib.OAuthSessionName)

	if err != nil {
		panic(err)
	}

	b := make([]byte, 16)
	rand.Read(b)

	state := base64.URLEncoding.EncodeToString(b)

	session.Values[lib.OAuthStateKey] = state
	err = session.Save(r, w)

	if err != nil {
		panic(err)
	}

	return state
}

type ExchangeTokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
}

type ExchangeTokenResult struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func (h *JiraOAuthServiceImpl) ExchangeToken(code string) (*ExchangeTokenResult, error) {
	client := &http.Client{}
	request := ExchangeTokenRequest{
		GrantType:    "authorization_code",
		ClientID:     h.oauth.ClientID,
		ClientSecret: h.oauth.ClientSecret,
		Code:         code,
		RedirectURI:  h.oauth.RedirectURL,
	}
	var result ExchangeTokenResult

	payload, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		"https://auth.atlassian.com/oauth/token",
		bytes.NewBuffer(payload),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
