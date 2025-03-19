package domainexceptions

import "fmt"

type InvalidAmountException struct {
	Amount int
}

func NewInvalidAmountException(amount int) *InvalidAmountException {
	return &InvalidAmountException{
		Amount: amount,
	}
}

func (e *InvalidAmountException) Error() string {
	return fmt.Sprintf("amount must be positive: %d", e.Amount)
}
