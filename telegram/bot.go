package telegram

import (
    "fmt"
    "log"
    "time"
    "kanban-calendar/internal/models"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
    bot    *tgbotapi.BotAPI
    ChatID string
}

func NewTelegramBot(token, chatID string) (*TelegramBot, error) {
    if token == "" {
        return nil, fmt.Errorf("токен Telegram не указан")
    }
    
    bot, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        return nil, fmt.Errorf("ошибка создания бота: %w", err)
    }
    
    bot.Debug = false
    log.Printf("Telegram бот авторизован как %s", bot.Self.UserName)
    
    return &TelegramBot{
        bot:    bot,
        ChatID: chatID,
    }, nil
}

// SendDeadlineNotification - отправляет уведомление о дедлайне
func (tb *TelegramBot) SendDeadlineNotification(task models.Task, hoursLeft int) error {
    var message string
    
    if hoursLeft <= 0 {
        // Просроченная задача
        overdue := time.Since(*task.Deadline)
        hours := int(overdue.Hours())
        
        message = fmt.Sprintf(
            "🚨 *ПРОСРОЧЕНА!* 🚨\n" +
            "*Задача:* %s\n" +
            "*Просрочено:* %d час(ов) назад\n" +
            "*Исполнитель:* %s\n" +
            "*Статус:* %s\n" +
            "*Приоритет:* %s\n\n" +
            "Ссылка: http://localhost:8080/tasks/%d",
            task.Title,
            hours,
            task.Assignee,
            task.Status,
            task.Priority,
            task.ID,
        )
    } else if hoursLeft <= 24 {
        // Дедлайн в течение 24 часов
        message = fmt.Sprintf(
            "⏰ *Скоро дедлайн!*\n" +
            "*Задача:* %s\n" +
            "*Осталось:* %d час(ов)\n" +
            "*Дедлайн:* %s\n" +
            "*Исполнитель:* %s\n" +
            "*Статус:* %s\n\n" +
            "Ссылка: http://localhost:8080/tasks/%d",
            task.Title,
            hoursLeft,
            task.Deadline.Format("02.01.2006 15:04"),
            task.Assignee,
            task.Status,
            task.ID,
        )
    } else {
        // Дедлайн через несколько дней
        daysLeft := hoursLeft / 24
        message = fmt.Sprintf(
            "📅 *Напоминание о дедлайне*\n" +
            "*Задача:* %s\n" +
            "*Осталось:* %d дней\n" +
            "*Дедлайн:* %s\n" +
            "*Исполнитель:* %s\n\n" +
            "Ссылка: http://localhost:8080/tasks/%d",
            task.Title,
            daysLeft,
            task.Deadline.Format("02.01.2006"),
            task.Assignee,
            task.ID,
        )
    }
    
    msg := tgbotapi.NewMessageToChannel(tb.ChatID, message)
    msg.ParseMode = "Markdown"
    
    _, err := tb.bot.Send(msg)
    return err
}

// SendStatusChangeNotification - отправляет уведомление об изменении статуса
func (tb *TelegramBot) SendStatusChangeNotification(task models.Task, oldStatus models.TaskStatus) error {
    message := fmt.Sprintf(
        "🔄 *Статус изменен*\n" +
        "*Задача:* %s\n" +
        "*Старый статус:* %s\n" +
        "*Новый статус:* %s\n" +
        "*Исполнитель:* %s\n\n" +
        "Ссылка: http://localhost:8080/tasks/%d",
        task.Title,
        oldStatus,
        task.Status,
        task.Assignee,
        task.ID,
    )
    
    msg := tgbotapi.NewMessageToChannel(tb.ChatID, message)
    msg.ParseMode = "Markdown"
    
    _, err := tb.bot.Send(msg)
    return err
}

// SendDailySummary - отправляет ежедневный отчет
func (tb *TelegramBot) SendDailySummary(
    totalTasks int,
    completedToday int,
    upcomingDeadlines []models.Task,
    overdueTasks []models.Task,
) error {
    message := fmt.Sprintf(
        "📊 *Ежедневный отчет*\n" +
        "*Всего задач:* %d\n" +
        "*Выполнено сегодня:* %d\n\n",
        totalTasks,
        completedToday,
    )
    
    if len(overdueTasks) > 0 {
        message += "🚨 *Просроченные задачи:*\n"
        for _, task := range overdueTasks {
            overdue := time.Since(*task.Deadline)
            message += fmt.Sprintf(
                "• %s (%s) - просрочено %dч\n",
                task.Title,
                task.Assignee,
                int(overdue.Hours()),
            )
        }
        message += "\n"
    }
    
    if len(upcomingDeadlines) > 0 {
        message += "⏰ *Ближайшие дедлайны (24ч):*\n"
        for _, task := range upcomingDeadlines {
            hoursLeft := int(time.Until(*task.Deadline).Hours())
            message += fmt.Sprintf(
                "• %s (%s) - через %dч\n",
                task.Title,
                task.Assignee,
                hoursLeft,
            )
        }
    }
    
    if len(overdueTasks) == 0 && len(upcomingDeadlines) == 0 {
        message += "✅ Все задачи в порядке! Нет просроченных и ближайших дедлайнов."
    }
    
    msg := tgbotapi.NewMessageToChannel(tb.ChatID, message)
    msg.ParseMode = "Markdown"
    
    _, err := tb.bot.Send(msg)
    return err
}

// SendTestMessage - отправляет тестовое сообщение
func (tb *TelegramBot) SendTestMessage() error {
    message := "✅ *Kanban Calendar Bot активирован!*\nБот готов отправлять уведомления о дедлайнах."
    
    msg := tgbotapi.NewMessageToChannel(tb.ChatID, message)
    msg.ParseMode = "Markdown"
    
    _, err := tb.bot.Send(msg)
    return err
}