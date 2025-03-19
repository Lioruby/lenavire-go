package valuesobjects

import "lenavire/internal/ledger/domain/exceptions"

type Amount struct {
	Value int
}

func NewAmount(value int) (Amount, error) {
	if value < 0 {
		return Amount{}, domainexceptions.NewInvalidAmountException(value)
	}
	return Amount{Value: value}, nil
}
