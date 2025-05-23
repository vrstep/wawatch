// backend/controller/main_test.go (or backend/main_test.go)
package controller_test // Or your test package name

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver for database/sql
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/vrstep/wawatch-backend/config"
	"github.com/vrstep/wawatch-backend/controller"
	"github.com/vrstep/wawatch-backend/middleware"
	"github.com/vrstep/wawatch-backend/routes"
	// Add models import if clearUserRelatedTables is in this file
)

var testRouter *gin.Engine
var testDB *gorm.DB

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Printf("INFO: Environment variable %s not set, using default: %s", key, fallback)
	return fallback
}

// findProjectRoot tries to find the project root by looking for go.mod
func findProjectRoot() string {
	_, currentFilePath, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("FATAL: Could not get current file path to determine project root.")
	}
	currentDir := filepath.Dir(currentFilePath) // Directory of the current test file

	// Walk up until go.mod is found
	for {
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			log.Printf("INFO: Project root found at: %s", currentDir)
			return currentDir // Found project root
		}
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir { // Reached filesystem root
			log.Fatalf("FATAL: Could not find go.mod by walking up from %s. Ensure tests are within the Go module.", filepath.Dir(currentFilePath))
		}
		currentDir = parentDir
	}
}

func setupTestEnvironment() {
	gin.SetMode(gin.TestMode)

	projectRoot := findProjectRoot() // Use the helper to find the module root

	envPath := filepath.Join(projectRoot, ".env.test")
	if err := godotenv.Load(envPath); err == nil {
		log.Printf("INFO: Loaded test environment from '%s'", envPath)
	} else {
		log.Printf("INFO: Could not load .env.test file from '%s' (Error: %v). Relying on OS environment variables or defaults for tests.", envPath, err)
	}

	testDbUser := getEnv("TEST_DB_USER", "postgres")
	testDbPassword := getEnv("TEST_DB_PASSWORD", "postgres")
	testDbHost := getEnv("TEST_DB_HOST", "localhost")
	testDbPort := getEnv("TEST_DB_PORT", "5435")
	testDbName := getEnv("TEST_DB_NAME", "wawatchdb_test_user_svc")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		testDbUser, testDbPassword, testDbHost, testDbPort, testDbName)

	var err error
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to test database with GORM: %v\nDSN: %s", err, dsn)
	}
	config.DB = testDB
	log.Println("INFO: Successfully connected to test database.")

	// --- Run Migrations ---
	migrationsDir := filepath.Join(projectRoot, "db", "migrations") // Path relative to determined project root
	absMigrationsDir, err := filepath.Abs(migrationsDir)
	if err != nil {
		log.Fatalf("FATAL: Could not get absolute path for migrations directory '%s': %v", migrationsDir, err)
	}

	if _, errStat := os.Stat(absMigrationsDir); os.IsNotExist(errStat) {
		log.Fatalf("FATAL: Migrations directory does not exist at resolved absolute path: %s. (Constructed from project root: %s, relative path: db/migrations)", absMigrationsDir, projectRoot)
	}

	migrationsPath := fmt.Sprintf("file://%s", absMigrationsDir) // Use absolute path
	log.Printf("INFO: Attempting to run migrations for user-service test DB from: %s", migrationsPath)

	mInstance, err := migrate.New(migrationsPath, dsn) // dsn should be defined earlier
	if err != nil {
		log.Fatalf("FATAL: Failed to create migrate instance: %v\nDSN: %s\nPath: %s", err, dsn, migrationsPath)
	}
	log.Println("INFO: Applying migrations up...")
	if err := mInstance.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("FATAL: Failed to run migrations up: %v", err)
	}
	version, dirty, _ := mInstance.Version()
	log.Printf("INFO: Test DB migrations applied. Version: %d, Dirty: %t", version, dirty)

	controller.InitAnimeServiceClient()
	log.Printf("INFO: AnimeServiceClient initialized.")

	testRouter = gin.New()
	testRouter.Use(gin.Recovery())
	testRouter.Use(middleware.RequestID())

	routes.UserRoutes(testRouter)
	routes.UserAnimeListRoutes(testRouter)
	routes.AnimePassThroughRoutes(testRouter)
	routes.UserHistoryRoutes(testRouter)
	log.Println("INFO: Test router configured.")
}

