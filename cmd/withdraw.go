package cmd

import (
	"agile-coder.com/atm-sim/internal"
	"fmt"
	"github.com/spf13/cobra"
)

// withdrawCmd represents the withdraw command
var withdrawCmd = &cobra.Command{
	Use:   "withdraw",
	Short: "withdraw funds",
	Long: `withdraw funds from the account
requires one parameter, the amount to withdraw
accounts are not allowed to overdraw so the requested amount must be less or equal to
the current account balance`,
	Run: func(cmd *cobra.Command, args []string) {
		var overdraftMessage string
		if len(args) != 1 {
			fmt.Println("withdraw takes one parameter - amount of the deposit")
			return
		}
		session := internal.GetSession()
		newBalance, err := internal.GetLedgerService().Withdraw(session.AccountId, args[0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if newBalance.WasOverdrawn {
			overdraftMessage = "You have been charged an overdraft fee of $5. "
		}
		fmt.Printf("Amount dispensed: $%.2f\n%sCurrent balance:%.2f\n", newBalance.AmountWithdrawn, overdraftMessage, newBalance.RemainingBalance)
	},
}

func init() {
	RootCmd.AddCommand(withdrawCmd)
}
