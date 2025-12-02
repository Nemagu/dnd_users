package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type WebConfig struct {
	DBHost     string
	DBPort     uint16
	DBUser     string
	DBPassword string
	DBName     string

	HTTPPort    uint16
	HTTPHost    string
	HTTPTimeout time.Duration

	JWTAccessLifetime  time.Duration
	JWTRefreshLifetime time.Duration
	JWTSecretKey       string

	Debug bool

	PasswordSecretKey string
	PasswordSalt      string
	PasswordCost      int

	EmailSecretKey     string
	EmailTokenLifetime time.Duration
	EmailFolderPath    string
	EmailTimeout       time.Duration
}

func MustNewWebConfig() *WebConfig {
	return &WebConfig{
		DBHost:     getEnv("DB_HOST", "localhost", true),
		DBPort:     getEnvAsUint16("DB_PORT", 5432, true),
		DBUser:     getEnv("DB_USER", "postgres", true),
		DBPassword: getEnv("DB_PASSWORD", "postgres", true),
		DBName:     getEnv("DB_NAME", "postgres", true),

		HTTPPort:    getEnvAsUint16("HTTP_PORT", 8080, true),
		HTTPHost:    getEnv("HTTP_HOST", "localhost", true),
		HTTPTimeout: time.Duration(getEnvAsInt("HTTP_TIMEOUT", 5, false)) * time.Second,

		JWTAccessLifetime:  time.Duration(getEnvAsInt("JWT_ACCESS_LIFETIME", 15, false)) * time.Minute,
		JWTRefreshLifetime: time.Duration(getEnvAsInt("JWT_REFRESH_LIFETIME", 1, false)) * 24 * time.Hour,
		JWTSecretKey:       getEnv("JWT_SECRET_KEY", "", true),

		Debug: getEnvAsBool("DEBUG", false, false),

		PasswordSecretKey: getEnv("PASSWORD_SECRET_KEY", "", true),
		PasswordSalt:      getEnv("PASSWORD_SALT", "", true),
		PasswordCost:      getEnvAsInt("PASSWORD_COST", 16, true),

		EmailSecretKey:     getEnv("EMAIL_SECRET_KEY", "", true),
		EmailTokenLifetime: time.Duration(getEnvAsInt("EMAIL_TOKEN_LIFETIME", 5*60, false)) * time.Second,
		EmailFolderPath:    getEnv("EMAIL_FOLDER_PATH", "./sent_emails", false),
		EmailTimeout:       time.Duration(getEnvAsInt("EMAIL_TIMEOUT", 15, false)) * time.Second,
	}
}

func getEnv(key, defaultValue string, required bool) string {
	value := os.Getenv(key)
	if value == "" && required {
		panic(fmt.Sprintf("Missing environment variable: %s", key))
	}
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsUint16(key string, defaultValue uint16, required bool) uint16 {
	valueStr := getEnv(key, fmt.Sprintf("%v", defaultValue), required)
	value, err := strconv.ParseUint(valueStr, 10, 16)
	if err != nil {
		panic(err)
	}
	return uint16(value)
}

func getEnvAsInt(key string, defaultValue int, required bool) int {
	valueStr := getEnv(key, fmt.Sprintf("%v", defaultValue), required)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		panic(err)
	}
	return value
}

func getEnvAsBool(key string, defaultValue bool, required bool) bool {
	const (
		trueStr  = "true"
		falseStr = "false"
	)
	defaultValueStr := falseStr
	if defaultValue {
		defaultValueStr = trueStr
	}
	valueStr := getEnv(key, defaultValueStr, required)
	switch valueStr {
	case trueStr:
		return true
	case falseStr:
		return false
	default:
		panic(fmt.Sprintf("Invalid boolean value for %s: %s", key, valueStr))
	}
}
