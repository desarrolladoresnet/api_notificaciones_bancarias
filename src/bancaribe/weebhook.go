package bancaribe

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/desarrolladoresnet/api_notificaciones_bancarias/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type notificationBancaribe struct {
	Amount               float64 `json:"amount"`
	BankName             string  `json:"bankName"`
	ClientPhone          string  `json:"clientPhone"`
	CommercePhone        string  `json:"commercePhone"`
	CreditorAccount      string  `json:"creditorAccount"`
	CurrencyCode         string  `json:"currencyCode"`
	Date                 string  `json:"date"`
	DebtorID             string  `json:"debtorID"`
	DestinyBankReference string  `json:"destinyBankReference"`
	OriginBankCode       string  `json:"originBankCode"`
	OriginBankReference  string  `json:"originBankReference"`
	PaymentType          string  `json:"paymentType"`
	Time                 string  `json:"time"`
}

func WeebHookBancaribe(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request notificationBancaribe

		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"codigo":         nil,
				"mensajeCliente": "error while receiving JSON data",
				"mensajeSistema": err.Error(),
				"success":        false,
			})
			return
		}

		model, err := tranformRequestToModel(request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"codigo":         nil,
				"mensajeCliente": "erro transforming to model",
				"mensajeSistema": err.Error(),
				"success":        false,
			})
			return
		}

		// ----- CHECK IF EXIST ----- //
		exist, err := CheckNotificationExists(*model, db)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"codigo":         nil,
				"mensajeCliente": "erro while checking the DB",
				"mensajeSistema": err.Error(),
				"success":        false,
			})
			return
		}
		if exist {
			c.JSON(http.StatusOK, gin.H{
				"mensajeCliente": "pago previamente recibido",
				"mensajeSistema": "Notificado",
				"success":        true,
			})
			return
		}

		// ----- SAVE DB ----- //
		result, err := saveNotification(model, db)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"codigo":         nil,
				"mensajeCliente": "erro while saving data in the DB",
				"mensajeSistema": err.Error(),
				"success":        false,
			})
			return
		}
		if !result { // ALREADY EXIST
			c.JSON(http.StatusBadRequest, gin.H{
				"codigo":         nil,
				"mensajeCliente": "can't save the data in the DB",
				"mensajeSistema": nil,
				"success":        false,
			})

		}

		c.JSON(http.StatusCreated, gin.H{
			"codigo":         "00",
			"mensajeCliente": "Aprobado",
			"mensajeSistema": "Notificado",
			"success":        true,
			"data":           model,
		})

	}

}

////////////////////////////////////////////////
////////////////////////////////////////////////
////////////////////////////////////////////////

//           FUNCIONES AUXILIARES             //

////////////////////////////////////////////////
////////////////////////////////////////////////
////////////////////////////////////////////////

func tranformRequestToModel(request notificationBancaribe) (*models.NotificationBancaribe, error) {
	fecha, err := TransformDate(request.Date)
	if err != nil {
		fmt.Println("---- Error en Date -----")
		return nil, err
	}

	hora, err := TransformHour(request.Time)
	if err != nil {
		fmt.Println("---- Error en Time -----")

		return nil, err
	}

	// monto, err := strconv.ParseFloat(request.Monto, 64)
	// if err != nil {
	// 	return nil, err
	// }
	if request.Amount <= 0.0 {
		return nil, fmt.Errorf("the amount of the payment can't be negative: %v", request.Amount)
	}

	notificacion := models.NotificationBancaribe{
		Amount:               request.Amount,
		BankName:             request.BankName,
		ClientPhone:          request.ClientPhone,
		CommercePhone:        request.CommercePhone,
		CreditorAccount:      request.CreditorAccount,
		CurrencyCode:         request.CurrencyCode,
		DateBancaribe:        request.Date,
		Date:                 *fecha,
		DebtorID:             request.DebtorID,
		DestinyBankReference: request.DestinyBankReference,
		OriginBankCode:       request.OriginBankCode,
		OriginBankReference:  request.OriginBankReference,
		PaymentType:          request.PaymentType,
		TimeBancaribe:        request.Time,
		Time:                 *hora,
	}

	return &notificacion, nil
}

/////////////////////////////////////////////////

func TransformDate(date string) (*time.Time, error) {
	// Intentar con formato día-mes-año primero
	parseDate, err := time.Parse("02-01-2006", date)
	if err != nil {
		// Si falla, intentar con otro formato si es necesario
		parseDate, err = time.Parse("2006-01-02", date)
		if err != nil {
			return nil, fmt.Errorf("fecha no válida: %v (formatos soportados: 'dd-mm-yyyy' o 'yyyy-mm-dd')", date)
		}
	}

	return &parseDate, nil
}

/////////////////////////////////////////////////

func TransformHour(timeStr string) (*time.Time, error) {
	// Intentar con diferentes formatos
	layouts := []string{"15.04", "15:04", "1504", "15 04", "15.04.00", "15:04:00", "150400", "15 04 00"}

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

//////////////////////////////////////////////////.
//////////////////////////////////////////////////.
//////////////////////////////////////////////////.

func CheckNotificationExists(model models.NotificationBancaribe, db *gorm.DB) (bool, error) {

	if db == nil {
		return false, fmt.Errorf("database connection cannot be nil")
	}

	var count int64

	result := db.Model(&models.NotificationBancaribe{}).
		Where("origin_bank_code = ? AND destiny_bank_reference = ? AND date = ? AND amount = ? AND origin_bank_reference = ?",
			model.OriginBankCode, model.DestinyBankReference, model.Date, model.Amount, model.OriginBankReference).
		Count(&count)

	if result.Error != nil {
		log.Printf("Error checking for existing BDV notification: %v", result.Error)
		return false, result.Error
	}

	// Si count > 0, significa que ya existe al menos una notificación con esos datos
	return count > 0, nil
}

//////////////////////////////////////////////////.
//////////////////////////////////////////////////.
//////////////////////////////////////////////////.

func saveNotification(model *models.NotificationBancaribe, db *gorm.DB) (bool, error) {
	fmt.Printf("\n ---%v \n\n", model)
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
