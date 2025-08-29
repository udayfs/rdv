package cmd

import (
	"github.com/spf13/cobra"
	"github.com/udayfs/rdv/utils"
)

var authCmd = &cobra.Command{
	Use:                   "auth",
	Short:                 "authenticate a drive user",
	Long:                  "authenticate a drive user",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {

		if err := utils.ClearScreen(); err != nil {
			utils.ExitOnError(err.Error())
		}
		
		utils.LogIn()
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
