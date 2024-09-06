package store

import (
	"os"
)

type Config struct {
	Port      string
	DBFile    string
	JWTSecret string
}

var Envs = initConfig()

func initConfig() Config {
	return Config{
		Port:      getEnv("PORT", "8080"),
		DBFile:    getEnv("DB_FILE", "cms.db"),  // SQLite database file
		JWTSecret: getEnv("JWT_SECRET", "randomjwtsecretkey"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
