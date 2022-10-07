package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tanema/cws/lib/gcloud"
	"github.com/tanema/cws/lib/term"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AccessCodeResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	Expiry       int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

var initCmd = &cobra.Command{
	Use:   "init [client-id] [client-secret]",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		state := genState()
		scopes := "https://www.googleapis.com/auth/chromewebstore"
		codeChan := make(chan string)
		go startAuthServer(state, scopes, codeChan)
		conf := &oauth2.Config{
			ClientID:     args[0],
			ClientSecret: args[1],
			RedirectURL:  "http://localhost:3333",
			Scopes:       strings.Split(scopes, " "),
			Endpoint:     google.Endpoint,
		}
		term.Println("Please visit {{. | blue}} to start auth", conf.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce))
		code := <-codeChan
		var token *AccessCodeResponse
		var err error
		cobra.CheckErr(term.Spinner("Exchanging Code", func() error {
			token, err = exchangeCode(args[0], args[1], code)
			return err
		}))

		cobra.CheckErr(term.Spinner("Saving config", func() error {
			conf := gcloud.Config{
				ID:           args[0],
				Secret:       args[1],
				RefreshToken: token.RefreshToken,
			}
			confBytes, err := json.MarshalIndent(conf, "", "\t")
			if err != nil {
				return err
			}
			return os.WriteFile("chrome_webstore.json", confBytes, 0666)
		}))
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Args = cobra.ExactArgs(2)
}

func startAuthServer(expectedState, expectedScopes string, codeChan chan string) {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()
		if expectedState != query.Get("state") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request: Unmatched State, please try again\n"))
		} else if expectedScopes != query.Get("scope") {
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
	if err := http.ListenAndServe("localhost:3333", nil); err != nil {
		log.Fatal(err)
	}
}

func exchangeCode(clientID, clientSecret, code string) (*AccessCodeResponse, error) {
	form := url.Values{}
	form.Set("code", code)
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("redirect_uri", "http://localhost:3333")
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
	codeResp := &AccessCodeResponse{}
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
