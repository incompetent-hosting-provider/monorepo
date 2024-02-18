package balance

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Balance Command
//
// Allows the user to check their current balance
var BalanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Add short description", // TODO: Add short description
	Long:  "Add long description",  // TODO: Add long description
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("'ihp balance' called")
	},
}

func init() {
	// Add subcommands
	BalanceCmd.AddCommand(purchaseCmd)
}
