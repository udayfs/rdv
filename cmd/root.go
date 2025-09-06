package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/udayfs/rdv/utils"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var (
	file     string
	dir      string
	outDir   string
	parent   string
	parentID string
	srv      *drive.Service
)

var rootCmd = &cobra.Command{
	Use:               "rdv",
	Short:             "Access your cloud drive storage from the terminal!",
	Long:              "rdv (Remote Drive View) is a cli tool that can fetch and upload files and directories to the specified drive.",
	Version:           utils.Version,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.ExitOnError(err.Error())
	}
}

func init() {
	rootCmd.SilenceErrors = true
	ctx := context.Background()

	// only google for now
	client, err := utils.LogIn(utils.Providers[0])
	if err != nil {
		utils.ExitOnError(err.Error())
	}

	srv, err = drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		utils.ExitOnError("Unable to retrieve Drive client: " + err.Error())
	}
}
