# План улучшений Planica BI

## ✅ Выполнено (26 февраля 2026)

### 1. Рефакторинг логгирования
- **Проблема:** `APP_ENV=production` не имел реального функционала
- **Решение:** Переключились на `APP_DEBUG=true/false`
- **Файлы:**
  - `backend/internal/config/config.go` — добавлено `AppDebug bool`
  - `backend/internal/logger/logger.go` — изменена сигнатура `Init(isDebug, logPath)`
  - `backend/cmd/api/main.go` — обновлён вызов логгера
  - `backend/.env` — `APP_DEBUG=true`, `APP_ENV` удалён
  - `backend/.env.example` — создан шаблон с комментариями
  - `.gitignore` — игнорирует все `.env` файлы

### 2. Безопасность (Приоритет 1)
- **JWT_SECRET** — добавлен и сгенерирован (32 символа, base64)
- **DB_PASSWORD** — заменён с `1234` на надёжный (48 символов, base64)
- **Файлы:**
  - `backend/.env` — добавлены `JWT_SECRET`, обновлён `DB_PASSWORD`
  - `backend/.env.example` — шаблон с инструкцией
  - `docker-compose.yml` — убраны значения по умолчанию для MySQL

### 3. Health Check (Приоритет 2)
- **Endpoint `/health`** — проверка статуса приложения
- **Endpoint `/ready`** — проверка готовности (для K8s/Docker)
- **Docker healthcheck** — для backend сервиса
- **Файлы:**
  - `backend/internal/handlers/health_handler.go` — новый хендлер
  - `backend/internal/router/router.go` — добавлены routes
  - `backend/Dockerfile` — добавлен `wget` для healthcheck
  - `docker-compose.yml` — healthcheck для backend

### 4. Валидация данных (Приоритет 2)
- **go-playground/validator v10** — добавлена валидация запросов
- **Middleware** — `middleware.ValidateRequest()` для проверки DTO
- **Теги валидации** — `required`, `email`, `min`, `max`, `oneof`
- **Файлы:**
  - `backend/internal/middleware/validator.go` — новый middleware
  - `backend/internal/services/auth_service.go` — теги на `RegisterRequest`, `LoginRequest`
  - `backend/internal/services/user_service.go` — теги на `CreateUserRequest`, `UpdateUserRequest`, `AssignRoleRequest`
  - `backend/internal/handlers/auth_handler.go` — вызов валидации
  - `backend/internal/handlers/user_handler.go` — вызов валидации
  - `backend/internal/handlers/project_user_handler.go` — вызов валидации
  - `backend/cmd/api/main.go` — инициализация валидатора

### 5. Пагинация API (Приоритет 2)
- **Query-параметры** — `page`, `per_page`, `sort`, `order`
- **Middleware** — `middleware.GetPagination()` для извлечения параметров
- **Ответ API** — `{ data: [...], total: N }`
- **Файлы:**
  - `backend/internal/middleware/pagination.go` — новый middleware
  - `backend/internal/repositories/project_repository.go` — `GetAllPaginated()`, `GetByUserIDPaginated()`
  - `backend/internal/repositories/user_repository.go` — `GetAllPaginated()`
  - `backend/internal/repositories/user_repository_interface.go` — добавлен метод в интерфейс
  - `backend/internal/services/project_service.go` — `GetAllProjectsPaginated()`
  - `backend/internal/services/user_service.go` — `GetAllUsersPaginated()`
  - `backend/internal/handlers/project_handler.go` — пагинация `/api/projects`
  - `backend/internal/handlers/user_handler.go` — пагинация `/api/users`
  - `backend/internal/handlers/user_service_interface.go` — добавлен метод в интерфейс
  - `backend/internal/handlers/project_handler.go` — добавлен метод в интерфейс

### 6. Ротация логов (Приоритет 3)
- **lumberjack** — автоматическая ротация, архивация, сжатие логов
- **Настройки** — размер файла, кол-во архивов, срок хранения, компрессия
- **Файлы:**
  - `backend/internal/logger/logger.go` — интеграция lumberjack
  - `backend/internal/config/config.go` — поля `LogMaxSize`, `LogMaxBackups`, `LogMaxAge`, `LogCompress`
  - `backend/cmd/api/main.go` — передача параметров в logger.Init()
  - `backend/.env.example` — переменные `LOG_MAX_SIZE`, `LOG_MAX_BACKUPS`, `LOG_MAX_AGE`, `LOG_COMPRESS`
  - `backend/.env` — настройки ротации

### 7. Token в sessionStorage (Приоритет 3)
- **Безопасность** — замена localStorage на sessionStorage для токенов
- **Причина** — защита от XSS, очистка при закрытии вкладки
- **Файлы:**
  - `frontend/src/services/api/authService.ts` — все методы заменены на sessionStorage
  - `frontend/src/services/api/apiClient.ts` — interceptor обновлён
  - `frontend/src/utils/projectStorage.ts` — хранение проекта в sessionStorage
  - `frontend/src/admin/dataProvider.ts` — получение токена
  - `frontend/src/admin/authProvider.ts` — комментарий обновлён
  - `frontend/src/admin/resources/UserProjects.tsx` — получение токена
  - `frontend/src/contexts/AuthContext.tsx` — комментарии обновлены
  - `frontend/src/components/Dashboard/Dashboard.tsx` — комментарий обновлён
  - **Исключение:** `frontend/src/contexts/ThemeContext.tsx` — оставлен localStorage (тема должна сохраняться между сессиями)

