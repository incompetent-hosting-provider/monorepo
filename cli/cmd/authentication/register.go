package authentication

import (
	"fmt"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
)

// Logout Command
//
// Allows the user to logout of the IHP CLI
var RegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Add short description", // TODO: Add short description
	Long: "Add long description", // TODO: Add long description
	Run: func(cmd *cobra.Command, args []string) {
		username, password := registerPrompt()
		if err := register(username, password); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("You have successfully registered")
			fmt.Println("You are now logged in as", currentUser.username)
			fmt.Println("Your user id is", currentUser.id)		
		}
	},
}

func init() {
}

/// Gets the username and password from the user
///
/// Returns the username and password
func registerPrompt() (string, string) {
	fmt.Println("Please enter your desired username:")
	var username string
	fmt.Scanln(&username)

	fmt.Println("Please enter your desired password:")
	var password string
	if pass, error := gopass.GetPasswdMasked(); error != nil {
		fmt.Println("Error: ", error)
		return registerPrompt()
	} else {
		password = string(pass)
	}


	fmt.Println("Please enter the provided password again:")
	var passwordConfirmation string
	if pass, error := gopass.GetPasswdMasked(); error != nil {
		fmt.Println("Error: ", error)
		return registerPrompt()
	} else {
		passwordConfirmation = string(pass)
	}

	if password != passwordConfirmation {
		fmt.Println("The provided passwords do not match")
		return registerPrompt()
	}

	return username, password
}