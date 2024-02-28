package messages

import "fmt"

// Prints a message that informs the user that they are not logged in.
func DisplayNotLoggedInMessage() {
	fmt.Println("You are not logged in.")
	fmt.Println("Please log in to continue.")
}

// Prints a message that informs the user that their session has expired.
// If an error is provided, it will be printed as well.
func DisplaySessionExpiredMessage() {
	fmt.Println("Your session has expired or is not valid.")
	fmt.Println("Please log in again and retry.")
}