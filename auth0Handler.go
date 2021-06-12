package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthHandler struct {
	ClientID string
	Domain   string
}

func NewAuthHandler(clientID string, domain string) (AuthHandler, error) {
	return AuthHandler{
		ClientID: clientID,
		Domain:   domain,
	}, nil
}

func (authHandler *AuthHandler) GetRedirectURL(callbackURL string) string {
	return fmt.Sprintf(
		"%s/authorize?response_type=token&client_id=%s&redirect_uri=%s&scope=%s&state=%s",
		authHandler.Domain, authHandler.ClientID, callbackURL,
		"openid%20profile", "xyzabc",
	)
}
func (authHandler *AuthHandler) GetUserInfo(token string) (*UserInfo, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", authHandler.Domain+"/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

type UserInfo struct {
	FamilyName string `json:"family_name"`
	GivenName  string `json:"given_name"`
	Locale     string `json:"locale"`
	Name       string `json:"name"`
	Nickname   string `json:"nickname"`
	Picture    string `json:"picture"`
	Sub        string `json:"sub"`
	UpdatedAt  string `json:"updated_at"`
}