func teardownTestEnvironment() {
	log.Println("INFO: Tearing down test environment...")

	// Use the same robust path finding for teardown
	projectRoot := findProjectRoot() // Use the helper
	migrationsDir := filepath.Join(projectRoot, "db", "migrations")
	absMigrationsDir, errPath := filepath.Abs(migrationsDir)
	if errPath != nil {
		log.Printf("WARNING: Could not get absolute path for migrations directory during teardown: %v", errPath)
		// Close DB and exit if path finding fails for teardown
		sqlDB, err := testDB.DB()
		if err == nil && sqlDB != nil {
			sqlDB.Close()
			log.Println("INFO: Closed test database connection (due to teardown path error).")
		}
		return
	}
	migrationsPath := fmt.Sprintf("file://%s", absMigrationsDir)

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getEnv("TEST_DB_USER", "postgres"),
		getEnv("TEST_DB_PASSWORD", "postgres"),
		getEnv("TEST_DB_HOST", "localhost"),
		getEnv("TEST_DB_PORT", "5435"),
		getEnv("TEST_DB_NAME", "wawatchdb_test_user_svc"))

	mInstance, err := migrate.New(migrationsPath, dsn)
	if err == nil {
		log.Println("INFO: Rolling back migrations (down)...")
		if err := mInstance.Down(); err != nil && err != migrate.ErrNoChange && !strings.Contains(err.Error(), "file does not exist") { // Be careful with "file does not exist"
			log.Printf("WARNING: Failed to run migrations down for test DB: %v", err)
		} else if err == nil || err == migrate.ErrNoChange {
			log.Println("INFO: Test DB migrations successfully rolled back or no changes.")
		}
		// Check for close errors from migrate instance
		if srcErr, dbErr := mInstance.Close(); srcErr != nil || dbErr != nil {
			log.Printf("WARNING: Error closing migrate instance during teardown. SourceErr: %v, DBErr: %v", srcErr, dbErr)
		}
	} else {
		log.Printf("WARNING: Could not create migrate instance for teardown: %v. DSN: %s, Path: %s", err, dsn, migrationsPath)
	}

	sqlDB, errDb := testDB.DB() // Renamed err to errDb to avoid conflict
	if errDb == nil && sqlDB != nil {
		sqlDB.Close()
		log.Println("INFO: Closed test database connection.")
	}
}

func TestMain(m *testing.M) {
	setupTestEnvironment()
	exitCode := m.Run()
	teardownTestEnvironment()
	os.Exit(exitCode)
}

// Helper function to clear specific tables before each test if needed.
// Call this at the beginning of your test functions.
func clearUserRelatedTables() { // Renamed from clearUserRelatedTables(db *gorm.DB) to use global testDB
	if testDB == nil {
		log.Println("WARNING: clearUserRelatedTables called but testDB is nil.")
		return
	}
	// Order is important if you have foreign key constraints and don't use CASCADE
	// Use Raw SQL for TRUNCATE for efficiency and to reset sequences if needed.
	// Or delete, which is slower but might respect GORM hooks if any.
	testDB.Exec("TRUNCATE TABLE user_view_histories RESTART IDENTITY CASCADE;")
	testDB.Exec("TRUNCATE TABLE user_anime_lists RESTART IDENTITY CASCADE;")
	testDB.Exec("TRUNCATE TABLE anime_caches RESTART IDENTITY CASCADE;") // This is the user-service's local anime cache
	testDB.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE;")
	// testDB.Exec("TRUNCATE TABLE schema_migrations RESTART IDENTITY CASCADE;") // Usually not needed to truncate this

	log.Println("INFO: User-related tables truncated for test.")
}

// You will also need the helper functions from the previous response in a _test.go file:
// - generateTestToken(...)
// - performRequest(...)
// - performAuthRequest(...)
// - createAndLoginTestUser(...)
