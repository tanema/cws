package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/tanema/cws/lib/gcloud"
	"github.com/tanema/cws/lib/term"
)

var rootCmd = &cobra.Command{
	Use:   "cws",
	Short: "A tool for managing chrome webstore extensions",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "./chrome_webstore.json", "id of extension to deploy")
	rootCmd.PersistentFlags().StringP("version", "v", "", "version to add to the manifest (default: yy.mm.dd.nn)")
}

func authenticate(cmd *cobra.Command) *gcloud.Client {
	var client *gcloud.Client
	var err error
	err = term.Spinner("Authenticating", func() error {
		client, err = gcloud.New(getString(cmd, "config"))
		return err
	})
	cobra.CheckErr(err)
	return client
}

func getVersion(cmd *cobra.Command) string {
	if version := getString(cmd, "version"); version != "" {
		return version
	}
	now := time.Now()
	return fmt.Sprintf("%v.%v", now.Format("06.1.2"), ((now.Hour()*60 + now.Minute()) / 10))
}

func getString(cmd *cobra.Command, key string) string {
	value, err := cmd.Flags().GetString(key)
	cobra.CheckErr(err)
	return value
}
