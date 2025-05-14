package api_key

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/desarrolladoresnet/api_notificaciones_bancarias/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

////////////////////////////////////////////
////////////////////////////////////////////
////////////////////////////////////////////

type apiUser struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Password2  string `json:"password2"`
	ApiKeyName string `json:"api_key_name"`
}

func GetApiKey(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request apiUser

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err,
				"message": "no se pudo parsear el json",
				"status":  http.StatusInternalServerError,
			})
			return
		}

		// ------- Check fields -------//
		if strings.TrimSpace(request.Username) == "" || strings.TrimSpace(request.Password) == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "username y password son requeridos",
				"status":  http.StatusBadRequest,
			})
			return
		}

		// ------ Exists User ------- //
		user, err := retrieveUser(db, request.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
				"message": "error al buscar usuario",
				"status":  http.StatusInternalServerError,
			})
			return
		}
		if user == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "usuario no encontrado o inactivo",
				"status":  http.StatusForbidden,
			})
			return
		}

		// ------ Check Pass ------- //
		can_pass := CheckPasswordHash(request.Password, user.Password)
		if !can_pass {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "credenciales inválidas",
				"status":  http.StatusForbidden,
			})
			return
		}

		// ------- Retrieve or Generate new API KEY ------- ///
		var apiKey models.APIKey
		request.ApiKeyName = user.APIKeyName

		// Buscar si ya existe una API key activa con este nombre para el usuario
		err = db.Where("name = ? AND active = ?", request.ApiKeyName, true).First(&apiKey).Error

		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error":   err.Error(),
					"message": "error al buscar API Key existente",
					"status":  http.StatusInternalServerError,
				})
				return
			}

			// No existe, crear una nueva
			newKey, err := generateSecureKey(32)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error":   err.Error(),
					"message": "error generando API Key",
					"status":  http.StatusInternalServerError,
				})
				return
			}

			apiKey = models.APIKey{
				Name:      request.ApiKeyName,
				Key:       newKey,
				Active:    true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err := db.Create(&apiKey).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error":   err.Error(),
					"message": "error guardando API Key",
					"status":  http.StatusInternalServerError,
				})
				return
			}
		}

		// ------- Return API KEY ------- //
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "API Key obtenida/generada exitosamente",
			"data": gin.H{
				"api_key":    apiKey.Key,
				"name":       apiKey.Name,
				"created_at": apiKey.CreatedAt,
				"active":     apiKey.Active,
			},
			"status": http.StatusOK,
		})
	}
}

////////////////////////////////////////////
////////////////////////////////////////////
////////////////////////////////////////////

type createApiKeyRequest struct {
	Name string `json:"name"`
}

func CreateApiKey(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request createApiKeyRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err,
				"message": "no se pudo parsear el json",
			})
			return
		}

		result, err := createAPIKey(db, request.Name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err,
				"message": "crear la entrada",
			})
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"error":   nil,
			"message": "api key creada",
			"data":    result,
		})

	}
}

////////////////////////////////////////////
////////////////////////////////////////////
////////////////////////////////////////////

func GetAPIKeys(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Query params
		name := c.DefaultQuery("name", "")
		activeOnly := c.DefaultQuery("active", "true")
		pageStr := c.DefaultQuery("page", "1")
		limitStr := c.DefaultQuery("limit", "10")

		// Parse pagination
		page, _ := strconv.Atoi(pageStr)
		if page < 1 {
			page = 1
		}

		limit, _ := strconv.Atoi(limitStr)
		if limit < 1 || limit > 100 {
			limit = 10
		}
		offset := (page - 1) * limit

		// Query builder
		query := db.Model(&models.APIKey{})
		if strings.TrimSpace(name) != "" {
			query = query.Where("name ILIKE ?", "%"+name+"%")
		}
		if activeOnly != "false" {
			query = query.Where("active = ?", true)
		}

		var apiKeys []models.APIKey
		var total int64

		query.Count(&total)
		if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&apiKeys).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
				"message": "error al obtener las API Keys",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    apiKeys,
			"meta": gin.H{
				"total": total,
				"page":  page,
				"limit": limit,
			},
		})
	}
}

////////////////////////////////////////////
////////////////////////////////////////////
////////////////////////////////////////////

func DeleteAPIKey(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "ID inválido",
			})
			return
		}

		var apiKey models.APIKey
		if err := db.First(&apiKey, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "API Key no encontrada",
			})
			return
		}

		if !apiKey.Active {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "La API Key ya está inactiva",
			})
			return
		}

		// Soft delete → simplemente marcamos como inactiva
		apiKey.Active = false
		apiKey.UpdatedAt = time.Now()

		if err := db.Save(&apiKey).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "No se pudo desactivar la API Key",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "API Key desactivada correctamente",
			"data":    apiKey,
		})
	}
}

////////////////////////////////////////////
////////////////////////////////////////////
////////////////////////////////////////////

//          FUNCIONES AUXILIARES          //

////////////////////////////////////////////
////////////////////////////////////////////
////////////////////////////////////////////

func createAPIKey(db *gorm.DB, name string) (*models.APIKey, error) {
	var existing models.APIKey
	err := db.Where("name = ? AND active = ?", name, true).First(&existing).Error
	if err == nil {
		return nil, fmt.Errorf("ya existe una API Key activa con ese nombre")
	} else if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("error buscando claves existentes: %w", err)
	}

	key, err := generateSecureKey(32)
	if err != nil {
		return nil, fmt.Errorf("error generando API Key: %w", err)
	}

	apiKey := &models.APIKey{
		Name:      name,
		Key:       key,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(apiKey).Error; err != nil {
		return nil, fmt.Errorf("error guardando API Key en la base de datos: %w", err)
	}

	return apiKey, nil
}

////////////////////////////////////////////
////////////////////////////////////////////
////////////////////////////////////////////

func generateSecureKey(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

////////////////////////////////////////////
////////////////////////////////////////////
////////////////////////////////////////////

func retrieveUser(db *gorm.DB, username string) (*models.APIUser, error) {
	var user_model models.APIUser
	err := db.Where("username = ? AND is_active = ?", username, true).First(&user_model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // User not found, return nil without error
		}
		return nil, fmt.Errorf("error searching for existing keys: %w", err)
	}

	// If we get here, user exists and is active
	return &user_model, nil
}
