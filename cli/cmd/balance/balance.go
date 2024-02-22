package balance

import (
	"cli/cmd"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	// Attach to root command
	cmd.RootCmd.AddCommand(balanceCmd)
}

// Balance Command
//
// Allows the user to check their current balance
var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Add short description", // TODO: Add short description
	Long:  "Add long description",  // TODO: Add long description
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("'ihp balance' called")
	},
}
