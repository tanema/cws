package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tanema/cws/lib/term"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "create an archive, upload, and publish it.",
	PreRun: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(cmd.MarkFlagDirname("dir"))
		cobra.CheckErr(cmd.MarkFlagRequired("dir"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		dirPath := getString(cmd, "dir")
		version := getVersion(cmd)
		term.Println("🚚 Deploying Version: {{. | bold}}", version)
		client := authenticate(cmd)
		archivePath := archiveExt(dirPath, version)
		defer os.Remove(archivePath)
		item, err := upload(client, archivePath)
		if err != nil {
			term.Println(`{{. | bold}}`, err)
			return
		}
		status, err := publish(client)
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
}
