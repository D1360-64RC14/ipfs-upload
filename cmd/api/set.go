package api

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	SetCmd = &cobra.Command{
		Use:   "set [flags] <api-url>",
		Short: "Defines a new api url as default",
		Example: `  ipfs-upload api set http://192.168.0.50
  ipfs-upload api set http://localhost:5001
  ipfs-upload api set -p 1337 http://10.0.0.10`,
		Args:   SetCmdValidateArgs,
		PreRun: SetCmdPreRun,
		Run:    func(cmd *cobra.Command, args []string) {},
	}

	apiUrl  string
	apiPort uint16
)

func init() {
	SetCmd.Flags().Uint16VarP(&apiPort, "port", "p", 5001, "API port of the local repository")
}

func SetCmdValidateArgs(cmd *cobra.Command, args []string) error {
	// Must be only 1 arg with an url
	// Validate length
	if err := cobra.ExactArgs(1)(cmd, args); err != nil {
		return err
	}

	// Validate URI standarts
	url, err := url.ParseRequestURI(args[0])
	if err != nil {
		return err
	}

	// Validate schema
	if sort.SearchStrings([]string{"http", "https"}, url.Scheme) != 0 {
		return fmt.Errorf(`scheme %s it's invalid`, url.Scheme)
	}

	switch url.Scheme {
	case "https":
		// SetCmd.Flags().Set("port", "443")
		break
	case "http":
		if url.Port() != "" {
			if _, err := strconv.ParseUint(url.Port(), 10, 16); err != nil {
				return err
			}
			// SetCmd.Flags().Set("port", url.Port())
		}
		break
	}

	return nil
}

func SetCmdPreRun(cmd *cobra.Command, args []string) {

}
