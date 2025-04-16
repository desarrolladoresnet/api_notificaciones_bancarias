package bdv

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/desarrolladoresnet/api_notificaciones_bancarias/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type getNotificationsBdv struct {
	Referencia    string `form:"referencia"`
	Fecha         string `form:"fecha"`
	NumeroCliente string `form:"numero_cliente"`
	IdCliente     string `form:"id_cliente"`
	Pagina        string `form:"pagina"`
}

func GetPayments(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request getNotificationsBdv
		if err := c.ShouldBind(&request); err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid request format", err.Error())
			return
		}

		// Validar parámetros (solo los que no están vacíos)
		if valid, errors := validateRequestParams(request); !valid {
			respondWithValidationError(c, errors)
			return
		}

		// Configurar paginación por defecto
		page, pageSize := configurePagination(request.Pagina)

		// Buscar notificaciones
		notifications, totalCount, err := searchWithFilters(db, request, page, pageSize)
		if err != nil {
			log.Printf("Error searching notifications: %v", err)
			respondWithError(c, http.StatusInternalServerError, "Database error", err.Error())
			return
		}

		// Calcular total de páginas
		totalPages := calculateTotalPages(totalCount, pageSize)

		// Responder con éxito
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Notifications retrieved successfully",
			"data":    notifications,
			"pagination": gin.H{
				"current_page": page,
				"total_pages":  totalPages,
				"total_items":  totalCount,
				"page_size":    pageSize,
			},
		})
	}
}

//////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////

// validateRequestParams valida solo los campos no vacíos
func validateRequestParams(request getNotificationsBdv) (bool, []string) {
	var errors []string

	if request.Referencia != "" {
		if _, err := strconv.Atoi(request.Referencia); err != nil {
			errors = append(errors, "La referencia debe contener solo números")
		}
	}

	if request.Pagina != "" {
		if _, err := strconv.Atoi(request.Pagina); err != nil {
			errors = append(errors, "La página debe ser un número válido")
		}
	}

	if request.Fecha != "" {
		if _, err := time.Parse("2006-01-02", request.Fecha); err != nil {
			errors = append(errors, "La fecha debe estar en formato YYYY-MM-DD")
		}
	}

	if request.NumeroCliente != "" {
		if !isValidVenezuelanPhone(request.NumeroCliente) {
			errors = append(errors, "Número de teléfono venezolano inválido")
		}
	}

	if request.IdCliente != "" {
		if !isValidVenezuelanID(request.IdCliente) {
			errors = append(errors, "ID de cliente inválido (debe comenzar con V, G, J, E o P)")
		}
	}

	return len(errors) == 0, errors
}

//////////////////////////////////////////////////////////

// searchWithFilters construye la consulta con los filtros proporcionados
func searchWithFilters(db *gorm.DB, request getNotificationsBdv, page, pageSize int) ([]models.NotificationBDV, int64, error) {
	query := db.Model(&models.NotificationBDV{})

	if request.Referencia != "" {
		query = query.Where("referencia_origen = ?", request.Referencia)
	}
	if request.Fecha != "" {
		query = query.Where("fecha_banco = ?", request.Fecha)
	}
	if request.NumeroCliente != "" {
		query = query.Where("numero_cliente = ?", request.NumeroCliente)
	}
	if request.IdCliente != "" {
		query = query.Where("id_cliente = ?", request.IdCliente)
	}

	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	var notifications []models.NotificationBDV
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&notifications).Error

	return notifications, totalCount, err
}

//////////////////////////////////////////////////////////

// Funciones auxiliares
func isValidVenezuelanPhone(phone string) bool {
	regex := regexp.MustCompile(`^(0)?(412|414|416|418|424|426)[0-9]{7}$`)
	return regex.MatchString(phone)
}

//////////////////////////////////////////////////////////

func isValidVenezuelanID(id string) bool {
	regex := regexp.MustCompile(`^[VGJEPvgjep][0-9]{5,9}$`)
	return regex.MatchString(id)
}

//////////////////////////////////////////////////////////

func configurePagination(pageStr string) (int, int) {
	const defaultPage = 1
	const defaultPageSize = 100

	page := defaultPage
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	return page, defaultPageSize
}

//////////////////////////////////////////////////////////

func calculateTotalPages(totalCount int64, pageSize int) int {
	if totalCount == 0 {
		return 1
	}
	return int((totalCount + int64(pageSize) - 1) / int64(pageSize))
}

//////////////////////////////////////////////////////////

func respondWithError(c *gin.Context, status int, message, detail string) {
	c.JSON(status, gin.H{
		"success": false,
		"message": message,
		"error":   detail,
	})
}

//////////////////////////////////////////////////////////

func respondWithValidationError(c *gin.Context, errors []string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"message": "Validation failed",
		"errors":  errors,
	})
}
