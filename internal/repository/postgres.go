package repository

import (
    "context"
    "database/sql"
    "fmt"
    "time"
    "kanban-calendar/internal/models"
)

// TaskRepository - репозиторий для работы с задачами
type TaskRepository struct {
    db *sql.DB
}

// NewTaskRepository - конструктор
func NewTaskRepository(db *sql.DB) *TaskRepository {
    return &TaskRepository{db: db}
}

func (r *TaskRepository) CreateTask(ctx context.Context, task *models.Task) error {
    query := `
        INSERT INTO tasks 
        (title, description, status, priority, deadline, start_date, end_date, assignee, last_notified_hours)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 100)
        RETURNING id, created_at, updated_at, last_notified_hours
    `
    
    return r.db.QueryRowContext(ctx, query,
        task.Title,
        task.Description,
        task.Status,
        task.Priority,
        task.Deadline,
        task.StartDate,
        task.EndDate,
        task.Assignee,
    ).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt, &task.LastNotifiedHours)
}

// GetTaskByID - получает задачу по ID (БЕЗ TAGS)
func (r *TaskRepository) GetTaskByID(ctx context.Context, id int) (*models.Task, error) {
    query := `
        SELECT id, title, description, status, priority, created_at, updated_at,
               deadline, start_date, end_date, assignee
        FROM tasks WHERE id = $1
    `
    
    task := &models.Task{}
    var deadline, startDate, endDate sql.NullTime
    
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &task.ID, &task.Title, &task.Description, &task.Status, &task.Priority,
        &task.CreatedAt, &task.UpdatedAt, &deadline, &startDate, &endDate,
        &task.Assignee,
        // УБРАЛИ: pq.Array(&tags),
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("задача с ID %d не найдена", id)
        }
        return nil, err
    }
    
    // Преобразуем NullTime в *time.Time
    if deadline.Valid {
        task.Deadline = &deadline.Time
    }
    if startDate.Valid {
        task.StartDate = &startDate.Time
    }
    if endDate.Valid {
        task.EndDate = &endDate.Time
    }
    // УБРАЛИ: task.Tags = tags
    
    return task, nil
}

