package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/udayfs/rdv/utils"
)

var fetchCmd = &cobra.Command{
	Use:                   "fetch [-f file | -d directory] [-o outputPath]",
	Short:                 "fetch a file or directory",
	Long:                  "fetch a file or directory",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {

		if (file == "" && dir == "") || (file != "" && dir != "") {
			utils.ExitOnError("you must provide either -f (file) or -d (directory), but not both")
		}

		if err := utils.ClearScreen(); err != nil {
			utils.ExitOnError(err.Error())
		}
		fmt.Println(utils.Colorize(utils.Gray, "[Info]"), "Fetching", file)
		// fetch logic
	},
}

func init() {
	fetchCmd.Flags().StringVarP(&file, "file", "f", "", "file to fetch")
	fetchCmd.Flags().StringVarP(&dir, "dir", "d", "", "directory to fetch")
	fetchCmd.Flags().StringVarP(&outDir, "outPath", "o", ".", "path for placing the fetched file/directory")
	rootCmd.AddCommand(fetchCmd)
}
