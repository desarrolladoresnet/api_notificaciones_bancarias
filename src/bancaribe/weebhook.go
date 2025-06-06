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
	ClientPhone          *string `json:"clientPhone"`   // Puntero para manejar null
	CommercePhone        *string `json:"commercePhone"` // Puntero para manejar null
	CreditorAccount      string  `json:"creditorAccount"`
	CurrencyCode         string  `json:"currencyCode"`
	Date                 string  `json:"date"`
	DebtorAccount        *string `json:"debtorAccount"` // Nuevo campo
	DebtorID             string  `json:"debtorID"`
	DestinyBankReference string  `json:"destinyBankReference"`
	OriginBankCode       *string `json:"originBankCode"` // Puntero para manejar null
	OriginBankReference  string  `json:"originBankReference"`
	PaymentType          string  `json:"paymentType"`
	Time                 string  `json:"time"`
}

func WeebHookBancaribe(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("PASO 1: Inicio de WeebHookBancaribe")
		var request notificationBancaribe

		fmt.Println("PASO 2: Bind JSON request")
		if err := c.ShouldBind(&request); err != nil {
			fmt.Println("! ERROR en bind JSON:", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"codigo":         nil,
				"message":        "error while receiving JSON data",
				"messageSystem,": err.Error(),
				"statusCode":     http.StatusBadRequest,
				"success":        false,
			})
			return
		}

		fmt.Println("PASO 3: Validar campos obligatorios")
		if err := request.Validate(); err != nil {
			fmt.Println("! ERROR en validación:", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"codigo":         nil,
				"message":        "validation error",
				"messageSystem,": err.Error(),
				"statusCode":     http.StatusBadRequest,
				"success":        false,
			})
			return
		}

		fmt.Println("PASO 4: Transformar request a modelo")
		model, err := tranformRequestToModel(request)
		if err != nil {
			fmt.Println("! ERROR transformando a modelo:", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"codigo":         nil,
				"message":        "erro transforming to model",
				"messageSystem,": err.Error(),
				"statusCode":     http.StatusBadRequest,
				"success":        false,
			})
			return
		}

		fmt.Println("PASO 5: Verificar duplicados en BD")
		exist, err := CheckNotificationExists(*model, db)
		if err != nil {
			fmt.Println("! ERROR verificando existencia:", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"codigo":         nil,
				"message":        "erro while checking the DB",
				"messageSystem,": err.Error(),
				"statusCode":     http.StatusBadRequest,
				"success":        false,
			})
			return
		}
		if exist {
			fmt.Println("-> Pago duplicado. Registro ya existe en BD")
			c.JSON(http.StatusOK, gin.H{
				"codigo":         "01",
				"message":        "pago previamente recibido",
				"messageSystem,": "Notificado",
				"success":        true,
				"statusCode":     http.StatusOK,
			})
			return
		}

		fmt.Println("PASO 6: Guardar notificación en BD")
		result, err := saveNotification(model, db)
		if err != nil {
			fmt.Println("! ERROR guardando en BD:", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"codigo":         nil,
				"message":        "erro while saving data in the DB",
				"messageSystem,": err.Error(),
				"success":        false,
				"statusCode":     http.StatusBadRequest,
			})
			return
		}
		if !result { // ALREADY EXIST
			fmt.Println("! ERROR: No se pudo guardar en BD (result=false)")
			c.JSON(http.StatusBadRequest, gin.H{
				"codigo":         nil,
				"message":        "can't save the data in the DB",
				"messageSystem,": nil,
				"success":        false,
				"statusCode":     http.StatusBadRequest,
			})
			return
		}

		fmt.Println("PASO 7: Notificación procesada exitosamente")
		fmt.Printf(">> ID generado: %v\n", model.ID)
		fmt.Printf(">> Monto: %v\n", model.Amount)
		fmt.Printf(">> Referencia: %v\n", model.OriginBankReference)

		c.JSON(http.StatusCreated, gin.H{
			"codigo":         "00",
			"message":        "Success",
			"messageSystem,": "Notificado",
			"success":        true,
			"data":           model,
			"statusCode":     http.StatusOK,
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
	fmt.Println("-> Transformando fecha...")
	fecha, err := TransformDate(request.Date)
	if err != nil {
		fmt.Println("---- Error en Date -----")
		return nil, err
	}

	fmt.Println("-> Transformando hora...")
	hora, err := TransformHour(request.Time)
	if err != nil {
		fmt.Println("---- Error en Time -----")
		return nil, err
	}

	if request.Amount <= 0.0 {
		return nil, fmt.Errorf("the amount of the payment can't be negative: %v", request.Amount)
	}

	// Manejar campos nulos
	clientPhone := ""
	if request.ClientPhone != nil {
		clientPhone = *request.ClientPhone
	}

	commercePhone := ""
	if request.CommercePhone != nil {
		commercePhone = *request.CommercePhone
	}

	debtorAccount := ""
	if request.DebtorAccount != nil {
		debtorAccount = *request.DebtorAccount
	}

	originBankCode := ""
	if request.OriginBankCode != nil {
		originBankCode = *request.OriginBankCode
	}

	fmt.Println("-> Creando modelo de notificación")
	notificacion := models.NotificationBancaribe{
		Amount:               &request.Amount,
		BankName:             &request.BankName,
		ClientPhone:          &clientPhone,
		CommercePhone:        &commercePhone,
		CreditorAccount:      &request.CreditorAccount,
		CurrencyCode:         &request.CurrencyCode,
		DateBancaribe:        &request.Date,
		Date:                 &*fecha,
		DebtorAccount:        &debtorAccount, // Nuevo campo
		DebtorID:             &request.DebtorID,
		DestinyBankReference: &request.DestinyBankReference,
		OriginBankCode:       &originBankCode,
		OriginBankReference:  &request.OriginBankReference,
		PaymentType:          &request.PaymentType,
		TimeBancaribe:        &request.Time,
		Time:                 hora,
	}

	return &notificacion, nil
}

/////////////////////////////////////////////////

func TransformDate(date string) (*time.Time, error) {
	fmt.Printf("-> Parseando fecha: '%s'\n", date)
	parseDate, err := time.Parse("02-01-2006", date)
	if err != nil {
		parseDate, err = time.Parse("2006-01-02", date)
		if err != nil {
			return nil, fmt.Errorf("fecha no válida: %v (formatos soportados: 'dd-mm-yyyy' o 'yyyy-mm-dd')", date)
		}
	}
	fmt.Printf("-> Fecha parseada: %v\n", parseDate.Format(time.RFC3339))
	return &parseDate, nil
}

/////////////////////////////////////////////////

func TransformHour(timeStr string) (*time.Time, error) {
	fmt.Printf("-> Parseando hora: '%s'\n", timeStr)
	layouts := []string{
		"15.04", "15:04", "1504", "15 04",
		"15.04.00", "15:04:00", "150400", "15 04 00",
		"15.04.05", "15:04:05", // Nuevos formatos con segundos
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, timeStr)
		if err == nil {
			if t.Hour() < 0 || t.Hour() > 23 || t.Minute() < 0 || t.Minute() > 59 {
				return nil, fmt.Errorf("hora o minutos fuera de rango")
			}
			fmt.Printf("-> Hora parseada (%s): %v\n", layout, t.Format("15:04:05"))
			return &t, nil
		}
	}

	return nil, fmt.Errorf("formato de hora inválido, formatos aceptados: HH.MM, HH:MM, HHMM, HH MM, etc")
}

