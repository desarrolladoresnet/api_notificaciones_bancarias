package bdv

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/desarrolladoresnet/api_notificaciones_bancarias/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type bdvRequest struct {
	BancoOrdenante string `json:"bancoOrdenante"`
	Referencia     string `json:"referenciaBancoOrdenante"`
	IdCliente      string `json:"idCliente"`
	IdComercio     string `json:"idComercio"`
	NumeroCliente  string `json:"numeroCliente"`
	NumeroComercio string `json:"numeroComercio"`
	Fecha          string `json:"fecha"`
	Hora           string `json:"hora"`
	Monto          string `json:"monto"`
}

// --- MEJORA: Función auxiliar para enviar respuestas de error consistentes ---
// Esto evita repetir el mismo c.JSON en cada punto de fallo.
func sendErrorResponse(c *gin.Context, httpStatus int, clientMsg, systemMsg string) {
	c.JSON(httpStatus, gin.H{
		"codigo":         nil,
		"mensajeCliente": clientMsg,
		"mensajeSistema": systemMsg,
		"success":        false,
	})
}

func WeebHookBDV(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request bdvRequest

		// 1. Parsear y verificar el JSON de entrada.
		if err := c.ShouldBindJSON(&request); err != nil {
			sendErrorResponse(c, http.StatusBadRequest, "error al recibir los datos JSON", err.Error())
			return
		}

		// 2. Validar que los campos obligatorios no estén vacíos.
		if err := request.Validate(); err != nil {
			sendErrorResponse(c, http.StatusBadRequest, "error de validación", err.Error())
			return
		}

		// 3. Transformar la solicitud a nuestro modelo de base de datos.
		model, err := tranformRequestToModel(request)
		if err != nil {
			log.Printf("Error al transformar el cuerpo de la solicitud: %v", err)
			sendErrorResponse(c, http.StatusBadRequest, "error al transformar los datos JSON", err.Error())
			return
		}

		// 4. Buscar si la notificación ya existe para evitar duplicados.
		exists, err := CheckNotificationExists(request.BancoOrdenante, request.Referencia, request.Fecha, request.IdCliente, db)
		if err != nil {
			log.Printf("Error al consultar la base de datos: %v", err)
			sendErrorResponse(c, http.StatusInternalServerError, "error al acceder a la base de datos", err.Error())
			return
		}
		if exists {
			c.JSON(http.StatusOK, gin.H{
				"codigo":         "01",
				"mensajeCliente": "pago previamente recibido",
				"mensajeSistema": "Notificado",
				"success":        true,
			})
			return
		}

		// 5. Si no existe, guardar la nueva notificación en la base de datos.
		saved, err := saveNotification(model, db)
		if err != nil || !saved {
			// El error ya se loguea dentro de saveNotification
			sendErrorResponse(c, http.StatusInternalServerError, "error al guardar los datos", err.Error())
			return
		}

		// 6. Si todo fue exitoso, enviar respuesta de éxito.
		c.JSON(http.StatusCreated, gin.H{
			"codigo":         "00",
			"mensajeCliente": "Aprobado",
			"mensajeSistema": "Notificado",
			"success":        true,
		})
	}
}

/////////////////////////////////////////////////

func tranformRequestToModel(request bdvRequest) (*models.NotificationBDV, error) {
	fecha, err := TransformDate(request.Fecha)
	if err != nil {
		return nil, err
	}

	hora, err := TransformHour(request.Hora)
	if err != nil {
		return nil, err
	}

	monto, err := strconv.ParseFloat(request.Monto, 64)
	if err != nil {
		// ANOTACIÓN: Error más específico para el parseo del monto.
		return nil, fmt.Errorf("el campo 'monto' no es un número válido: %w", err)
	}
	if monto <= 0 {
		return nil, fmt.Errorf("el monto del pago no puede ser cero o negativo: %v", monto)
	}

	notificacion := models.NotificationBDV{
		BancoOrigen:      request.BancoOrdenante,
		ReferenciaOrigen: request.Referencia,
		IdCliente:        request.IdCliente,
		NumeroCliente:    request.NumeroCliente,
		NumeroComercio:   request.NumeroComercio,
		IdComercio:       request.IdComercio,
		FechaBanco:       request.Fecha,
		FechaTranformada: *fecha,
		HoraBanco:        request.Hora,
		HoraTransformada: *hora,
		Monto:            monto,
	}

	return &notificacion, nil
}

/////////////////////////////////////////////////

