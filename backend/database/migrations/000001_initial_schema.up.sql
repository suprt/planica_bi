-- +migrate Up
-- Create users table
CREATE TABLE IF NOT EXISTS `users` (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL UNIQUE,
    `password` VARCHAR(255) NOT NULL,
    `email_verified_at` TIMESTAMP NULL,
    `timezone` VARCHAR(255) DEFAULT 'Europe/Moscow',
    `language` ENUM('ru', 'en') DEFAULT 'ru',
    `is_active` BOOLEAN DEFAULT TRUE,
    `last_login_at` TIMESTAMP NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX `idx_users_email` (`email`),
    INDEX `idx_users_is_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create projects table
CREATE TABLE IF NOT EXISTS `projects` (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `name` TEXT NOT NULL,
    `slug` VARCHAR(255) NOT NULL UNIQUE,
    `public_token` VARCHAR(64) UNIQUE,
    `timezone` VARCHAR(255) DEFAULT 'Europe/Moscow',
    `currency` ENUM('RUB') DEFAULT 'RUB',
    `is_active` BOOLEAN DEFAULT TRUE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX `idx_projects_slug` (`slug`),
    INDEX `idx_projects_public_token` (`public_token`),
    INDEX `idx_projects_is_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create user_project_roles table (pivot)
CREATE TABLE IF NOT EXISTS `user_project_roles` (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `user_id` INT UNSIGNED NOT NULL,
    `project_id` INT UNSIGNED NOT NULL,
    `role` ENUM('admin', 'manager', 'client') NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX `idx_user_project` (`user_id`, `project_id`),
    INDEX `idx_project_roles_project_id` (`project_id`),
    INDEX `idx_project_roles_role` (`role`),
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`project_id`) REFERENCES `projects`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create yandex_counters table
CREATE TABLE IF NOT EXISTS `yandex_counters` (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `project_id` INT UNSIGNED NOT NULL,
    `counter_id` BIGINT NOT NULL,
    `name` VARCHAR(255),
    `is_primary` BOOLEAN DEFAULT FALSE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX `idx_counter_project_counter` (`project_id`, `counter_id`),
    INDEX `idx_counters_project_id` (`project_id`),
    FOREIGN KEY (`project_id`) REFERENCES `projects`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create yandex_direct_accounts table
CREATE TABLE IF NOT EXISTS `yandex_direct_accounts` (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `project_id` INT UNSIGNED NOT NULL,
    `client_login` VARCHAR(255) NOT NULL,
    `account_name` VARCHAR(255),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX `idx_direct_project_id` (`project_id`),
    INDEX `idx_direct_client_login` (`client_login`),
    FOREIGN KEY (`project_id`) REFERENCES `projects`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create goals table
CREATE TABLE IF NOT EXISTS `goals` (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `counter_id` INT UNSIGNED NOT NULL,
    `goal_id` BIGINT NOT NULL,
    `name` VARCHAR(255),
    `is_conversion` BOOLEAN DEFAULT FALSE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX `idx_goal_counter_goal` (`counter_id`, `goal_id`),
    INDEX `idx_goals_counter_id` (`counter_id`),
    FOREIGN KEY (`counter_id`) REFERENCES `yandex_counters`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create metrics_monthly table
CREATE TABLE IF NOT EXISTS `metrics_monthly` (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `project_id` INT UNSIGNED NOT NULL,
    `year` INT NOT NULL,
    `month` INT NOT NULL,
    `visits` INT NOT NULL DEFAULT 0,
    `users` INT NOT NULL DEFAULT 0,
    `bounce_rate` DECIMAL(5,2),
    `avg_session_duration_sec` INT NOT NULL DEFAULT 0,
    `conversions` INT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE INDEX `idx_metrics_project_month` (`project_id`, `year`, `month`),
    INDEX `idx_metrics_monthly_project_id` (`project_id`),
    FOREIGN KEY (`project_id`) REFERENCES `projects`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create metrics_age_monthly table
CREATE TABLE IF NOT EXISTS `metrics_age_monthly` (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `project_id` INT UNSIGNED NOT NULL,
    `year` INT NOT NULL,
    `month` INT NOT NULL,
    `age_group` VARCHAR(50) NOT NULL,
    `visits` INT NOT NULL DEFAULT 0,
    `users` INT NOT NULL DEFAULT 0,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE INDEX `idx_metrics_age_project_month_age` (`project_id`, `year`, `month`, `age_group`),
    INDEX `idx_metrics_age_monthly_project_id` (`project_id`),
    FOREIGN KEY (`project_id`) REFERENCES `projects`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create direct_campaigns table
CREATE TABLE IF NOT EXISTS `direct_campaigns` (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `project_id` INT UNSIGNED NOT NULL,
    `campaign_id` BIGINT NOT NULL,
    `name` VARCHAR(255),
    `status` VARCHAR(50),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX `idx_campaign_project_campaign` (`project_id`, `campaign_id`),
    INDEX `idx_direct_campaigns_project_id` (`project_id`),
    FOREIGN KEY (`project_id`) REFERENCES `projects`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create direct_campaign_monthly table
CREATE TABLE IF NOT EXISTS `direct_campaign_monthly` (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `campaign_id` INT UNSIGNED NOT NULL,
    `year` INT NOT NULL,
    `month` INT NOT NULL,
    `impressions` INT NOT NULL DEFAULT 0,
    `clicks` INT NOT NULL DEFAULT 0,
    `ctr` DECIMAL(5,2),
    `cpc` DECIMAL(10,2),
    `conversions` INT,
    `cpa` DECIMAL(10,2),
    `cost` DECIMAL(10,2) NOT NULL DEFAULT 0,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE INDEX `idx_campaign_month_campaign_month` (`campaign_id`, `year`, `month`),
    INDEX `idx_direct_campaign_monthly_campaign_id` (`campaign_id`),
    FOREIGN KEY (`campaign_id`) REFERENCES `direct_campaigns`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create direct_totals_monthly table
CREATE TABLE IF NOT EXISTS `direct_totals_monthly` (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `project_id` INT UNSIGNED NOT NULL,
    `year` INT NOT NULL,
    `month` INT NOT NULL,
    `impressions` INT NOT NULL DEFAULT 0,
    `clicks` INT NOT NULL DEFAULT 0,
    `ctr` DECIMAL(5,2),
    `cpc` DECIMAL(10,2),
    `conversions` INT,
    `cpa` DECIMAL(10,2),
    `cost` DECIMAL(10,2) NOT NULL DEFAULT 0,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE INDEX `idx_totals_project_month` (`project_id`, `year`, `month`),
    INDEX `idx_direct_totals_monthly_project_id` (`project_id`),
    FOREIGN KEY (`project_id`) REFERENCES `projects`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create seo_queries_monthly table
CREATE TABLE IF NOT EXISTS `seo_queries_monthly` (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `project_id` INT UNSIGNED NOT NULL,
    `year` INT NOT NULL,
    `month` INT NOT NULL,
    `impressions` INT NOT NULL DEFAULT 0,
    `clicks` INT NOT NULL DEFAULT 0,
    `ctr` DECIMAL(5,2),
    `avg_position` DECIMAL(5,2),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE INDEX `idx_seo_project_month` (`project_id`, `year`, `month`),
    INDEX `idx_seo_queries_monthly_project_id` (`project_id`),
    FOREIGN KEY (`project_id`) REFERENCES `projects`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
