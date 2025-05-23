package dbconnection

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct{
	Host string
	Port string
	Password string
	User string
	DBName string
	SSLMode string
}

func GetConfigs() *PostgresConfig {
	fmt.Println(os.Getwd())
	
	err:= godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	
	pgconfig := &PostgresConfig{
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Password: os.Getenv ("DB_PASS"),
		User: os.Getenv ("DB_USER"),
		SSLMode: os.Getenv ("DB_SSLMODE"),
		DBName: os.Getenv ("DB_NAME"),
	}

	return pgconfig
}

func NewConnection(config *PostgresConfig)(*gorm.DB, error){
	//data source name
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return db,err
	}
	return db, nil
}

