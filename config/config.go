package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

type CollectionConfig []Collection

type Collection struct {
	Collection     string       `yaml:"collection"`
	GitPath        string       `yaml:"git_path"`
	FileNameFormat string       `yaml:"file_name_format"`
	MetadataSchema []SchemaItem `yaml:"metadata_schema"`
}

type SchemaItem struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Type        string `yaml:"type"`
	Required    bool   `yaml:"required"`
	Filter      bool   `yaml:"filter"`
	ItemsType   string `yaml:"items_type,omitempty"`
}

func (c *Collection) ToYAMLString() (string, error) {
	data, err := yaml.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("error marshaling to YAML: %w", err)
	}
	return string(data), nil
}

func (c *CollectionConfig) GetCollectionConfig(collection string) *Collection {
	for _, col := range *c {
		if col.Collection == collection {
			return &col
		}
	}
	return nil
}

type UserConfig []User

type User struct {
	Email                 string                 `yaml:"email"`
	IsAdmin               bool                   `yaml:"is_admin"`
	CollectionPermissions []CollectionPermission `yaml:"collection_permissions"`
}

type CollectionPermission struct {
	Name            string                  `yaml:"name"`
	Permission      string                  `yaml:"permission"`
	DefaultMetadata map[string]MetadataItem `yaml:"default_metadata,omitempty"`
}

type MetadataItem struct {
	Value  interface{} `yaml:"value"`
	Strict bool        `yaml:"strict"`
}

func (c *CollectionPermission) ToYAMLString() (string, error) {
	data, err := yaml.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("error marshaling to YAML: %w", err)
	}
	return string(data), nil
}

func (u *UserConfig) GetUserConfig(email string) *User {
	for _, user := range *u {
		if user.Email == email {
			return &user
		}
	}
	return nil
}

func (u *User) GetCollectionPermission(collection string) *CollectionPermission {
	for _, col := range u.CollectionPermissions {
		if col.Name == collection {
			return &col
		}
	}

	return nil
}

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
	CollectionConfig        *CollectionConfig
	UserConfig              *UserConfig
}

const (
	twoDaysInSeconds = 60 * 60 * 24 * 2
)

var Envs Config

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
	configData, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(fmt.Sprintf("Error reading config.yaml: %v", err))
	}
	var collectionConfig CollectionConfig
	if err := yaml.Unmarshal(configData, &collectionConfig); err != nil {
		panic(fmt.Sprintf("Error unmarshaling config.yml: %v", err))
	}

	userConfigData, err := os.ReadFile("users.yaml")
	var userConfig UserConfig
	if err := yaml.Unmarshal(userConfigData, &userConfig); err != nil {
		panic(fmt.Sprintf("Error unmarshaling users.yaml: %v", err))
	}

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
		CollectionConfig:        &collectionConfig,
		UserConfig:              &userConfig,
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
