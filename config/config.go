package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost              string
	Port                    string
	CookiesAuthSecret       string
	CookiesAuthAgeInSeconds int
	CookiesAuthIsSecure     bool
	CookiesAuthIsHttpOnly   bool
	AzureADClientID         string
	AzureADClientSecret     string
	AzureADTenantID         string
	GitHubToken             string
	BlogPostConfig          string
}

const (
	twoDaysInSeconds = 60 * 60 * 24 * 2
)

var Envs Config
var blogConfig string

func init() {
	// Load the .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	} else {
		fmt.Println("Loaded .env file")
	}
	Envs = initConfig()
}

func initConfig() Config {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
    log.Fatalf("Error reading config.yaml: %v", err)
	}
	blogConfig = string(data)

	return Config{
		PublicHost:              getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                    getEnv("PORT", "7000"),
		CookiesAuthSecret:       getEnv("COOKIES_AUTH_SECRET", "some-secret-key"),
		CookiesAuthAgeInSeconds: getEnvAsInt("COOKIES_AUTH_AGE_IN_SECONDS", twoDaysInSeconds),
		CookiesAuthIsSecure:     getEnvAsBool("COOKIES_AUTH_IS_SECURE", false),
		CookiesAuthIsHttpOnly:   getEnvAsBool("COOKIES_AUTH_IS_HTTP_ONLY", false),
		AzureADClientID:         getEnvOrError("AZURE_AD_CLIENT_ID"),
		AzureADClientSecret:     getEnvOrError("AZURE_AD_CLIENT_SECRET"),
		AzureADTenantID:         getEnvOrError("AZURE_AD_TENANT_ID"),
		GitHubToken:             getEnvOrError("GITHUB_TOKEN"),
		BlogPostConfig:          blogConfig,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvOrError(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	// should server panic or not?
	// I thing it should since the environment variable is required to run the application
	panic(fmt.Sprintf("Environment variable %s is not set", key))
}

func getEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fallback
		}
		return b
	}

	return fallback
}
