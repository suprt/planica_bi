USE `planica_bi`;

-- SEO данные согласно ТЗ (разделы 1.2, 4.1 и 8)
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
    
    -- Уникальность для идемпотентности как в ТЗ
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

-- Аудит создания SEO таблиц
INSERT INTO `schema_audit` (`migration_version`, `change_type`, `object_type`, `object_name`) 
VALUES 
('005_create_seo_tables', 'CREATE', 'TABLE', 'seo_queries_monthly');