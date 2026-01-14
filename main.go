package main

import (
    "log"
    "os"
    "kanban-calendar/internal/config"
    "kanban-calendar/internal/database"
    "kanban-calendar/internal/handlers"
    "kanban-calendar/internal/repository"
    "kanban-calendar/scheduler"
    "kanban-calendar/telegram"
    "github.com/gin-gonic/gin"
)

func main() {
    // Загружаем конфигурацию
    cfg := config.Load()
    
    log.Printf("Конфигурация загружена")
    log.Printf("Порт сервера: %s", cfg.ServerPort)
    log.Printf("БД: %s@%s:%s/%s", 
        cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)
    
    if cfg.TelegramToken != "" {
        log.Println("Telegram токен найден")
    } else {
        log.Println("Telegram токен не указан, уведомления отключены")
    }
    
    // Подключаемся к БД
    db, err := database.Connect(cfg)
    if err != nil {
        log.Fatalf("Ошибка подключения к БД: %v", err)
    }
    defer db.Close()
    
    // Проверяем/создаем таблицы
    if err := database.Migrate(db); err != nil {
        log.Printf("Предупреждение миграций: %v", err)
    }

    _, _ = db.Exec(`ALTER TABLE tasks DROP CONSTRAINT IF EXISTS tasks_external_uid_key;`)
    _, _ = db.Exec(`DROP INDEX IF EXISTS tasks_external_uid_key;`)
    
    // Создаем репозиторий
    repo := repository.NewTaskRepository(db)
    
    // Инициализируем Telegram бота (если токен указан)
    var telegramBot *telegram.TelegramBot
    if cfg.TelegramToken != "" && cfg.TelegramChatID != "" {
        // Получаем URL фронтенда из окружения (или ставим дефолт)
        frontendURL := os.Getenv("FRONTEND_URL")
        if frontendURL == "" {
            frontendURL = "http://localhost:3000" // Дефолт для разработки
        }

        telegramBot, err = telegram.NewTelegramBot(cfg.TelegramToken, cfg.TelegramChatID, frontendURL)
        if err != nil {
            log.Printf("Telegram бот не запущен: %v", err)
        } else {
            log.Println("Telegram бот инициализирован")
            
            // Запускаем планировщик уведомлений
            sched := scheduler.NewScheduler(repo, telegramBot)
            sched.Start()
            telegramBot.SendTestMessage()
            log.Println("Планировщик уведомлений запущен")
        }
    }
    
    // Настраиваем Gin
    if cfg.ServerPort == "8080" {
        gin.SetMode(gin.ReleaseMode)
    }
    r := gin.Default()
    
    // Настраиваем маршруты
    handlers.SetupRoutes(r, repo)
    
    // Запуск сервера
    log.Printf("Сервер запущен на http://localhost:%s", cfg.ServerPort)
    log.Println("Документация API доступна по адресу http://localhost:" + cfg.ServerPort)
    
    if err := r.Run(":" + cfg.ServerPort); err != nil {
        log.Fatal("Ошибка запуска сервера:", err)
    }
}