package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tanema/cws/lib/gcloud"
	"github.com/tanema/cws/lib/term"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload [dir-path]",
	Args:  cobra.ExactArgs(1),
	Short: "Upload a new package",
	Run: func(cmd *cobra.Command, args []string) {
		version := getVersion(cmd)
		term.Println("🚚 Uploading Version: {{. | bold}}", version)
		client := authenticate(cmd)
		archivePath := archiveExt(args[0], version)
		defer os.Remove(archivePath)
		item, err := upload(client, archivePath)
		if err != nil {
			term.Println(`{{. | bold}}`, err)
			return
		}
		term.Println(`✅ {{.Version | bold}} {{"Upload Successful" | green}} Upload State: {{.State | bold}}`, struct {
			State   string
			Version string
		}{item.UploadState, version})
		term.Println("See package status at: {{. | blue}}", "https://chrome.google.com/webstore/devconsole")
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().StringP("config", "c", "./chrome_webstore.json", "id of extension to deploy")
	uploadCmd.Flags().StringP("version", "v", "", "version to add to the manifest (default: yy.mm.dd.nn)")
}

func upload(client *gcloud.Client, archivePath string) (item gcloud.WebStoreItem, err error) {
	term.Spinner("Uploading", func() error {
		item, err = client.UploadExtension(archivePath)
		return err
	})
	return
}
