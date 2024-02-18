package balance

import (
	"cli/internal/authentication"
	"cli/internal/backend"
	"cli/internal/utils"
	"errors"
	"fmt"
	"os"

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
		utils.DisplayNotLoggedInMessage()
		return
		}

		amount, err := backend.GetBalance(*tokens)
		if err != nil {
			if errors.Is(err, backend.ErrNotAuthenticated) {
				// Refresh tokens and try again
				newTokens, err := authentication.RefreshTokens()
				if err != nil {
					displayUnableToGetBalanceMessage(err)
					os.Exit(1)
				}

				if newTokens == nil {
					utils.DisplaySessionExpiredMessage()
					return
				}

				amount, err = backend.GetBalance(*newTokens)
				if err != nil {
					displayUnableToGetBalanceMessage(err)
					os.Exit(1)
				}
			} else {
				displayUnableToGetBalanceMessage(err)
				os.Exit(1)
			}
		}

		fmt.Println("Your current balance is:", amount, "credits")
		fmt.Println("To purchase more credits, run 'ihp balance purchase'")
	},
}

func displayUnableToGetBalanceMessage(err error) {
	fmt.Println("Unable to get your balance:", err)
	fmt.Println("Please try again later.")
}
