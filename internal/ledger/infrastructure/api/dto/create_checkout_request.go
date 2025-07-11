package dto

import "lenavire/internal/ledger/domain/entities"

type CreateCheckoutRequest struct {
	Amount           int                           `json:"amount" validate:"required,gt=0"`
	SubscriptionType entities.SubscriptionType    `json:"subscription_type" validate:"required,oneof=one_time subscription"`
}