// Получение всех задач (обновлено: добавлено поле last_notified_hours)
func (r *TaskRepository) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	query := `
		SELECT id, title, description, status, priority, created_at, updated_at, 
		       deadline, start_date, end_date, assignee, last_notified_hours
		FROM tasks
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var deadline, startDate, endDate sql.NullTime
		
		err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority,
			&task.CreatedAt, &task.UpdatedAt, &deadline, &startDate, &endDate,
			&task.Assignee, &task.LastNotifiedHours, // Сканируем новое поле
		)
		if err != nil {
			return nil, err
		}
		
		if deadline.Valid { task.Deadline = &deadline.Time }
		if startDate.Valid { task.StartDate = &startDate.Time }
		if endDate.Valid { task.EndDate = &endDate.Time }
		
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// UpdateTask - обновляет задачу (БЕЗ TAGS)
func (r *TaskRepository) UpdateTask(ctx context.Context, task *models.Task) error {
    query := `
        UPDATE tasks 
        SET title = $1, description = $2, status = $3, priority = $4,
            deadline = $5, start_date = $6, end_date = $7, 
            assignee = $8, updated_at = CURRENT_TIMESTAMP
        WHERE id = $9
        RETURNING updated_at
    `
    
    return r.db.QueryRowContext(ctx, query,
        task.Title,
        task.Description,
        task.Status,
        task.Priority,
        task.Deadline,
        task.StartDate,
        task.EndDate,
        task.Assignee,
        // УБРАЛИ: pq.Array(task.Tags),
        task.ID,
    ).Scan(&task.UpdatedAt)
}

// DeleteTask - удаляет задачу
func (r *TaskRepository) DeleteTask(ctx context.Context, id int) error {
    query := `DELETE FROM tasks WHERE id = $1`
    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return err
    }
    
    rows, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rows == 0 {
        return fmt.Errorf("задача с ID %d не найдена", id)
    }
    
    return nil
}

// GetTasksByStatus - получает задачи по статусу (БЕЗ TAGS)
func (r *TaskRepository) GetTasksByStatus(ctx context.Context, status models.TaskStatus) ([]models.Task, error) {
    query := `
        SELECT id, title, description, status, priority, created_at, updated_at,
               deadline, start_date, end_date, assignee
        FROM tasks 
        WHERE status = $1
        ORDER BY created_at DESC
    `
    
    rows, err := r.db.QueryContext(ctx, query, status)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var tasks []models.Task
    for rows.Next() {
        var task models.Task
        var deadline, startDate, endDate sql.NullTime
        
        err := rows.Scan(
            &task.ID, &task.Title, &task.Description, &task.Status, &task.Priority,
            &task.CreatedAt, &task.UpdatedAt, &deadline, &startDate, &endDate,
            &task.Assignee,
            // УБРАЛИ: pq.Array(&tags),
        )
        if err != nil {
            return nil, err
        }
        
        if deadline.Valid {
            task.Deadline = &deadline.Time
        }
        if startDate.Valid {
            task.StartDate = &startDate.Time
        }
        if endDate.Valid {
            task.EndDate = &endDate.Time
        }
        // УБРАЛИ: task.Tags = tags
        
        tasks = append(tasks, task)
    }
    
    return tasks, nil
}

// GetCalendarEvents - получает события для календаря
func (r *TaskRepository) GetCalendarEvents(ctx context.Context, startDate, endDate time.Time) ([]models.CalendarEvent, error) {
    query := `
        SELECT id, title, description, status, 
               COALESCE(start_date, created_at) as start,
               COALESCE(end_date, deadline, created_at + INTERVAL '1 day') as end
        FROM tasks 
        WHERE (start_date BETWEEN $1 AND $2) 
           OR (end_date BETWEEN $1 AND $2)
           OR (deadline BETWEEN $1 AND $2)
           OR (created_at BETWEEN $1 AND $2)
        ORDER BY start
    `
    
    rows, err := r.db.QueryContext(ctx, query, startDate, endDate)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var events []models.CalendarEvent
    for rows.Next() {
        var event models.CalendarEvent
        var start, end time.Time
        
        err := rows.Scan(
            &event.ID,
            &event.Title,
            &event.Description,
            &event.Status,
            &start,
            &end,
        )
        if err != nil {
            return nil, err
        }
        
        event.Start = start
        event.End = end
        
        // Устанавливаем цвет по статусу
        if event.Status == models.StatusDone {
            event.Color = "#28a745"
        } else if event.Status == models.StatusInProgress {
            event.Color = "#ffc107"
        } else {
            event.Color = "#3174ad"
        }
        
        events = append(events, event)
    }
    
    return events, nil
}
// GetUpcomingDeadlines - получает задачи с приближающимися дедлайнами
func (r *TaskRepository) GetUpcomingDeadlines(ctx context.Context, hoursBefore int) ([]models.Task, error) {
    query := `
        SELECT id, title, description, status, priority, created_at, updated_at,
               deadline, start_date, end_date, assignee
        FROM tasks 
        WHERE deadline IS NOT NULL 
          AND deadline > NOW()
          AND deadline <= NOW() + INTERVAL '1 hour' * $1
          AND status != $2
        ORDER BY deadline ASC
    `
    
    rows, err := r.db.QueryContext(ctx, query, hoursBefore, models.StatusDone)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var tasks []models.Task
    for rows.Next() {
        var task models.Task
        var deadline, startDate, endDate sql.NullTime
        
        err := rows.Scan(
            &task.ID, &task.Title, &task.Description, &task.Status, &task.Priority,
            &task.CreatedAt, &task.UpdatedAt, &deadline, &startDate, &endDate,
            &task.Assignee,
        )
        if err != nil {
            return nil, err
        }
        
        if deadline.Valid {
            task.Deadline = &deadline.Time
        }
        if startDate.Valid {
            task.StartDate = &startDate.Time
        }
        if endDate.Valid {
            task.EndDate = &endDate.Time
        }
        
        tasks = append(tasks, task)
    }
    
    return tasks, nil
}

// GetOverdueTasks - получает просроченные задачи
func (r *TaskRepository) GetOverdueTasks(ctx context.Context) ([]models.Task, error) {
    query := `
        SELECT id, title, description, status, priority, created_at, updated_at,
               deadline, start_date, end_date, assignee
        FROM tasks 
        WHERE deadline IS NOT NULL 
          AND deadline < NOW()
          AND status != $1
        ORDER BY deadline ASC
    `
    
    rows, err := r.db.QueryContext(ctx, query, models.StatusDone)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var tasks []models.Task
    for rows.Next() {
        var task models.Task
        var deadline, startDate, endDate sql.NullTime
        
        err := rows.Scan(
            &task.ID, &task.Title, &task.Description, &task.Status, &task.Priority,
            &task.CreatedAt, &task.UpdatedAt, &deadline, &startDate, &endDate,
            &task.Assignee,
        )
        if err != nil {
            return nil, err
        }
        
        if deadline.Valid {
            task.Deadline = &deadline.Time
        }
        if startDate.Valid {
            task.StartDate = &startDate.Time
        }
        if endDate.Valid {
            task.EndDate = &endDate.Time
        }
        
        tasks = append(tasks, task)
    }
    
    return tasks, nil
}

// GetTasksCompletedToday - получает задачи, выполненные сегодня
func (r *TaskRepository) GetTasksCompletedToday(ctx context.Context) ([]models.Task, error) {
    query := `
        SELECT id, title, description, status, priority, created_at, updated_at,
               deadline, start_date, end_date, assignee
        FROM tasks 
        WHERE status = $1 
          AND DATE(updated_at) = CURRENT_DATE
        ORDER BY updated_at DESC
    `
    
    rows, err := r.db.QueryContext(ctx, query, models.StatusDone)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var tasks []models.Task
    for rows.Next() {
        var task models.Task
        var deadline, startDate, endDate sql.NullTime
        
        err := rows.Scan(
            &task.ID, &task.Title, &task.Description, &task.Status, &task.Priority,
            &task.CreatedAt, &task.UpdatedAt, &deadline, &startDate, &endDate,
            &task.Assignee,
        )
        if err != nil {
            return nil, err
        }
        
        if deadline.Valid {
            task.Deadline = &deadline.Time
        }
        if startDate.Valid {
            task.StartDate = &startDate.Time
        }
        if endDate.Valid {
            task.EndDate = &endDate.Time
        }
        
        tasks = append(tasks, task)
    }
    
    return tasks, nil
}

// CreateNotification - создает запись об уведомлении
func (r *TaskRepository) CreateNotification(ctx context.Context, notification *models.Notification) error {
    query := `
        INSERT INTO notifications 
        (task_id, type, message, sent_at, is_sent, chat_id)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `
    
    return r.db.QueryRowContext(ctx, query,
        notification.TaskID,
        notification.Type,
        notification.Message,
        notification.SentAt,
        notification.IsSent,
        notification.ChatID,
    ).Scan(&notification.ID)
}

// MarkNotificationSent - помечает уведомление как отправленное
func (r *TaskRepository) MarkNotificationSent(ctx context.Context, id int) error {
    query := `UPDATE notifications SET is_sent = true WHERE id = $1`
    _, err := r.db.ExecContext(ctx, query, id)
    return err
}

func (r *TaskRepository) UpdateLastNotified(ctx context.Context, taskID int, hours int) error {
    query := `UPDATE tasks SET last_notified_hours = $1, updated_at = NOW() WHERE id = $2`
    _, err := r.db.ExecContext(ctx, query, hours, taskID)
    return err
}