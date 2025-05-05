package models

import (
	"time"
)

type NotificationBDV struct {
	BancoOrigen      string    `gorm:"column:banco_origen;size:4" json:"banco_origen"`                // banco del cliente
	ReferenciaOrigen string    `gorm:"column:referencia_origen;size:15" json:"referencia_origen"`     // referencia del banco del cliente
	IdCliente        string    `gorm:"column:id_cliente;size:15" json:"id_cliente"`                   // CI/Rif cliente
	NumeroCliente    string    `gorm:"column:numero_cliente;size:15" json:"numero_cliente"`           // tlf cliente
	IdComercio       string    `gorm:"column:id_comercio;size:23" json:"id_comercio"`                 // Rif Comercio
	NumeroComercio   string    `gorm:"column:numero_comercio;size:15" json:"numero_comercio"`         // Tlf Comercio
	FechaBanco       string    `gorm:"column:fecha_banco;size:11" json:"fecha_banco"`                 // Fecha en str
	FechaTranformada time.Time `gorm:"column:fecha_transformada;type:date" json:"fecha_transformada"` // transformar para crear busquedas
	HoraBanco        string    `gorm:"column:hora_banco;size:7" json:"hora_banco"`                    // hora en str
	HoraTransformada time.Time `gorm:"column:hora_transformada;type:time" json:"hora_transformada"`
	Monto            float64   `gorm:"column:monto;type:decimal(13,2)" json:"monto"` // previendo futuras conversiones monetarias
}
