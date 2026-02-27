-- Тестовые данные для разработки
-- Применять только на dev окружении!

USE `reports`;

-- Вставка тестовых метрик для проекта test-project (project_id = 1)
INSERT INTO `metrics_monthly` (`project_id`, `year`, `month`, `visits`, `users`, `bounce_rate`, `avg_session_duration_sec`, `conversions`, `created_at`) VALUES
(1, 2025, 12, 15420, 12350, 35.5, 245, 1230, NOW()),
(1, 2026, 1, 18200, 14500, 32.0, 280, 1450, NOW()),
(1, 2026, 2, 16800, 13200, 33.2, 265, 1340, NOW())
ON DUPLICATE KEY UPDATE `visits` = VALUES(`visits`);

-- Тестовые данные по возрастам
INSERT INTO `metrics_age_monthly` (`project_id`, `year`, `month`, `age_group`, `visits`, `users`, `bounce_rate`, `avg_session_duration_sec`, `created_at`) VALUES
(1, 2026, 2, '18-24', 3200, 2800, 40.5, 180, NOW()),
(1, 2026, 2, '25-34', 5600, 4900, 30.2, 290, NOW()),
(1, 2026, 2, '35-44', 4200, 3500, 32.0, 260, NOW()),
(1, 2026, 2, '45-54', 2400, 1600, 35.5, 240, NOW()),
(1, 2026, 2, '55+', 1400, 400, 38.0, 220, NOW())
ON DUPLICATE KEY UPDATE `visits` = VALUES(`visits`);

-- Тестовые кампании Директа
INSERT INTO `direct_campaigns` (`project_id`, `direct_account_id`, `campaign_id`, `name`, `status`, `created_at`, `updated_at`) VALUES
(1, 1, 1001, 'Поиск - Бренд', 'on', NOW(), NOW()),
(1, 1, 1002, 'РСЯ - Товары', 'on', NOW(), NOW())
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- Тестовые данные по кампаниям
INSERT INTO `direct_campaign_monthly` (`project_id`, `direct_campaign_id`, `year`, `month`, `impressions`, `clicks`, `ctr_pct`, `cpc`, `conversions`, `cpa`, `cost`, `created_at`) VALUES
(1, 1, 2026, 2, 50000, 2500, 5.0, 45.50, 120, 950.00, 113750.00, NOW()),
(1, 2, 2026, 2, 80000, 3200, 4.0, 35.20, 95, 1180.00, 112640.00)
ON DUPLICATE KEY UPDATE `impressions` = VALUES(`impressions`);

-- Тестовые totals
INSERT INTO `direct_totals_monthly` (`project_id`, `year`, `month`, `impressions`, `clicks`, `ctr_pct`, `cpc`, `conversions`, `cpa`, `cost`, `created_at`) VALUES
(1, 2026, 2, 130000, 5700, 4.38, 39.50, 215, 1070.00, 225140.00, NOW())
ON DUPLICATE KEY UPDATE `impressions` = VALUES(`impressions`);

-- Тестовые SEO данные
INSERT INTO `seo_queries_monthly` (`project_id`, `year`, `month`, `query`, `position`, `url`, `created_at`) VALUES
(1, 2026, 2, 'купить электронику', 12, 'https://example.ru/electronics', NOW()),
(1, 2026, 2, 'строительные услуги', 8, 'https://example.ru/construction', NOW())
ON DUPLICATE KEY UPDATE `position` = VALUES(`position`);
