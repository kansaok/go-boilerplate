package db

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

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

func ConnectDB(config *Config) (interface{}, error) {
	var err error

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

func GetDB() interface{} {
	if dbConn == nil {
		log.Fatal("Database connection is not initialized. Call ConnectDB first.")
	}
	return dbConn
}

func connectSQL(config *Config) (*gorm.DB, error) {
	var dsn string
	var dialector gorm.Dialector

	switch config.DBConnection {
	case "postgres":
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort, config.DBSSLMode)
		dialector = postgres.Open(dsn)
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
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

func connectMongoDB(config *Config) (*mongo.Client, error) {
	uri := fmt.Sprintf("mongodb://%s:%s", config.DBHost, config.DBPort)

	clientOptions := options.Client().ApplyURI(uri)

	if config.DBUser != "" && config.DBPassword != "" {
		credential := options.Credential{
			Username: config.DBUser,
			Password: config.DBPassword,
		}
		clientOptions.SetAuth(credential)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB database successfully!")
	return client, nil
}
