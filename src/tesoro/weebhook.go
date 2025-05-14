package tesoro

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/desarrolladoresnet/api_notificaciones_bancarias/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type notificacionTesoro struct {
	PaymentType string  `json:"payment_type"`
	Reference   string  `json:"reference"`
	SourceBank  string  `json:"source_bank"`
	Amount      float64 `json:"amount"`
	SourcePhone string  `json:"source_phone"`
	PaymenDate  string  `json:"payment_date"`
	Document    string  `json:"document"`
}

func WeebhookTesoro(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req notificacionTesoro

		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"codigo":         nil,
				"message":        "error while receiving JSON data",
				"messageSystem,": err.Error(),
				"statusCode":     http.StatusBadRequest,
				"success":        false,
			})
			return
		}

		flag, err := tesoroIsValid(req)
		if err != nil || !flag {
			fmt.Printf("%v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"codigo":         nil,
				"message":        "error while transforming JSON data",
				"messageSystem,": err.Error(),
				"statusCode":     http.StatusBadRequest,
				"success":        false,
			})
			return
		}

		notificacion_model, err := tranformToTesoroModel(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"codigo":         nil,
				"message":        "error while transforming JSON data",
				"messageSystem,": err.Error(),
				"statusCode":     http.StatusInternalServerError,
				"success":        false,
			})
			return
		}

		result, err := saveNotification(notificacion_model, db)
		if err != nil {
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
			c.JSON(http.StatusBadRequest, gin.H{
				"codigo":         nil,
				"message":        "can't save the data in the DB",
				"messageSystem,": nil,
				"success":        false,
				"statusCode":     http.StatusBadRequest,
			})

		}

		c.JSON(http.StatusCreated, gin.H{
			"codigo":         "00",
			"message":        "Success",
			"messageSystem,": "Notificado",
			"success":        true,
			"data":           notificacion_model,
			"statusCode":     http.StatusOK,
		})

	}
}

////////////////////////////////////////////////
////////////////////////////////////////////////
////////////////////////////////////////////////

// 			FUNCIONES AUXILIARES 			  //

////////////////////////////////////////////////
////////////////////////////////////////////////
////////////////////////////////////////////////

func saveNotification(model *models.NotificacionTesoro, db *gorm.DB) (bool, error) {
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

////////////////////////////////////////////////
////////////////////////////////////////////////
////////////////////////////////////////////////

func tesoroIsValid(req notificacionTesoro) (bool, error) {
	var error_string string
	var flag bool = true

	if req.Amount <= 0.0 {
		error_string += "el monto fue igual o menor a cero\n"
		flag = false
	}

	if strings.TrimSpace(req.Document) == "" {
		error_string += "el campo document esta vacio\n"
		flag = false
	}

	if strings.TrimSpace(req.PaymentType) == "" {
		error_string += "el campo payment_type esta vacio\n"
		flag = false
	}

	if strings.TrimSpace(req.Reference) == "" {
		error_string += "el campor reference esta vacio\n"
		flag = false
	}

	if strings.TrimSpace(req.SourceBank) == "" {
		error_string += "el campo source bank esta vacio"
	}

	if !flag {
		return flag, errors.New(error_string)
	}

	return flag, nil
}

////////////////////////////////////////////////
////////////////////////////////////////////////
////////////////////////////////////////////////

func tranformToTesoroModel(req notificacionTesoro) (*models.NotificacionTesoro, error) {
	fecha, err := TransformDate(req.PaymenDate)
	if err != nil {
		fmt.Println("---- Error en Date -----")
		return nil, err
	}

	if req.Amount <= 0.0 {
		return nil, fmt.Errorf("the amount of the payment can't be negative: %v", req.Amount)
	}

	notificacion := models.NotificacionTesoro{
		PaymentType:     req.PaymentType,
		Reference:       req.Reference,
		SourceBank:      req.SourceBank,
		Amount:          req.Amount,
		SourcePhone:     req.SourcePhone,
		PaymenDate:      req.PaymenDate,
		PaymentDateTime: *fecha,
		Document:        req.Document,
	}

	return &notificacion, nil
}

/////////////////////////////////////////////////

func TransformDate(date string) (*time.Time, error) {
	// Try RFC3339 format (which is what your JSON appears to be using)
	fmt.Println("Date: ")
	fmt.Println(date)

	parseDate, err := time.Parse(time.RFC3339, date)
	if err == nil {
		return &parseDate, nil
	}

	// Try yyyy-mm-dd format
	parseDate, err = time.Parse("2006-01-02", date)
	if err == nil {
		return &parseDate, nil
	}

	// Try dd-mm-yyyy format
	parseDate, err = time.Parse("02-01-2006", date)
	if err == nil {
		return &parseDate, nil
	}

	return nil, fmt.Errorf("fecha no válida: %v (formatos soportados: 'dd-mm-yyyy' o 'yyyy-mm-dd')", date)
}
