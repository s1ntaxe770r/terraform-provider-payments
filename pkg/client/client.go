package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
)

type ApiClient struct {
	ApiToken   string
	Email      string
	l          logrus.Logger
	BaseUrl    string
	AuthToken  string
	httpClient *http.Client
	AccountNo  string
}

type APIRequest struct {
	ServiceType string                 `json:"servicetype"`
	RequestRef  string                 `json:"requestref"`
	Data        map[string]interface{} `json:"data", omitempty`
}

type Bank struct {
	BankCode string `json:"bankCode"`
	BankName string `json:"bankName"`
}

type Banks []Bank

type BankListResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Data    struct {
		Banks Banks `json:"banks"`
	} `json:"data"`
}

type NameEnquiryResponse struct {
	Message string          `json:"message"`
	Status  bool            `json:"status"`
	Data    NameEnquiryData `json:"data"`
}

type NameEnquiryData struct {
	BeneficiaryAccountNumber string      `json:"beneficiaryAccountNumber"`
	BeneficiaryName          string      `json:"beneficiaryName"`
	SenderAccountNumber      string      `json:"senderAccountNumber"`
	SenderName               interface{} `json:"senderName"`
	BeneficiaryCustomerID    int         `json:"beneficiaryCustomerID"`
	BeneficiaryBankCode      string      `json:"beneficiaryBankCode"`
	NameEnquiryID            int         `json:"nameEnquiryID"`
	ResponseCode             string      `json:"responseCode"`
	TransferCharge           int         `json:"transferCharge"`
	SessionID                string      `json:"sessionID"`
}

type SingleFundTransferResponse struct {
	Message          string      `json:"message"`
	Status           bool        `json:"status"`
	RequestReference string      `json:"requestReference"`
	ResponseCode     string      `json:"responseCode"`
	Data             interface{} `json:"data" omitempty`
}

// enum to represnt service types
type ServiceType string

const (
	BankList    ServiceType = "BANK_LIST"
	NameEnquiry ServiceType = "NAME_ENQUIRY"
	Transfer    ServiceType = "SINGLE_FUND_TRANSFER"
)

func NewClient(apiToken string, email string, accountNo string) *ApiClient {
	client := &http.Client{}
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	return &ApiClient{
		ApiToken:   apiToken,
		Email:      email,
		BaseUrl:    "https://kuda-openapi.kuda.com/v2.1",
		httpClient: client,
		AccountNo:  accountNo,
		l:          *l,
	}
}

func (c *ApiClient) PostRequest(url string, data interface{}, authToken string) (*http.Response, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	return c.httpClient.Do(req)
}

func (c *ApiClient) GetAuthToken() string {
	// 	curl \
	// -H 'Content-Type: application/JSON' \
	// -d '{
	// 			"email": "user@example.com",
	// 			"apiKey": "abcd1234keyexample" //this is the key copied from the dashboard
	// 		}'
	// -X POST https://kuda-openapi.kuda.com/v2.1/Account/GetToken
	url := fmt.Sprintf("%s/Account/GetToken", c.BaseUrl)
	data := map[string]string{
		"email":  c.Email,
		"apiKey": c.ApiToken,
	}
	// make a request using net/http
	resp, err := c.PostRequest(url, data, "")
	if err != nil {
		c.l.Error(err.Error())
	}
	defer resp.Body.Close()
	// read response as bytes and convert to string
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respStr := buf.String()
	return respStr
}

func (c *ApiClient) GetBankList(authToken string) (*BankListResponse, error) {
	data := APIRequest{
		ServiceType: string(BankList),
		RequestRef:  "0",
	}
	c.l.Debug("Getting bank list")

	resp, err := c.PostRequest(c.BaseUrl, data, authToken)
	if err != nil {
		c.l.Error(err.Error())
	}
	defer resp.Body.Close()
	// read response as bytes and convert to string
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	var bankListResponse BankListResponse
	err = json.Unmarshal(buf.Bytes(), &bankListResponse)
	if err != nil {
		c.l.Error(err.Error())
	}
	return &bankListResponse, nil

}

func (c *ApiClient) GetSenderName(beneficiaryAccountNumber string, beneficiaryBankCode string, authToken string) (NameEnquiryResponse, error) {
	data := APIRequest{
		ServiceType: string(NameEnquiry),
		RequestRef:  "",
	}
	data.Data = make(map[string]interface{}) // Initialize the map

	data.Data["beneficiaryAccountNumber"] = beneficiaryAccountNumber
	data.Data["beneficiaryBankCode"] = beneficiaryBankCode
	data.Data["SenderTrackingReference"] = " "
	data.Data["isRequestFromVirtualAccount"] = "false"

	resp, err := c.PostRequest(c.BaseUrl, data, authToken)
	if err != nil {
		return NameEnquiryResponse{}, err
	}
	defer resp.Body.Close()

	// read response as bytes and convert to NameEnquiryResponse
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	var enquireResp NameEnquiryResponse

	err = json.Unmarshal(buf.Bytes(), &enquireResp)
	if err != nil {
		return NameEnquiryResponse{}, err
	}

	if !enquireResp.Status {
		return NameEnquiryResponse{}, errors.New(enquireResp.Message)
	}
	return enquireResp, nil
}

func (c *ApiClient) SingleFundTransfer(beneficiaryAccountNumber string, amount string, bankCode string, authToken string) (SingleFundTransferResponse, error) {
	id, _ := shortid.Generate()
	beneficiaryInfo, err := c.GetSenderName(beneficiaryAccountNumber, bankCode, authToken)
	if err != nil {
		return SingleFundTransferResponse{}, errors.New("Error beneficiar info")
	}

	data := APIRequest{
		ServiceType: string(Transfer),
		RequestRef:  id,
	}
	data.Data = make(map[string]interface{})

	data.Data["beneficiaryAccount"] = beneficiaryAccountNumber
	data.Data["beneficiaryBankCode"] = bankCode
	data.Data["beneficiaryName"] = beneficiaryInfo.Data.BeneficiaryName
	data.Data["amount"] = amount
	data.Data["narration"] = "Sent from Terraform"
	data.Data["nameEnquiryySessionID"] = beneficiaryInfo.Data.SessionID
	data.Data["senderName"] = "Terraform"
	data.Data["clientFeeCharge"] = 0
	data.Data["ClientAccountNumber"] = c.AccountNo

	resp, err := c.PostRequest(c.BaseUrl, data, authToken)
	if err != nil {
		return SingleFundTransferResponse{}, errors.New("Error making transfer: " + err.Error())
	}

	defer resp.Body.Close()

	var transferResp SingleFundTransferResponse

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	err = json.Unmarshal(buf.Bytes(), &transferResp)

	if err != nil {
		return SingleFundTransferResponse{}, errors.New("Error making transfer: " + err.Error())
	}

	if !transferResp.Status {
		return SingleFundTransferResponse{}, errors.New("Error making transfer: " + transferResp.Message)
	}

	return transferResp, nil

}
