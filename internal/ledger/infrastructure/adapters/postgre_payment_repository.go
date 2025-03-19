package adapters

import (
	"lenavire/internal/ledger/domain/entities"
	"lenavire/internal/ledger/infrastructure/database/mappers"

	"gorm.io/gorm"
)

type PostgrePaymentRepository struct {
	db     *gorm.DB
	mapper *mappers.PaymentMapper
}

func NewPostgrePaymentRepository(db *gorm.DB) *PostgrePaymentRepository {
	mapper := mappers.NewPaymentMapper()

	return &PostgrePaymentRepository{db: db, mapper: mapper}
}

func (r *PostgrePaymentRepository) Create(payment *entities.Payment) error {
	paymentModel := r.mapper.ToPersistence(payment)

	return r.db.Create(paymentModel).Error
}
