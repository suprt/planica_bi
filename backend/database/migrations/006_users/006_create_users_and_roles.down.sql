USE `reports`;

DROP TABLE IF EXISTS `user_project_roles`;
DROP TABLE IF EXISTS `users`;

-- Запись в аудит об удалении таблиц
INSERT INTO `schema_audit` (`migration_version`, `change_type`, `object_type`, `object_name`) 
VALUES 
('006_create_users_and_roles', 'DROP', 'TABLE', 'user_project_roles'),
('006_create_users_and_roles', 'DROP', 'TABLE', 'users');