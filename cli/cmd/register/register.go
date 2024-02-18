package register

import (
	"cli/internal/authentication"
	"cli/internal/utils"
	"fmt"

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
		if auth, _ := authentication.GetCurrentAuthentication(); auth != nil {
			// We don't handle the error since this means that we are unable
			// to read the current authentication state. In this case we just
			// assume that the user is not logged in.

			fmt.Println("You are currently logged in. Please log out first!")
			return
		}

		server, err := utils.GetCallbackServer()
		if err != nil {
			fmt.Println("Failed to start a callback server to register. Please try again.")
			return
		}
		defer server.Close()

		result := make(chan error)
		go authentication.PerformTokenExchange(server, result)

		addr := server.Addr().String()
		redirectURL := fmt.Sprintf("http://%s", addr)
		url:= authentication.DefaultKeycloakConfig.GetRegisterURL(redirectURL)
		
		err = utils.OpenBrowser(url)
		if err != nil {
			fmt.Println("Failed to open the browser. Please open the following URL manually:")
			fmt.Println(url)
		}

		err = <- result
		if err != nil {
			fmt.Println("Something went wrong during the login process.")
			fmt.Println(err.Error())
			fmt.Println("Please try again. If the problem persists, please contact the support.")
		}
		
		fmt.Println("Successfully registered.")
		fmt.Println("You can now use the IHP-CLI!.")
	},
}