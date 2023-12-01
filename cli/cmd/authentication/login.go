package authentication

import (
	"fmt"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
)

// Login Command
//
// Allows the user to login to the IHP CLI
var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Add short description", // TODO: Add short description
	Long: "Add long description", // TODO: Add long description
	Run: func(cmd *cobra.Command, args []string) {
		var username, password string = loginPrompt()
		if err := login(username, password); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("You are now logged in as", currentUser.username)
			fmt.Println("Your user id is", currentUser.id)		
		}
	},
}

func init() {
}

func loginPrompt() (string, string) {
	fmt.Println("Please enter your username:")
	var username string
	fmt.Scanln(&username)

	fmt.Println("Please enter your password:")
	var password string
	if pass, error := gopass.GetPasswdMasked(); error != nil {
		fmt.Println("Error: ", error)
		return registerPrompt()
	} else {
		password = string(pass)
	}

	return username, password
}