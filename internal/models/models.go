package models

import (
    "time"
)

// TaskStatus - тип для статуса задачи
type TaskStatus string

// Енум статусов
const (
    StatusTodo       TaskStatus = "todo"
    StatusInProgress TaskStatus = "in_progress"
    StatusDone       TaskStatus = "done"
)

// Task - основная структура задачи
type Task struct {
    ID          int         `json:"id"`                    // ID задачи
    Title       string      `json:"title"`                 // Заголовок
    Description string      `json:"description,omitempty"` // Описание (может быть пустым)
    Status      TaskStatus  `json:"status"`                // Статус из констант выше
    Priority    string      `json:"priority,omitempty"`    // Приоритет: low, medium, high
    CreatedAt   time.Time   `json:"created_at"`           // Когда создана
    UpdatedAt   time.Time   `json:"updated_at"`           // Когда обновлена
    Deadline    *time.Time  `json:"deadline,omitempty"`   // Дедлайн (может быть nil)
    StartDate   *time.Time  `json:"start_date,omitempty"` // Дата начала (для календаря)
    EndDate     *time.Time  `json:"end_date,omitempty"`   // Дата окончания (для календаря)
    Assignee    string      `json:"assignee,omitempty"`   // Исполнитель
    Tags        []string    `json:"tags,omitempty"`       // Теги (массив строк)
    LastNotifiedHours int    `json:"last_notified_hours"`
}

// CalendarEvent - структура для отображения в календаре
type CalendarEvent struct {
    ID          int         `json:"id"`
    Title       string      `json:"title"`
    Description string      `json:"description"`
    Start       time.Time   `json:"start"`
    End         time.Time   `json:"end"`
    Status      TaskStatus  `json:"status"`
    Color       string      `json:"color,omitempty"` // Цвет события в календаре
}

// CreateTaskRequest - структура для запроса создания задачи
type CreateTaskRequest struct {
    Title       string     `json:"title" binding:"required"`
    Description string     `json:"description"`
    Status      TaskStatus `json:"status"`
    Priority    string     `json:"priority"`
    Deadline    string     `json:"deadline"`  // Будем парсить из строки
    StartDate   string     `json:"start_date"`
    EndDate     string     `json:"end_date"`
    Assignee    string     `json:"assignee"`
    Tags        []string   `json:"tags"`
}

// UpdateTaskRequest - структура для запроса обновления задачи
type UpdateTaskRequest struct {
    Title       string     `json:"title"`
    Description string     `json:"description"`
    Status      TaskStatus `json:"status"`
    Priority    string     `json:"priority"`
    Deadline    string     `json:"deadline"`
    StartDate   string     `json:"start_date"`
    EndDate     string     `json:"end_date"`
    Assignee    string     `json:"assignee"`
    Tags        []string   `json:"tags"`
}

// Метод для преобразования Task в CalendarEvent
func (t *Task) ToCalendarEvent() CalendarEvent {
    // Выбираем цвет в зависимости от статуса
    color := "#3174ad" // Синий по умолчанию
    if t.Status == StatusDone {
        color = "#28a745" // Зеленый для выполненных
    } else if t.Status == StatusInProgress {
        color = "#ffc107" // Желтый для в работе
    }
    
    // Даты начала и окончания
    start := time.Now()
    if t.StartDate != nil {
        start = *t.StartDate
    }
    
    end := start.Add(24 * time.Hour) // По умолчанию 1 день
    if t.EndDate != nil {
        end = *t.EndDate
    } else if t.Deadline != nil {
        end = *t.Deadline
    }
    
    return CalendarEvent{
        ID:          t.ID,
        Title:       t.Title,
        Description: t.Description,
        Start:       start,
        End:         end,
        Status:      t.Status,
        Color:       color,
    }
}