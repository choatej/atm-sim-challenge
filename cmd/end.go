package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// endCmd represents the end command
var endCmd = &cobra.Command{
	Use:   "end",
	Short: "Exit the application",
	Long:  "exit the application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("exiting...")
		os.Exit(0)
	},
}

func init() {
	RootCmd.AddCommand(endCmd)
}
