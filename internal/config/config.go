package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
	ResumeID     string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	AccessToken  string
	RefreshToken string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		ClientID:     os.Getenv("HH_CLIENT_ID"),
		ClientSecret: os.Getenv("HH_CLIENT_SECRET"),
		Username:     os.Getenv("HH_USERNAME"),
		Password:     os.Getenv("HH_PASSWORD"),
		ResumeID:     os.Getenv("HH_RESUME_ID"),
		DBHost:       os.Getenv("DB_HOST"),
		DBPort:       os.Getenv("DB_PORT"),
		DBUser:       os.Getenv("DB_USER"),
		DBPassword:   os.Getenv("DB_PASSWORD"),
		DBName:       os.Getenv("DB_NAME"),
		AccessToken:  os.Getenv("HH_ACCESS_TOKEN"),
		RefreshToken: os.Getenv("HH_REFRESH_TOKEN"),
	}, nil
}
