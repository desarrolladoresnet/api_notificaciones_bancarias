package models

import "time"

type NotificationBancaribe struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	Amount               float64   `gorm:"type:decimal(13,2);column:amount" json:"amount"`
	BankName             string    `gorm:"size:30;column:bank_name" json:"bank_name"`
	ClientPhone          string    `gorm:"size:16;column:client_phone" json:"client_phone"`
	CommercePhone        string    `gorm:"size:16;column:commerce_phone" json:"commerce_phone"`
	CreditorAccount      string    `gorm:"size:50;column:creditor_account" json:"creditor_account"`
	CurrencyCode         string    `gorm:"size:5;column:currency_code" json:"currency_code"`
	DateBancaribe        string    `gorm:"size:12;column:date_bancaribe" json:"date_bancaribe"`
	Date                 time.Time `gorm:"type:date;column:date" json:"date"`
	DebtorID             string    `gorm:"size:15;column:debtor_id" json:"debtor_id"`
	DestinyBankReference string    `gorm:"size:15;column:destiny_bank_reference" json:"destiny_bank_reference"`
	OriginBankCode       string    `gorm:"size:5;column:origin_bank_code" json:"origin_bank_code"`
	OriginBankReference  string    `gorm:"size:15;column:origin_bank_reference" json:"origin_bank_reference"`
	PaymentType          string    `gorm:"size:6;column:payment_type" json:"payment_type"`
	TimeBancaribe        string    `gorm:"size:10;column:time_bancaribe" json:"time_bancaribe"`
	Time                 time.Time `gorm:"type:time;column:time" json:"time"`
}
