package cmd

import (
	"agile-coder.com/atm-sim/internal"
	"fmt"
	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "log out the user",
	Long:  `Logs out the current user. To perform any functions the user will need to re-authorize`,
	Run: func(cmd *cobra.Command, args []string) {
		session := internal.GetSession()
		if session.IsAuthenticated {
			currentAccountId := session.AccountId
			session.IsAuthenticated = false
			session.AccountId = ""
			fmt.Printf("Account %s logged out.\n", currentAccountId)
		} else {
			fmt.Println("No account is currently authorized.")
		}
	},
}

func init() {
	RootCmd.AddCommand(logoutCmd)
}
