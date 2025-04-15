package models

import (
	"time"
)

type NotificationBDV struct {
	BancoOrigen      string    `gorm:"size:4"`    // banco del cliente
	ReferenciaOrigen string    `gorm:"size:15"`   // referencia del banco del cliente
	IdCliente        string    `gorm:"size:15"`   // CI/Rif cliente
	NumeroCliente    string    `gorm:"size:15"`   // tlf cliente
	IdComercio       string    `gorm:"size:23"`   // Rif Comercio
	NumeroComercio   string    `gorm:"size:15"`   // Tlf Comercio
	FechaBanco       string    `gorm:"size:11"`   // Fecha en str
	FechaTranformada time.Time `gorm:"type:date"` // transformar para crear busquedas
	HoraBanco        string    `gorm:"size:7"`    // hora en str
	HoraTransformada time.Time `gorm:"type:time"`
	Monto            float64   `gorm:"type:decimal(13,2)"` // previendo futuras conversiones monetarias
}
