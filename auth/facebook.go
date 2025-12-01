package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const facebookPrefix = "fb:"

func GetFacebookID(r *http.Request, appID, appSecret string) (string, error) {
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

	debugURL := fmt.Sprintf(
		"https://graph.facebook.com/v24.0/debug_token?input_token=%s&access_token=%s|%s",
		inputToken, appID, appSecret,
	)

	resp, err := http.Get(debugURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	if !tokenResp.Data.IsValid {
		return "", fmt.Errorf("facebook token invalid")
	}

	if tokenResp.Data.AppID != appID {
		return "", fmt.Errorf("token not issued for this app")
	}

	if time.Now().Unix() > tokenResp.Data.ExpiresAt {
		return "", fmt.Errorf("token expired")
	}

	return facebookPrefix + tokenResp.Data.UserID, nil
}
