package instances

import (
	"cli/cmd"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	cmd.RootCmd.AddCommand(instancesCmd)
}

// Instances Command
//
// Lists all instances belonging to a user
var instancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "Add short description", // TODO: Add short description
	Long:  "Add long description",  // TODO: Add long description
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("'ihp instances' called")
	},
}
