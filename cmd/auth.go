package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/udayfs/rdv/utils"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"os"
)

var authCmd = &cobra.Command{
	Use:                   "auth [-p <provider>] [-r]",
	Short:                 "authenticate a drive user",
	Long:                  "authenticate a drive user",
	DisableFlagsInUseLine: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := utils.ClearScreen(); err != nil {
			utils.ExitOnError(err.Error())
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		revoke, err := cmd.Flags().GetBool("revoke")
		if err != nil {
			utils.ExitOnError(err.Error())
		}
		if revoke {
			_, err = os.Stat(utils.TokenFilePath)
			if err != nil {
				if os.IsNotExist(err) {
					utils.ExitOnSuccess("User session already revoked!")
				}
				utils.ExitOnError(err.Error())
			}

			err = os.Remove(utils.TokenFilePath)
			if err != nil {
				utils.ExitOnError(err.Error())
			}
			utils.ExitOnSuccess("User auth session revoked successfully!")
		}
	},
}

func init() {
	authCmd.Flags().StringVarP(&provider, "provider", "p", "gdrive", "authorize rdv with a particular remote drive")
	authCmd.Flags().BoolP("revoke", "r", false, "revokes existing logged in user session")
	rootCmd.AddCommand(authCmd)
}

func auth() error {
	ctx := context.Background()

	client, err := utils.LogIn(provider)
	if err != nil {
		return err
	}

	srv, err = drive.NewService(ctx, option.WithHTTPClient(client))
	return err
}
