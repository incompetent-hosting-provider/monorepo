package logout

import (
	"cli/cmd"
	"cli/internal/authentication"
	"cli/internal/utils"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	cmd.RootCmd.AddCommand(logoutCmd)
}

// Logout Command
//
// Allows the user to logout from the IHP CLI
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of the IHP CLI",
	Long:  "Log out of the IHP CLI. This will clear your current session.",
	Run: func(cmd *cobra.Command, args []string) {
		auth := authentication.GetCurrentAuthentication()
		if auth == nil {
			fmt.Println("You are not logged in.")
			return
		}

		server, err := utils.GetCallbackServer()
		if err != nil {
			fmt.Println("Failed to start a callback server to logout. Please try again.")
			return
		}
		defer server.Close()

		result := make(chan error)
		go authentication.PerformLogout(server, result)

		addr := server.Addr().String()
		redirectURL := fmt.Sprintf("http://%s", addr)
		url := authentication.DefaultKeycloakConfig.GetLogoutURL(redirectURL)

		err = utils.OpenBrowser(url)
		if err != nil {
			fmt.Println("Failed to open the browser. Please open the following URL manually:")
			fmt.Println(url)
		}

		err = <-result
		if err != nil {
			fmt.Println("Something went wrong during the logout process.")
			fmt.Println(err.Error())
			fmt.Println("Please try again. If the problem persists, please contact the support.")
		}

		fmt.Println("Successfully logged out.")
	},
}
