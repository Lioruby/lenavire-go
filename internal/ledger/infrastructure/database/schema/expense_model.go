package schema

import (
	"time"

	"gorm.io/gorm"
)

type ExpenseModel struct {
	gorm.Model
	ID     string `gorm:"primaryKey"`
	Amount int    `gorm:"not null default:0"`
	Date   time.Time
}

func (ExpenseModel) TableName() string {
	return "expenses"
}
