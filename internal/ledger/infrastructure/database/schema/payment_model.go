package schema

import (
	"time"

	"gorm.io/gorm"
)

type PaymentModel struct {
	gorm.Model
	ID          string `gorm:"primaryKey"`
	Amount      int    `gorm:"not null default:0"`
	Date        time.Time
	Name        string `gorm:"not null"`
	Email       string `gorm:"not null"`
	PaymentType string `gorm:"not null"`
}

func (PaymentModel) TableName() string {
	return "payments"
}
