package login

import (
	"cli/internal/authentication"
	"cli/internal/utils"
	"fmt"

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
		if auth, _ := authentication.GetCurrentAuthentication(); auth != nil {
			// We don't handle the error since this means that we are unable
			// to read the current authentication state. In this case we just
			// assume that the user is not logged in.

			fmt.Println("You are already logged in.")
			return
		}

		server, err := utils.GetCallbackServer()
		if err != nil {
			fmt.Println("Failed to start a callback server to login. Please try again.")
			return
		}
		defer server.Close()

		result := make(chan error)
		go authentication.PerformTokenExchange(server, result)

		addr := server.Addr().String()
		redirectURL := fmt.Sprintf("http://%s", addr)
		url := authentication.DefaultKeycloakConfig.GetLoginURL(redirectURL)

		err = utils.OpenBrowser(url)
		if err != nil {
			fmt.Println("Failed to open the browser. Please open the following URL manually:")
			fmt.Println(url)
		}

		err = <-result
		if err != nil {
			fmt.Println("Something went wrong during the login process.")
			fmt.Println(err.Error())
			fmt.Println("Please try again. If the problem persists, please contact the support.")
		}

		fmt.Println("Successfully logged in.")
	},
}
