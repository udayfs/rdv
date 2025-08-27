package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/udayfs/rdv/utils"
)

var uploadCmd = &cobra.Command{
	Use:                   "upload [-f file | -d directory]",
	Short:                 "uploads a file or directory to the drive",
	Long:                  "uploads a file or directory to the drive",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {

		if (file == "" && dir == "") || (file != "" && dir != "") {
			utils.ExitOnError("you must provide either -f (file) or -d (directory), but not both")
		}

		if err := utils.ClearScreen(); err != nil {
			utils.ExitOnError(err.Error())
		}
		fmt.Println(utils.Colorize(utils.Gray, "[Info]"), "Uploading", file)
		// upload logic

	},
}

func init() {
	uploadCmd.Flags().StringVarP(&file, "file", "f", "", "file to upload")
	uploadCmd.Flags().StringVarP(&dir, "dir", "d", "", "directory to upload")
	rootCmd.AddCommand(uploadCmd)
}