//////////////////////////////////////////////////.
//////////////////////////////////////////////////.
//////////////////////////////////////////////////.

func CheckNotificationExists(model models.NotificationBancaribe, db *gorm.DB) (bool, error) {
	fmt.Println("-> Buscando duplicados en BD...")
	fmt.Printf(">> Parámetros: OriginBankCode=%s | DestinyBankRef=%s | Date=%v | Amount=%v | OriginBankRef=%s\n",
		model.OriginBankCode, model.DestinyBankReference, model.Date, model.Amount, model.OriginBankReference)

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

	fmt.Printf("-> Registros encontrados: %d\n", count)
	return count > 0, nil
}

//////////////////////////////////////////////////.
//////////////////////////////////////////////////.
//////////////////////////////////////////////////.

func saveNotification(model *models.NotificationBancaribe, db *gorm.DB) (bool, error) {
	fmt.Println("-> Guardando en BD...")
	fmt.Printf(">> Datos: %+v\n", model)

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

	fmt.Printf("-> Filas afectadas: %d\n", result.RowsAffected)
	fmt.Printf("-> ID generado: %v\n", model.ID)

	if result.RowsAffected == 0 {
		log.Println("No rows were affected when saving notification")
		return false, fmt.Errorf("no rows affected")
	}

	return true, nil
}

//////////////////////////////////////////////////.
//////////////////////////////////////////////////.
//////////////////////////////////////////////////.

func (n *notificationBancaribe) Validate() error {
	fmt.Println("-> Validando campos obligatorios...")

	// Campos siempre requeridos
	requiredFields := map[string]string{
		"amount":               fmt.Sprintf("%v", n.Amount),
		"bankName":             n.BankName,
		"creditorAccount":      n.CreditorAccount,
		"currencyCode":         n.CurrencyCode,
		"date":                 n.Date,
		"destinyBankReference": n.DestinyBankReference,
		"originBankReference":  n.OriginBankReference,
		"paymentType":          n.PaymentType,
		"time":                 n.Time,
	}

	for field, value := range requiredFields {
		if value == "" || value == "0" {
			return fmt.Errorf("%s es obligatorio", field)
		}
	}

	fmt.Println("-> Validación exitosa")
	return nil
}
