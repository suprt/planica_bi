# Задачи бэкенда для админ-панели

**Дата:** 2025-11-15  
**Статус:** ❌ **НЕ РЕАЛИЗОВАНО**

---

## 🔴 Критичные задачи для админ-панели

Админ-панель требует от бэкенда следующие API endpoints, которые **отсутствуют**:

### 1. ❌ Управление пользователями (User Management)

**Текущее состояние:**
- ✅ Модели `User` и `UserProjectRole` существуют
- ✅ `UserRepository` с базовыми методами (Create, GetByEmail, GetByID)
- ✅ `AuthService` с Register и Login
- ❌ **НЕТ API endpoints для управления пользователями**

**Что нужно добавить:**

#### 1.1. Список всех пользователей (для админа)
```
GET /api/users
```
- Требует роль: `admin`
- Возвращает список всех пользователей с их ролями
- Пагинация (опционально)
- Фильтрация по email, имени, статусу (опционально)

#### 1.2. Получение пользователя по ID
```
GET /api/users/:id
```
- Требует роль: `admin`
- Возвращает детальную информацию о пользователе с его ролями в проектах

#### 1.3. Создание пользователя (админом)
```
POST /api/users
```
- Требует роль: `admin`
- Создает нового пользователя
- Пароль генерируется автоматически или задается админом
- Возвращает созданного пользователя (без пароля)

#### 1.4. Обновление пользователя
```
PUT /api/users/:id
```
- Требует роль: `admin`
- Обновляет данные пользователя (имя, email, timezone, language, is_active)
- Смена пароля (опционально)

#### 1.5. Удаление пользователя
```
DELETE /api/users/:id
```
- Требует роль: `admin`
- Удаляет пользователя (каскадно удаляются роли)

---

### 2. ❌ Управление ролями пользователей в проектах

**Текущее состояние:**
- ✅ Модель `UserProjectRole` существует
- ✅ `UserRepository.GetUserProjectRole()` - получение роли
- ✅ `UserRepository.GetUserProjects()` - получение проектов пользователя
- ❌ **НЕТ методов для создания/обновления/удаления ролей**
- ❌ **НЕТ API endpoints для управления ролями**

**Что нужно добавить:**

#### 2.1. Назначение роли пользователю в проекте
```
POST /api/projects/:id/users
```
- Требует роль: `admin` или `manager` (для своего проекта)
- Назначает роль пользователю в проекте
- Body: `{ "user_id": 1, "role": "client" }`
- Если роль уже существует - обновляет её

#### 2.2. Изменение роли пользователя в проекте
```
PUT /api/projects/:id/users/:userId
```
- Требует роль: `admin` или `manager` (для своего проекта)
- Изменяет роль пользователя в проекте
- Body: `{ "role": "manager" }`

#### 2.3. Удаление роли пользователя из проекта
```
DELETE /api/projects/:id/users/:userId
```
- Требует роль: `admin` или `manager` (для своего проекта)
- Удаляет роль пользователя из проекта (убирает доступ)

#### 2.4. Список пользователей проекта
```
GET /api/projects/:id/users
```
- Требует роль: `admin`, `manager` или `client` (для своего проекта)
- Возвращает список всех пользователей проекта с их ролями

#### 2.5. Список проектов пользователя
```
GET /api/users/:id/projects
```
- Требует роль: `admin` или сам пользователь
- Возвращает список всех проектов пользователя с его ролями

---

### 3. ⚠️ Дополнительные методы в репозиториях

**Что нужно добавить в `UserRepository`:**

```go
// AssignRole assigns a role to user in project
AssignRole(ctx context.Context, userID, projectID uint, role string) error

// UpdateRole updates user's role in project
UpdateRole(ctx context.Context, userID, projectID uint, role string) error

// RemoveRole removes user's role from project
RemoveRole(ctx context.Context, userID, projectID uint) error

// GetProjectUsers retrieves all users for a project with their roles
GetProjectUsers(ctx context.Context, projectID uint) ([]models.UserProjectRole, error)

// GetAllUsers retrieves all users (for admin)
GetAllUsers(ctx context.Context) ([]models.User, error)

// Update updates user information
Update(ctx context.Context, user *models.User) error

// Delete deletes a user
Delete(ctx context.Context, userID uint) error
```

