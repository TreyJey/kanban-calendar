package database

import (
    "database/sql"
    "fmt"
    "log"
    "time"
    "kanban-calendar/internal/config"
    _ "github.com/lib/pq" // Драйвер PostgreSQL
)

// Connect - подключается к PostgreSQL с повторными попытками
func Connect(cfg *config.Config) (*sql.DB, error) {
    // Строка подключения
    connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
    )
    
    log.Printf("Подключение к БД: %s@%s:%s", 
        cfg.DBUser, cfg.DBHost, cfg.DBPort)
    
    var db *sql.DB
    var err error
    
    // Пытаемся подключиться несколько раз
    maxAttempts := 10
    for attempt := 1; attempt <= maxAttempts; attempt++ {
        // Открываем соединение
        db, err = sql.Open("postgres", connStr)
        if err != nil {
            log.Printf("Попытка %d/%d: ошибка подключения к БД: %v", 
                attempt, maxAttempts, err)
            time.Sleep(5 * time.Second)
            continue
        }
        
        // Проверяем подключение
        if err := db.Ping(); err != nil {
            log.Printf("Попытка %d/%d: ошибка ping БД: %v", 
                attempt, maxAttempts, err)
            db.Close()
            time.Sleep(5 * time.Second)
            continue
        }
        
        log.Printf("Подключение к БД установлено (попытка %d)", attempt)
        
        // Настройки пула соединений
        db.SetMaxOpenConns(25)      // Максимум открытых соединений
        db.SetMaxIdleConns(5)       // Максимум неактивных соединений
        db.SetConnMaxLifetime(5 * time.Minute) // Время жизни соединения
        
        return db, nil
    }
    
    return nil, fmt.Errorf("не удалось подключиться к БД после %d попыток: %w", maxAttempts, err)
}

// Migrate - выполняет миграции
func Migrate(db *sql.DB) error {
    // Пока просто проверяем таблицу
    query := `
        SELECT EXISTS (
            SELECT FROM information_schema.tables 
            WHERE table_name = 'tasks'
        )
    `
    
    var exists bool
    err := db.QueryRow(query).Scan(&exists)
    if err != nil {
        return fmt.Errorf("ошибка проверки таблицы: %w", err)
    }
    
    if !exists {
        log.Println("Таблица 'tasks' не существует. Нужно выполнить миграции.")
    } else {
        log.Println("Таблица 'tasks' существует")
    }
    
    return nil
}