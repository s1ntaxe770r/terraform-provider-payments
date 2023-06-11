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

type APIRequest struct {
	ServiceType string                 `json:"service type"`
	RequestRef  string                 `json:"requestref"`
	Data        map[string]interface{} `json:"data", omitempty`
}

// BankList response struct
//
//	{
//		"status": true,
//		"message": "Completed Successfully",
//		"data": {
//			"banks": [
//				{
//					"bankCode":"",
//					"bankName":""
//				},
//			]
//		}
//	}
type BankListResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Banks []struct {
			BankCode string `json:"bankCode"`
			BankName string `json:"bankName"`
		} `json:"banks"`
	} `json:"data"`
}

// enum to represnt service types
type ServiceType string

const (
	BankList ServiceType = "BankList"
)

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
	c.l.Info("Getting auth token")
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

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImN0eSI6IkpXVCJ9.eyJmdWxsbmFtZSI6Ik9MVVdBU0VHVU4gSlVCUklMIE9ZRVRVTkpJIiwiZW1haWwiOiJoZWxsb0BqdWJyaWwueHl6IiwiY2xpZW50LWtleSI6ImFmMFAxNGRVZU15OEpnRGw1elR0IiwiZXhwIjoxNjg2NTE1NTA0LCJpc3MiOiJDb3JlSWRlbnRpdHkiLCJhdWQiOiJDb3JlSWRlbnRpdHkifQ.vFJIIiYYkjVjn69Kl3hS3YemN05UdoB5zwhujeppemY

func (c *ApiClient) GetBankList() (*BankListResponse, error) {
	data := APIRequest{
		ServiceType: string(BankList),
		RequestRef:  "",
	}
	c.l.Debug("Getting bank list")
	c.GetAuthToken()
	c.l.Info("Auth token: %s", c.AuthToken)

	resp, err := c.PostRequest(c.BaseUrl, data)
	if err != nil {
		c.l.Error(err.Error())
	}
	defer resp.Body.Close()
	// read response as bytes and convert to string
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	// print response as string
	fmt.Println(buf.String())

	var bankListResponse BankListResponse
	err = json.Unmarshal(buf.Bytes(), &bankListResponse)
	if err != nil {
		c.l.Error(err.Error())
	}
	return &bankListResponse, nil

}
