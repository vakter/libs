package dbhandler

import (
    "log"
    "os"

    "github.com/go-pg/pg/v10"
    "github.com/golang-migrate/migrate/v4"
    migratepg "github.com/golang-migrate/migrate/v4/database/pg"
//"github.com/golang-migrate/migrate/v4/database/pgx"
    //"github.com/golang-migrate/migrate/v4/source/file"
)

// DBHandler handles the PostgreSQL database connection
type DBHandler struct {
    Conn *pg.DB
}

// ConnectPg establishes a connection to the PostgreSQL database.
func (db *DBHandler) ConnectPg() {
    options := &pg.Options{
        Addr:     getEnv("DB_ADDR", "localhost:5432"),
        User:     getEnv("DB_USER", "postgres"),
        Password: getEnv("DB_PASSWORD", ""),
        Database: getEnv("DB_NAME", "mydb"),
    }

    db.Conn = pg.Connect(options)

    // Check connection
    _, err := db.Conn.Exec("SELECT 1")
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
}

// RunMigrations runs database migrations using golang-migrate.
func (db *DBHandler) RunMigrations() {
    // Initialize the migrate instance
    m, err := migrate.New(
        "file://migrations", // Path to your migration files
        "postgres://"+getEnv("DB_USER", "postgres")+":"+getEnv("DB_PASSWORD", "")+
            "@"+getEnv("DB_ADDR", "localhost:5432")+"/"+getEnv("DB_NAME", "mydb")+"?sslmode=disable",
    )
    if err != nil {
        log.Fatalf("Failed to initialize migration: %v", err)
    }

    // Run the migrations
    err = m.Up()
    if err != nil && err != migrate.ErrNoChange {
        log.Fatalf("Failed to run migrations: %v", err)
    }
}

// Close closes the database connection
func (db *DBHandler) Close() {
    if db.Conn != nil {
        db.Conn.Close()
    }
}

// getEnv retrieves the value of the environment variable named by the key or fallback to the default value
func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}
