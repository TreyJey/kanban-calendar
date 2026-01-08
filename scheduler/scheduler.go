package scheduler

import (
	"context"
	"log"
	"time"
	"kanban-calendar/internal/repository"
	"kanban-calendar/telegram"
)

type Scheduler struct {
	repo     *repository.TaskRepository
	telegram *telegram.TelegramBot
	interval time.Duration
}

func NewScheduler(repo *repository.TaskRepository, telegram *telegram.TelegramBot) *Scheduler {
	return &Scheduler{
		repo:     repo,
		telegram: telegram,
		interval: 5 * time.Minute, // Оптимально проверять раз в 5 минут
	}
}

func (s *Scheduler) Start() {
	log.Println("🚀 Планировщик запущен с интервалами: 48, 24, 12, 6, 3, 0 ч.")
	go s.runDeadlineChecker()
}

func (s *Scheduler) runDeadlineChecker() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.checkDeadlines() // Первая проверка при запуске

	for range ticker.C {
		s.checkDeadlines()
	}
}

func (s *Scheduler) checkDeadlines() {
	ctx := context.Background()
	
	// Определяем пороги уведомлений (в часах)
	thresholds := []int{48, 24, 12, 6, 3, 0}

	tasks, err := s.repo.GetAllTasks(ctx)
	if err != nil {
		log.Printf("❌ Ошибка получения задач: %v", err)
		return
	}

	for _, task := range tasks {
		// Пропускаем выполненные или задачи без дедлайна
		if task.Deadline == nil || task.Status == "done" {
			continue
		}

		// Считаем сколько часов осталось до дедлайна
		hoursLeft := int(time.Until(*task.Deadline).Hours())

		// Проверяем пороги по порядку
		for _, t := range thresholds {
			// Если время пришло (осталось <= порога) 
			// И мы еще не уведомляли именно об этом пороге (LastNotifiedHours > t)
			if hoursLeft <= t && task.LastNotifiedHours > t {
				
				log.Printf("🔔 Отправка уведомления: задача '%s' (порог %d ч.)", task.Title, t)
				
				err := s.telegram.SendDeadlineNotification(task, hoursLeft)
				if err != nil {
					log.Printf("❌ Ошибка отправки в ТГ: %v", err)
					continue
				}

				// Запоминаем в базе, что этот порог пройден
				err = s.repo.UpdateLastNotified(ctx, task.ID, t)
				if err != nil {
					log.Printf("❌ Ошибка обновления порога в БД: %v", err)
				}
				
				break // Для одной задачи отправляем только одно уведомление за раз
			}
		}
	}
}