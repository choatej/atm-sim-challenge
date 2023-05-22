package cmd

import (
	"agile-coder.com/atm-sim/internal"
	"fmt"
	"github.com/spf13/cobra"
)

// balanceCmd represents the balance command
var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "return the balance",
	Long:  `This command returns the account balance in US dollars`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return fmt.Errorf("the balance command does not take any parameters\n")
		}
		session := internal.GetSession()
		fmt.Printf("balance: $%.2f\n", internal.GetLedgerService().GetBalance(session.AccountId))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(balanceCmd)
}