// --- CORRECCIÓN: Se actualizó para aceptar múltiples formatos de fecha ---
func TransformDate(dateStr string) (*time.Time, error) {
	// Lista de layouts de fecha aceptados.
	layouts := []string{
		"2006-01-02", // Formato YYYY-MM-DD
		"20060102",   // Formato YYYYMMDD
	}

	for _, layout := range layouts {
		parsedDate, err := time.Parse(layout, dateStr)
		if err == nil {
			// Si el parseo es exitoso, retornamos el resultado.
			return &parsedDate, nil
		}
	}

	// Si ningún formato funcionó, retornamos un error claro.
	return nil, fmt.Errorf("formato de fecha inválido. Formatos aceptados: YYYY-MM-DD, YYYYMMDD")
}

/////////////////////////////////////////////////

// ANOTACIÓN: Esta función ya era correcta y aceptaba el formato "2311" (HHMM).
// No se necesitan cambios aquí.
func TransformHour(timeStr string) (*time.Time, error) {
	layouts := []string{"15.04", "15:04", "1504", "15 04"}

	for _, layout := range layouts {
		t, err := time.Parse(layout, timeStr)
		if err == nil {
			// La verificación de rango es una buena práctica, aunque time.Parse ya la realiza.
			if t.Hour() < 0 || t.Hour() > 23 || t.Minute() < 0 || t.Minute() > 59 {
				return nil, fmt.Errorf("hora o minutos fuera de rango")
			}
			return &t, nil
		}
	}

	return nil, fmt.Errorf("formato de hora inválido, formatos aceptados: HH.MM, HH:MM, HHMM, HH MM")
}

/////////////////////////////////////////////////

func saveNotification(model *models.NotificationBDV, db *gorm.DB) (bool, error) {
	if model == nil {
		return false, fmt.Errorf("el modelo de notificación no puede ser nulo")
	}
	if db == nil {
		return false, fmt.Errorf("la conexión a la base de datos no puede ser nula")
	}

	result := db.Create(model)
	if result.Error != nil {
		log.Printf("Error al guardar notificación BDV: %v", result.Error)
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		log.Println("No se afectaron filas al guardar la notificación")
		return false, fmt.Errorf("no se afectaron filas en la base de datos")
	}

	return true, nil
}

/////////////////////////////////////////////////

func CheckNotificationExists(bancoOrigen string, referenciaOrigen string, fechaBanco string,
	idCliente string, db *gorm.DB) (bool, error) {

	if db == nil {
		return false, fmt.Errorf("la conexión a la base de datos no puede ser nula")
	}

	var count int64
	// ANOTACIÓN: Es buena práctica usar Model(&models.NotificationBDV{}) para ser explícito.
	result := db.Model(&models.NotificationBDV{}).
		Where("banco_origen = ? AND referencia_origen = ? AND fecha_banco = ? AND id_cliente = ?",
			bancoOrigen, referenciaOrigen, fechaBanco, idCliente).
		Count(&count)

	if result.Error != nil {
		log.Printf("Error al verificar si la notificación BDV existe: %v", result.Error)
		return false, result.Error
	}

	return count > 0, nil
}

/////////////////////////////////////////////////

// Validate verifica que ningún campo de la solicitud esté vacío.
func (r *bdvRequest) Validate() error {
	// ANOTACIÓN: Este método es simple y efectivo para este caso.
	if r.BancoOrdenante == "" {
		return fmt.Errorf("el campo 'bancoOrdenante' es obligatorio")
	}
	if r.Referencia == "" {
		return fmt.Errorf("el campo 'referenciaBancoOrdenante' es obligatorio")
	}
	if r.IdCliente == "" {
		return fmt.Errorf("el campo 'idCliente' es obligatorio")
	}
	if r.IdComercio == "" {
		return fmt.Errorf("el campo 'idComercio' es obligatorio")
	}
	if r.NumeroCliente == "" {
		return fmt.Errorf("el campo 'numeroCliente' es obligatorio")
	}
	if r.NumeroComercio == "" {
		return fmt.Errorf("el campo 'numeroComercio' es obligatorio")
	}
	if r.Fecha == "" {
		return fmt.Errorf("el campo 'fecha' es obligatorio")
	}
	if r.Hora == "" {
		return fmt.Errorf("el campo 'hora' es obligatorio")
	}
	if r.Monto == "" {
		return fmt.Errorf("el campo 'monto' es obligatorio")
	}
	return nil
}
