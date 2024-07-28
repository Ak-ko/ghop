package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_HOST string
	DB_PORT string

	DB_USER string
	DB_PASSWORD string
	DB_NAME string
}

var ENV = initConfig()

func initConfig() Config {
	err := godotenv.Load(); if err != nil {
		log.Fatal(err)
	}

	return Config{
		DB_HOST: getEnv("DB_HOST", "127.0.0.1"),
		DB_PORT: getEnv("DB_PORT", "3306"),
		DB_USER: getEnv("DB_USER", "root"),
		DB_PASSWORD: getEnv("DB_PASSWORD", ""),
		DB_NAME: getEnv("DB_NAME", "go"),
	}
}

func getEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key); if ok {
		return value
	}
	return fallback
}