package cmd

import (
	"cli/cmd/balance"
	"cli/cmd/instance"
	"cli/cmd/instances"
	"cli/cmd/login"
	"cli/cmd/logout"
	"cli/cmd/register"
	"cli/internal/authentication"
	"cli/internal/backend"
	"cli/internal/utils"
	_ "embed"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(balance.BalanceCmd)
	RootCmd.AddCommand(instance.InstanceCmd)
	RootCmd.AddCommand(instances.InstancesCmd)
	RootCmd.AddCommand(login.LoginCmd)
	RootCmd.AddCommand(logout.LogoutCmd)
	RootCmd.AddCommand(register.RegisterCmd)
}


var (
	//go:embed banner.txt
	banner string
)

// The IHP CLI root command
//
// This command is the root command for the IHP CLI.
// If called without any subcommands, it will print
// an overview about the current user, their balance
// and their running instances.
var RootCmd = &cobra.Command{
	Use:   "ihp",
	Short: "Add short description", // TODO: Add short description
	Long:  "Add long description",  // TODO: Add long description
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(banner)
		fmt.Println("Version: 0.1.0")
		fmt.Println()

		tokens := authentication.GetCurrentAuthentication()
		if tokens == nil {
			displayNotLoggedInMessage()
			return
		}

		userInfo, err := backend.GetUserInfo(*tokens)
		if err != nil {
			if errors.Is(err, backend.ErrNotAuthenticated) {
				// Refresh tokens and try again
				newTokens, err := authentication.RefreshTokens()
				if err != nil {
					displayUnableToGetUserInfoMessage(err)
					os.Exit(1)
				}

				if newTokens == nil {
					displaySessionExpiredMessage()
					return
				}

				userInfo, err = backend.GetUserInfo(*newTokens)
				if err != nil {
					displayUnableToGetUserInfoMessage(err)
					os.Exit(1)
				}
			} else {
				displayUnableToGetUserInfoMessage(err)
				os.Exit(1)
			}
		}

		fmt.Println("Welcome", userInfo.Email)
		fmt.Println("Your balance is", userInfo.Balance)
	},
}

func displayNotLoggedInMessage() {
	fmt.Println("You are not logged in.")
	fmt.Println("Please log in to continue.")
}

func displaySessionExpiredMessage() {
	fmt.Println("Your session has expired and you have been logged out.")
	fmt.Println("Please log in again to continue.")
}

func displayUnableToGetUserInfoMessage(err error) {
	fmt.Println("Unable to get your user info:", err)
	fmt.Println("Please try to log out and log in again.")
	fmt.Println("If the problem persists, please contact support.")
}
