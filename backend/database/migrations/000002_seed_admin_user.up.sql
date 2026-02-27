-- +migrate Up
-- Создание дефолтного пользователя-администратора
-- Email: admin@test.ru
-- Пароль: password123 (bcrypt hash)

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

-- Назначение роли admin на все существующие проекты
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
