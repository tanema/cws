package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tanema/cws/lib/gcloud"
	"github.com/tanema/cws/lib/term"
)

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish the extension to the chrome webstore",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚚 Publishing")
		public, _ := cmd.Flags().GetBool("public")
		status, err := publish(authenticate(cmd), public)
		if err != nil {
			term.Println(`{{. | bold}}`, err)
			return
		}
		term.Println(`✅ {{"Publish Successfully" | green}} Publication Status: {{. | cyan}}`, status.Status)
		term.Println("See package status at: {{. | blue}}", "https://chrome.google.com/webstore/devconsole")
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)
	publishCmd.Flags().StringP("config", "c", "./chrome_webstore.json", "id of extension to deploy")
	publishCmd.Flags().BoolP("public", "p", false, "Deploy to public, otherwise default to trustedTesters")
}

func publish(client *gcloud.Client, public bool) (status gcloud.WebStoreItem, err error) {
	audience := "test users"
	if public {
		audience = "public"
	}
	term.Spinner(fmt.Sprintf("Publishing to %v", audience), func() error {
		status, err = client.PublishExtension(public)
		return err
	})
	return
}
