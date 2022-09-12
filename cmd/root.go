package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ipfs-upload",
		Short: "ipfs-upload is a CLI tool made to speed up your file uploading process to a local repository.",
	}

	configFile string
	apiUrl     string
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func defaultConfigPath() string {
	dotConfigPath, err := os.UserConfigDir()
	cobra.CheckErr(err)

	return path.Join(dotConfigPath, "config.yaml")
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(
		&configFile, "config", "c", defaultConfigPath(), "Configuration file path",
	)
	rootCmd.PersistentFlags().StringVarP(
		&apiUrl, "api-url", "u", "", "Use a different api url",
	)

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.AddCommand(
		fileCmd,
		clipboardCmd,
		apiCmd,
	)
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		cfgPath, err := os.UserConfigDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(path.Join(cfgPath, "/ipfs-upload"))
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
