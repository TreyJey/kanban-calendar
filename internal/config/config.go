package config

import "os"

type Config struct {
    DBHost         string
    DBPort         string
    DBUser         string
    DBPassword     string
    DBName         string
    ServerPort     string
    TelegramToken  string
    TelegramChatID string
    DBRetryDelay   string
}

func Load() *Config {
    return &Config{
        DBHost:         getEnv("DB_HOST", "kanban-postgres"),
        DBPort:         getEnv("DB_PORT", "5432"),
        DBUser:         getEnv("DB_USER", "postgres"),
        DBPassword:     getEnv("DB_PASSWORD", "password"),
        DBName:         getEnv("DB_NAME", "kanban"),
        ServerPort:     getEnv("SERVER_PORT", "8080"),
        TelegramToken:  getEnv("TELEGRAM_TOKEN", ""),
        TelegramChatID: getEnv("TELEGRAM_CHAT_ID", ""),
        DBRetryDelay:   getEnv("DB_RETRY_DELAY", "5"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}