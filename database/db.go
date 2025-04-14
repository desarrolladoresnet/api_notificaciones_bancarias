package db

import (
	"fmt"
	"log"
	"os"

	"github.com/desarrolladoresnet/api_notificaciones_bancarias/src/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func Database() (*gorm.DB, error) {
	fmt.Println("Connecting to the database")
	// Obtener las variables de entorno
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	timezone := os.Getenv("DB_TIMEZONE")
	development := os.Getenv("DEVELOPMENT")
	automigrate := os.Getenv("AUTOMIGRATE")

	// Crear el DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		host, user, password, dbname, port, timezone)
	fmt.Println(dsn)

	// Configuración de Gorm (puedes ajustar el logger según tus necesidades)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Puedes cambiar el nivel del logger si es necesario
	})
	if err != nil {
		fmt.Printf("could not connect to the database:\n %v", err)
		return nil, fmt.Errorf("could not connect to the database: %w", err)
	}

	// Opcionalmente, configurar la conexión a la base de datos (por ejemplo, conexión máxima, tiempo de espera, etc.)
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("could not configure the database connection: %w", err)
	}

	log.Println("connected")
	// db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("running migrations")

	log.Printf("Development: %v\nAutomigrate: %v\n", development, automigrate)

	if development == "true" {
		log.Println("Development mode")
		if automigrate == "true" {
			log.Println("Automigrate enabled")
			err = AutoMigrateDB(db)
			if err != nil {
				return nil, fmt.Errorf("error running migrations: %w", err)
			}
		}
	}

	// CreateDefaultAdmin(db)

	// Configurar los parámetros de la conexión, como máximo número de conexiones abiertas, etc.
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(50)
	sqlDB.SetConnMaxLifetime(5 * 60) // Ejemplo de 5 minutos

	return db, nil
}

func AutoMigrateDB(db *gorm.DB) error {
	//* ELIMINAR TABLAS *//
	log.Println("Eliminando tabla intermedia si existe")
	err := db.Migrator().DropTable("notificacion_bdv")
	if err != nil {
		log.Fatalf("Error eliminando la tabla user_clients: %v", err)
	}

	//* MIGRAR TABLAS *//

	err = db.AutoMigrate(&models.NotificationBDV{})
	if err != nil {
		return fmt.Errorf("error al migrar las tablas: %w", err)
	}

	log.Println("Migraciones completadas")
	return nil
}
