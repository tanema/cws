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
		status, err := publish(authenticate(cmd))
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
}

func publish(client *gcloud.Client) (status gcloud.WebStoreItem, err error) {
	term.Spinner("Publishing", func() error {
		status, err = client.PublishExtension()
		return err
	})
	return
}
