package internal

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"
)

/*
moneyPattern defines the valid formats for money in commands. The regex requires that the money:
  - starts with either a $ or a non-zero digit
  - decimals are defined using the US standard of using a '.'
  - either has exactly zero or two decimal places
  - commas are not allowed
  - no magnitude notations (K for thousands for instance) are allowed
*/
const moneyPattern = "^(\\$)?([1-9]\\d*(\\.\\d\\d)?)$"
const overdraftFee = 5.00

// LedgerHistoryEntry holds the transaction history for accounts
type LedgerHistoryEntry struct {
	Date    time.Time
	Amount  float64
	Balance float64
}

// WithdrawResult since we need multiple pieces of info for a withdrawal, wrap it in a struct
type WithdrawResult struct {
	AmountWithdrawn  float64
	RemainingBalance float64
	WasOverdrawn     bool
}

// Ledger holds the account balances and history
type Ledger struct {
	// amount able to be dispensed
	availableCash float64
	// map of account # to balance
	balances  map[string]float64
	histories map[string][]LedgerHistoryEntry
}

// the shared Ledger instance
var ledger = &Ledger{}

func GetLedgerService() *Ledger {
	return ledger
}

func (ledger *Ledger) GetAvailableCash() float64 {
	return ledger.availableCash
}

// SetInitialBalances sets the starting balances in the Ledger
func (ledger *Ledger) SetInitialBalances(availableCash float64, balances map[string]float64) {
	ledger.balances = balances
	ledger.availableCash = availableCash
}

// GetBalance returns the current balance for a given account
func (ledger *Ledger) GetBalance(account string) (balance float64) {
	return ledger.balances[account]
}

// StringToMoney validates that a given string is an allowed money value and returns the amount as a float64
func StringToMoney(input string) (float64, error) {
	regex := regexp.MustCompile(moneyPattern)
	matches := regex.FindStringSubmatch(input)
	if matches != nil {
		data, err := strconv.ParseFloat(matches[2], 64)
		if err != nil {
			return 0, err
		}
		return data, nil
	}
	return 0, &InvalidAmountError{message: fmt.Sprintf("invalid number format %s", input)}
}

// Deposit adds funds to a given account
func (ledger *Ledger) Deposit(accountId string, amount string) (float64, error) {
	currentBalance := ledger.balances[accountId]
	dollarAmount, err := StringToMoney(amount)
	if err != nil {
		return currentBalance, err
	}
	newValue := currentBalance + dollarAmount
	ledger.balances[accountId] = newValue
	ledger.addHistory(accountId, dollarAmount, newValue)
	return newValue, nil
}

// Withdraw removes funds from a given account
func (ledger *Ledger) Withdraw(accountId string, amount string) (*WithdrawResult, error) {
	currentBalance := ledger.balances[accountId]

	// customer is already overdrawn
	if currentBalance <= 0 {
		return &WithdrawResult{RemainingBalance: currentBalance, WasOverdrawn: true}, &OverdrawnError{}
	}

	// the machine is empty
	if ledger.availableCash == 0 {
		return &WithdrawResult{RemainingBalance: currentBalance}, &NoMoneyLeftError{}
	}

	// covert the request to a float64
	dollarAmount, err := StringToMoney(amount)
	if err != nil {
		return &WithdrawResult{RemainingBalance: currentBalance}, err
	}

	// can only dispense in units of $20
	if math.Mod(dollarAmount, 20) != 0 {
		return &WithdrawResult{RemainingBalance: currentBalance}, &InvalidAmountError{message: "Withdrawals must be in units of $20."}
	}
	result := WithdrawResult{}
	// can only dispense partial amount
	if dollarAmount > ledger.availableCash {
		dollarAmount = ledger.availableCash
	}
	newValue := currentBalance - dollarAmount
	ledger.addHistory(accountId, dollarAmount*-1, newValue)
	if newValue < 0 {
		newValue = newValue - overdraftFee
		ledger.addHistory(accountId, overdraftFee*-1, newValue)
		result.WasOverdrawn = true
	}
	ledger.balances[accountId] = newValue
	ledger.availableCash = ledger.availableCash - dollarAmount
	result.RemainingBalance = newValue
	result.AmountWithdrawn = dollarAmount
	return &result, nil
}

// addHistory updates the ledger history with a new transacion
func (ledger *Ledger) addHistory(accountId string, amount float64, balance float64) {
	// lazy initialization of Ledger.histories
	newEntry := LedgerHistoryEntry{Date: time.Now(), Amount: amount, Balance: balance}
	Logger.Printf("adding history for %s %.2f\n", accountId, newEntry.Amount)
	if ledger.histories == nil {
		ledger.histories = map[string][]LedgerHistoryEntry{}
	}
	entry, ok := ledger.histories[accountId]
	if ok {
		Logger.Printf("appending history to account %s\n", accountId)
		ledger.histories[accountId] = append(entry, newEntry)
	} else {
		Logger.Printf("starting history for %s\n", accountId)
		ledger.histories[accountId] = []LedgerHistoryEntry{newEntry}
	}
}

// GetHistory returns the transaction history for a given account
func (ledger *Ledger) GetHistory(accountId string) []LedgerHistoryEntry {
	retVal := ledger.histories[accountId]
	Logger.Printf("returning %d histories\n", len(retVal))
	return retVal
}
