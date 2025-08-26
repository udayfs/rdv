package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:                   "fetch [-f file | -d directory] [-o outputPath]",
	Short:                 "fetch a file or directory",
	Long:                  "fetch a file or directory",
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
	fetchCmd.Flags().StringVarP(&file, "file", "f", "", "file to fetch")
	fetchCmd.Flags().StringVarP(&dir, "dir", "d", "", "directory to fetch")
	fetchCmd.Flags().StringVarP(&outDir, "outPath", "o", ".", "path for placing the fetched file/directory")
	rootCmd.AddCommand(fetchCmd)
}