### 8. Unit-тесты (Приоритет 3)
- **Table-driven tests** — стандарт для Go
- **Mock репозиториев** — ручные моки для изоляции сервисов
- **Покрытие:**
  - `AuthService` — 17 тестов (Register, Login, ValidateToken)
  - `ProjectService` — 9 тестов (Create, Get, GetAll, Update, Delete)
  - `UserService` — 14 тестов (Create, Update, Delete, AssignRole)
  - `MarketingService` — 2 теста (GetMarketingData)
  - `CounterService` — 7 тестов (CreateCounter, GetCountersByProject)
  - `MetricsService` — 2 теста (GetMetricsWithData)
  - `DirectService` — 6 тестов (CreateAccount, GetAccountsByProject)
  - `GoalService` — 9 тестов (CreateGoal, GetGoal, GetGoalsByCounter, DeleteGoal)
- **Файлы:**
  - `backend/internal/services/*_test.go` — 8 файлов с тестами
  - `backend/internal/services/interfaces.go` — интерфейсы для всех репозиториев
- **Исправления в сервисах:**
  - `auth_service.go` — проверка на nil при GetByEmail
  - `user_service.go` — проверка на nil при GetByID (UpdateUser, DeleteUser, AssignRole)
  - `goal_service.go` — проверка на nil при GetByID (CreateGoal, GetGoalsByCounter)

### 9. Критичные улучшения (Приоритет 1) — 26 февраля 2026
- **Rate Limiting** — защита от brute force и DDoS
  - `golang.org/x/time/rate` — token bucket алгоритм
  - Общий limiter: 10 req/s, burst 20
  - Строгий limiter для `/api/auth/*`: 2 req/s, burst 5
  - Очистка старых ключей каждые 1 минуту
  - **Файлы:**
    - `backend/internal/middleware/rate_limiter.go` — новый middleware
    - `backend/internal/router/router.go` — применение limiter
    - `backend/cmd/api/main.go` — shutdown limiter
- **CORS настройка** — whitelist доменов
  - Разрешены только запросы с `cfg.FrontendURL`
  - Разрешённые методы: GET, POST, PUT, DELETE, PATCH
  - Разрешённые заголовки: Origin, Content-Type, Accept, Authorization
  - **Файлы:**
    - `backend/internal/router/router.go` — `CORSWithConfig`
- **Graceful shutdown для очередей**
  - Контекст с таймаутом 30 секунд
  - Ожидание завершения текущих задач
  - **Файлы:**
    - `backend/internal/queue/worker.go` — улучшенный `Shutdown()`
    - `backend/cmd/api/main.go` — вызов shutdown
- **Миграции БД** — golang-migrate
  - Версионирование миграций
  - Up/Down миграции
  - Первая миграция: полная схема БД
  - **Файлы:**
    - `backend/database/migrate.go` — функции RunMigrations, RollbackMigrations
    - `backend/database/migrations/000001_initial_schema.up.sql` — up миграция
    - `backend/database/migrations/000001_initial_schema.down.sql` — down миграция
    - `backend/internal/database/database.go` — AutoMigrate остаётся для совместимости

---

## 🔄 В процессе / Следующие шаги

### Приоритет 1 (Критично)
- [x] **Секреты в .env** — реальные токены в `.env`, нужно очистить
- [x] **Слабый пароль MySQL** — `1234` → сгенерировать случайный
- [x] **JWT_SECRET по умолчанию** — требует обязательной установки
- [x] **Rate Limiting** — добавлен
- [x] **CORS** — настроен
- [x] **Graceful shutdown** — добавлен
- [x] **Миграции БД** — добавлены

### Приоритет 2 (Высокий)
- [x] **Health check** — endpoint `/health` + docker-compose healthcheck
- [x] **Валидация данных** — добавлен go-playground/validator
- [x] **Пагинация API** — `/api/projects`, `/api/users`

### Приоритет 3 (Средний)
- [x] **Ротация логов** — lumberjack для zap
- [x] **Unit-тесты** — покрыты основные сервисы (66 тестов)
- [x] **Token в sessionStorage** — вместо localStorage (frontend)

---

## 📝 Заметки

### Логирование
```bash
# Разработка (цветные логи, DEBUG уровень)
APP_DEBUG=true

# Production (JSON, только INFO+)
APP_DEBUG=false
```

### Сборка
```bash
cd backend
go build -o api.exe ./cmd/api
```

### Запуск в Docker
```bash
docker-compose up -d
docker-compose logs -f backend
```

### Rate Limiting
```bash
# Общий limiter: 10 запросов в секунду, burst 20
# Auth endpoints: 2 запроса в секунду, burst 5

# При превышении: HTTP 429 Too Many Requests
```

### Миграции
```bash
# Запуск миграций (в разработке через AutoMigrate)
# Для production использовать migrate CLI:
migrate -path database/migrations -database "mysql://user:pass@tcp(host:3306)/db" up
```
