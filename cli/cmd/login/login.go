package login

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Login Command
//
// Allows the user to login to the IHP CLI using keycloak
var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Add short description", // TODO: Add short description
	Long:  "Add long description",  // TODO: Add long description
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("'ihp login' called")
	},
}