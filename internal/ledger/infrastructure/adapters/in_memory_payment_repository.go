package adapters

import "lenavire/internal/ledger/domain/entities"

type InMemoryPaymentRepository struct {
	Payments []*entities.Payment
}

func NewInMemoryPaymentRepository() *InMemoryPaymentRepository {
	return &InMemoryPaymentRepository{
		Payments: []*entities.Payment{},
	}
}

func (r *InMemoryPaymentRepository) Create(payment *entities.Payment) error {
	r.Payments = append(r.Payments, payment)
	return nil
}
