-- Seeds для создания администратора
-- Тестовый пользователь-администратор для входа в систему

USE `reports`;

-- Установка кодировки UTF-8 для корректного отображения русских символов
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- Вставка администратора
-- Пароль: password123 (bcrypt hash: $2a$10$ROHJgkE4RZOlwPp9hTw4Uu9pZ2ga1eeUuuFQ4lMg06UotIFBiarRK)
-- Примечание: Если хэш в базе обрезан (начинается с точки), используйте UPDATE ниже
INSERT INTO `users` (
    `name`, `email`, `password`, `timezone`, `language`, `is_active`, 
    `created_at`, `updated_at`
) VALUES (
    'Администратор',
    'admin@test.ru',
    '$2a$10$ROHJgkE4RZOlwPp9hTw4Uu9pZ2ga1eeUuuFQ4lMg06UotIFBiarRK',
    'Europe/Moscow',
    'ru',
    1,
    NOW(),
    NOW()
)
ON DUPLICATE KEY UPDATE
    `name` = VALUES(`name`),
    `password` = VALUES(`password`),
    `updated_at` = NOW();

-- Исправление пароля, если он был обрезан при первоначальной вставке
-- Это нужно для случаев, когда MySQL интерпретирует $ как переменную
UPDATE `users` 
SET `password` = '$2a$10$ROHJgkE4RZOlwPp9hTw4Uu9pZ2ga1eeUuuFQ4lMg06UotIFBiarRK' 
WHERE `email` = 'admin@test.ru' 
  AND (`password` NOT LIKE '$2a$%' OR `password` IS NULL OR LENGTH(`password`) < 60);

-- Назначение администратору роли admin на все проекты
-- Получаем ID пользователя admin@test.ru и назначаем роль admin на все существующие проекты
INSERT INTO `user_project_roles` (
    `user_id`, `project_id`, `role`, `created_at`, `updated_at`
)
SELECT 
    (SELECT `id` FROM `users` WHERE `email` = 'admin@test.ru' LIMIT 1) AS `user_id`,
    `id` AS `project_id`,
    'admin' AS `role`,
    NOW() AS `created_at`,
    NOW() AS `updated_at`
FROM `projects`
WHERE (SELECT `id` FROM `users` WHERE `email` = 'admin@test.ru' LIMIT 1) IS NOT NULL
ON DUPLICATE KEY UPDATE
    `role` = VALUES(`role`),
    `updated_at` = NOW();

