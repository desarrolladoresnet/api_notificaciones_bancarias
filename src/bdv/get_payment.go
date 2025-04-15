package bdv

import (
	"fmt"
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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Verificar si hay parámetros para validar
		hasParams := request.Referencia != "" || request.Fecha != "" ||
			request.NumeroCliente != "" || request.IdCliente != ""

		// Solo validar los parámetros si al menos uno fue proporcionado
		if hasParams {
			if valid, errors := checkParamsFields(request); !valid {
				c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
				return
			}
		}

		// Establecer página por defecto si no se proporciona
		if request.Pagina == "" {
			request.Pagina = "1"
		}

		// Buscar en BD con los parámetros proporcionados
		notifications, totalCount, err := SearchNotifications(request, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar notificaciones"})
			log.Printf("Error searching notifications: %v", err)
			return
		}

		// Determinar la página actual
		page := 1
		if pageNum, err := strconv.Atoi(request.Pagina); err == nil && pageNum > 0 {
			page = pageNum
		}

		// Calcular el total de páginas
		totalPages := (totalCount + 99) / 100 // Redondeo hacia arriba

		// Responder con los resultados y la información de paginación
		c.JSON(http.StatusOK, gin.H{
			"data": notifications,
			"pagination": gin.H{
				"current_page": page,
				"total_pages":  totalPages,
				"total_items":  totalCount,
				"page_size":    100,
			},
		})
	}
}

/////////////////////////////////////////////////
/////////////////////////////////////////////////
/////////////////////////////////////////////////

func checkParamsFields(request getNotificationsBdv) (bool, []string) {
	var errorsList []string
	var isValid = true

	// Validar Referencia (solo números)
	if _, err := strconv.Atoi(request.Referencia); err != nil {
		isValid = false
		errorsList = append(errorsList, "La referencia solo puede contener números")
	}

	if _, err := strconv.Atoi(request.Pagina); err != nil {
		isValid = false
		errorsList = append(errorsList, "La pagina solo puede contener números")
	}

	// Validar Fecha (formato YYYY-MM-DD)
	if _, err := time.Parse("2006-01-02", request.Fecha); err != nil {
		isValid = false
		errorsList = append(errorsList, "La fecha debe estar en formato YYYY-MM-DD (ej: 2023-01-01)")
	}

	// Validar Número de Cliente (teléfono venezolano)
	phoneRegex := regexp.MustCompile(`^(0)?(412|414|416|418|424|426|416|424)[0-9]{7}$`)
	if !phoneRegex.MatchString(request.NumeroCliente) {
		isValid = false
		errorsList = append(errorsList, "El número de cliente debe ser un teléfono válido para Venezuela (ej: 04141234567 o 4123456789)")
	}

	// Validar ID Cliente (documento venezolano)
	idRegex := regexp.MustCompile(`^[VGJEPvgjep][0-9]{5,9}$`)
	if !idRegex.MatchString(request.IdCliente) {
		isValid = false
		errorsList = append(errorsList, "El ID de cliente debe ser un documento válido de Venezuela (debe comenzar con V, G, J, E o P)")
	}

	return isValid, errorsList
}

/////////////////////////////////////////////////

func SearchNotifications(request getNotificationsBdv, db *gorm.DB) ([]models.NotificationBDV, int64, error) {
	if db == nil {
		return nil, 0, fmt.Errorf("database connection cannot be nil")
	}

	// Inicializar la consulta base
	query := db.Model(&models.NotificationBDV{})

	// Aplicar filtros solo si los parámetros no están vacíos
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

	// Contar el total de registros que coinciden con los filtros
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		log.Printf("Error counting matching notifications: %v", err)
		return nil, 0, err
	}

	// Configurar paginación
	const pageSize = 100
	page := 1

	// Convertir página a entero si no está vacía
	if request.Pagina != "" {
		if pageNum, err := strconv.Atoi(request.Pagina); err == nil && pageNum > 0 {
			page = pageNum
		}
	}

	// Calcular offset para la paginación
	offset := (page - 1) * pageSize

	// Obtener resultados paginados
	var notifications []models.NotificationBDV
	if err := query.Offset(offset).Limit(pageSize).Find(&notifications).Error; err != nil {
		log.Printf("Error fetching notifications: %v", err)
		return nil, 0, err
	}

	return notifications, totalCount, nil
}

/////////////////////////////////////////////////

func GetNotificationExists(bancoOrigen string, referenciaOrigen string, fechaBanco string,
	id_cliente string, db *gorm.DB) (bool, error) {

	if db == nil {
		return false, fmt.Errorf("database connection cannot be nil")
	}

	var count int64

	result := db.Model(&models.NotificationBDV{}).
		Where("banco_origen = ? AND referencia_origen = ? AND fecha_banco = ? AND id_cliente = ?",
			bancoOrigen, referenciaOrigen, fechaBanco, id_cliente).
		Count(&count)

	if result.Error != nil {
		log.Printf("Error checking for existing BDV notification: %v", result.Error)
		return false, result.Error
	}

	// Si count > 0, significa que ya existe al menos una notificación con esos datos
	return count > 0, nil
}

/////////////////////////////////////////////////
