# Админ-панель: React-admin

**Дата:** 2025-11-15  
**Решение:** React-admin  
**Статус:** В разработке

---

## Выбранное решение

**React-admin** - готовое решение для создания админ-панелей на React.

**Почему React-admin:**
- ✅ Готовое решение с готовыми компонентами
- ✅ Работает с любым REST API (наш Go бэкенд)
- ✅ Быстрая разработка
- ✅ Активное сообщество и документация
- ✅ Подходит для MVP

**Документация:**
- Официальный сайт: https://marmelab.com/react-admin/
- Документация: https://marmelab.com/react-admin/Documentation.html
- GitHub: https://github.com/marmelab/react-admin

---

## Требования к API

React-admin работает с REST API и ожидает следующие форматы:

### Список ресурсов
```
GET /api/users
Response: { "data": [...], "total": 100 }
```

### Получение одного ресурса
```
GET /api/users/:id
Response: { "data": {...} }
```

### Создание ресурса
```
POST /api/users
Body: { ... }
Response: { "data": {...} }
```

### Обновление ресурса
```
PUT /api/users/:id
Body: { ... }
Response: { "data": {...} }
```

### Удаление ресурса
```
DELETE /api/users/:id
Response: { "data": {...} }
```

---

## Endpoints для реализации

### 1. Управление пользователями
- `GET /api/users` - список всех пользователей
- `GET /api/users/:id` - детали пользователя
- `POST /api/users` - создание пользователя
- `PUT /api/users/:id` - обновление пользователя
- `DELETE /api/users/:id` - удаление пользователя

### 2. Управление ролями в проектах
- `GET /api/projects/:id/users` - пользователи проекта
- `POST /api/projects/:id/users` - назначение роли
- `PUT /api/projects/:id/users/:userId` - изменение роли
- `DELETE /api/projects/:id/users/:userId` - удаление роли
- `GET /api/users/:id/projects` - проекты пользователя

---

## Примечания

- Все endpoints требуют авторизации
- Endpoints для управления пользователями требуют роль `admin`
- Endpoints для управления ролями требуют роль `admin` или `manager` (для своего проекта)
- Формат ответов должен соответствовать ожиданиям react-admin


