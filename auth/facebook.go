package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	FacebookAPIVersion = "v24.0"
	facebookPrefix     = "fb:"
)

type FacebookProvider struct {
	AppID      string
	AppSecret  string
	HTTPClient *http.Client
}

func (p *FacebookProvider) UserID(r *http.Request) (string, error) {
	var tokenResp struct {
		Data struct {
			UserID    string `json:"user_id"`
			IsValid   bool   `json:"is_valid"`
			AppID     string `json:"app_id"`
			ExpiresAt int64  `json:"expires_at"`
		} `json:"data"`
	}

	if err := r.ParseForm(); err != nil {
		return "", err
	}

	inputToken := r.FormValue("token")
	if inputToken == "" {
		return "", fmt.Errorf("empty token")
	}

	client := p.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}

	debugURL := fmt.Sprintf(
		"https://graph.facebook.com/%s/debug_token?input_token=%s&access_token=%s|%s",
		FacebookAPIVersion,
		inputToken,
		p.AppID,
		p.AppSecret,
	)

	resp, err := client.Get(debugURL)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	if !tokenResp.Data.IsValid {
		return "", fmt.Errorf("facebook token invalid")
	}

	if tokenResp.Data.AppID != p.AppID {
		return "", fmt.Errorf("token not issued for this app")
	}

	if time.Now().Unix() > tokenResp.Data.ExpiresAt {
		return "", fmt.Errorf("token expired")
	}

	return facebookPrefix + tokenResp.Data.UserID, nil
}
