package instance

import (
	"cli/internal/authentication"
	"cli/internal/backend"
	"cli/internal/messages"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// Instances Delete Command
//
// Allows the user to delete an instance
var deleteCmd = &cobra.Command{
	Use:   "delete [instance_id]",
	Short: "Deletes an instance.", 
	Long:  "Deletes an instance.",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		
		tokens := authentication.GetCurrentAuthentication()
		if tokens == nil {
			messages.DisplayNotLoggedInMessage()
			return
		}

		err := backend.DefaultBackendClient.DeleteInstance(tokens.AccessToken, id, true)
		if errors.Is(err, backend.ErrNotAuthenticated) {
			messages.DisplaySessionExpiredMessage()
			return
		} else if err != nil {
			fmt.Println("Unable to delete your instance:", err)
			fmt.Println("Please try again later.")
			return
		}

		fmt.Printf("Instance <<%s>> successfully deleted.", id)
	},
}
