package cmd

import (
	"agile-coder.com/atm-sim/internal"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

// authorizeCmd represents the authorize command
var authorizeCmd = &cobra.Command{
	Use:   "authorize",
	Short: "authorizing the user to perform transactions",
	Long: `Authorizes the user to perform account activities such as
- get balance
- deposit
- withdrawal
- view transaction history
The command takes two inputs - account number and pin`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// parameter validation
		params := []string{"account number", "pin"}
		if len(args) != len(params) {
			return fmt.Errorf("%s requires %d parameters: %s\n", cmd.Name(), len(params), strings.Join(params, ", "))
		}

		accountId := args[0]
		pin := args[1]
		return authCommand(accountId, pin)
	},
}

func authCommand(accountId string, pin string) error {

	authService := internal.GetAuthorizationService()
	if authService == nil {
		return fmt.Errorf("authService is nil\n")
	}
	ok, err := authService.Authenticate(accountId, pin)

	if ok {
		fmt.Printf("%s successfully authorized.\n", accountId)
		internal.Logger.Printf("successful login for %s\n", accountId)
		session := internal.GetSession()
		session.IsAuthenticated = true
		session.AccountId = accountId
	} else {
		fmt.Println("Authorization failed.")
		internal.Logger.Printf("invalid login attempt for %s\n", accountId)
	}
	return err
}

func init() {
	RootCmd.AddCommand(authorizeCmd)
}
