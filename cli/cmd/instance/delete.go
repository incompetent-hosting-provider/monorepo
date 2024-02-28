package instance

import (
	"cli/internal/authentication"
	"cli/internal/backend"
	"cli/internal/messages"
	"errors"
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// Instances Delete Command
//
// Allows the user to delete an instance
var deleteCmd = &cobra.Command{
	Use:   "delete [instance_id]",
	Short: "Deletes an instance.", 
	Long:  "Deletes an instance.",
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tokens := authentication.GetCurrentAuthentication()
		if tokens == nil {
			messages.DisplayNotLoggedInMessage()
			return
		}

		userInstances, err := backend.DefaultBackendClient.GetUserInstances(tokens.AccessToken, true)
		if err != nil {
			handleDeleteError(err)
			return
		}

		if len(userInstances) == 0 {
			fmt.Println("You don't have any instances to delete.")
			return
		}

		instanceToDeletePrompt := promptui.Select{
			Label: "Instance to delete",
			Items: userInstances,
		}

		instanceIndex, _, err := instanceToDeletePrompt.Run()
		if err != nil {
			handleDeleteError(err)
			return
		}
		
		err = backend.DefaultBackendClient.DeleteInstance(tokens.AccessToken, userInstances[instanceIndex].ID, true)
		if err != nil {
			handleDeleteError(err)
			return	
		}

		fmt.Printf("Instance %s - %s successfully deleted.", userInstances[instanceIndex].ID, userInstances[instanceIndex].Name)
	},
}

func handleDeleteError(err error) {
	if errors.Is(err, backend.ErrNotAuthenticated) {
		messages.DisplaySessionExpiredMessage()
		return
	} else {
		fmt.Println("Unable to delete your instance:", err)
		fmt.Println("Please try again later.")
		return
	}
}
