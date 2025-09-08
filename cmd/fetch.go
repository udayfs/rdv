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
	Aliases:               []string{"get"},
	PreRun: func(cmd *cobra.Command, args []string) {
		err := auth()
		if err != nil {
			utils.ExitOnError("Authorization error: " + err.Error())
		}

		if (file == "" && dir == "") || (file != "" && dir != "") {
			utils.ExitOnError("You must provide either -f (file) or -d (directory), but not both")
		}
		if err = utils.ClearScreen(); err != nil {
			utils.ExitOnError(err.Error())
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var message string
		var driveFName string

		if file == "" {
			driveFName = dir
			isDir = true
			message = fmt.Sprintf("Fetching directory %s", dir)
		} else {
			driveFName = file
			isDir = false
			message = fmt.Sprintf("Fetching file %s", file)
		}

		_, err := utils.Spinner(func() (struct{}, error) {
			return struct{}{}, fetch(driveFName)
		}, message)

		if err != nil {
			utils.ExitOnError("Unable to complete the fetch operation: " + err.Error())
		}

		utils.ExitOnSuccess("Fetch operation for " + driveFName + " successfully completed!")
	},
}

func init() {
	fetchCmd.Flags().StringVarP(&file, "file", "f", "", "file to fetch")
	fetchCmd.Flags().StringVarP(&dir, "dir", "d", "", "directory to fetch")
	fetchCmd.Flags().StringVarP(&outDir, "outPath", "o", ".", "path for placing the fetched file/directory")
	rootCmd.AddCommand(fetchCmd)
}

func fetch(driveFileName string) error {
	var err error

	if isDir {
		err = fetchDriveFolder(driveFileName)
	} else {
		err = fetchDriveFile(driveFileName)
	}

	return err
}

func fetchDriveFile(f string) error {
	panic("unimplemented")
}

func fetchDriveFolder(fol string) error {
	panic("unimplemented")
}
