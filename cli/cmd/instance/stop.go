package instance

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	instanceCmd.AddCommand(stopCmd)
}

// Instances Create Command
//
// Allows the user to stop an instance
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Add short description", // TODO: Add short description
	Long:  "Add long description",  // TODO: Add long description
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("'ihp instance stop' called")
	},
}
