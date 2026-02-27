# Planica BI

Система бизнес-аналитики для агрегации и анализа данных из Яндекс.Метрики и Яндекс.Директа с использованием AI для выявления трендов и рекомендаций.

## 🚀 Основные возможности

- **Интеграция с Яндекс.Метрикой и Яндекс.Директом** — автоматическая синхронизация данных
- **Автоматическая генерация отчетов** — ежемесячные отчеты с анализом метрик
- **AI-аналитика** — анализ данных с помощью Ollama (встроенный Go-модуль) и выдача рекомендаций
- **Административная панель** — управление проектами, пользователями и ролями
- **Публичные отчеты** — генерация публичных ссылок для отчетов
- **Система ролей** — гибкое управление доступом на уровне проектов
- **Асинхронные задачи** — фоновая обработка через Asynq

## 🛠 Технологии

**Backend:**
- Go 1.24+ (Echo framework)
- MySQL 9.5
- Redis 7
- Asynq (асинхронные задачи)
- GORM (ORM)
- JWT аутентификация

**Frontend:**
- React 19
- TypeScript
- React Router
- React Admin (админ-панель)
- Axios

**AI:**
- Ollama API (Go-интеграция)

**DevOps:**
- Docker & Docker Compose
- Nginx

## 📋 Требования

- Docker & Docker Compose
- Git

## ⚙️ Быстрый старт

1. **Клонируйте репозиторий:**
```bash
git clone https://github.com/suprt/planica_bi.git
cd planica_bi
```

2. **Настройте переменные окружения:**
Скопируйте `backend/.env.example` в `backend/.env` и заполните необходимые значения:
```bash
# Database
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USERNAME=reports
DB_PASSWORD=your-secure-password-here
DB_DATABASE=reports

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# App
APP_PORT=8080
APP_ENV=production

# JWT
JWT_SECRET=your-jwt-secret-key-change-in-production

# Ollama AI
OLLAMA_API_KEY=your-api-key
OLLAMA_MODEL=glm-4.6

# Yandex OAuth (заполните для интеграции)
YANDEX_CLIENT_ID=your-client-id
YANDEX_CLIENT_SECRET=your-client-secret
```

3. **Запустите приложение:**
```bash
docker-compose up -d
```

4. **Проверьте статус:**
```bash
docker-compose ps
```

## 🌐 Доступ к приложению

- **Frontend:** http://localhost:3000
- **Backend API:** http://localhost:8080
- **Adminer (DB):** http://localhost:8081
- **Redis:** localhost:6379

## 🔐 Учетные данные (для разработки)

После первого запуска миграций создаётся тестовый аккаунт:

- Email: `admin@test.ru`
- Password: `password123`

> **Примечание:** В production-окружении создайте своего пользователя через API или добавьте свою миграцию.

## 📁 Структура проекта

```
planica_bi/
├── backend/              # Go backend
│   ├── cmd/api/         # Точка входа
│   ├── internal/        # Внутренние пакеты
│   │   ├── ai/          # AI-интеграция (Ollama)
│   │   ├── cache/       # Redis кэширование
│   │   ├── config/      # Конфигурация приложения
│   │   ├── cron/        # Планировщик задач
│   │   ├── database/    # Работа с БД
│   │   ├── handlers/    # HTTP handlers
│   │   ├── services/    # Бизнес-логика
│   │   ├── repositories/# Доступ к данным
│   │   ├── models/      # Модели данных
│   │   ├── router/      # Маршрутизация и middleware
│   │   ├── queue/       # Асинхронные задачи (Asynq)
│   │   ├── middleware/  # Middleware (валидация, rate limiter, pagination)
│   │   ├── logger/      # Логирование
│   │   └── integrations/# Внешние интеграции
│   ├── database/        # SQL миграции
│   │   ├── migrate.go   # Миграции
│   │   └── migrations/  # Файлы миграций
│   ├── pkg/             # Публичные пакеты
│   └── storage/         # Логи и файлы
├── frontend/            # React frontend
│   ├── src/
│   │   ├── pages/       # Страницы приложения
│   │   ├── components/  # React компоненты
│   │   ├── admin/       # React Admin панель
│   │   ├── services/    # API клиенты
│   │   ├── contexts/    # React Context
│   │   ├── hooks/       # Custom hooks
│   │   ├── types/       # TypeScript типы
│   │   └── utils/       # Утилиты
└── docker-compose.yml   # Docker конфигурация
```

## 🔄 API Endpoints

### Аутентификация
- `POST /api/auth/register` - Регистрация
- `POST /api/auth/login` - Вход

### Проекты
- `GET /api/projects` - Список проектов
- `GET /api/projects/:id` - Детали проекта
- `GET /api/projects/:id/public-link` - Публичная ссылка на отчет

### Отчеты
- `GET /api/reports/:projectId` - Получить отчет
- `POST /api/reports/:projectId/generate` - Сгенерировать отчет
- `GET /api/reports/:projectId/status` - Статус генерации

### Синхронизация
- `POST /api/sync/:projectId` - Принудительная синхронизация

### OAuth
- `GET /api/oauth/yandex` - Инициализация OAuth
- `GET /api/oauth/yandex/callback` - OAuth callback

## 🔧 Разработка

### Backend

```bash
cd backend
go run cmd/api/main.go
```

### Frontend

```bash
cd frontend
npm install
npm start
```

## 📝 Основные функции

1. **Автоматическая синхронизация** — ежедневная синхронизация данных из Яндекс.Метрики и Яндекс.Директа
2. **Генерация отчетов** — автоматическая генерация отчетов с анализом метрик за 3 месяца
3. **AI анализ** — автоматический анализ трендов и выдача рекомендаций
4. **Управление доступом** — система ролей (admin, manager, viewer) на уровне проектов
5. **Кэширование** — Redis кэш для ускорения работы с отчетами

## 🐳 Docker команды

```bash
# Запуск
docker-compose up -d

# Остановка
docker-compose down

# Пересборка
docker-compose build

# Логи
docker-compose logs -f backend
docker-compose logs -f frontend

# Перезапуск сервиса
docker-compose restart backend
```

## 📚 Дополнительная информация

- Логи приложения: `backend/storage/logs/app.log`
- База данных инициализируется автоматически при первом запуске через миграции (`backend/database/migrations/`)
- Для запуска миграций вручную: `go run backend/database/migrate.go`
