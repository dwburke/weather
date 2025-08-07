package cmd

import (
	"fmt"

	"github.com/dwburke/weather/types"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(test)
}

var test = &cobra.Command{
	Use:   "test",
	Short: "Run a test command",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		t := &types.Test{}
		t.Amount = 12.34
		t.DateTime = "2024-10-01 12:34:56"

		if err := t.Create(); err != nil {
			return fmt.Errorf("Error creating Test record: %w", err)
		}
		return nil
	},
}
