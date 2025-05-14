package api_key

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/desarrolladoresnet/api_notificaciones_bancarias/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateApiUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request apiUser

		// Parsear el JSON de entrada
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
				"message": "No se pudo parsear el JSON",
				"status":  http.StatusBadRequest,
			})
			return
		}

		// Validar campos obligatorios
		if strings.TrimSpace(request.Username) == "" || strings.TrimSpace(request.Password) == "" || strings.TrimSpace(request.ApiKeyName) == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Username y password son requeridos",
				"status":  http.StatusBadRequest,
			})
			return
		}

		// Validar que las contraseñas coincidan
		if request.Password != request.Password2 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Las contraseñas no coinciden",
				"status":  http.StatusBadRequest,
			})
			return
		}

		// Validar longitud mínima de contraseña
		if len(request.Password) < 8 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "La contraseña debe tener al menos 8 caracteres",
				"status":  http.StatusBadRequest,
			})
			return
		}

		// Verificar si el usuario ya existe
		var existingUser models.APIUser
		if err := db.Where("username = ?", request.Username).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "El nombre de usuario ya existe",
				"status":  http.StatusConflict,
			})
			return
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
				"message": "Error al verificar usuario existente",
				"status":  http.StatusInternalServerError,
			})
			return
		}

		// Hashear la contraseña
		hashedPassword, err := HashPassword(request.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
				"message": "Error al hashear la contraseña",
				"status":  http.StatusInternalServerError,
			})
			return
		}

		// Crear el nuevo usuario
		newUser := models.APIUser{
			Username:   request.Username,
			Password:   string(hashedPassword),
			IsActive:   true,
			APIKeyName: request.ApiKeyName,
		}

		if err := db.Create(&newUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
				"message": "Error al crear el usuario",
				"status":  http.StatusInternalServerError,
			})
			return
		}

		// Respuesta exitosa (sin incluir la contraseña)
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Usuario creado exitosamente",
			"data": gin.H{
				"id":        newUser.ID,
				"username":  newUser.Username,
				"is_active": newUser.IsActive,
			},
			"status": http.StatusCreated,
		})
	}
}

func CreateDefaultUser(db *gorm.DB, username, password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return fmt.Errorf("error creating user password: %v", err.Error())
	}

	newUser := models.APIUser{
		Username:   username,
		Password:   hashedPassword,
		IsActive:   true,
		APIKeyName: "desarrollo_netcom",
	}

	if err := db.Create(&newUser).Error; err != nil {
		return fmt.Errorf("error creating default User: %v", err.Error())
	}
	return nil
}
