package balance

import (
	"cli/internal/authentication"
	"cli/internal/backend"
	"cli/internal/messages"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	BalanceCmd.AddCommand(purchaseCmd)
}

// Balance Command
//
// Allows the user to check their current balance
var BalanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Add short description", // TODO: Add short description
	Long:  "Add long description",  // TODO: Add long description
	Run: func(cmd *cobra.Command, args []string) {
		tokens := authentication.GetCurrentAuthentication()
		if tokens == nil {
			messages.DisplayNotLoggedInMessage()
			return
		}

		amount, err := backend.DefaultBackendClient.GetBalance(tokens.AccessToken, true)
		if errors.Is(err, backend.ErrNotAuthenticated) {
			messages.DisplaySessionExpiredMessage()
			return
		} else if err != nil {
			displayUnableToGetBalanceMessage(err)
			return
		}

		fmt.Println("Your current balance is:", amount, "credits")
		fmt.Println("To purchase more credits, run 'ihp balance purchase'")
	},
}

func displayUnableToGetBalanceMessage(err error) {
	fmt.Println("Unable to get your balance:", err)
	fmt.Println("Please try again later.")
}
