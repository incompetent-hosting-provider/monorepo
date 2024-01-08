package register

import (
	"cli/internal/services/authentication"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Register Command
//
// Allows the user to register for the IHP CLI using keycloak
var RegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register for the IHP CLI",
	Long: "Register for the IHP CLI via keycloak",
	Run: func(cmd *cobra.Command, args []string) {
		if err := authentication.Register(); err != nil {
			fmt.Println("Something went wrong while registering")
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Register successful!")
		auth, err := authentication.GetSessionToken()
		if err != nil {
			fmt.Println("Something went wrong while getting the session token after register")
			fmt.Println("Please try to login again")
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Session token: " + auth)
	},
}