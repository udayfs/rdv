package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/udayfs/rdv/utils"
	"os"
)

var (
	file   string
	dir    string
	outDir string
)

var rootCmd = &cobra.Command{
	Use:   "rdv [Commands] [Flags]",
	Short: "Access your cloud drive storage from the terminal!",
	Long:  "rdv (Remote Drive View) is a cli tool that can fetch and upload files and directories to the specified drive.",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	Version: utils.Version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(utils.Colorize(utils.Red, "[Error] "), err.Error())
		os.Exit(1)
	}
}

func init() {}
