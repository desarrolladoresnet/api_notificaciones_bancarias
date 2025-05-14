package db

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/desarrolladoresnet/api_notificaciones_bancarias/src/api_key"
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
	log.Println("Eliminando tablas si existen...")

	tablesToDrop := []string{"api_key", "api_user", "notificacion_bdv", "notificacion_bancaribe", "notificacion_tesoro", "user_tesoro"}
	for _, table := range tablesToDrop {
		if err := db.Migrator().DropTable(table); err != nil {
			log.Printf("Error eliminando la tabla %s: %v", table, err)
			// No usamos log.Fatalf para evitar que termine el programa abruptamente
			return fmt.Errorf("error eliminando tabla %s: %w", table, err)
		}
	}

	log.Println("Iniciando migración de tablas...")

	if err := db.AutoMigrate(&models.NotificationBDV{}, &models.APIKey{}, models.NotificationBancaribe{}, models.APIUser{}, models.NotificacionTesoro{}); err != nil {
		return fmt.Errorf("error al migrar las tablas: %w", err)
	}

	err := createDefaultUser(db)
	if err != nil {
		return err
	}

	log.Println("✅ Migraciones completadas con éxito")
	return nil
}

func createDefaultUser(db *gorm.DB) error {
	fmt.Println("Creating default user...")

	// Obtener las variables de entorno
	username := os.Getenv("DEFAULT_USERNAME")
	password := os.Getenv("DEFAULT_PASSWORD")

	if username == "" || password == "" {
		return errors.New("DEFAULT_USERNAME and/or DEFAULT_PASSWORD environment variables not set")
	}

	// Verificar si ya existe un usuario con ese username
	var existingUser models.APIUser
	result := db.Where("username = ? AND is_active = ?", username, true).First(&existingUser)
	if result.Error == nil {
		fmt.Printf("User %s already exists, skipping creation\n", username)
		return nil
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error checking existing user: %v", result.Error)
	}

	// Crear el usuario
	err := api_key.CreateDefaultUser(db, username, password)
	if err != nil {
		return fmt.Errorf("failed to create default user: %w", err)
	}

	log.Printf("✅ Default user %s created successfully\n", username)
	return nil
}
