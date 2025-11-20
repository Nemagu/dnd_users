package config

import (
	"fmt"
	"os"
	"strconv"
)

type WebConfig struct {
	DBHost            string
	DBPort            uint16
	DBUser            string
	DBPassword        string
	DBName            string
	HTTPPort          uint16
	Debug             bool
	PasswordSecretKey string
	PasswordSalt      string
	EmailSecretKey    string
	EmailSalt         string
}

func MustNewWebConfig() *WebConfig {
	return &WebConfig{
		DBHost:            getEnv("DB_HOST", "localhost", true),
		DBPort:            getEnvAsUint16("DB_PORT", 5432, true),
		DBUser:            getEnv("DB_USER", "postgres", true),
		DBPassword:        getEnv("DB_PASSWORD", "postgres", true),
		DBName:            getEnv("DB_NAME", "postgres", true),
		HTTPPort:          getEnvAsUint16("HTTP_PORT", 8080, true),
		Debug:             getEnvAsBool("DEBUG", false, false),
		PasswordSecretKey: getEnv("PASSWORD_SECRET_KEY", "", true),
		PasswordSalt:      getEnv("PASSWORD_SALT", "", true),
		EmailSecretKey:    getEnv("EMAIL_SECRET_KEY", "", true),
		EmailSalt:         getEnv("EMAIL_SALT", "", true),
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
