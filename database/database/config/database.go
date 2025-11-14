package config

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "time"

    _ "github.com/go-sql-driver/mysql"
    "github.com/joho/godotenv"
)

// Config represents database configuration
type Config struct {
    Host        string
    Port        string
    User        string
    Password    string
    Database    string
    Charset     string
    ParseTime   bool
    Loc         string
    MaxOpenCons int
    MaxIdleCons int
    MaxLifetime time.Duration
}

// NewConfig creates configuration based on environment
func NewConfig(environment string) *Config {
    // Load environment variables from .env file
    godotenv.Load()
    
    baseConfig := &Config{
        Charset:     "utf8mb4",
        ParseTime:   true,
        Loc:         "Local",
        MaxOpenCons: 25,
        MaxIdleCons: 5,
        MaxLifetime: 5 * time.Minute,
    }
    
    switch environment {
    case "production":
        baseConfig.Host = getEnv("PRODUCTION_DB_HOST", "localhost")
        baseConfig.Port = getEnv("PRODUCTION_DB_PORT", "3306")
        baseConfig.User = getEnv("PRODUCTION_DB_USER", "planica_user")
        baseConfig.Password = getEnv("PRODUCTION_DB_PASSWORD", "")
        baseConfig.Database = getEnv("PRODUCTION_DB_NAME", "planica_bi")
    case "staging":
        baseConfig.Host = getEnv("STAGING_DB_HOST", "localhost")
        baseConfig.Port = getEnv("STAGING_DB_PORT", "3306")
        baseConfig.User = getEnv("STAGING_DB_USER", "planica_user")
        baseConfig.Password = getEnv("STAGING_DB_PASSWORD", "")
        baseConfig.Database = getEnv("STAGING_DB_NAME", "planica_bi")
    default: // development
        baseConfig.Host = getEnv("DB_HOST", "localhost")
        baseConfig.Port = getEnv("DB_PORT", "3306")
        baseConfig.User = getEnv("DB_USER", "planica_user")
        baseConfig.Password = getEnv("DB_PASSWORD", "root")
        baseConfig.Database = getEnv("DB_NAME", "planica_bi")
    }
    
    return baseConfig
}

// GetConnectionString returns MySQL DSN
func (c *Config) GetConnectionString() string {
    return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
        c.User, c.Password, c.Host, c.Port, c.Database, c.Charset, c.ParseTime, c.Loc)
}

// Connect establishes database connection with connection pooling
func Connect(environment string) (*sql.DB, error) {
    config := NewConfig(environment)
    
    db, err := sql.Open("mysql", config.GetConnectionString())
    if err != nil {
        return nil, fmt.Errorf("failed to open database connection: %v", err)
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(config.MaxOpenCons)
    db.SetMaxIdleConns(config.MaxIdleCons)
    db.SetConnMaxLifetime(config.MaxLifetime)
    
    // Verify connection
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %v", err)
    }
    
    log.Printf("✅ Successfully connected to %s database: %s", environment, config.Database)
    return db, nil
}

// Repository provides database operations
type Repository struct {
    db *sql.DB
}

// NewRepository creates new repository instance
func NewRepository(db *sql.DB) *Repository {
    return &Repository{db: db}
}

// HealthCheck verifies database connectivity
func (r *Repository) HealthCheck() error {
    return r.db.Ping()
}

// GetDatabaseVersion returns MySQL version
func (r *Repository) GetDatabaseVersion() (string, error) {
    var version string
    err := r.db.QueryRow("SELECT VERSION()").Scan(&version)
    return version, err
}

// GetMigrationStatus returns current migration status
func (r *Repository) GetMigrationStatus() ([]map[string]interface{}, error) {
    query := `
        SELECT version, description, applied_at, applied_by, status, execution_time_ms
        FROM schema_migrations 
        ORDER BY applied_at DESC
        LIMIT 10
    `
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var migrations []map[string]interface{}
    for rows.Next() {
        var (
            version       string
            description   string
            appliedAt     time.Time
            appliedBy     string
            status        string
            executionTime *int
        )
        
        err := rows.Scan(&version, &description, &appliedAt, &appliedBy, &status, &executionTime)
        if err != nil {
            return nil, err
        }
        
        migration := map[string]interface{}{
            "version":        version,
            "description":    description,
            "applied_at":     appliedAt,
            "applied_by":     appliedBy,
            "status":         status,
            "execution_time": executionTime,
        }
        migrations = append(migrations, migration)
    }
    
    return migrations, nil
}

// Helper function to get environment variable with default
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}