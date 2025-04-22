package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type APIKey struct {
	ID     uint
	Key    string
	Name   string
	Active bool
}

func APIKeyAuthMiddlewareBDV(db *gorm.DB) gin.HandlerFunc {
	// Lista de nombres permitidos
	allowedNames := map[string]bool{
		"desarrollo_netcom":  true,
		"bdv_notificaciones": true,
	}

	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if strings.TrimSpace(apiKey) == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key requerida"})
			c.Abort()
			return
		}

		var key APIKey
		if err := db.Where("key = ? AND active = ?", apiKey, true).First(&key).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key inválida o inactiva"})
			c.Abort()
			return
		}

		// Verificar si el nombre está en la lista permitida
		if !allowedNames[key.Name] {
			c.JSON(http.StatusForbidden, gin.H{"error": "API key no autorizada para este recurso"})
			c.Abort()
			return
		}

		c.Set("api_key_id", key.ID)
		c.Next()
	}
}
