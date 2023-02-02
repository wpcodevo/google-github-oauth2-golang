package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/wpcodevo/google-github-oath2-golang/initializers"
)

type GitHubOauthToken struct {
	Access_token string
}

type GitHubUserResult struct {
	Name  string
	Photo string
	Email string
}

func GetGitHubOauthToken(code string) (*GitHubOauthToken, error) {
	const rootURl = "https://github.com/login/oauth/access_token"

	config, _ := initializers.LoadConfig(".")
	values := url.Values{}
	values.Add("code", code)
	values.Add("client_id", config.GitHubClientID)
	values.Add("client_secret", config.GitHubClientSecret)

	query := values.Encode()

	queryString := fmt.Sprintf("%s?%s", rootURl, bytes.NewBufferString(query))
	req, err := http.NewRequest("POST", queryString, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve token")
	}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, res.Body)
	if err != nil {
		return nil, err
	}

	parsedQuery, err := url.ParseQuery(resBody.String())
	if err != nil {
		return nil, err
	}

	tokenBody := &GitHubOauthToken{
		Access_token: parsedQuery["access_token"][0],
	}

	return tokenBody, nil
}

func GetGitHubUser(access_token string) (*GitHubUserResult, error) {
	rootUrl := "https://api.github.com/user"

	req, err := http.NewRequest("GET", rootUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))

	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve user")
	}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, res.Body)
	if err != nil {
		return nil, err
	}

	var GitHubUserRes map[string]interface{}

	if err := json.Unmarshal(resBody.Bytes(), &GitHubUserRes); err != nil {
		return nil, err
	}

	userBody := &GitHubUserResult{
		Email: GitHubUserRes["email"].(string),
		Name:  GitHubUserRes["login"].(string),
		Photo: GitHubUserRes["avatar_url"].(string),
	}

	return userBody, nil
}
