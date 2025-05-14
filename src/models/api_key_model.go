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

type APIUser struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Username   string `gorm:"column:username;size:18;unique;not null" json:"username"` // Unique y not null para evitar duplicados
	Password   string `gorm:"column:password;size:512" json:"-"`                       // `json:"-"` omite el campo en respuestas JSON
	IsActive   bool   `gorm:"column:is_active;default:true" json:"_"`                  // Mejor nombre de columna y campo JSON
	APIKeyName string `gorm:"column:api_keyname;size:30;unique;not null" json:"api_keyname"`
}
