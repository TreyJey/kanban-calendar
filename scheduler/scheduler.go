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
	telegram *telegram.TelegramBot // Исправлено: TelegramBot вместо Bot
}

func NewScheduler(repo *repository.TaskRepository, tg *telegram.TelegramBot) *Scheduler {
	return &Scheduler{
		repo:     repo,
		telegram: tg,
	}
}

func (s *Scheduler) Start() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			s.CheckDeadlines()
		}
	}()
}

func (s *Scheduler) CheckDeadlines() {
	// ВНИМАНИЕ: Если здесь будет ошибка "undefined: s.repo.GetTasks",
	// проверь postgres.go и замени GetTasks на GetAllTasks (или как он там называется)
	tasks, err := s.repo.GetAllTasks(context.Background())
	if err != nil {
		log.Printf("❌ Ошибка получения задач: %v", err)
		return
	}

	// Пороги уведомлений: 48ч, 24ч, 12ч, 6ч, 3ч, 0ч
	thresholds := []int{0, 3, 6, 12, 24, 48}

	for _, task := range tasks {
		// Пропускаем выполненные задачи и задачи без дедлайна
		if task.Status == "done" || task.Deadline == nil || task.Deadline.IsZero() {
			continue
		}

		hoursLeft := time.Until(*task.Deadline).Hours()

		for _, t := range thresholds {
			// Логика: если осталось меньше порога T и мы еще не уведомляли об этом пороге
			if hoursLeft <= float64(t) && task.LastNotifiedHours > t {
				
				err := s.telegram.SendDeadlineNotification(task, int(hoursLeft))
				if err != nil {
					log.Printf("❌ Ошибка отправки в TG: %v", err)
					break 
				}

				// Обновляем состояние в базе, чтобы не было дублей
				err = s.repo.UpdateLastNotified(context.Background(), task.ID, t)
				if err != nil {
					log.Printf("❌ Ошибка обновления порога в БД: %v", err)
				} else {
					log.Printf("✅ Уведомление отправлено для '%s' (порог %d ч.)", task.Title, t)
				}
				
				break // Выходим из цикла порогов для этой задачи
			}
		}
	}
}