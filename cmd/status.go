package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tanema/cws/lib/gcloud"
	"github.com/tanema/cws/lib/term"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "check the publication status of your extension",
	Run: func(cmd *cobra.Command, args []string) {
		term.Println(`🕵️  {{"Status" | green}}{{with .Draft.UploadState}} Draft: {{. | bold}}{{end}}{{with .Published.UploadState}} Published: {{. | bold}}{{end}}`, status(authenticate(cmd)))
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().StringP("config", "c", "./chrome_webstore.json", "id of extension to deploy")
}

func status(client *gcloud.Client) (status gcloud.WebStoreItemStatus) {
	cobra.CheckErr(term.Spinner("Fetching Status", func() error {
		var err error
		status, err = client.ExtensionStatus()
		return err
	}))
	return
}
