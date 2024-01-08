package login

import (
	"cli/internal/services/authentication"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Login Command
//
// Allows the user to login to the IHP CLI using keycloak
var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the IHP CLI",
	Long:  "Login to the IHP CLI via Keycloak",
	Run: func(cmd *cobra.Command, args []string) {
		if err := authentication.Login(); err != nil {
			fmt.Println("Something went wrong while logging out")
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Login successful!")
		auth, err := authentication.GetSessionToken()
		if err != nil {
			fmt.Println("Something went wrong while getting the session token after login")
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Session token: " + auth)
	},
}