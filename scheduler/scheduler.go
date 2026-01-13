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
	tasks, err := s.repo.GetAllTasks(context.Background())
	if err != nil {
		log.Printf("Ошибка получения задач: %v", err)
		return
	}

	// Идем ОТ БОЛЬШЕГО К МЕНЬШЕМУ
	thresholds := []int{48, 24, 12, 6, 3, 0}
	
	loc := time.FixedZone("UTC+5", 5*60*60)
	now := time.Now().In(loc)

	for _, task := range tasks {
		if task.Status == "done" || task.Deadline == nil || task.Deadline.IsZero() {
			continue
		}

		deadline := task.Deadline.In(loc)
		timeLeft := deadline.Sub(now)
		hoursLeft := timeLeft.Hours()

		// Если задача просрочена более чем на 1 час, перестаем спамить
		if hoursLeft < -1.0 {
			continue
		}

		for i, t := range thresholds {
			isCurrentThreshold := hoursLeft <= float64(t)
			if i+1 < len(thresholds) {
				if hoursLeft <= float64(thresholds[i+1]) {
					continue
				}
			}

			// Если мы в этом пороге И мы о нем еще не уведомляли
			if isCurrentThreshold && task.LastNotifiedHours > t {
				err := s.telegram.SendDeadlineNotification(task, int(hoursLeft))
				if err != nil {
					log.Printf("Ошибка отправки в TG: %v", err)
					break 
				}

				// Фиксируем в базе
				err = s.repo.UpdateLastNotified(context.Background(), task.ID, t)
				if err != nil {
					log.Printf("Ошибка обновления порога в БД: %v", err)
				}
				
				// ПРЕРЫВАЕМ цикл порогов
				break 
			}
		}
	}
}