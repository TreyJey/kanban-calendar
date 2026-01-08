package handlers

import (
    "net/http"
    "time"
    "kanban-calendar/internal/models"
    "kanban-calendar/internal/repository"
    "github.com/gin-gonic/gin"
)

// GetCalendarEvents - получает события для календаря
func GetCalendarEvents(repo *repository.TaskRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Получаем параметры дат из запроса
        startStr := c.DefaultQuery("start", time.Now().AddDate(0, -1, 0).Format(time.RFC3339))
        endStr := c.DefaultQuery("end", time.Now().AddDate(0, 1, 0).Format(time.RFC3339))
        
        // Парс дат
        start, err := time.Parse(time.RFC3339, startStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error":   "Неверный формат даты начала",
                "details": err.Error(),
                "example": "2024-01-01T00:00:00Z",
            })
            return
        }
        
        end, err := time.Parse(time.RFC3339, endStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error":   "Неверный формат даты окончания",
                "details": err.Error(),
                "example": "2024-12-31T23:59:59Z",
            })
            return
        }
        
        // Проверяем, что даты валидны
        if end.Before(start) {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "Дата окончания должна быть после даты начала",
            })
            return
        }
        
        // Получаем события из БД
        events, err := repo.GetCalendarEvents(c.Request.Context(), start, end)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error":   "Ошибка получения событий календаря",
                "details": err.Error(),
            })
            return
        }
        
        if events == nil {
            events = []models.CalendarEvent{}
        }
        
        c.JSON(http.StatusOK, gin.H{
            "events": events,
            "count":  len(events),
            "start":  start.Format(time.RFC3339),
            "end":    end.Format(time.RFC3339),
        })
    }
}