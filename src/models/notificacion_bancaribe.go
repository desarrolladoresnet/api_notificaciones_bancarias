package models

import "time"

type NotificationBancaribe struct {
	// El ID principal sigue siendo uint, ya que generalmente es autoincremental
	// y no nulo. Si también puede ser nulo, cámbialo a *uint.
	ID uint `gorm:"primaryKey" json:"id"`

	// El resto de los campos se cambian a punteros para ser opcionales/nulables.
	Amount               *float64   `gorm:"type:decimal(13,2);column:amount" json:"amount,omitempty"`
	BankName             *string    `gorm:"size:50;column:bank_name" json:"bank_name,omitempty"`
	ClientPhone          *string    `gorm:"size:20;column:client_phone" json:"client_phone,omitempty"`
	CommercePhone        *string    `gorm:"size:20;column:commerce_phone" json:"commerce_phone,omitempty"`
	CreditorAccount      *string    `gorm:"size:50;column:creditor_account" json:"creditor_account,omitempty"`
	CurrencyCode         *string    `gorm:"size:5;column:currency_code" json:"currency_code,omitempty"`
	DateBancaribe        *string    `gorm:"size:12;column:date_bancaribe" json:"date_bancaribe,omitempty"`
	Date                 *time.Time `gorm:"type:date;column:date" json:"date,omitempty"`
	DebtorAccount        *string    `gorm:"size:50;column:debtor_account" json:"debtor_account,omitempty"`
	DebtorID             *string    `gorm:"size:20;column:debtor_id" json:"debtor_id,omitempty"`
	DestinyBankReference *string    `gorm:"size:20;column:destiny_bank_reference" json:"destiny_bank_reference,omitempty"`
	OriginBankCode       *string    `gorm:"size:10;column:origin_bank_code" json:"origin_bank_code,omitempty"`
	OriginBankReference  *string    `gorm:"size:20;column:origin_bank_reference" json:"origin_bank_reference,omitempty"`
	PaymentType          *string    `gorm:"size:10;column:payment_type" json:"payment_type,omitempty"`
	TimeBancaribe        *string    `gorm:"size:12;column:time_bancaribe" json:"time_bancaribe,omitempty"`
	Time                 *time.Time `gorm:"type:time;column:time" json:"time,omitempty"`
}
