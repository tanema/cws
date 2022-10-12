package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/akyoto/tty"
	"github.com/spf13/cobra"

	"github.com/tanema/cws/lib/gcloud"
	"github.com/tanema/cws/lib/term"
)

var rootCmd = &cobra.Command{
	Use:     "cws",
	Version: "0.0.1",
	Short:   "A tool for managing chrome webstore extensions",
	Long: `A cli for automating or locally managing and publishing chrome webstore
extensions. Best used in CI.

Env Vars:
  CWS_EXTENSION_ID     chrome webstore id of the extension
  CWS_CLIENT_ID        google oauth client id
  CWS_CLIENT_SECRET    google oauth client secret
  CWS_REFRESH_TOKEN    google oauth client refresh token. Run cws init to get this value
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
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

func spinner(message string, fn func() error) error {
	if tty.IsTerminal(os.Stdout.Fd()) {
		return term.Spinner(message, fn)
	}
	fmt.Printf("🔄: %v", message)
	err := fn()
	if err != nil {
		fmt.Printf("🔥: %v", message)
	} else {
		fmt.Printf("✅: %v", message)
	}
	return err
}
