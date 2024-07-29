package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_HOST string
	DB_PORT string

	DB_USER string
	DB_PASSWORD string
	DB_NAME string
	JWT_EXPIRATION int64
	JWT_SECRET string
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
		JWT_SECRET: getEnv("JWT_SECRET", "secret"),
		JWT_EXPIRATION: getEnvAsInt("JWT_EXPIRATION", 3600 * 24 * 7),
	}
}

func getEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key); if ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}