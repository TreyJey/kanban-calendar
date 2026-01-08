package handlers

import (
    "time"
    "kanban-calendar/internal/repository"
    "github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, repo *repository.TaskRepository) {
    // CORS middleware
    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    })
    
    // Группа API маршрутов
    api := r.Group("/api")
    {
        // Задачи
        tasks := api.Group("/tasks")
        {
            tasks.GET("", GetTasks(repo))                    // GET /api/tasks
            tasks.GET("/:id", GetTaskByID(repo))            // GET /api/tasks/:id
            tasks.POST("", CreateTask(repo))                // POST /api/tasks
            tasks.PUT("/:id", UpdateTask(repo))             // PUT /api/tasks/:id
            tasks.DELETE("/:id", DeleteTask(repo))          // DELETE /api/tasks/:id
            tasks.GET("/status/:status", GetTasksByStatus(repo)) // GET /api/tasks/status/:status
        }
        
        // Календарь
        calendar := api.Group("/calendar")
        {
            calendar.GET("/events", GetCalendarEvents(repo)) // GET /api/calendar/events
        }
        
        // Системные
        api.GET("/health", func(c *gin.Context) {
            c.JSON(200, gin.H{
                "status":   "healthy",
                "service":  "kanban-calendar",
                "database": "connected",
                "time":     time.Now().Format(time.RFC3339),
            })
        })
        
        api.GET("/version", func(c *gin.Context) {
            c.JSON(200, gin.H{
                "version": "1.0.0",
                "name":    "Kanban Calendar API",
            })
        })
    }
    
    // Главная страница
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Kanban Calendar API",
            "version": "1.0.0",
            "docs":    "Доступные эндпоинты:",
            "endpoints": []gin.H{
                {"method": "GET",    "path": "/api/tasks",           "description": "Получить все задачи"},
                {"method": "GET",    "path": "/api/tasks/:id",       "description": "Получить задачу по ID"},
                {"method": "POST",   "path": "/api/tasks",           "description": "Создать новую задачу"},
                {"method": "PUT",    "path": "/api/tasks/:id",       "description": "Обновить задачу"},
                {"method": "DELETE", "path": "/api/tasks/:id",       "description": "Удалить задачу"},
                {"method": "GET",    "path": "/api/tasks/status/:status", "description": "Получить задачи по статусу"},
                {"method": "GET",    "path": "/api/calendar/events", "description": "Получить события календаря"},
                {"method": "GET",    "path": "/api/health",          "description": "Проверка здоровья сервиса"},
            },
        })
    })
}