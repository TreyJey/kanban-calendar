package models

import "time"

type Notification struct {
    ID        int       `json:"id"`
    TaskID    int       `json:"task_id"`
    Type      string    `json:"type"` // "deadline", "reminder", "status_change"
    Message   string    `json:"message"`
    SentAt    time.Time `json:"sent_at"`
    IsSent    bool      `json:"is_sent"`
    ChatID    string    `json:"chat_id,omitempty"`
}

const (
    NotificationTypeDeadline    = "deadline"
    NotificationTypeReminder    = "reminder"
    NotificationTypeStatusChange = "status_change"
    NotificationTypeDailyReport = "daily_report"
)