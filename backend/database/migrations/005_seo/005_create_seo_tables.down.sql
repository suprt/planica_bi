USE `planica_bi`;

DROP TABLE IF EXISTS `seo_queries_monthly`;

-- Запись в аудит об удалении таблиц
INSERT INTO `schema_audit` (`migration_version`, `change_type`, `object_type`, `object_name`) 
VALUES 
('005_create_seo_tables', 'DROP', 'TABLE', 'seo_queries_monthly');