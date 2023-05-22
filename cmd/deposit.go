package cmd

import (
	"agile-coder.com/atm-sim/internal"
	"fmt"
	"github.com/spf13/cobra"
)

// depositCmd represents the deposit command
var depositCmd = &cobra.Command{
	Use:   "deposit",
	Short: "make a deposit",
	Long: `Deposit funds in the account
required parameter: amount to deposit in dollars and cents`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("deposit takes one parameter - amount of the deposit")
			return
		}
		session := internal.GetSession()
		newBalance, err := internal.GetLedgerService().Deposit(session.AccountId, args[0])
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("Current balance: $%.2f\n", newBalance)
		}
	},
}

func init() {
	RootCmd.AddCommand(depositCmd)
}
