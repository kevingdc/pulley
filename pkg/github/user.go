package github

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/kevingdc/pulley/pkg/config"
)

type User struct {
	ID int64 `json:"id"`
}

func GetUser(config *config.Config, code string) (*User, error) {
	accessToken, err := getUserAccessToken(config, code)
	if err != nil {
		return nil, err
	}

	user, err := getUserInfo(accessToken)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func getUserAccessToken(config *config.Config, code string) (string, error) {
	values := map[string]string{
		"client_id":     config.GithubClientID,
		"client_secret": config.GithubClientSecret,
		"code":          code,
	}

	json_data, err := json.Marshal(values)
	if err != nil {
		return "", err
	}

	res, err := http.Post("https://github.com/login/oauth/access_token", "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	accessToken := parseAccessToken(string(body))

	return accessToken, nil
}

func getUserInfo(accessToken string) (*User, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var user User
	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func parseAccessToken(response string) string {
	params := strings.Split(response, "&")

	for _, param := range params {
		splitParam := strings.Split(param, "=")
		if splitParam[0] == "access_token" {
			return splitParam[1]
		}
	}

	return ""
}
