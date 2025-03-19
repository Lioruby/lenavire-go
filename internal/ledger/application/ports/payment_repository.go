package ports

import "lenavire/internal/ledger/domain/entities"

type PaymentRepository interface {
	Create(payment *entities.Payment) error
}
