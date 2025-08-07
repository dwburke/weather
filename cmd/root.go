package cmd

import (
	"github.com/spf13/cobra"
)

var cfgFile string
var verbose bool

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.entity.yaml)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "verbose output")
}

var rootCmd = &cobra.Command{
	Use:   "weather",
	Short: "weather tracker",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}
