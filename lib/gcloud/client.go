package gcloud

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"

	"github.com/tanema/cws/lib/term"
)

type (
	// Client acts as a client to gcloud apis
	Client struct {
		token  string
		Config *Config
	}
	WebStoreItemError struct {
		Code   string `json:"error_code"`
		Detail string `json:"error_detail"`
	}
	WebStoreItemStatus struct {
		Draft     WebStoreItem
		Published WebStoreItem
	}
	WebStoreItem struct {
		ID          string              `json:"id"`
		CRXVersion  string              `json:"crxVersion"`
		Kind        string              `json:"kind"`
		PublicKey   string              `json:"publicKey"`
		UploadState string              `json:"uploadState"`
		ItemError   []WebStoreItemError `json:"itemError"`
		Status      []string            `json:"status"`
		Detail      []string            `json:"statusDetail"`
	}
	gcloudTokenResp struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
		Expiry      int    `json:"expires_in"`
		Type        string `json:"token_type"`
		IDToken     string `json:"id_token"`
	}
	webStoreErrorMessage struct {
		Message string `json:"message"`
		Reason  string `json:"reason"`
	}
	webStoreError struct {
		Code    int                    `json:"code"`
		Message string                 `json:"message"`
		Status  string                 `json:"status"`
		Errors  []webStoreErrorMessage `json:"errors"`
	}
	webStoreErrorResp struct {
		Error webStoreError `json:"error"`
	}
)

// New creates a new gcloud client
func New(configPath string) (*Client, error) {
	config, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}
	client := &Client{Config: config}
	return client, client.authenticate()
}

func (client *Client) authenticate() error {
	params := url.Values{}
	params.Set("client_id", client.Config.ID)
	params.Set("client_secret", client.Config.Secret)
	params.Set("refresh_token", client.Config.RefreshToken)
	params.Set("grant_type", "refresh_token")
	resp := gcloudTokenResp{}
	err := client.doRequest(http.MethodPost, "https://oauth2.googleapis.com/token?"+params.Encode(), nil, &resp)
	client.token = resp.AccessToken
	return err
}

func (client *Client) ExtensionStatus() (WebStoreItemStatus, error) {
	status := WebStoreItemStatus{
		Draft:     WebStoreItem{},
		Published: WebStoreItem{},
	}
	draftErr := client.doRequest(http.MethodGet, "https://www.googleapis.com/chromewebstore/v1.1/items/"+client.Config.ExtID+"?projection=DRAFT", nil, &status.Draft)
	pubErr := client.doRequest(http.MethodGet, "https://www.googleapis.com/chromewebstore/v1.1/items/"+client.Config.ExtID+"?projection=PUBLISHED", nil, &status.Published)
	if draftErr != nil && pubErr != nil {
		return status, draftErr
	}
	return status, nil
}

func (client *Client) CreateExtension(archivePath string) (WebStoreItem, error) {
	resp := WebStoreItem{}
	archive, err := os.Open(archivePath)
	if err != nil {
		return resp, err
	}
	defer archive.Close()
	if err := client.doRequest(http.MethodPost, "https://www.googleapis.com/upload/chromewebstore/v1.1/items?uploadType=media", archive, &resp); err != nil {
		return resp, err
	}
	if resp.UploadState != "SUCCESS" {
		return resp, resp
	}
	return resp, nil
}

func (client *Client) UploadExtension(archivePath string) (WebStoreItem, error) {
	resp := WebStoreItem{}
	archive, err := os.Open(archivePath)
	if err != nil {
		return resp, err
	}
	defer archive.Close()
	url := "https://www.googleapis.com/upload/chromewebstore/v1.1/items/" + client.Config.ExtID + "?uploadType=media"
	if err := client.doRequest(http.MethodPut, url, archive, &resp); err != nil {
		return resp, err
	}
	if resp.UploadState != "SUCCESS" {
		return resp, resp
	}
	return resp, nil
}

func (client *Client) PublishExtension(public bool) (WebStoreItem, error) {
	resp := WebStoreItem{}
	target := "trustedTesters"
	if public {
		target = "default"
	}
	url := "https://www.googleapis.com/chromewebstore/v1.1/items/" + client.Config.ExtID + "/publish?publishTarget=" + target
	if err := client.doRequest(http.MethodPost, url, nil, &resp); err != nil {
		return resp, err
	}
	if !reflect.DeepEqual(resp.Status, []string{"OK"}) {
		return resp, fmt.Errorf("Failed to publish extension with status: %v errors: %v", strings.Join(resp.Status, ", "), strings.Join(resp.Detail, ", "))
	}
	return resp, nil
}

func (client *Client) doRequest(method, url string, body io.Reader, respData interface{}) error {
	if client.Config.Debug {
		fmt.Println("REQUESTION:", method, url)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("constructing new request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+client.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request: %v", err)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading resp body: %v", err)
	}

	if client.Config.Debug {
		fmt.Println("RESPONSE:", string(bodyBytes))
	}

	storeErr := &webStoreErrorResp{}
	if err := json.Unmarshal(bodyBytes, storeErr); err == nil {
		if storeErr.Error.Code != 200 && (storeErr.Error.Status != "" || storeErr.Error.Message != "") {
			return &storeErr.Error
		}
	}

	if respData != nil {
		if err := json.Unmarshal(bodyBytes, respData); err != nil {
			return fmt.Errorf("unmarshalling resp: %v", err)
		}
	}
	return nil
}

func (err *webStoreError) Error() string {
	return term.String(`{{printf "(%v)%v" .Code .Status | yellow}} {{.Message | bold}}`, err)
}

// in the event that ItemError is available
func (item WebStoreItem) Error() string {
	return term.String(`{{.UploadState | yellow}} {{range .ItemError}}
  Code: {{.Code}}
  Detail: {{.Detail}}{{end}}`, item)
}
