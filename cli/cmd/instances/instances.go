package instances

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Instances Command
//
// Lists all instances belonging to a user
var InstancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "Add short description", // TODO: Add short description
	Long:  "Add long description",  // TODO: Add long description
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("'ihp instances' called")
	},
}
