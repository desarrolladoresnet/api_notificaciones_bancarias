package router_module

import (
	"github.com/desarrolladoresnet/api_notificaciones_bancarias/src/bdv"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Router(api *gin.RouterGroup, db *gorm.DB) {
	bdv_routes := api.Group("/bdv")

	bdv_routes.POST("/webhook", bdv.WeebHookBDV(db))
}
