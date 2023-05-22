package internal

import "fmt"

// InvalidInputError is used when an input fails input validation
type InvalidInputError struct {
	message string
}

func (e *InvalidInputError) Error() string {
	return fmt.Sprintf("invalid input: %s", e.message)
}

func (e *InvalidInputError) Is(target error) bool {
	other, ok := target.(*InvalidInputError)
	if !ok {
		return false
	}
	return e.message == other.message
}

// InvalidAmountError is used when an amount is logically invalid
type InvalidAmountError struct {
	message string
}

func (e *InvalidAmountError) Error() string {
	return fmt.Sprintf("invalid input: %s", e.message)
}

func (e *InvalidAmountError) Is(target error) bool {
	other, ok := target.(*InvalidAmountError)
	if !ok {
		return false
	}
	return e.message == other.message
}

// InsufficientFundsError is used when a withdrawal is attempted but there are insufficient funds to complete it
type InsufficientFundsError struct {
}

func (e *InsufficientFundsError) Is(target error) bool {
	_, ok := target.(*InsufficientFundsError)
	return ok
}

func (e *InsufficientFundsError) Error() string {
	return "Insufficient funds. Withdrawal amount exceeds account balance"
}

// NoMoneyLeftError is used when the machine is empty and a withdrawal is attempted
type NoMoneyLeftError struct {
}

func (e *NoMoneyLeftError) Error() string {
	return "Unable to process your withdrawal at this time."
}

// OverdrawnError is used when an account is already overdrawn and a withdrawal is attempted
type OverdrawnError struct {
}

func (e *OverdrawnError) Error() string {
	return "Your account is overdrawn! You may not make withdrawals at this time."
}

func (e *OverdrawnError) Is(target error) bool {
	_, ok := target.(*OverdrawnError)
	return ok
}
