package instance

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Instances Start Command
//
// Allows the user to start an instance
var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Add short description", // TODO: Add short description
	Long:  "Add long description",  // TODO: Add long description
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("'ihp instance start' called")
	},
}
