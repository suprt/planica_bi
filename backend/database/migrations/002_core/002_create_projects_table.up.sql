USE `reports`;

CREATE TABLE `projects` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `uuid` CHAR(36) NOT NULL DEFAULT (UUID()),
    `name` VARCHAR(255) NOT NULL,
    `slug` VARCHAR(100) NOT NULL,
    `description` TEXT NULL,
    
    -- Настройки проекта согласно ТЗ
    `timezone` VARCHAR(50) NOT NULL DEFAULT 'Europe/Moscow',
    `currency` ENUM('RUB', 'USD', 'EUR') NOT NULL DEFAULT 'RUB',
    `language` ENUM('ru', 'en') NOT NULL DEFAULT 'ru',
    
    -- Статус проекта
    `status` ENUM('active', 'paused', 'archived') NOT NULL DEFAULT 'active',
    `is_public` TINYINT(1) NOT NULL DEFAULT 0,
    
    -- Даты проекта
    `start_date` DATE NULL,
    `end_date` DATE NULL,
    
    -- Метаданные
    `created_by` BIGINT UNSIGNED NULL,
    `updated_by` BIGINT UNSIGNED NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uq_uuid` (`uuid` ASC),
    UNIQUE INDEX `uq_slug` (`slug` ASC),
    INDEX `idx_status` (`status` ASC),
    INDEX `idx_public` (`is_public` ASC),
    INDEX `idx_created_at` (`created_at` ASC),
    INDEX `idx_dates` (`start_date` ASC, `end_date` ASC),
    INDEX `idx_deleted` (`deleted_at` ASC),
    
    CONSTRAINT `chk_dates` CHECK (`end_date` IS NULL OR `start_date` IS NULL OR `end_date` > `start_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8;

-- Аудит создания таблицы
INSERT INTO `schema_audit` (`migration_version`, `change_type`, `object_type`, `object_name`) 
VALUES ('002_create_projects_table', 'CREATE', 'TABLE', 'projects');