---

### 4. ❌ UserService для бизнес-логики

**Что нужно создать:**

```go
// backend/internal/services/user_service.go

type UserService struct {
    userRepo repositories.UserRepositoryInterface
}

// GetAllUsers returns all users (admin only)
GetAllUsers(ctx context.Context) ([]UserPublic, error)

// GetUserByID returns user by ID (admin only)
GetUserByID(ctx context.Context, userID uint) (*UserPublic, error)

// CreateUser creates a new user (admin only)
CreateUser(ctx context.Context, req *CreateUserRequest) (*UserPublic, error)

// UpdateUser updates user information (admin only)
UpdateUser(ctx context.Context, userID uint, req *UpdateUserRequest) (*UserPublic, error)

// DeleteUser deletes a user (admin only)
DeleteUser(ctx context.Context, userID uint) error

// AssignProjectRole assigns role to user in project
AssignProjectRole(ctx context.Context, projectID, userID uint, role string) error

// UpdateProjectRole updates user's role in project
UpdateProjectRole(ctx context.Context, projectID, userID uint, role string) error

// RemoveProjectRole removes user's role from project
RemoveProjectRole(ctx context.Context, projectID, userID uint) error

// GetProjectUsers returns all users for a project
GetProjectUsers(ctx context.Context, projectID uint) ([]ProjectUserResponse, error)

// GetUserProjects returns all projects for a user
GetUserProjects(ctx context.Context, userID uint) ([]UserProjectResponse, error)
```

---

### 5. ❌ UserHandler для HTTP endpoints

**Что нужно создать:**

```go
// backend/internal/handlers/user_handler.go

type UserHandler struct {
    userService UserServiceInterface
}

// GetAllUsers handles GET /api/users
GetAllUsers(c echo.Context) error

// GetUser handles GET /api/users/:id
GetUser(c echo.Context) error

// CreateUser handles POST /api/users
CreateUser(c echo.Context) error

// UpdateUser handles PUT /api/users/:id
UpdateUser(c echo.Context) error

// DeleteUser handles DELETE /api/users/:id
DeleteUser(c echo.Context) error

// GetProjectUsers handles GET /api/projects/:id/users
GetProjectUsers(c echo.Context) error

// AssignProjectRole handles POST /api/projects/:id/users
AssignProjectRole(c echo.Context) error

// UpdateProjectRole handles PUT /api/projects/:id/users/:userId
UpdateProjectRole(c echo.Context) error

// RemoveProjectRole handles DELETE /api/projects/:id/users/:userId
RemoveProjectRole(c echo.Context) error

// GetUserProjects handles GET /api/users/:id/projects
GetUserProjects(c echo.Context) error
```

---

## 📋 Итоговый список задач

### Критичные (для работы админ-панели):

1. ✅ Модели и репозитории (частично готовы)
2. ❌ **Добавить методы в `UserRepository` для работы с ролями**
3. ❌ **Создать `UserService` с бизнес-логикой**
4. ❌ **Создать `UserHandler` с HTTP endpoints**
5. ❌ **Добавить routes в `router.go`:**
   - `GET /api/users` - список пользователей (admin)
   - `GET /api/users/:id` - детали пользователя (admin)
   - `POST /api/users` - создание пользователя (admin)
   - `PUT /api/users/:id` - обновление пользователя (admin)
   - `DELETE /api/users/:id` - удаление пользователя (admin)
   - `GET /api/projects/:id/users` - пользователи проекта
   - `POST /api/projects/:id/users` - назначение роли
   - `PUT /api/projects/:id/users/:userId` - изменение роли
   - `DELETE /api/projects/:id/users/:userId` - удаление роли
   - `GET /api/users/:id/projects` - проекты пользователя

---

## 🎯 Вывод

**Админ-панель НЕ может работать без этих endpoints!**

Frontend разработчик может создать интерфейс, но без этих API endpoints админ-панель не сможет:
- Показывать список пользователей
- Создавать/редактировать/удалять пользователей
- Назначать роли пользователям в проектах
- Управлять доступом к проектам

**Это задачи бэкенда, которые нужно реализовать!**

---

## 📝 Приоритет

**Критичный** - без этих endpoints админ-панель не сможет выполнять основные функции управления пользователями и ролями.

