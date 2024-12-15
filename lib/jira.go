package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/andygrunwald/go-jira"
)

func JiraClient(cloudId, token string) (*jira.Client, error) {
	base := fmt.Sprintf("https://api.atlassian.com/ex/jira/%s/", cloudId)
	tp := jira.BearerAuthTransport{
		Token: token,
	}

	jiraClient, err := jira.NewClient(tp.Client(), base)
	if err != nil {
		return nil, err
	}

	return jiraClient, nil
}

type AccessibleResource struct {
	ID        string   `json:"id"`
	URL       string   `json:"url"`
	Name      string   `json:"name"`
	Scopes    []string `json:"scopes"`
	AvatarURL string   `json:"avatarUrl"`
}

func AccessibleResources(token string) (*AccessibleResource, error) {
	var results []AccessibleResource

	url := "https://api.atlassian.com/oauth/token/accessible-resources"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	fmt.Println(string(body))

	if err := json.Unmarshal(body, &results); err != nil {
		return nil, err
	}

	// At the moment we only request access for one cloud uri, in case we do support more we need to
	// let the user select the preferable cloud id
	return &results[0], nil
}
