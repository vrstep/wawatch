package config

import (
	"database/sql" // Import database/sql
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres" // Alias migrate's postgres driver
	_ "github.com/golang-migrate/migrate/v4/source/file"                     // Import file source driver
	gormpostgres "gorm.io/driver/postgres"                                   // Import GORM's postgres driver
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// GORM connection
	// dsn := "postgres://postgres:postgres@localhost:5432/wawatchdb?sslmode=disable"

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Build the DSN from environment variables
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname)

	gormDB, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{})

	// gormDB, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	log.Fatalf("Failed to connect to database with GORM: %v", err)
	// }

	// Run migrations using a **separate** sql.DB
	rawDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open sql.DB for migration: %v", err)
	}
	defer rawDB.Close() // okay to close, it's not used by GORM

	log.Println("Running database migrations...")
	if err := runMigrations(rawDB); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}
	log.Println("Database migrations completed successfully.")

	DB = gormDB
}

// runMigrations function executes the database migrations using an existing sql.DB
func runMigrations(sqlDB *sql.DB) error {
	// Use the aliased migrate postgres driver here
	driver, err := migratepostgres.WithInstance(sqlDB, &migratepostgres.Config{})
	if err != nil {
		return err
	}

	// Point to your existing directory containing migration files
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations", // Correct path to your migrations
		"wawatchdb",            // Database name/identifier (keep as "postgres")
		driver,                 // The database instance driver
	)
	if err != nil {
		return err
	}

	// Apply all available "up" migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	// Check for migration source/database errors after running
	sourceErr, dbErr := m.Close()
	if sourceErr != nil {
		log.Printf("Migration source error on close: %v", sourceErr)
	}
	if dbErr != nil {
		log.Printf("Migration database error on close: %v", dbErr)
	}

	if err == migrate.ErrNoChange {
		log.Println("No new migrations to apply.")
		return nil // Not an actual error
	}

	return err // Return the original error from m.Up() if it wasn't ErrNoChange
}
