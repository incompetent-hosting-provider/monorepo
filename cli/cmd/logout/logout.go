package logout

import (
	"cli/internal/services/authentication"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Logout Command
//
// Allows the user to logout from the IHP CLI
var LogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of the IHP CLI",
	Long: "Log out of the IHP CLI. This will clear your current session.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := authentication.Logout(); err != nil {
			fmt.Println("Something went wrong while logging out")
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Logout successful!")
	},
}

func init() {}