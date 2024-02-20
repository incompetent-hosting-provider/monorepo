package balance

import (
	"cli/internal/authentication"
	"cli/internal/backend"
	"cli/internal/messages"
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// Purchase Command
//
// Allows the user to purchaseCmd more credits.
//
// NOTE: This command is only a mocked command and does not actually
// perform a purchaseCmd flow. Currently it will just add credits to the users balance.
var purchaseCmd = &cobra.Command{
	Use:   "purchase [amount]",
	Short: "Purchase additional credits",
	Long:  `Purchase additional credits.
This command allows the user to purchase more credits. 
The amount parameter specifies the number of credits to purchase.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tokens := authentication.GetCurrentAuthentication()
		if tokens == nil {
			messages.DisplayNotLoggedInMessage()
			return
		}

		purchaseAmount, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid amount:",  args[0])
			fmt.Println("Please make sure the amount is a valid (non decimal) number.")
			return
		}
		
		fmt.Printf("Purchasing %d credits...\n", purchaseAmount)
		newAmount, err := backend.Purchase(*tokens, purchaseAmount)
		if err != nil {
			if errors.Is(err, backend.ErrNotAuthenticated) {
				messages.DisplaySessionExpiredMessage()
				return
			} else {
				displayUnableToPurchaseMessage(err)
				return
			}
		}

		fmt.Println("Successfully purchased credits!")
		fmt.Println("Your new balance is:", newAmount, "credits")
	},
}

func displayUnableToPurchaseMessage(err error) {
	fmt.Println("Unable to purchase credits:", err)
	fmt.Println("Please try again later.")
}

