USE `planica_bi`;

-- Таблица счетчиков Яндекс.Метрики
CREATE TABLE `yandex_counters` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL,
    `counter_id` BIGINT UNSIGNED NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `site_name` VARCHAR(255) NULL,
    
    -- Настройки интеграции
    `is_primary` TINYINT(1) NOT NULL DEFAULT 0,
    `auto_import` TINYINT(1) NOT NULL DEFAULT 1,
    `import_goals` TINYINT(1) NOT NULL DEFAULT 1,
    
    -- Статус интеграции
    `status` ENUM('active', 'paused', 'error') NOT NULL DEFAULT 'active',
    `last_sync_at` TIMESTAMP NULL,
    `sync_status` TEXT NULL,
    
    -- Токены (в продакшене должны быть в vault)
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
    
    -- Настройки интеграции
    `auto_import_campaigns` TINYINT(1) NOT NULL DEFAULT 1,
    `import_statistics` TINYINT(1) NOT NULL DEFAULT 1,
    
    -- Статус
    `status` ENUM('active', 'paused', 'error') NOT NULL DEFAULT 'active',
    `last_sync_at` TIMESTAMP NULL,
    `sync_status` TEXT NULL,
    
    -- Токены
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
    
    -- Тип цели
    `goal_type` ENUM('url', 'action', 'step', 'email') NOT NULL DEFAULT 'url',
    `is_conversion` TINYINT(1) NOT NULL DEFAULT 0,
    `conversion_weight` DECIMAL(5,2) NOT NULL DEFAULT 1.00,
    
    -- Настройки
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

-- Аудит создания таблиц аналитики
INSERT INTO `schema_audit` (`migration_version`, `change_type`, `object_type`, `object_name`) 
VALUES 
('003_create_yandex_tables', 'CREATE', 'TABLE', 'yandex_counters'),
('003_create_yandex_tables', 'CREATE', 'TABLE', 'direct_accounts'),
('003_create_yandex_tables', 'CREATE', 'TABLE', 'direct_campaigns'),
('003_create_yandex_tables', 'CREATE', 'TABLE', 'goals');