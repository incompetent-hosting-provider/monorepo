package instance

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	InstanceCmd.AddCommand(createCmd)
	InstanceCmd.AddCommand(deleteCmd)
	InstanceCmd.AddCommand(startCmd)
	InstanceCmd.AddCommand(stopCmd)
}

// Instance Command
//
// Displays information about a specific instance
var InstanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "Add short description", // TODO: Add short description
	Long:  "Add long description",  // TODO: Add long description
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("'ihp instance' called")
	},
}
