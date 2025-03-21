package config

import (
	"os"
)

// Config armazena todas as configurações da aplicação
type Config struct {
	DatabaseURL       string
	BacenUsername     string
	BacenPassword     string
	JWTSecret         string
	ServerPort        string
	TokenExpiryHours  int
}

// NewConfig cria uma nova instância de configuração
func NewConfig() *Config {
	return &Config{
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		BacenUsername:     os.Getenv("usernameBC"),
		BacenPassword:     os.Getenv("passwordBC"),
		JWTSecret:         getEnvOrDefault("JWT_SECRET", "zH4NRP1HMALxxCFnRZABFA7GOJtzU_gIj02alfL1lvI"),
		ServerPort:        getEnvOrDefault("PORT", "8080"),
		TokenExpiryHours:  24, // Token válido por 24 horas
	}
}

// getEnvOrDefault retorna o valor da variável de ambiente ou o valor padrão
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}