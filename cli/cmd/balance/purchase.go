package balance

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Purchase Command
//
// Allows the user to purchase more credits.
//
// NOTE: This command is only a mocked command and does not actually
// perform a purchase flow. Currently it will just add credits to the users balance.
var addCmd = &cobra.Command{
	Use:   "purchase",
	Short: "Add short description", // TODO: Add short description
	Long: "Add long description", // TODO: Add long description
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("'ihp balance add' called")
	},
}

func init() {
}