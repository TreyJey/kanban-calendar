## Kanban Calendar API Integration

Backend-сервис для управления задачами канбан-доски с поддержкой импорта из внешних календарей и системой уведомлений через Telegram.

## Стек технологий

Language: Go (Golang)

Database: PostgreSQL

Containerization: Docker / Docker-compose

Standard: iCalendar (RFC 5545)

Notifications: Telegram Bot API

## Telegram Бот

Сервис включает в себя планировщик уведомлений (Notification Scheduler), который отслеживает дедлайны задач.

Функционал бота:

Автоматическая отправка уведомлений о приближающихся дедлайнах.

Логирование отправленных уведомлений в таблицу notifications для предотвращения повторов.

Настройка уведомлений:
Для корректной работы системы в файле .env необходимо указать:

 - TELEGRAM_BOT_TOKEN — токен, полученный от @BotFather.

 - TELEGRAM_CHAT_ID — ID чата (личного или группового), куда бот будет слать алерты.

## Быстрый запуск

Клонируйте репозиторий:

git clone <ссылка-на-репозиторий>
cd kanban-calendar

Настройте окружение:
Создайте файл .env на основе примера:

DB_HOST=kanban-postgres
DB_PORT=5432
TELEGRAM_TOKEN=your_telegram_bot_token_here
TELEGRAM_CHAT_ID=your_telegram_chat_id_here
FRONTEND_URL=http://localhost:3000

Запустите через Docker:

docker-compose up --build
Сервер будет доступен по адресу: http://localhost:8080

## API Documentation

## API Endpoints

Метод,Путь,Описание,Тело запроса (Request Body)
GET,/api/tasks,Получить все задачи,—
GET,/api/tasks/:id,Получить задачу по ID,—
GET,/api/tasks/status/:status,Получить задачи по статусу,—
GET,/api/calendar/events,Задачи в формате событий календаря,—
GET,/api/health,Проверка состояния сервиса,—
POST,/api/tasks,Создать новую задачу,JSON (см. структуру ниже)
POST,/api/tasks/import,Импорт календаря (.ics),multipart/form-data (key: calendar)
PUT,/api/tasks/:id,Обновить существующую задачу,JSON (см. структуру ниже)
DELETE,/api/tasks/:id,Удалить задачу,—



Структура JSON (для POST и PUT)
При создании или обновлении задачи используйте следующий формат:
{
  "title": "Собрать NPM модуль",       // (string) Обязательно
  "description": "Подготовить проект", // (string)
  "status": "in_progress",            // "todo", "in_progress", "done"
  "deadline": "2026-01-20T15:00:00Z", // (string, ISO 8601)
  "start_date": "2026-01-20T10:00:00Z",
  "end_date": "2026-01-20T11:00:00Z",
  "assignee": "Frontend Dev"
}

## Структура БД

Система автоматически создает и управляет двумя таблицами:

tasks — хранение данных о задачах и внешних ID.

notifications — история отправленных уведомлений, связанная с задачами.