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
	dbHost := os.Getenv("DB_HOST_WRITER")
	dbPort := os.Getenv("DB_PORT")
	if username == "" {
		fmt.Println("Missing username")
	}
	if password == "" {
		fmt.Println("Missing password")
	}
	if dbName == "" {
		fmt.Println("Missing dbName")
	}
	if dbHost == "" {
		fmt.Println("Missing dbHost")
	}
	if dbPort == "" {
		fmt.Println("Missing dbPort")
	}
	if username == "" || password == "" || dbName == "" || dbHost == "" || dbPort == "" {
		return nil, fmt.Errorf("Missing environment variables for connecting to database")
	}
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable password=%s", dbHost, username, dbName, dbPort, password)
	return gorm.Open("postgres", dbURI)
}
