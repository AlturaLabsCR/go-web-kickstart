package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const githubPrefix = "gh:"

type GitHubProvider struct {
	ClientID     string
	ClientSecret string
	HTTPClient   *http.Client
}

func (p *GitHubProvider) UserID(r *http.Request) (string, error) {
	if err := r.ParseForm(); err != nil {
		return "", err
	}

	code := r.FormValue("code")
	if code == "" {
		return "", fmt.Errorf("empty code")
	}

	client := p.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}

	form := url.Values{}
	form.Set("client_id", p.ClientID)
	form.Set("client_secret", p.ClientSecret)
	form.Set("code", code)

	req, err := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("no access token returned")
	}

	userReq, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		return "", err
	}

	userReq.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)

	userResp, err := client.Do(userReq)
	if err != nil {
		return "", err
	}
	defer userResp.Body.Close()

	var user struct {
		ID int64 `json:"id"`
	}
	if err := json.NewDecoder(userResp.Body).Decode(&user); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%d", githubPrefix, user.ID), nil
}
