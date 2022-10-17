package cmd

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"github.com/tanema/cws/lib/gcloud"
	"github.com/tanema/cws/lib/term"
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
		auth := gcloud.NewAuthenticator(args[0], args[1], "https://www.googleapis.com/auth/chromewebstore")
		term.Println("Please visit {{. | blue}} to start auth", auth.URL())
		var conf *gcloud.Config
		var err error
		cobra.CheckErr(term.Spinner("Waiting for response", func() error {
			conf, err = auth.ListForResponse()
			return err
		}))

		cobra.CheckErr(term.Spinner("Saving config", func() error {
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
