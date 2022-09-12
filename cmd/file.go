package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	fileCmd = &cobra.Command{
		Use:     "file [flags] <filename>",
		Aliases: []string{"f"},
		Short:   "Uploads a file to your local repository",
		Example: `  ipfs-upload file Screenshot_2022-07-12_12-08-29.png
	  ipfs-upload file --new-name my-precious.txt unnamed.txt`,
		Args: cobra.ExactArgs(1),
		Run:  fileRun,
	}

	newName string
)

func init() {
	fileCmd.Flags().StringVarP(&newName, "new-name", "n", "", "New name for the file")
}

func fileRun(cmd *cobra.Command, args []string) {
	fmt.Println(args)
}
