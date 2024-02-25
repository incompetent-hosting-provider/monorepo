package main

import (
	"cli/cmd"
	"fmt"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println("Something went wrong during the execution of the program.")
		fmt.Println(err.Error())
		fmt.Println("Please try to reinstall the IHP CLI. If the problem persists, please contact the support.")
	}
}
