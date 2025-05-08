package config

import (
	"database/sql"
	"fmt"
	"log"
	"os" // Import os package to read environment variables

	"github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Import file source driver
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB holds the database connection pool for the anime service
var DB *gorm.DB

func ConnectDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Build DSN from environment variables
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname)

	// Set ANIME_DB_DSN for compatibility with existing code
	os.Setenv("ANIME_DB_DSN", dsn)

	// Use the DSN directly rather than reading ANIME_DB_DSN again
	gormDB, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{})

	// dsn := os.Getenv("ANIME_DB_DSN")
	// if dsn == "" {
	// 	dsn = "postgres://postgres:postgres@localhost:5442/wawatch_animeservicedb?sslmode=disable" // Default DSN
	// 	log.Println("Warning: ANIME_DB_DSN environment variable not set, using default DSN.")
	// }

	// gormDB, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	log.Fatalf("ANIME_SERVICE: Failed to connect to database with GORM: %v", err)
	// }

	// --- Run Migrations ---
	// Use the same DSN for a separate connection for migrations
	rawDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("ANIME_SERVICE: Failed to open sql.DB for migration: %v", err)
	}
	defer rawDB.Close() // Close this connection after migrations are done

	log.Println("ANIME_SERVICE: Running database migrations...")
	// Pass the rawDB connection to the migration runner
	if err := runMigrations(rawDB); err != nil {
		log.Fatalf("ANIME_SERVICE: Failed to run database migrations: %v", err)
	}
	log.Println("ANIME_SERVICE: Database migrations completed successfully.")

	// Assign the GORM DB instance to the global variable
	DB = gormDB
}

// runMigrations executes the database migrations for the anime service
func runMigrations(sqlDB *sql.DB) error {
	driver, err := migratepostgres.WithInstance(sqlDB, &migratepostgres.Config{})
	if err != nil {
		log.Printf("ANIME_SERVICE: Error creating migrate driver instance: %v", err)
		return err
	}

	// *** IMPORTANT: Point to the migrations directory WITHIN the anime-service ***
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations", // Assumes anime-service/db/migrations exists
		"postgres",             // Target database type (use "postgres" for identifier)
		driver,
	)
	if err != nil {
		log.Printf("ANIME_SERVICE: Error creating migrate instance: %v", err)
		return err
	}

	// Apply all available "up" migrations
	log.Println("ANIME_SERVICE: Applying migrations...")
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Printf("ANIME_SERVICE: Error applying migrations: %v", err)
		return err
	}

	// Check for migration source/database errors after running
	sourceErr, dbErr := m.Close()
	if sourceErr != nil {
		log.Printf("ANIME_SERVICE: Migration source error on close: %v", sourceErr)
		// Decide if this should be a fatal error or just a warning
	}
	if dbErr != nil {
		log.Printf("ANIME_SERVICE: Migration database error on close: %v", dbErr)
		// Decide if this should be a fatal error or just a warning
	}

	if err == migrate.ErrNoChange {
		log.Println("ANIME_SERVICE: No new migrations to apply.")
		return nil // Return nil as no change is not an error
	}

	// If err was not nil and not ErrNoChange, it's returned above
	// If we reach here and err was nil (meaning migrations were applied), return nil
	return nil
}
