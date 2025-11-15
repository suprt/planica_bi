USE `planica_bi`;

DROP TABLE IF EXISTS `goals`;
DROP TABLE IF EXISTS `direct_campaigns`;
DROP TABLE IF EXISTS `direct_accounts`;
DROP TABLE IF EXISTS `yandex_counters`;

-- Запись в аудит об удалении таблиц
INSERT INTO `schema_audit` (`migration_version`, `change_type`, `object_type`, `object_name`) 
VALUES 
('003_create_yandex_tables', 'DROP', 'TABLE', 'goals'),
('003_create_yandex_tables', 'DROP', 'TABLE', 'direct_campaigns'),
('003_create_yandex_tables', 'DROP', 'TABLE', 'direct_accounts'),
('003_create_yandex_tables', 'DROP', 'TABLE', 'yandex_counters');