package bdv

import (
	"fmt"
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
	IdComercio     string `json:"numeroComercio"`
	NumeroCliente  string `json:"numeroCliente"`
	NumeroComercio string `json:"tranformRequestToModelroComercio"`
	Fecha          string `json:"fecha"`
	Hora           string `json:"hora"`
	Monto          string `json:"monto"`
}

func WeebHookBDV(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request bdvRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "error while parsing JSON data",
				"success": false,
			})
		}
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
		NumeroCliente:    request.NumeroCliente,
		IdComercio:       request.NumeroComercio,
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
