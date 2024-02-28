package instances

import (
	"cli/internal/authentication"
	"cli/internal/backend"
	"cli/internal/messages"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// Instances Command
//
// Lists all instances belonging to a user
var InstancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "Lists all instances belonging to a user.",
	Long:  "Lists all instances belonging to a user. \n Format is: [STATUS] ID - NAME, IMAGE",
	Run: func(cmd *cobra.Command, args []string) {
		tokens := authentication.GetCurrentAuthentication()
		if tokens == nil {
			messages.DisplayNotLoggedInMessage()
			return
		}

		instances, err := backend.DefaultBackendClient.GetUserInstances(tokens.AccessToken, true)
		if errors.Is(err, backend.ErrNotAuthenticated) {
			messages.DisplaySessionExpiredMessage()
			return
		} else if err != nil {
			fmt.Println("Unable to get your instances:", err)
			fmt.Println("Please try again later.")
			return
		}

		if len(instances) < 1 {
			fmt.Println("You currently have no instances.")
			return
		}

		for _, instance := range instances {
			fmt.Println(instance.String())
		}
	},
}
