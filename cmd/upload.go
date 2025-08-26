package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:                   "upload [-f file | -d directory]",
	Short:                 "uploads a file or directory to the drive",
	Long:                  "uploads a file or directory to the drive",
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		if (file == "" && dir == "") || (file != "" && dir != "") {
			return fmt.Errorf("you must provide either -f (file) or -d (directory), but not both")
		}

		// fetch logic

		return nil
	},
}

func init() {
	uploadCmd.Flags().StringVarP(&file, "file", "f", "", "file to upload")
	uploadCmd.Flags().StringVarP(&dir, "dir", "d", "", "directory to upload")
	rootCmd.AddCommand(uploadCmd)
}
