package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tanema/cws/lib/term"
)

var deployCmd = &cobra.Command{
	Use:   "deploy [dir-path]",
	Args:  cobra.ExactArgs(1),
	Short: "create an archive, upload, and publish it.",
	Run: func(cmd *cobra.Command, args []string) {
		version := getVersion(cmd)
		term.Println("🚚 Deploying Version: {{. | bold}}", version)
		client := authenticate(cmd)
		archivePath := archiveExt(args[0], version)
		defer os.Remove(archivePath)
		item, err := upload(client, archivePath)
		if err != nil {
			term.Println(`{{. | bold}}`, err)
			return
		}
		status, err := publish(client, false)
		if err != nil {
			term.Println(`{{. | bold}}`, err)
			return
		}
		term.Println(`✅ {{.Version | bold}} {{"Deployed Successfully" | green}}
  Upload State      : {{.State | bold}}
  Publication Status: {{.Status | bold}}`, struct {
			State   string
			Status  string
			Version string
		}{item.UploadState, strings.Join(status.Status, ", "), version})
		term.Println("See package status at: {{. | blue}}", "https://chrome.google.com/webstore/devconsole")
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().StringP("config", "c", "./chrome_webstore.json", "id of extension to deploy")
	deployCmd.Flags().StringP("version", "v", "", "version to add to the manifest (default: yy.mm.dd.nn)")
}
