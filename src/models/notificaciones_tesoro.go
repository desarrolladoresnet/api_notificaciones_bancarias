package models

import "time"

type NotificacionTesoro struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	PaymentType     string    `gorm:"column:payment_type;size:15" json:"payment_type"`
	Reference       string    `gorm:"column:reference;size:50" json:"reference"`
	SourceBank      string    `gorm:"column:source_bank;size:10" json:"source_bank"`
	Amount          float64   `gorm:"column:amount;type:decimal(13,2)" json:"amount"`
	SourcePhone     string    `gorm:"column:source_phone;size:20" json:"source_phone"`
	PaymenDate      string    `gorm:"column:payment_date;size:32" json:"payment_date"`
	PaymentDateTime time.Time `gorm:"column:payment_date_time" json:"payment_date_time"` // fecha en formato time
	Document        string    `gorm:"column:document;size:32" json:"document"`
}
