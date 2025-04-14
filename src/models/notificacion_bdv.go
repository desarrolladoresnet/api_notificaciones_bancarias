package models

import "time"

type NotificationBDV struct {
	BancoOrigen      string    `gorm:"size:4"`
	ReferenciaOrigen string    `gorm:"size:15"`
	NumeroCliente    string    `gorm:"size:15"`
	IdComercio       string    `gorm:"size:23"`
	NumeroComercio   string    `gorm:"size:15"`
	FechaBanco       string    `gorm:"size:11"`
	FechaTranformada time.Time `gorm:"type:date"` // transformar para crear busquedas
	HoraBanco        string    `gorm:"size:7"`
	HoraTransformada time.Time `gorm:"type:time"`
	Monto            float64   `gorm:"type:decimal(13,2)"` // previendo futuras conversiones monetarias
}
