# setup.ps1 - Скрипт настройки Kanban Calendar
Write-Host "==========================================" -ForegroundColor Green
Write-Host "   Настройка Kanban Calendar + Telegram   " -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Green
Write-Host ""

# Проверяем наличие .env файла
if (-not (Test-Path .env)) {
    Write-Host "Файл .env не найден. Создаем из примера..." -ForegroundColor Yellow
    Copy-Item .env.example .env -ErrorAction SilentlyContinue
    Write-Host "Файл .env создан. Отредактируй его и добавь Telegram токены!" -ForegroundColor Cyan
    Write-Host ""
}

# Создаем директории если их нет
$directories = @("logs", "migrations")
foreach ($dir in $directories) {
    if (-not (Test-Path $dir)) {
        New-Item -ItemType Directory -Path $dir -Force | Out-Null
        Write-Host "Создана директория: $dir" -ForegroundColor Gray
    }
}

# Проверяем Docker
Write-Host "`nПроверка Docker..." -ForegroundColor Blue
try {
    $dockerVersion = docker --version 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Docker установлен: $dockerVersion" -ForegroundColor Green
    } else {
        Write-Host "✗ Docker не найден. Установите Docker Desktop для Windows" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "✗ Docker не найден. Установите Docker Desktop для Windows" -ForegroundColor Red
    exit 1
}

# Показываем меню
Write-Host "`nДоступные команды:" -ForegroundColor Cyan
Write-Host "  1. Запустить приложение через Docker Compose" -ForegroundColor White
Write-Host "  2. Запустить вручную (go run)" -ForegroundColor White
Write-Host "  3. Выполнить миграции" -ForegroundColor White
Write-Host "  4. Создать Telegram бота" -ForegroundColor White
Write-Host "  5. Выход" -ForegroundColor White
Write-Host ""

$choice = Read-Host "Выберите действие (1-5)"

switch ($choice) {
    "1" {
        Write-Host "`nЗапуск через Docker Compose..." -ForegroundColor Blue
        Write-Host "Используются переменные из .env файла" -ForegroundColor Gray
        docker-compose up -d
        Write-Host "`nПриложение запущено!" -ForegroundColor Green
        Write-Host "API доступно по: http://localhost:8080" -ForegroundColor Green
        Write-Host "База данных: localhost:5432" -ForegroundColor Green
    }
    "2" {
        Write-Host "`nЗапуск вручную..." -ForegroundColor Blue
        Write-Host "Убедись что PostgreSQL запущена на localhost:5432" -ForegroundColor Yellow
        
        # Проверяем зависимости
        Write-Host "Проверка зависимостей..." -ForegroundColor Gray
        go mod tidy
        
        # Запускаем
        Write-Host "`nЗапуск приложения..." -ForegroundColor Green
        go run main.go
    }
    "3" {
        Write-Host "`nВыполнение миграций..." -ForegroundColor Blue
        if (Test-Path "migrations") {
            Get-ChildItem migrations/*.sql | ForEach-Object {
                Write-Host "  Миграция: $($_.Name)" -ForegroundColor Gray
            }
        } else {
            Write-Host "Директория migrations не найдена" -ForegroundColor Red
        }
    }
    "4" {
        Write-Host "`nИнструкция по созданию Telegram бота:" -ForegroundColor Blue
        Write-Host "1. Откройте Telegram и найдите @BotFather" -ForegroundColor White
        Write-Host "2. Отправьте /newbot и следуйте инструкциям" -ForegroundColor White
        Write-Host "3. Получите токен вида: 1234567890:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw" -ForegroundColor White
        Write-Host "4. Добавьте бота в группу или напишите ему" -ForegroundColor White
        Write-Host "5. Откройте ссылку в браузере:" -ForegroundColor White
        Write-Host "   https://api.telegram.org/bot<ВАШ_ТОКЕН>/getUpdates" -ForegroundColor Cyan
        Write-Host "6. Найдите 'chat':{'id':XXXXXX} в ответе" -ForegroundColor White
        Write-Host "7. Добавьте токен и chat_id в файл .env" -ForegroundColor White
        Write-Host ""
        Read-Host "Нажмите Enter для продолжения"
    }
    "5" {
        Write-Host "Выход..." -ForegroundColor Gray
        exit 0
    }
    default {
        Write-Host "Неверный выбор" -ForegroundColor Red
    }
}