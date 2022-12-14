package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tanema/cws/lib/manifest"
	"github.com/tanema/cws/lib/term"
)

var manifestCmd = &cobra.Command{
	Use:   "manifest [dir-path]",
	Args:  cobra.ExactArgs(1),
	Short: "Update the manifest version, and remove any dev keys",
	Run: func(cmd *cobra.Command, args []string) {
		update_manifest(args[0], getVersion(cmd), getString(cmd, "json"))
	},
}

func init() {
	rootCmd.AddCommand(manifestCmd)
	manifestCmd.Flags().StringP("version", "v", "", "version to add to the manifest (default: yy.mm.dd.nn)")
	manifestCmd.Flags().StringP("json", "j", "", "json changes to the manifest. Should be formatted by key:value comma separated")
}

func update_manifest(path, version, jsonChangeset string) {
	cobra.CheckErr(term.Spinner(fmt.Sprintf("Updating manifest version to %v", version), func() error {
		return manifest.Update(filepath.Join(path, "manifest.json"), version, jsonChangeset)
	}))
}
