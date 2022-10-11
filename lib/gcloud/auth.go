package gcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type (
	AuthAccess struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		IDToken      string `json:"id_token"`
		Expiry       int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
	}
	Authenticator struct {
		id     string
		secret string
		scopes string
		state  string
		server http.Server
		conf   *oauth2.Config
	}
)

const serverhost = "localhost:3333"

func NewAuthenticator(clientID, clientSecret, scopes string) *Authenticator {
	return &Authenticator{
		id:     clientID,
		secret: clientSecret,
		scopes: scopes,
		server: http.Server{Addr: serverhost},
		conf: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  "http://" + serverhost,
			Scopes:       strings.Split(scopes, " "),
			Endpoint:     google.Endpoint,
		},
	}
}

func (auth *Authenticator) URL() string {
	auth.state = genState()
	return auth.conf.AuthCodeURL(auth.state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

func (auth *Authenticator) ListForResponse() (*Config, error) {
	codeChan := make(chan string)
	go auth.startAuthServer(codeChan)
	defer auth.server.Shutdown(context.Background())
	access, err := auth.exchangeCode(<-codeChan)
	if err != nil {
		return nil, err
	}
	return &Config{
		ID:           auth.id,
		Secret:       auth.secret,
		RefreshToken: access.RefreshToken,
	}, nil
}

func (auth *Authenticator) startAuthServer(codeChan chan string) {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()
		if auth.state != query.Get("state") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request: Unmatched State, please try again\n"))
		} else if auth.scopes != query.Get("scope") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request: required scope was not granted, please try again.\n"))
		} else if code := query.Get("code"); code == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request: No code found in response, please try again\n"))
		} else {
			codeChan <- code
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Succeeded you can now close this tab\n"))
		}
	})
	if err := auth.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (auth *Authenticator) exchangeCode(code string) (*AuthAccess, error) {
	form := url.Values{}
	form.Set("code", code)
	form.Set("client_id", auth.id)
	form.Set("client_secret", auth.secret)
	form.Set("redirect_uri", "http://"+serverhost)
	form.Set("grant_type", "authorization_code")
	resp, err := http.Post(
		"https://oauth2.googleapis.com/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("request: %v", err)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading resp body: %v", err)
	}
	codeResp := &AuthAccess{}
	return codeResp, json.Unmarshal(bodyBytes, codeResp)
}

func genState() string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, 10)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
