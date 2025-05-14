package router_module

import (
	"github.com/desarrolladoresnet/api_notificaciones_bancarias/src/api_key"
	"github.com/desarrolladoresnet/api_notificaciones_bancarias/src/bancaribe"
	"github.com/desarrolladoresnet/api_notificaciones_bancarias/src/bdv"
	"github.com/desarrolladoresnet/api_notificaciones_bancarias/src/middleware"
	"github.com/desarrolladoresnet/api_notificaciones_bancarias/src/tesoro"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Router(api *gin.RouterGroup, db *gorm.DB) {
	bdv_routes := api.Group("/bdv")

	// Aplicar el middleware solo a rutas que deben protegerse
	// ---- BDV ---- //
	bdv_routes.Use(middleware.APIKeyAuthMiddlewareBDV(db))
	{
		bdv_routes.POST("/webhook", bdv.WeebHookBDV(db))
		bdv_routes.GET("/notificaciones", bdv.GetPayments(db))
	}

	// ---- BANCARIBE ---- //
	bancaribe_routes := api.Group("/bancaribe")
	bdv_routes.Use(middleware.APIKeyAuthMiddlewareBancaribe(db))
	{
		bancaribe_routes.POST("/webhook", bancaribe.WeebHookBancaribe(db))
		bancaribe_routes.GET("/notificaciones", bancaribe.GetPayments(db))
	}

	// ---- TESORO ---- //
	tesoro_routes := api.Group("/tesoro")
	tesoro_routes.Use(middleware.APIKeyAuthMiddlewareTesoro(db))
	{
		tesoro_routes.POST("/webhook", tesoro.WeebhookTesoro(db))
		// tesoro_routes.GET("/notificaciones", bancaribe.GetPayments(db))
	}

	// ---- API KEY ROUTES ---- //
	api_key_routes := api.Group("/access")
	api_key_routes.POST("/login", api_key.GetApiKey(db))
	api_key_routes.Use(middleware.APIKeyAuthMiddlewareApiKey(db))
	{
		api_key_routes.POST("/user", api_key.CreateApiUser(db))
		api_key_routes.POST("/api-key", api_key.CreateApiKey(db))
		api_key_routes.GET("/api-key", api_key.GetAPIKeys(db))
		api_key_routes.DELETE("/api-key:id", api_key.DeleteAPIKey(db))
	}

}
