package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// EnvConfig carga las variables de entorno desde el archivo .env
// y retorna el valor de la variable de entorno solicitada
// Requiere el nombre de la variable a buscar
func EnvConfig(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv(key)
}
