package cmd

import (
	"cli/cmd/balance"
	"cli/cmd/instance"
	"cli/cmd/instances"
	"cli/cmd/login"
	"cli/cmd/logout"
	"cli/cmd/register"
	"cli/internal/services/authentication"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// The IHP CLI root command
//
// This command is the root command for the IHP CLI.
// If called without any subcommands, it will print
// an overview about the current user, their balance
// and their running instances.
var rootCmd = &cobra.Command{
	Use:   "ihp",
	Short: "Add short description", // TODO: Add short description
	Long: "Add long description", // TODO: Add long description
	Run: func(cmd *cobra.Command, args []string) {
		sessionToken, err := authentication.GetSessionToken()
		if err != nil {
			fmt.Println("Something went wrong while getting the session token after login")
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Session token: " + sessionToken)
	},
}

func init() {
	rootCmd.AddCommand(balance.BalanceCmd)
	rootCmd.AddCommand(instance.InstanceCmd)
	rootCmd.AddCommand(instances.InstancesCmd)
	rootCmd.AddCommand(login.LoginCmd)
	rootCmd.AddCommand(logout.LogoutCmd)
	rootCmd.AddCommand(register.RegisterCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

