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
		test, _ := cmd.Flags().GetBool("test")
		status, err := publish(authenticate(cmd), !test)
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
	publishCmd.Flags().BoolP("test", "t", false, "Deploy to test users, otherwise default")
}

func publish(client *gcloud.Client, public bool) (status gcloud.WebStoreItem, err error) {
	audience := ""
	if !public {
		audience = "to test users"
	}
	term.Spinner(fmt.Sprintf("Publishing %v", audience), func() error {
		status, err = client.PublishExtension(public)
		return err
	})
	return
}
