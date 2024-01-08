package logout

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Logout Command
//
// Allows the user to logout from the IHP CLI
var LogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Add short description", // TODO: Add short description
	Long: "Add long description", // TODO: Add long description
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("'ihp logout' called")
	},
}