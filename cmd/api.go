package cmd

import (
	"github.com/D1360-64RC14/ipfs-upload/cmd/api"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Configure your default repository api url",
}

func init() {
	apiCmd.AddCommand(
		api.SetCmd,
		api.ViewCmd,
	)
}
