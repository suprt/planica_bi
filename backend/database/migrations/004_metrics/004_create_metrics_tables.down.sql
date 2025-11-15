USE `reports`;

DROP TABLE IF EXISTS `direct_totals_monthly`;
DROP TABLE IF EXISTS `direct_campaign_monthly`;
DROP TABLE IF EXISTS `metrics_age_monthly`;
DROP TABLE IF EXISTS `metrics_monthly`;

-- Запись в аудит об удалении таблиц
INSERT INTO `schema_audit` (`migration_version`, `change_type`, `object_type`, `object_name`) 
VALUES 
('004_create_metrics_tables', 'DROP', 'TABLE', 'direct_totals_monthly'),
('004_create_metrics_tables', 'DROP', 'TABLE', 'direct_campaign_monthly'),
('004_create_metrics_tables', 'DROP', 'TABLE', 'metrics_age_monthly'),
('004_create_metrics_tables', 'DROP', 'TABLE', 'metrics_monthly');