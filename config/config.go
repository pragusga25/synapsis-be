package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBName            string `mapstructure:"DB_NAME"`
	JWTSecret         string `mapstructure:"JWT_SECRET"`
	MidtransServerKey string `mapstructure:"MIDTRANS_SERVER_KEY"`
	Env               string `mapstructure:"ENV"`
	Port              int    `mapstructure:"PORT"`
}

func LoadConfig() (Config, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		fmt.Println("Error parsing PORT")
	}

	return Config{
		DBHost:            os.Getenv("DB_HOST"),
		DBPort:            os.Getenv("DB_PORT"),
		DBUser:            os.Getenv("DB_USER"),
		DBPassword:        os.Getenv("DB_PASSWORD"),
		DBName:            os.Getenv("DB_NAME"),
		JWTSecret:         os.Getenv("JWT_SECRET"),
		MidtransServerKey: os.Getenv("MIDTRANS_SERVER_KEY"),
		Env:               os.Getenv("ENV"),
		Port:              port,
	}, nil

}
