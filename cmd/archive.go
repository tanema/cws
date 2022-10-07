package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tanema/cws/lib/archive"
	"github.com/tanema/cws/lib/term"
)

var archiveCmd = &cobra.Command{
	Use:   "archive [dir-path]",
	Short: "zip the dist directory, update the manifest version at the same time",
	Run: func(cmd *cobra.Command, args []string) {
		version := getVersion(cmd)
		path := archiveExt(args[0], version)
		term.Println(`✅ {{.Version | bold}} {{"Archive Created At:" | green}} {{.Path | cyan}}`, struct {
			Version string
			Path    string
		}{version, path})
	},
}

func init() {
	rootCmd.AddCommand(archiveCmd)
	archiveCmd.Args = cobra.ExactArgs(1)
}

func archiveExt(dirPath, version string) (archivePath string) {
	var err error
	cobra.CheckErr(term.Spinner("Creating Archive", func() error {
		archivePath, err = archive.Zip(dirPath, version)
		return err
	}))
	return archivePath
}
