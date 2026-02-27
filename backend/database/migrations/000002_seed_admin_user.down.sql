-- +migrate Down
-- Удаление дефолтного пользователя-администратора

DELETE FROM `user_project_roles` WHERE `user_id` = (SELECT `id` FROM `users` WHERE `email` = 'admin@test.ru');
DELETE FROM `users` WHERE `email` = 'admin@test.ru';
