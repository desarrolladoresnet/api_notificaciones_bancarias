package router_module

import (
	"github.com/desarrolladoresnet/api_notificaciones_bancarias/src/api_key"
	"github.com/desarrolladoresnet/api_notificaciones_bancarias/src/bdv"
	"github.com/desarrolladoresnet/api_notificaciones_bancarias/src/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Router(api *gin.RouterGroup, db *gorm.DB) {
	bdv_routes := api.Group("/bdv")

	// Aplicar el middleware solo a rutas que deben protegerse
	bdv_routes.Use(middleware.APIKeyAuthMiddlewareBDV(db))
	{
		bdv_routes.POST("/webhook", bdv.WeebHookBDV(db))
		bdv_routes.GET("/notificaciones", bdv.GetPayments(db))
	}

	// ---- API KEY ROUTES ---- //
	api_key_routes := api.Group("/api-key")

	api_key_routes.POST("/", api_key.CreateApiKey(db))
	api_key_routes.GET("/", api_key.GetAPIKeys(db))
	api_key_routes.DELETE("/:id", api_key.DeleteAPIKey(db))

}
