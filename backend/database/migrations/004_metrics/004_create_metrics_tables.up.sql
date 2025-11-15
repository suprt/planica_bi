USE `reports`;

-- Основные метрики Метрики (помесячно) согласно ТЗ
CREATE TABLE `metrics_monthly` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL,
    `year` SMALLINT UNSIGNED NOT NULL,
    `month` TINYINT UNSIGNED NOT NULL,
    
    -- Основные метрики по ТЗ
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

-- Метрики по возрастным группам согласно ТЗ
CREATE TABLE `metrics_age_monthly` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL,
    `year` SMALLINT UNSIGNED NOT NULL,
    `month` TINYINT UNSIGNED NOT NULL,
    `age_group` ENUM('18-24', '25-34', '35-44', '45-54', '55+', 'unknown') NOT NULL,
    
    -- Метрики по ТЗ
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

-- Метрики Директа по кампаниям согласно ТЗ
CREATE TABLE `direct_campaign_monthly` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL,
    `direct_campaign_id` BIGINT UNSIGNED NOT NULL,
    `year` SMALLINT UNSIGNED NOT NULL,
    `month` TINYINT UNSIGNED NOT NULL,
    
    -- Метрики Директа по ТЗ
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

-- Итоговые метрики Директа по проекту согласно ТЗ
CREATE TABLE `direct_totals_monthly` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL,
    `year` SMALLINT UNSIGNED NOT NULL,
    `month` TINYINT UNSIGNED NOT NULL,
    
    -- Итоговые метрики по ТЗ
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

-- Аудит создания таблиц метрик
INSERT INTO `schema_audit` (`migration_version`, `change_type`, `object_type`, `object_name`) 
VALUES 
('004_create_metrics_tables', 'CREATE', 'TABLE', 'metrics_monthly'),
('004_create_metrics_tables', 'CREATE', 'TABLE', 'metrics_age_monthly'),
('004_create_metrics_tables', 'CREATE', 'TABLE', 'direct_campaign_monthly'),
('004_create_metrics_tables', 'CREATE', 'TABLE', 'direct_totals_monthly');