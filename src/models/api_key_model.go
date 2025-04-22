package models

import (
	"time"

	"gorm.io/gorm"
)

type APIKey struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`             // Nombre del cliente o propósito
	Key       string `gorm:"uniqueIndex;not null"` // La clave en sí
	Active    bool   `gorm:"default:true"`         // Si está habilitada o no
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` // Para borrado suave (soft delete)
}
