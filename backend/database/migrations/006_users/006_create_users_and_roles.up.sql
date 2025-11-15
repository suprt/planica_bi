USE `planica_bi`;

-- Таблица пользователей согласно ТЗ (раздел 7)
CREATE TABLE `users` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    `email_verified_at` TIMESTAMP NULL,
    
    -- Дополнительные поля для системы
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

-- Таблица ролей пользователей в проектах согласно ТЗ
CREATE TABLE `user_project_roles` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `project_id` BIGINT UNSIGNED NOT NULL,
    
    -- Роли согласно ТЗ: admin, manager, client
    `role` ENUM('admin', 'manager', 'client') NOT NULL,
    
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    PRIMARY KEY (`id`),
    
    -- Уникальность как в ТЗ: один пользователь - одна роль в проекте
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

-- Аудит создания таблиц пользователей
INSERT INTO `schema_audit` (`migration_version`, `change_type`, `object_type`, `object_name`) 
VALUES 
('006_create_users_and_roles', 'CREATE', 'TABLE', 'users'),
('006_create_users_and_roles', 'CREATE', 'TABLE', 'user_project_roles');