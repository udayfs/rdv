package cmd

import (
	"github.com/spf13/cobra"
	"github.com/udayfs/rdv/utils"
	"google.golang.org/api/drive/v3"
	"os"
	"path/filepath"
)

var uploadCmd = &cobra.Command{
	Use:                   "upload [-f file | -d directory]",
	Short:                 "uploads a file or directory to the drive",
	Long:                  "uploads a file or directory to the drive",
	DisableFlagsInUseLine: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		if (file == "" && dir == "") || (file != "" && dir != "") {
			utils.ExitOnError("You must provide either -f (file) or -d (directory), but not both")
		}
		if err := utils.ClearScreen(); err != nil {
			utils.ExitOnError(err.Error())
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if file != "" {
			f, err := os.Open(file)
			if err != nil {
				utils.ExitOnError(err.Error())
			}
			defer f.Close()

			filename := filepath.Base(file)
			if filename == "." || filename == "/" || filename == "\\" {
				utils.ExitOnError("Invalid file name or path")
			}

			utils.Info("Uploading " + filename)
			uploadedFile, err := srv.Files.Create(&drive.File{Name: filename}).Media(f).Do()
			if err != nil {
				utils.ExitOnError("Unable to upload file: " + err.Error())
			}

			utils.ExitOnSuccess("File uploaded successfully: " + uploadedFile.Name)
		}
	},
}

func init() {
	uploadCmd.Flags().StringVarP(&file, "file", "f", "", "file to upload")
	uploadCmd.Flags().StringVarP(&dir, "dir", "d", "", "directory to upload")
	rootCmd.AddCommand(uploadCmd)
}
