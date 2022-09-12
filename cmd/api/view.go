package api

import "github.com/spf13/cobra"

var ViewCmd = &cobra.Command{
	Use:     "view",
	Aliases: []string{"ls", "show"},
	Short:   "View your current api url",
	Run:     func(cmd *cobra.Command, args []string) {},
}

func init() {

}
