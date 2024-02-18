package instance

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Instances Create Command
//
// Runs the create instance prompt so the user can create a new instance
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Add short description", // TODO: Add short description
	Long:  "Add long description",  // TODO: Add long description
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("'ihp instance create' called")
	},
}
