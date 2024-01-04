package register

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Register Command
//
// Allows the user to register for the IHP CLI using keycloak
var RegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Add short description", // TODO: Add short description
	Long: "Add long description", // TODO: Add long description
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("'ihp register' called")
	},
}

func init() {}