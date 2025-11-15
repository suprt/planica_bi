USE `reports`;

-- Установка кодировки UTF-8
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- Таблица для отслеживания миграций
CREATE TABLE IF NOT EXISTS `schema_migrations` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `version` VARCHAR(50) NOT NULL,
    `description` TEXT NOT NULL,
    `checksum` VARCHAR(64) NOT NULL,
    `applied_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `applied_by` VARCHAR(255) NOT NULL DEFAULT (CURRENT_USER()),
    `execution_time_ms` INT UNSIGNED,
    `status` ENUM('pending', 'success', 'failed', 'rolled_back') NOT NULL DEFAULT 'success',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uq_version` (`version` ASC),
    INDEX `idx_applied_at` (`applied_at` ASC),
    INDEX `idx_status` (`status` ASC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Таблица аудита изменений схемы
CREATE TABLE IF NOT EXISTS `schema_audit` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `migration_version` VARCHAR(50) NULL,
    `change_type` ENUM('CREATE', 'ALTER', 'DROP', 'TRUNCATE', 'OTHER') NOT NULL,
    `object_type` ENUM('DATABASE', 'TABLE', 'INDEX', 'COLUMN', 'CONSTRAINT', 'TRIGGER') NOT NULL,
    `object_name` VARCHAR(255) NOT NULL,
    `sql_statement` TEXT NULL,
    `changed_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `changed_by` VARCHAR(255) NOT NULL DEFAULT (CURRENT_USER()),
    PRIMARY KEY (`id`),
    INDEX `idx_change_type` (`change_type` ASC),
    INDEX `idx_object_type` (`object_type` ASC),
    INDEX `idx_changed_at` (`changed_at` ASC),
    INDEX `idx_migration` (`migration_version` ASC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Таблица проектов (соответствует GORM модели)
CREATE TABLE `projects` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` LONGTEXT NOT NULL,
    `slug` VARCHAR(191) NOT NULL,
    `public_token` VARCHAR(64) NULL,
    `timezone` VARCHAR(191) NULL DEFAULT 'Europe/Moscow',
    `currency` ENUM('RUB') NULL DEFAULT 'RUB',
    `is_active` TINYINT(1) NULL DEFAULT 1,
    `created_at` DATETIME(3) NULL,
    `updated_at` DATETIME(3) NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uq_slug` (`slug` ASC),
    UNIQUE INDEX `uq_public_token` (`public_token` ASC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Таблица пользователей
CREATE TABLE `users` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    `email_verified_at` TIMESTAMP NULL,
    `timezone` VARCHAR(50) NOT NULL DEFAULT 'Europe/Moscow',
    `language` ENUM('ru', 'en') NOT NULL DEFAULT 'ru',
    `is_active` TINYINT(1) NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `last_login_at` TIMESTAMP NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uq_email` (`email` ASC),
    INDEX `idx_email` (`email` ASC),
    INDEX `idx_created` (`created_at` ASC),
    INDEX `idx_active` (`is_active` ASC),
    INDEX `idx_last_login` (`last_login_at` ASC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Таблица ролей пользователей в проектах
CREATE TABLE `user_project_roles` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `project_id` BIGINT UNSIGNED NOT NULL,
    `role` ENUM('admin', 'manager', 'client') NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uq_user_project` (`user_id` ASC, `project_id` ASC),
    INDEX `fk_user_roles_user` (`user_id` ASC),
    INDEX `fk_user_roles_project` (`project_id` ASC),
    INDEX `idx_role` (`role` ASC),
    INDEX `idx_created` (`created_at` ASC),
    CONSTRAINT `fk_user_roles_user`
        FOREIGN KEY (`user_id`)
        REFERENCES `users` (`id`)
        ON DELETE CASCADE
        ON UPDATE NO ACTION,
    CONSTRAINT `fk_user_roles_project`
        FOREIGN KEY (`project_id`)
        REFERENCES `projects` (`id`)
        ON DELETE CASCADE
        ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Таблица счетчиков Яндекс.Метрики
CREATE TABLE `yandex_counters` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL,
    `counter_id` BIGINT UNSIGNED NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `site_name` VARCHAR(255) NULL,
    `is_primary` TINYINT(1) NOT NULL DEFAULT 0,
    `auto_import` TINYINT(1) NOT NULL DEFAULT 1,
    `import_goals` TINYINT(1) NOT NULL DEFAULT 1,
    `status` ENUM('active', 'paused', 'error') NOT NULL DEFAULT 'active',
    `last_sync_at` TIMESTAMP NULL,
    `sync_status` TEXT NULL,
    `access_token_encrypted` TEXT NULL,
    `token_expires_at` TIMESTAMP NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uq_project_counter` (`project_id` ASC, `counter_id` ASC),
    INDEX `fk_yandex_counters_project` (`project_id` ASC),
    INDEX `idx_counter_id` (`counter_id` ASC),
    INDEX `idx_status` (`status` ASC),
    INDEX `idx_primary` (`is_primary` ASC),
    INDEX `idx_last_sync` (`last_sync_at` ASC),
    CONSTRAINT `fk_yandex_counters_project`
        FOREIGN KEY (`project_id`)
        REFERENCES `projects` (`id`)
        ON DELETE CASCADE
        ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Таблица аккаунтов Яндекс.Директа
CREATE TABLE `direct_accounts` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL,
    `client_login` VARCHAR(255) NOT NULL,
    `account_name` VARCHAR(255) NOT NULL,
    `account_id` BIGINT UNSIGNED NOT NULL,
    `auto_import_campaigns` TINYINT(1) NOT NULL DEFAULT 1,
    `import_statistics` TINYINT(1) NOT NULL DEFAULT 1,
    `status` ENUM('active', 'paused', 'error') NOT NULL DEFAULT 'active',
    `last_sync_at` TIMESTAMP NULL,
    `sync_status` TEXT NULL,
    `access_token_encrypted` TEXT NULL,
    `token_expires_at` TIMESTAMP NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uq_project_login` (`project_id` ASC, `client_login` ASC),
    INDEX `fk_direct_accounts_project` (`project_id` ASC),
    INDEX `idx_account_id` (`account_id` ASC),
    INDEX `idx_status` (`status` ASC),
    CONSTRAINT `fk_direct_accounts_project`
        FOREIGN KEY (`project_id`)
        REFERENCES `projects` (`id`)
        ON DELETE CASCADE
        ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Таблица кампаний Директа
CREATE TABLE `direct_campaigns` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `direct_account_id` BIGINT UNSIGNED NOT NULL,
    `campaign_id` BIGINT UNSIGNED NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `status` VARCHAR(50) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uq_account_campaign` (`direct_account_id` ASC, `campaign_id` ASC),
    INDEX `fk_direct_campaigns_account` (`direct_account_id` ASC),
    INDEX `idx_campaign_id` (`campaign_id` ASC),
    INDEX `idx_status` (`status` ASC),
    CONSTRAINT `fk_direct_campaigns_account`
        FOREIGN KEY (`direct_account_id`)
        REFERENCES `direct_accounts` (`id`)
        ON DELETE CASCADE
        ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Таблица целей Яндекс.Метрики
CREATE TABLE `goals` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `counter_id` BIGINT UNSIGNED NOT NULL,
    `goal_id` BIGINT UNSIGNED NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `description` TEXT NULL,
    `goal_type` ENUM('url', 'action', 'step', 'email') NOT NULL DEFAULT 'url',
    `is_conversion` TINYINT(1) NOT NULL DEFAULT 0,
    `conversion_weight` DECIMAL(5,2) NOT NULL DEFAULT 1.00,
    `is_active` TINYINT(1) NOT NULL DEFAULT 1,
    `target_url` VARCHAR(500) NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uq_counter_goal` (`counter_id` ASC, `goal_id` ASC),
    INDEX `fk_goals_counter` (`counter_id` ASC),
    INDEX `idx_goal_id` (`goal_id` ASC),
    INDEX `idx_conversion` (`is_conversion` ASC),
    INDEX `idx_active` (`is_active` ASC),
    CONSTRAINT `fk_goals_counter`
        FOREIGN KEY (`counter_id`)
        REFERENCES `yandex_counters` (`id`)
        ON DELETE CASCADE
        ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Основные метрики Метрики (помесячно)
CREATE TABLE `metrics_monthly` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL,
    `year` SMALLINT UNSIGNED NOT NULL,
    `month` TINYINT UNSIGNED NOT NULL,
    `visits` INT UNSIGNED NOT NULL DEFAULT 0,
    `users` INT UNSIGNED NOT NULL DEFAULT 0,
    `bounce_rate` DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    `avg_session_duration_sec` INT UNSIGNED NOT NULL DEFAULT 0,
    `conversions` INT UNSIGNED NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uq_project_month` (`project_id` ASC, `year` ASC, `month` ASC),
    INDEX `fk_metrics_monthly_project` (`project_id` ASC),
    INDEX `idx_project_period` (`project_id` ASC, `year` ASC, `month` ASC),
    INDEX `idx_period` (`year` ASC, `month` ASC),
    CONSTRAINT `fk_metrics_monthly_project`
        FOREIGN KEY (`project_id`)
        REFERENCES `projects` (`id`)
        ON DELETE CASCADE
        ON UPDATE NO ACTION,
    CONSTRAINT `chk_month` CHECK (`month` >= 1 AND `month` <= 12),
    CONSTRAINT `chk_year` CHECK (`year` >= 2020 AND `year` <= 2030)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=COMPRESSED;

-- Метрики по возрастным группам
CREATE TABLE `metrics_age_monthly` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL,
    `year` SMALLINT UNSIGNED NOT NULL,
    `month` TINYINT UNSIGNED NOT NULL,
    `age_group` ENUM('18-24', '25-34', '35-44', '45-54', '55+', 'unknown') NOT NULL,
    `visits` INT UNSIGNED NOT NULL DEFAULT 0,
    `users` INT UNSIGNED NOT NULL DEFAULT 0,
    `bounce_rate` DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    `avg_session_duration_sec` INT UNSIGNED NOT NULL DEFAULT 0,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uq_project_age_month` (`project_id` ASC, `year` ASC, `month` ASC, `age_group` ASC),
    INDEX `fk_metrics_age_project` (`project_id` ASC),
    INDEX `idx_project_period` (`project_id` ASC, `year` ASC, `month` ASC),
    INDEX `idx_age_group` (`age_group` ASC),
    CONSTRAINT `fk_metrics_age_project`
        FOREIGN KEY (`project_id`)
        REFERENCES `projects` (`id`)
        ON DELETE CASCADE
        ON UPDATE NO ACTION,
    CONSTRAINT `chk_age_month` CHECK (`month` >= 1 AND `month` <= 12),
    CONSTRAINT `chk_age_year` CHECK (`year` >= 2020 AND `year` <= 2030)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=COMPRESSED;

-- Метрики Директа по кампаниям
CREATE TABLE `direct_campaign_monthly` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL,
    `direct_campaign_id` BIGINT UNSIGNED NOT NULL,
    `year` SMALLINT UNSIGNED NOT NULL,
    `month` TINYINT UNSIGNED NOT NULL,
    `impressions` INT UNSIGNED NOT NULL DEFAULT 0,
    `clicks` INT UNSIGNED NOT NULL DEFAULT 0,
    `ctr_pct` DECIMAL(6,2) NOT NULL DEFAULT 0.00,
    `cpc` DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    `conversions` INT UNSIGNED NULL,
    `cpa` DECIMAL(12,2) NULL,
    `cost` DECIMAL(14,2) NOT NULL DEFAULT 0.00,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uq_campaign_month` (`direct_campaign_id` ASC, `year` ASC, `month` ASC),
    INDEX `fk_direct_campaign_monthly_project` (`project_id` ASC),
    INDEX `fk_direct_campaign_monthly_campaign` (`direct_campaign_id` ASC),
    INDEX `idx_project_period` (`project_id` ASC, `year` ASC, `month` ASC),
    CONSTRAINT `fk_direct_campaign_monthly_project`
        FOREIGN KEY (`project_id`)
        REFERENCES `projects` (`id`)
        ON DELETE CASCADE
        ON UPDATE NO ACTION,
    CONSTRAINT `fk_direct_campaign_monthly_campaign`
        FOREIGN KEY (`direct_campaign_id`)
        REFERENCES `direct_campaigns` (`id`)
        ON DELETE CASCADE
        ON UPDATE NO ACTION,
    CONSTRAINT `chk_direct_month` CHECK (`month` >= 1 AND `month` <= 12),
    CONSTRAINT `chk_direct_year` CHECK (`year` >= 2020 AND `year` <= 2030)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=COMPRESSED;

-- Итоговые метрики Директа по проекту
CREATE TABLE `direct_totals_monthly` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL,
    `year` SMALLINT UNSIGNED NOT NULL,
    `month` TINYINT UNSIGNED NOT NULL,
    `impressions` INT UNSIGNED NOT NULL DEFAULT 0,
    `clicks` INT UNSIGNED NOT NULL DEFAULT 0,
    `ctr_pct` DECIMAL(6,2) NOT NULL DEFAULT 0.00,
    `cpc` DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    `conversions` INT UNSIGNED NULL,
    `cpa` DECIMAL(12,2) NULL,
    `cost` DECIMAL(14,2) NOT NULL DEFAULT 0.00,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uq_project_month` (`project_id` ASC, `year` ASC, `month` ASC),
    INDEX `fk_direct_totals_project` (`project_id` ASC),
    INDEX `idx_project_period` (`project_id` ASC, `year` ASC, `month` ASC),
    CONSTRAINT `fk_direct_totals_project`
        FOREIGN KEY (`project_id`)
        REFERENCES `projects` (`id`)
        ON DELETE CASCADE
        ON UPDATE NO ACTION,
    CONSTRAINT `chk_totals_month` CHECK (`month` >= 1 AND `month` <= 12),
    CONSTRAINT `chk_totals_year` CHECK (`year` >= 2020 AND `year` <= 2030)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=COMPRESSED;

-- SEO данные
CREATE TABLE `seo_queries_monthly` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL,
    `year` SMALLINT UNSIGNED NOT NULL,
    `month` TINYINT UNSIGNED NOT NULL,
    `query` VARCHAR(500) NOT NULL,
    `position` DECIMAL(5,1) NOT NULL,
    `url` VARCHAR(500) NULL,
    `impressions` INT UNSIGNED NULL DEFAULT 0,
    `clicks` INT UNSIGNED NULL DEFAULT 0,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `fk_seo_queries_project` (`project_id` ASC),
    INDEX `idx_project_period` (`project_id` ASC, `year` ASC, `month` ASC),
    INDEX `idx_query` (`query`(255) ASC),
    INDEX `idx_position` (`position` ASC),
    INDEX `idx_period` (`year` ASC, `month` ASC),
    UNIQUE INDEX `uq_project_query_month` (`project_id` ASC, `query`(200) ASC, `year` ASC, `month` ASC),
    CONSTRAINT `fk_seo_queries_project`
        FOREIGN KEY (`project_id`)
        REFERENCES `projects` (`id`)
        ON DELETE CASCADE
        ON UPDATE NO ACTION,
    CONSTRAINT `chk_seo_month` CHECK (`month` >= 1 AND `month` <= 12),
    CONSTRAINT `chk_seo_year` CHECK (`year` >= 2020 AND `year` <= 2030),
    CONSTRAINT `chk_position` CHECK (`position` >= 0 AND `position` <= 100)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=COMPRESSED;

-- Запись в аудит о создании таблиц
INSERT INTO `schema_audit` (`change_type`, `object_type`, `object_name`, `sql_statement`) 
VALUES 
('CREATE', 'TABLE', 'schema_migrations', 'System migrations table'),
('CREATE', 'TABLE', 'schema_audit', 'System audit table'),
('CREATE', 'TABLE', 'projects', 'Projects table'),
('CREATE', 'TABLE', 'users', 'Users table'),
('CREATE', 'TABLE', 'user_project_roles', 'User project roles table'),
('CREATE', 'TABLE', 'yandex_counters', 'Yandex Metrica counters table'),
('CREATE', 'TABLE', 'direct_accounts', 'Yandex Direct accounts table'),
('CREATE', 'TABLE', 'direct_campaigns', 'Yandex Direct campaigns table'),
('CREATE', 'TABLE', 'goals', 'Yandex Metrica goals table'),
('CREATE', 'TABLE', 'metrics_monthly', 'Monthly metrics table'),
('CREATE', 'TABLE', 'metrics_age_monthly', 'Age group metrics table'),
('CREATE', 'TABLE', 'direct_campaign_monthly', 'Direct campaign metrics table'),
('CREATE', 'TABLE', 'direct_totals_monthly', 'Direct totals metrics table'),
('CREATE', 'TABLE', 'seo_queries_monthly', 'SEO queries table');