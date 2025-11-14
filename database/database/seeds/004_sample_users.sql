USE `planica_bi`;

-- Вставка тестовых пользователей
INSERT INTO `users` (`name`, `email`, `password`, `timezone`, `language`, `is_active`) VALUES 
('Администратор системы', 'admin@planica.ru', '$2y$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Europe/Moscow', 'ru', 1),
('Менеджер проектов', 'manager@planica.ru', '$2y$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Europe/Moscow', 'ru', 1),
('Клиент Иванов', 'client@electronics-store.ru', '$2y$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Europe/Moscow', 'ru', 1),
('Клиент Петров', 'client@construction-company.ru', '$2y$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Europe/Moscow', 'ru', 1);

-- Назначение ролей
INSERT INTO `user_project_roles` (`user_id`, `project_id`, `role`) VALUES 
(1, 1, 'admin'),   -- Админ на все проекты
(1, 2, 'admin'),
(1, 3, 'admin'),
(2, 1, 'manager'), -- Менеджер на проекты
(2, 2, 'manager'),
(3, 1, 'client'),  -- Клиент только на свой проект
(4, 2, 'client');  -- Клиент только на свой проект