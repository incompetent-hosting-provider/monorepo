package instance

import (
	"cli/cmd"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	cmd.RootCmd.AddCommand(instanceCmd)
}

// Instance Command
//
// Displays information about a specific instance
var instanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "Add short description", // TODO: Add short description
	Long:  "Add long description",  // TODO: Add long description
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("'ihp instance' called")
	},
}
