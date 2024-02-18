package cmd

import (
	"cli/cmd/balance"
	"cli/cmd/instance"
	"cli/cmd/instances"
	"cli/cmd/login"
	"cli/cmd/logout"
	"cli/cmd/register"
	_ "embed"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	//go:embed banner.txt
	banner string
)

func init() {
	rootCmd.AddCommand(balance.BalanceCmd)
	rootCmd.AddCommand(instance.InstanceCmd)
	rootCmd.AddCommand(instances.InstancesCmd)
	rootCmd.AddCommand(login.LoginCmd)
	rootCmd.AddCommand(logout.LogoutCmd)
	rootCmd.AddCommand(register.RegisterCmd)
}

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
		fmt.Println(banner)
		fmt.Println("Version: 0.1.0")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

