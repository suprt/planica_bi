-- Инициализация базы данных Planica BI
SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION';

-- Создание базы данных
CREATE DATABASE IF NOT EXISTS `planica_bi` 
DEFAULT CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

USE `planica_bi`;

-- Таблица для отслеживания миграций
CREATE TABLE IF NOT EXISTS `schema_migrations` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `version` VARCHAR(50) NOT NULL,
    `description` TEXT NOT NULL,
    `checksum` VARCHAR(64) NOT NULL,
    `applied_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `applied_by` VARCHAR(255) NOT NULL DEFAULT USER(),
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
    `changed_by` VARCHAR(255) NOT NULL DEFAULT USER(),
    PRIMARY KEY (`id`),
    INDEX `idx_change_type` (`change_type` ASC),
    INDEX `idx_object_type` (`object_type` ASC),
    INDEX `idx_changed_at` (`changed_at` ASC),
    INDEX `idx_migration` (`migration_version` ASC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Запись в аудит о создании базы
INSERT INTO `schema_audit` (`change_type`, `object_type`, `object_name`, `sql_statement`) 
VALUES ('CREATE', 'DATABASE', 'planica_bi', 'Database initialization');

SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;