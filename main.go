package main

import (
	"log"
	"net/http"
	"time"

	db "github.com/desarrolladoresnet/api_notificaciones_bancarias/database"
	router_module "github.com/desarrolladoresnet/api_notificaciones_bancarias/src/router"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Connecting to the BD...")

	database, err := db.Database()
	if err != nil {
		log.Printf("Error while conecting to the db: %v", err)
		log.Println("Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
		database, err = db.Database()
		if err != nil {
			log.Fatalf("Failed to connect to the database after retry: %v", err)
		}
	}

	log.Println(database)

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "api notificaciones UP!",
		})
	})
	api := r.Group("api/v1")
	router_module.Router(api, database)

	log.Println("Starting server on :5000...")
	if err := r.Run("0.0.0.0:5000"); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
