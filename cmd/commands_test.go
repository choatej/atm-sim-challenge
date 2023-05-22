package cmd

import (
	"agile-coder.com/atm-sim/internal"
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestAuthorizeCmd(t *testing.T) {
	accountId := "jc123"
	testCases := []struct {
		name           string
		args           []string
		expectedOutput string
	}{
		{name: "no params", args: []string{}, expectedOutput: "authorize requires 2 parameters: account number, pin\n"},
		{name: "too many params", args: []string{accountId, "0000", "foo"}, expectedOutput: "authorize requires 2 parameters: account number, pin\n"},
		{name: "good pin", args: []string{accountId, "0000"}, expectedOutput: "jc123 successfully authorized.\n"},
		{name: "bad pin", args: []string{accountId, "1111"}, expectedOutput: "Authorization failed.\n"},
	}

	for _, test := range testCases {
		authService := internal.GetAuthorizationService()
		encryptedPin, err := internal.EncryptPin("0000")
		if err != nil {
			t.Fatal(err)
		}
		authService.SetAuthData(map[string]internal.EncryptedPin{accountId: encryptedPin})
		capturedText, err := runAndGetOutput(authorizeCmd, "authorize", test.args)
		if err != nil {
			assert.Equal(t, test.expectedOutput, err.Error())
		} else {
			// Assert the captured output against expected output
			assert.Equal(t, test.expectedOutput, capturedText, "%s failed. expected: %s got: %s", test.name, test.expectedOutput, capturedText)
		}
	}
}

func TestBalanceCmd(t *testing.T) {
	accountId := "jc123"
	testCases := []struct {
		name           string
		args           []string
		balance        float64
		expectedOutput string
	}{
		{name: "too many params", args: []string{"foo"}, balance: 40.00, expectedOutput: "the balance command does not take any parameters\n"},
		{name: "value", args: []string{}, balance: 40.00, expectedOutput: "balance: $40.00\n"},
	}

	for _, test := range testCases {
		session := internal.GetSession()
		session.IsAuthenticated = true
		session.AccountId = accountId

		ledger := internal.GetLedgerService()
		ledger.SetInitialBalances(10000, map[string]float64{
			accountId: test.balance,
		})
		capturedText, err := runAndGetOutput(balanceCmd, "balance", test.args)
		if err != nil {
			assert.Equal(t, test.expectedOutput, err.Error())
		} else {
			// Assert the captured output against expected output
			assert.Equal(t, test.expectedOutput, capturedText, "%s failed. expected: %s got: %s", test.name, test.expectedOutput, capturedText)
		}
	}
}

func TestDepositCmd(t *testing.T) {

	testCases := []struct {
		name           string
		args           []string
		expectedOutput string
	}{
		{name: "good deposit",
			args:           []string{"20.00"},
			expectedOutput: "Current balance: $60.00\n"},
		{name: "no args",
			args:           []string{},
			expectedOutput: "deposit takes one parameter - amount of the deposit\n"},
		{name: "too many args",
			args:           []string{"20.00", "30.00"},
			expectedOutput: "deposit takes one parameter - amount of the deposit\n"},
		{name: "not a number",
			args:           []string{"xyzzy"},
			expectedOutput: "invalid input: invalid number format xyzzy\n"},
	}

	for _, test := range testCases {
		accountId := "jc123"
		session := internal.GetSession()
		session.IsAuthenticated = true
		session.AccountId = accountId

		ledger := internal.GetLedgerService()
		ledger.SetInitialBalances(10000, map[string]float64{
			accountId: 40.00,
		})

		// Get the captured output
		capturedText, err := runAndGetOutput(depositCmd, "deposit", test.args)
		if err != nil {
			assert.Equal(t, test.expectedOutput, err.Error())
		} else {
			// Assert the captured output against expected output
			assert.Equal(t, test.expectedOutput, capturedText, "%s failed. expected: %s got: %s", test.name, test.expectedOutput, capturedText)
		}
	}
}

func TestLogoutCmd(t *testing.T) {
	accountId := "jc123"
	session := internal.GetSession()
	session.IsAuthenticated = true
	session.AccountId = accountId
	capturedText, err := runAndGetOutput(logoutCmd, "logout", []string{})
	if err != nil {
		t.Fatal(err)
	}
	expectedOutput := fmt.Sprintf("Account %s logged out.\n", accountId)
	assert.Equal(t, expectedOutput, capturedText)
}

func TestLogoutCmdNoUser(t *testing.T) {
	session := internal.GetSession()
	session.IsAuthenticated = false
	capturedText, err := runAndGetOutput(logoutCmd, "logout", []string{})
	if err != nil {
		t.Fatal(err)
	}
	expectedOutput := "No account is currently authorized.\n"
	assert.Equal(t, expectedOutput, capturedText)
}

func TestHistoryCmd(t *testing.T) {

	accountId := "jc456"

	testCases := []struct {
		name           string
		args           []string
		expectedOutput string
	}{
		{name: "too many params", args: []string{"foo"}, expectedOutput: "the history command does not take any parameters\n"},
		{name: "no history", args: []string{}, expectedOutput: "No history found\n"},
		//{name: "value", args: []string{}, transaction: ledger.Deposit, expectedOutput: "date\\t\\t\\tamount\\t\\tbalance\\n2023-05-28 14:27:10Z\\t\\t20.00\\t\\t60.00\\n2023-05-28 14:27:10Z\\t\\t20.00\\t\\t60.00\\n\n"},
	}

	for _, test := range testCases {
		session := internal.GetSession()
		session.IsAuthenticated = true
		session.AccountId = accountId

		capturedText, err := runAndGetOutput(historyCmd, "history", test.args)
		if err != nil {
			assert.Equal(t, test.expectedOutput, err.Error())
		} else {
			// Assert the captured output against expected output
			assert.Equal(t, test.expectedOutput, capturedText, "%s failed. expected: %s got: %s", test.name, test.expectedOutput, capturedText)
		}
	}
}

func TestHistoryDetail(t *testing.T) {
	accountId := "jc678"
	session := internal.GetSession()
	session.IsAuthenticated = true
	session.AccountId = accountId

	ledger := internal.GetLedgerService()
	ledger.SetInitialBalances(10000, map[string]float64{
		accountId: 40.00,
	})
	_, err := ledger.Deposit(accountId, "40.00")
	if err != nil {
		t.Fatal(err.Error())
	}
	capturedText, err := runAndGetOutput(historyCmd, "history", []string{})
	if err != nil {
		t.Error("history command failed")
	} else {
		lines := strings.Split(capturedText, "\n")
		assert.Equal(t, "date\t\t\t\tamount\t\tbalance", lines[0])
		historyLine := strings.Split(lines[1], "\t\t")
		assert.Equal(t, "40.00", historyLine[1])
		assert.Equal(t, "80.00", historyLine[2])
	}
}

func TestUnauthorizedCmd(t *testing.T) {
	session := internal.GetSession()
	session.IsAuthenticated = false
	_, err := runAndGetOutput(depositCmd, "deposit", []string{})
	if err == nil {
		t.Fatal("should have had an error")
	}
	assert.Equal(t, err.Error(), "Authorization required.\n")

}

func TestWithdrawCommand(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		expectedOutput string
		startingCash   float64
		endingCash     float64
	}{
		{name: "good withdrawal",
			args:           []string{"20.00"},
			expectedOutput: "Amount dispensed: $20.00\nCurrent balance:20.00\n",
			startingCash:   10000.00,
			endingCash:     9980.00,
		},

		{name: "overdraft",
			args:           []string{"60.00"},
			expectedOutput: "Amount dispensed: $60.00\nYou have been charged an overdraft fee of $5. Current balance:-25.00\n",
			startingCash:   10000.00,
			endingCash:     9940.00,
		},
		{name: "no args",
			args:           []string{},
			expectedOutput: "withdraw takes one parameter - amount of the deposit\n",
			startingCash:   10000.00,
			endingCash:     10000.00,
		},
		{name: "too many args",
			args:           []string{"20.00", "30.00"},
			expectedOutput: "withdraw takes one parameter - amount of the deposit\n",
			startingCash:   10000.00,
			endingCash:     10000.00,
		},
		{name: "not a number",
			args:           []string{"xyzzy"},
			expectedOutput: "invalid input: invalid number format xyzzy\n",
			startingCash:   10000.00,
			endingCash:     10000.00,
		},
	}

	for _, test := range testCases {
		accountId := "jc123"
		session := internal.GetSession()
		session.IsAuthenticated = true
		session.AccountId = accountId
		availableCash := 10000.00

		ledger := internal.GetLedgerService()
		ledger.SetInitialBalances(availableCash, map[string]float64{
			accountId: 40.00,
		})

		// Get the captured output
		capturedText, err := runAndGetOutput(withdrawCmd, "withdraw", test.args)
		if err != nil {
			assert.Equal(t, test.expectedOutput, err.Error())
		} else {
			// Assert the captured output against expected output
			assert.Equal(t, test.expectedOutput, capturedText, "%s failed. expected: %s got: %s", test.name, test.expectedOutput, capturedText)
			assert.Equal(t, test.endingCash, ledger.GetAvailableCash())
		}
	}
}

func runAndGetOutput(cmd *cobra.Command, commandName string, args []string) (string, error) {
	internal.InitLogger("", true)
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	// Start capturing stdout
	var capturedOutput bytes.Buffer
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, _ = io.Copy(&capturedOutput, r)
	}()

	// Execute the command
	RootCmd.AddCommand(cmd)
	cmd.SetArgs(args)
	RootCmd.SetArgs(append([]string{commandName}, args...))
	err := cmd.Execute()
	if err != nil {
		return "", err
	}
	// Close the writer and wait for the goroutine to finish
	err = w.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close: %s", err.Error())
	}
	wg.Wait()
	// Restore the original stdout
	os.Stdout = oldStdout

	// Get the captured output
	return capturedOutput.String(), nil
}
