package utils

import "fmt"

func DisplayNotLoggedInMessage() {
	fmt.Println("You are not logged in.")
	fmt.Println("Please log in to continue.")
}

func DisplaySessionExpiredMessage() {
	fmt.Println("Your session has expired and you have been logged out.")
	fmt.Println("Please log in again to continue.")
}