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

func WeebHookBDV(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request bdvRequest

		// ------- Parsear y verificar el JSON ------ //
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"codigo":         nil,
				"mensajeCliente": "error while receiving JSON data",
				"mensajeSistema": err.Error(),
				"success":        false,
			})
			return
		}

		if err := request.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"codigo":         nil,
				"mensajeCliente": "validation error",
				"mensajeSistema": err.Error(),
				"success":        false,
			})
			return
		}

		// ------- Parsear y verificar el JSON ------ //
		// fmt.Println("--- Transformando ---")
		model, err := tranformRequestToModel(request)
		if err != nil {
			log.Printf("Error while parsing the request body\n%v", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"codigo":         nil,
				"mensajeCliente": "error while transforming JSON data",
				"mensajeSistema": err.Error(),
				"success":        false,
			})
			return
		}

		// ---- Buscar si ya se habia reportado la notificacion ----- //
		// CheckNotificationExists arroja true si ya existe una entrada con los mismos datos
		result, err := CheckNotificationExists(request.BancoOrdenante, request.Referencia, request.Fecha, request.IdCliente, db)
		if err != nil {
			log.Printf("Error while parsing the request body\n%v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"codigo":         nil,
				"mensajeCliente": "error while accesong to the database",
				"mensajeSistema": err.Error(),
				"success":        false,
			})
			return
		}
		if result {
			c.JSON(http.StatusOK, gin.H{
				"codigo":         "01",
				"mensajeCliente": "pago previamente recibido",
				"mensajeSistema": "Notificado",
				"success":        true,
			})
			return
		}

		// ---- Si no existe la notificacion se guarda en BD ----- //
		result, err = saveNotification(model, db)
		if err != nil || !result {
			c.JSON(http.StatusInternalServerError, gin.H{
				"codigo":         nil,
				"mensajeCliente": "error saving the data",
				"mensajeSistema": err.Error(),
				"success":        false,
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"codigo":         "00",
			"mensajeCliente": "Aprobado",
			"mensajeSistema": "Notificado",
			"success":        result,
		})

	}
}

/////////////////////////////////////////////////
/////////////////////////////////////////////////
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
		return nil, err
	}
	if monto <= 0 {
		return nil, fmt.Errorf("the amount of the payment can't be negative: %v", monto)
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

func TransformDate(date string) (*time.Time, error) {
	parseDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	return &parseDate, nil
}

/////////////////////////////////////////////////

func TransformHour(timeStr string) (*time.Time, error) {
	// Intentar con diferentes formatos
	layouts := []string{"15.04", "15:04", "1504", "15 04"}

	for _, layout := range layouts {
		t, err := time.Parse(layout, timeStr)
		if err == nil {
			// Verificar rangos si el parseo fue exitoso
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
		return false, fmt.Errorf("notification model cannot be nil")
	}

	if db == nil {
		return false, fmt.Errorf("database connection cannot be nil")
	}

	result := db.Create(model)
	if result.Error != nil {
		log.Printf("Error saving BDV notification: %v", result.Error)
		return false, result.Error
	}

	if result.RowsAffected == 0 {
		log.Println("No rows were affected when saving notification")
		return false, fmt.Errorf("no rows affected")
	}

	return true, nil
}

/////////////////////////////////////////////////

func CheckNotificationExists(bancoOrigen string, referenciaOrigen string, fechaBanco string,
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

// Validate verifica que ningún campo esté vacío.
func (r *bdvRequest) Validate() error {
	if r.BancoOrdenante == "" {
		return fmt.Errorf("bancoOrdenante es obligatorio")
	}
	if r.Referencia == "" {
		return fmt.Errorf("referenciaBancoOrdenante es obligatorio")
	}
	if r.IdCliente == "" {
		return fmt.Errorf("idCliente es obligatorio")
	}
	if r.IdComercio == "" {
		return fmt.Errorf("idComercio es obligatorio")
	}
	if r.NumeroCliente == "" {
		return fmt.Errorf("numeroCliente es obligatorio")
	}
	if r.NumeroComercio == "" {
		return fmt.Errorf("numeroComercio es obligatorio")
	}
	if r.Fecha == "" {
		return fmt.Errorf("fecha es obligatorio")
	}
	if r.Hora == "" {
		return fmt.Errorf("hora es obligatorio")
	}
	if r.Monto == "" {
		return fmt.Errorf("monto es obligatorio")
	}
	return nil // Todos los campos están correctos
}
