package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/exp/slog"
)

type ApiClient struct {
	ApiToken   string
	Email      string
	l          slog.Logger
	BaseUrl    string
	AuthToken  string
	httpClient *http.Client
}

func NewClient(apiToken string, email string) *ApiClient {
	client := &http.Client{}
	return &ApiClient{
		ApiToken:   apiToken,
		Email:      email,
		BaseUrl:    "https://kuda-openapi.kuda.com/v2.1",
		httpClient: client,
	}
}

func (c *ApiClient) PostRequest(url string, data interface{}) (*http.Response, error) {
	c.l.Debug("Making post request to: %s", url)
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AuthToken))
	return c.httpClient.Do(req)
}

func (c *ApiClient) GetAuthToken() {
	// 	curl \
	// -H 'Content-Type: application/JSON' \
	// -d '{
	// 			"email": "user@example.com",
	// 			"apiKey": "abcd1234keyexample" //this is the key copied from the dashboard
	// 		}'
	// -X POST https://kuda-openapi.kuda.com/v2.1/Account/GetToken
	c.l.Debug("Getting auth token")
	url := fmt.Sprintf("%s/Account/GetToken", c.BaseUrl)
	data := map[string]string{
		"email":  c.Email,
		"apiKey": c.ApiToken,
	}
	resp, err := c.PostRequest(url, data)
	if err != nil {
		c.l.Error(err.Error())
	}
	defer resp.Body.Close()
	// read response as bytes and convert to string
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	c.AuthToken = buf.String()
}
