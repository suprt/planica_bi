-- +migrate Down
DROP TABLE IF EXISTS `seo_queries_monthly`;
DROP TABLE IF EXISTS `direct_totals_monthly`;
DROP TABLE IF EXISTS `direct_campaign_monthly`;
DROP TABLE IF EXISTS `direct_campaigns`;
DROP TABLE IF EXISTS `metrics_age_monthly`;
DROP TABLE IF EXISTS `metrics_monthly`;
DROP TABLE IF EXISTS `goals`;
DROP TABLE IF EXISTS `yandex_direct_accounts`;
DROP TABLE IF EXISTS `yandex_counters`;
DROP TABLE IF EXISTS `user_project_roles`;
DROP TABLE IF EXISTS `projects`;
DROP TABLE IF EXISTS `users`;
