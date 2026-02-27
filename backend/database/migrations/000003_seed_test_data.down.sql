-- +migrate Down
-- Удаление тестовых данных

DELETE FROM `seo_queries_monthly` WHERE `project_id` = 1;
DELETE FROM `direct_totals_monthly` WHERE `project_id` = 1;
DELETE FROM `direct_campaign_monthly` WHERE `campaign_id` IN (SELECT `id` FROM `direct_campaigns` WHERE `project_id` = 1);
DELETE FROM `direct_campaigns` WHERE `project_id` = 1;
DELETE FROM `metrics_age_monthly` WHERE `project_id` = 1;
DELETE FROM `metrics_monthly` WHERE `project_id` = 1;
DELETE FROM `projects` WHERE `id` = 1;
