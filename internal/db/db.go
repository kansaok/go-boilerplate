package db

import (
	"fmt"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
    dbConn interface{}
    once   sync.Once
)

// Config contains configurations for database connections
type Config struct {
    DBHost       string
    DBPort       string
    DBUser       string
    DBPassword   string
    DBName       string
    DBSSLMode    string
    DBConnection string
    DBFile       string
}

// ConnectDB initializes a database connection with the provided configuration
func ConnectDB(config *Config) (interface{}, error) {
    var err error

    // Ensure that the connection is initialized only once
    once.Do(func() {
        switch config.DBConnection {
        case "postgres", "mysql", "sqlite":
            dbConn, err = connectSQL(config)
        case "mongodb":
            dbConn, err = connectMongoDB(config)
        default:
            err = fmt.Errorf("unsupported database driver: %s", config.DBConnection)
        }
    })

    return dbConn, err
}

// GetDB retrieves the initialized database connection
func GetDB() interface{} {
    if dbConn == nil {
        log.Fatal("Database connection is not initialized. Call ConnectDB first.")
    }
    return dbConn
}

// connectSQL initializes a connection to PostgreSQL/MySQL/SQLite using GORM
func connectSQL(config *Config) (*gorm.DB, error) {
    var dsn string
    var dialector gorm.Dialector

    switch config.DBConnection {
    case "postgres":
        dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
            config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort, config.DBSSLMode)
        dialector = postgres.Open(dsn)
    case "mysql":
        dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
            config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
        dialector = mysql.Open(dsn)
    case "sqlite":
        dsn = config.DBFile
        dialector = sqlite.Open(dsn)
    default:
        return nil, fmt.Errorf("unsupported SQL database driver: %s", config.DBConnection)
    }

    db, err := gorm.Open(dialector, &gorm.Config{})
    if err != nil {
        return nil, err
    }

    log.Printf("Connected to %s database successfully!", config.DBConnection)
    return db, nil
}

// connectMongoDB initializes a connection to MongoDB
func connectMongoDB(config *Config) (*mongo.Client, error) {
    uri := fmt.Sprintf("mongodb://%s:%s", config.DBHost, config.DBPort)
    clientOptions := options.Client().ApplyURI(uri)

    client, err := mongo.Connect(nil, clientOptions)
    if err != nil {
        return nil, err
    }

    if err := client.Ping(nil, nil); err != nil {
        return nil, err
    }

    log.Println("Connected to MongoDB database successfully!")
    return client, nil
}
