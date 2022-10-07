package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tanema/cws/lib/gcloud"
	"github.com/tanema/cws/lib/term"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new extension by uploading a brand new archive",
	PreRun: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(cmd.MarkFlagDirname("dir"))
		cobra.CheckErr(cmd.MarkFlagRequired("dir"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		dirPath := getString(cmd, "dir")
		version := getVersion(cmd)
		term.Println("🚚 Creating Version: {{. | bold}}", version)
		client := authenticate(cmd)
		archivePath := archiveExt(dirPath, version)
		defer os.Remove(archivePath)
		status, err := create(client, archivePath)
		if err != nil {
			term.Println(`{{. | bold}}`, err)
			return
		}
		term.Println(`✅ {{. | bold}} {{"Created Successfully" | green}}`, version)
		term.Println(`ID: {{.ID}}
Kind: {{.Kind}}
State: {{.UploadState}}`, status)
		term.Println("See package status at: {{. | blue}}", "https://chrome.google.com/webstore/devconsole")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func create(client *gcloud.Client, archivePath string) (status gcloud.WebStoreItem, err error) {
	term.Spinner("Creating", func() error {
		status, err = client.CreateExtension(archivePath)
		return err
	})
	return
}
