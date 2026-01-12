package handlers

import (
    "github.com/arran4/golang-ical"
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

func ImportCalendar(repo *repository.TaskRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Читаем файл из формы (ключ должен быть "calendar")
		fileHeader, err := c.FormFile("calendar")
		if err != nil {
			c.JSON(400, gin.H{"error": "Файл не найден в запросе"})
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(500, gin.H{"error": "Не удалось открыть файл"})
			return
		}
		defer file.Close()

		// 2. Парсим содержимое .ics
		cal, err := ics.ParseCalendar(file)
		if err != nil {
			c.JSON(400, gin.H{"error": "Ошибка формата файла .ics"})
			return
		}

		imported := 0
		skipped := 0

		// 3. Проходим по всем событиям в файле
		for _, event := range cal.Events() {
			uid := event.Id()
			summary := ""
			if prop := event.GetProperty(ics.ComponentPropertySummary); prop != nil {
				summary = prop.Value
			}

			description := ""
			if prop := event.GetProperty(ics.ComponentPropertyDescription); prop != nil {
				description = prop.Value
			}

			// Простейший парсинг дат (формат 20260102T150405Z)
			var start, end *time.Time
			if prop := event.GetProperty(ics.ComponentPropertyDtStart); prop != nil {
				t, _ := time.Parse("20060102T150405Z", prop.Value)
				if !t.IsZero() {
					start = &t
				}
			}
			if prop := event.GetProperty(ics.ComponentPropertyDtEnd); prop != nil {
				t, _ := time.Parse("20060102T150405Z", prop.Value)
				if !t.IsZero() {
					end = &t
				}
			}

			// Создаем объект задачи для базы
			task := &models.Task{
				ExternalUID:       uid,
				Title:             summary,
				Description:       description,
				Status:            models.StatusTodo,
				StartDate:         start,
				EndDate:           end,
				Deadline:          end,
				LastNotifiedHours: 100, // Чтобы бот начал отсчет заново
			}

			// 4. Пробуем сохранить в базу
			if err := repo.CreateTask(c.Request.Context(), task); err != nil {
				// Если ошибка (например, такой UID уже есть), пропускаем
				skipped++
				continue
			}
			imported++
		}

		c.JSON(200, gin.H{
			"status":   "success",
			"imported": imported,
			"skipped":  skipped,
		})
	}
}