package internal

import (
	"errors"
	"reflect"
	"testing"
)

const account = "jc123"

func TestDeposit(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected float64
		err      error
	}{
		{name: "not a number", value: "xyzzy", expected: 50.00, err: &InvalidAmountError{message: "invalid number format xyzzy"}},
		{name: "leading dollar sign", value: "$25.95", expected: 75.95, err: nil},
		{name: "too many decimals", value: "25.222", expected: 50.00, err: &InvalidAmountError{message: "invalid number format 25.222"}},
		{name: "good value", value: "150.15", expected: 200.15, err: nil},
	}

	for _, test := range tests {
		InitLogger("", true)
		var balance = 50.00
		const accountId = "jc123"
		testLedger := GetLedgerService()
		testLedger.SetInitialBalances(0, map[string]float64{
			accountId: balance,
		})
		_, err := testLedger.Deposit(account, test.value)
		newBalance := testLedger.GetBalance(accountId)
		if newBalance != test.expected {
			t.Errorf("%s: incorrect balance after deposit expected %.2f got %.2f\n", test.name, test.expected, newBalance)
		}
		if err != nil && !errors.Is(err, test.err) {
			t.Errorf("%s: unexpected error %+v\n", test.name, err)
		}
		if test.err != nil && err == nil {
			t.Errorf("%s: expected error %s but none was found\n", test.name, reflect.TypeOf(err))
		}
	}
}

func TestWithdraw(t *testing.T) {
	tests := []struct {
		name          string
		value         string
		expected      float64
		err           error
		availableCash float64
	}{
		{name: "not a number", value: "xyzzy", expected: 50.00, err: &InvalidAmountError{message: "invalid number format xyzzy"}, availableCash: 500},
		{name: "leading dollar sign", value: "$20.00", expected: 30.00, err: nil, availableCash: 500},
		{name: "too many decimals", value: "25.222", expected: 50.00, err: &InvalidAmountError{message: "invalid number format 25.222"}, availableCash: 500},
		{name: "overdraw", value: "80.00", expected: -35.00, err: nil, availableCash: 500},
		{name: "empty machine", value: "20.00", expected: 50.00, err: &NoMoneyLeftError{}, availableCash: 0},
		{name: "partial fulfillment", value: "40.00", expected: 30.00, err: nil, availableCash: 20.00},
		{name: "not a multiple of 20", value: "25.00", expected: 50.00, err: &InvalidAmountError{message: "Withdrawals must be in units of $20."}, availableCash: 500},
		{name: "good value", value: "20.00", expected: 30.00, err: nil, availableCash: 500},
	}
	InitLogger("", true)
	for _, test := range tests {
		var balance = 50.00
		const accountId = "jc123"
		var testLedger = Ledger{
			availableCash: test.availableCash,
			balances: map[string]float64{
				accountId: balance,
			},
		}
		result, err := testLedger.Withdraw(account, test.value)
		if result != nil {
			newBalance := result.RemainingBalance
			if newBalance != test.expected {
				t.Errorf("%s: incorrect balance after withdrawl expected %.2f got %.2f\n", test.name, test.expected, newBalance)
			}
		}
		if err != nil && !errors.Is(err, test.err) {
			t.Errorf("%s: unexpected error %+v\n", test.name, err)
		}
		if test.err != nil && err == nil {
			t.Errorf("%s: expected error %s but none was found\n", test.name, reflect.TypeOf(err))
		}
	}
}

func TestAlreadyOverdrawn(t *testing.T) {
	accountId := "jc123"
	InitLogger("", true)
	ledger := GetLedgerService()
	ledger.SetInitialBalances(500, map[string]float64{
		accountId: -20,
	})
	_, err := ledger.Withdraw(accountId, "20.00")
	if err == nil {
		t.Errorf("expected an error but did not get it")
	} else {
		if !errors.Is(err, &OverdrawnError{}) {
			t.Errorf("error was not an OverdrawnError")
		}
	}
}

func TestHistory(t *testing.T) {
	accountId := "jc456"
	InitLogger("", true)
	ledger := GetLedgerService()
	ledger.histories = map[string][]LedgerHistoryEntry{}
	ledger.SetInitialBalances(5000, map[string]float64{
		accountId: 0,
	})
	_, err := ledger.Deposit(accountId, "20.00")
	if err != nil {
		t.Fatal("failed to deposit")
	}
	_, err = ledger.Withdraw(accountId, "40.00")
	if err != nil {
		t.Fatal("failed to withdraw")
	}
	history := ledger.GetHistory(accountId)
	if len(history) != 3 {
		t.Fatalf("expected 3 history entries but got %d %v", len(history), history)
	}
	depositHistoryRecord := history[0]
	compareHistoryEntries(LedgerHistoryEntry{Amount: 20.00, Balance: 20.00}, depositHistoryRecord, t)

	withdrawalHistoryRecord := history[1]
	compareHistoryEntries(LedgerHistoryEntry{Amount: -40.00, Balance: -20.00}, withdrawalHistoryRecord, t)

	overdraftHistoryRecord := history[2]
	compareHistoryEntries(LedgerHistoryEntry{Amount: -5.00, Balance: -25.00}, overdraftHistoryRecord, t)
}

func compareHistoryEntries(expected LedgerHistoryEntry, actual LedgerHistoryEntry, t *testing.T) {
	if actual.Amount != expected.Amount {
		t.Errorf("amount mismatch: expected: %.2f, got: %.2f", expected.Amount, actual.Amount)
	}
	if actual.Balance != expected.Balance {
		t.Errorf("balance mismatch: expected: %.2f, got: %.2f", expected.Balance, actual.Balance)
	}
}
