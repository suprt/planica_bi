USE `reports`;

DROP TABLE IF EXISTS `projects`;

-- Запись в аудит об удалении таблицы
INSERT INTO `schema_audit` (`migration_version`, `change_type`, `object_type`, `object_name`) 
VALUES ('002_create_projects_table', 'DROP', 'TABLE', 'projects');