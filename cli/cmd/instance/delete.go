package instance

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	instanceCmd.AddCommand(deleteCmd)
}

// Instances Delete Command
//
// Allows the user to delete an instance
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Add short description", // TODO: Add short description
	Long:  "Add long description",  // TODO: Add long description
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("'ihp instance delete' called")
	},
}
