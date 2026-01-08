package handlers

import (
    "net/http"
    "strconv"
    "time"
    "kanban-calendar/internal/models"
    "kanban-calendar/internal/repository"
    "github.com/gin-gonic/gin"
)

// GetTasks - получает все задачи
func GetTasks(repo *repository.TaskRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        tasks, err := repo.GetAllTasks(c.Request.Context())
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error":   "Ошибка получения задач",
                "details": err.Error(),
            })
            return
        }
        
        if tasks == nil {
            tasks = []models.Task{} // Пустой массив вместо nil
        }
        
        c.JSON(http.StatusOK, gin.H{
            "tasks": tasks,
            "count": len(tasks),
        })
    }
}

// GetTaskByID - получает задачу по ID
func GetTaskByID(repo *repository.TaskRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        id, err := strconv.Atoi(c.Param("id"))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "Неверный формат ID задачи",
            })
            return
        }
        
        task, err := repo.GetTaskByID(c.Request.Context(), id)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{
                "error":   "Задача не найдена",
                "details": err.Error(),
            })
            return
        }
        
        c.JSON(http.StatusOK, task)
    }
}

// CreateTask - создает новую задачу
func CreateTask(repo *repository.TaskRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req models.CreateTaskRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error":   "Неверный формат данных",
                "details": err.Error(),
            })
            return
        }
        
        // Создаем задачу из запроса
        task := &models.Task{
            Title:       req.Title,
            Description: req.Description,
            Status:      req.Status,
            Priority:    req.Priority,
            Assignee:    req.Assignee,
            Tags:        req.Tags,
        }
        
        // Парсим даты если они переданы
        parseTime := func(timeStr string) (*time.Time, error) {
            if timeStr == "" {
                return nil, nil
            }
            t, err := time.Parse(time.RFC3339, timeStr)
            if err != nil {
                return nil, err
            }
            return &t, nil
        }
        
        if deadline, err := parseTime(req.Deadline); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error":   "Неверный формат даты дедлайна",
                "details": err.Error(),
            })
            return
        } else {
            task.Deadline = deadline
        }
        
        if startDate, err := parseTime(req.StartDate); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error":   "Неверный формат даты начала",
                "details": err.Error(),
            })
            return
        } else {
            task.StartDate = startDate
        }
        
        if endDate, err := parseTime(req.EndDate); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error":   "Неверный формат даты окончания",
                "details": err.Error(),
            })
            return
        } else {
            task.EndDate = endDate
        }
        
        // Создаем задачу в БД
        if err := repo.CreateTask(c.Request.Context(), task); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error":   "Ошибка создания задачи",
                "details": err.Error(),
            })
            return
        }
        
        c.JSON(http.StatusCreated, task)
    }
}

// UpdateTask - обновляет задачу
func UpdateTask(repo *repository.TaskRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        id, err := strconv.Atoi(c.Param("id"))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "Неверный формат ID задачи",
            })
            return
        }
        
        // Получаем существующую задачу
        task, err := repo.GetTaskByID(c.Request.Context(), id)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{
                "error":   "Задача не найдена",
                "details": err.Error(),
            })
            return
        }
        
        // Парсим обновления
        var req models.UpdateTaskRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error":   "Неверный формат данных",
                "details": err.Error(),
            })
            return
        }
        
        // Обновляем поля если они переданы
        if req.Title != "" {
            task.Title = req.Title
        }
        if req.Description != "" {
            task.Description = req.Description
        }
        if req.Status != "" {
            task.Status = req.Status
        }
        if req.Priority != "" {
            task.Priority = req.Priority
        }
        if req.Assignee != "" {
            task.Assignee = req.Assignee
        }
        if req.Tags != nil {
            task.Tags = req.Tags
        }
        
        // Парсим даты если они переданы
        parseTime := func(timeStr string) (*time.Time, error) {
            if timeStr == "" {
                return nil, nil
            }
            t, err := time.Parse(time.RFC3339, timeStr)
            if err != nil {
                return nil, err
            }
            return &t, nil
        }
        
        if req.Deadline != "" {
            if deadline, err := parseTime(req.Deadline); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{
                    "error":   "Неверный формат даты дедлайна",
                    "details": err.Error(),
                })
                return
            } else {
                task.Deadline = deadline
            }
        }
        
        if req.StartDate != "" {
            if startDate, err := parseTime(req.StartDate); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{
                    "error":   "Неверный формат даты начала",
                    "details": err.Error(),
                })
                return
            } else {
                task.StartDate = startDate
            }
        }
        
        if req.EndDate != "" {
            if endDate, err := parseTime(req.EndDate); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{
                    "error":   "Неверный формат даты окончания",
                    "details": err.Error(),
                })
                return
            } else {
                task.EndDate = endDate
            }
        }
        
        // Сохраняем изменения
        if err := repo.UpdateTask(c.Request.Context(), task); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error":   "Ошибка обновления задачи",
                "details": err.Error(),
            })
            return
        }
        
        c.JSON(http.StatusOK, task)
    }
}

// DeleteTask - удаляет задачу
func DeleteTask(repo *repository.TaskRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        id, err := strconv.Atoi(c.Param("id"))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "Неверный формат ID задачи",
            })
            return
        }
        
        if err := repo.DeleteTask(c.Request.Context(), id); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error":   "Ошибка удаления задачи",
                "details": err.Error(),
            })
            return
        }
        
        c.JSON(http.StatusOK, gin.H{
            "message": "Задача успешно удалена",
            "id":      id,
        })
    }
}

// GetTasksByStatus - получает задачи по статусу
func GetTasksByStatus(repo *repository.TaskRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        status := models.TaskStatus(c.Param("status"))
        
        // Проверяем валидность статуса
        validStatuses := map[models.TaskStatus]bool{
            models.StatusTodo:       true,
            models.StatusInProgress: true,
            models.StatusDone:       true,
        }
        
        if !validStatuses[status] {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "Неверный статус задачи. Допустимые значения: todo, in_progress, done",
            })
            return
        }
        
        tasks, err := repo.GetTasksByStatus(c.Request.Context(), status)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error":   "Ошибка получения задач",
                "details": err.Error(),
            })
            return
        }
        
        if tasks == nil {
            tasks = []models.Task{}
        }
        
        c.JSON(http.StatusOK, gin.H{
            "tasks": tasks,
            "count": len(tasks),
            "status": status,
        })
    }
}