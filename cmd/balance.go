package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var balanceCmd = &cobra.Command{
	Use: "balance",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("to be implemented")
	},
	Short: "Show balance info",
}

func init() {
	rootCmd.AddCommand(balanceCmd)
}
