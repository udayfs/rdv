package cmd

import (
	"github.com/spf13/cobra"
	"github.com/udayfs/rdv/utils"
)

var (
	file   string
	dir    string
	outDir string
)

var rootCmd = &cobra.Command{
	Use:               "rdv",
	Short:             "Access your cloud drive storage from the terminal!",
	Long:              "rdv (Remote Drive View) is a cli tool that can fetch and upload files and directories to the specified drive.",
	SilenceErrors:     true,
	Version:           utils.Version,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.ExitOnError(err.Error())
	}
}

func init() {}
