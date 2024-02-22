package instance

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	instanceCmd.AddCommand(startCmd)
}

// Instances Start Command
//
// Allows the user to start an instance
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Add short description", // TODO: Add short description
	Long:  "Add long description",  // TODO: Add long description
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("'ihp instance start' called")
	},
}
