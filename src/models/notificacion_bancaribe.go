package models

import "time"

type NotificationBancaribe struct {
	ID                   uint      `gorm:"primaryKey"`
	Amount               float64   `gorm:"type:decimal(13,2)column:amount"`
	BankName             string    `gorm:"size:30;column:bank_name"`
	ClientPhone          string    `gorm:"size:16;column:client_phone"`
	CommercePhone        string    `gorm:"size:16column:commerce_phone"`
	CreditorAccount      string    `gorm:"size:50column:creditor_account"`
	CurrencyCode         string    `gorm:"size:5column:currency_code"`
	DateBancaribe        string    `gorm:"size:12column:date_bancaribe"`
	Date                 time.Time `gorm:"type:datecolumn:date"`
	DebtorID             string    `gorm:"size:15;column:debtor_id"`
	DestinyBankReference string    `gorm:"size:15;column:destiny_bank_reference"`
	OriginBankCode       string    `gorm:"size:5;column:origin_bank_code"`
	OriginBankReference  string    `gorm:"size:15;colum:origin_bank_reference"`
	PaymentType          string    `gorm:"size:6;column:payment_type"`
	TimeBancaribe        string    `gorm:"size:10column:time_bancaribe"`
	Time                 time.Time `gorm:"type:timecolumn:time"`
}
