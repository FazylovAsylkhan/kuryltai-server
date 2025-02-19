package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DbURL     string
	SecretKey string
}

var (
	config Config
	once   sync.Once
)

const minSecretKeySize = 32

func Get() *Config {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Fatal("No .env file found")
		}
		portString, exists := os.LookupEnv("PORT")
		if !exists {
			log.Fatal("PORT is not found in the environment")
		}
		dbURL, exists := os.LookupEnv("DB_URL")
		if !exists {
			log.Fatal("DB_URL is not found in the environment")
		}
		secretKey, exists := os.LookupEnv("SECRET_KEY")
		if !exists {
			log.Fatal("SECRET_KEY is not found in the environment")
		}
		if len(secretKey) < minSecretKeySize {
			log.Fatalf("SECRET_KEY must be at least %d characters", minSecretKeySize)
		}
		config = Config{
			Port:      portString,
			DbURL:     dbURL,
			SecretKey: secretKey,
		}
	})

	return &config
}
