package application

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
)

type PostgresConnector struct {
}

func (p *PostgresConnector) GetConnection() (db *gorm.DB, err error) {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST_READER")
	dbPort := os.Getenv("DB_PORT")
	if username == "" || password == "" || dbName == "" || dbHost == "" || dbPort == "" {
		return nil, fmt.Errorf("Missing environment variables for connecting to database")
	}
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable password=%s", dbHost, username, dbName, dbPort, password)
	return gorm.Open("postgres", dbURI)
}
