package cmd

import (
	"agile-coder.com/atm-sim/internal"
	"fmt"
	"github.com/spf13/cobra"
)

// historyCmd represents the history command
var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "view transaction history",
	Long:  `shows a history of all deposits and withdrawals`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return fmt.Errorf("the history command does not take any parameters\n")
		}
		session := internal.GetSession()
		historyEntries := internal.GetLedgerService().GetHistory(session.AccountId)
		if len(historyEntries) == 0 {
			fmt.Println("No history found")
			return nil
		}
		fmt.Println("date\t\t\t\tamount\t\tbalance")
		for _, entry := range historyEntries {
			formattedDate := entry.Date.Format("2006-01-02 15:04:05Z")
			fmt.Printf("%s\t\t%.2f\t\t%.2f\n", formattedDate, entry.Amount, entry.Balance)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(historyCmd)
}
