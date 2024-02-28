package instance

import (
	"cli/internal/authentication"
	"cli/internal/backend"
	"cli/internal/messages"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	InstanceCmd.AddCommand(createCmd)
	InstanceCmd.AddCommand(deleteCmd)
}

// Instance Command
//
// Displays information about a specific instance
var InstanceCmd = &cobra.Command{
	Use:   "instance [instance_id]",
	Short: "Displays information about a specific instance", 
	Long:  "Displays information about a specific instance",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		
		tokens := authentication.GetCurrentAuthentication()
		if tokens == nil {
			messages.DisplayNotLoggedInMessage()
			return
		}

		instance, err := backend.DefaultBackendClient.GetUserInstance(tokens.AccessToken, id, true)
		if errors.Is(err, backend.ErrNotAuthenticated) {
			messages.DisplaySessionExpiredMessage()
			return
		} else if err != nil {
			fmt.Println("Unable to get your instance:", err)
			fmt.Println("Please try again later.")
			return
		}

		fmt.Println(instance.String())
	},
}