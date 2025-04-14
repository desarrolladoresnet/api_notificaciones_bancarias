package bdv

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type bdvRequest struct {
	BancoOrdenante string `json:"bancoOrdenante"`
	Referencia     string `json:"referenciaBancoOrdenante"`
	IdComercio     string `json:"numeroCliente"`
	NumeroComercio string `json:"numeroCOmercio"`
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
